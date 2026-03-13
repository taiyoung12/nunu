package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// MemorySearcher is the interface that the memory service must implement.
type MemorySearcher interface {
	SearchSimilar(ctx context.Context, query string, limit int) ([]MemoryResult, error)
}

type MemoryResult struct {
	Question   string  `json:"question"`
	Summary    string  `json:"summary"`
	SQLUsed    string  `json:"sql_used"`
	Similarity float64 `json:"similarity"`
}

type MemorySearchTool struct {
	searcher MemorySearcher
}

func NewMemorySearchTool(searcher MemorySearcher) *MemorySearchTool {
	return &MemorySearchTool{searcher: searcher}
}

func (t *MemorySearchTool) Name() string {
	return "memory_search"
}

func (t *MemorySearchTool) Description() string {
	return "과거에 성공적으로 답변한 유사 질문과 사용된 SQL을 검색합니다. 새로운 질문을 받으면 항상 먼저 이 도구를 사용하세요."
}

func (t *MemorySearchTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "검색할 질문 텍스트",
			},
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "반환할 최대 결과 수 (기본값: 3)",
			},
		},
		"required": []string{"query"},
	}
}

func (t *MemorySearchTool) Execute(ctx context.Context, args string) (string, error) {
	var params struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	if params.Limit <= 0 {
		params.Limit = 3
	}

	results, err := t.searcher.SearchSimilar(ctx, params.Query, params.Limit)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "유사한 과거 기억을 찾지 못했습니다.", nil
	}

	output, _ := json.MarshalIndent(results, "", "  ")
	return string(output), nil
}
