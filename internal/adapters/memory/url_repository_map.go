package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"sync"

	"github.com/physicist2018/url-shortener-go/internal/core/models/urlmodels"
)

var (
	ErrorOpeningFileWhenRestore = errors.New("ошибка открытия файла при восстановлении базы")
	ErrorCreatingFileWhenDump   = errors.New("ошибка создания файла при сохранении базы")
	ErrorShortURLAlreadyInDB    = errors.New("короткая ссылка уже есть в базе")
	ErrorShortURLNotFound       = errors.New("короткая ссылка не найдена в базе")
	ErrorSyncDB                 = errors.New("ошибка синхронизации базы")
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
		return urlmodels.URL{}, ErrorShortURLAlreadyInDB
	}

	r.urls[url.Short] = url
	return url, nil
}

func (r *URLRepositoryMap) FindByShort(shortURL string) (urlmodels.URL, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	url, exists := r.urls[shortURL]
	if !exists {
		return urlmodels.URL{}, ErrorShortURLNotFound
	}

	return url, nil
}

// Этот метод реализует сохранение мапы в файл, при этом, если возникли проблемы при открытии
// файла, мы просто возвращаем ошибку, которая в вызывающем коде должна вызвать панику
func (r *URLRepositoryMap) Dump(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return ErrorCreatingFileWhenDump
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

// Этот метод должен реализовать поддержку загрузки списка сопоставлений из файла
// в случае возниконвания ошибки невозможности загрузки, мы просто инициализируем чистую мапу
// при этом сигнализируем об этом
func (r *URLRepositoryMap) Restore(filename string) error {
	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		return ErrorOpeningFileWhenRestore
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
