package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/physicist2018/url-shortener-go/internal/domain"
)

type InMemoryLinkRepository struct {
	links  map[string]*domain.URLLink
	mu     sync.RWMutex
	dbfile *os.File
}

func NewInMemoryLinkRepository(dbFilePath string) (*InMemoryLinkRepository, error) {
	repo := &InMemoryLinkRepository{
		links: make(map[string]*domain.URLLink),
	}

	// Открываем файл для добавления данных
	file, err := os.OpenFile(dbFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	repo.dbfile = file

	if err := repo.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return repo, nil
}

func (m *InMemoryLinkRepository) Store(ctx context.Context, urllink *domain.URLLink) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.links[urllink.ShortURL] = urllink

	// Добавляем данные в файл
	data, err := json.Marshal(urllink)
	if err != nil {
		return err
	}

	_, err = m.dbfile.Write(append(data, '\n'))

	return err
}

func (m *InMemoryLinkRepository) Find(ctx context.Context, shortURL string) (*domain.URLLink, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	urllink, exists := m.links[shortURL]
	if !exists {
		return nil, errors.New("ссылка не найдена")
	}

	return urllink, nil
}

func (m *InMemoryLinkRepository) Ping(ctx context.Context) error {
	return nil
}

func (m *InMemoryLinkRepository) load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.dbfile.Name())
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var urllink domain.URLLink
		if err := json.Unmarshal([]byte(line), &urllink); err != nil {
			return err
		}
		m.links[urllink.ShortURL] = &urllink
	}
	return nil
}

func (m *InMemoryLinkRepository) Close() error {
	return m.dbfile.Close()
}
