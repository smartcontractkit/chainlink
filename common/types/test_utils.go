package types

import (
	"math"
	"math/big"
	mrand "math/rand"
)

func RandomID() ID {
	id := mrand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

func NewIDFromInt(id int64) ID {
	return big.NewInt(id)
}
