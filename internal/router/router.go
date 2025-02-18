package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/physicist2018/url-shortener-go/internal/handler"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/compressor"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/httplogger"
	"github.com/rs/zerolog"
)

func NewRouter(linkHandler *handler.URLLinkHandler, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Мидлвары
	r.Use(httplogger.LoggerMiddleware(&logger))
	r.Use(compressor.RequestDecompressionMiddleware)
	r.Use(compressor.ResponseCompressionMiddleware(compressor.BestCompression))
	r.Use(middleware.AllowContentType("text/plain", "application/json", "text/html", "application/x-gzip"))
	r.Use(middleware.Recoverer)

	// Маршруты
	r.Post("/", linkHandler.ShortenURL)
	r.Post("/api/shorten", linkHandler.HandleGenerateShortURLJson)
	r.Post("/api/shorten/batch", linkHandler.HandleGenerateShortURLJsonBatch)
	r.Get("/{shortURL}", linkHandler.Redirect)
	r.Get("/ping", linkHandler.PingHandler)

	return r
}
