package memory

import (
	"errors"

	"sync"

	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
)

type URLRepositoryMap struct {
	mutex sync.RWMutex
	urls  map[string]urlmodels.URL
}

func NewURLRepositoryMap() *URLRepositoryMap {
	return &URLRepositoryMap{
		urls: make(map[string]urlmodels.URL),
	}
}

// Реализуем интерфеййс
func (r *URLRepositoryMap) Save(url urlmodels.URL) (urlmodels.URL, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urls[url.Short]; exists {
		return urlmodels.URL{}, errors.New("короткая ссылка уже есть")
	}

	r.urls[url.Short] = url
	return url, nil
}

func (r *URLRepositoryMap) FindByShort(shortURL string) (urlmodels.URL, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	url, exists := r.urls[shortURL]
	if !exists {
		return urlmodels.URL{}, errors.New("короткая ссылка не найдена")
	}

	return url, nil
}
