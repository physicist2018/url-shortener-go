package urlshortenerserver

import (
	"compress/gzip"
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/physicist2018/url-shortener-go/internal/adapters/httphandlers"
	"github.com/physicist2018/url-shortener-go/internal/adapters/memory"
	"github.com/physicist2018/url-shortener-go/internal/config"
	"github.com/physicist2018/url-shortener-go/internal/core/services/url"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/compressor"
	"github.com/physicist2018/url-shortener-go/internal/middlewares/httplogger"
	"github.com/physicist2018/url-shortener-go/pkg/utils"
	"github.com/rs/zerolog"
	//"go.uber.org/zap"
)

type URLShortenerServer struct {
	Config  *config.Config
	Logger  *zerolog.Logger //*zap.SugaredLogger
	Handler http.Handler
	HTTP    *http.Server
	URLRepo *memory.URLRepositoryMap
}

// NewServer initializes the server with configuration and logger
func NewURLShortenerServer(config *config.Config, logger *zerolog.Logger) *URLShortenerServer {
	// генератор случайных строк
	randomStringGenerator := utils.NewRandomString(config.MaxShortURLLength,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	)

	// хранилище ссылок в памяти
	urlRepo := memory.NewURLRepositoryMap(config.FileStoragePath)

	// Восстановление данных из файла
	logger.Info().Msg(strings.Join([]string{"Загрузка ссылок из файла", config.FileStoragePath}, " "))

	if err := urlRepo.RestoreFromFile(); err != nil {
		if !errors.Is(err, memory.ErrorOpeningFileWhenRestore) {
			logger.Panic().Err(err)
		}
		logger.Info().Msg(strings.Join([]string{"При восстановлении хранилища из файла оный не обнаружен (будет создан при закрытии):", err.Error()}, " "))
		//logger.Infof("При восстановлении хранилища из файла оный не обнаружен (будет создан при закрытии): %s", err.Error())
	}

	urlService := url.NewURLService(urlRepo, randomStringGenerator)
	urlHandler := httphandlers.NewURLHandler(urlService, config.BaseURLServer)

	// Настраиваем маршруты и middlewares
	r := chi.NewRouter()
	r.Use(compressor.RequestDecompressionMiddleware)
	r.Use(compressor.ResponseCompressionMiddleware(gzip.BestCompression))
	r.Use(middleware.AllowContentType("text/plain", "application/json", "text/html", "application/x-gzip"))
	r.Use(httplogger.LoggerMiddleware(logger))
	r.Use(middleware.Recoverer)

	// Определение эндпоинтов
	r.Post("/", urlHandler.HandleGenerateShortURL)
	r.Post("/api/shorten", urlHandler.HandleGenerateShortURLJson)
	r.Get("/{shortURL}", urlHandler.HandleRedirect)

	// Создаем и возвращаем сервер
	return &URLShortenerServer{
		Config:  config,
		Logger:  logger,
		Handler: r,
		URLRepo: urlRepo,
		HTTP: &http.Server{
			Addr:    config.ServerAddr,
			Handler: r,
		},
	}
}

// Start launches the server and listens for incoming connections
func (s *URLShortenerServer) Start() {
	// Запуск HTTP-сервера в горутине
	go func() {
		if err := s.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Err(err)
			//s.Logger.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()
	s.Logger.Info().Msg(strings.Join([]string{"Сервер запущен на ", s.Config.ServerAddr}, " "))
}

// Shutdown gracefully stops the server
func (s *URLShortenerServer) Shutdown() {
	// Создаем контекст с таймаутом для завершения сервера
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(),
		time.Duration(s.Config.MaxShutdownTime)*time.Second)
	defer cancelShutdown()

	// Закрытие сервера
	s.Logger.Info().Msg("Выключение сервера...")
	if err := s.HTTP.Shutdown(shutdownCtx); err != nil {
		s.Logger.Err(err)
	}

	// Сохранение ссылок в файл
	s.Logger.Info().Msg(strings.Join([]string{"Сохраняем базу ссылок на диск в файл", s.Config.FileStoragePath}, " "))
	if err := s.URLRepo.DumpToFile(); err != nil {
		s.Logger.Err(err)
	}

	s.Logger.Info().Msg("Сервер корректно завершил работу.")
}
