package url

import (
	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
	"github.com/physicist2018/url-shortener-go/internal/core/ports/randomstring"
	ports "github.com/physicist2018/url-shortener-go/internal/core/repository/url"
)

type URLService struct {
	randomStringGenerator randomstring.RandomStringGenerator
	urlRepo               ports.URLRepository
}

func NewURLService(urlRepo ports.URLRepository, rndStringGenerator randomstring.RandomStringGenerator) *URLService {
	return &URLService{
		randomStringGenerator: rndStringGenerator,
		urlRepo:               urlRepo,
	}
}

// Создаем короткую ссылку
func (s *URLService) GenerateShortURL(originalURL string) (urlmodels.URL, error) {

	var shortURL string
	for {
		shortURL = s.randomStringGenerator.GenerateRandomString()
		_, err := s.urlRepo.FindByShort(shortURL)

		if err != nil {
			break
		}
	}

	url := urlmodels.URL{
		Original: originalURL,
		Short:    shortURL,
	}
	return s.urlRepo.Save(url)
}

// Получаем ссылку по короткому URL
func (s *URLService) GetOriginalURL(shortURL string) (urlmodels.URL, error) {
	return s.urlRepo.FindByShort(shortURL)
}
