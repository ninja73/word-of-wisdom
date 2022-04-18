package config

import (
	"github.com/BurntSushi/toml"
	"time"
)

type Duration struct{ time.Duration }

type Server struct {
	Address     string   `toml:"address"`
	BitStrength int32    `toml:"bit-strength"`
	SecretKey   string   `toml:"secret-key"`
	Timeout     Duration `toml:"timeout"`
	Expiration  Duration `toml:"expiration"`
	RateLimit   int32    `toml:"rate-limit"`
}

type Redis struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	PoolSize int    `toml:"pool-size"`
}

type Cache struct {
	Expiration Duration `toml:"expiration"`
}

type ServerConfig struct {
	StoreFile  string `toml:"store-file"`
	Cache      Cache  `toml:"cache"`
	Server     Server `toml:"server"`
	CacheRedis Redis  `toml:"cache-redis"`
}

func ParseServerConfig(configFile string) (*ServerConfig, error) {
	var config ServerConfig
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
