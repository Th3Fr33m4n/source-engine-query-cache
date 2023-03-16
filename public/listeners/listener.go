package listeners

import (
	"net"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
	"github.com/Th3Fr33m4n/source-engine-query-cache/domain/a2s"
	"github.com/Th3Fr33m4n/source-engine-query-cache/public/handlers"
	"github.com/Th3Fr33m4n/source-engine-query-cache/public/ratelimit"
	log "github.com/sirupsen/logrus"
)

func Listen(g domain.GameServer) {
	// listen to incoming udp packets
	p := ":" + g.Port
	udpServer, err := net.ListenPacket("udp", p)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("listening on " + p)

	defer udpServer.Close()

	for {
		buf := make([]byte, config.Get().ReadBufferSize)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}

		go response(udpServer, addr, buf[:n], g)
	}
}

func response(conn net.PacketConn, addr net.Addr, req []byte, sv domain.GameServer) {
	glbrtl := ratelimit.GetGlobalLimiter()
	clrtl := ratelimit.GetLimiterForAddress(addr.String())

	if !glbrtl.Allow() || !clrtl.Allow() {
		conn.WriteTo(a2s.TooManyRequests, addr)
		return
	}
	cat, hasChallenge := a2s.CategorizeRequest(req)

	if cat == domain.InvalidQuery {
		return
	}

	if !hasChallenge {
		handlers.SendChallenge(conn, addr)
	} else if cat == domain.A2sInfo {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, RawQuery: req, Sv: sv, A2sq: a2s.InfoQuery,
		})
	} else if cat == domain.A2sPlayers {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, RawQuery: req, Sv: sv, A2sq: a2s.PlayersQuery,
		})
	} else if cat == domain.A2sRules {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, RawQuery: req, Sv: sv, A2sq: a2s.RulesQuery,
		})
	}
}
