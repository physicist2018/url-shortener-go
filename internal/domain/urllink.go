package domain

type URLLink struct {
	UserID   string `json:"user_id" db:"user_id"`
	ShortURL string `json:"short_url" db:"short_url"`
	LongURL  string `json:"original_url" db:"original_url"`
}
