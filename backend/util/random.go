package util

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	str := StringWithCharset(length, charset)
	return str
}

func RondomBio(length int) *string {
	bio := RandomString(length)
	return &bio
}

func RandomEmail(length int) string {
	var emailBuilder strings.Builder
	emailBuilder.WriteString(RandomString(length))
	emailBuilder.WriteString("@mail.com")
	return emailBuilder.String()
}
