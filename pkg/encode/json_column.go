package encode

import (
	"database/sql"
	"encoding/json"
	"io"
)

// Encode rows into a json column-oriented structure
func ToJsonColumns(w io.Writer, rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnData := make(map[string][]any)

	for _, col := range columns {
		columnData[col] = []any{}
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		for i, col := range columns {
			var v any
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}

			columnData[col] = append(columnData[col], v)
		}
	}

	b, err := json.Marshal(columnData)
	if err != nil {
		return err
	}

	w.Write(b)
	return nil
}
