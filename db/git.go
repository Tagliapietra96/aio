// db git.go
package db

import (
	"aio/helpers"
	"aio/logger"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
)

// linkRepo function adds a link to a remote repository.
// it is used to link the database to a remote repository for versioning.
// this action is optional and can be skipped by the user and performed later.
func linkRepo() {
	if helpers.RunConfirm("Do you want to add a link to a remote repository?") {
		logger.Line("Please enter the remote repository name:")
		remote := helpers.RunInput("YourUsername/repo-name")
		output, err := cmdExec("git", "remote", "add", "origin", "git@github.com:"+remote)
		logger.Fatal("Failed to add remote repository", err, "output", output)
	}
}

// remoteExists function checks if a remote repository is linked to the database.
// it is used to check if a remote repository is linked to the database for versioning.
func remoteExists() bool {
	if !remoteexixtsCheck {
		output, err := cmdExec("git", "remote")
		logger.Fatal("Failed to get remote repository", err, "output", output)
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

// gitPull function pulls the remote repository.
// it is used to pull the remote repository before making changes to the database.
func gitPull() {
	if remoteExists() && hasRemoteCommits() && !isAligned {
		gitMain()
		output, err := cmdExec("git", "pull", "origin", "main")
		logger.Fatal("Failed to pull remote repository", err, "output", output)
		isAligned = true
	}
}

// gitMain function checks out the main branch.
// it is used to checkout the main branch before making changes to the database.
func gitMain() {
	if !hasmainCheck {
		_, err := cmdExec("git", "show-ref", "--verify", "refs/heads/main")
		hasmain = err == nil
		hasmainCheck = true
	}

	if !hasmain {
		output, err := cmdExec("git", "checkout", "-b", "main")
		logger.Fatal("Failed to create main branch", err, "output", output)
		hasmain = true
	} else {
		output, err := cmdExec("git", "checkout", "main")
		logger.Fatal("Failed to checkout main branch", err, "output", output)
	}
}

// initialCommit function commits the database file.
// it is used to commit the database file to the local repository.
// it also renames the branch to main and pushes the changes to the remote repository if it exists.
func initialCommit() {
	dbfile := getPath("data.db")

	if !haschangesCheck {
		output, err := cmdExec("git", "status", "--porcelain", dbfile)
		logger.Fatal("Failed to check changes", err, "output", output)

		haschanges = strings.TrimSpace(output) != ""
		haschangesCheck = true
	}

	if haschanges {
		var output string
		var err error

		if !haslocalCheck {
			_, err := cmdExec("git", "rev-parse", "--verify", "HEAD")
			haslocal = err == nil // Return true if there are local commits
			haslocalCheck = true
		}

		if !hasRemoteCommits() && !haslocal {
			output, err = cmdExec("git", "add", dbfile)
			logger.Fatal("Failed to add database file", err, "output", output)

			// do the initial commit
			output, err = cmdExec("git", "commit", "-m", "initial commit")
			logger.Fatal("Failed to commit changes", err, "output", output)
			output, err = cmdExec("git", "branch", "-M", "main")
			logger.Fatal("Failed to rename branch", err, "output", output)

			haslocal = true
		}
	}
}

// gitFlow function is a wrapper function that executes a series of git commands.
// it is used to execute a series of git commands in a transaction to ensure the integrity of the data.
// if an error occurs, it rolls back the changes.
func gitFlow(action func() error) error {
	// pull the remote repository (if it exists) before making changes
	gitPull()

	// create a new branch
	branch := time.Now().Format("20060102150405")
	output, err := cmdExec("git", "checkout", "-b", branch)
	logger.Fatal("Failed to create branch", err, "output", output)

	// execute the actions on the database, if an error occurs, delete the branch and return the error
	err = action()
	if err != nil {
		gitMain()
		output, err = cmdExec("git", "branch", "-D", branch)
		logger.Fatal("Failed to delete branch", err, "output", output)
		return err
	}

	// commit the changes to the database
	output, err = cmdExec("git", "add", getPath("data.db"))
	logger.Fatal("Failed to add database file", err, "output", output)

	output, err = cmdExec("git", "commit", "-m", "changes-"+branch)
	logger.Fatal("Failed to commit changes", err, "output", output)

	gitMain()
	output, err = cmdExec("git", "merge", branch, "--no-ff")
	logger.Fatal("Failed to merge branches", err, "output", output)
	return nil
}

// revertDB function reverts the database to a previous version.
// it is used to revert the database to a previous version.
// it gets the commit hash of the version to revert to and reverts the database to that version.
// it also commits the changes.
func revertDB() {
	// get the commit hash
	output, err := cmdExec("git", "log", "--pretty=format:%h %ad %s", "--date=short", "data.db")
	logger.Fatal("Failed to get log", err, "output", output)
	history := strings.Split(output, "\n")
	logger.Line("Select the version to revert to:\n")
	choice := helpers.RunSelect(history)
	logger.Line("")
	commitHash := strings.Split(choice, " ")[0]
	version := strings.Split(choice, " ")[2]

	// do a back up of the database
	backup()

	// revert the database
	output, err = cmdExec("git", "checkout", commitHash, "--", "data.db")
	logger.Fatal("Failed to revert database", err, "output", output)

	// commit the changes
	output, err = cmdExec("git", "add", "data.db")
	logger.Fatal("Failed to add database file", err, "output", output)

	output, err = cmdExec("git", "commit", "-m", "revert-database-to-"+version)
	logger.Fatal("Failed to commit changes", err, "output", output)

	logger.Info("Database reverted to version " + version)
}

// save function saves the changes made to the database.
// it is used to save the changes made to the database and push them to the remote repository if it exists.
// it also creates a new branch for the changes.
func save() {
	if remoteExists() {
		// add the new push schedule log
		err := gitFlow(func() error {
			return do("push_schedules_create")
		})
		logger.Fatal("Failed to save new push schedule log", err)

		// check if the remote branch exists
		output, err := cmdExec("git", "ls-remote", "--heads", "origin", "main")
		logger.Fatal("Failed to check remote branch existence", err, "output", output)

		if strings.TrimSpace(output) != "" {
			// check if there are commits to push
			output, err = cmdExec("git", "rev-list", "--count", "origin/main..HEAD")
			logger.Fatal("Failed to check for unpushed commits", err)

			// Convert the output to an integer
			commitsToPush, err := strconv.Atoi(strings.TrimSpace(output))
			logger.Fatal("Failed to parse git rev-list output", err, "output", output)

			if commitsToPush < 1 {
				return
			}
		}

		// push the changes to the remote repository if it exists and if there are commits to push
		logger.Debug("Pushing changes to the remote repository...")
		output, err = cmdExec("git", "push", "-u", "origin", "main")
		logger.Fatal("Failed to push changes", err, "output", output)
		hasremote = true
		logger.Info("Changes pushed to the remote repository")
	}
}
