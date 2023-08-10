package a2s

import (
	"testing"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"

	"github.com/stretchr/testify/assert"
)

func TestParseGoldsrcMultipacketResponse(t *testing.T) {
	p1 := []byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xB4, 0x27, 0x00, 0x00, 0x02, 0xFF,
		0xFF, 0xFF, 0xFF, 0x45, 0x6D, 0x00, 0x61, 0x64, 0x6D,
	}

	p2 := []byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xB4, 0x27, 0x00, 0x00, 0x12, 0x61,
		0x72, 0x74, 0x6D, 0x6F, 0x6E, 0x65, 0x79, 0x00, 0x38,
	}

	pNumber, totalPackets := ParseGoldsrcMultipacketResponse(p1)

	assert.Equal(t, byte(0), pNumber)
	assert.Equal(t, byte(2), totalPackets)

	pNumber, totalPackets = ParseGoldsrcMultipacketResponse(p2)

	assert.Equal(t, byte(1), pNumber)
	assert.Equal(t, byte(2), totalPackets)
}

func TestParseSourceMultipacketResponse(t *testing.T) {
	p1 := []byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xB4, 0x27, 0x00, 0x00, 0x00, 0x02,
		0xFF, 0xFF, 0xFF, 0x45, 0x6D, 0x00, 0x61, 0x64, 0x6D,
	}
	p2 := []byte{
		0xFE, 0xFF, 0xFF, 0xFF, 0xB4, 0x27, 0x00, 0x00, 0x01, 0x02,
		0x72, 0x74, 0x6D, 0x6F, 0x6E, 0x65, 0x79, 0x00, 0x38,
	}

	pNumber, totalPackets := ParseSourceMultipacketResponse(p1)

	assert.Equal(t, byte(0), pNumber)
	assert.Equal(t, byte(2), totalPackets)

	pNumber, totalPackets = ParseSourceMultipacketResponse(p2)

	assert.Equal(t, byte(1), pNumber)
	assert.Equal(t, byte(2), totalPackets)
}

func TestCategorizeRequest(t *testing.T) {
	cat, hasChallenge := CategorizeRequest(a2sInfoRequest)
	assert.Equal(t, domain.A2sInfo, cat)
	assert.False(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest(a2sPlayersRequest)
	assert.Equal(t, domain.A2sPlayers, cat)
	assert.False(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest(a2sRulesRequest)
	assert.Equal(t, domain.A2sRules, cat)
	assert.False(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest(append(a2sInfoRequest, 0xff, 0x45, 0x23, 0x89))
	assert.Equal(t, domain.A2sInfo, cat)
	assert.True(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest(append(a2sPlayersPrefix, 0xff, 0x45, 0x23, 0x89))
	assert.Equal(t, domain.A2sPlayers, cat)
	assert.True(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest(append(a2sRulesPrefix, 0xff, 0x45, 0x23, 0x89))
	assert.Equal(t, domain.A2sRules, cat)
	assert.True(t, hasChallenge)
	cat, hasChallenge = CategorizeRequest([]byte{0xff, 0x45, 0x23, 0x89})
	assert.Equal(t, domain.InvalidQuery, cat)
	assert.False(t, hasChallenge)
}
