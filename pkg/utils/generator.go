package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortLength = 6

// generateShortURL генерирует короткую ссылку из случайных символов
func GenerateShortURL() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	shortURL := make([]byte, shortLength)
	for i := range shortURL {
		shortURL[i] = charset[rnd.Intn(len(charset))]
	}
	return string(shortURL)
}
