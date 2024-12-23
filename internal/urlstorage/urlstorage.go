package urlstorage

import (
	"errors"

	"github.com/physicist2018/url-shortener-go/internal/randomstring"
)

const (
	maxAttempts            = 3
	shortURLLength         = 6
	defaultStorageCapacity = 20
	TooManyAttempts        = "не удалось создать уникальную короткую ссылку за 3 попытки"
	LongURLNotFound        = "длинная ссылка не найдена"
	ShortURLNotFound       = "короткая ссылка не найдена"
	ShortURLExists         = "ссылка уже существует"
	ShortURLCreated        = "ссылка создана успешно"
)

// Юаза представляет собой список объектов типа URLItem
// по-сути - аналог таблицы
type URLStorage struct {
	Store []URLItem
}

var defaultURLStorage *URLStorage

// NewURLStorage creates a new instance of URLStorage with an empty slice of URLItem and a default capacity.
func NewURLStorage() *URLStorage {
	return &URLStorage{
		Store: make([]URLItem, 0, defaultStorageCapacity),
	}
}

// CreateShortURL creates a new short URL in the storage
func (s *URLStorage) CreateShortURL(longURL string) (string, error) {
	// добавляем новый короткий URL в хранилище
	for i := 0; i < maxAttempts; i++ {
		shortURL := randomstring.RandomString(shortURLLength)

		if _, err := s.FindShortURL(shortURL); err != nil {
			s.Store = append(s.Store, URLItem{
				LongURL:  longURL,
				ShortURL: shortURL,
			})
			return shortURL, nil
		}
	}

	return "", errors.New(TooManyAttempts)
}

// GetLongURL returns a long URL by a short URL
func (s *URLStorage) GetLongURL(shortURL string) (string, error) {
	// возвращаем длинную ссылку по короткой
	for i, v := range s.Store {
		if v.ShortURL == shortURL {
			return s.Store[i].LongURL, nil
		}
	}
	return "", errors.New(ShortURLNotFound)

}

// FindShortURL returns a short URL by a long URL
func (s *URLStorage) FindShortURL(longURL string) (string, error) {
	for i, v := range s.Store {
		if v.LongURL == longURL {
			return s.Store[i].ShortURL, nil
		}
	}
	return "", errors.New(LongURLNotFound)
}

// GetDefaultURLStorage returns a default URLStorage instance
func GetDefaultURLStorage() *URLStorage {
	return defaultURLStorage
}

// SetDefaultURLStorage sets a default URLStorage instance
func SetDefaultURLStorage(s *URLStorage) {
	defaultURLStorage = s
}

func init() {
	defaultURLStorage = NewURLStorage()
}
