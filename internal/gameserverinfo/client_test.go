package gameserverinfo

import (
	"testing"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/packets"
)

func TestGetSomething(t *testing.T) {
	config.Init()
	sv := config.GameServer{IP: "216.52.148.19", Port: "27016", Engine: "goldsrc"}
	response, err := ConnectAndQuery(sv, packets.A2sInfo)
	if err != nil {
		println(err)
		t.Error(err)
	}
	println(string(response[0]))
}

func TestBits(t *testing.T) {
	var n1 byte = 0x02

	n2 := n1 & 0xf0
	n3 := n1 & 0x0f
	n2 = n2 >> 4

	println(n2)
	println(n3)
}

/*
7 6 5 4 3 2 1 0
0 0 0 1 0 0 1 0  => 0x12
&
1 1 1 1 0 0 0 0 => 0xF0
________________________
0 0 0 1 0 0 0 0 => 0x10


&
0 0 0 0 1 1 1 1 => 0x0F
________________________
0 0 0 0 0 0 1 0


16 2

8 4 2
*/
