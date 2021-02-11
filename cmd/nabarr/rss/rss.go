package rss

import (
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"time"
)

type Client struct {
	cron             *cron.Cron
	cache            *cache.Client
	cacheDuration    time.Duration
	cacheFiltersHash string
	pvrs             map[string]pvr.PVR

	log zerolog.Logger
}

func New(c Config, cc *cache.Client, cfh string, pvrs map[string]pvr.PVR) *Client {
	return &Client{
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
		cache:            cc,
		cacheDuration:    24 * time.Hour,
		cacheFiltersHash: cfh,
		pvrs:             pvrs,

		log: nabarr.GetLogger(c.Verbosity).With().Logger(),
	}
}

func (c *Client) Start() {
	c.cron.Start()
}

func (c *Client) Stop() {
	ctx := c.cron.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
	}
}
