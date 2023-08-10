package a2s

import (
	"bytes"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
)

type Query struct {
	Type           domain.QueryType
	Body           []byte
	Header         []byte
	ResponseType   domain.ResponseType
	ResponseHeader []byte
}

func (q *Query) MatchResponse(res []byte) domain.ResponseType {
	if bytes.HasPrefix(res, q.ResponseHeader) {
		return q.ResponseType
	}

	rt := domain.InvalidResponse
	if bytes.HasPrefix(res, a2sChallengePrefix) {
		rt = domain.A2sChallengeResponse
	} else if q.Type == domain.A2sRules && bytes.HasPrefix(res, a2sMultipacketResponsePrefix) {
		rt = domain.A2sRulesSplitResponse
	}
	return rt
}

func (q *Query) IsQueryWChallenge(res []byte) bool {
	return bytes.HasPrefix(res, q.Header) && len(res) == len(q.Header)+challengeBytes
}

func (q *Query) Build(challenge []byte) []byte {
	return append(q.Header, challenge...)
}

func (q *Query) GetChallengeFromRequest(req []byte) []byte {
	hLen := len(q.Header)
	return req[hLen : hLen+challengeBytes]
}

func (q *Query) GetChallengeFromResponse(res []byte) []byte {
	hLen := len(a2sChallengePrefix)
	return res[hLen : hLen+challengeBytes]
}

var (
	InfoQuery = Query{
		Type:           domain.A2sInfo,
		Body:           a2sInfoRequest,
		Header:         a2sInfoRequest,
		ResponseType:   domain.A2sInfoResponse,
		ResponseHeader: a2sInfoResponsePrefix,
	}

	PlayersQuery = Query{
		Type:           domain.A2sPlayers,
		Body:           a2sPlayersRequest,
		Header:         a2sPlayersPrefix,
		ResponseType:   domain.A2sPlayersResponse,
		ResponseHeader: a2sPlayersResponsePrefix,
	}

	RulesQuery = Query{
		Type:           domain.A2sRules,
		Body:           a2sRulesRequest,
		Header:         a2sRulesPrefix,
		ResponseType:   domain.A2sRulesResponse,
		ResponseHeader: a2sRulesResponsePrefix,
	}
)

func BuildChallengeResponse(ch []byte) []byte {
	return append(a2sChallengePrefix, ch...)
}
