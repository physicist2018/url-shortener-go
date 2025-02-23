package domain

import "context"

type URLLinkService interface {
	CreateShortURL(ctx context.Context, link URLLink) (URLLink, error)
	GetOriginalURL(ctx context.Context, link URLLink) (URLLink, error)
	MarkURLAsDeleted(ctx context.Context, links []URLLink) error
	//DeleteShortLinksForUser(ctx context.Context, userID string, shortURLs []string) error
	FindAll(ctx context.Context, userID string) ([]URLLink, error)
	Ping(ctx context.Context) error
}
