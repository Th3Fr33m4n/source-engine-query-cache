package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMap_Get(t *testing.T) {
	m := NewConcurrentMap()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(mm *ConcurrentMap, wg *sync.WaitGroup) {
		m.Set("k1", "v1")
		wg.Done()
	}(m, wg)
	go func(mm *ConcurrentMap, wg *sync.WaitGroup) {
		m.Set("k2", "v2")
		wg.Done()
	}(m, wg)
	wg.Wait()
	assert.Equal(t, "v1", m.Get("k1"))
	assert.Equal(t, "v2", m.Get("k2"))
}
