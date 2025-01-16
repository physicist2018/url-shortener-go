package urlports

import "github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"

type URLService interface {
	GenerateShortURL(originalURL string) (*urlmodels.URL, error)
	GetOriginalURL(shortURL string) (*urlmodels.URL, error)
}
