package shortener

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/randomstring"
)

type UrlStorage struct {
	Store      map[string]string
	sync.Mutex // для синхронизации
}

func (s *UrlStorage) HasLongUrl(longUrl string) (string, bool) {
	ok := false
	var shortUrl string

	for key, val := range s.Store {
		if val == longUrl {
			ok = true
			shortUrl = key
			break
		}
	}

	return shortUrl, ok
}

func (s *UrlStorage) HasShortUrl(shortUrl string) bool {
	_, ok := s.Store[shortUrl]
	return ok
}

func (s *UrlStorage) AddUrl(longUrl string) (string, error) {
	var shortUrl string
	var ok bool

	for i := 0; i < 3; i++ {
		shortUrl = randomstring.RandomString(10)
		if _, ok = s.Store[shortUrl]; !ok {
			break
		}
	}
	if !ok {
		s.Lock()
		s.Store[shortUrl] = longUrl
		s.Unlock()
		return shortUrl, nil
	}
	return "", errors.New("too many attempts")
}

func (s *UrlStorage) GetUrl(shortUrl string) (string, error) {
	if val, ok := s.Store[shortUrl]; ok {
		return val, nil
	}
	return "", errors.New("not found")
}

var urlStorage *UrlStorage

func RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	return http.ListenAndServe(`:8080`, mux)
}

// mainHandler is the handler for the main route
// it examines method and call necessary callback function
// in order to handle request properly
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Срзу отсекаем другий тип контента кроме text/plain
	if (r.Method == http.MethodPost) && (r.Header.Get("Content-Type") == "text/plain") {
		postRoute(w, r)

	} else if r.Method == http.MethodGet {
		getRoute(w, r)

	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
	}
}

// postRoute is the handler for POST request
func postRoute(w http.ResponseWriter, r *http.Request) {
	url, err := bufio.NewReader(r.Body).ReadString('\r')

	if (err != nil) && (err != io.EOF) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	url = strings.TrimSpace(url)
	if shortUrl, ok := urlStorage.HasLongUrl(url); ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortUrl))
		return
	}

	if shortUrl, err := urlStorage.AddUrl(url); err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortUrl))
		return
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
}

func getRoute(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	if longUrl, err := urlStorage.GetUrl(shortUrl); err == nil {
		w.Header().Set("Location", longUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound) // 404
}

func init() {
	urlStorage = &UrlStorage{
		Store: map[string]string{},
	}
}
