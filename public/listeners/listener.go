package listeners

import (
	"net"

	"github.com/Th3Fr33m4n/source-engine-query-cache/public/handlers"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/packets"
	"github.com/Th3Fr33m4n/source-engine-query-cache/ratelimit"
)

func Listen(g config.GameServer) {
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

func response(conn net.PacketConn, addr net.Addr, buf []byte, sv config.GameServer) {
	glbrtl := ratelimit.GetGlobalLimiter()
	clrtl := ratelimit.GetLimiterForAddress(addr.String())

	if !glbrtl.Allow() || !clrtl.Allow() {
		conn.WriteTo(packets.TooManyRequests, addr)
		return
	}

	if packets.IsA2sInfoRequest(buf) || packets.IsA2sPlayersRequest(buf) || packets.IsA2sRulesRequest(buf) {
		handlers.SendChallenge(conn, addr)
	} else if packets.IsA2sInfoWChallenge(buf) {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, Query: buf, Sv: sv,
			GetChallenge: packets.GetChallengeFromA2sInfo,
			QType:        packets.A2sInfo,
		})
	} else if packets.IsA2sPlayersWChallenge(buf) {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, Query: buf, Sv: sv,
			GetChallenge: packets.GetChallengeFromA2sPlayers,
			QType:        packets.A2sPlayers,
		})
	} else if packets.IsA2sRulesWChallenge(buf) {
		handlers.A2sQueryHandler(handlers.A2sQueryContext{
			Conn: conn, Addr: addr, Query: buf, Sv: sv,
			GetChallenge: packets.GetChallengeFromA2sRules,
			QType:        packets.A2sRules,
		})
	}
}
