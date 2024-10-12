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
	"sync"
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

// cmdExec function executes a command and returns the output.
// it is used to execute a command and return the output.
// it also sets the working directory to the executable directory.
func cmdExec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = getExecDir()
	output, err := cmd.Output()
	return string(output), err
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// git functions
//

var (
	remoteexixtsCheck bool
	remoteexixts      bool
	hasremoteCheck    bool
	hasremote         bool
	haslocalCheck     bool
	haslocal          bool
	haschangesCheck   bool
	haschanges        bool
	hasmainCheck      bool
	hasmain           bool
	isAligned         bool
	wu                sync.WaitGroup
)

// linkRepo function adds a link to a remote repository.
// it is used to link the database to a remote repository for versioning.
// this action is optional and can be skipped by the user and performed later.
func linkRepo() {
	if helpers.RunConfirm("Do you want to add a link to a remote repository?") {
		fmt.Println("Please enter the remote repository name:")
		remote := helpers.RunInput("YourUsername/repo-name")
		output, err := cmdExec("git", "remote", "add", "origin", "git@github.com:"+remote)
		if err != nil {
			log.Fatal("Failed to add remote repository", "output", output, "error", err)
		}
	}
}

// remoteExists function checks if a remote repository is linked to the database.
// it is used to check if a remote repository is linked to the database for versioning.
func remoteExists() bool {
	if !remoteexixtsCheck {
		output, err := cmdExec("git", "remote")
		if err != nil {
			log.Fatal("Failed to get remote repository", "output", output, "error", err)
		}
		remoteexixts = strings.TrimSpace(output) != ""
		remoteexixtsCheck = true
	}
	return remoteexixts
}

// hasRemoteCommits function checks if there are remote commits.
// it is used to check if there are remote commits before making changes to the database.
func hasRemoteCommits() bool {
	if !hasremoteCheck {
		cmd := exec.Command("git", "ls-remote", "--heads", "origin")
		cmd.Dir = getExecDir()
		output, err := cmd.Output()
		hasremote = err == nil && len(output) > 0 // Restituisce true se ci sono branch remoti
		hasremoteCheck = true
	}
	return hasremote
}

// hasLocalCommits function checks if there are local commits.
// it is used to check if there are local commits before making changes to the database.
func hasLocalCommits() bool {
	if !haslocalCheck {
		_, err := cmdExec("git", "rev-parse", "--verify", "HEAD")
		haslocal = err == nil // Return true if there are local commits
		haslocalCheck = true
	}
	return haslocal
}

// hasChanges function checks if there are changes in the database file.
func hasChanges() bool {
	if !haschangesCheck {
		dbfile := getPath("data.db")
		output, err := cmdExec("git", "status", "--porcelain", dbfile)
		if err != nil {
			log.Fatal("Failed to check changes", "output", output, "error", err)
		}

		haschanges = strings.TrimSpace(output) != ""
		haschangesCheck = true
	}
	return haschanges
}

// hasMain function checks if the main branch exists.
// it is used to check if the main branch exists before making changes to the database.
func hasMain() bool {
	if !hasmainCheck {
		_, err := cmdExec("git", "show-ref", "--verify", "refs/heads/main")
		hasmain = err == nil
		hasmainCheck = true
	}
	return hasmain
}

// gitFetch function fetches the remote repository.
// it is used to fetch the remote repository before making changes to the database.
func gitFetch() {
	if remoteExists() && hasRemoteCommits() && !isAligned {
		output, err := cmdExec("git", "fetch")
		if err != nil {
			log.Fatal("Failed to fetch remote repository", "output", output, "error", err)
		}
	}
}

// gitPull function pulls the remote repository.
// it is used to pull the remote repository before making changes to the database.
func gitPull() {
	if remoteExists() && hasRemoteCommits() && !isAligned {
		gitMain()
		output, err := cmdExec("git", "pull", "origin", "main")
		if err != nil {
			log.Fatal("Failed to pull remote repository", "output", output, "error", err)
		}
		isAligned = true
	}
}

// gitMain function checks out the main branch.
// it is used to checkout the main branch before making changes to the database.
func gitMain() {
	if !hasMain() {
		output, err := cmdExec("git", "checkout", "-b", "main")
		if err != nil {
			log.Fatal("Failed to create main branch", "output", output, "error", err)
		}
		hasmain = true
	} else {
		output, err := cmdExec("git", "checkout", "main")
		if err != nil {
			log.Fatal("Failed to checkout main branch", "output", output, "error", err)
		}
	}
}

// initialCommit function commits the database file.
// it is used to commit the database file to the local repository.
// it also renames the branch to main and pushes the changes to the remote repository if it exists.
func initialCommit() {
	dbfile := getPath("data.db")
	var output string
	var err error

	if !hasRemoteCommits() && !hasLocalCommits() {
		output, err = cmdExec("git", "add", dbfile)
		if err != nil {
			log.Fatal("Failed to add database file", "output", output, "error", err)
		}

		// do the initial commit
		output, err = cmdExec("git", "commit", "-m", "initial commit")
		if err != nil {
			log.Fatal("Failed to commit changes", "output", output, "error", err)
		}

		output, err = cmdExec("git", "branch", "-M", "main")
		if err != nil {
			log.Fatal("Failed to rename branch", "output", output, "error", err)
		}

		haslocal = true
	}
}

