package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
	"github.com/rs/zerolog"
)

const (
	RequestResponseTimeout = 5 * time.Second
	PingTimeout            = 3 * time.Second
	MaxQueueCapacity       = 100
	processingInterval     = 5 * time.Second
	batchSize              = 10
)

type URLLinkHandler struct {
	service     domain.URLLinkService
	baseURL     string
	log         zerolog.Logger
	deleteQueue chan domain.DeleteRecordTask
	mu          sync.Mutex
}

func NewURLLinkHandler(service domain.URLLinkService, baseURL string, logger zerolog.Logger, ctx context.Context, wg *sync.WaitGroup) *URLLinkHandler {
	h := &URLLinkHandler{
		service:     service,
		baseURL:     baseURL,
		log:         logger,
		deleteQueue: make(chan domain.DeleteRecordTask, MaxQueueCapacity),
	}

	ticker := time.NewTicker(processingInterval)

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer ticker.Stop() // Останавливаем таймер при завершении горутины

		for {
			select {
			case req, ok := <-h.deleteQueue:
				if !ok {
					// Канал закрыт, завершаем горутину
					h.log.Info().Msg("Канал удаления закрыт, завершаем горутину")
					return
				}

				if err := h.service.MarkURLsAsDeleted(ctx, req); err != nil {
					h.log.Error().Err(err).Msg("Ошибка при пометке ссылок на удаление")
				} else {
					h.log.Info().Int("количество удаленных ссылок", len(req.ShortURLs))
				}

			case <-ctx.Done():
				// Получен сигнал завершения через контекст
				h.log.Info().Msg("Получен сигнал завершения через контекст")
				return
			}

		}
	}()
	h.log.Info().Msg("Инициализация хэндлеров прошла успешно")
	return h
}

func (h *URLLinkHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDKey).(string)
	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()

	longURLBytes, err := io.ReadAll(r.Body)
	if err != nil || len(longURLBytes) == 0 {
		h.log.Info().Msg("Пустое тело запроса или ошибка чтения тела запроса")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL := string(longURLBytes)
	urllink, err := h.service.CreateShortURL(ctx, domain.URLLink{LongURL: longURL, UserID: userID})
	if err != nil {
		h.log.Info().Msg(err.Error())
		if errors.Is(err, repoerrors.ErrorShortLinkAlreadyInDB) {
			fullURL := strings.Join([]string{h.baseURL, urllink.ShortURL}, "/")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fullURL))
			return
		} else if errors.Is(err, repoerrors.ErrorSQLInternal) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			// Остальные ошибки трактуем как BadRequest
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	} else {
		fullURL := strings.Join([]string{h.baseURL, urllink.ShortURL}, "/")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fullURL))
	}
}

func (h *URLLinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	path := r.URL.Path
	shortURL := strings.TrimPrefix(path, "/")
	urllink, err := h.service.GetOriginalURL(ctx, domain.URLLink{ShortURL: shortURL})

	if err != nil {

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if urllink.DeletedFlag {
		http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
		return
	}

	w.Header().Set("Location", urllink.LongURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLLinkHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), PingTimeout)
	defer cancel()

	err := h.service.Ping(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *URLLinkHandler) Close() {
	if h.deleteQueue != nil {
		close(h.deleteQueue)
	}
}
