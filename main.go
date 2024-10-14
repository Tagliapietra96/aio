////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// File Name: main.go
// Created by: Matteo Tagliapietra 2024-09-01
// Last Update: 2024-10-14

// This is the main entry point for the application.
// It initializes the database and checks if the user exists.
// If the user does not exist, it initializes the user.

// Version: 0.0.1

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// main package is the entry point for the application.
package main

// imports the necessary packages
// cmd package is used to execute commands
// db package is used to interact with the database
import (
	"aio/cmd"
	"aio/db"
)

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// main function is the entry point for the application.
func main() {
	// initializes the database (and the user if it does not exist)
	db.Init()

	// execute commands
	cmd.Execute()

	// push the db file to the git repository if it has not been pushed today yet
	db.AutoSave()
}
