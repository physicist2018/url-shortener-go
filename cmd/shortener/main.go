package main

import (
	// "github.com/physicist2018/url-shortener-go/internal/config"
	// "github.com/physicist2018/url-shortener-go/internal/shortener"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/physicist2018/url-shortener-go/adapters/httphandlers"
	"github.com/physicist2018/url-shortener-go/adapters/memory"
	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/core/services/url"
)

// func main() {
// 	_ = config.ConfigApp()
// 	//fmt.Println(config.DefaultConfig)
// 	if err := shortener.RunServer(); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	configuration := config.MakeConfig()
	configuration.Parse()

	urlRepo := memory.NewURLRepositoryMap()
	urlService := url.NewURLService(urlRepo)
	urlHandler := httphandlers.NewURLHandler(urlService)

	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("text/plain"))
	r.Post("/", urlHandler.HandleGenerateShortURL)
	r.Get("/{shortUrl}", urlHandler.HandleRedirect)

	log.Fatal(http.ListenAndServe(configuration.ServerAddr, r))
}
