package gameserverinfo

import (
	"net"
	"time"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"
	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/internal/utils"
	"github.com/Th3Fr33m4n/source-engine-query-cache/packets"
	"golang.org/x/exp/maps"
)

var challenges = utils.NewConcurrentMap()

func ConnectAndQuery(sv config.GameServer, qt packets.QueryType) ([][]byte, error) {
	conn, err := connect(sv)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	log.Debug("sending query")
	chKey := sv.String()
	ch := challenges.Get(chKey)
	var q []byte

	if ch == nil {
		// no challenge set for this client
		q = packets.GetQuery(qt)
	} else {
		// a challenge for this client has already been set, use it
		q = packets.BuildQuery(qt, ch.([]byte))
	}

	response, err := sendAndGet(conn, q)
	if err != nil {
		return nil, err
	}

	cat := packets.CategorizeResponse(response)

	if cat == packets.A2sChallengeResponse {
		log.Debug("server response is a challenge")
		ch = packets.GetChallengeFromServerResponse(response)
		challenges.Set(chKey, ch)
		// increase deadline because this is a second request
		addDeadline(conn)
		q = packets.BuildQuery(qt, ch.([]byte))
		log.Debug("sending query with challenge")
		response, err = sendAndGet(conn, q)
		if err != nil {
			return nil, err
		}
		cat = packets.CategorizeResponse(response)
	}

	if cat == packets.A2sRulesSplitResponse {
		mpResponse, err := handleSplitPacket(conn, response, sv.Engine)
		if err != nil {
			return nil, err
		}
		return mpResponse, nil
	} else if cat == packets.GetQueryResponseType(qt) {
		log.Debug("response matches expected structure")
		log.Debug(string(response))
		return [][]byte{response}, nil
	} else {
		log.Error("invalid server response")
		log.Debug(string(response))
		return nil, ErrInvalidResponse
	}
}

func connect(g config.GameServer) (*net.UDPConn, error) {
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

func handleSplitPacket(conn *net.UDPConn, msg []byte, engine domain.EngineType) ([][]byte, error) {
	responsePackets := make(map[byte][]byte)
	pNum, totalPackets := getPacketCount(msg, engine)
	responsePackets[pNum] = msg
	var addedPackets byte = 1
	for addedPackets < totalPackets {
		l1 := len(responsePackets)
		p := make([]byte, config.Get().ReadBufferSize)
		addDeadline(conn)
		rb, err := conn.Read(p)
		if err != nil {
			return nil, err
		}
		pNum, _ = getPacketCount(p, engine)
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
		return packets.ParseGoldsrcMultipacketResponse(msg)
	} else {
		return packets.ParseSourceMultipacketResponse(msg)
	}
}
