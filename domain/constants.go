package domain

type EngineType string

type (
	QueryType    int
	ResponseType int
)

const (
	GoldSrc EngineType = "goldsrc"
	Source  EngineType = "source"
)

const (
	A2sInfo      QueryType = 0
	A2sPlayers   QueryType = 1
	A2sRules     QueryType = 2
	InvalidQuery QueryType = 3
)

const (
	A2sInfoResponse       ResponseType = 0
	A2sPlayersResponse    ResponseType = 1
	A2sRulesResponse      ResponseType = 2
	A2sChallengeResponse  ResponseType = 3
	A2sRulesSplitResponse ResponseType = 4
	InvalidResponse       ResponseType = -1
)
