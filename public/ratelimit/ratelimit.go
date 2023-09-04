package ratelimit

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	ttlMap "github.com/leprosus/golang-ttl-map"
	"golang.org/x/time/rate"
)

const (
	randomIntervalThreshold = 20
	baseExpirationTime      = 60
)

var (
	limiterMap    = ttlMap.New()
	globalLimiter = rate.NewLimiter(
		rate.Every(time.Millisecond*time.Duration(config.Get().RateLimitGlobal)),
		int(config.Get().RateLimitGlobalBurst))
)

func GetGlobalLimiter() *rate.Limiter {
	return globalLimiter
}

func GetLimiterForAddress(addr string) *rate.Limiter {
	savedLimiter, ok := limiterMap.Get(addr)

	if !ok {
		log.Debug("missing rate limiter for client " + addr)
	}

	if savedLimiter != nil {
		return savedLimiter.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(
		rate.Every(time.Millisecond*time.Duration(config.Get().RateLimitClient)),
		int(config.Get().RateLimitClientBurst))

	ttl := baseExpirationTime + int64(rand.Intn(randomIntervalThreshold))
	limiterMap.Set(addr, limiter, ttl)

	return limiter
}
