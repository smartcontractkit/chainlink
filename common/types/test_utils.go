package types

import (
	"math"
	"math/big"
	"math/rand"
)

func RandomID() ID {
	id := rand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

func NewIDFromInt(id int64) ID {
	return big.NewInt(id)
}
