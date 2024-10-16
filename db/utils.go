// db package contains the functions to interact with the database.
package db

import (
	"aio/logger"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

// loadQuery function reads the content of a sql file and returns it as a string.
// it is used to load the content of the sql files that contain the queries to execute.
func loadQuery(filename string) string {
	query, err := sqlFiles.ReadFile("queries/" + filename + ".sql")
	logger.Fatal("Failed to read query file", err, "file", filename)
	return string(query)
}

// getExecDir function returns the directory of the executable file.
// it is used to run all the commands in the directory of the executable file.
// this maintaning the integrity of the data.
func getExecDir() string {
	execPath, err := os.Executable()
	logger.Fatal("Failed to get executable path", err)
	return filepath.Dir(execPath)
}

// getPath function returns the full path of a file.
// it is use to retrieve the path from the executable file.
func getPath(path string) string {
	execDir := getExecDir()
	return filepath.Join(execDir, path)
}

// cmdExec function executes a command and returns the output.
// it is used to execute a command and return the output.
// it also sets the working directory to the executable directory.
func cmdExec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = getExecDir()
	output, err := cmd.Output()
	return string(output), err
}

// getDb function returns a pointer to a sql.DB object.
// it is used to open the database file and return a pointer to the database object.
// if the database file does not exist, it creates a new one.
func getDb() (*sql.DB, error) {
	dbfile := getPath("data.db")
	return sql.Open("sqlite3", dbfile)
}

// backup function creates a backup of the database file.
// it is used to create a backup of the database file before making changes to the database.
func backup() {
	dbfile := getPath("data.db")
	backupfile := getPath(fmt.Sprintf("data_backup_%s.db", time.Now().Format("20060102150405")))
	input, err := os.Open(dbfile)
	logger.Fatal("Failed to open database file", err)
	defer input.Close()

	output, err := os.Create(backupfile)
	logger.Fatal("Failed to create backup file", err)
	defer output.Close()

	_, err = io.Copy(output, input)
	logger.Fatal("Failed to copy database file", err)
}

// do function executes a query on the database.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func do(query string, args ...any) error {
	db, err := getDb()
	if err != nil {
		return err
	}

	defer db.Close()

	q := loadQuery(query)
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return err
}

// gets function executes a query on the database and returns the result.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func gets(query string, args ...any) (*sql.Rows, error) {
	db, err := getDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	q := loadQuery(query)
	rows, err := db.Query(q, args...)

	return rows, err
}

// get function executes a query on the database and returns the result.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func get(query string, args ...any) (*sql.Row, error) {
	db, err := getDb()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	q := loadQuery(query)
	row := db.QueryRow(q, args...)

	return row, err
}
