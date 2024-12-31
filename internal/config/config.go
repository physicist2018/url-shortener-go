package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr    string
	BaseURLServer string
}

var DefaultConfig Config

func NewConfig() *Config {
	return &Config{
		ServerAddr:    "localhost:8080",
		BaseURLServer: "http://localhost:8000",
	}
}

func MakeConfig(serverAddr string, baseURLServer string) *Config {
	return &Config{
		ServerAddr:    serverAddr,
		BaseURLServer: baseURLServer,
	}
}

func ConfigApp() *Config {
	flag.StringVar(&DefaultConfig.ServerAddr, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&DefaultConfig.BaseURLServer, "b", "http://localhost:8080", "префикс короткого URL")
	flag.Parse()
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		DefaultConfig.ServerAddr = envServerAddr
	}

	if envBaseURLServer := os.Getenv("BASE_URL"); envBaseURLServer != "" {
		DefaultConfig.BaseURLServer = envBaseURLServer
	}
	return &DefaultConfig
}
