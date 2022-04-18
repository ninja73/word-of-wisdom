package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type ServerSetting struct {
	Address     string
	TTL         int64
	BitStrength int32
	Timeout     time.Duration
	SecretKey   string
}

type CacheSetting struct {
	CleanupInterval time.Duration
	CacheTTL        int64
}

type Config struct {
	StoreFile     string
	CacheSetting  CacheSetting
	ServerSetting ServerSetting
}

func ParseConfig(configFile string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
