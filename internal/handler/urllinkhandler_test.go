package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/mocks"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

func TestShortenURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	h := NewURLLinkHandler(mockService, "http://localhost", logger)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey, "test-user"))

	expectedURLLink := domain.URLLink{
		LongURL:  "https://example.com",
		ShortURL: "abc123",
		UserID:   "test-user",
	}

	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{
			LongURL: "https://example.com",
			UserID:  "test-user",
		}).
		Return(expectedURLLink, nil)

	h.ShortenURL(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost/abc123", w.Body.String())
}

func TestShortenURL_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	h := NewURLLinkHandler(mockService, "http://localhost", logger)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/shorten", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey, "test-user"))

	existingURLLink := domain.URLLink{
		LongURL:  "https://example.com",
		ShortURL: "abc123",
		UserID:   "test-user",
	}

	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{
			LongURL: "https://example.com",
			UserID:  "test-user",
		}).
		Return(existingURLLink, repoerrors.ErrorShortLinkAlreadyInDB)

	h.ShortenURL(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, "http://localhost/abc123", w.Body.String())
}

func TestShortenURL_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	h := NewURLLinkHandler(mockService, "http://localhost", logger)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "Empty body",
			body: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body := bytes.NewBufferString(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/", body)
			r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey, "test-user"))

			h.ShortenURL(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestShortenURL_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	h := NewURLLinkHandler(mockService, "http://localhost", logger)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/shorten", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey, "test-user"))

	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{
			LongURL: "https://example.com",
			UserID:  "test-user",
		}).
		Return(domain.URLLink{}, repoerrors.ErrorSQLInternal)

	h.ShortenURL(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// func TestURLLinkHandler_ShortenURL(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		body           string
// 		mockSetup      func(*mocks.MockURLLinkService)
// 		expectedStatus int
// 		expectedBody   string
// 	}{
// 		{
// 			name: "Success - URL shortened",
// 			body: "https://example.com",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().CreateShortURL(gomock.Any(), domain.URLLink{LongURL: "https://example.com"}).
// 					Return(domain.URLLink{ShortURL: "abc123"}, nil)
// 			},
// 			expectedStatus: http.StatusCreated,
// 			expectedBody:   "http://localhost/abc123",
// 		},
// 		{
// 			name: "Conflict - URL already exists",
// 			body: "https://example.com",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().CreateShortURL(gomock.Any(), "https://example.com").
// 					Return(domain.URLLink{ShortURL: "abc123"}, repoerrors.ErrorShortLinkAlreadyInDB)
// 			},
// 			expectedStatus: http.StatusConflict,
// 			expectedBody:   "http://localhost/abc123",
// 		},
// 		{
// 			name:           "Bad Request - Empty body",
// 			body:           "",
// 			mockSetup:      func(m *mocks.MockURLLinkService) {},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   http.StatusText(http.StatusBadRequest) + "\n",
// 		},
// 	}

// 	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockURLLinkService(ctrl)
// 			tt.mockSetup(mockService)

// 			handler := NewURLLinkHandler(mockService, "http://localhost", logger)

// 			req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(tt.body))
// 			w := httptest.NewRecorder()

// 			handler.ShortenURL(w, req)

// 			resp := w.Result()
// 			defer resp.Body.Close()

// 			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

// 			body, _ := io.ReadAll(resp.Body)
// 			assert.Equal(t, tt.expectedBody, string(body))
// 		})
// 	}
// }

// func TestURLLinkHandler_Redirect(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		shortURL       string
// 		mockSetup      func(*mocks.MockURLLinkService)
// 		expectedStatus int
// 		expectedHeader string
// 	}{
// 		{
// 			name:     "Success - Redirect to original URL",
// 			shortURL: "abc123",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().GetOriginalURL(gomock.Any(), "abc123").
// 					Return("https://example.com", nil)
// 			},
// 			expectedStatus: http.StatusTemporaryRedirect,
// 			expectedHeader: "https://example.com",
// 		},
// 		{
// 			name:     "Not Found - Short URL not found",
// 			shortURL: "invalid",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().GetOriginalURL(gomock.Any(), "invalid").
// 					Return("", errors.New("not found"))
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedHeader: "",
// 		},
// 	}

// 	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockURLLinkService(ctrl)
// 			tt.mockSetup(mockService)

// 			handler := NewURLLinkHandler(mockService, "http://localhost", logger)

// 			req := httptest.NewRequest(http.MethodGet, "/"+tt.shortURL, nil)
// 			w := httptest.NewRecorder()

// 			handler.Redirect(w, req)

// 			resp := w.Result()
// 			defer resp.Body.Close()
// 			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
// 			assert.Equal(t, tt.expectedHeader, resp.Header.Get("Location"))
// 		})
// 	}
// }

// func TestURLLinkHandler_PingHandler(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		mockSetup      func(*mocks.MockURLLinkService)
// 		expectedStatus int
// 	}{
// 		{
// 			name: "Success - Ping OK",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().Ping(gomock.Any()).Return(nil)
// 			},
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name: "Internal Server Error - Ping failed",
// 			mockSetup: func(m *mocks.MockURLLinkService) {
// 				m.EXPECT().Ping(gomock.Any()).Return(errors.New("ping failed"))
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			mockService := mocks.NewMockURLLinkService(ctrl)
// 			tt.mockSetup(mockService)

// 			handler := NewURLLinkHandler(mockService, "http://localhost", logger)

// 			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
// 			w := httptest.NewRecorder()

// 			handler.PingHandler(w, req)

// 			resp := w.Result()
// 			defer resp.Body.Close()

// 			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
// 		})
// 	}
// }
