// git package provides utility functions for working with git.
package git

import (
	"aio/pkg/inputs"
	"aio/pkg/log"
	"aio/pkg/utils/cmd"
	"aio/pkg/utils/fs"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const gitignore = `*
!data.db
`

// remoteExists function checks if a remote repository is linked to the database.
// it is used to check if a remote repository is linked to the database for versioning.
func remoteExists() (bool, error) {
	output, err := cmd.Output("git", "remote")
	if err != nil {
		log.Err("failed to check remote repository", "output", string(output))
		return false, err
	}
	return strings.TrimSpace(string(output)) != "", nil // return true if a remote repository is linked
}

// hasLocalCommits function checks if there are local commits.
// it is used to check if there are local commits before making changes to the database.
func hasLocalCommits() bool {
	_, err := cmd.Output("git", "rev-parse", "--verify", "HEAD")
	return err == nil // return true if there are local commits
}

// hasRemoteCommits function checks if there are remote commits.
// it is used to check if there are remote commits before making changes to the database.
func hasRemoteCommits() bool {
	output, err := cmd.Output("git", "ls-remote", "--heads", "origin")
	return err == nil && len(output) > 0 // return true if there are remote commits
}

// hasChanges function checks if there are changes to the database.
// it is used to check if there are changes to the database before making changes.
func hasChanges() (bool, error) {
	dbfile, err := fs.DBfile() // Get the database file path
	if err != nil {
		log.Err("failed to get database file path")
		return false, err
	}

	output, err := cmd.Output("git", "status", "--porcelain", dbfile) // Check the status of the database file
	if err != nil {
		log.Err("failed to check if database has been changed", "output", string(output))
		return false, err
	}

	return strings.TrimSpace(string(output)) != "", nil // Return true if there are changes to the database
}

// linkRepo function adds a link to a remote repository.
// it is used to link the database to a remote repository for versioning.
// this action is optional and can be skipped by the user and performed later.
func LinkRepo() error {
	log.Deb("checking if remote repository exists...")
	re, err := remoteExists() // Check if a remote repository is linked to the database
	if err != nil {
		log.Err("failed to check if remote repository exists")
		return err
	}

	if !re {
		log.Deb("linking remote repository...")
		if inputs.RunConfirm("Do you want to add a link to a remote repository?") {
			log.Print("Please enter the remote repository name:")                                 // Ask the user to enter the remote repository name
			remote := inputs.RunInput("YourUsername/repo-name")                                   // Get the remote repository name from the user
			output, err := cmd.Output("git", "remote", "add", "origin", "git@github.com:"+remote) // Add the remote repository
			if err != nil {
				log.Err("failed to add remote repository", "output", string(output))
				return err
			}

			log.Info("remote repository added successfully!", "repository", remote)
			return nil
		}
	}

	log.Warn("remote repository already exists")
	return nil // Return nil if the remote repository is already linked
}

// InitialCommit function commits the database file.
// it is used to commit the database file to the local repository.
// it also renames the branch to main and pushes the changes to the remote repository if it exists.
func InitialCommit() error {
	dbfile, err := fs.DBfile() // Get the database file path
	if err != nil {
		log.Err("failed to get database file path")
		return err
	}

	log.Deb("checking if database has been changed...")
	ch, err := hasChanges() // Check if there are changes to the database
	if err != nil {
		log.Err("failed to check if database has been changed")
		return err
	}

	if ch {
		log.Deb("checking if there are local commits...")
		if !hasRemoteCommits() && !hasLocalCommits() {
			// add the database file
			output, err := cmd.Output("git", "add", dbfile)
			if err != nil {
				log.Err("failed to add database file", "output", string(output))
				return err
			}

			// do the initial commit
			output, err = cmd.Output("git", "commit", "-m", "initial commit")
			if err != nil {
				log.Err("failed to commit database file", "output", string(output))
				return err
			}

			// rename the branch to main
			output, err = cmd.Output("git", "branch", "-M", "main")
			if err != nil {
				log.Err("failed to rename branch to main", "output", string(output))
				return err
			}

			log.Info("database committed successfully!")
			return nil
		}
	}

	log.Warn("no changes to commit")
	return nil
}

func Init() error {
	log.Deb("checking if git is initialized...")

	// get the path to the .git directory
	git, err := fs.Path(".git")
	if err != nil {
		log.Err("failed to get .git directory path")
		return err
	}

	// check if git is initialized
	if _, err := os.Stat(git); os.IsNotExist(err) {
		log.Deb("initializing git...")
		output, err := cmd.Output("git", "init")
		if err != nil {
			log.Err("failed to initialize git", "output", string(output))
			return err
		}

		// if git is not initialized, link the repository
		LinkRepo()

		re, err := remoteExists() // Check if a remote repository is linked to the database
		if err != nil {
			log.Err("failed to check if remote repository exists")
			return err
		}

		if re {
			log.Deb("fetching remote repository...")
			if re && hasRemoteCommits() {
				output, err := cmd.Output("git", "fetch")
				if err != nil {
					log.Err("failed to fetch remote repository", "output", string(output))
					return err
				}
			}
		}

		log.Deb("checking if .gitignore file exists...")

		// if the .gitignore file does not exist, create a new one
		gitIgnorePath, err := fs.Path(".gitignore")
		if err != nil {
			log.Err("failed to get .gitignore file path")
			return err
		}

		if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
			log.Warn(".gitignore file not found, creating a new one", "file", gitIgnorePath)
			err := os.WriteFile(gitIgnorePath, []byte(gitignore), 0644)
			if err != nil {
				return errors.New("failed to create .gitignore file: " + err.Error())
			}

			log.Info(".gitignore file created", "file", gitIgnorePath)
		}

		log.Info("git initialized successfully!")
	}

	err = Pull()
	return err
}

