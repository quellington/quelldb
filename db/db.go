// Database package for a key-value store
package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thirashapw/quelldb/constants"
)

type DB struct {
	memStorage *MemStorage
	wal        *WAL
	basePath   string
}

func Open(path string) (*DB, error) {
	os.MkdirAll(path, 0755)
	wal, err := NewWAL(filepath.Join(path, constants.LOG_FILE))
	if err != nil {
		return nil, err
	}
	db := &DB{
		memStorage: NewMemStorage(),
		wal:        wal,
		basePath:   path,
	}
	return db, nil
}

func (db *DB) Put(key, value string) error {
	db.memStorage.Put(key, value)
	return db.wal.Write("PUT", key, value)
}

func (db *DB) Get(key string) (string, bool) {
	// check MemTable first
	if val, ok := db.memStorage.Get(key); ok {
		return val, true
	}

	// check from newest SSTable to oldest
	files, _ := os.ReadDir(db.basePath)
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) {
			path := filepath.Join(db.basePath, f.Name())
			data, _ := ReadSSStorage(path)
			if val, ok := data[key]; ok {
				return val, true
			}
		}
	}
	return "", false
}

func (db *DB) Delete(key string) error {
	db.memStorage.Delete(key)
	return db.wal.Write("DEL", key, "")
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
	return WriteSSStorage(path, db.memStorage.All())
}

func (db *DB) Close() error {
	return db.wal.Close()
}
