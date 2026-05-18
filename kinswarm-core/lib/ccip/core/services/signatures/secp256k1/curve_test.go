package secp256k1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var group = &Secp256k1{}

func TestSecp256k1_String(t *testing.T) {
	require.Equal(t, group.String(), "Secp256k1")
}

func TestSecp256k1_Constructors(t *testing.T) {
	require.Equal(t, group.ScalarLen(), 32)
	require.Equal(t, ToInt(group.Scalar()), bigZero)
	require.Equal(t, group.PointLen(), 33)
	require.Equal(t, group.Point(), &secp256k1Point{fieldZero, fieldZero})
}
