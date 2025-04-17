// Copyright 2025 The QuellDB Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the LICENSE file.

// Memory storage implementation for key-value pairs
package base

import (
	"sync"
	"time"
)

type MemStorage struct {
	data map[string]string
	mu   sync.RWMutex
	ttl  map[string]time.Time
}

// NewMemStorage creates a new instance of MemStorage.
// It initializes the internal map to store key-value pairs.
// The map is protected by a read-write mutex to allow concurrent access.
func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
	}
	go ms.ttlEvictionLoop()
	return ms
}

// Get retrieves the value associated with the given key.
// It returns the value and a boolean indicating whether the key exists in the storage.
// The method uses a read lock to allow concurrent reads.
func (m *MemStorage) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Put stores the key-value pair in the storage.
// If the key already exists, it updates the value.
// The method uses a write lock to ensure exclusive access during the operation.
func (m *MemStorage) Put(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Delete removes the key-value pair associated with the given key.
// The method uses a write lock to ensure exclusive access during the operation.
// If the key does not exist, it does nothing.
func (m *MemStorage) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// All returns a copy of all key-value pairs in the storage.
// It uses a read lock to allow concurrent reads.
// The returned map is a shallow copy, so modifications to it do not affect the original storage.
// This method is useful for iterating over all entries without locking the storage.
// It returns a new map containing all key-value pairs.
func (m *MemStorage) All() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cloned := make(map[string]string)
	for k, v := range m.data {
		cloned[k] = v
	}
	return cloned
}

// SetTTL sets a time-to-live (TTL) for the given key.
// The key will be automatically deleted after the specified duration.
// The method uses a write lock to ensure exclusive access during the operation.
func (m *MemStorage) ttlEvictionLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		m.mu.Lock()
		for k, exp := range m.ttl {
			if now.After(exp) {
				delete(m.data, k)
				delete(m.ttl, k)
			}
		}
		m.mu.Unlock()
	}
}
