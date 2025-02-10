package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/handler"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/compressor"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/httplogger"
	"github.com/physicist2018/url-shortener-go/internal/repository/repofactorymethod"
	"github.com/physicist2018/url-shortener-go/internal/service"
	randomstringgenerator "github.com/physicist2018/url-shortener-go/pkg/randomstring_generator"
	"github.com/rs/zerolog"
)

func main() {
	var err error
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	logger.Info().Msg("конфигурирование сервера")

	cfg, err := config.Load()
	if err != nil {
		logger.Error().Err(err).Send()
		panic(err)
	}

	logger.Info().Msg(cfg.String())

	logger.Info().Msg("инициализация генератора случайных ссылок")
	randomStringGenerator := randomstringgenerator.NewRandomStringDefault()

	repofactory := repofactorymethod.NewRepofactorymethod()
	var linkRepo domain.URLLinkRepo

	// Надо реализовать InJsonFile storage
	if cfg.DatabaseDSN != "" {
		linkRepo, err = repofactory.CreateRepo("postgres", cfg.DatabaseDSN)
	} else if cfg.FileStoragePath != "" {
		linkRepo, err = repofactory.CreateRepo("inmemory", cfg.FileStoragePath)
	} else {
		linkRepo, err = repofactory.CreateRepo("inmemory", cfg.FileStoragePath)
	}

	//inmemory.NewInMemoryLinkRepository(cfg.FileStoragePath)
	if err != nil {
		logger.Error().Err(err).Send()
		panic(err)
	}
	defer linkRepo.Close()

	linkService := service.NewURLLinkService(linkRepo, randomStringGenerator)
	linkHandler := handler.NewURLLinkHandler(linkService, cfg.BaseURLServer)

	r := chi.NewRouter()
	r.Use(httplogger.LoggerMiddleware(&logger))
	r.Use(compressor.RequestDecompressionMiddleware)
	r.Use(compressor.ResponseCompressionMiddleware(compressor.BestCompression))
	r.Use(middleware.AllowContentType("text/plain", "application/json", "text/html", "application/x-gzip"))
	r.Use(middleware.Recoverer)

	///
	r.Post("/", linkHandler.ShortenURL)
	r.Post("/api/shorten", linkHandler.HandleGenerateShortURLJson)
	r.Post("/api/shorten/batch", linkHandler.HandleGenerateShortURLJsonBatch)
	r.Get("/{shortURL}", linkHandler.Redirect)
	r.Get("/ping", linkHandler.PingHandler)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}
	logger.Info().Msg("Запуск сервера")
	log.Fatal(server.ListenAndServe())
}
