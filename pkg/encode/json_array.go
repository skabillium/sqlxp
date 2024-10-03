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

	// Write the opening bracket of the JSON array
	w.Write([]byte("["))

	isFirst := true
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

		// Marshal row data
		b, err := json.Marshal(row)
		if err != nil {
			return err
		}

		// Write a comma before the next row if it's not the first row
		if !isFirst {
			w.Write([]byte(","))
		}
		isFirst = false

		w.Write(b)
	}

	// Write the closing bracket of the JSON array
	w.Write([]byte("]"))

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
