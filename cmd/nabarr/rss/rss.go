package rss

import (
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"time"
)

type Client struct {
	cron *cron.Cron
	pvrs map[string]pvr.PVR

	feeds map[string]string
	log   zerolog.Logger
}

func New(c Config, pvrs map[string]pvr.PVR) *Client {
	return &Client{
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
		pvrs: pvrs,

		feeds: make(map[string]string, 0),
		log:   nabarr.GetLogger(c.Verbosity).With().Logger(),
	}
}

func (c *Client) Start() {
	c.cron.Start()
}

func (c *Client) Stop() {
	ctx := c.cron.Stop()
	select {
	case <-ctx.Done():
		c.log.Info().
			Strs("feed_names", c.feedNames()).
			Msg("Gracefully stopped")
	case <-time.After(5 * time.Second):
		c.log.Warn().
			Strs("feed_names", c.feedNames()).
			Msg("Forcefully stopped")
	}
}

func (c *Client) feedNames() []string {
	feeds := make([]string, 0)
	for f, _ := range c.feeds {
		feeds = append(feeds, f)
	}
	return feeds
}
