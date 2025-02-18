package domain

import "context"

type URLLinkService interface {
	CreateShortURL(ctx context.Context, longURL string) (*URLLink, error)
	GetOriginalURL(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) error
}
