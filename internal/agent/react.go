package agent

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// AgentRequest is the input to the ReAct agent.
type AgentRequest struct {
	ConversationID string
	UserID         string
	ChannelID      string
	Question       string
}

// AgentResponse is the output from the ReAct agent.
type AgentResponse struct {
	Answer    string
	SQLUsed   string
	ToolsUsed []string
	Steps     []Step
}

// Step records a single ReAct iteration.
type Step struct {
	ToolName  string
	ToolArgs  string
	Result    string
}

// ReactAgent implements the ReAct (Reason + Act) loop.
type ReactAgent struct {
	llm      LLMClient
	tools    *ToolRegistry
	logger   *zap.Logger
	maxSteps int
}

func NewReactAgent(llm LLMClient, tools *ToolRegistry, logger *zap.Logger, maxSteps int) *ReactAgent {
	if maxSteps <= 0 {
		maxSteps = 15
	}
	return &ReactAgent{
		llm:      llm,
		tools:    tools,
		logger:   logger,
		maxSteps: maxSteps,
	}
}

// Run executes the ReAct loop.
func (a *ReactAgent) Run(ctx context.Context, req AgentRequest) (*AgentResponse, error) {
	messages := []Message{
		{Role: "system", Content: SystemPrompt},
		{Role: "user", Content: req.Question},
	}

	toolList := a.tools.List()
	response := &AgentResponse{}
	var lastSQL string

	for step := 0; step < a.maxSteps; step++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		a.logger.Info("react step",
			zap.Int("step", step+1),
			zap.String("conversation_id", req.ConversationID),
		)

		llmResp, err := a.llm.Chat(ctx, messages, toolList)
		if err != nil {
			return nil, fmt.Errorf("llm chat error at step %d: %w", step+1, err)
		}

		// No tool calls → final answer
		if len(llmResp.ToolCalls) == 0 {
			response.Answer = llmResp.Content
			response.SQLUsed = lastSQL
			return response, nil
		}

		// Add assistant message with tool calls
		assistantMsg := Message{
			Role:      "assistant",
			Content:   llmResp.Content,
			ToolCalls: llmResp.ToolCalls,
		}
		messages = append(messages, assistantMsg)

		// Execute each tool call
		for _, tc := range llmResp.ToolCalls {
			a.logger.Info("executing tool",
				zap.String("tool", tc.Name),
				zap.String("args", tc.Arguments),
			)

			tool, err := a.tools.Get(tc.Name)
			if err != nil {
				toolResult := fmt.Sprintf("Error: tool %q not found", tc.Name)
				messages = append(messages, Message{
					Role:       "tool",
					Content:    toolResult,
					ToolCallID: tc.ID,
				})
				continue
			}

			result, err := tool.Execute(ctx, tc.Arguments)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			// Track tool usage
			response.ToolsUsed = append(response.ToolsUsed, tc.Name)
			response.Steps = append(response.Steps, Step{
				ToolName: tc.Name,
				ToolArgs: tc.Arguments,
				Result:   result,
			})

			// Track SQL for memory
			if tc.Name == "query_postgres" {
				lastSQL = tc.Arguments
			}

			messages = append(messages, Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
	}

	// Max steps exceeded — ask LLM for a final summary
	messages = append(messages, Message{
		Role:    "user",
		Content: "최대 단계에 도달했습니다. 지금까지의 정보를 바탕으로 최선의 답변을 제공해주세요.",
	})

	llmResp, err := a.llm.Chat(ctx, messages, nil) // no tools to force text response
	if err != nil {
		return nil, fmt.Errorf("final summary error: %w", err)
	}

	response.Answer = llmResp.Content
	response.SQLUsed = lastSQL
	return response, nil
}
