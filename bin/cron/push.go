// main package, push schedule cron job
package main

import (
	"aio/pkg/git"
	"aio/pkg/log"

	"github.com/robfig/cron/v3"
)

// pushSchedule function adds a cron job to the cron service.
// the cron job coomit the changes to the database and push them to the remote repository every 5 minutes.
func pushSchedule(c *cron.Cron, bin string) {
	_, err := c.AddFunc("@every 5m", func() {
		checkForMainBinary(c, bin) // check if the main binary exists
		wg.Wait()                  // wait for the previous job to finish
		log.Deb("--- save cron job started ---")

		err := git.Main() // check out the main branch
		if err != nil {
			log.Err("failed to check out the main branch", "err", err)
			return
		}

		err = git.Commit() // commit the changes to the database
		if err != nil {
			log.Err("failed to commit the changes", "err", err)
			return
		}

		err = git.Push() // push the database to the remote repository
		if err != nil {
			log.Err("failed to push the database to the remote repository", "err", err)
			return
		}

		log.Deb("--- save cron job ended ---")
	})

	if err != nil {
		log.Err("failed to add pushSchedule cron job")
		log.Fat(err)
	}
}
