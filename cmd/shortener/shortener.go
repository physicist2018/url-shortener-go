package main

import (
	"context"
	"os"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/deleter"
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/handler"
	"github.com/physicist2018/url-shortener-go/internal/repository/repofactorymethod"
	"github.com/physicist2018/url-shortener-go/internal/router"
	"github.com/physicist2018/url-shortener-go/internal/server"
	"github.com/physicist2018/url-shortener-go/internal/service"
	stringgenstategy "github.com/physicist2018/url-shortener-go/internal/stringgenstrategy"
	uniquestring "github.com/physicist2018/url-shortener-go/pkg/uniquestring"
	"github.com/rs/zerolog"
)

func main() {
	var err error
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	logger.Info().Msg("конфигурирование сервера")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка загрузки конфигурации")
	}

	logger.Info().Msg(cfg.String())

	logger.Info().Msg("инициализация генератора случайных ссылок")
	randomStringStrategy := uniquestring.NewRandomStringDefault()

	stringGeneratorContext := stringgenstategy.StringGeneratorContext{}
	stringGeneratorContext.SetStrategy(randomStringStrategy)

	repofactory := repofactorymethod.NewRepoFactoryMethod()
	var linkRepo domain.URLLinkRepo

	if cfg.DatabaseDSN != "" {
		linkRepo, err = repofactory.CreateRepo("postgres", cfg.DatabaseDSN)
	} else {
		linkRepo, err = repofactory.CreateRepo("inmemory", cfg.FileStoragePath)
	}

	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка инициализации репозитория")
	}
	defer func() {
		if err := linkRepo.Close(); err != nil {
			logger.Error().Err(err).Msg("Ошибка при закрытии репозитория")
		}
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	linkService := service.NewURLLinkService(linkRepo, stringGeneratorContext, logger)
	linkDeleter := deleter.NewDeleter(linkService, logger)
	linkDeleter.Start(ctx, &wg) //Запускаем горутину асинхронного удаления ссылок

	linkHandler := handler.NewURLLinkHandler(linkService, cfg.BaseURLServer, logger, linkDeleter)

	r := router.NewRouter(linkHandler, logger)

	srv := server.NewServer(cfg.ServerAddr, r, logger)
	srv.Start()

	linkHandler.Close() // Закрываем канал обмена с горутиной, что приводит к очистке очереди и завершению
	logger.Info().Msg("Closing link handler")
	wg.Wait()
}
