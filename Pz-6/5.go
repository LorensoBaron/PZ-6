package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	value     string
	expiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
	go c.cleanup()
	return c
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	fmt.Printf(" Кэш: %s = %s (живет %v)\n", key, value, c.ttl)
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return "", false
	}

	if time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return "", false
	}

	return item.value, true
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.expiresAt) {
				delete(c.items, k)
				fmt.Printf(" Очищен устаревший ключ: %s\n", k)
			}
		}
		c.mu.Unlock()
	}
}

func main() {
	cache := NewCache(3 * time.Second)

	cache.Set("user1", "Анна")
	cache.Set("user2", "Борис")

	// Проверяем сразу
	if val, ok := cache.Get("user1"); ok {
		fmt.Printf(" Чтение: user1 = %s\n", val)
	}

	// Ждем 4 секунды
	fmt.Println(" Ждем 4 секунды...")
	time.Sleep(4 * time.Second)

	// Проверяем снова
	if _, ok := cache.Get("user1"); !ok {
		fmt.Println(" user1: истек срок жизни")
	}
}
