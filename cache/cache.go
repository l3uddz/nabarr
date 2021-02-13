package cache

import (
	"context"
	"fmt"
	"github.com/l3uddz/nabarr/logger"
	"github.com/lefelys/state"
	"github.com/rs/zerolog"
	"github.com/xujiajun/nutsdb"
	"time"
)

type Client struct {
	log zerolog.Logger

	st state.State
	db *nutsdb.DB
}

func New(path string) (*Client, error) {
	db, err := nutsdb.Open(nutsdb.Options{
		Dir:                  path,
		EntryIdxMode:         nutsdb.HintKeyValAndRAMIdxMode,
		SegmentSize:          8 * 1024 * 1024,
		NodeNum:              1,
		RWMode:               nutsdb.FileIO,
		SyncEnable:           true,
		StartFileLoadingMode: nutsdb.MMap,
	})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	log := logger.New("trace").With().Logger()

	// start cleaner
	st, tail := state.WithShutdown()
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-tail.End():
				ticker.Stop()
				tail.Done()
				return
			case <-ticker.C:
				// clean cache
				err := db.Update(func(tx *nutsdb.Tx) error {
					return db.Merge()
				})

				switch {
				case err == nil:
					log.Info().Msg("Cleaned cache")
				case err.Error() == "the number of files waiting to be merged is at least 2":
				// there were no data files to be merged
				default:
					// unexpected error
					log.Error().
						Err(err).
						Msg("Failed cleaning cache")
				}
			}
		}
	}()

	return &Client{
		log: log,
		st:  st,
		db:  db,
	}, nil
}

func (c *Client) Close() error {
	// shutdown cleaner
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := c.st.Shutdown(ctx); err != nil {
		c.log.Error().
			Err(err).
			Msg("Failed shutting down cache cleaner gracefully")
	}

	// close cache
	return c.db.Close()
}
