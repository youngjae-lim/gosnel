package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// setup sets a rootpath of cli application and database type(postgres or mysql)
func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		gos.RootPath = path
		gos.DB.DbType = os.Getenv("DATABASE_TYPE")
	}
}

// getDSN returns a database connection string for postgres, mysql
func getDSN() string {
	dbType := gos.DB.DbType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		}
		return dsn
	}
	return "mysql://" + gos.BuildDSN()
}

func checkForDB() {
	dbType := gos.DB.DbType

	if dbType == "" {
		exitGracefully(errors.New("no database connection provided in .env"))
	}

	if !fileExists(gos.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("config/database.yml does not exist"))
	}
}

func showHelp() {
	color.Yellow(`Available commands:

	help			- show the help commands
	version			- print application version
	migrate			- runs all up migrations that have not been run previsouly
	migrate down		- reverses the most recent migration
	migrate reset		- runs all down migrations in reverse order, and then all up migrations
	make auth		- creates and runs migrations for authentication tables, and creates models and middlewares
	make migration <name> <format> 	- creates two new up and down migrations in the /migrations folder; format=sql/fizz (default fizz)
	make handler <name>	- creates a stub handler in the /handlers directory
	make model <name> 	- creates a new model in the /data directory
	make session 		- creates a table in the database as a session store
	make mail <name>    	- creates two starter mail templates in the mail directory

	`)
}

// updateSourceFiles is a function that replaces 'myapp' with appURL in a *.go file
func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	// check for an error before doing anything else
	if err != nil {
		return err
	}

	// check if current file is directory
	if fi.IsDir() {
		return nil
	}

	// only check go files
	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}

	if matched {
		// read the file contents
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}

		newContents := strings.Replace(string(read), "myapp", appURL, -1)

		// write the changed file
		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}

// updateSource walks through the entire project folders and updates each import statement
// based on the appURL used to create a gosnel skeleton app.
func updateSource() {
	// walk entire project folder, including subfolders recursively
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGracefully(err)
	}
}
