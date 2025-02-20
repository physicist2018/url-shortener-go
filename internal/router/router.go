package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/physicist2018/url-shortener-go/internal/handler"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/authenticator"
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
	r.Post("/", authenticator.AuthMiddlewareFunc(linkHandler.ShortenURL))
	r.Post("/api/shorten", authenticator.AuthMiddlewareFunc(linkHandler.HandleGenerateShortURLJson))
	r.Post("/api/shorten/batch", authenticator.AuthMiddlewareFunc(linkHandler.HandleGenerateShortURLJsonBatch))
	r.Get("/{shortURL}", linkHandler.Redirect)
	r.Get("/ping", linkHandler.PingHandler)
	r.Get("/api/user/urls", authenticator.AuthMiddlewareFunc(linkHandler.HandleGetAllShortedURLsForUserJson))
	return r
}
