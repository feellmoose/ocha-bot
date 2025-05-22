package helper

import (
	"math/rand/v2"
)

func RandomDensity(min, max float64, refactor func(float64) float64) float64 {
	random := rand.Float64()
	return max + (max-min)*refactor(random)
}

func Num(min, max int) int {
	return rand.IntN(max-min+1) + min
}

const ascii = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ascii[rand.IntN(len(ascii))]
	}
	return string(b)
}
