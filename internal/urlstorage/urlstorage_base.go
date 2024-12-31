package urlstorage

// URLItem - то, как хранится запись сопоставления короткого и длинного URL
type URLItem struct {
	ShortURL string
	LongURL  string
}

// Интерфейс репозиторий для хранения сопоставления коротких и длинных URL
type URLStorageRepository interface {
	CreateShortURL(longURL string) (string, error)
	GetLongURL(shortURL string) (string, error)
	FindShortURL(longURL string) (string, error)
}
