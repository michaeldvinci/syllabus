package cache

import (
	"sync"
	"time"

	"github.com/michaeldvinci/syllabus/internal/models"
)

// Cache provides thread-safe caching with TTL
type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
	ttl   time.Duration
}

type cacheItem struct {
	val       models.SeriesInfo
	expiresAt time.Time
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (models.SeriesInfo, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(it.expiresAt) {
		return models.SeriesInfo{}, false
	}
	return it.val, true
}

// Set stores an item in the cache
func (c *Cache) Set(key string, v models.SeriesInfo) {
	c.mu.Lock()
	c.items[key] = cacheItem{
		val:       v,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}