package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijkmlnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min,max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomStrings(n int) string {
	var sb strings.Builder  // 字符串生成对象
	k := len(alphabet)

	for i:=0;i<n;i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a random owner name
func RandomOwner() string {
	return RandomStrings(6)
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0,1000)
}

// RanomCurrency generate a random currency code
func RandomCurrrency() string {
	currencies := []string{"USD","EUR","CAD"}

	n := len(currencies)
	return currencies[rand.Intn(n)]
}