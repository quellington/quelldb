// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

package quelldb

import (
	"sort"
	"strings"
)

type Iterator struct {
	keys   []string
	values map[string]string
	index  int
}

// NewIterator creates a new iterator for the database.
// It retrieves all keys from the in-memory storage, sorts them,
// and initializes the iterator with the sorted keys and their corresponding values.
// The iterator starts at index -1, indicating that it is before the first element.
// The caller can use the Next() method to advance the iterator and access keys and values.
// The iterator is not thread-safe and should be used by a single goroutine at a time.
// The caller is responsible for closing the iterator when done.
// The iterator does not require any additional resources to be closed.
func (db *DB) Iterator() *Iterator {
	all := db.memStorage.All()
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &Iterator{
		keys:   keys,
		values: all,
		index:  -1,
	}
}


// NewPrefixIterator creates a new iterator for the database with a specific prefix.
// It retrieves all keys from the in-memory storage that start with the given prefix,
// sorts them, and initializes the iterator with the sorted keys and their corresponding values.
func (db *DB) PrefixIterator(prefix string) *Iterator {
	all := db.memStorage.All()
	filtered := make(map[string]string)
	for k, v := range all {
		if strings.HasPrefix(k, prefix) {
			filtered[k] = v
		}
	}
	keys := make([]string, 0, len(filtered))
	for k := range filtered {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &Iterator{
		keys:   keys,
		values: filtered,
		index:  -1,
	}
}


// Next advances the iterator to the next key-value pair.
func (it *Iterator) Next() bool {
	it.index++
	return it.index < len(it.keys)
}

func (it *Iterator) Key() string {
	return it.keys[it.index]
}

func (it *Iterator) Value() string {
	return it.values[it.keys[it.index]]
}
