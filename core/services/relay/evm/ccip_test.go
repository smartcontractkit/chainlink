package evm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CCIPSubjectUUID(t *testing.T) {
	// We want the function to be
	// (1) an actual function (i.e., deterministic)
	assert.Equal(t, chainToUUID(big.NewInt(1)), chainToUUID(big.NewInt(1)))
	// (2) injective (produce different results for different inputs)
	assert.NotEqual(t, chainToUUID(big.NewInt(1)), chainToUUID(big.NewInt(2)))
	// (3) stable across runs
	assert.Equal(t, "c980e777-c95c-577b-83f6-ceb26a1a982d", chainToUUID(big.NewInt(1)).String())
}
