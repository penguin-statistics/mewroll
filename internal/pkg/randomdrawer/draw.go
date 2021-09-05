package randomdrawer

import (
	"crypto/rand"
	"math/big"
)

func Draw(lower, upper int) int {
	diff := upper - lower
	res, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		panic("failed to generate random number from crypto source: " + err.Error())
	}
	return int(res.Add(res, big.NewInt(int64(lower))).Int64())
}
