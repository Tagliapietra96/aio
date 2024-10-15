// db package is used to interact with the database.
package db

import (
	"aio/helpers"
	"aio/style"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

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

// getExecDir function returns the directory of the executable file.
// it is used to run all the commands in the directory of the executable file.
// this maintaning the integrity of the data.
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
	if err != nil {
		log.Fatal("Failed to open database file", "error", err)
	}
	defer input.Close()

	output, err := os.Create(backupfile)
	if err != nil {
		log.Fatal("Failed to create backup file", "error", err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		log.Fatal("Failed to copy database file", "error", err)
	}
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

		// if git is not initialized, link the repository
		linkRepo()
		if remoteExists() {
			log.Info("Fetching remote repository...")
			if remoteExists() && hasRemoteCommits() && !isAligned {
				output, err := cmdExec("git", "fetch")
				if err != nil {
					log.Fatal("Failed to fetch remote repository", "output", output, "error", err)
				}
			}
			gitPull()
		}

		// if the .gitignore file does not exist, create a new one
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

	// if the database file does not exist, create a new one
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

	// create the tables, indexes, and triggers
	err = do("tables")
	if err != nil {
		log.Fatal("Failed to create tables", "error", err)
	}

	// check if there are characters in the database
	var exists bool
	row, err := get("characters_exists")
	if err != nil {
		log.Fatal("Failed to check if characters exist", "error", err)
	}

	err = row.Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check if characters exist", "error", err)
	}

	// if characters do not exist, create the initial character
	if !exists {
		style.PrintTitle("\nWelcome to AIO - Your Life, Gamified!")
		style.Print(`
Turn your tasks, goals, and habits into an adventure.
Track progress, manage your finances, boost productivity, and level up in all aspects of life.
Ready to make self-improvement fun? Your journey starts now!
`)
		style.Print("Welcome, traveler! Before we begin your adventure, we need to know your name.")
		style.Print("What is your first name, brave soul?")
		fn := helpers.RunInput("Jhon")

		style.Print("A strong name indeed! Now, please tell us your family name, the one that will echo through the halls of history.")
		style.Print("What is your last name, worthy adventurer?")
		ln := helpers.RunInput("Smith")

		style.Print("Every hero has a title that the bards will sing of! Choose a nickname, one that will strike fear into your foes or inspire your allies.")
		style.Print("What shall your unique nickname be?")
		nn := helpers.RunInput("The Reaper")

		style.Print("Even legends have a beginning. We need to know when your story began.")
		style.Print("Please provide your date of birth in the form of 02 Jan 06 (day month year).")
		style.Print("When were you born, chosen one?")
		dob := helpers.RunInputWithValidation("02 Jan 06", helpers.TimeValidate)

		style.Print("Every great adventurer must wisely manage their resources, not just in battle, but also in life.")
		style.Print("Set your monthly budget—this will guide how you manage your gold throughout your journey!")
		style.Print(("How much gold will you allocate each month for your expenses? (Enter a numeric value)"))
		b := helpers.RunInputWithValidation("1500.00", helpers.NumberValidate)

		err := gitFlow(func() error {
			return do("characters_create", fn, ln, nn, helpers.TimeDBReformat(dob), helpers.NumberParse(b))
		})
		if err != nil {
			log.Fatal("Failed to create character", "error", err)
		}

		style.Print("🎉 Your character has been created! 🎉")
		style.Print("Welcome, " + fn + " " + ln + ", also known as " + nn + "!")
		style.Print("Now, go forth and conquer the challenges ahead!")
		style.Print("\nTo view an help text use the command 'aio --help', or 'aio -h'.\n")
	}

	// check if there are daily logins in the database for today
	row, err = get("daily_logins_today_exists")
	if err != nil {
		log.Fatal("Failed to check if daily logins exist", "error", err)
	}

	err = row.Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check if daily logins exist", "error", err)
	}

	// if daily logins do not exist, create the daily login and update the character stats
	if !exists {
		err = gitFlow(func() error { return do("characters_daily_login") })
		if err != nil {
			log.Fatal("Failed to create daily login", "error", err)
		}
	}

	initialCommit()
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

// Revert function reverts the database to a previous version.
func Revert() {
	confirm := helpers.RunConfirm("Are you sure you want to revert the database?")
	if !confirm {
		return
	}
	revertDB()
}
