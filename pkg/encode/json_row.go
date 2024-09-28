package encode

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Encode rows into a json row-oriented structure
func ToJsonRows(w io.Writer, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var results []map[string]any

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			var v any
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			rowMap[col] = v
		}

		results = append(results, rowMap)
	}

	b, err := json.Marshal(results)
	if err != nil {
		return err
	}

	w.Write(b)
	return nil
}
