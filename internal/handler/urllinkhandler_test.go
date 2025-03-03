package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/physicist2018/url-shortener-go/internal/deleter"
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/mocks"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

func TestShortenURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	linkDeleter := deleter.NewDeleter(mockService, logger)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linkDeleter.Start(ctx, &wg) //Запускаем горутину асинхронного удаления ссылок
	h := NewURLLinkHandler(mockService, "http://localhost", logger, linkDeleter)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey{}, "test-user"))

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
	h.Close()
	wg.Wait()
}

func TestShortenURL_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	linkDeleter := deleter.NewDeleter(mockService, logger)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linkDeleter.Start(ctx, &wg) //Запускаем горутину асинхронного удаления ссылок

	h := NewURLLinkHandler(mockService, "http://localhost", logger, linkDeleter)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/shorten", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey{}, "test-user"))

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

	h.Close()
	wg.Wait()
}

func TestShortenURL_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	linkDeleter := deleter.NewDeleter(mockService, logger)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linkDeleter.Start(ctx, &wg) //Запускаем горутину асинхронного удаления ссылок

	h := NewURLLinkHandler(mockService, "http://localhost", logger, linkDeleter)

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
			r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey{}, "test-user"))

			h.ShortenURL(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}

	h.Close()
	wg.Wait()
}

func TestShortenURL_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockURLLinkService(ctrl)
	logger := zerolog.New(nil)
	linkDeleter := deleter.NewDeleter(mockService, logger)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linkDeleter.Start(ctx, &wg) //Запускаем горутину асинхронного удаления ссылок

	h := NewURLLinkHandler(mockService, "http://localhost", logger, linkDeleter)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString("https://example.com")
	r := httptest.NewRequest(http.MethodPost, "/shorten", body)
	r = r.WithContext(context.WithValue(r.Context(), domain.UserIDKey{}, "test-user"))

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

	h.Close()
	wg.Wait()
}
