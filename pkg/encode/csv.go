package encode

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// Encode rows into a csv
func ToCSV(w io.Writer, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write the CSV headers (column names)
	if err := writer.Write(columns); err != nil {
		return err
	}

	for rows.Next() {
		// Prepare a slice to hold the raw data for each column
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		rowData := make([]string, len(columns))
		for i, val := range values {
			switch v := val.(type) {
			case nil:
				rowData[i] = "NULL"
			case []byte:
				rowData[i] = string(v)
			case int64:
				rowData[i] = strconv.FormatInt(v, 10)
			case float64:
				rowData[i] = strconv.FormatFloat(v, 'f', -1, 64)
			default:
				rowData[i] = fmt.Sprintf("%v", v)
			}
		}

		if err := writer.Write(rowData); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
