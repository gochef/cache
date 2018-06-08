package cache

import (
	"sync"
	"time"
)

type (
	// MemoryCacheItem represents an item to be put in cache
	MemoryCacheItem struct {
		value     interface{}
		expiresAt int64
	}

	// MemoryCache represents a memory cache driver instance
	MemoryCache struct {
		sync.RWMutex
		store map[string]*MemoryCacheItem
	}
)

// NewMemoryCache creates and returns a memory cache driver instance
func NewMemoryCache() Driver {
	return &MemoryCache{
		store: make(map[string]*MemoryCacheItem),
	}
}

// Get fetches an item from the cache
// returns the item and a boolean indicating whether the item was found
// false if not found, true if found
func (m *MemoryCache) Get(key string) (interface{}, bool) {
	if data := m.store[key]; data != nil {
		return data.value, true
	}

	return nil, false
}

// Put puts an item into the cache for the specified duration in seconds
// An expiration of less than 1 leaves the item in cache forever
func (m *MemoryCache) Put(key string, data interface{}, duration int64) {
	d := &MemoryCacheItem{
		value: data,
	}

	if duration < 1 {
		d.expiresAt = 0
	} else {
		d.expiresAt = time.Now().Unix() + duration
	}

	m.Lock()
	m.store[key] = d
	m.Unlock()
}

// Remove removes an item from the cache
func (m *MemoryCache) Remove(key string) {
	m.Lock()
	delete(m.store, key)
	m.Unlock()
}

// Clear empties the cache
func (m *MemoryCache) Clear() {
	m.Lock()
	m.store = make(map[string]*MemoryCacheItem)
	m.Unlock()
}
