package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"wow/config"
	"wow/internal/client"
	"wow/internal/logger"
)

func main() {
	configFile := flag.String("config", "./client-config.toml", "config file")

	flag.Parse()

	logger.InitLogger(os.Stderr, log.InfoLevel)

	cfg, err := config.ParseClientConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < cfg.Clients; i++ {
		go func() {
			cln := client.NewClient(cfg.ServerAddress, cfg.Timeout.Duration)
			for {
				quote, err := cln.GetQuote()
				if err != nil {
					log.Error(err)
					continue
				}

				log.Infof("Quote: %s", quote)
			}
		}()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
