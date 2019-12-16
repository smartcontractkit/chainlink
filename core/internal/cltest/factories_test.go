package cltest

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func i(x int64) *big.Int { return big.NewInt(x) }

func TestBigHexInt(t *testing.T) {
	asBig := i(0).Sub(i(0).Exp(i(2), i(64), big.NewInt(0)), i(1)) // 2**64-1
	x := asBig.Uint64()
	newBig := BigHexInt(x)
	assert.Equal(t, (*big.Int)(&newBig).Uint64(), x)
}
