package randoms

import (
	"crypto/rand"
	"math/big"
)

func GetRandomNumber(l, r int64) int64 {
	num, _ := rand.Int(rand.Reader, big.NewInt(r-l+1))
	return num.Int64() + l
}
