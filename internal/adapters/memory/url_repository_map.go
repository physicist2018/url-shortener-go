package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"sync"

	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
)

var ErrorOpeningFile = errors.New("ошибка открытия файла")
var ErrorCreatingFile = errors.New("ошибка создания файла")

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

func (r *URLRepositoryMap) DumpToFile(fullFilePath string) error {
	file, err := os.Create(fullFilePath)
	if err != nil {
		return ErrorCreatingFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "")
	for _, url := range r.urls {
		if err := encoder.Encode(url); err != nil {
			return fmt.Errorf("error encoding data to json: %w", err)
		}
	}

	return nil
}

func (r *URLRepositoryMap) RestoreFromFile(fullFilePath string) error {
	// Открываем файл для чтения
	file, err := os.Open(fullFilePath)
	if err != nil {
		return ErrorOpeningFile
	}
	defer file.Close()

	// Чтение файла построчно
	decoder := json.NewDecoder(file)
	for {
		var url urlmodels.URL
		// Декодируем одну строку в структуру
		if err := decoder.Decode(&url); err != nil {
			// Если достигнут конец файла (EOF), выходим из цикла
			if err.Error() == "EOF" {
				break
			}
			// В случае другой ошибки
			return fmt.Errorf("error decoding json: %w", err)
		}
		// Добавляем URL в срез
		r.urls[url.Short] = url
	}

	return nil
}
