package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	PostgresConfig PostgresConfig `yaml:"PostgresConfig"`
}

func MustLoad(configPath string) (config Config) {
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set!")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file is not exist %s", configPath)
	}
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	return
}
