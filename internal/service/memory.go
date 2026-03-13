package service

import (
	"context"
	"encoding/json"

	"nunu/internal/agent/tools"
	"nunu/internal/model"
	"nunu/internal/repository"
	"nunu/pkg/embedding"
	"nunu/pkg/log"

	"github.com/pgvector/pgvector-go"
	"go.uber.org/zap"
)

type MemoryService interface {
	tools.MemorySearcher
	Save(ctx context.Context, conversationID, question, summary, sqlUsed string, toolsUsed []string) error
}

type memoryService struct {
	repo      repository.MemoryRepository
	embedding *embedding.Client
	logger    *log.Logger
}

func NewMemoryService(
	repo repository.MemoryRepository,
	embClient *embedding.Client,
	logger *log.Logger,
) MemoryService {
	return &memoryService{
		repo:      repo,
		embedding: embClient,
		logger:    logger,
	}
}

func (s *memoryService) SearchSimilar(ctx context.Context, query string, limit int) ([]tools.MemoryResult, error) {
	vec, err := s.embedding.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	memories, err := s.repo.SearchByEmbedding(ctx, vec, limit)
	if err != nil {
		return nil, err
	}

	results := make([]tools.MemoryResult, 0, len(memories))
	for _, m := range memories {
		results = append(results, tools.MemoryResult{
			Question: m.UserQuestion,
			Summary:  m.Summary,
			SQLUsed:  m.SQLUsed,
		})
	}

	return results, nil
}

func (s *memoryService) Save(ctx context.Context, conversationID, question, summary, sqlUsed string, toolsUsed []string) error {
	// Generate embedding for the question + summary
	text := question
	if summary != "" {
		text += "\n" + summary
	}

	vec, err := s.embedding.Embed(ctx, text)
	if err != nil {
		s.logger.Warn("failed to generate embedding for memory", zap.Error(err))
		return err
	}

	toolsJSON, _ := json.Marshal(toolsUsed)

	memory := &model.Memory{
		ConversationID: conversationID,
		UserQuestion:   question,
		Summary:        summary,
		SQLUsed:        sqlUsed,
		ToolsUsed:      string(toolsJSON),
		Embedding:      pgvector.NewVector(vec),
	}

	return s.repo.Create(ctx, memory)
}
