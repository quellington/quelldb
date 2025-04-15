// Database package for a key-value store
// v1.0.1
package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thirashapw/quelldb/constants"
)

type Options struct {
	EncryptionKey []byte
}

type DB struct {
	memStorage *MemStorage
	wal        *WAL
	basePath   string
	key        []byte
}

func Open(path string, opts *Options) (*DB, error) {
	os.MkdirAll(path, 0755)

	wal, err := NewWAL(filepath.Join(path, constants.LOG_FILE))
	if err != nil {
		return nil, err
	}

	db := &DB{
		memStorage: NewMemStorage(),
		basePath:   path,
		wal:        wal,
	}

	if opts != nil && len(opts.EncryptionKey) > 0 {
		if len(opts.EncryptionKey) != 32 {
			return nil, fmt.Errorf("encryption key must be 32 bytes (AES-256)")
		}
		db.key = opts.EncryptionKey
	}

	return db, nil
}

func (db *DB) Put(key, value string) error {
	db.memStorage.Put(key, value)
	return db.wal.Write(constants.PUT, key, value)
}

func (db *DB) Get(key string) (string, bool) {
	// check MemS first
	if val, ok := db.memStorage.Get(key); ok {
		return val, true
	}

	// check from newest SSS to oldest
	files, _ := os.ReadDir(db.basePath)
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) {
			path := filepath.Join(db.basePath, f.Name())
			data, _ := ReadSSStorage(path, db.key)
			if val, ok := data[key]; ok {
				return val, true
			}
		}
	}
	return "", false

}

func (db *DB) Delete(key string) error {
	db.memStorage.Delete(key)
	return db.wal.Write(constants.DELETE, key, "")
}

func (db *DB) Flush() error {
	files, _ := os.ReadDir(db.basePath)
	id := 0
	for _, f := range files {
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) {
			numStr := strings.TrimSuffix(strings.TrimPrefix(f.Name(), constants.SSS_PREFIX), constants.SSS_SUFFIX)
			num, _ := strconv.Atoi(numStr)
			if num >= id {
				id = num + 1
			}
		}
	}

	filename := fmt.Sprintf("%s%05d%s", constants.SSS_PREFIX, id, constants.SSS_SUFFIX)
	path := filepath.Join(db.basePath, filename)
	return WriteSSStorage(path, db.memStorage.All(), db.key)
}

func (db *DB) Close() error {
	return db.wal.Close()
}
