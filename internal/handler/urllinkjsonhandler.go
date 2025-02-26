package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

const (
	timeToDeleteTimeout = 5 * time.Second
)

type (
	requestBody struct {
		URL string `json:"url"`
	}

	responseBody struct {
		Result string `json:"result"`
	}

	batchRequestItem struct {
		ID  string `json:"correlation_id"`
		URL string `json:"original_url"`
	}

	batchResponseItem struct {
		ID     string `json:"correlation_id"`
		Result string `json:"short_url"`
	}

	batchResponseListPerUser struct {
		ShortURL string `json:"short_url"`
		LongURL  string `json:"original_url"`
	}
)

func (h *URLLinkHandler) HandleGenerateShortURLJson(w http.ResponseWriter, r *http.Request) {

	if !h.isContentTypeJSON(r) {
		http.Error(w, "Content-Type должен быть application/json", http.StatusBadRequest)
		return
	}

	var reqBody requestBody
	if err := h.decodeJSONBody(r, &reqBody); err != nil || reqBody.URL == "" {
		h.log.Error().Err(err).Msg("Ошибка декодирования JSON")
		http.Error(w, "Некорректное тело запроса. url должно быть json", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()

	parsedURL, err := url.ParseRequestURI(reqBody.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Некорректный URL", http.StatusBadRequest)
		return
	}

	urlModel, err := h.service.CreateShortURL(ctx, domain.URLLink{LongURL: reqBody.URL})
	if err != nil {
		if errors.Is(err, repoerrors.ErrorShortLinkAlreadyInDB) {
			h.sendJSONResponse(w, http.StatusConflict, urlModel.ShortURL)
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.sendJSONResponse(w, http.StatusCreated, urlModel.ShortURL)
}

func (h *URLLinkHandler) HandleGenerateShortURLJsonBatch(w http.ResponseWriter, r *http.Request) {
	if !h.isContentTypeJSON(r) {
		http.Error(w, "Content-Type должен быть application/json", http.StatusBadRequest)
		return
	}

	var reqBody []batchRequestItem
	if err := h.decodeJSONBody(r, &reqBody); err != nil || len(reqBody) == 0 {
		http.Error(w, "Некорректное тело запроса. url должно быть json", http.StatusBadRequest)
		return
	}

	respBody := make([]batchResponseItem, len(reqBody))
	for i, req := range reqBody {
		urlModel, err := h.service.CreateShortURL(r.Context(), domain.URLLink{LongURL: req.URL})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		respBody[i] = batchResponseItem{
			ID:     req.ID,
			Result: fmt.Sprintf("%s/%s", h.baseURL, urlModel.ShortURL),
		}
	}

	h.sendBatchJSONResponse(w, http.StatusCreated, respBody)
}

func (h *URLLinkHandler) HandleGetAllShortedURLsForUserJSON(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()

	urls, err := h.service.FindAll(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
			return
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	urlsPerUser := make([]batchResponseListPerUser, len(urls))
	for i, url := range urls {
		urlsPerUser[i] = batchResponseListPerUser{
			ShortURL: fmt.Sprintf("%s/%s", h.baseURL, url.ShortURL),
			LongURL:  url.LongURL,
		}
	}

	if len(urlsPerUser) > 0 {
		h.sendBatchJSONResponseForUser(w, http.StatusOK, urlsPerUser)
		return
	}
	h.sendBatchJSONResponseForUser(w, http.StatusNoContent, urlsPerUser)

}

func (h *URLLinkHandler) HandleDeleteShortedURLsForUserJSON(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(domain.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "UserID is missing or invalid", http.StatusUnauthorized)
		return
	}

	h.log.Info().Str("userID", userID).Msg("Пользователь запросивший удаление")
	var shortLinks []string
	if err := h.decodeArrayOfShortLinks(r, &shortLinks); err != nil || len(shortLinks) == 0 {
		http.Error(w, "Некорректное тело запроса. Это должен быть список котортких ссылок в json", http.StatusBadRequest)
		return
	}

	urlstodelete := make([]domain.URLLink, len(shortLinks))
	for i, ullink := range shortLinks {
		urlstodelete[i] = domain.URLLink{
			ShortURL: ullink,
			UserID:   userID,
		}
	}

	h.log.Info().Int("длина очереди на удаление", h.deleter.Size()).Send()

	h.deleter.Enqueue(urlstodelete...)
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusAccepted)

}

// Вспомогательные методы

func (h *URLLinkHandler) isContentTypeJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

func (h *URLLinkHandler) decodeJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *URLLinkHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, shortURL string) {
	respBody := responseBody{
		Result: strings.Join([]string{h.baseURL, shortURL}, "/"),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(respBody)
}

func (h *URLLinkHandler) sendBatchJSONResponse(w http.ResponseWriter, statusCode int, respBody []batchResponseItem) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(respBody)
}

func (h *URLLinkHandler) sendBatchJSONResponseForUser(w http.ResponseWriter, statusCode int, respBody []batchResponseListPerUser) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(respBody) > 0 {
		json.NewEncoder(w).Encode(respBody)
	}
}

func (h *URLLinkHandler) decodeArrayOfShortLinks(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
