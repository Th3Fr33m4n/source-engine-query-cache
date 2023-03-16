package scraper

import (
	"sync"
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain/a2s"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/go-co-op/gocron"
	ttl_map "github.com/leprosus/golang-ttl-map"
)

var (
	infoMap   = ttl_map.New()
	scheduler *gocron.Scheduler
)

type ServerInfo struct {
	Info    [][]byte
	Rules   [][]byte
	Players [][]byte
	Server  domain.GameServer
}

func (s *ServerInfo) AddInfo(qt domain.QueryType, info [][]byte) {
	switch qt {
	case domain.A2sInfo:
		s.Info = info
	case domain.A2sRules:
		s.Rules = info
	case domain.A2sPlayers:
		s.Players = info
	}
}

func (s *ServerInfo) GetInfo(qt domain.QueryType) [][]byte {
	var res [][]byte
	switch qt {
	case domain.A2sInfo:
		res = s.Info
	case domain.A2sRules:
		res = s.Rules
	case domain.A2sPlayers:
		res = s.Players
	}
	return res
}

func getTTL() int64 {
	return int64(3 * config.Get().ServerInfoUpdateInterval)
}

func RegisterServer(g domain.GameServer) {
	infoMap.Set(g.String(), &ServerInfo{Server: g}, getTTL())
}

func GetServerInfo(g domain.GameServer) (*ServerInfo, error) {
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
	go FillInfo(a2s.InfoQuery, s, &wg)
	go FillInfo(a2s.PlayersQuery, s, &wg)
	go FillInfo(a2s.RulesQuery, s, &wg)
	wg.Wait()
	infoMap.Set(s.Server.String(), s, getTTL())
}

func FillInfo(a2sq a2s.A2sQuery, serverInfo *ServerInfo, wg *sync.WaitGroup) {
	ctx := &QueryContext{A2sQ: a2sq, Sv: serverInfo.Server}
	response, err := ConnectAndQuery(ctx)
	if err == nil {
		serverInfo.AddInfo(a2sq.QueryT, response)
	}
	wg.Done()
}
