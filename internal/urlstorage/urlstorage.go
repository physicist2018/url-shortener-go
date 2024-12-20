package urlstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/randomstring"
)

const (
	maxAttempts     = 3
	shortURLLength  = 6
	TooManyAttempts = "не удалось создать уникальную короткую ссылку за 3 попытки"
	LongURLNotFound = "длинная ссылка не найдена"
	ShortURLExists  = "ссылка уже существует"
	ShortURLCreated = "ссылка создана успешно"
)

// Структура для хранения ссылок
type URLStorage struct {
	Store      map[string]string
	sync.Mutex // для синхронизации
}

var defaultURLStorage *URLStorage = &URLStorage{
	Store: make(map[string]string, 0),
}

// CreateShortURL создает уникальную короткую ссылку и возвращает ее, а также error
func (s *URLStorage) CreateShortURL(longURL string) (string, error) {

	if shortURL, ok := s.hasLongURL(longURL); ok {
		return shortURL, nil
	}

	return s.addURL(longURL)
}

// GetURL возвращает короткую ссылку из хранилища и nil, если ссылка не найдена - nil и ошибка
func (s *URLStorage) GetURL(shortURL string) (string, error) {
	if val, ok := s.Store[shortURL]; ok {
		return val, nil
	}
	return "", errors.New(LongURLNotFound)
}

func (s *URLStorage) FindShortURL(longURL string) (string, error) {

	if shortURL, ok := s.hasLongURL(longURL); ok {
		return shortURL, nil
	} else {
		return "", errors.New(LongURLNotFound)
	}
}

// hasLongURL проверяет есть ли длинная ссылка в хранилище, если есть, возвращает
// короткую сслку и true, если нет, возвращает false, при этом shortURL не валиден
func (s *URLStorage) hasLongURL(longURL string) (string, bool) {
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

// addURL добавляет новую ссылку в хранилище, возвращает короткую ссылку и nil,
// если ссылка добавилась успешно иначе возвращает ошибку
// На генерацию случайной ссылки есть три попытки
func (s *URLStorage) addURL(longURL string) (string, error) {
	var shortURL string
	var ok bool

	for i := 0; i < maxAttempts; i++ {
		shortURL = randomstring.RandomString(shortURLLength)
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
	return "", errors.New(TooManyAttempts)
}

func GetDefaultUrlStorage() *URLStorage {
	return defaultURLStorage
}

func SetDefaultUrlStorage(s *URLStorage) {
	fmt.Println("New storage setup")
	defaultURLStorage = s
}
