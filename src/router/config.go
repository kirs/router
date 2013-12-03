package router

import (
	"github.com/BurntSushi/toml"
	"log"
)

type tomlConfig struct {
	Routers []configRouter
}

type configRouter struct {
	Alias     string
	Source    string
	LaunchCmd string
	RailsPort string
}

func ReadConfig(path string) tomlConfig {
	var config tomlConfig

	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
