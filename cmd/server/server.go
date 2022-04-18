package main

import (
	"flag"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"wow/config"
	"wow/internal/cache"
	"wow/internal/server"
	"wow/internal/store"
)

func main() {
	logFile := flag.String("log", "./out.txt", "log file")
	configFile := flag.String("config", "config.toml", "config file")
	flag.Parse()

	file, err := os.Open(*logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	initLogger(file)

	cfg, err := config.ParseConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	fileStore, err := store.NewFileStore(cfg.StoreFile)
	if err != nil {
		log.Fatal(err)
	}

	mCache := cache.NewInMemoryCache(cfg.CacheSetting.CleanupInterval, cfg.CacheSetting.CacheTTL)

	srv := server.NewServer(
		fileStore,
		mCache,
		server.WithBitStrength(cfg.ServerSetting.BitStrength),
		server.WithTimeout(cfg.ServerSetting.Timeout),
		server.WithSecretKey(cfg.ServerSetting.SecretKey),
	)

	if err := srv.TCPListen(cfg.ServerSetting.Address); err != nil {
		log.Fatal(err)
	}
}

func initLogger(logFile io.Writer) {
	log.SetOutput(logFile)
	log.SetReportCaller(true)

	log.SetFormatter(&nested.Formatter{})
}
