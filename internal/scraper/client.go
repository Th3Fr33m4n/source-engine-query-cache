package scraper

import (
	"net"
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
	"github.com/Th3Fr33m4n/source-engine-query-cache/domain/a2s"
	"github.com/Th3Fr33m4n/source-engine-query-cache/internal/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type QueryContext struct {
	Sv           domain.GameServer
	A2sQ         a2s.Query
	conn         *net.UDPConn
	lastResponse []byte
}

var challenges = utils.NewConcurrentMap()

func ConnectAndQuery(ctx *QueryContext) ([][]byte, error) {
	var err error
	ctx.conn, err = connect(ctx.Sv)
	if err != nil {
		return nil, err
	}

	defer ctx.conn.Close()

	log.Debug("sending query")
	ch := challenges.Get(ctx.Sv.String())
	var q []byte

	if ch == nil {
		// no challenge set for this client
		q = ctx.A2sQ.Body
	} else {
		// a challenge for this client has already been set, use it
		q = ctx.A2sQ.Build(ch.([]byte))
	}

	ctx.lastResponse, err = sendAndGet(ctx.conn, q)
	if err != nil {
		return nil, err
	}

	cat := ctx.A2sQ.MatchResponse(ctx.lastResponse)

	if cat == domain.A2sChallengeResponse {
		cat, err = addChallengeAndResend(ctx)
		if err != nil {
			return nil, err
		}
	}

	if cat == domain.A2sRulesSplitResponse {
		mpResponse, err := handleSplitPacket(ctx)
		if err != nil {
			return nil, err
		}
		return mpResponse, nil
	} else if cat != domain.InvalidResponse {
		log.Debug("response matches expected structure")
		log.Debug(string(ctx.lastResponse))
		return [][]byte{ctx.lastResponse}, nil
	} else {
		log.Error("invalid server response")
		log.Debug(string(ctx.lastResponse))
		return nil, ErrInvalidResponse
	}
}

func connect(g domain.GameServer) (*net.UDPConn, error) {
	udpServer, err := net.ResolveUDPAddr("udp", g.String())
	if err != nil {
		log.Errorf("resolveUDPAddr failed: %v", err.Error())
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		log.Errorf("listen failed: %v", err.Error())
		return nil, err
	}
	addDeadline(conn)
	return conn, nil
}

func sendAndGet(conn *net.UDPConn, msg []byte) ([]byte, error) {
	_, err := conn.Write(msg)
	if err != nil {
		log.Errorf("Write data failed: %v", err.Error())
		return nil, err
	}
	received := make([]byte, config.Get().ReadBufferSize)
	n, err := conn.Read(received)
	if err != nil {
		log.Errorf("Read data failed: %v", err.Error())
		return nil, err
	}
	return received[:n], nil
}

func addDeadline(conn *net.UDPConn) {
	conn.SetDeadline(time.Now().Add(config.Get().ServerInfoUpdateTimeout))
}

func addChallengeAndResend(ctx *QueryContext) (domain.ResponseType, error) {
	log.Debug("server response is a challenge")
	ch := ctx.A2sQ.GetChallengeFromResponse(ctx.lastResponse)
	challenges.Set(ctx.Sv.String(), ch)
	// increase deadline because this is a second request
	addDeadline(ctx.conn)
	q := ctx.A2sQ.Build(ch)
	log.Debug("sending query with challenge")
	var err error
	ctx.lastResponse, err = sendAndGet(ctx.conn, q)
	if err != nil {
		return domain.InvalidResponse, err
	}
	return ctx.A2sQ.MatchResponse(ctx.lastResponse), nil
}

func handleSplitPacket(ctx *QueryContext) ([][]byte, error) {
	responsePackets := make(map[byte][]byte)
	pNum, totalPackets := getPacketCount(ctx.lastResponse, ctx.Sv.Engine)
	responsePackets[pNum] = ctx.lastResponse
	var addedPackets byte = 1
	for addedPackets < totalPackets {
		l1 := len(responsePackets)
		p := make([]byte, config.Get().ReadBufferSize)
		addDeadline(ctx.conn)
		rb, err := ctx.conn.Read(p)
		if err != nil {
			return nil, err
		}
		pNum, _ = getPacketCount(p, ctx.Sv.Engine)
		responsePackets[pNum] = p[:rb]
		l2 := len(responsePackets)
		if l1 != l2 {
			addedPackets++
		}
	}
	return maps.Values(responsePackets), nil
}

func getPacketCount(msg []byte, engine domain.EngineType) (byte, byte) {
	if engine == domain.GoldSrc {
		return a2s.ParseGoldsrcMultipacketResponse(msg)
	} else {
		return a2s.ParseSourceMultipacketResponse(msg)
	}
}
