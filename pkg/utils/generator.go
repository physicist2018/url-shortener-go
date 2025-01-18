package utils

import (
	"math/rand"
)

type RandomString struct {
	length    int
	generator *rand.Rand
}

func NewRandomString(length int, rnd *rand.Rand) *RandomString {
	return &RandomString{
		length:    length,
		generator: rnd,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString генерирует короткую ссылку из случайных символов
func (rs *RandomString) GenerateRandomString() string {
	shortURL := make([]byte, rs.length)
	for i := range shortURL {
		shortURL[i] = charset[rs.generator.Intn(len(charset))]
	}
	return string(shortURL)
}
