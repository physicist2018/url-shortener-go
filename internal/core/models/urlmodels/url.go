package urlmodels

// это данные нашего приложения, model
type URL struct {
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}
