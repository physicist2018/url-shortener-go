package httphandlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/adapters/httphandlers"
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
	"github.com/physicist2018/url-shortener-go/internal/core/services/url/mocks"
	"github.com/stretchr/testify/assert"
)

func TestURLHandler_HandleGenerateShortURL(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		body               string
		expectedStatusCode int
		expectedResponse   string
		mockGenerateError  error
		mockShortURL       *urlmodels.URL
	}{

		{
			name:               "Valid URL but failed to generate short URL",
			method:             "POST",
			url:                "/",
			body:               "https://example.com",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   http.StatusText(http.StatusBadRequest) + "\n",
			mockGenerateError:  errors.New("failed to generate short URL"),
			mockShortURL: &urlmodels.URL{
				Short:    "/",
				Original: "https://example.com",
			},
		},
		{
			name:               "Successful URL shortening",
			method:             "POST",
			url:                "/",
			body:               "https://example.com",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "http://localhost:8080/abc123",
			mockShortURL: &urlmodels.URL{
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
			mockService.On("GenerateShortURL", tt.body).Return(tt.mockShortURL, tt.mockGenerateError)
			// Создаем запрос
			req := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			// Обрабатываем запрос с помощью обработчика
			handler.HandleGenerateShortURL(w, req)

			// Проверяем код статуса ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем содержимое ответа
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}

func TestURLHandler_HandleRedirect(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		shortURL           string
		expectedStatusCode int
		expectedLocation   string
		mockGetURLResult   *urlmodels.URL
		mockGetURLError    error
	}{
		{
			name:               "Valid short URL",
			method:             "GET",
			shortURL:           "http://localhost:8080/abc123",
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLocation:   "https://example.com",
			mockGetURLResult: &urlmodels.URL{
				Short:    "abc123",
				Original: "",
			},
			mockGetURLError: nil,
		},
		{
			name:               "Invalid short URL",
			method:             "GET",
			shortURL:           "http://localhost:8080/invalid123",
			expectedStatusCode: http.StatusNotFound, // 404 Not Found
			expectedLocation:   "",
			mockGetURLResult: &urlmodels.URL{
				Short:    "invalid123",
				Original: "",
			},
			mockGetURLError: errors.New("Not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockURLService{}
			handler := httphandlers.NewURLHandler(mockService, "http://localhost:8080")
			// Настроим мок для генерации короткой ссылки
			mockService.On("GetOriginalURL", tt.mockGetURLResult.Short).Return(tt.mockGetURLResult, tt.mockGetURLError)
			// Создаем запрос

			req := httptest.NewRequest(tt.method, tt.shortURL, nil)
			w := httptest.NewRecorder()

			// Обрабатываем запрос с помощью обработчика
			handler.HandleRedirect(w, req)

			// Проверяем код статуса ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем код статуса ответа
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Проверяем, что в случае перенаправления был правильный Location
			if tt.expectedStatusCode == http.StatusFound {
				location := w.Header().Get("Location")
				assert.Equal(t, tt.expectedLocation, location)
			}

			// Если URL не найден, проверяем, что в ответе ошибка "URL not found"
			if tt.expectedStatusCode == http.StatusNotFound {
				assert.Equal(t, http.StatusText(http.StatusNotFound)+"\n", w.Body.String())
			}
		})
	}
}
