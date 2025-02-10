package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
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
