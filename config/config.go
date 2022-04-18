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
	Limit       int32    `toml:"limit"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Cache struct {
	Expiration Duration `toml:"expiration"`
}

type Config struct {
	StoreFile  string `toml:"store-file"`
	Cache      Cache  `toml:"cache"`
	Server     Server `toml:"server"`
	CacheRedis Redis  `toml:"cache-redis"`
}

func ParseConfig(configFile string) (*Config, error) {
	var config Config
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
