package util

import (
	"math/rand"
)

// Random int generates a random integer between min and max (inclusive).
func RandomInt(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n using letters from the English alphabet.
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "AUD", "JPY"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}
