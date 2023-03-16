package a2s

import (
	"bytes"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
)

const challengeBytes = 4

const (
	goldSrcMultipacketByte      = 8
	sourceMultipacketPacketByte = 8
	sourceMultipacketTotalByte  = 9
)

var (
	a2sInfoRequest = []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E,
		0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00,
	}
	a2sChallengePrefix           = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x41}
	a2sPlayersRequest            = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF}
	a2sPlayersPrefix             = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55}
	a2sRulesRequest              = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
	a2sInfoResponsePrefix        = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x49}
	a2sRulesPrefix               = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56}
	a2sPlayersResponsePrefix     = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x44}
	a2sRulesResponsePrefix       = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x45}
	a2sMultipacketResponsePrefix = []byte{0xFE, 0xFF, 0xFF, 0xFF}
)

var TooManyRequests = []byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0x58, 0x54, 0x6f, 0x6f, 0x20, 0x6d, 0x61,
	0x6e, 0x79, 0x20, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x00,
}

func CategorizeRequest(req []byte) (domain.QueryType, bool) {
	if bytes.Equal(req, a2sInfoRequest) {
		return domain.A2sInfo, false
	} else if bytes.HasPrefix(req, a2sInfoRequest) && len(req) == len(a2sInfoRequest)+challengeBytes {
		return domain.A2sInfo, true
	} else if bytes.Equal(req, a2sPlayersRequest) {
		return domain.A2sPlayers, false
	} else if bytes.HasPrefix(req, a2sPlayersPrefix) && len(req) == len(a2sPlayersPrefix)+challengeBytes {
		return domain.A2sPlayers, true
	} else if bytes.Equal(req, a2sRulesRequest) {
		return domain.A2sRules, false
	} else if bytes.HasPrefix(req, a2sRulesPrefix) && len(req) == len(a2sRulesPrefix)+challengeBytes {
		return domain.A2sRules, true
	} else {
		return domain.InvalidQuery, false
	}
}

func ParseGoldsrcMultipacketResponse(r []byte) (byte, byte) {
	pNumber := r[goldSrcMultipacketByte]
	/* Upper 4 bits represent the number of the current packet (starting at 0)
	and bottom 4 bits represent the total number of packets (2 to 15).*/
	// get the upper 4 bits
	packetNumber := pNumber & 0xf0
	// get the lower 4 bits
	totalPackets := pNumber & 0x0f

	// shift bits to match the correct size
	packetNumber = packetNumber >> 4

	return packetNumber, totalPackets
}

func ParseSourceMultipacketResponse(r []byte) (byte, byte) {
	return r[sourceMultipacketPacketByte], r[sourceMultipacketTotalByte]
}
