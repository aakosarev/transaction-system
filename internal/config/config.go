package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	HTTP struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
	} `yaml:"http"`
	PostgreSQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"postgresql"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
