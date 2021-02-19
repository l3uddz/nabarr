package cache

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"time"
)

func (c *Client) Put(bucket string, key string, val []byte, ttl time.Duration) error {
	if err := c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(fmt.Sprintf("%s_%s", bucket, key)), val)
		if ttl > 0 {
			e = e.WithTTL(ttl)
		}
		return txn.SetEntry(e)
	}); err != nil {
		return fmt.Errorf("%v: %v; put: %w", bucket, key, err)
	}
	return nil
}
