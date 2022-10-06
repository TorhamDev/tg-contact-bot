package models

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Bot struct {
		Token      string `yaml:"token"`
		PollerTime int    `yaml:"poller_time"`
		DbName string `yaml:"db_name"`
	}
}

func LoadConfigs() Config {

	ConfigFile, err := os.Open("./config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	defer ConfigFile.Close()

	var cfg Config

	decoder := yaml.NewDecoder(ConfigFile)

	err = decoder.Decode(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
