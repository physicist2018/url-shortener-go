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
	"github.com/physicist2018/url-shortener-go/pkg/utils"
)

const maxShortURLLength = 6

func main() {
	configuration := config.MakeConfig()
	configuration.Parse()

	randomStringGenerator := utils.NewRandomString(maxShortURLLength, rand.New(rand.NewSource(time.Now().UnixNano())))
	urlRepo := memory.NewURLRepositoryMap()
	urlService := url.NewURLService(urlRepo, randomStringGenerator)
	urlHandler := httphandlers.NewURLHandler(urlService, configuration.BaseURLServer)

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("text/plain"))
	r.Post("/", urlHandler.HandleGenerateShortURL)
	r.Get("/{shortURL}", urlHandler.HandleRedirect)

	if err := http.ListenAndServe(configuration.ServerAddr, r); err != nil {
		log.Fatal(err)
	}
}
