package main

import (
	"time"
	"sync"
)

type Record struct {
	deadline    time.Time
	value       string
}

type Cache struct {
	mu     sync.Mutex
	state  map[string]Record
	ttl    time.Duration
}

func newCache(ttl int) (c Cache) {
	if ttl == 0 {
		ttl = 5
	}
	c = Cache{
		state: make(map[string]Record),
		ttl:   time.Minute * time.Duration(ttl),
	}
	return
}

func (c *Cache) Set(key string, value string) {
	c.mu.Lock()
	c.state[key] = Record{
		time.Now().Add(c.ttl),
		value,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (ret string, ok bool) {
	var rec Record

	c.mu.Lock()
	rec, ok = c.state[key]
	if ok {
		if time.Now().After(rec.deadline) {
			delete(c.state, key)
			ok = false
		} else {
			ret = rec.value
		}
	}
	c.mu.Unlock()
	return
}