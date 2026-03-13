package service

import (
	"context"
	"fmt"
	"strings"

	"nunu/internal/agent/tools"
	"nunu/internal/repository"
	"nunu/pkg/log"
)

type QueryService interface {
	tools.QueryExecutor
}

type queryService struct {
	engine  repository.QueryEngine
	maxRows int
	logger  *log.Logger
}

func NewQueryService(
	engine repository.QueryEngine,
	logger *log.Logger,
	maxRows int,
) QueryService {
	if maxRows <= 0 {
		maxRows = 100
	}
	return &queryService{
		engine:  engine,
		maxRows: maxRows,
		logger:  logger,
	}
}

func (s *queryService) ExecuteReadOnly(ctx context.Context, query string, maxRows int) ([]map[string]interface{}, error) {
	// Additional validation at service layer
	normalized := strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(normalized, "SELECT") && !strings.HasPrefix(normalized, "WITH") {
		return nil, fmt.Errorf("only SELECT queries are allowed")
	}

	if maxRows <= 0 || maxRows > s.maxRows {
		maxRows = s.maxRows
	}

	// Add LIMIT if not present
	if !strings.Contains(normalized, "LIMIT") {
		query = fmt.Sprintf("%s LIMIT %d", strings.TrimRight(query, "; \n\t"), maxRows)
	}

	rows, err := s.engine.Execute(ctx, query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
