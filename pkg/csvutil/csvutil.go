package csvutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// WriteCSV writes rows of data to a CSV file and returns the file path.
func WriteCSV(storagePath string, rows []map[string]interface{}) (string, error) {
	if len(rows) == 0 {
		return "", fmt.Errorf("no data to write")
	}

	// Collect and sort column names
	colSet := make(map[string]struct{})
	for _, row := range rows {
		for k := range row {
			colSet[k] = struct{}{}
		}
	}
	columns := make([]string, 0, len(colSet))
	for k := range colSet {
		columns = append(columns, k)
	}
	sort.Strings(columns)

	// Generate filename
	filename := fmt.Sprintf("query_%s.csv", time.Now().Format("20060102_150405"))
	filePath := filepath.Join(storagePath, filename)

	// Ensure directory exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("create csv directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("create csv file: %w", err)
	}
	defer file.Close()

	// Write BOM for Excel compatibility
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(columns); err != nil {
		return "", fmt.Errorf("write csv header: %w", err)
	}

	// Write rows
	for _, row := range rows {
		record := make([]string, len(columns))
		for i, col := range columns {
			if v, ok := row[col]; ok {
				record[i] = fmt.Sprintf("%v", v)
			}
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("write csv row: %w", err)
		}
	}

	return filename, nil
}
