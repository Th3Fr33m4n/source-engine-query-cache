package main

import (
	"os"
	"sync"

	"github.com/Th3Fr33m4n/source-engine-query-cache/domain"

	log "github.com/sirupsen/logrus"

	"github.com/Th3Fr33m4n/source-engine-query-cache/config"
	"github.com/Th3Fr33m4n/source-engine-query-cache/internal/scraper"
	"github.com/Th3Fr33m4n/source-engine-query-cache/public/listeners"
)

func main() {
	displayBanner()
	log.Info("Source engine query cache starting" +
		"g up...")
	config.Init()
	log.Info("Config loaded")
	setUpLogger()

	for _, sv := range config.Get().Servers {
		if sv.Engine == domain.Source || sv.Engine == domain.GoldSrc {
			scraper.RegisterServer(sv)
		} else {
			log.Panicf("Invalid engine type for server %s", sv.String())
		}
	}

	log.Info("Starting background game server info scraper...")
	scraper.Init()
	defer scraper.Shutdown()

	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Info("Starting listener interfaces...")
	for _, g := range config.Get().Servers {
		go listeners.Listen(g)
	}

	wg.Wait()
}

func setUpLogger() {
	lvl, err := log.ParseLevel(config.Get().LogLevel)
	if err != nil {
		lvl = log.ErrorLevel
	}
	log.SetLevel(lvl)
}

func displayBanner() {
	dat, err := os.ReadFile("./banner.txt")
	if err != nil {
		return
	}
	println(string(dat))
}
