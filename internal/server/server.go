package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

func NewServer(addr string, handler http.Handler, logger zerolog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}

func (s *Server) Start() {
	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info().Msg("Запуск сервера")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("Ошибка запуска сервера")
		}
	}()

	<-done
	s.logger.Info().Msg("Остановка сервера")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Ошибка при остановке сервера")
	}
	s.logger.Info().Msg("Сервер остановлен")
}
