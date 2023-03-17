package ratelimit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGlobalLimiter(t *testing.T) {
	rl := GetGlobalLimiter()
	assert.NotNil(t, rl)
}

func TestGetLimiterForAddress(t *testing.T) {
	rl := GetLimiterForAddress("192.168.1.2")
	rl2 := GetLimiterForAddress("192.168.4.3")
	rl3 := GetLimiterForAddress("192.168.1.2")
	assert.NotNil(t, rl)
	assert.NotNil(t, rl2)
	assert.NotNil(t, rl3)
	assert.True(t, rl == rl3)
	assert.False(t, rl == rl2)
}
