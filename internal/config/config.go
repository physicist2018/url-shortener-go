package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerAddr         string
	BaseURLServer      string
	FileStoragePath    string
	MaxShortURLLength  int
	MaxShutdownTime    int
	MaxTimeBetweenSync int
}

const (
	maxShortURLLen     = 5
	maxShutdownTime    = 5
	maxTimeBetweenSync = 120
)

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
	flag.IntVar(&cfg.MaxShortURLLength, "max-short-url-len", maxShortURLLen, "максимально допустимая длина короткой ссылки")
	flag.IntVar(&cfg.MaxShutdownTime, "max-shutdown-time", maxShutdownTime, "время в секундах, кторое мы ждем прежде чем прекратим выключать сервер")
	flag.IntVar(&cfg.MaxTimeBetweenSync, "max-time-between-sync-db", maxTimeBetweenSync, "интервал синхронизации бд в секундах")
	return cfg
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"ServerAddr: %s\nBaseURLServer: %s\nFileStoragePath: %s\nMaxShortURLLength: %d\nMaxShutdownTime: %d\nMaxTimeBetweenSync: %d",
		c.ServerAddr,
		c.BaseURLServer,
		c.FileStoragePath,
		c.MaxShortURLLength,
		c.MaxShutdownTime,
		c.MaxTimeBetweenSync,
	)
}
