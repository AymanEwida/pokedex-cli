package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()

	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(internal time.Duration) {
	ch := time.Tick(internal)

	for {
		<-ch

		now := time.Now()

		c.mu.Lock()

		shouldDeleteKeys := []string{}

		for key, entry := range c.cache {
			if now.Sub(entry.createdAt) > internal {
				shouldDeleteKeys = append(shouldDeleteKeys, key)
			}
		}

		for _, key := range shouldDeleteKeys {
			delete(c.cache, key)
		}

		c.mu.Unlock()
	}
}
