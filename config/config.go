package config

import "github.com/BurntSushi/toml"

type Config struct {
	Storage StorageConfig `toml:"storage"`
	Server  ServerConfig  `toml:"server"`
}

type ServerConfig struct {
	Port         string `toml:"port"`
	ReadTimeout  int    `toml:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout"`
}

type StorageConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	DBName   string `toml:"db_name" `
	Password string `toml:"password"`
}

func ParseConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	return &config, err
}
