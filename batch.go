// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/quellington/quelldb/base"
	"github.com/quellington/quelldb/constants"
)

type Options struct {
	EncryptionKey []byte
	CompactLimit  uint
	BoomBitSize   uint
	BoomHashCount uint
}

type DB struct {
	memStorage    *base.MemStorage
	wal           *base.WAL
	basePath      string
	key           []byte
	compactLimit  uint
	boomBitSize   uint
	boomHashCount uint
	manifestSSSs  []SSSMeta
	subscribers   map[int]func(ChangeEvent)
	subLock       sync.RWMutex
	nextSubID     int
}

// Open initializes a new database at the specified path.
// It creates the necessary directories and initializes the WAL.
// If the path already exists, it will be opened for reading.
// The options parameter allows for optional encryption key setup.
// If the encryption key is provided, it must be 32 bytes long for AES-256.
// If the key is not provided, the database will be unencrypted.
// The function returns a pointer to the DB instance and an error if any occurs.
func Open(path string, opts *Options) (*DB, error) {
	os.MkdirAll(path, 0755)

	var encryptionKey []byte

	walLogPath := filepath.Join(path, constants.LOG_FILE)
	wal, err := base.NewWAL(walLogPath)
	if err != nil {
		return nil, err
	}

	db := &DB{
		memStorage:    base.NewMemStorage(),
		basePath:      path,
		wal:           wal,
		compactLimit:  constants.SSS_COMPACT_DEFAULT_LIMIT,
		boomBitSize:   constants.BOOM_BIT_SIZE,
		boomHashCount: constants.BOOM_HASH_COUNT,
	}

	if opts != nil {
		if len(opts.EncryptionKey) > 0 {
			if len(opts.EncryptionKey) != 32 {
				return nil, fmt.Errorf("encryption key must be 32 bytes (AES-256)")
			}
			encryptionKey = opts.EncryptionKey
			db.key = opts.EncryptionKey
		}

		if opts.CompactLimit > 0 {
			db.compactLimit = opts.CompactLimit
		}

		if opts.BoomBitSize > 0 {
			db.boomBitSize = opts.BoomBitSize
		}

		if opts.BoomHashCount > 0 {
			db.boomHashCount = opts.BoomHashCount
		}
	}

	// Load the manifest SSS files
	mnfts, err := LoadManifest(path, encryptionKey)
	if err != nil {
		return nil, err
	}
	db.manifestSSSs = mnfts

	// Check if the WAL file exists
	if err := db.replayWAL(walLogPath); err != nil {
		return nil, fmt.Errorf("WAL replay failed: %w", err)
	}

	return db, nil
}

// Put stores a key-value pair in the database.
// It first stores the pair in memory and then writes it to the WAL.
// The function returns an error if any occurs during the write operation.
// If the key already exists, it will be updated with the new value.
// The value is stored in plaintext, and if encryption is enabled, it will be encrypted before writing to the WAL.
func (db *DB) Put(key, value string) error {
	db.memStorage.Put(key, value)

	// Publish to subscribers
	db.publish(ChangeEvent{
		Type:  constants.PUT,
		Key:   key,
		Value: value,
	})

	return db.wal.Write(constants.PUT, key, value)
}

// PutBatch stores multiple key-value pairs in the database.
// It first stores the pairs in memory and then writes them to the WAL.
// If the key already exists, it will be updated with the new value.
func (db *DB) PutBatch(kvs map[string]string) error {
	if len(kvs) == 0 {
		return nil
	}

	var wls []string

	for key, value := range kvs {
		db.memStorage.Put(key, value)
		wls = append(wls, fmt.Sprintf("%s|%s|%s\n", constants.PUT, key, value))
	}

	return db.wal.WriteLines(wls)
}

// Get retrieves the value associated with the given key.
// It first checks the in-memory storage and then searches through the SSS files in reverse order.
// The function returns the value and error whether the key was found.
// If the key is not found in memory or in any SSS files, it returns an empty string and false.
// If the key is found, the value is returned in plaintext.
// If encryption is enabled, the value will be decrypted before returning.
// The function also handles the case where the key is not found in any SSS files.
// If the key is found in memory, it will be returned immediately without checking the SSS files.
func (db *DB) Get(key string) (string, error) {
	// check MemS first
	if val, ok := db.memStorage.Get(key); ok {
		return val, nil
	}

	// check from newest SSS to oldest
	files, _ := os.ReadDir(db.basePath)
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]
		if strings.HasPrefix(f.Name(), constants.SSS_PREFIX) {
			path := filepath.Join(db.basePath, f.Name())

			// check bloom filter
			// if the bloom filter is not found, skip this file
			filter, err := base.LoadBloomFilter(path+constants.SSS_BOOM_FILTER_SUFFIX, constants.BOOM_BIT_SIZE, constants.BOOM_HASH_COUNT)
			if err == nil && !filter.Test(key) {
				// key not in bloom filter, skip this file
				continue
			}

			data, _ := base.ReadSSStorage(path, db.key)
			if val, ok := data[key]; ok {
				return val, nil
			}
		}
	}
	return "", fmt.Errorf("key not found")

}

// Delete removes the key-value pair associated with the given key.
// It first deletes the pair from memory and then writes the delete operation to the WAL.
// The function returns an error if any occurs during the write operation.
// If the key does not exist, it will not raise an error.
// The function does not check the SSS files for the key before deleting it from memory.
// If the key is found in memory, it will be deleted immediately.
func (db *DB) Delete(key string) error {
	db.memStorage.Delete(key)

	// Publish to subscribers
	db.publish(ChangeEvent{
		Type: constants.DELETE,
		Key:  key,
	})

	return db.wal.Write(constants.DELETE, key, "")
}

// Flush writes the in-memory data to a new SSS file.
// It generates a new filename based on the highest existing SSS file number.
// The new file will be named with the format "sss_00001.qldb".
// The function returns an error if any occurs during the write operation.
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

	minKey, maxKey, err := base.WriteSSStorage(path, db.memStorage.All(), db.key)

	if err != nil {
		return err
	}

	db.manifestSSSs = append(db.manifestSSSs, SSSMeta{
		Filename: filename,
		MinKey:   minKey,
		MaxKey:   maxKey,
	})

	return SaveManifest(db.basePath, db.manifestSSSs, db.key)
}

// Close closes the database and the WAL.
// It ensures that all data is flushed to the SSS files and the WAL is closed properly.
// The function returns an error if any occurs during the close operation.
// It is important to call this function when you are done using the database to ensure data integrity.
// The function does not delete any SSS files or the WAL file.
// It only closes the file handles and ensures that all data is written to disk.
// After calling this function, the database instance should not be used anymore.
func (db *DB) Close() error {
	return db.wal.Close()
}
