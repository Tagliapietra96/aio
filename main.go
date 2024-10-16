////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// File Name: main.go
// Created by: Matteo Tagliapietra 2024-09-01

// This is the main entry point for the application.
// It initializes the database and checks if the user exists.
// If the user does not exist, it initializes the user.

// App Version: 0.1.1

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// main package is the entry point for the application.
package main

import (
	"aio/cmd"
	"aio/logger"
)

// main function is the entry point for the application.
func main() {
	defer logger.Close()                                  // close the logger when the main function exits
	cmd.Execute()                                         // execute the commands
	logger.Debug("All processes completed successfully!") // notify the user that all processes have completed
	logger.Debug("-----------------------------")         // print a separator
}