// gitNewBranch function creates a new branch.
// it is used to create a new branch before making changes to the database.
func gitNewBranch() string {
	branchName := time.Now().Format("2006-01-02-15-04-05")
	output, err := cmdExec("git", "checkout", "-b", branchName)
	if err != nil {
		log.Fatal("Failed to create branch", "output", output, "error", err)
	}
	return branchName
}

// push function pushes the changes to the remote repository.
func push() {
	wu.Wait()
	// push the changes to the remote repository if it exists
	if remoteExists() {
		log.Info("Pushing changes to the remote repository...")
		output, err := cmdExec("git", "push", "-u", "origin", "main")
		if err != nil {
			log.Fatal("Failed to push changes", "output", output, "error", err)
		}
		hasremote = true
	}
}

// gitCommit function commits the changes made to the database.
// it is used to commit the changes made to the database to the local repository.
// it also merges the branch with the main branch and pushes the changes to the remote repository if it exists.
func gitCommit(branchName string) {
	defer wu.Done()
	output, err := cmdExec("git", "add", getPath("data.db"))
	if err != nil {
		log.Fatal("Failed to add database file", "output", output, "error", err)
	}

	output, err = cmdExec("git", "commit", "-m", fmt.Sprintf("Version %s", branchName))
	if err != nil {
		log.Fatal("Failed to commit changes", "output", output, "error", err)
	}

	gitMain()
	output, err = cmdExec("git", "merge", branchName, "--no-ff")
	if err != nil {
		log.Fatal("Failed to merge branches", "output", output, "error", err)
	}
}

// gitRollback function rolls back the changes made to the database.
// it is used to roll back the changes made to the database if an error occurs.
func gitRollback(branchName string) {
	output, err := cmdExec("git", "checkout", "main")
	if err != nil {
		log.Fatal("Failed to checkout main branch", "output", output, "error", err)
	}

	output, err = cmdExec("git", "branch", "-D", branchName)
	if err != nil {
		log.Fatal("Failed to delete branch", "output", output, "error", err)
	}
}

// gitFlow function is a wrapper function that executes a series of git commands.
// it is used to execute a series of git commands in a transaction to ensure the integrity of the data.
// if an error occurs, it rolls back the changes.
func gitFlow(action func() error) error {
	wu.Add(1)
	gitPull()
	branch := gitNewBranch()
	err := action()
	if err != nil {
		gitRollback(branch)
		return err
	}

	go gitCommit(branch)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////

//
// Database functions
//

// getDb function returns a pointer to a sql.DB object.
// it is used to open the database file and return a pointer to the database object.
// if the database file does not exist, it creates a new one.
func getDb() (*sql.DB, error) {
	dbfile := getPath("data.db")
	return sql.Open("sqlite3", dbfile)
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

// Init function initializes the database.
// the funciton initialize also git for the db versioning
// it is used to create the database file and tables if they do not exist.
// it loads the main sql file that contains the queries to create the tables, indexes, and triggers.
func Init() {

	// check if git is initialized
	if _, err := os.Stat(getPath(".git")); os.IsNotExist(err) {
		log.Info("Initializing git...")
		output, err := cmdExec("git", "init")
		if err != nil {
			log.Fatal("Failed to initialize git", "output", output, "error", err)
		}

		linkRepo()
		if remoteExists() {
			log.Info("Fetching remote repository...")
			gitFetch()
			gitPull()
		}

		gitIgnorePath := getPath(".gitignore")
		if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
			log.Warn(".gitignore file not found, creating a new one", "file", gitIgnorePath)

			content := `*
!data.db
`
			err := os.WriteFile(gitIgnorePath, []byte(content), 0644)
			if err != nil {
				log.Fatal("Failed to create .gitignore file", "error", err)
			}

			log.Info(".gitignore file created", "file", gitIgnorePath)
		}

		log.Info("Git initialized successfully!\n")
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
	err = do("tables")
	if err != nil {
		log.Fatal("Failed to create tables", "error", err)
	}

	if hasChanges() {
		initialCommit()
	}

	log.Info("Database initialized successfully!")
}

// save function saves the changes made to the database.
// it is used to save the changes made to the database and push them to the remote repository if it exists.
// it also creates a new branch for the changes.
func save() {
	err := gitFlow(func() error {
		return do("push_schedules_create")
	})
	if err != nil {
		log.Fatal("Failed to save new push schedule log", "error", err)
	}
	push()
}

// AutoSave function saves the changes made to the database automatically.
// check if there are push schedules in the database that have the created_at field se to the current day
// if there aren't any, it launch the save function
func AutoSave() {
	var exists bool
	row, err := get("push_schedules_today_exists")
	if err != nil {
		log.Fatal("Failed to get push schedules", "error", err)
	}

	err = row.Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check if push schedules exist", "error", err)
	}

	if !exists {
		save()
	}
}

//////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////
