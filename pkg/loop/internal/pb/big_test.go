package loop

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigToBytes(t *testing.T) {
	for _, tc := range []struct {
		name string
		i    *big.Int
	}{
		{"nil", nil},
		{"zero", big.NewInt(0)},
		{"pos", big.NewInt(1)},
		{"neg", big.NewInt(-1)},
		{"big", new(big.Int).SetBytes(bytes.Repeat([]byte{0xFA}, 100))},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b := NewBigInt(tc.i)
			got := BigFromBigInt(b)
			require.Equal(t, tc.i, got)
		})
	}
}
