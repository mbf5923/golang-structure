package util

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateToken(userId string, email string) string {
	// generate accessToken sha256 from user id and email and random string and time
	apiToken := randStringRunes(32)
	sum := sha256.Sum256([]byte(apiToken + email + userId + time.DateTime))

	return hex.EncodeToString(sum[:])
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
