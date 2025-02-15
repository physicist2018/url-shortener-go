package handler

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

const (
	RequestResponseTimeout = 5 * time.Second
)

type URLLinkHandler struct {
	service domain.URLLinkServicer
	baseURL string
}

func NewURLLinkHandler(service domain.URLLinkServicer, baseURL string) *URLLinkHandler {
	return &URLLinkHandler{
		service: service,
		baseURL: baseURL,
	}
}

func (h *URLLinkHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()

	longURLBytes, err := io.ReadAll(r.Body)
	if err != nil || len(longURLBytes) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL := string(longURLBytes)
	urllink, err := h.service.CreateShortURL(ctx, longURL)

	switch {
	case errors.Is(err, repoerrors.ErrURLAlreadyInDB):
		fullURL := strings.Join([]string{h.baseURL, urllink.ShortURL}, "/")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fullURL))
	case err != nil:
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	default:
		fullURL := strings.Join([]string{h.baseURL, urllink.ShortURL}, "/")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fullURL))
	}
}

func (h *URLLinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	//shortURL := chi.URLParam(r, "shortURL")
	path := r.URL.Path
	shortURL := strings.TrimPrefix(path, "/")

	originalURL, err := h.service.GetOriginalURL(ctx, shortURL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
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
