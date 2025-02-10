package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/repository/repofactorymethod"
	"github.com/physicist2018/url-shortener-go/internal/service"
	randomstringgenerator "github.com/physicist2018/url-shortener-go/pkg/randomstring_generator"
	"github.com/stretchr/testify/assert"
)

func TestURLLinkHandler_ShortenURL(t *testing.T) {
	randomStringGenerator := randomstringgenerator.NewRandomStringFixed()

	repofactory := repofactorymethod.NewRepofactorymethod()
	linkRepo, _ := repofactory.CreateRepo("inmemory", "test.db")
	defer linkRepo.Close()
	urlService := service.NewURLLinkService(linkRepo, randomStringGenerator)
	type fields struct {
		service *service.URLLinkService
		baseURL string
	}
	type args struct {
		expectedStatusCode int
		body               string
		expectedResponse   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Invalid URL",
			fields: fields{
				service: urlService,
				baseURL: "http://localhost:8080",
			},
			args: args{
				expectedStatusCode: http.StatusBadRequest,
				body:               "",
				expectedResponse:   "",
			},
		},
		{
			name: "Successfully created URL",
			fields: fields{
				service: urlService,
				baseURL: "http://localhost:8080",
			},
			args: args{
				expectedStatusCode: http.StatusCreated,
				body:               "http://ya.ru",
				expectedResponse:   "http://localhost:8080/wSv9w",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &URLLinkHandler{
				service: tt.fields.service,
				baseURL: tt.fields.baseURL,
			}
			req := httptest.NewRequest("POST", "http://localhost:8080/", strings.NewReader(tt.args.body))
			w := httptest.NewRecorder()

			h.ShortenURL(w, req)
			respBytes, _ := io.ReadAll(w.Body)
			log.Println(string(respBytes))
			assert.Equal(t, tt.args.expectedStatusCode, w.Code)
			if tt.args.body != "" {
				assert.Equal(t, tt.args.expectedResponse, string(respBytes))
			}
		})
	}
}
