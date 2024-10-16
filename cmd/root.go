// cmd package is used to execute commands
package cmd

import (
	"aio/db"
	"aio/logger"

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
		db.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		revert, err := cmd.Flags().GetBool("revert")
		logger.Fatal("Failed to get flag revert", err)

		if revert {
			db.Revert()
		}
	},
	// push the changes on repository
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		db.AutoSave()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	logger.Fatal("Failed to execute command", err)
}

func init() {
	rootCmd.Flags().BoolP("revert", "r", false, "Revert the db version")
}
