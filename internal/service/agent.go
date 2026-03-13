package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"nunu/internal/agent"
	"nunu/internal/agent/tools"
	"nunu/internal/model"
	"nunu/internal/repository"
	"nunu/pkg/log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AgentService interface {
	HandleQuestion(ctx context.Context, channelID, userID, threadTS, question string) (string, error)
	HandleFeedback(ctx context.Context, channelID, threadTS, feedback string) error
}

type agentService struct {
	agent    *agent.ReactAgent
	memorySvc  MemoryService
	convRepo   repository.ConversationRepository
	conf       *viper.Viper
	logger     *log.Logger
}

func NewAgentService(
	conf *viper.Viper,
	logger *log.Logger,
	memorySvc MemoryService,
	knowledgeSvc KnowledgeService,
	querySvc QueryService,
	csvSvc CSVService,
	convRepo repository.ConversationRepository,
) AgentService {
	// Build LLM client
	llmClient := agent.NewLLMClient(conf, logger.Logger)

	// Build tool registry
	registry := agent.NewToolRegistry()
	registry.Register(tools.NewMemorySearchTool(memorySvc))
	registry.Register(tools.NewKnowledgeSearchTool(knowledgeSvc))
	registry.Register(tools.NewQueryPostgresTool(querySvc, conf.GetInt("agent.max_query_rows")))
	registry.Register(tools.NewCSVDownloadTool(csvSvc))

	maxSteps := conf.GetInt("agent.max_steps")
	reactAgent := agent.NewReactAgent(llmClient, registry, logger.Logger, maxSteps)

	return &agentService{
		agent:    reactAgent,
		memorySvc:  memorySvc,
		convRepo:   convRepo,
		conf:       conf,
		logger:     logger,
	}
}

func (s *agentService) HandleQuestion(ctx context.Context, channelID, userID, threadTS, question string) (string, error) {
	conversationID := uuid.New().String()
	start := time.Now()

	// Save conversation record
	conv := &model.Conversation{
		ChannelID: channelID,
		UserID:    userID,
		ThreadTS:  threadTS,
		Question:  question,
	}
	if err := s.convRepo.Create(ctx, conv); err != nil {
		s.logger.Warn("failed to save conversation", zap.Error(err))
	}

	// Run the ReAct agent
	req := agent.AgentRequest{
		ConversationID: conversationID,
		UserID:         userID,
		ChannelID:      channelID,
		Question:       question,
	}

	resp, err := s.agent.Run(ctx, req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		// Update conversation with failure
		conv.Success = false
		conv.Duration = duration
		conv.Answer = fmt.Sprintf("오류가 발생했습니다: %v", err)
		s.convRepo.Update(ctx, conv)
		return "", err
	}

	// Update conversation with success
	toolsJSON, _ := json.Marshal(resp.ToolsUsed)
	conv.Answer = resp.Answer
	conv.SQLUsed = resp.SQLUsed
	conv.ToolsUsed = string(toolsJSON)
	conv.Success = true
	conv.Duration = duration
	if err := s.convRepo.Update(ctx, conv); err != nil {
		s.logger.Warn("failed to update conversation", zap.Error(err))
	}

	// Save successful interaction to memory (feedback loop)
	go func() {
		bgCtx := context.Background()
		summary := s.generateSummary(resp)
		if err := s.memorySvc.Save(bgCtx, conversationID, question, summary, resp.SQLUsed, resp.ToolsUsed); err != nil {
			s.logger.Warn("failed to save memory", zap.Error(err))
		}
	}()

	return resp.Answer, nil
}

func (s *agentService) HandleFeedback(ctx context.Context, channelID, threadTS, feedback string) error {
	conv, err := s.convRepo.GetByThreadTS(ctx, channelID, threadTS)
	if err != nil {
		return fmt.Errorf("conversation not found: %w", err)
	}
	conv.Feedback = feedback
	return s.convRepo.Update(ctx, conv)
}

func (s *agentService) generateSummary(resp *agent.AgentResponse) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("답변: %s", truncate(resp.Answer, 200)))
	if resp.SQLUsed != "" {
		parts = append(parts, fmt.Sprintf("SQL: %s", resp.SQLUsed))
	}
	if len(resp.ToolsUsed) > 0 {
		parts = append(parts, fmt.Sprintf("사용 도구: %s", strings.Join(resp.ToolsUsed, ", ")))
	}
	return strings.Join(parts, "\n")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
