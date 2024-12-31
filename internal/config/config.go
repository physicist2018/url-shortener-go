package config

import "flag"

type Config struct {
	ServerAddr    string
	BaseURLServer string
}

var DefaultConfig Config

func NewConfig() *Config {
	return &Config{
		ServerAddr:    "localhost:8080",
		BaseURLServer: "http://localhost:8000/qsd54gFg",
	}
}

func MakeConfig(serverAddr string, baseURLServer string) *Config {
	return &Config{
		ServerAddr:    serverAddr,
		BaseURLServer: baseURLServer,
	}
}

func init() {
	flag.StringVar(&DefaultConfig.ServerAddr, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&DefaultConfig.BaseURLServer, "b", "http://localhost:8080/aaa", "префикс короткого URL")
}
