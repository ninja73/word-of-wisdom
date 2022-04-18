package config

import "github.com/BurntSushi/toml"

type ClientConfig struct {
	ServerAddress string   `toml:"server-address"`
	Timeout       Duration `toml:"timeout"`
	Clients       int      `toml:"clients"`
}

func ParseClientConfig(configFile string) (*ClientConfig, error) {
	var config ClientConfig
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
