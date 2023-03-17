package utils

import "sync"

type ConcurrentMap struct {
	l sync.Mutex
	m map[any]any
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{sync.Mutex{}, make(map[any]any)}
}

func (c *ConcurrentMap) Get(key any) any {
	c.l.Lock()
	defer c.l.Unlock()
	return c.m[key]
}

func (c *ConcurrentMap) Set(key any, value any) {
	c.l.Lock()
	c.m[key] = value
	c.l.Unlock()
}
