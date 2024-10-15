////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// File Name: main.go
// Created by: Matteo Tagliapietra 2024-09-01
// Last Update: 2024-10-15

// This is the main entry point for the application.
// It initializes the database and checks if the user exists.
// If the user does not exist, it initializes the user.

// Version: 0.1.0

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// main package is the entry point for the application.
package main

// imports the necessary packages
// cmd package is used to execute commands
import (
	"aio/cmd"
)

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// main function is the entry point for the application.
func main() {
	// execute commands
	cmd.Execute()
}
