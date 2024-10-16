// db package is used to interact with the database.
package db

import (
	"aio/helpers"
	"aio/logger"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Init function initializes the database.
// the funciton initialize also git for the db versioning
// it is used to create the database file and tables if they do not exist.
// it loads the main sql file that contains the queries to create the tables, indexes, and triggers.
func Init() {
	// check if git is initialized
	if _, err := os.Stat(getPath(".git")); os.IsNotExist(err) {
		logger.Debug("Initializing git...")
		output, err := cmdExec("git", "init")
		logger.Fatal("Failed to initialize git", err, "output", output)

		// if git is not initialized, link the repository
		linkRepo()
		if remoteExists() {
			logger.Debug("Fetching remote repository...")
			if remoteExists() && hasRemoteCommits() && !isAligned {
				output, err := cmdExec("git", "fetch")
				logger.Fatal("Failed to fetch remote repository", err, "output", output)
			}
			gitPull()
		}

		// if the .gitignore file does not exist, create a new one
		gitIgnorePath := getPath(".gitignore")
		if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
			logger.Warn(".gitignore file not found, creating a new one", "file", gitIgnorePath)

			content := `*
!data.db
`
			err := os.WriteFile(gitIgnorePath, []byte(content), 0644)
			logger.Fatal("Failed to create .gitignore file", err)
			logger.Info(".gitignore file created", "file", gitIgnorePath)
		}

		logger.Info("Git initialized successfully!")
	}

	logger.Debug("Initializing database...")

	// if the database file does not exist, create a new one
	dbfile := getPath("data.db")
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		logger.Warn("Database file not found, creating a new one", "file", dbfile)

		var file *os.File
		file, err = os.Create(dbfile)
		logger.Fatal("Failed to create database file", err)

		file.Close()
		logger.Info("Database file created", "file", dbfile)
	}

	// create the tables, indexes, and triggers
	err = do("tables")
	logger.Fatal("Failed to create tables", err)
	logger.Info("Database initialized successfully!")

	// do the initial commit
	initialCommit()

	// check if there are characters in the database
	var exists bool
	logger.Debug("Initializing character...")
	row, err := get("characters_exists")
	logger.Fatal("Failed to check if characters exist", err)

	err = row.Scan(&exists)
	logger.Fatal("Failed to check if characters exist", err)

	// if characters do not exist, create the initial character
	if !exists {
		logger.Warn("No characters found in the database")
		logger.Debug("Creating initial character...")
		logger.Print("Welcome to AIO - Your Life, Gamified!", logger.TitleStyle)
		logger.Line(`
Turn your tasks, goals, and habits into an adventure.
Track progress, manage your finances, boost productivity, and level up in all aspects of life.
Ready to make self-improvement fun? Your journey starts now!
`)
		logger.Line("Welcome, traveler! Before we begin your adventure, we need to know your name.")
		logger.Line("What is your first name, brave soul?")
		fn := helpers.RunInput("Jhon")

		logger.Line("A strong name indeed! Now, please tell us your family name, the one that will echo through the halls of history.")
		logger.Line("What is your last name, worthy adventurer?")
		ln := helpers.RunInput("Smith")

		logger.Line("Every hero has a title that the bards will sing of! Choose a nickname, one that will strike fear into your foes or inspire your allies.")
		logger.Line("What shall your unique nickname be?")
		nn := helpers.RunInput("The Reaper")

		logger.Line("Even legends have a beginning. We need to know when your story began.")
		logger.Line("Please provide your date of birth in the form of 02 Jan 2006.")
		logger.Line("When were you born, chosen one?")
		dob := helpers.RunInputWithValidation("02 Jan 2006", helpers.TimeValidate)

		logger.Line("Every great adventurer must wisely manage their resources, not just in battle, but also in life.")
		logger.Line("Set your monthly budget—this will guide how you manage your gold throughout your journey!")
		logger.Line(("How much gold will you allocate each month for your expenses? (Enter a numeric value)"))
		b := helpers.RunInputWithValidation("1500.00", helpers.NumberValidate)

		err := do("characters_create", fn, ln, nn, helpers.TimeDBReformat(dob), helpers.NumberParse(b))
		logger.Fatal("Failed to create character", err)

		logger.Line("🎉 Your character has been created! 🎉")
		logger.Line("Welcome, %s %s, also known as %s!", fn, ln, nn)
		logger.Line("Now, go forth and conquer the challenges ahead!")
		logger.Line("\nTo view an help text use the command 'aio --help', or 'aio -h'.\n")
	}
	logger.Info("Character initialized successfully!")

	// check if there are daily logins in the database for today
	row, err = get("daily_logins_today_exists")
	logger.Fatal("Failed to check if daily logins exist", err)

	err = row.Scan(&exists)
	logger.Fatal("Failed to check if daily logins exist", err)

	// if daily logins do not exist, create the daily login and update the character stats
	if !exists {
		logger.Warn("No daily logins found for today")
		logger.Debug("Creating daily login...")
		err = gitFlow(func() error { return do("characters_daily_login") })
		logger.Fatal("Failed to create daily login", err)
		logger.Info("Daily login created successfully!")
	}
}

// AutoSave function saves the changes made to the database automatically.
// check if there are push schedules in the database that have the created_at field se to the current day
// if there aren't any, it launch the save function
func AutoSave() {
	var exists bool
	row, err := get("push_schedules_today_exists")
	logger.Fatal("Failed to get push schedules", err)
	err = row.Scan(&exists)
	logger.Fatal("Failed to check if push schedules exist", err)

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
