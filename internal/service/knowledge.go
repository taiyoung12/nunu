package service

import (
	"context"

	"nunu/internal/agent/tools"
	"nunu/internal/repository"
	"nunu/pkg/embedding"
	"nunu/pkg/log"

	"go.uber.org/zap"
)

type KnowledgeService interface {
	tools.KnowledgeSearcher
}

type knowledgeService struct {
	repo      repository.KnowledgeRepository
	embedding *embedding.Client
	logger    *log.Logger
}

func NewKnowledgeService(
	repo repository.KnowledgeRepository,
	embClient *embedding.Client,
	logger *log.Logger,
) KnowledgeService {
	return &knowledgeService{
		repo:      repo,
		embedding: embClient,
		logger:    logger,
	}
}

func (s *knowledgeService) Search(ctx context.Context, query string, category string, limit int) ([]tools.KnowledgeResult, error) {
	// Try semantic search first
	vec, err := s.embedding.Embed(ctx, query)
	if err != nil {
		s.logger.Warn("embedding failed, falling back to keyword search", zap.Error(err))
		return s.keywordSearch(ctx, query, category, limit)
	}

	knowledges, err := s.repo.SearchByEmbedding(ctx, vec, category, limit)
	if err != nil {
		return nil, err
	}

	// Fall back to keyword search if no results
	if len(knowledges) == 0 {
		return s.keywordSearch(ctx, query, category, limit)
	}

	results := make([]tools.KnowledgeResult, 0, len(knowledges))
	for _, k := range knowledges {
		results = append(results, tools.KnowledgeResult{
			Category: k.Category,
			Title:    k.Title,
			Content:  k.Content,
		})
	}

	return results, nil
}

func (s *knowledgeService) keywordSearch(ctx context.Context, query string, category string, limit int) ([]tools.KnowledgeResult, error) {
	knowledges, err := s.repo.SearchByKeyword(ctx, query, category, limit)
	if err != nil {
		return nil, err
	}

	results := make([]tools.KnowledgeResult, 0, len(knowledges))
	for _, k := range knowledges {
		results = append(results, tools.KnowledgeResult{
			Category: k.Category,
			Title:    k.Title,
			Content:  k.Content,
		})
	}

	return results, nil
}
