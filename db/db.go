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
// exec package is used to execute commands
// filepath package is used to manipulate file paths
// log package is used to log messages to the console
// go-sqlite3 package is the driver used to interact with SQLite databases
import (
	"aio/helpers"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
// Path functions
//

func getExecDir() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Failed to get executable path", "error", err)
	}
	return filepath.Dir(execPath)
}

// getPath function returns the full path of a file.
// it is use to retrieve the path from the executable file.
func getPath(path string) string {
	execDir := getExecDir()
	return filepath.Join(execDir, path)
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// git functions
//

// linkRepo function adds a link to a remote repository.
// it is used to link the database to a remote repository for versioning.
// this action is optional and can be skipped by the user and performed later.
func linkRepo() {
	if helpers.RunConfirm("Do you want to add a link to a remote repository?") {
		fmt.Println("Please enter the remote repository SSH link:")
		remote := helpers.RunInput("SSH link")
		cmd := exec.Command("git", "remote", "add", "origin", remote)
		cmd.Dir = getExecDir()
		err := cmd.Run()
		if err != nil {
			log.Fatal("Failed to add remote repository", "error", err)
		}
	}
}

// remoteExists function checks if a remote repository is linked to the database.
// it is used to check if a remote repository is linked to the database for versioning.
func remoteExists() bool {
	cmd := exec.Command("git", "remote")
	cmd.Dir = getExecDir()
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to get remote repository", "error", err)
	}
	return strings.TrimSpace(string(output)) != ""
}

// gitPreVersioning function creates a new branch for versioning.
// it is used to create a new branch before making changes to the database.
func gitPreVersioning() string {
	var cmd *exec.Cmd
	var err error
	execDir := getExecDir()

	if remoteExists() {
		cmd := exec.Command("git", "fetch")
		cmd.Dir = execDir
		err := cmd.Run()
		if err != nil {
			log.Fatal("Failed to fetch remote repository", "error", err)
		}

		cmd = exec.Command("git", "pull")
		cmd.Dir = execDir
		err = cmd.Run()
		if err != nil {
			log.Fatal("Failed to pull remote repository", "error", err)
		}
	}

	branchName := time.Now().Format("2006-01-02-15-04-05")
	cmd = exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to create branch", "error", err)
	}

	return branchName
}

// gitPostVersioning function commits and pushes the changes to the database.
// it is used to commit and push the changes to the database after making changes.
func gitPostVersioning(branchName string) {
	re := remoteExists()
	execDir := getExecDir()
	dbfile := getPath("data.db")

	cmd := exec.Command("git", "add", dbfile)
	cmd.Dir = execDir
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to add database file", "error", err)
	}

	cm := fmt.Sprintf("Version %s", branchName)
	cmd = exec.Command("git", "commit", "-m", cm)
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to commit changes", "error", err)
	}

	cmd = exec.Command("git", "rev-parse", "--verify", "main")
	cmd.Dir = execDir
	err = cmd.Run()

	if err != nil {
		cmd = exec.Command("git", "checkout", "-b", "main")
		cmd.Dir = execDir
		err = cmd.Run()
		if err != nil {
			log.Fatal("Failed to create main branch", "error", err)
		}
	}

	if re {
		cmd = exec.Command("git", "push", "-u", "origin", branchName)
		cmd.Dir = execDir
		err = cmd.Run()
		if err != nil {
			log.Fatal("Failed to push changes", "error", err)
		}
	}

	cmd = exec.Command("git", "checkout", "main")
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to checkout main branch", "error", err)
	}

	cmd = exec.Command("git", "merge", branchName)
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to merge branches", "error", err)
	}

	if re {
		cmd = exec.Command("git", "push")
		cmd.Dir = execDir
		err = cmd.Run()
		if err != nil {
			log.Fatal("Failed to push changes", "error", err)
		}
	}
}

// gitRollback function rolls back the changes made to the database.
// it is used to roll back the changes made to the database if an error occurs.
func gitRollback(branchName string) {
	execDir := getExecDir()

	cmd := exec.Command("git", "checkout", "main")
	cmd.Dir = execDir
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to checkout main branch", "error", err)
	}

	cmd = exec.Command("git", "branch", "-D", branchName)
	cmd.Dir = execDir
	err = cmd.Run()
	if err != nil {
		log.Fatal("Failed to delete branch", "error", err)
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// Database functions
//

// getDb function returns a pointer to a sql.DB object.
// it is used to open the database file and return a pointer to the database object.
// if the database file does not exist, it creates a new one.
func getDb() *sql.DB {
	dbfile := getPath("data.db")

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal("Failed to open database", "error", err)
	}
	return db
}

// Init function initializes the database.
// the funciton initialize also git for the db versioning
// it is used to create the database file and tables if they do not exist.
// it loads the main sql file that contains the queries to create the tables, indexes, and triggers.
func Init() {
	var cmd *exec.Cmd
	execDir := getExecDir()

	// check if git is initialized
	if _, err := os.Stat(getPath(".git")); os.IsNotExist(err) {
		log.Info("Initializing git...")
		cmd = exec.Command("git", "init")
		cmd.Dir = execDir
		err = cmd.Run()
		if err != nil {
			log.Fatal("Failed to initialize git", "error", err)
		}

		linkRepo()

		log.Info("...git initialized successfully")
	}

	dbfile := getPath("data.db")
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		log.Warn("Database file not found, creating a new one", "file", dbfile)

		var file *os.File
		file, err = os.Create(dbfile)
		if err != nil {
			log.Fatal("Failed to create database file", "error", err)
		}

		file.Close()
		log.Info("Database file created", "file", dbfile)

	}

	log.Info("Initializing database...")
	db := getDb()
	defer db.Close()

	query := loadQuery("tables")

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to execute query", "query", query, "error", err)
	}
	log.Info("...database initialized successfully")
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
