package pokecache

import (
	"sync"
	"time"
)

type (
	cacheEntry struct {
		createdAt time.Time
		val       []byte
	}
	Cache struct {
		cacheEntries map[string]cacheEntry
		mu           *sync.Mutex
	}
)

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntries: make(map[string]cacheEntry),
		mu:           &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cacheEntry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.cacheEntries[key] = cacheEntry
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.cacheEntries[key]
	return entry.val, ok
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		deadline := time.Now().Add(-interval)
		for key, entry := range cache.cacheEntries {
			if entry.createdAt.Before(deadline) {
				cache.mu.Lock()
				delete(cache.cacheEntries, key)
				cache.mu.Unlock()
			}
		}
	}
}
