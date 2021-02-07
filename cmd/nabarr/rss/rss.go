package rss

import (
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Client struct {
	cron *cron.Cron
	pvrs map[string]pvr.PVR

	log zerolog.Logger
}

func New(c Config, pvrs map[string]pvr.PVR) *Client {
	return &Client{
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
		pvrs: pvrs,

		log: nabarr.GetLogger(c.Verbosity).With().Logger(),
	}
}

func (c *Client) Start() {
	c.cron.Start()
}