package httphandlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/adapters/httphandlers"
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
	"github.com/physicist2018/url-shortener-go/internal/core/services/url/mocks"
	"github.com/stretchr/testify/assert"
)

func TestURLHandler_HandleGenerateShortURLJson(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		body               string
		contentType        string
		expectedStatusCode int
		expectedResponse   string
		mockGenerateError  error
		mockShortURL       urlmodels.URL
	}{
		{
			name:               "Successful URL shortening",
			method:             "POST",
			url:                "/api/shorten",
			body:               `{"url":"https://example.com"}`,
			contentType:        "application/json",
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"result":"http://localhost:8080/abc123"}` + "\n",
			mockShortURL: urlmodels.URL{
				Short:    "abc123",
				Original: "https://example.com",
			},
		},
		{
			name:               "Wrong content type",
			method:             "POST",
			url:                "/api/shorten",
			body:               `{"url":"https://example.com"}`,
			contentType:        "text/plain",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Content-Type must be application/json\n",
			mockShortURL: urlmodels.URL{
				Short:    "abc123",
				Original: "https://example.com",
			},
		},
		{
			name:               "Invalid json",
			method:             "POST",
			url:                "/api/shorten",
			body:               `{url:"https://example.com"}`,
			contentType:        "application/json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Некорректное тело запроса. url должно быть строкой\n",
			mockShortURL: urlmodels.URL{
				Short:    "abc123",
				Original: "https://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockURLService{}
			handler := httphandlers.NewURLHandler(mockService, "http://localhost:8080")
			// Настроим мок для генерации короткой ссылки
			mockService.On("GenerateShortURL", tt.mockShortURL.Original).Return(tt.mockShortURL, tt.mockGenerateError)
			// Создаем запрос
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			// Обрабатываем запрос с помощью обработчика
			handler.HandleGenerateShortURLJson(w, req)

			// Проверяем код статуса ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем содержимое ответа
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
