package domain

import "context"

type URLLinkRepo interface {
	Store(ctx context.Context, urlLink URLLink) (URLLink, error)
	Find(ctx context.Context, shortURL string) (URLLink, error)
	FindAll(ctx context.Context, userID string) ([]URLLink, error)
	Ping(context.Context) error
	Close() error
}
