package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerAddr        string
	BaseURLServer     string
	FileStoragePath   string
	DatabaseDSN       string
	MaxShortURLLength int
	MaxShutdownTime   int
}

func NewConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&cfg.BaseURLServer, "b", "http://localhost:8080", "префикс короткого URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "dbase.json", "имя файла персистентного хранилища коротких URL")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "параметры подключения к базе данных")
	flag.IntVar(&cfg.MaxShortURLLength, "max-short-url-len", 5, "максимально допустимая длина короткой ссылки")
	flag.IntVar(&cfg.MaxShutdownTime, "max-shutdown-time", 5, "время в секундах, кторое мы ждем прежде чем прекратим выключать сервер")
	return cfg
}

func Load() (*Config, error) {
	cfg := NewConfig()
	cfg.Parse()
	return cfg, nil
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

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		c.DatabaseDSN = envDatabaseDSN
	}
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"ServerAddr: %s, \nBaseURLServer: %s, \nFileStoragePath: %s, \nDatabaseDSN: %s, \nMaxShortURLLength: %d, \nMaxShutdownTime: %d",
		c.ServerAddr,
		c.BaseURLServer,
		c.FileStoragePath,
		c.DatabaseDSN,
		c.MaxShortURLLength,
		c.MaxShutdownTime,
	)
}
