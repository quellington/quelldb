package db

import (
	"os"
	"path/filepath"
)

type DB struct {
	memStorage *MemStorage
	wal        *WAL
	basePath   string
}

func Open(path string) (*DB, error) {
	os.MkdirAll(path, 0755)
	wal, err := NewWAL(filepath.Join(path, "wal.log"))
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
	return db.memStorage.Get(key)
}

func (db *DB) Delete(key string) error {
	db.memStorage.Delete(key)
	return db.wal.Write("DEL", key, "")
}

func (db *DB) Flush() error {
	// write SSTable
	data := db.memStorage.All()
	path := filepath.Join(db.basePath, "sstable.data")
	return WriteSSStorage(path, data)
}

func (db *DB) Close() error {
	return db.wal.Close()
}
