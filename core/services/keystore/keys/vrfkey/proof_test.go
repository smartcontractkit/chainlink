package vrfkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVRF_VerifyProof(t *testing.T) {
	sk, err := NewV2()
	require.NoError(t, err)
	seed, nonce := big.NewInt(2), big.NewInt(3)
	p, err := sk.GenerateProofWithNonce(seed, nonce)
	require.NoError(t, err, "could not generate proof")
	p.Seed = big.NewInt(0).Add(seed, big.NewInt(1))
	valid, err := p.VerifyVRFProof()
	require.NoError(t, err, "could not validate proof")
	require.False(t, valid, "invalid proof was found valid")
}
