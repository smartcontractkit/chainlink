package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestMax(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	m := Max(big.NewInt(1), big.NewInt(2))
	require.Equal(t, 0, big.NewInt(2).Cmp(m))
}

func TestMin(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	m := Min(big.NewInt(1), big.NewInt(2))
	require.Equal(t, 0, big.NewInt(1).Cmp(m))
}

func TestAccumulate(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	s := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
		big.NewInt(4),
		big.NewInt(5),
	}
	expected := big.NewInt(15)
	require.Equal(t, expected, Accumulate(s))
	s = []*big.Int{}
	expected = big.NewInt(0)
	require.Equal(t, expected, Accumulate(s))
}
