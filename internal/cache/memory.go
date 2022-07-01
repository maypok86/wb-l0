package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/maypok86/wb-l0/internal/entity"
)

var ErrItemNotFound = errors.New("cache: item not found")

type item struct {
	value     *entity.Order
	createdAt int64
	ttl       int64
}

type MemoryCache struct {
	mutex sync.RWMutex
	cache map[int]*item
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{cache: make(map[int]*item)}
	go c.setTTLTimer()

	return c
}

func (c *MemoryCache) setTTLTimer() {
	for {
		now := time.Now().Unix()
		c.mutex.Lock()
		for k, v := range c.cache {
			if now-v.createdAt > v.ttl {
				delete(c.cache, k)
			}
		}
		c.mutex.Unlock()

		<-time.After(time.Second)
	}
}

func (c *MemoryCache) Set(key int, value *entity.Order, ttl int64) error {
	c.mutex.Lock()
	c.cache[key] = &item{
		value:     value,
		createdAt: time.Now().Unix(),
		ttl:       ttl,
	}
	c.mutex.Unlock()

	return nil
}

func (c *MemoryCache) Get(key int) (*entity.Order, error) {
	c.mutex.RLock()
	item, ok := c.cache[key]
	c.mutex.RUnlock()

	if !ok {
		return nil, ErrItemNotFound
	}

	return item.value, nil
}
