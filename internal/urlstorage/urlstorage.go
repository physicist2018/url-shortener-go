package urlstorage

import (
	"errors"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/randomstring"
)

// Структура для хранения ссылок
type URLStorage struct {
	Store      map[string]string
	sync.Mutex // для синхронизации
}

// HasLongURL проверяет есть ли длинная ссылка в хранилище, если есть, возвращает
// короткую сслку и true, если нет, возвращает false, при этом shortURL не валиден
func (s *URLStorage) HasLongURL(longURL string) (string, bool) {
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

// HasShortURL проверяет есть ли короткая ссылка в хранилище,
// если есть - возвращает true, если нет, возвращает false
func (s *URLStorage) HasShortURL(shortURL string) bool {
	_, ok := s.Store[shortURL]
	return ok
}

// AddURL добавляет новую ссылку в хранилище, возвращает короткую ссылку и nil,
// если ссылка добавилась успешно иначе возвращает ошибку
// На генерацию случайной ссылки есть три попытки
func (s *URLStorage) AddURL(longURL string) (string, error) {
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

// GetURL возвращает короткую ссылку из хранилища и nil, если ссылка не найдена - nil и ошибка
func (s *URLStorage) GetURL(shortURL string) (string, error) {
	if val, ok := s.Store[shortURL]; ok {
		return val, nil
	}
	return "", errors.New("not found")
}
