// cmd package, completion command file
package cmd

import (
	"aio/pkg/log"
	"errors"

	"github.com/spf13/cobra"
)

const completionLongDesc = `
Completion (aio completion [bash|zsh|fish|powershell]) generates shell completion scripts for the specified shell.
The shell code must be evaluated to provide interactive completion of commands.
This can be done by sourcing it.

To load completions:

# Bash:
  source <(aio completion bash)

# Zsh:
  source <(aio completion zsh)

# Fish:
  aio completion fish | source

# Powershell:
  aio completion powershell | Out-String | Invoke-Expression
`

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Args:  cobra.ExactArgs(1),
	Short: "Generate completion scripts for your shell",
	Long:  completionLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			err := rootCmd.GenBashCompletion(cmd.OutOrStdout())
			if err != nil {
				log.Err("failed to generate bash completion")
				log.Fat(err)
			}
		case "zsh":
			err := rootCmd.GenZshCompletion(cmd.OutOrStdout())
			if err != nil {
				log.Err("failed to generate zsh completion")
				log.Fat(err)
			}
		case "fish":
			err := rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
			if err != nil {
				log.Err("failed to generate fish completion")
				log.Fat(err)
			}
		case "powershell":
			err := rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
			if err != nil {
				log.Err("failed to generate powershell completion")
				log.Fat(err)
			}
		default:
			log.Err("unsupported shell type")
			log.Fat(errors.New("unsupported shell type, use bash, zsh, fish, or powershell as the argument"))
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
