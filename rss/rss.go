package rss

import (
	"time"

	"github.com/lefelys/state"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/l3uddz/nabarr/logger"
)

type Client struct {
	cron  *cron.Cron
	cache *cache.Client
	pvrs  map[string]pvr.PVR

	log zerolog.Logger
}

func New(c Config, cc *cache.Client, pvrs map[string]pvr.PVR) *Client {
	return &Client{
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
		cache: cc,
		pvrs:  pvrs,

		log: logger.Child(logger.WithLevel(c.Verbosity)),
	}
}

func (c *Client) Start() state.State {
	c.cron.Start()

	st, tail := state.WithShutdown()
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-tail.End():
				ticker.Stop()

				// shutdown cron
				ctx := c.cron.Stop()
				select {
				case <-ctx.Done():
				case <-time.After(5 * time.Second):
				}

				tail.Done()
				return
			case <-ticker.C:
			}
		}
	}()

	return st
}
