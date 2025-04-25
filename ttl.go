// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"time"

	"github.com/thirashapw/quelldb/constants"
)

// PutWithTTL stores a key-value pair in the database with a specified time-to-live (TTL).
// It first stores the pair in memory with the specified TTL and then writes it to the WAL.
// The function takes a key, a value, and a TTL duration as parameters.
// The TTL duration specifies how long the key-value pair should be valid.
// After the TTL expires, the key-value pair will be automatically removed from the in-memory storage.
// The function returns an error if any occurs during the write operation.
func (db *DB) PutTTL(key, value string, ttl time.Duration) error {
	db.memStorage.PutWithTTL(key, value, ttl)

	// publish to subscribers
	db.publish(ChangeEvent{
		Type:  constants.PUT,
		Key:   key,
		Value: value,
	})

	return db.wal.Write(constants.PUT, key, value)
}
