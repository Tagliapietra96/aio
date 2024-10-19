// cmd package provides utility functions for working with commands.
package cmd

import (
	"aio/pkg/utils/fs"
	"os/exec"
)

// Output function executes a command and returns the output.
// it is used to execute a command and return the output.
// it also sets the working directory to the executable directory.
func Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	dir, err := fs.ExecDir()
	if err != nil {
		return nil, err
	}

	cmd.Dir = dir
	output, err := cmd.Output()
	return output, err
}

// Start function starts a command.
// it is used to start a command.
// it also sets the working directory to the executable directory.
func Start(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	dir, err := fs.ExecDir()
	if err != nil {
		return err
	}

	cmd.Dir = dir
	return cmd.Start()
}
