package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"skabillium.io/sqlxp/pkg/cli"
	"skabillium.io/sqlxp/pkg/encode"
)

func main() {
	args, err := cli.ParseArguments()
	if err != nil {
		cli.Fatal(err)
	}

	db, err := sql.Open(args.Database.Driver, args.Database.DSN())
	if err != nil {
		cli.Fatal("could not connect to database", err)
	}
	defer db.Close()

	if args.Ping {
		err := db.Ping()
		if err != nil {
			cli.Fatal(err)
		} else {
			fmt.Println("Pong")
		}
		return
	}

	rows, err := db.Query(args.Query)
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	switch args.OutputFormat {
	case cli.FormatCSV:
		err = encode.ToCSV(&builder, rows)
		if err != nil {
			cli.Fatal(err)
		}
	case cli.FormatJSON:
		if args.Orientation == cli.OrientationRow {
			err := encode.ToJsonRows(&builder, rows)
			if err != nil {
				cli.Fatal(err)
			}
		} else if args.Orientation == cli.OrientationColumn {
			err := encode.ToJsonColumns(&builder, rows)
			if err != nil {
				cli.Fatal(err)
			}
		} else {
			err := encode.ToJsonArray(&builder, rows)
			if err != nil {
				cli.Fatal(err)
			}
		}

	default:
		cli.Fatal("unreachable")
	}

	content := builder.String()

	if args.Print {
		fmt.Print(content)
		return
	}

	file, err := os.Create(args.OutputFile)
	if err != nil {
		cli.Fatal("could not create file", err)
	}
	// Ensure the file is closed after writing
	defer file.Close()

	// Write the string to the file
	_, err = file.WriteString(content)
	if err != nil {
		cli.Fatal(err)
	}
}
