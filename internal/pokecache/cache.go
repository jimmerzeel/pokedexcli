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
	items map[string]cacheEntry
	mtx   sync.Mutex
	done  chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	// create a new Cache instance with an initialized map and mutex
	cache := &Cache{
		items: make(map[string]cacheEntry),
		done:  make(chan struct{}),
	}

	// start the reap loop as a background Goroutine
	go cache.reapLoop(interval)

	// return the fully prepared cache instance
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.items[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if item, ok := c.items[key]; ok {
		return item.val, true
	}
	return nil, false
}

func (c *Cache) Stop() {
	// Safely close the `done` channel once
	select {
	case <-c.done: // If already closed, skip
		return
	default:
		close(c.done)
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.mtx.Lock()
			for key, value := range c.items {
				if time.Since(value.createdAt) > interval {
					delete(c.items, key)
				}
			}
			c.mtx.Unlock()
		}
	}
}
