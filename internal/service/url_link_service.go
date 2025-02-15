package service

import (
	"context"
	"errors"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/ports/randomstring"
	"github.com/physicist2018/url-shortener-go/internal/repository/repoerrors"
)

type URLLinkService struct {
	generator randomstring.RandomStringGenerator
	repo      domain.URLLinkRepo
}

func NewURLLinkService(repo domain.URLLinkRepo, generator randomstring.RandomStringGenerator) *URLLinkService {
	return &URLLinkService{
		repo:      repo,
		generator: generator,
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

	if errors.Is(err, repoerrors.ErrURLAlreadyInDB) {
		return link, err
	} else if err != nil {
		return nil, err
	}

	return link, nil
}

// метод получения оригинальной ссылки
func (u *URLLinkService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	link, err := u.repo.Find(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return link.LongURL, nil
}

// проверка соединения
func (u *URLLinkService) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}
