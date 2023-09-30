package util

import (
	"math/rand"
	"strings"
	"time"
)

const aplhabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(aplhabet)
	for i := 0; i < n; i++ {
		c := aplhabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()

}

func RandomOwner() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{
		"IDR", "USD", "EUR",
	}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
