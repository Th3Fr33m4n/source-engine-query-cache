package handlers

import (
	"bytes"
	"net"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
	"github.com/Th3Fr33m4n/source-engine-query-cache/domain/a2s"
	"github.com/Th3Fr33m4n/source-engine-query-cache/internal/scraper"
	"github.com/Th3Fr33m4n/source-engine-query-cache/public/challenge"
	log "github.com/sirupsen/logrus"
)

type A2sQueryContext struct {
	Conn     net.PacketConn
	Addr     net.Addr
	RawQuery []byte
	A2sq     a2s.A2sQuery
	Sv       domain.GameServer
}

func SendChallenge(udpServer net.PacketConn, addr net.Addr) {
	ch := challenge.GetForClient(addr.String())
	udpServer.WriteTo(a2s.BuildChallengeResponse(ch), addr)
}

func A2sQueryHandler(ctx A2sQueryContext) {
	clch := challenge.GetForClient(ctx.Addr.String())
	ch := ctx.A2sq.GetChallengeFromRequest(ctx.RawQuery)

	if bytes.Equal(clch, ch) {
		si, err := scraper.GetServerInfo(ctx.Sv)
		if err != nil {
			log.Println(err.Error())
		}
		v := si.GetInfo(ctx.A2sq.QueryT)
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
