package shortener

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/physicist2018/url-shortener-go/internal/urlstorage"
)

func TestPostRoute(t *testing.T) {

	//urls := urlstorage.GetDefaultUrlStorage()
	// Первый тест, проверяем штатное функционирование обработчика
	// смотрим результат
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString("yandex.ru\r\n"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	postRoute(rr, req)

	// respBody, _ := rr.Body.ReadString('\n')
	// // Check the status code is what we expect.
	// if status := rr.Code; status != http.StatusCreated {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusCreated)
	// }

	// if respBody != "http://localhost:8080/wSv9wq" {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		respBody, "http://localhost:8080/wSv9wq")
	// }

	req, err = http.NewRequest("POST", "/", bytes.NewBufferString("yandex.ru\r\n"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	postRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	req, err = http.NewRequest("PUT", "/", bytes.NewBufferString("yandex.ru\r\n"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	postRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	req, err = http.NewRequest("POST", "", bytes.NewBufferString("yandex.ru\r\n"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	postRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	req, err = http.NewRequest("POST", "/", bytes.NewBufferString("yandex.ru\nmail.ru"))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	postRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

}

func TestGetRoute(t *testing.T) {
	// Первый  тест, проверяем что при запросе несуществующего url возвращается 404
	req, err := http.NewRequest("GET", "/qwerty", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	getRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	urlstorage.SetDefaultURLStorage(&urlstorage.URLStorage{
		Store: []urlstorage.URLItem{
			{
				LongURL:  "yandex.ru",
				ShortURL: "qwerty",
			},
		},
	})
	// Второй тест - запрос нужного короткого URL возвращает длинный URL
	req, err = http.NewRequest("GET", "http://localhost:8080/qwerty", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	getRoute(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusTemporaryRedirect)
	}

	data, err := io.ReadAll(rr.Body)
	if err != nil {
		if string(data) != "yandex.ru" {
			t.Errorf("handler returned wrong answer: got %v want %v",
				string(data), "yandex.ru")
		}
	}
}

// func TestMainRoute(t *testing.T) {
// 	type want struct {
// 		code     int
// 		response string
// 	}
// 	tests := []struct {
// 		urls        *urlstorage.URLStorage
// 		name        string
// 		method      string
// 		url         string
// 		want        want
// 		contenttype string
// 	}{
// 		{
// 			urls: &urlstorage.URLStorage{
// 				Store: map[string]string{},
// 			},
// 			name:        "test_get_nothing",
// 			method:      "GET",
// 			url:         "http://localhost:8080/",
// 			contenttype: "text/plain",
// 			want: want{
// 				code:     400,
// 				response: "Bad Request\n",
// 			},
// 		},
// 		{
// 			urls: &urlstorage.URLStorage{
// 				Store: map[string]string{
// 					"qwerty": "yandex.ru",
// 				},
// 			},
// 			name:        "test_get_yandex",
// 			method:      "GET",
// 			url:         "http://localhost:8080/qwerty",
// 			contenttype: "text/plain",
// 			want: want{
// 				code:     307,
// 				response: "",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var r *http.Request
// 			if tt.method == "POST" {
// 				r = httptest.NewRequest(tt.method, tt.url, strings.NewReader("yandex.ru"))
// 			} else {
// 				r = httptest.NewRequest(tt.method, tt.url, nil)
// 			}
// 			r.Header.Set("Content-Type", tt.contenttype)
// 			w := httptest.NewRecorder()
// 			urlstorage.SetDefaultURLStorage(tt.urls)
// 			mainHandler(w, r)
// 			assert.Equal(t, tt.want.response, w.Body.String(), "Ответ не совпадает")
// 			assert.Equal(t, tt.want.code, w.Code, "Код запроса не совпадает")
// 		})
// 	}
// }
