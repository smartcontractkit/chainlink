package starknet

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	s, err := SignatureFromBigInts(big.NewInt(7),
		big.NewInt(11))

	require.NoError(t, err)

	x, y, err := s.Ints()
	require.NoError(t, err)
	require.Equal(t, big.NewInt(7), x)
	require.Equal(t, big.NewInt(11), y)

	b, err := s.Bytes()
	require.NoError(t, err)
	require.NotNil(t, b)

	roundTrip, err := SignatureFromBytes(b)
	require.NoError(t, err)
	// compare fields because comparing proto objects directly gave incorrect failures
	require.EqualValues(t, s.sig.X.GetValue(), roundTrip.sig.X.GetValue())
	require.EqualValues(t, s.sig.X.GetNegative(), roundTrip.sig.X.GetNegative())
	require.EqualValues(t, s.sig.Y.GetValue(), roundTrip.sig.Y.GetValue())
	require.EqualValues(t, s.sig.Y.GetNegative(), roundTrip.sig.Y.GetNegative())

	// no negative allowed
	s, err = SignatureFromBigInts(big.NewInt(-7),
		big.NewInt(11))
	require.Error(t, err)
	require.Nil(t, s)
}
