package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// QueryExecutor is the interface that the query service must implement.
type QueryExecutor interface {
	ExecuteReadOnly(ctx context.Context, query string, maxRows int) ([]map[string]interface{}, error)
}

type QueryPostgresTool struct {
	executor QueryExecutor
	maxRows  int
}

func NewQueryPostgresTool(executor QueryExecutor, maxRows int) *QueryPostgresTool {
	if maxRows <= 0 {
		maxRows = 100
	}
	return &QueryPostgresTool{executor: executor, maxRows: maxRows}
}

func (t *QueryPostgresTool) Name() string {
	return "query_postgres"
}

func (t *QueryPostgresTool) Description() string {
	return "PostgreSQL 데이터베이스에서 SELECT 쿼리를 실행합니다. INSERT, UPDATE, DELETE는 허용되지 않습니다."
}

func (t *QueryPostgresTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "실행할 SQL SELECT 쿼리",
			},
		},
		"required": []string{"query"},
	}
}

func (t *QueryPostgresTool) Execute(ctx context.Context, args string) (string, error) {
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}
	if params.Query == "" {
		return "", fmt.Errorf("query is required")
	}

	rows, err := t.executor.ExecuteReadOnly(ctx, params.Query, t.maxRows)
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "쿼리 결과가 없습니다. (0 rows)", nil
	}

	result := map[string]interface{}{
		"row_count": len(rows),
		"rows":      rows,
	}
	output, _ := json.MarshalIndent(result, "", "  ")
	return string(output), nil
}
