package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer  HTTPServer
	DbURL string `env:"DATABASE_URL"`
}

type HTTPServer struct {
	Address string        `env:"ADDRESS" env-default:"localhost:8080"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"4"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60"`
}

const configPath = "config/config.env"

func MustLoad() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}