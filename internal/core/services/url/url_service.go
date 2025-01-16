package url

import (
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
	ports "github.com/physicist2018/url-shortener-go/internal/core/repository/url"
	"github.com/physicist2018/url-shortener-go/pkg/utils"
)

type URLService struct {
	urlRepo ports.URLRepository
}

func NewURLService(urlRepo ports.URLRepository) *URLService {
	return &URLService{
		urlRepo: urlRepo,
	}
}

// Создаем короткую ссылку
func (s *URLService) GenerateShortURL(originalURL string) (*urlmodels.URL, error) {
	shortURL := utils.GenerateShortURL()
	url := &urlmodels.URL{
		Original: originalURL,
		Short:    shortURL,
	}

	return s.urlRepo.Save(url)
}

// Получаем ссылку по короткому URL
func (s *URLService) GetOriginalURL(shortURL string) (*urlmodels.URL, error) {
	return s.urlRepo.FindByShort(shortURL)
}