// Main function checks out the main branch.
// it is used to checkout the main branch before making changes to the database.
func Main() error {
	rerr := func(err error, output []byte) error {
		if err != nil {
			log.Err("failed to checkout main branch", "output", string(output))
			return err
		}
		return nil
	}

	_, err := cmd.Output("git", "show-ref", "--verify", "refs/heads/main")
	if err != nil {
		output, err := cmd.Output("git", "checkout", "-b", "main")
		return rerr(err, output)
	} else {
		output, err := cmd.Output("git", "checkout", "main")
		return rerr(err, output)
	}
}

// Pull function pulls the remote repository.
// it is used to pull the remote repository before making changes to the database.
func Pull() error {
	re, err := remoteExists() // Check if a remote repository is linked to the database
	if err != nil {
		log.Err("failed to check if remote repository exists")
		return err
	}

	if re && hasRemoteCommits() {
		Main()
		output, err := cmd.Output("git", "pull", "origin", "main")
		if err != nil {
			log.Err("failed to pull remote repository", "output", string(output))
			return err
		}

	}

	return nil
}

// Commit function commits the changes made to the database.
// it is used to commit the changes made to the database to the local repository.
func Commit() error {
	log.Deb("checking if database has been changed...")
	ch, err := hasChanges() // Check if there are changes to the database
	if err != nil {
		log.Err("failed to check if database has been changed")
		return err
	}

	if ch {
		dbfile, err := fs.DBfile() // Get the database file path
		if err != nil {
			log.Err("failed to get database file path")
			return err
		}

		// add the database file
		output, err := cmd.Output("git", "add", dbfile)
		if err != nil {
			log.Err("failed to add database file", "output", string(output))
			return err
		}

		// commit the changes
		output, err = cmd.Output("git", "commit", "-m", "changes-"+time.Now().Format("20060102150405"))
		if err != nil {
			log.Err("failed to commit database file", "output", string(output))
			return err
		}

		log.Info("database committed successfully!")
		return nil
	}

	log.Info("no changes to commit")
	return nil
}

func Push() error {
	re, err := remoteExists() // Check if a remote repository is linked to the database
	if err != nil {
		log.Err("failed to check if remote repository exists")
		return err
	}

	if re {
		// check if the remote branch exists
		output, err := cmd.Output("git", "ls-remote", "--heads", "origin", "main")
		if err != nil {
			log.Err("failed to check remote branch existence", "output", string(output))
			return err
		}

		if strings.TrimSpace(string(output)) != "" {
			log.Deb("checking if there are local commits...")

			// check if there are commits to push
			output, err = cmd.Output("git", "rev-list", "--count", "origin/main..HEAD")
			if err != nil {
				log.Err("failed to check for unpushed commits", "output", string(output))
				return err
			}

			// Convert the output to an integer
			commitsToPush, err := strconv.Atoi(strings.TrimSpace(string(output)))
			if err != nil {
				log.Err("failed to parse git rev-list output", "output", string(output))
				return err
			}

			if commitsToPush > 0 {
				// push the changes to the remote repository if it exists and if there are commits to push
				output, err = cmd.Output("git", "push", "-u", "origin", "main")
				if err != nil {
					log.Err("failed to push changes", "output", string(output))
					return err
				}
				log.Info("changes pushed successfully!")
				return nil
			}

			log.Info("no changes to push")
			return nil
		}
	}

	log.Warn("no remote repository linked")
	return nil
}

// Revert function reverts the database to a previous version.
// it is used to revert the database to a previous version.
// it gets the commit hash of the version to revert to and reverts the database to that version.
// it also commits the changes.
func Revert() error {
	confirm := inputs.RunConfirm("Are you sure you want to revert the database?")
	if !confirm {
		return nil
	}

	log.Deb("getting commit history...")

	// get the commit hash
	output, err := cmd.Output("git", "log", "--pretty=format:%h %ad %s", "--date=short", "data.db")
	if err != nil {
		log.Err("failed to get commit history", "output", string(output))
		return err
	}

	history := strings.Split(string(output), "\n")
	log.Print("Select the version to revert to:\n")
	choice := inputs.RunSelect(history)
	log.Print("")
	commitHash := strings.Split(choice, " ")[0]
	version := strings.Split(choice, " ")[2]

	log.Deb("do a back up of the database...")
	// do a back up of the database
	err = fs.Backup()
	if err != nil {
		log.Err("failed to back up the database")
		return err
	}

	log.Deb("reverting database to " + version + "...")
	// revert the database
	output, err = cmd.Output("git", "checkout", commitHash, "--", "data.db")
	if err != nil {
		log.Err("failed to revert database", "output", string(output))
		return err
	}

	log.Deb("committing changes...")
	// commit the changes
	output, err = cmd.Output("git", "add", "data.db")
	if err != nil {
		log.Err("failed to add database file", "output", string(output))
		return err
	}

	// commit the changes
	output, err = cmd.Output("git", "commit", "-m", "revert-database-to-"+version)
	if err != nil {
		log.Err("failed to commit database file", "output", string(output))
		return err
	}

	log.Info("database reverted successfully!")
	return nil
}
