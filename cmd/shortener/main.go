package main

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
const maxShutdownTime = 5

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

	sugar.Info("Загрузка ссылок из файла... ", configuration.FileStoragePath)
	if err = urlRepo.RestoreFromFile(configuration.FileStoragePath); err != nil {
		if !errors.Is(err, memory.ErrorOpeningFile) {
			sugar.Panic(err)
		}
		sugar.Infof("При восстановлении хранилища из файла оный не обнаружен (будет создан при закритии): %s", err.Error())
	}

	urlService := url.NewURLService(urlRepo, randomStringGenerator)
	urlHandler := httphandlers.NewURLHandler(urlService, configuration.BaseURLServer)

	r := chi.NewRouter()

	// Устанавливаем необходимые middlwares
	//r.Use(compressor.CompressionMiddleware(gzip.BestSpeed))
	r.Use(compressor.RequestDecompressionMiddleware)
	r.Use(compressor.ResponseCompressionMiddleware(gzip.BestCompression))
	// Устанавливаем допустимые типв контента
	// application/x-gzip без него тесты не проходят
	r.Use(middleware.AllowContentType("text/plain", "application/json", "text/html", "application/x-gzip"))
	r.Use(httplogger.LoggerMiddleware(sugar))
	r.Use(middleware.Recoverer)

	// Определяем наши эндпоинты
	r.Post("/", urlHandler.HandleGenerateShortURL)
	r.Post("/api/shorten", urlHandler.HandleGenerateShortURLJson)
	r.Get("/{shortURL}", urlHandler.HandleRedirect)

	// Создаем сервер
	server := http.Server{
		Addr:    configuration.ServerAddr,
		Handler: r,
	}

	// Используем канал для перехвата сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Создаем контекст для управления временем завершения
	ctx := context.WithoutCancel(context.Background())

	// Используем sync.WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup
	wg.Add(1)

	// Запуск HTTP-сервера в горутине
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Ошибка при запуске сервера:", err)
		}
	}()
	fmt.Println("\nСервер зaпущен!")
	// Ожидаем получения сигнала завершения
	<-stop
	fmt.Println("\nПолучен сигнал завершения, выключаем сервер...")

	// Создаем контекст с таймаутом для корректного завершения сервера
	shutdownCtx, cancelShutdown := context.WithTimeout(ctx, maxShutdownTime*time.Second)
	defer cancelShutdown()

	// Закрытие HTTP-сервера
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Println("Ошибка при выключении сервера:", err)
	}

	// Ожидаем завершения всех горутин (по факту, к нас в wg она одна)
	wg.Wait()

	fmt.Println("Сервер адекватно выключен.")

	sugar.Info("Сохраняем базу ссылок на диск в файл ", configuration.FileStoragePath)
	if err = urlRepo.DumpToFile(configuration.FileStoragePath); err != nil {
		sugar.Error("При записи на диск возникла ошибка: ", err)
	}
}
