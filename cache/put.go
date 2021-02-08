package cache

import (
	"fmt"
	"github.com/xujiajun/nutsdb"
	"time"
)

func (c *Client) Put(bucket string, key string, val []byte, ttl time.Duration) error {
	if err := c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, []byte(key), val, uint32(ttl.Seconds()))
	}); err != nil {
		return fmt.Errorf("%v: %v; put: %w", bucket, key, err)
	}
	return nil
}
