package httphandlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	originalURL, err := bufio.NewReader(r.Body).ReadString('\n')

	if errors.Is(err, nil) && !errors.Is(err, io.EOF) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
		return
	}

	originalURL = strings.TrimSpace(originalURL)
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
	shortURLparts := strings.Split(r.URL.String(), "/")
	shortURL := shortURLparts[len(shortURLparts)-1]
	url, err := h.urlService.GetOriginalURL(shortURL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Set("Location", url.Original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
