// cmd package is used to execute commands
package cmd

import (
	"aio/pkg/db"
	"aio/pkg/git"
	"aio/pkg/log"

	"github.com/spf13/cobra"
)

var rootLongDesc = `
aio (All In One) is a multi-purpose CLI app that manages tasks, notes, finances, 
health, and productivity, transforming life into a video game. 
Track progress, achieve goals, and level up in all aspects of your journey—making productivity fun and rewarding.
It is a fun way to keep track of your life and improve yourself.

For more information, visit the project page at https://github.com/Tagliapietra96/aio`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aio",
	Short: "A all in one application",
	Long:  rootLongDesc,
	// init the database and user if they do not exist
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := db.Init()
		if err != nil {
			log.Err("failed to initialize the database")
			log.Fat(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		revert, err := cmd.Flags().GetBool("revert")
		if err != nil {
			log.Err("failed to get flag revert")
			log.Fat(err)
		}

		if revert {
			err := git.Revert()
			if err != nil {
				log.Err("failed to revert the db version")
				log.Fat(err)
			}
		}

		addRemote, err := cmd.Flags().GetBool("link-remote")
		if err != nil {
			log.Err("failed to get flag link-remote")
			log.Fat(err)
		}

		if addRemote {
			err := git.LinkRepo()
			if err != nil {
				log.Err("failed to link the remote repository")
				log.Fat(err)

			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Err("failed to execute the root command")
		log.Fat(err)
	}
}

func init() {
	rootCmd.Flags().BoolP("revert", "r", false, "Revert the db version")
	rootCmd.Flags().BoolP("link-remote", "l", false, "Add a remote repository")
}
