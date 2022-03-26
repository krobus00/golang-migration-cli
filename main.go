package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		dir      = flags.String("dir", fmt.Sprintf("./%s", os.Getenv("MIGRATION_PATH")), "directory with migration files")
		dbstring = flags.String("dbstring",
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
				os.Getenv("DATABASE_USERNAME"),
				os.Getenv("DATABASE_PASSWORD"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
			),
			"connection string",
		)
	)
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) > 1 && args[0] == "create" {
		arguments := []string{}
		arguments = append(arguments, args[1:]...)
		arguments = append(arguments, "sql")
		if err := goose.Run("create", nil, *dir, arguments...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	driver, command := os.Getenv("DATABASE_DIALECT"), args[0]

	switch driver {
	case "postgres", "mysql", "sqlite3", "redshift":
		if err := goose.SetDialect(driver); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("%q driver not supported\n", driver)
	}

	switch *dbstring {
	case "":
		log.Fatalf("-dbstring=%q not supported\n", *dbstring)
	default:
	}

	if driver == "redshift" {
		driver = "postgres"
	}

	db, err := sql.Open(driver, *dbstring)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", *dbstring, err)
	}

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
Drivers:
    postgres
    mysql
    sqlite3
    redshift
Examples:
    goose status
    goose create init
    goose create add_some_column
    goose create fetch_user_data
    goose up
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
`
)
