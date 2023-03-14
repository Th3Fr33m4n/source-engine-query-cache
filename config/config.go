package config

import (
	"fmt"
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"

	"github.com/spf13/viper"
)

type GameServer struct {
	IP     string
	Port   string
	Engine domain.EngineType
}

func (g *GameServer) String() string {
	return fmt.Sprintf("%s:%s", g.IP, g.Port)
}

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
	Servers                  map[string]GameServer
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
