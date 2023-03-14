package packets

import "bytes"

const challengeBytes = 4

type (
	QueryType    int
	ResponseType int
)

const (
	goldSrcMultipacketByte      = 8
	sourceMultipacketPacketByte = 8
	sourceMultipacketTotalByte  = 9
)

const (
	A2sInfo    QueryType = 0
	A2sPlayers QueryType = 1
	A2sRules   QueryType = 2
)

const (
	A2sInfoResponse       ResponseType = 0
	A2sPlayersResponse    ResponseType = 1
	A2sRulesResponse      ResponseType = 2
	A2sChallengeResponse  ResponseType = 3
	A2sRulesSplitResponse ResponseType = 4
	Invalid                            = -1
)

var (
	a2sInfoRequest = []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65,
		0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00,
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

func BuildChallengeResponse(challenge []byte) []byte {
	return append(a2sChallengePrefix, challenge...)
}

func IsA2sInfoRequest(p []byte) bool {
	return bytes.Equal(a2sInfoRequest, p)
}

func IsA2sInfoWChallenge(p []byte) bool {
	hasPrefix := bytes.HasPrefix(p, a2sInfoRequest)
	expectedLen := len(a2sInfoRequest) + challengeBytes
	return hasPrefix && len(p) == expectedLen
}

func GetChallengeFromA2sInfo(p []byte) []byte {
	return p[len(a2sInfoRequest) : len(a2sInfoRequest)+challengeBytes]
}

func IsA2sPlayersRequest(p []byte) bool {
	return bytes.Equal(a2sPlayersRequest, p)
}

func IsA2sPlayersWChallenge(p []byte) bool {
	hasPrefix := bytes.HasPrefix(p, a2sPlayersPrefix)
	expectedLen := len(a2sPlayersPrefix) + challengeBytes
	return hasPrefix && len(p) == expectedLen
}

func GetChallengeFromA2sPlayers(p []byte) []byte {
	return p[len(a2sPlayersPrefix) : len(a2sPlayersPrefix)+challengeBytes]
}

func GetChallengeFromA2sRules(p []byte) []byte {
	return p[len(a2sRulesPrefix) : len(a2sRulesPrefix)+challengeBytes]
}

func GetChallengeFromServerResponse(res []byte) []byte {
	l := len(a2sChallengePrefix)
	return res[l : l+challengeBytes]
}

func BuildQuery(qt QueryType, challenge []byte) []byte {
	var query []byte
	switch qt {
	case A2sInfo:
		query = append(a2sInfoRequest, challenge...)
	case A2sPlayers:
		query = append(a2sPlayersPrefix, challenge...)
	default:
		query = append(a2sRulesPrefix, challenge...)
	}
	return query
}

func CategorizeResponse(r []byte) ResponseType {
	var qt ResponseType
	if bytes.HasPrefix(r, a2sInfoResponsePrefix) {
		qt = A2sInfoResponse
	} else if bytes.HasPrefix(r, a2sPlayersResponsePrefix) {
		qt = A2sPlayersResponse
	} else if bytes.HasPrefix(r, a2sRulesPrefix) {
		qt = A2sRulesResponse
	} else if bytes.HasPrefix(r, a2sChallengePrefix) {
		qt = A2sChallengeResponse
	} else if bytes.HasPrefix(r, a2sMultipacketResponsePrefix) {
		qt = A2sRulesSplitResponse
	} else {
		qt = Invalid
	}
	return qt
}

func GetQuery(qt QueryType) []byte {
	var query []byte
	if qt == A2sInfo {
		query = a2sInfoRequest
	} else if qt == A2sPlayers {
		query = a2sPlayersRequest
	} else {
		query = a2sRulesRequest
	}
	return query
}

func GetQueryResponseType(qt QueryType) ResponseType {
	var rt ResponseType
	if qt == A2sInfo {
		rt = A2sInfoResponse
	} else if qt == A2sPlayers {
		rt = A2sPlayersResponse
	} else if qt == A2sRules {
		rt = A2sRulesResponse
	} else {
		rt = Invalid
	}
	return rt
}

func IsA2sRulesRequest(p []byte) bool {
	return bytes.Equal(p, a2sRulesRequest)
}

func IsA2sRulesWChallenge(p []byte) bool {
	hasPrefix := bytes.HasPrefix(p, a2sRulesPrefix)
	expectedLen := len(a2sRulesPrefix) + challengeBytes
	return hasPrefix && len(p) == expectedLen
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
