// db package is used to interact with the database.
package db

import (
	"aio/pkg/git"
	"aio/pkg/inputs"
	"aio/pkg/log"
	"aio/pkg/utils/cmd"
	"aio/pkg/utils/fs"
	"aio/pkg/utils/num"
	"aio/pkg/utils/tm"
	"database/sql"
	"embed"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed queries/*.sql
var sqlFiles embed.FS

// loadQuery function reads the content of a sql file and returns it as a string.
// it is used to load the content of the sql files that contain the queries to execute.
func loadQuery(filename string) (string, error) {
	query, err := sqlFiles.ReadFile("queries/" + filename + ".sql")
	if err != nil {
		log.Err("failed to read query file")
	}

	return string(query), err
}

// getDb function returns a pointer to a sql.DB object.
// it is used to open the database file and return a pointer to the database object.
// if the database file does not exist, it creates a new one.
func getDb() (*sql.DB, error) {
	// get the database file path
	dbfile, err := fs.DBfile()
	if err != nil {
		log.Err("failed to get database file path")
		return nil, err
	}

	// open the database
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Err("failed to open database")
		return nil, err
	}

	return db, nil
}

// do function executes a query on the database.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func do(query string, args ...any) error {
	// open the database
	db, err := getDb()
	if err != nil {
		log.Err("failed to open database")
		return err
	}

	defer db.Close()

	// load the query
	q, err := loadQuery(query)
	if err != nil {
		log.Err("failed to load query")
		return err
	}

	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Err("failed to start transaction")
		return err
	}

	// execute the query
	_, err = tx.Exec(q, args...)
	if err != nil {
		tx.Rollback()
		log.Err("failed to execute query")
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Err("failed to commit transaction")
		return err
	}

	return nil
}

// gets function executes a query on the database and returns the result.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func gets(query string, args ...any) (*sql.Rows, error) {
	// open the database
	db, err := getDb()
	if err != nil {
		log.Err("failed to open database")
		return nil, err
	}

	defer db.Close()

	// load the query
	q, err := loadQuery(query)
	if err != nil {
		log.Err("failed to load query")
		return nil, err
	}

	// execute the query
	rows, err := db.Query(q, args...)
	if err != nil {
		log.Err("failed to execute query")
		return nil, err
	}

	return rows, nil
}

// get function executes a query on the database and returns the result.
// the execution is done in a transaction to ensure the integrity of the data.
// if the transaction fails, it rolls back the changes.
// every step is logged in case of errors and stop the execution
func get(query string, args ...any) (*sql.Row, error) {
	// open the database
	db, err := getDb()
	if err != nil {
		log.Err("failed to open database")
		return nil, err
	}

	defer db.Close()

	// load the query
	q, err := loadQuery(query)
	if err != nil {
		log.Err("failed to load query")
		return nil, err
	}

	// execute the query
	row := db.QueryRow(q, args...)
	return row, nil
}

// Init function initializes the database.
// the funciton initialize also git for the db versioning
// it is used to create the database file and tables if they do not exist.
// it loads the main sql file that contains the queries to create the tables, indexes, and triggers.
func Init() error {
	git.Init()

	log.Deb("initializing database...")

	// if the database file does not exist, create a new one
	dbfile, err := fs.Path("data.db")
	if err != nil {
		log.Err("failed to get database file path")
		return err
	}

	_, err = os.Stat(dbfile)
	if os.IsNotExist(err) {
		log.Warn("database file not found, creating a new one", "file", dbfile)

		var file *os.File
		file, err = os.Create(dbfile)
		if err != nil {
			log.Err("failed to create database file")
			return err
		}

		file.Close()
		log.Info("database file created", "file", dbfile)
	}

	// create the tables, indexes, and triggers
	err = do("tables")
	if err != nil {
		log.Err("failed to create tables")
		return err
	}

	log.Info("database initialized successfully!")

	// do the initial commit
	git.InitialCommit()

	// check if there are characters in the database
	var exists bool
	log.Deb("initializing character...")
	row, err := get("characters_exists")
	if err != nil {
		log.Err("failed to check if characters exist")
		return err
	}

	err = row.Scan(&exists)
	if err != nil {
		log.Err("failed to check if characters exist")
		return err
	}

	// if characters do not exist, create the initial character
	if !exists {
		log.Warn("no characters found in the database")
		log.Deb("creating initial character...")
		log.PrintS("Welcome to AIO - Your Life, Gamified!", log.TitleStyle)
		log.Print(`
Turn your tasks, goals, and habits into an adventure.
Track progress, manage your finances, boost productivity, and level up in all aspects of life.
Ready to make self-improvement fun? Your journey starts now!
`)
		log.Print("Welcome, traveler! Before we begin your adventure, we need to know your name.")
		log.Print("What is your first name, brave soul?")
		fn := inputs.RunInput("Jhon")

		log.Print("A strong name indeed! Now, please tell us your family name, the one that will echo through the halls of history.")
		log.Print("What is your last name, worthy adventurer?")
		ln := inputs.RunInput("Smith")

		log.Print("Every hero has a title that the bards will sing of! Choose a nickname, one that will strike fear into your foes or inspire your allies.")
		log.Print("What shall your unique nickname be?")
		nn := inputs.RunInput("The Reaper")

		log.Print("Even legends have a beginning. We need to know when your story began.")
		log.Print("Please provide your date of birth in the form of 02 Jan 2006.")
		log.Print("When were you born, chosen one?")
		dob := inputs.RunInputWithValidation("02 Jan 2006", tm.ValidateDate)

		log.Print("Every great adventurer must wisely manage their resources, not just in battle, but also in life.")
		log.Print("Set your monthly budget this will guide how you manage your gold throughout your journey!")
		log.Print(("How much gold will you allocate each month for your expenses? (Enter a numeric value)"))
		b := inputs.RunInputWithValidation("1500.00", num.Validate)

		budget, err := num.ParseFloat(b)
		if err != nil {
			log.Err("failed to parse budget")
			return err
		}

		birth, err := tm.DBReformat(dob)
		if err != nil {
			log.Err("failed to reformat birth date")
			return err
		}

		err = do("characters_create", fn, ln, nn, birth, budget)
		if err != nil {
			log.Err("failed to create character")
			return err
		}

		log.Print("🎉 Your character has been created! 🎉")
		log.Print("Welcome, %s %s, also known as %s!", fn, ln, nn)
		log.Print("Now, go forth and conquer the challenges ahead!")
		log.Print("\nTo view an help text use the command 'aio --help', or 'aio -h'.\n")
	}
	log.Info("character initialized successfully!")

	// check if there are daily logins in the database for today
	row, err = get("daily_logins_today_exists")
	if err != nil {
		log.Err("failed to check if daily logins exist")
		return err
	}

	err = row.Scan(&exists)
	if err != nil {
		log.Err("failed to check if daily logins exist")
		return err
	}

	// if daily logins do not exist, create the daily login and update the character stats
	if !exists {
		log.Warn("no daily logins found for today")
		log.Deb("creating daily login...")
		err = do("characters_daily_login")
		if err != nil {
			log.Err("failed to create daily login")
			return err
		}
		log.Info("daily login created successfully!")
	}

	// launch the cron service
	bin, err := fs.Path("cron")
	if err != nil {
		log.Err("failed to get cron binary path")
		return err
	}

	running := false
	out, err := cmd.Output("pgrep", "-f", bin)
	if err == nil {
		running = strings.TrimSpace(string(out)) != ""
	}

	if !running {
		err = os.Chmod(bin, 0755)
		if err != nil {
			log.Err("failed to change cron binary permissions")
			return err
		}

		err := cmd.Start("caffeinate", "-s", bin)
		if err != nil {
			log.Err("failed to start cron service")
			return err
		}

		log.Info("started cron service.")
	}

	log.Info("database initialized successfully!")
	return nil
}
