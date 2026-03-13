package repository

import (
	"context"
	"fmt"
	"strings"
)

type QueryEngine interface {
	Execute(ctx context.Context, query string) ([]map[string]interface{}, error)
	EngineName() string
}

type postgresQueryEngine struct {
	queryDB *QueryDB
}

func NewPostgresQueryEngine(queryDB *QueryDB) QueryEngine {
	return &postgresQueryEngine{queryDB: queryDB}
}

func (e *postgresQueryEngine) EngineName() string {
	return "postgres"
}

func (e *postgresQueryEngine) Execute(ctx context.Context, query string) ([]map[string]interface{}, error) {
	// Validate: only SELECT statements allowed
	normalized := strings.TrimSpace(strings.ToUpper(query))
	if !strings.HasPrefix(normalized, "SELECT") && !strings.HasPrefix(normalized, "WITH") {
		return nil, fmt.Errorf("only SELECT (or WITH ... SELECT) queries are allowed")
	}

	// Block dangerous keywords
	forbidden := []string{"INSERT ", "UPDATE ", "DELETE ", "DROP ", "ALTER ", "CREATE ", "TRUNCATE ", "GRANT ", "REVOKE "}
	for _, keyword := range forbidden {
		if strings.Contains(normalized, keyword) {
			return nil, fmt.Errorf("forbidden SQL keyword detected: %s", strings.TrimSpace(keyword))
		}
	}

	rows, err := e.queryDB.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns error: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for readability
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}
