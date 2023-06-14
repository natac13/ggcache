package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Cache is a simple in-memory key-value store.
type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

// New returns a new Cache.
func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	log.Printf("GET %s\n", string(key))

	v, ok := c.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	log.Printf("GOT %s from key: %s\n", v, string(key))
	return v, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value
	log.Printf("SET %s to %s\n", key, value)
	go func() {
		<-time.After(ttl)
		delete(c.data, string(key))
	}()

	return nil
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.data[string(key)]
	return ok
}
