package challenge

import (
	"crypto/rand"
	mrand "math/rand"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	ttlMap "github.com/leprosus/golang-ttl-map"
)

const randomIntervalThreshold = 20

var challengeMap = ttlMap.New()

func GenerateRandom() []byte {
	ch := make([]byte, 4)
	_, err := rand.Read(ch)
	if err != nil {
		panic(err)
	}
	return ch
}

func GetForClient(addr string) []byte {
	val, ok := challengeMap.Get(addr)
	var ch []byte
	if !ok {
		ch = GenerateRandom()
	} else {
		ch = val.([]byte)
	}
	challengeMap.Set(addr, ch, getTTL())
	return ch
}

func getTTL() int64 {
	return int64(config.Get().ChallengeTTL) + int64(mrand.Intn(randomIntervalThreshold))
}
