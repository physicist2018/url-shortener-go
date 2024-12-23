package shortener

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/urlstorage"
)

func Test_postRoute(t *testing.T) {
	type args struct {
		storage     *urlstorage.URLStorage
		contenttype string
		content     string
		want        string
		code        int
		url         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "yandex.ru usual response", args: args{
				storage:     urlstorage.GetDefaultURLStorage(),
				contenttype: "text/plain",
				content:     "yandex.ru",
				want:        "http://example.com/wSv9wq",
				code:        201,
				url:         "/",
			},
		},

		{
			name: "yandex.ru bad content type", args: args{
				storage:     urlstorage.GetDefaultURLStorage(),
				contenttype: "text/html",
				content:     "yandex.ru",
				want:        "http://example.com/wSv9wq",
				code:        400,
				url:         "/",
			},
		},

		{
			name: "yandex.ru bad url", args: args{
				storage:     urlstorage.GetDefaultURLStorage(),
				contenttype: "text/plain",
				content:     "yandex.ru",
				want:        "http://example.com/wSv9wq",
				code:        400,
				url:         "/bad",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlstorage.SetDefaultURLStorage(tt.args.storage)
			request := httptest.NewRequest("POST", tt.args.url, strings.NewReader(tt.args.content))
			request.Header.Add("Content-Type", tt.args.contenttype)
			response := httptest.NewRecorder()
			postRoute(response, request)

			if tt.args.code != http.StatusBadRequest {
				if response.Code != tt.args.code {
					t.Errorf("postRoute() = %v, want %v", response.Code, tt.args.code)
				}

				if response.Body.String() != tt.args.want {
					t.Errorf("postRoute() = %v, want %v", response.Body.String(), tt.args.want)
				}
			} else {
				if response.Code != tt.args.code {
					t.Errorf("postRoute() = %v, want %v", response.Code, tt.args.code)
				}
			}

		})
	}
}

func Test_getRoute(t *testing.T) {
	type args struct {
		storage *urlstorage.URLStorage

		want string
		code int
		url  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "yandex.ru usual response", args: args{
				storage: &urlstorage.URLStorage{
					Store: []urlstorage.URLItem{
						{
							LongURL:  "yandex.ru",
							ShortURL: "wSv9wq",
						},
					},
				},
				want: "yandex.ru",
				code: http.StatusTemporaryRedirect,
				url:  "http://localhost:8080/wSv9wq",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlstorage.SetDefaultURLStorage(tt.args.storage)
			request := httptest.NewRequest("GET", tt.args.url, nil)

			response := httptest.NewRecorder()
			getRoute(response, request)

			if response.Code != tt.args.code {
				t.Errorf("getRoute() = %v, want %v", response.Code, tt.args.code)
			}

			if response.Header().Get("Location") != tt.args.want {
				t.Errorf("getRoute() Location is = %v, want: %v", response.Header(), tt.args.want)
			}

		})
	}
}
