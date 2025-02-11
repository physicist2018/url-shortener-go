package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

type (
	requestBody struct {
		URL string `json:"url"`
	}

	responseBody struct {
		Result string `json:"result"`
	}
)

func (h *URLLinkHandler) HandleGenerateShortURLJson(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" { // не JSON
		http.Error(w, "Content-Type должен быть application/json", http.StatusBadRequest)
		return
	}
	// Парсим тело запроса
	var reqBody requestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.URL == "" {
		http.Error(w, "Некорректное тело запроса. url должно быть строкой", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), RequestResponseTimeout)
	defer cancel()
	urlModel, err := h.service.CreateShortURL(ctx, reqBody.URL)

	if errors.Is(err, repoerrors.ErrUrlAlreadyInDB) {
		respBody := responseBody{
			Result: strings.Join([]string{h.baseURL, urlModel.ShortURL}, "/"),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(respBody)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Формируем ответ
	respBody := responseBody{
		Result: strings.Join([]string{h.baseURL, urlModel.ShortURL}, "/"),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respBody)
}

func (h *URLLinkHandler) HandleGenerateShortURLJsonBatch(w http.ResponseWriter, r *http.Request) {
	type (
		requestBody struct {
			ID  string `json:"correlation_id"`
			URL string `json:"original_url"`
		}

		responseBody struct {
			ID     string `json:"correlation_id"`
			Result string `json:"short_url"`
		}
	)

	if r.Header.Get("Content-Type") != "application/json" { // не JSON
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	// Парсим тело запроса
	var reqBody []requestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || len(reqBody) == 0 {
		http.Error(w, "Некорректное тело запроса или пустой запрос", http.StatusBadRequest)
		return
	}

	respBody := make([]responseBody, len(reqBody))

	for i, req := range reqBody {
		urlModel, err := h.service.CreateShortURL(r.Context(), req.URL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
			return
		}
		respBody[i].ID = req.ID
		respBody[i].Result = fmt.Sprintf("%s/%s", h.baseURL, urlModel.ShortURL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respBody)
}
