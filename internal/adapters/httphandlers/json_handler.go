package httphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	RequestBody struct {
		URL string `json:"url"`
	}

	ResponseBody struct {
		Result string `json:"result"`
	}
)

func (h *URLHandler) HandleGenerateShortURLJson(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("Content-Type") != "application/json" { // не JSON
	// 	http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
	// 	return
	// }
	// Парсим тело запроса
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.URL == "" {
		http.Error(w, "Некорректное тело запроса. url должно быть строкой", http.StatusBadRequest)
		return
	}

	urlModel, err := h.urlService.GenerateShortURL(reqBody.URL)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	// Формируем ответ
	respBody := ResponseBody{
		Result: fmt.Sprintf("%s/%s", h.baseURL, urlModel.Short),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respBody)
}
