// main package, logs clean cron job
package main

import (
	"aio/pkg/log"
	"aio/pkg/utils/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

func cleanLogs(c *cron.Cron, bin string) {
	_, err := c.AddFunc("@every 24h", func() {
		checkForMainBinary(c, bin)                     // stop the cron job if the main binary is deleted
		log.Deb("--- clean logs cron job started ---") // print a separator

		today := time.Now().Format("2006-01-02") // Today's date
		logDir, err := fs.Path("logs")           // Path to the logs directory
		if err != nil {
			log.Err("failed to get the logs directory", "err", err)
			return
		}

		files, err := os.ReadDir(logDir)
		if err != nil {
			log.Err("failed to read log directory", "err", err)
			return
		}

		log.Deb("cleaning logs directory...")

		for _, file := range files {
			// Ignore directories and non-log files
			if file.IsDir() || filepath.Ext(file.Name()) != ".log" {
				continue
			}

			// If the file is not today's log file, delete it
			if file.Name() != today+".log" {
				filePath := filepath.Join(logDir, file.Name())
				if err := os.Remove(filePath); err != nil {
					log.Err("failed to remove old log file: "+filePath, "err", err)
				} else {
					log.Info("deleted old log file: " + filePath)
				}
			}
		}

		log.Deb("--- clean logs cron job ended ---") // print a separator
	})

	if err != nil {
		log.Err("failed to add cleanLogs cron job")
		log.Fat(err)
	}
}
