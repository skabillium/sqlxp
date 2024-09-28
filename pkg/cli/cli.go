package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	flag "github.com/spf13/pflag"
)

const (
	// Database driver types
	DriverMysql    = "mysql"
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite3"

	// Output file formats
	FormatCSV = iota
	FormatJSON

	// JSON orientation options
	OrientationRow = iota
	OrientationColumn
	OrientationArray

	CurrentVersion = "1.0.0"
)

type Database struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	dsn      string
}

func (d *Database) DSN() string {
	if d.dsn != "" {
		return d.dsn
	}

	switch d.Driver {
	case DriverMysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.Name)
	default:
		log.Fatalf("database driver '%s' not supported", d.Driver)
	}

	return ""
}

type Config struct {
	Database     Database
	OutputFile   string
	OutputFormat int
	Orientation  int
	Query        string
	Ping         bool
	Print        bool
}

func ParseArguments() (*Config, error) {
	var (
		dbDriver, user, password, host, db, output, query, file, dsn, orientationStr, format string
		ping, print, help, version                                                           bool
		port                                                                                 int
	)

	flag.StringVarP(&dbDriver, "driver", "D", "", "Database driver")
	flag.StringVarP(&user, "user", "u", "", "Database user")
	flag.StringVarP(&password, "password", "p", "", "Database password")
	flag.StringVarP(&host, "host", "h", "localhost", "Database host")
	flag.IntVarP(&port, "port", "P", 0, "Database port")
	flag.StringVarP(&db, "database", "d", "", "Database name")
	flag.StringVarP(&output, "out", "o", "", "Output file")
	flag.StringVarP(&query, "query", "q", "", "Query to execute")
	flag.StringVarP(&file, "file", "f", "", "File to read query from")
	flag.StringVar(&orientationStr, "orientation", "", "Orientation for json encoding, options: 'row', 'column', 'array'")
	flag.BoolVar(&ping, "ping", false, "Validate the connection without running a query")
	flag.BoolVar(&print, "print", false, "Print the output instead of writing to a file")
	flag.StringVar(&format, "format", "csv", "Format for the output, if you are outputting to a file it will be inferred from the extension")
	flag.BoolVar(&help, "help", false, "Show usage")
	flag.BoolVarP(&version, "version", "v", false, "Print version")

	flag.Parse()

	if len(os.Args) == 1 || help {
		flag.Usage()
		os.Exit(0)
	}

	if version {
		ExitWithMessage("sqlxp version", CurrentVersion)
	}

	if flag.NArg() == 1 && dsn == "" {
		dsn = flag.Args()[0]
	}

	if port < 0 {
		return nil, errors.New("port must be a positive integer")
	}

	var driver string
	if dsn != "" {
		driver = getDriverFromDSN(dsn)
	} else {
		if dbDriver == "" || user == "" || password == "" || host == "" || port == 0 || db == "" {
			return nil, errors.New("either dsn or all of the connection options are required")
		}

		var ok bool
		driver, ok = getDriver(dbDriver)
		if !ok {
			return nil, fmt.Errorf("unsupported driver '%s'", dbDriver)
		}
	}

	if file != "" {
		if !strings.HasSuffix(file, ".sql") {
			return nil, errors.New("File input should be sql file")
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		query = string(content)
	}

	if !ping {
		if query == "" {
			return nil, errors.New("query cannot be empty")
		}

		if !isSelectQuery(query) {
			return nil, errors.New("only 'select' queries can be executed")
		}
	}

	if output != "" {
		format = getFileExtension(output)
	} else {
		if format == "" {
			return nil, errors.New("format is required when no output is specified")
		}
	}

	var outputFormat int
	switch format {
	case "csv":
		outputFormat = FormatCSV
	case "json":
		outputFormat = FormatJSON
	default:
		return nil, fmt.Errorf("unsupported output format '%s'", format)
	}

	var orientation int
	if outputFormat == FormatJSON {
		switch orientationStr {
		case "", "r", "row", "rows":
			orientation = OrientationRow
		case "c", "col", "column", "columns":
			orientation = OrientationColumn
		case "a", "arr", "array", "arrays":
			orientation = OrientationArray
		default:
			return nil, fmt.Errorf("invalid json orientation '%s'", orientationStr)
		}
	}

	if !print && output == "" {
		return nil, errors.New("no output file specified")
	}

	return &Config{
		Database: Database{
			Driver:   driver,
			User:     user,
			Password: password,
			Host:     host,
			Port:     strconv.Itoa(port),
			Name:     db,
			dsn:      dsn,
		},
		OutputFile:   output,
		OutputFormat: outputFormat,
		Orientation:  orientation,
		Query:        query,
		Ping:         ping,
		Print:        print,
	}, nil
}

func Fatal(a ...any) {
	PrintError(a...)
	os.Exit(1)
}

func PrintError(a ...any) {
	p := append([]any{"Error:"}, a...)
	fmt.Fprintln(os.Stderr, p...)
}

func ExitWithMessage(a ...any) {
	fmt.Println(a...)
	os.Exit(0)
}

func getFileExtension(file string) string {
	var start int
	runes := []rune(file)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == '.' {
			start = i + 1
			break
		}
	}
	return string(runes[start:])
}

func getDriver(driver string) (string, bool) {
	d := strings.ToLower(driver)
	switch d {
	case DriverPostgres, "pg":
		return DriverPostgres, true
	case DriverSqlite, "sqlite":
		return DriverSqlite, true
	case DriverMysql, "mariadb", "maria":
		return DriverMysql, true
	default:
		return "", false
	}
}

func getDriverFromDSN(dsn string) string {
	if strings.HasPrefix(dsn, "postgres://") {
		return DriverPostgres
	} else if strings.Contains(dsn, "@tcp(") {
		return DriverMysql
	} else {
		return DriverSqlite
	}
}

func isSelectQuery(query string) bool {
	for i, c := range query {
		if unicode.IsSpace(c) {
			continue
		}
		return strings.ToLower(query[i:i+6]) == "select"
	}
	return false
}
