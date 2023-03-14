package handlers

import (
	"bytes"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/internal/gameserverinfo"
	"github.com/Th3Fr33m4n/source-engine-query-cache/packets"
	"github.com/Th3Fr33m4n/source-engine-query-cache/public/challenge"
)

type A2sQueryContext struct {
	Conn         net.PacketConn
	Addr         net.Addr
	Query        []byte
	Sv           config.GameServer
	GetChallenge func([]byte) []byte
	QType        packets.QueryType
}

func SendChallenge(udpServer net.PacketConn, addr net.Addr) {
	ch := challenge.GetForClient(addr.String())
	udpServer.WriteTo(packets.BuildChallengeResponse(ch), addr)
}

func A2sQueryHandler(ctx A2sQueryContext) {
	clch := challenge.GetForClient(ctx.Addr.String())
	ch := ctx.GetChallenge(ctx.Query)

	if bytes.Equal(clch, ch) {
		si, err := gameserverinfo.GetServerInfo(ctx.Sv)
		if err != nil {
			log.Println(err.Error())
		}
		v := si.GetInfo(packets.A2sRules)
		if v != nil {
			log.Println(v)
			for _, p := range v {
				ctx.Conn.WriteTo(p, ctx.Addr)
			}
		}
	} else {
		SendChallenge(ctx.Conn, ctx.Addr)
	}
}
