package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/physicist2018/url-shortener-go/internal/adapters/urlshortenerserver"
	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/rs/zerolog"
)

func main() {
	// Инициализация конфигурации
	configuration := config.MakeConfig()
	configuration.Parse()

	// Инициализация логгера

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	//logger, err := zap.NewDevelopment()
	//if err != nil {
	//	panic(err)
	//}
	//defer logger.Sync()
	//sugar := logger.Sugar()

	// Инициализация сервера
	server := urlshortenerserver.NewURLShortenerServer(configuration, &logger)

	// Запуск сервера
	server.Start()

	// Ожидаем сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Завершаем сервер
	server.Shutdown()

}
