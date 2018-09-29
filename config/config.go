package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Server   string
	Database string
}

func (conf *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", conf); err != nil {
		log.Fatal(err)
	}
}
