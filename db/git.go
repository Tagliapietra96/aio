// db git.go
package db

import (
	"aio/helpers"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
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
	if !hasmainCheck {
		_, err := cmdExec("git", "show-ref", "--verify", "refs/heads/main")
		hasmain = err == nil
		hasmainCheck = true
	}

	if !hasmain {
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
	if !haschangesCheck {
		dbfile := getPath("data.db")
		output, err := cmdExec("git", "status", "--porcelain", dbfile)
		if err != nil {
			log.Fatal("Failed to check changes", "output", output, "error", err)
		}

		haschanges = strings.TrimSpace(output) != ""
		haschangesCheck = true
	}

	if haschanges {

		dbfile := getPath("data.db")
		var output string
		var err error

		if !haslocalCheck {
			_, err := cmdExec("git", "rev-parse", "--verify", "HEAD")
			haslocal = err == nil // Return true if there are local commits
			haslocalCheck = true
		}

		if !hasRemoteCommits() && !haslocal {
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
}

// gitFlow function is a wrapper function that executes a series of git commands.
// it is used to execute a series of git commands in a transaction to ensure the integrity of the data.
// if an error occurs, it rolls back the changes.
func gitFlow(action func() error) error {
	wu.Add(1)
	// pull the remote repository (if it exists) before making changes
	gitPull()

	// create a new branch
	branch := time.Now().Format("20060102150405")
	output, err := cmdExec("git", "checkout", "-b", branch)
	if err != nil {
		log.Fatal("Failed to create branch", "output", output, "error", err)
	}

	// execute the actions on the database, if an error occurs, delete the branch and return the error
	err = action()
	if err != nil {
		gitMain()
		output, err = cmdExec("git", "branch", "-D", branch)
		if err != nil {
			log.Fatal("Failed to delete branch", "output", output, "error", err)
		}
		return err
	}

	// commit the changes to the database
	go func() {
		defer wu.Done()
		output, err := cmdExec("git", "add", getPath("data.db"))
		if err != nil {
			log.Fatal("Failed to add database file", "output", output, "error", err)
		}

		output, err = cmdExec("git", "commit", "-m", "changes-"+branch)
		if err != nil {
			log.Fatal("Failed to commit changes", "output", output, "error", err)
		}

		gitMain()
		output, err = cmdExec("git", "merge", branch, "--no-ff")
		if err != nil {
			log.Fatal("Failed to merge branches", "output", output, "error", err)
		}
	}()
	return nil
}

// revertDB function reverts the database to a previous version.
// it is used to revert the database to a previous version.
// it gets the commit hash of the version to revert to and reverts the database to that version.
// it also commits the changes.
func revertDB() {
	// get the commit hash
	output, err := cmdExec("git", "log", "--pretty=format:%h %ad %s", "--date=short", "data.db")
	if err != nil {
		log.Fatal("Failed to get log", "output", output, "error", err)
	}
	history := strings.Split(output, "\n")
	fmt.Println("Select the version to revert to:\n")
	choice := helpers.RunSelect(history)
	fmt.Println("")
	commitHash := strings.Split(choice, " ")[0]
	version := strings.Split(choice, " ")[2]

	// do a back up of the database
	backup()

	// revert the database
	output, err = cmdExec("git", "checkout", commitHash, "--", "data.db")
	if err != nil {
		log.Fatal("Failed to revert database", "output", output, "error", err)
	}

	// commit the changes
	output, err = cmdExec("git", "add", "data.db")
	if err != nil {
		log.Fatal("Failed to add database file", "output", output, "error", err)
	}

	output, err = cmdExec("git", "commit", "-m", "revert-database-to-"+version)
	if err != nil {
		log.Fatal("Failed to commit changes", "output", output, "error", err)
	}

	log.Info("Database reverted to version " + version)
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
		if err != nil {
			log.Fatal("Failed to save new push schedule log", "error", err)
		}

		wu.Wait()
		// check if the remote branch exists
		output, err := cmdExec("git", "ls-remote", "--heads", "origin", "main")
		if err != nil {
			log.Fatal("Failed to check remote branch existence", "output", output, "error", err)
		}

		if strings.TrimSpace(output) != "" {
			// check if there are commits to push
			output, err = cmdExec("git", "rev-list", "--count", "origin/main..HEAD")
			if err != nil {
				log.Fatal("Failed to check for unpushed commits", "error", err)
			}

			// Convert the output to an integer
			commitsToPush, err := strconv.Atoi(strings.TrimSpace(output))
			if err != nil {
				log.Fatal("Failed to parse git rev-list output", "output", output, "error", err)
			}

			if commitsToPush < 1 {
				return
			}
		}

		// push the changes to the remote repository if it exists and if there are commits to push
		log.Info("Pushing changes to the remote repository...")
		output, err = cmdExec("git", "push", "-u", "origin", "main")
		if err != nil {
			log.Fatal("Failed to push changes", "output", output, "error", err)
		}
		hasremote = true
	}
}
