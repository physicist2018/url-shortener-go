package main

import (
	"log"
	"net/http"
	"time"

	"math/rand"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/physicist2018/url-shortener-go/internal/adapters/httphandlers"
	"github.com/physicist2018/url-shortener-go/internal/adapters/memory"
	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/core/services/url"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/compressor"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/httplogger"
	"github.com/physicist2018/url-shortener-go/pkg/utils"
	"go.uber.org/zap"
)

const maxShortURLLength = 6

func main() {
	configuration := config.MakeConfig()
	configuration.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	randomStringGenerator := utils.NewRandomString(maxShortURLLength, rand.New(rand.NewSource(time.Now().UnixNano())))
	urlRepo := memory.NewURLRepositoryMap()
	urlService := url.NewURLService(urlRepo, randomStringGenerator)
	urlHandler := httphandlers.NewURLHandler(urlService, configuration.BaseURLServer)

	r := chi.NewRouter()
	//r.Use(compressor.CompressionMiddleware())
	r.Use(middleware.AllowContentType("text/plain", "application/json", "text/html", "application/x-gzip"))
	r.Use(httplogger.LoggerMiddleware(sugar))
	r.Post("/", urlHandler.HandleGenerateShortURL)
	r.Post("/api/shorten", compressor.GzipMiddleware(urlHandler.HandleGenerateShortURLJson))
	r.Get("/{shortURL}", urlHandler.HandleRedirect)

	if err := http.ListenAndServe(configuration.ServerAddr, r); err != nil {
		log.Fatal(err)
	}
}
