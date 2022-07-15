package SID

import (
	"math/rand"
	"time"
)

func SIDgenerator(n int) string {
	rand.Seed((time.Now().UnixNano() + int64(time.Now().Second()-time.Now().Hour()+time.Now().Day())) / int64(time.Now().Year()))
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
