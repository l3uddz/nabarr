package cache

import (
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/rs/zerolog"
	"github.com/xujiajun/nutsdb"
)

type Client struct {
	log zerolog.Logger

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

	return &Client{
		log: nabarr.GetLogger("trace").With().Logger(),
		db:  db,
	}, nil
}
func (c *Client) Close() error {
	return c.db.Close()
}
