package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	m := Max(big.NewInt(1), big.NewInt(2))
	require.Equal(t, 0, big.NewInt(2).Cmp(m))
}

func TestMin(t *testing.T) {
	m := Min(big.NewInt(1), big.NewInt(2))
	require.Equal(t, 0, big.NewInt(1).Cmp(m))
}

func TestAccumulate(t *testing.T) {
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

func TestAddPercentage(t *testing.T) {
	tests := []struct {
		name       string
		value      *big.Int
		percentage uint16
		want       *big.Int
	}{
		{
			name:       "add 10% to 100",
			value:      big.NewInt(100),
			percentage: 10,
			want:       big.NewInt(110),
		},
		{
			name:       "add 0% to 1000",
			value:      big.NewInt(1000),
			percentage: 0,
			want:       big.NewInt(1000),
		},
		{
			name:       "add 10% to 0",
			value:      big.NewInt(0),
			percentage: 10,
			want:       big.NewInt(0),
		},
		{
			name:       "add 13% to 1998",
			value:      big.NewInt(1998),
			percentage: 13,
			want:       big.NewInt(2257), // Rounds down to nearest integer
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddPercentage(tt.value, tt.percentage); got.Cmp(tt.want) != 0 {
				t.Errorf("AddPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
