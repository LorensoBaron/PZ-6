package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	fmt.Printf("Кэш обновлен: %s=%s\n", key, value)
}

func main() {
	cache := NewCache()
	var wg sync.WaitGroup

	cache.Set("name", "Alice")

	// 5 читателей одновременно читают
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if val, ok := cache.Get("name"); ok {
				fmt.Printf("Читатель %d: %s\n", id, val)
			}
		}(i)
	}

	wg.Wait()
}
