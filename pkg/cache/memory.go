package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrItemNotFound = errors.New("cache: item not found")

type item struct {
	value     any
	createdAt int64
	ttl       int64
}

type MemoryCache struct {
	mutex sync.RWMutex
	cache map[any]*item
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{cache: make(map[any]*item)}
	c.setTTLTimer()

	return c
}

func (c *MemoryCache) setTTLTimer() {
	for {
		c.mutex.Lock()
		for k, v := range c.cache {
			if time.Now().Unix()-v.createdAt > v.ttl {
				delete(c.cache, k)
			}
		}
		c.mutex.Unlock()

		<-time.After(time.Second)
	}
}

func (c *MemoryCache) Set(key, value any, ttl int64) error {
	c.mutex.Lock()
	c.cache[key] = &item{
		value:     value,
		createdAt: time.Now().Unix(),
		ttl:       ttl,
	}
	c.mutex.Unlock()

	return nil
}

func (c *MemoryCache) Get(key any) (any, error) {
	c.mutex.RLock()
	item, ok := c.cache[key]
	c.mutex.RUnlock()

	if !ok {
		return nil, ErrItemNotFound
	}

	return item.value, nil
}
