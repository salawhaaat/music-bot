package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Token string
	Host  string
	Dsn   string
	Addr  string
}

func New() *Config {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse() // -config=config/local.toml
	cfg := &Config{}
	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}
	return cfg
}
