package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/deleter"
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
	service domain.URLLinkService
	baseURL string
	log     zerolog.Logger
	//deleteQueue chan domain.DeleteRecordTask
	deleter *deleter.Deleter
	//mu          sync.Mutex
}

func NewURLLinkHandler(service domain.URLLinkService, baseURL string, logger zerolog.Logger, deleter *deleter.Deleter) *URLLinkHandler {
	h := &URLLinkHandler{
		service: service,
		baseURL: baseURL,
		log:     logger,
		deleter: deleter,
	}

	h.log.Info().Msg("Инициализация хэндлеров прошла успешно")
	return h
}

func (h *URLLinkHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDKey).(string)
	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()

	longURLBytes, err := io.ReadAll(r.Body)
	if err != nil || len(longURLBytes) == 0 {
		h.log.Info().
			Msg("Пустое тело запроса или ошибка чтения тела запроса")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL := string(longURLBytes)
	longURL = strings.TrimSpace(longURL)
	_, err = url.ParseRequestURI(longURL)
	if err != nil {
		h.log.Error().Msg("Некорректный URL")
		http.Error(w, "Некорректный URL", http.StatusBadRequest)
		return
	}

	urllink, err := h.service.CreateShortURL(ctx, domain.URLLink{LongURL: longURL, UserID: userID})

	if err != nil {
		switch {
		case errors.Is(err, repoerrors.ErrorShortLinkAlreadyInDB):
			fullURL := strings.Join([]string{h.baseURL, urllink.ShortURL}, "/")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fullURL))

		case errors.Is(err, repoerrors.ErrorSQLInternal):
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))

		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))

		}
		return
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
	h.log.Info().
		Str("shortURL", shortURL).
		Str("longURL", urllink.LongURL).
		Msg("Перенаправление выполнено успешно")
}

func (h *URLLinkHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), PingTimeout)
	defer cancel()

	err := h.service.Ping(ctx)
	if err != nil {
		h.log.Info().
			Err(err).
			Msg("При проверке соединения возникла ошибка")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.log.Info().
		Msg("Соединение с БД успешно проверено")
}

func (h *URLLinkHandler) Close() {
	if h.deleter != nil {
		h.deleter.Close()
	}
}
