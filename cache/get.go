package cache

import (
	"fmt"
	"github.com/xujiajun/nutsdb"
)

func (c *Client) Get(bucket string, key string) ([]byte, error) {
	var value []byte
	if err := c.db.View(func(tx *nutsdb.Tx) error {
		e, err := tx.Get(bucket, []byte(key))
		if err != nil {
			return fmt.Errorf("%v: %v; get: %w", bucket, key, err)
		}
		value = e.Value
		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}
