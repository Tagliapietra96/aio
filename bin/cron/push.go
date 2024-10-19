// main package, push schedule cron job
package main

import (
	"aio/pkg/git"
	"aio/pkg/log"

	"github.com/robfig/cron/v3"
)

// pushSchedule function adds a cron job to the cron service.
// the cron job coomit the changes to the database and push them to the remote repository every 5 minutes.
func pushSchedule(c *cron.Cron, check func()) {
	c.AddFunc("@every 5m", func() {
		check()   // check if the main binary exists
		wg.Wait() // wait for the previous job to finish
		log.Deb("--- save cron job started ---")

		err := git.Main() // check out the main branch
		if err != nil {
			log.Err("failed to check out the main branch", "err", err)
			return
		}

		log.Deb("committing changes to the database...")
		err = git.Commit() // commit the changes to the database
		if err != nil {
			log.Err("failed to commit the changes", "err", err)
			return
		}
		log.Info("changes committed successfully")

		err = git.Push() // push the database to the remote repository
		if err != nil {
			log.Err("failed to push the database to the remote repository", "err", err)
			return
		}
		log.Info("database pushed to the remote repository successfully")

		log.Deb("--- save cron job ended ---")
	})
}
