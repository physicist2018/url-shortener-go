package randomstring

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	return RandomStringWithAlphabet(n, alphabet)
}

func RandomStringWithAlphabet(n int, alphabet string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return b.String()
}
