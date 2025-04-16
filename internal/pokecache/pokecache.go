package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Interval time.Duration
	CacheMap map[string]cacheEntry
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CacheMap[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if v, ok := c.CacheMap[key]; ok {
		return v.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for v, t := range c.CacheMap {
		difference := time.Now().Sub(t.createdAt)
		if difference > c.Interval {
			delete(c.CacheMap, v)
		}
	}
}

func NewCache(i time.Duration) *Cache {
	ticker := time.NewTicker(i)
	c := &Cache{
		Interval: i,
		CacheMap: map[string]cacheEntry{},
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				c.reapLoop()
			}
		}
	}()
	return c
}
