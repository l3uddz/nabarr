package cache

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func (c *Client) Delete(bucket string, key string) error {
	if err := c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(fmt.Sprintf("%s_%v", bucket, key)))
	}); err != nil {
		return fmt.Errorf("%v: %v; delete: %w", bucket, key, err)
	}
	return nil
}
