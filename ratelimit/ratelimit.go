package ratelimit

import (
	"log"
	"math/rand"
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/leprosus/golang-ttl-map"
	"golang.org/x/time/rate"
)

const (
	randomIntervalThreshold = 20
	baseExpirationTime      = 60
)

var (
	limiterMap    = ttl_map.New()
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
		log.Println("error obtaining rate limiter for client " + addr)
	}
	if savedLimiter != nil {
		return savedLimiter.(*rate.Limiter)
	}
	limiter := rate.NewLimiter(
		rate.Every(time.Millisecond*time.Duration(config.Get().RateLimitClient)),
		int(config.Get().RateLimitClientBurst))
	limiterMap.Set(addr, limiter, baseExpirationTime+int64(rand.Intn(randomIntervalThreshold)))
	return limiter
}
