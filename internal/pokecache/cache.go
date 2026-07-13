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
	cacheMap map[string]cacheEntry
	sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cacheMap := make(map[string]cacheEntry)
	cache := &Cache{
		cacheMap: cacheMap,
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.Lock()
	c.cacheMap[key] = entry
	c.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.cacheMap[key]
	if !ok {
		return []byte{}, false
	}
	return val.val, true

}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Lock()
		for key, entry := range c.cacheMap {
			if time.Since(entry.createdAt) > interval {
				delete(c.cacheMap, key)
			}
		}
		c.Unlock()
	}
}
