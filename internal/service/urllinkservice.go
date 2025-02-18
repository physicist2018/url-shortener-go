package service

import (
	"context"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/ports/randomstring"
	"github.com/rs/zerolog"
)

type URLLinkService struct {
	log       zerolog.Logger
	generator randomstring.RandomStringGenerator
	repo      domain.URLLinkRepo
}

func NewURLLinkService(repo domain.URLLinkRepo, generator randomstring.RandomStringGenerator, logger zerolog.Logger) *URLLinkService {
	return &URLLinkService{
		repo:      repo,
		generator: generator,
		log:       logger,
	}
}

// Метод создания короткой ссылки
func (u *URLLinkService) CreateShortURL(ctx context.Context, longURL string) (*domain.URLLink, error) {
	shortURL := u.generator.GenerateRandomString()
	link := &domain.URLLink{
		ShortURL: shortURL,
		LongURL:  longURL,
	}
	err := u.repo.Store(ctx, link)

	return link, err
}

// метод получения оригинальной ссылки
func (u *URLLinkService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	link, err := u.repo.Find(ctx, shortURL)

	if err != nil {
		u.log.Info().Err(err)
		return "", err
	}

	return link.LongURL, nil
}

// проверка соединения
func (u *URLLinkService) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}
