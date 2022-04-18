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
	configFile := flag.String("config", "./server-config.toml", "config file")

	flag.Parse()

	logger.InitLogger(os.Stderr, log.InfoLevel)

	cfg, err := config.ParseServerConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	quoteStore, err := store.NewFileStore(cfg.StoreFile)
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := redis.InitRedis(
		cfg.CacheRedis.Host,
		cfg.CacheRedis.Port,
		cfg.CacheRedis.Password,
		cfg.CacheRedis.DB,
		cfg.CacheRedis.PoolSize,
	)
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
			Limit:       cfg.Server.RateLimit,
		},
	)

	if err := srv.TCPListen(cfg.Server.Address); err != nil {
		log.Fatal(err)
	}
}
