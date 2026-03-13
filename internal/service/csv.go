package service

import (
	"context"
	"fmt"
	"path"

	"nunu/internal/agent/tools"
	"nunu/internal/repository"
	"nunu/pkg/csvutil"
	"nunu/pkg/log"
)

type CSVService interface {
	tools.CSVGenerator
}

type csvService struct {
	engine      repository.QueryEngine
	storagePath string
	baseURL     string
	logger      *log.Logger
}

func NewCSVService(
	engine repository.QueryEngine,
	logger *log.Logger,
	storagePath string,
	baseURL string,
) CSVService {
	return &csvService{
		engine:      engine,
		storagePath: storagePath,
		baseURL:     baseURL,
		logger:      logger,
	}
}

func (s *csvService) GenerateFromQuery(ctx context.Context, query string) (string, error) {
	rows, err := s.engine.Execute(ctx, query)
	if err != nil {
		return "", fmt.Errorf("csv query execution error: %w", err)
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("no data to export")
	}

	filename, err := csvutil.WriteCSV(s.storagePath, rows)
	if err != nil {
		return "", err
	}

	downloadURL := path.Join(s.baseURL, filename)
	return downloadURL, nil
}
