package httphandlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	ports "github.com/physicist2018/url-shortener-go/internal/core/ports/urlports"
)

type URLHandler struct {
	urlService ports.URLService
	baseURL    string
}

func NewURLHandler(service ports.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: service,
		baseURL:    baseURL,
	}
}

// HandleGenerateShortURL is a function that handles the generation of a short URL.
// It checks if the request path is correct, reads the original URL from the request body, trims it and checks if it's not empty.
// If everything is correct, it generates a short URL using the urlService.
// If there is an error during the process, it returns a 400 Bad Request error.
func (h *URLHandler) HandleGenerateShortURL(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	originalURL, err := bufio.NewReader(r.Body).ReadString('\n')
	log.Println(originalURL)
	if (err != nil) && (err != io.EOF) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	originalURL = strings.TrimSpace(originalURL)
	log.Println(originalURL)
	if len(originalURL) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}
	shortURL, err := h.urlService.GenerateShortURL(originalURL)
	fullURL := fmt.Sprintf("%s/%s", h.baseURL, shortURL.Short)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fullURL))

}

// HandleRedirect is a function that handles the redirection to the original URL.
// It extracts the short URL from the request path, retrieves the original URL using the urlService,
// and redirects the user to the original URL.
// If there is an error during the process, it returns a 404 Not Found error.
func (h *URLHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.String()
	//shortURL := chi.URLParam(r, "shortURL")
	url, err := h.urlService.GetOriginalURL(shortURL)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.Original, http.StatusFound)
}
