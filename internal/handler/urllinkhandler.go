package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
	"github.com/rs/zerolog"
)

const (
	RequestResponseTimeout = 5 * time.Second
)

type URLLinkHandler struct {
	service domain.URLLinkService
	baseURL string
	log     zerolog.Logger
}

func NewURLLinkHandler(service domain.URLLinkService, baseURL string, logger zerolog.Logger) *URLLinkHandler {
	return &URLLinkHandler{
		service: service,
		baseURL: baseURL,
		log:     logger,
	}
}

func (h *URLLinkHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
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
	userID := r.Context().Value("userID").(string)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	//shortURL := chi.URLParam(r, "shortURL")
	path := r.URL.Path
	shortURL := strings.TrimPrefix(path, "/")
	urllink, err := h.service.GetOriginalURL(ctx, domain.URLLink{ShortURL: shortURL, UserID: userID})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", urllink.LongURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *URLLinkHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Ping(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
