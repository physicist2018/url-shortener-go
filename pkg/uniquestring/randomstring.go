package uniquestring

import (
	"time"

	"math/rand"
)

const (
	RandomStringLength = 5
)

// Стратегия для генерации случайной строки на основе rand
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

func NewRandomStringDefault() *RandomString {
	return &RandomString{
		length:    RandomStringLength,
		generator: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func NewRandomStringFixed() *RandomString {
	return &RandomString{
		length:    RandomStringLength,
		generator: rand.New(rand.NewSource(10)),
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Реализация метода интерфейса для генерации случайной строки
func (rs *RandomString) Generate() string {
	shortURL := make([]byte, rs.length)
	for i := range shortURL {
		shortURL[i] = charset[rs.generator.Intn(len(charset))]
	}
	return string(shortURL)
}
