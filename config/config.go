package config

import (
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
	"github.com/spf13/viper"
)

type configWrapper struct {
	LogLevel                 string
	ReadBufferSize           uint
	Port                     string
	BufferSize               string
	ChallengeTTL             uint
	RateLimitGlobal          uint
	RateLimitGlobalBurst     uint
	RateLimitClient          uint
	RateLimitClientBurst     uint
	ServerInfoUpdateTimeout  time.Duration
	ServerInfoUpdateInterval int
	Servers                  map[string]domain.GameServer
}

var cfg configWrapper

func Init() {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.AddConfigPath("./")
	v.SetConfigName("config")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&cfg)

	if err != nil {
		panic(err)
	}
}

func Get() configWrapper {
	return cfg
}
