package rss

import (
	"fmt"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/robfig/cron/v3"
)

func (c *Client) AddJob(feed feedItem) error {
	// prepare job
	if feed.Cron == "" {
		feed.Cron = "*/15 * * * *"
	}

	// create job
	job := &rssJob{
		name: feed.Name,
		log:  c.log.With().Str("feed_name", feed.Name).Logger(),
		url:  feed.URL,
		pvrs: make(map[string]pvr.PVR, 0),

		attempts: 0,
		errors:   make([]error, 0),
	}

	// add pvrs
	for _, p := range feed.Pvrs {
		po, exists := c.pvrs[p]
		if !exists {
			return fmt.Errorf("pvr object does not exist: %v", p)
		}
		job.pvrs[p] = po
	}

	// schedule job
	if id, err := c.cron.AddJob(feed.Cron, cron.NewChain(
		cron.SkipIfStillRunning(cron.DiscardLogger)).Then(job),
	); err != nil {
		return fmt.Errorf("add job: %w", err)
	} else {
		job.jobID = id
	}

	job.log.Info().Msg("Initialised")
	return nil
}

func (j *rssJob) Run() {
	// increase attempt counter
	j.attempts++

	// run job
	err := j.process()

	// handle job response
	switch {
	case err == nil:
		// job completed successfully
		j.attempts = 0
		j.errors = j.errors[:0]
		return

	default:
		j.log.Warn().
			Err(err).
			Int("attempts", j.attempts).
			Msg("Unexpected error occurred")
		j.errors = append(j.errors, err)
	}

	if j.attempts > 5 {
		j.log.Error().
			Errs("error", j.errors).
			Int("attempts", j.attempts).
			Msg("Consecutive errors occurred while refreshing rss, job has been stopped...")
		j.cron.Remove(j.jobID)
	}
}
