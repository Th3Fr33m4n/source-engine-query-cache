package challenge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandom(t *testing.T) {
	ch := GenerateRandom()
	ch2 := GenerateRandom()

	assert.NotNil(t, ch)
	assert.NotNil(t, ch2)
	assert.Equal(t, 4, len(ch))
	assert.Equal(t, 4, len(ch2))
	assert.NotEqual(t, ch, ch2)
}

func TestGetForClient(t *testing.T) {
	clientId := "192.168.2.63"
	ch, _ := challengeMap.Get(clientId)
	ch2 := GetForClient(clientId)
	ch3 := GetForClient(clientId)

	assert.Nil(t, ch)
	assert.NotNil(t, ch2)
	assert.NotNil(t, ch3)
	assert.Equal(t, 4, len(ch2))
	assert.Equal(t, 4, len(ch3))
	assert.Equal(t, ch2, ch3)
}
