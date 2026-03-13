package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Message represents a chat message.
type Message struct {
	Role       string     `json:"role"` // "system", "user", "assistant", "tool"
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// ToolCall represents a tool invocation by the LLM.
type ToolCall struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Arguments string `json:"arguments"`
}

// LLMResponse is the response from the LLM.
type LLMResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// LLMClient abstracts LLM providers.
type LLMClient interface {
	Chat(ctx context.Context, messages []Message, tools []Tool) (*LLMResponse, error)
}

// NewLLMClient creates an LLM client based on config.
func NewLLMClient(conf *viper.Viper, logger *zap.Logger) LLMClient {
	provider := conf.GetString("llm.provider")
	switch provider {
	case "anthropic":
		return NewAnthropicClient(conf, logger)
	default:
		return NewOpenAIClient(conf, logger)
	}
}

// --- OpenAI Client ---

type OpenAIClient struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float32
	logger      *zap.Logger
}

func NewOpenAIClient(conf *viper.Viper, logger *zap.Logger) *OpenAIClient {
	apiKey := conf.GetString("llm.openai.api_key")
	return &OpenAIClient{
		client:      openai.NewClient(apiKey),
		model:       conf.GetString("llm.openai.model"),
		maxTokens:   conf.GetInt("llm.max_tokens"),
		temperature: float32(conf.GetFloat64("llm.temperature")),
		logger:      logger,
	}
}

func (c *OpenAIClient) Chat(ctx context.Context, messages []Message, tools []Tool) (*LLMResponse, error) {
	oaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, msg := range messages {
		oaiMsg := openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
		if msg.ToolCallID != "" {
			oaiMsg.ToolCallID = msg.ToolCallID
		}
		if len(msg.ToolCalls) > 0 {
			oaiMsg.ToolCalls = make([]openai.ToolCall, len(msg.ToolCalls))
			for i, tc := range msg.ToolCalls {
				oaiMsg.ToolCalls[i] = openai.ToolCall{
					ID:   tc.ID,
					Type: openai.ToolTypeFunction,
					Function: openai.FunctionCall{
						Name:      tc.Name,
						Arguments: tc.Arguments,
					},
				}
			}
		}
		oaiMessages = append(oaiMessages, oaiMsg)
	}

	var oaiTools []openai.Tool
	for _, t := range tools {
		paramBytes, _ := json.Marshal(t.Parameters())
		oaiTools = append(oaiTools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  json.RawMessage(paramBytes),
			},
		})
	}

	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    oaiMessages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}
	if len(oaiTools) > 0 {
		req.Tools = oaiTools
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("openai chat error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	choice := resp.Choices[0]
	result := &LLMResponse{
		Content: choice.Message.Content,
	}

	for _, tc := range choice.Message.ToolCalls {
		result.ToolCalls = append(result.ToolCalls, ToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: tc.Function.Arguments,
		})
	}

	return result, nil
}

// --- Anthropic Client ---

type AnthropicClient struct {
	client      anthropic.Client
	model       string
	maxTokens   int
	temperature float64
	logger      *zap.Logger
}

func NewAnthropicClient(conf *viper.Viper, logger *zap.Logger) *AnthropicClient {
	apiKey := conf.GetString("llm.anthropic.api_key")
	client := anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))
	return &AnthropicClient{
		client:      client,
		model:       conf.GetString("llm.anthropic.model"),
		maxTokens:   conf.GetInt("llm.max_tokens"),
		temperature: conf.GetFloat64("llm.temperature"),
		logger:      logger,
	}
}

func (c *AnthropicClient) Chat(ctx context.Context, messages []Message, tools []Tool) (*LLMResponse, error) {
	// Separate system message
	var systemText string
	var anthropicMsgs []anthropic.MessageParam
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			systemText = msg.Content
		case "user":
			anthropicMsgs = append(anthropicMsgs, anthropic.NewUserMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				blocks := []anthropic.ContentBlockParamUnion{}
				if msg.Content != "" {
					blocks = append(blocks, anthropic.NewTextBlock(msg.Content))
				}
				for _, tc := range msg.ToolCalls {
					var input interface{}
					json.Unmarshal([]byte(tc.Arguments), &input)
					if input == nil {
						input = map[string]interface{}{}
					}
					blocks = append(blocks, anthropic.NewToolUseBlock(tc.ID, input, tc.Name))
				}
				anthropicMsgs = append(anthropicMsgs, anthropic.NewAssistantMessage(blocks...))
			} else {
				anthropicMsgs = append(anthropicMsgs, anthropic.NewAssistantMessage(
					anthropic.NewTextBlock(msg.Content),
				))
			}
		case "tool":
			anthropicMsgs = append(anthropicMsgs, anthropic.NewUserMessage(
				anthropic.NewToolResultBlock(msg.ToolCallID, msg.Content, false),
			))
		}
	}

	// Build tool definitions
	var anthropicTools []anthropic.ToolUnionParam
	for _, t := range tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        t.Name(),
				Description: anthropic.String(t.Description()),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: t.Parameters(),
				},
			},
		})
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: int64(c.maxTokens),
		Messages:  anthropicMsgs,
	}
	if systemText != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: systemText},
		}
	}
	if len(anthropicTools) > 0 {
		params.Tools = anthropicTools
	}

	resp, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic chat error: %w", err)
	}

	result := &LLMResponse{}
	for _, block := range resp.Content {
		switch block.Type {
		case "text":
			result.Content += block.Text
		case "tool_use":
			argsBytes, _ := json.Marshal(block.Input)
			result.ToolCalls = append(result.ToolCalls, ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: string(argsBytes),
			})
		}
	}

	return result, nil
}
