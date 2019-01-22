package rsavrf

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSAVRF_PrimeKeySize(t *testing.T) {
	bitLen := 1024
	p, err := rand.Prime(rand.Reader, bitLen)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, p.BitLen(), bitLen)
}

func TestRSAVRF_safePrime(t *testing.T) {
	// Short, because this is slow. Greater than 64+1, because rand.Prime
	// logic changes for smaller bit lengths
	bitLen := 66
	p := safePrime(uint32(bitLen))
	assert.Equal(t, p.BitLen(), bitLen)
}

func TestRSAVRF_MakeKey(t *testing.T) {
	k, err := MakeKey(150)
	if err != nil {
		panic(err)
	}
	assert.True(t, k.Primes[0].Cmp(big.NewInt(0)) == 1)
}
