// cmd package provides utility functions for working with commands.
package cmd

import (
	"aio/pkg/utils/fs"
	"errors"
	"os/exec"
	"runtime"
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

// StartBinaryWithInhibitSystemSleep function starts a binary with an inhibit system sleep.
// it is used to start a binary with an inhibit system sleep.
// it prevents the system from going to sleep while the binary is running.
// it returns an error if the operating system is not supported.
// it returns an error if the command fails to start.
func StartBinaryWithInhibitSystemSleep(binPath string) error {
	var cmd *exec.Cmd

	// check the operating system and run the appropriate command
	switch runtime.GOOS {
	case "darwin":
		// -i: prevent idle sleep -s: prevent system sleep
		cmd = exec.Command("caffeinate", "-i", "-s", binPath)
	case "linux":
		// --why: reason for inhibition --mode: block the system sleep
		cmd = exec.Command("systemd-inhibit", "--why=Prevent sleep", "--mode=block", binPath)
	case "windows":
		// powershell command to prevent system sleep
		powershellCmd := `
        Add-Type -TypeDefinition 'using System; using System.Runtime.InteropServices; public class Sleep { 
        [DllImport("kernel32.dll", CharSet = CharSet.Auto, SetLastError = true)] 
        public static extern uint SetThreadExecutionState(uint esFlags); }'; 
        [Sleep]::SetThreadExecutionState(0x80000002); 
        Start-Process -FilePath '` + binPath + `' -NoNewWindow -Wait
        `

		// run the powershell command
		cmd = exec.Command("powershell", "-Command", powershellCmd)
	default:
		// at this moment, only darwin, linux, and windows are supported
		// if the operating system is not supported, return an error
		return errors.New("unsupported operating system: " + runtime.GOOS)
	}

	// start the command
	// return an error if the command fails to start
	err := cmd.Run()
	if err != nil {
		return errors.New("failed to start binary with inhibit system sleep: " + err.Error())
	}

	return nil
}
