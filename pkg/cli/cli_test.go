package cli

import "testing"

func Test_getFileExtension(t *testing.T) {
	tests := [][2]string{
		{"example.json", "json"},
		{"example.csv", "csv"},
		{".env", "env"},
		{"no-ext", "no-ext"},
	}

	for _, test := range tests {
		file := test[0]
		expected := test[1]

		result := getFileExtension(file)
		if result != expected {
			t.Errorf("error testing '%s': expected '%s', got '%s'", file, expected, result)
		}
	}
}

func Test_getDriverFromDSN(t *testing.T) {
	tests := [][2]string{
		{"postgres://user:password@localhost/dbname?sslmode=disable", DriverPostgres},
		{"user:password@tcp(127.0.0.1:3306)/dbname", DriverMysql},
		{"other", DriverSqlite},
	}

	for _, test := range tests {
		file := test[0]
		expected := test[1]

		result := getDriverFromDSN(file)
		if result != expected {
			t.Errorf("error testing '%s': expected '%s', got '%s'", file, expected, result)
		}
	}
}
