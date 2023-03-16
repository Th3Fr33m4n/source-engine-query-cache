package scraper

import (
	"testing"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain/a2s"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
)

func TestGetSomething(t *testing.T) {
	config.Init()
	sv := domain.GameServer{IP: "216.52.148.19", Port: "27016", Engine: "goldsrc"}
	response, err := ConnectAndQuery(&QueryContext{Sv: sv, A2sQ: a2s.InfoQuery})
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
