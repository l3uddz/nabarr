package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/l3uddz/nabarr/logger"
	"github.com/lefelys/state"
	"github.com/rs/zerolog"
)

type Client struct {
	log zerolog.Logger
	st  state.State

	db *badger.DB
}

func New(path string) (*Client, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	log := logger.New("").With().Logger()

	// start cleaner
	st, tail := state.WithShutdown()
	ticker := time.NewTicker(6 * time.Hour)
	go func() {
		for {
			select {
			case <-tail.End():
				ticker.Stop()
				tail.Done()
				return
			case <-ticker.C:
				// clean cache
				for {
					if db.RunValueLogGC(0.5) == nil {
						continue
					}
					break
				}
				log.Debug().Msg("Cleaned cache")
			}
		}
	}()

	return &Client{
		log: log,
		st:  st,

		db: db,
	}, nil
}

func (c *Client) Close() error {
	// shutdown cleaner
	if c.st != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := c.st.Shutdown(ctx); err != nil {
			c.log.Error().
				Err(err).
				Msg("Failed shutting down cache cleaner gracefully")
		}
	}

	// close cache
	return c.db.Close()
}
