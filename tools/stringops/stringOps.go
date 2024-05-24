package stringops

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rng.Intn(len(chars))])
	}
	return b.String()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

