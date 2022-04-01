package cache

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func (c *Client) Get(bucket string, key string) ([]byte, error) {
	var value []byte

	if err := c.db.View(func(txn *badger.Txn) error {
		// get key
		v, err := txn.Get([]byte(fmt.Sprintf("%s_%s", bucket, key)))
		if err != nil {
			return fmt.Errorf("%v: %v: get error: %w", bucket, key, err)

		}

		// validate value
		if v.IsDeletedOrExpired() {
			return fmt.Errorf("%v: %v: get: key does not exist", bucket, key)
		}

		// read value
		err = v.Value(func(val []byte) error {
			value = val
			return nil
		})
		if err != nil {
			return fmt.Errorf("%v: %v: get: value error: %w", bucket, key, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}
