package cache

import (
	"github.com/dgraph-io/badger/v3"
	"testing"
)

func newDb(t *testing.T) *badger.DB {
	opts := badger.DefaultOptions("").WithInMemory(true)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("newDb() open error: %v", err)
	}
	return db
}
