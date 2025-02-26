package service

import (
	"context"

	"github.com/physicist2018/url-shortener-go/internal/domain"
	"github.com/physicist2018/url-shortener-go/internal/stringgenstategy"
	"github.com/rs/zerolog"
)

type URLLinkService struct {
	log       zerolog.Logger
	generator stringgenstategy.StringGeneratorContext
	repo      domain.URLLinkRepo
}

func NewURLLinkService(repo domain.URLLinkRepo, generator stringgenstategy.StringGeneratorContext, logger zerolog.Logger) *URLLinkService {
	return &URLLinkService{
		repo:      repo,
		generator: generator,
		log:       logger,
	}
}

// Метод создания короткой ссылки
func (u *URLLinkService) CreateShortURL(ctx context.Context, link domain.URLLink) (domain.URLLink, error) {
	shortURL := u.generator.GenerateString()
	urllink := domain.URLLink{
		ShortURL: shortURL,
		LongURL:  link.LongURL,
		UserID:   link.UserID,
	}

	return u.repo.Store(ctx, urllink)
}

// метод получения оригинальной ссылки
func (u *URLLinkService) GetOriginalURL(ctx context.Context, link domain.URLLink) (domain.URLLink, error) {
	link, err := u.repo.Find(ctx, link.ShortURL)

	if err != nil {
		u.log.Info().Err(err)
		return domain.URLLink{}, err
	}

	return link, nil
}

func (u *URLLinkService) FindAll(ctx context.Context, userID string) ([]domain.URLLink, error) {
	return u.repo.FindAll(ctx, userID)
}

// проверка соединения
func (u *URLLinkService) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}

func (u *URLLinkService) MarkURLsAsDeleted(ctx context.Context, links []domain.URLLink) error {
	return u.repo.MarkDeletedBatch(ctx, links)
}
