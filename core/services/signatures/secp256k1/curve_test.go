package secp256k1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var group = &Secp256k1{}

func TestSecp256k1_String(t *testing.T) {
	require.Equal(t, "Secp256k1", group.String())
}

func TestSecp256k1_Constructors(t *testing.T) {
	require.Equal(t, 32, group.ScalarLen())
	require.Equal(t, ToInt(group.Scalar()), bigZero)
	require.Equal(t, 33, group.PointLen())
	require.Equal(t, &secp256k1Point{fieldZero, fieldZero}, group.Point())
}
