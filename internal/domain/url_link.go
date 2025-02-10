package domain

type URLLink struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}
