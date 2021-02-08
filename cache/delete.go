package cache

import (
	"fmt"
	"github.com/xujiajun/nutsdb"
)

func (c *Client) Delete(bucket string, key string) error {
	if err := c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	}); err != nil {
		return fmt.Errorf("%v: %v; delete: %w", bucket, key, err)
	}
	return nil
}
