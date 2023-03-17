package domain

import (
	"fmt"
)

type GameServer struct {
	IP     string
	Port   string
	Engine EngineType
}

func (g *GameServer) String() string {
	return fmt.Sprintf("%s:%s", g.IP, g.Port)
}
