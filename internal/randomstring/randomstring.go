package randomstring

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a random string of length n.
func RandomString(n int) string {
	return RandomStringWithAlphabet(n, alphabet)
}

// RandomStringWithAlphabet returns a random string of length n, with the given alphabet.
func RandomStringWithAlphabet(n int, alphabet string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}
	return b.String()
}
