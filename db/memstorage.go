// Memory storage implementation for key-value pairs

package db

import "sync"

type MemStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func (m *MemStorage) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func (m *MemStorage) Put(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *MemStorage) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *MemStorage) All() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cloned := make(map[string]string)
	for k, v := range m.data {
		cloned[k] = v
	}
	return cloned
}
