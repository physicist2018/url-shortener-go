package ports

import (
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
)

// URLRepository описывает репозиторий для хранения URL
type URLRepository interface {
	Save(url *urlmodels.URL) (*urlmodels.URL, error)
	FindByShort(shortURL string) (*urlmodels.URL, error)
}
