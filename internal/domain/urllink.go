package domain

type URLLink struct {
	ShortURL string `json:"short_url" db:"short_url"`
	LongURL  string `json:"original_url" db:"original_url"`
}
