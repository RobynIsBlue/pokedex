package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.cacheMap[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if v, ok := c.cacheMap[key]; ok {
		return v.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	
}

func NewCache(interval int) {

}
