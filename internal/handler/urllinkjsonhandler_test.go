package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
	"github.com/physicist2018/url-shortener-go/internal/service"
)

func TestURLLinkHandler_HandleGenerateShortURLJson(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		mockSetup      func(*service.MockURLLinkServicer)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success - URL shortened",
			contentType: "application/json",
			body:        `{"url": "https://example.com"}`,
			mockSetup: func(m *service.MockURLLinkServicer) {
				m.EXPECT().CreateShortURL(gomock.Any(), "https://example.com").
					Return(&domain.URLLink{ShortURL: "abc123"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"result":"http://localhost/abc123"}` + "\n",
		},
		{
			name:        "Conflict - URL already exists",
			contentType: "application/json",
			body:        `{"url": "https://example.com"}`,
			mockSetup: func(m *service.MockURLLinkServicer) {
				m.EXPECT().CreateShortURL(gomock.Any(), "https://example.com").
					Return(&domain.URLLink{ShortURL: "abc123"}, repoerrors.ErrorShortLinkAlreadyInDB)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"result":"http://localhost/abc123"}` + "\n",
		},
		{
			name:           "Bad Request - Invalid Content-Type",
			contentType:    "text/plain",
			body:           `{"url": "https://example.com"}`,
			mockSetup:      func(m *service.MockURLLinkServicer) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `Content-Type должен быть application/json` + "\n",
		},
		{
			name:           "Bad Request - Invalid JSON body",
			contentType:    "application/json",
			body:           `{"url": ""}`,
			mockSetup:      func(m *service.MockURLLinkServicer) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `Некорректное тело запроса. url должно быть json` + "\n",
		},
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockURLLinkServicer(ctrl)
			tt.mockSetup(mockService)

			handler := NewURLLinkHandler(mockService, "http://localhost", logger)

			req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handler.HandleGenerateShortURLJson(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.expectedBody, string(body))
		})
	}
}

func TestURLLinkHandler_HandleGenerateShortURLJsonBatch(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		mockSetup      func(*service.MockURLLinkServicer)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success - Batch URLs shortened",
			contentType: "application/json",
			body:        `[{"correlation_id": "1", "original_url": "https://example.com"}]`,
			mockSetup: func(m *service.MockURLLinkServicer) {
				m.EXPECT().CreateShortURL(gomock.Any(), "https://example.com").
					Return(&domain.URLLink{ShortURL: "abc123"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `[{"correlation_id":"1","short_url":"http://localhost/abc123"}]` + "\n",
		},
		{
			name:           "Bad Request - Invalid Content-Type",
			contentType:    "text/plain",
			body:           `[{"correlation_id": "1", "original_url": "https://example.com"}]`,
			mockSetup:      func(m *service.MockURLLinkServicer) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Content-Type должен быть application/json\n",
		},
		{
			name:           "Bad Request - Invalid JSON body",
			contentType:    "application/json",
			body:           `[]`,
			mockSetup:      func(m *service.MockURLLinkServicer) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Некорректное тело запроса. url должно быть json\n",
		},
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := service.NewMockURLLinkServicer(ctrl)
			tt.mockSetup(mockService)

			handler := NewURLLinkHandler(mockService, "http://localhost", logger)

			req := httptest.NewRequest(http.MethodPost, "/shorten/batch", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handler.HandleGenerateShortURLJsonBatch(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.expectedBody, string(body))
		})
	}
}
