package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// CSVGenerator is the interface that the CSV service must implement.
type CSVGenerator interface {
	GenerateFromQuery(ctx context.Context, query string) (string, error) // returns download URL
}

type CSVDownloadTool struct {
	generator CSVGenerator
}

func NewCSVDownloadTool(generator CSVGenerator) *CSVDownloadTool {
	return &CSVDownloadTool{generator: generator}
}

func (t *CSVDownloadTool) Name() string {
	return "csv_download"
}

func (t *CSVDownloadTool) Description() string {
	return "SQL 쿼리 결과를 CSV 파일로 생성하고 다운로드 URL을 반환합니다. 결과 행이 20개를 초과할 때 사용하세요."
}

func (t *CSVDownloadTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "CSV로 내보낼 SQL SELECT 쿼리",
			},
		},
		"required": []string{"query"},
	}
}

func (t *CSVDownloadTool) Execute(ctx context.Context, args string) (string, error) {
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	if params.Query == "" {
		return "", fmt.Errorf("query is required")
	}

	url, err := t.generator.GenerateFromQuery(ctx, params.Query)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CSV 파일이 생성되었습니다.\n다운로드 URL: %s", url), nil
}
