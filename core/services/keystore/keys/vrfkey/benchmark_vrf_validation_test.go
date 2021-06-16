package vrfkey

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

// Run with `go test -bench BenchmarkProofValidation`
func BenchmarkProofValidation(b *testing.B) {
	key := NewPrivateKeyXXXTestingOnly(big.NewInt(1))
	var proofs []Proof
	for i := 0; i < b.N; i++ {
		p, err := key.GenerateProof(big.NewInt(int64(i)))
		require.NoError(b, err, "failed to generate proof number %d", i)
		proofs = append(proofs, p)
	}
	b.ResetTimer()
	for i, p := range proofs {
		isValid, err := p.VerifyVRFProof()
		require.NoError(b, err, "failed to check proof number %d", i)
		require.True(b, isValid, "proof number %d is invalid", i)
	}
}
