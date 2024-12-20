package shortener

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/urlstorage"
	"github.com/stretchr/testify/assert"
)

func TestPostRoute(t *testing.T) {

	//urls := urlstorage.GetDefaultUrlStorage()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString("yandex.ru\r\n"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	mainHandler(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

}

func TestMainRoute(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		urls        urlstorage.URLStorage
		name        string
		method      string
		url         string
		want        want
		contenttype string
	}{
		{
			urls: urlstorage.URLStorage{
				Store: map[string]string{},
			},
			name:        "test_get_nothing",
			method:      "GET",
			url:         "http://localhost:8080/",
			contenttype: "text/plain",
			want: want{
				code:     400,
				response: "Bad Request\n",
			},
		},
		{
			urls: urlstorage.URLStorage{
				Store: map[string]string{
					"qwerty": "yandex.ru",
				},
			},
			name:        "test_get_yandex",
			method:      "GET",
			url:         "http://localhost:8080/qwerty",
			contenttype: "text/plain",
			want: want{
				code:     307,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *http.Request
			if tt.method == "POST" {
				r = httptest.NewRequest(tt.method, tt.url, strings.NewReader("yandex.ru"))
			} else {
				r = httptest.NewRequest(tt.method, tt.url, nil)
			}
			r.Header.Set("Content-Type", tt.contenttype)
			w := httptest.NewRecorder()
			urlstorage.SetDefaultUrlStorage(&tt.urls)
			mainHandler(w, r)
			assert.Equal(t, tt.want.response, w.Body.String(), "Ответ не совпадает")
			assert.Equal(t, tt.want.code, w.Code, "Код запроса не совпадает")
		})
	}
}
