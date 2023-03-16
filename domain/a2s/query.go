package a2s

import (
	"bytes"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
)

type A2sQuery struct {
	QueryT         domain.QueryType
	Query          []byte
	Header         []byte
	ResponseT      domain.ResponseType
	ResponseHeader []byte
}

func (q *A2sQuery) MatchResponse(res []byte) domain.ResponseType {
	if bytes.HasPrefix(res, q.ResponseHeader) {
		return q.ResponseT
	}

	rt := domain.InvalidResponse
	if bytes.HasPrefix(res, a2sChallengePrefix) {
		rt = domain.A2sChallengeResponse
	} else if q.QueryT == domain.A2sRules && bytes.HasPrefix(res, a2sMultipacketResponsePrefix) {
		rt = domain.A2sRulesSplitResponse
	}
	return rt
}

func (q *A2sQuery) IsQueryWChallenge(res []byte) bool {
	return bytes.HasPrefix(res, q.Header) && len(res) == len(q.Header)+challengeBytes
}

func (q *A2sQuery) Build(challenge []byte) []byte {
	return append(q.Header, challenge...)
}

func (q *A2sQuery) GetChallengeFromRequest(req []byte) []byte {
	hLen := len(q.Header)
	return req[hLen : hLen+challengeBytes]
}

func (q *A2sQuery) GetChallengeFromResponse(res []byte) []byte {
	hLen := len(a2sChallengePrefix)
	return res[hLen : hLen+challengeBytes]
}

var (
	InfoQuery = A2sQuery{
		QueryT:         domain.A2sInfo,
		Query:          a2sInfoRequest,
		Header:         a2sInfoRequest,
		ResponseT:      domain.A2sInfoResponse,
		ResponseHeader: a2sInfoResponsePrefix,
	}

	PlayersQuery = A2sQuery{
		QueryT:         domain.A2sPlayers,
		Query:          a2sPlayersRequest,
		Header:         a2sPlayersPrefix,
		ResponseT:      domain.A2sPlayersResponse,
		ResponseHeader: a2sPlayersResponsePrefix,
	}

	RulesQuery = A2sQuery{
		QueryT:         domain.A2sRules,
		Query:          a2sRulesRequest,
		Header:         a2sRulesPrefix,
		ResponseT:      domain.A2sRulesResponse,
		ResponseHeader: a2sRulesResponsePrefix,
	}
)

func BuildChallengeResponse(ch []byte) []byte {
	return append(a2sChallengePrefix, ch...)
}
