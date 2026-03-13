package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// KnowledgeSearcher is the interface that the knowledge service must implement.
type KnowledgeSearcher interface {
	Search(ctx context.Context, query string, category string, limit int) ([]KnowledgeResult, error)
}

type KnowledgeResult struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type KnowledgeSearchTool struct {
	searcher KnowledgeSearcher
}

func NewKnowledgeSearchTool(searcher KnowledgeSearcher) *KnowledgeSearchTool {
	return &KnowledgeSearchTool{searcher: searcher}
}

func (t *KnowledgeSearchTool) Name() string {
	return "knowledge_search"
}

func (t *KnowledgeSearchTool) Description() string {
	return "테이블 스키마, 예시 SQL, 비즈니스 용어 등 지식 베이스를 검색합니다. SQL 작성 전에 관련 스키마를 확인하세요."
}

func (t *KnowledgeSearchTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "검색할 키워드 또는 질문",
			},
			"category": map[string]interface{}{
				"type":        "string",
				"description": "카테고리 필터: schema, example_sql, glossary (비워두면 전체 검색)",
				"enum":        []string{"schema", "example_sql", "glossary", ""},
			},
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "반환할 최대 결과 수 (기본값: 5)",
			},
		},
		"required": []string{"query"},
	}
}

func (t *KnowledgeSearchTool) Execute(ctx context.Context, args string) (string, error) {
	var params struct {
		Query    string `json:"query"`
		Category string `json:"category"`
		Limit    int    `json:"limit"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	if params.Limit <= 0 {
		params.Limit = 5
	}

	results, err := t.searcher.Search(ctx, params.Query, params.Category, params.Limit)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "관련 지식을 찾지 못했습니다.", nil
	}

	output, _ := json.MarshalIndent(results, "", "  ")
	return string(output), nil
}
