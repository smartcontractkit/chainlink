package bigmath

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	testCases := []struct {
		x        interface{}
		y        interface{}
		expected *big.Int
	}{
		{
			x:        int32(1),
			y:        int32(2),
			expected: big.NewInt(2),
		},
		{
			x:        big.NewInt(1),
			y:        big.NewInt(2),
			expected: big.NewInt(2),
		},
		{
			x:        float64(1.0),
			y:        float64(2.0),
			expected: big.NewInt(2),
		},
		{
			x:        "1",
			y:        "2",
			expected: big.NewInt(2),
		},
		{
			x:        uint(1),
			y:        uint(2),
			expected: big.NewInt(2),
		},
	}
	for _, testCase := range testCases {
		m := Max(testCase.x, testCase.y)
		require.Equal(t, 0, testCase.expected.Cmp(m))
	}
}

func TestMin(t *testing.T) {
	testCases := []struct {
		x        interface{}
		y        interface{}
		expected *big.Int
	}{
		{
			x:        int32(1),
			y:        int32(2),
			expected: big.NewInt(1),
		},
		{
			x:        big.NewInt(1),
			y:        big.NewInt(2),
			expected: big.NewInt(1),
		},
		{
			x:        float64(1.0),
			y:        float64(2.0),
			expected: big.NewInt(1),
		},
		{
			x:        "2",
			y:        "1",
			expected: big.NewInt(1),
		},
		{
			x:        uint(2),
			y:        uint(1),
			expected: big.NewInt(1),
		},
	}
	for _, testCase := range testCases {
		m := Min(testCase.x, testCase.y)
		require.Equal(t, 0, testCase.expected.Cmp(m))
	}
}

func TestAccumulate(t *testing.T) {
	s := []interface{}{1, 2, 3, 4, 5}
	expected := big.NewInt(15)
	require.Equal(t, expected, Accumulate(s))
	s = []interface{}{}
	expected = big.NewInt(0)
	require.Equal(t, expected, Accumulate(s))
}
