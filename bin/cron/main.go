////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// File Name: cron_script.go
// Created by: Matteo Tagliapietra 2024-10-17

// This is the main entry point for the cron service.
// It initializes the cron service and adds cron jobs.
// This script is used to run cron jobs in the background.

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// main package is the entry point for the cron service.
package main

import (
	"aio/pkg/log"
	"aio/pkg/utils/fs"
	"sync"

	"os"

	"github.com/robfig/cron/v3"
)

var wg sync.WaitGroup

// checkForMainBinary function checks if the main binary exists.
// if the main binary does not exist, it stops the cron job and exits the program.
func checkForMainBinary(c *cron.Cron, bin string) {
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		c.Stop()
		log.Err("main cron binary not found, stopping cron job")
		log.Fat(err)
	}
}

// monitorMainBinary function monitors the main binary.
// if the main binary is deleted, it stops the cron job and exits the program.
func monitorMainBinary(c *cron.Cron, bin string) {
	_, err := c.AddFunc("@every 10s", func() {
		checkForMainBinary(c, bin)
	})
	if err != nil {
		log.Err("failed to add monitorMainBinary cron job")
		log.Fat(err)
	}
}

// main function is the entry point for the cron service.
// it initializes the cron service and adds cron jobs.
// it starts the cron service and keeps it running.
func main() {
	bin, err := fs.Path("cron") // Path to the main binary
	if err != nil {
		log.Err("failed to get the main binary path", "err", err)
	}

	c := cron.New()

	// jobs list
	pushSchedule(c, bin) // push db to the remote repository every 5 minutes
	cleanLogs(c, bin)    // clean the logs directory every 24 hours

	monitorMainBinary(c, bin) // monitor the main binary every 10 seconds
	c.Start()                 // start the cron service
	select {}                 // keep the cron service running
}
