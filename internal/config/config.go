package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr        string
	BaseURLServer     string
	FileStoragePath   string
	MaxShortURLLength int
	MaxShutdownTime   int
}

func (c *Config) Parse() {
	flag.Parse()
	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		c.ServerAddr = envServerAddr
	}

	if envBaseURLServer := os.Getenv("BASE_URL"); envBaseURLServer != "" {
		c.BaseURLServer = envBaseURLServer
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		c.FileStoragePath = envFileStoragePath
	}
}

func MakeConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&cfg.BaseURLServer, "b", "http://localhost:8080", "префикс короткого URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "dbase.json", "имя файла персистентного хранилища коротких URL")
	flag.IntVar(&cfg.MaxShortURLLength, "max-short-url-len", 5, "максимально допустимая длина короткой ссылки")
	flag.IntVar(&cfg.MaxShutdownTime, "max-shutdown-time", 5, "время в секундах, кторое мы ждем прежде чем прекратим выключать сервер")
	return cfg
}
