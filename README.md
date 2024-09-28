# SQLXP

sqlxp is a CLI for exporting sql queries to multiple file formats.

## Usage
Select users from an sqlite database and write the result as a csv
```sh
sqlxp -q 'SELECT * FROM users' -o users.csv example.db
```

For more usage information run `sqlxp --help`.

## Features

- Supports multiple SQL databases (PostgreSQL, MySQL, SQLite)
- Execute raw SQL queries and process results
- Outputs results in json & csv
- Option to print output to the console or write to a file

## Installation
To install the tool, make sure you have [Go](https://golang.org/doc/install) installed. Then, clone this repository and run:

```bash
git clone https://github.com/skabillium/sqlxp.git
cd sqlxp
go mod tidy
go build -o sqlxp ./cmd # Or make build
```
