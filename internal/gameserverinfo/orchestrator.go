package gameserverinfo

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/packets"
	"github.com/go-co-op/gocron"
	"github.com/leprosus/golang-ttl-map"
)

var (
	infoMap   = ttl_map.New()
	scheduler *gocron.Scheduler
)

type ServerInfo struct {
	Info    [][]byte
	Rules   [][]byte
	Players [][]byte
	Server  config.GameServer
}

func (s *ServerInfo) AddInfo(qt packets.QueryType, info [][]byte) {
	switch qt {
	case packets.A2sInfo:
		s.Info = info
	case packets.A2sRules:
		s.Rules = info
	case packets.A2sPlayers:
		s.Players = info
	}
}

func (s *ServerInfo) GetInfo(qt packets.QueryType) [][]byte {
	var res [][]byte
	switch qt {
	case packets.A2sInfo:
		res = s.Info
	case packets.A2sRules:
		res = s.Rules
	case packets.A2sPlayers:
		res = s.Players
	}
	return res
}

func getTTL() int64 {
	return int64(3 * config.Get().ServerInfoUpdateInterval)
}

func RegisterServer(g config.GameServer) {
	infoMap.Set(g.String(), &ServerInfo{Server: g}, getTTL())
}

func GetServerInfo(g config.GameServer) (*ServerInfo, error) {
	info, ok := infoMap.Get(g.String())
	if !ok {
		return nil, ErrMissingServerInfo
	}
	return info.(*ServerInfo), nil
}

func Init() {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.
		Every(config.Get().ServerInfoUpdateInterval).Seconds().
		Do(UpdateServersInfo)

	scheduler.StartAsync()
}

func Shutdown() {
	scheduler.Stop()
}

func UpdateServersInfo() {
	log.Info("checking for updates")
	infoMap.Range(func(key string, value interface{}, ttl int64) {
		log.Info("updating server info for: " + key)
		si := value.(*ServerInfo)
		go ObtainServerInfo(si)
	})
}

func ObtainServerInfo(s *ServerInfo) {
	log.Debug("obtaining info from server")
	var wg sync.WaitGroup
	wg.Add(3)
	go FillInfo(packets.A2sInfo, s, &wg)
	go FillInfo(packets.A2sPlayers, s, &wg)
	go FillInfo(packets.A2sRules, s, &wg)
	wg.Wait()
	infoMap.Set(s.Server.String(), s, getTTL())
}

func FillInfo(qt packets.QueryType, serverInfo *ServerInfo, wg *sync.WaitGroup) {
	response, err := ConnectAndQuery(serverInfo.Server, qt)
	if err == nil {
		serverInfo.AddInfo(qt, response)
	}
	wg.Done()
}
