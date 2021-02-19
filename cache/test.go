package cache

import (
	"github.com/xujiajun/nutsdb"
	"os"
	"path/filepath"
	"testing"
)

func newDb(t *testing.T, dir string) *nutsdb.DB {
	db, err := nutsdb.Open(nutsdb.Options{
		Dir:                  filepath.Join(os.TempDir(), dir),
		EntryIdxMode:         nutsdb.HintKeyValAndRAMIdxMode,
		SegmentSize:          8 * 1024 * 1024,
		NodeNum:              1,
		RWMode:               nutsdb.FileIO,
		SyncEnable:           true,
		StartFileLoadingMode: nutsdb.MMap,
	})
	if err != nil {
		t.Fatalf("newDb(dir: %v) open error: %v", dir, err)
	}
	return db
}
