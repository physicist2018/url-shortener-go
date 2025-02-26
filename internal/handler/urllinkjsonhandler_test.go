package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/physicist2018/url-shortener-go/internal/deleter"
	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/mocks"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestHandleGenerateShortURLJson_Success(t *testing.T) {
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

	// Данные запроса
	requestBody := requestBody{
		URL: "https://example.com",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBodyBytes))
	r.Header.Set("Content-Type", "application/json")

	// Ожидаемая модель URLLink
	expectedURLLink := domain.URLLink{
		LongURL:  "https://example.com",
		ShortURL: "abc123",
	}

	// Ожидания для вызова mock
	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{LongURL: "https://example.com"}).
		Return(expectedURLLink, nil)

	w := httptest.NewRecorder()
	h.HandleGenerateShortURLJson(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respBody responseBody
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "http://localhost/abc123", respBody.Result)
	h.Close()
	wg.Wait()
}

func TestHandleGenerateShortURLJson_Conflict(t *testing.T) {
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

	// Данные запроса
	requestBody := requestBody{
		URL: "https://example.com",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBodyBytes))
	r.Header.Set("Content-Type", "application/json")

	// Ожидаемая модель URLLink (уже существует)
	existingURLLink := domain.URLLink{
		LongURL:  "https://example.com",
		ShortURL: "abc123",
	}

	// Ожидания для вызова mock
	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{LongURL: "https://example.com"}).
		Return(existingURLLink, repoerrors.ErrorShortLinkAlreadyInDB)

	w := httptest.NewRecorder()
	h.HandleGenerateShortURLJson(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var respBody responseBody
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "http://localhost/abc123", respBody.Result)

	h.Close()
	wg.Wait()
}

func TestHandleGenerateShortURLJsonBatch_Success(t *testing.T) {
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

	// Данные запроса
	requestItems := []batchRequestItem{
		{ID: "1", URL: "https://example.com"},
		{ID: "2", URL: "https://test.com"},
	}
	reqBodyBytes, _ := json.Marshal(requestItems)
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(reqBodyBytes))
	r.Header.Set("Content-Type", "application/json")

	// Ожидаемая модель URLLink
	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{LongURL: "https://example.com"}).
		Return(domain.URLLink{ShortURL: "abc123"}, nil)
	mockService.
		EXPECT().
		CreateShortURL(gomock.Any(), domain.URLLink{LongURL: "https://test.com"}).
		Return(domain.URLLink{ShortURL: "xyz789"}, nil)

	w := httptest.NewRecorder()
	h.HandleGenerateShortURLJsonBatch(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respBody []batchResponseItem
	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, 2, len(respBody))
	assert.Equal(t, "1", respBody[0].ID)
	assert.Equal(t, "http://localhost/abc123", respBody[0].Result)
	assert.Equal(t, "2", respBody[1].ID)
	assert.Equal(t, "http://localhost/xyz789", respBody[1].Result)

	h.Close()
	wg.Wait()
}

func TestHandleGenerateShortURLJson_BadRequest(t *testing.T) {
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

	// Некорректное тело запроса
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte("{invalid-json")))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.HandleGenerateShortURLJson(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	h.Close()
	wg.Wait()
}
