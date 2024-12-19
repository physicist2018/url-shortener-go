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

type URLStorage struct {
	Store      map[string]string
	sync.Mutex // для синхронизации
}

func (s *URLStorage) HasLongUrl(longURL string) (string, bool) {
	ok := false
	var shortURL string

	for key, val := range s.Store {
		if val == longURL {
			ok = true
			shortURL = key
			break
		}
	}

	return shortURL, ok
}

func (s *URLStorage) HasShortUrl(shortURL string) bool {
	_, ok := s.Store[shortURL]
	return ok
}

func (s *URLStorage) AddUrl(longURL string) (string, error) {
	var shortURL string
	var ok bool

	for i := 0; i < 3; i++ {
		shortURL = randomstring.RandomString(10)
		if _, ok = s.Store[shortURL]; !ok {
			break
		}
	}
	if !ok {
		s.Lock()
		s.Store[shortURL] = longURL
		s.Unlock()
		return shortURL, nil
	}
	return "", errors.New("too many attempts")
}

func (s *URLStorage) GetUrl(shortURL string) (string, error) {
	if val, ok := s.Store[shortURL]; ok {
		return val, nil
	}
	return "", errors.New("not found")
}

var urlStorage *URLStorage

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

	if (r.Method == http.MethodPost) && strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
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
	if len(url) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

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
	urlStorage = &URLStorage{
		Store: map[string]string{},
	}
}
