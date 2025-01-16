package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr    string
	BaseURLServer string
}

func (c *Config) Parse() {
	flag.Parse()
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		c.ServerAddr = envServerAddr
	}

	if envBaseURLServer := os.Getenv("BASE_URL"); envBaseURLServer != "" {
		c.BaseURLServer = envBaseURLServer
	}
}

func MakeConfig() *Config {
	cfg := &Config{
		ServerAddr:    "localhost:8080",
		BaseURLServer: "http://localhost:8000",
	}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&cfg.BaseURLServer, "b", "http://localhost:8080", "префикс короткого URL")

	return cfg
}
