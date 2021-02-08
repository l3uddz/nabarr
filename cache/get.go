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

func (c *Client) Get(bucket string, key string) ([]byte, error) {
	resp := make([]byte, 0)
	if err := c.db.View(func(tx *nutsdb.Tx) error {
		e, err := tx.Get(bucket, []byte(key))
		if err != nil {
			return fmt.Errorf("%v: %v; get: %w", bucket, key, err)
		}
		resp = e.Value
		return nil
	}); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Delete(bucket string, key string) error {
	if err := c.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(bucket, []byte(key))
	}); err != nil {
		return fmt.Errorf("%v: %v; delete: %w", bucket, key, err)
	}
	return nil
}
