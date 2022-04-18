package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"wow/config"
	"wow/internal/cache"
	"wow/internal/logger"
	"wow/internal/redis"
	"wow/internal/server"
	"wow/internal/store"
)

func main() {
	logFilePath := flag.String("log", "./out.txt", "log file")
	configFile := flag.String("config", "config.toml", "config file")
	flag.Parse()

	logFile, err := os.Open(*logFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	if err := logger.InitLogger(logFile); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	quoteStore, err := store.NewFileStore(cfg.StoreFile)
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := redis.InitRedis(cfg.CacheRedis.Host, cfg.CacheRedis.Port, cfg.CacheRedis.Password, cfg.CacheRedis.DB)
	if err != nil {
		log.Fatal(err)
	}

	powCache := cache.NewRedisCache(rdb, cfg.Cache.Expiration.Duration)

	srv := server.NewServer(
		quoteStore,
		powCache,
		&server.Options{
			BitStrength: cfg.Server.BitStrength,
			Timeout:     cfg.Server.Timeout.Duration,
			SecretKey:   cfg.Server.SecretKey,
			Expiration:  cfg.Server.Expiration.Duration,
			Limit:       cfg.Server.Limit,
		},
	)

	if err := srv.TCPListen(cfg.Server.Address); err != nil {
		log.Fatal(err)
	}
}
