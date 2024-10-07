//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// File Name: db.go
// Created by: Matteo Tagliapietra 2024-09-01
// Last Update: 2024-10-05

// This is the main file of the db package.
// It contains the Init function that initializes the database.
// It is used to create the database file and tables if they do not exist.

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

// db package is used to interact with the database.
package db

// imports the necessary packages
// sql package is used to interact with the database
// embed package is used to embed files in the binary
// os package is used to read and create files
// filepath package is used to manipulate file paths
// log package is used to log messages to the console
// go-sqlite3 package is the driver used to interact with SQLite databases
import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// Embed files
//

//go:embed queries/*.sql
var sqlFiles embed.FS

// loadQuery function reads the content of a sql file and returns it as a string.
// it is used to load the content of the sql files that contain the queries to execute.
func loadQuery(filename string) string {
	query, err := sqlFiles.ReadFile("queries/" + filename + ".sql")
	if err != nil {
		log.Fatal("Failed to read query file", "file", filename, "error", err)
	}
	return string(query)
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// Database functions
//

// getPath function returns the full path of a file.
// it is use to retrieve the path from the executable file.
func getPath(path string) string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path", "error", err)
	}
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, path)
}

// getDb function returns a pointer to a sql.DB object.
// it is used to open the database file and return a pointer to the database object.
// if the database file does not exist, it creates a new one.
func getDb() *sql.DB {
	dbfile := getPath("data.db")
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		log.Warn("Database file not found, creating a new one", "file", dbfile)
		file, err := os.Create(dbfile)
		if err != nil {
			log.Fatal("Failed to create database file", "error", err)
		}
		file.Close()
		log.Info("Database file created", "file", dbfile)
	}

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal("Failed to open database", "error", err)
	}
	return db
}

// Init function initializes the database.
// it is used to create the database file and tables if they do not exist.
// it loads the main sql file that contains the queries to create the tables, indexes, and triggers.
func Init() {
	log.Info("Initializing database...")
	db := getDb()
	defer db.Close()

	query := loadQuery("tables")

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to execute query", "query", query, "error", err)
	}
	log.Info("...database initialized successfully")
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
