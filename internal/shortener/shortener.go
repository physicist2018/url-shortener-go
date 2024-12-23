package shortener

import (
	"bufio"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/physicist2018/url-shortener-go/internal/urlstorage"
)

// RunServer starts the server
func RunServer() error {
	router := chi.NewMux()
	router.Use(middleware.AllowContentType("text/plain"))
	router.Route("/", func(r chi.Router) {
		r.Post("/", postRoute)
		r.Get("/{shortURL}", getRoute)
	})

	//router.HandleFunc("POST /", postRoute)
	//router.HandleFunc("GET /", getRoute)

	return http.ListenAndServe(`:8080`, router)
}

// postRoute is the handler for POST request
func postRoute(w http.ResponseWriter, r *http.Request) {
	url, err := bufio.NewReader(r.Body).ReadString('\n')
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	if (err != nil) && (err != io.EOF) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	url = strings.TrimSpace(url)
	if len(url) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	if shortURL, err := urlstorage.GetDefaultURLStorage().FindShortURL(url); err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + shortURL))
		return
	}

	if shortURL, err := urlstorage.GetDefaultURLStorage().CreateShortURL(url); err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + shortURL))
		return
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
}

func getRoute(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	if longURL, err := urlstorage.GetDefaultURLStorage().GetLongURL(shortURL); err == nil {
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 404
}
