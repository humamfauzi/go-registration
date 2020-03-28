package libtest

import (
	"math/rand"
	"time"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	stringBytes := make([]byte, length)
	var randomInteger int
	for i := 0; i < length; i++ {
		randomInteger = rand.Intn(len(letters))
		stringBytes[i] = letters[randomInteger]
	}
	return string(stringBytes)
}
