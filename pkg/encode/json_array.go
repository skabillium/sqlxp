package encode

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Encode rows into a json array of arrays structure
func ToJsonArray(w io.Writer, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var results [][]any

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		row := make([]any, len(columns))
		for i := range columns {
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				row[i] = string(b)
			} else {
				row[i] = val
			}
		}

		results = append(results, row)
	}

	b, err := json.Marshal(results)
	if err != nil {
		return err
	}

	w.Write(b)
	return nil
}
