package values

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListUnwrapTo(t *testing.T) {
	t.Run("[]int", func(t *testing.T) {
		expected := []int{1, 2, 3}
		got := []int{}
		sliceTest[int](t, expected, got)
	})

	t.Run("[]int64", func(t *testing.T) {
		expected := []int64{1, 2, 3}
		got := []int64{}
		sliceTest[int64](t, expected, got)
	})

	t.Run("[]string", func(t *testing.T) {
		expected := []string{"hello", "world"}
		got := []string{}
		sliceTest[string](t, expected, got)
	})

	t.Run("[][]byte", func(t *testing.T) {
		expected := [][]byte{[]byte("hello"), []byte("world")}
		got := [][]byte{}
		sliceTest[[]byte](t, expected, got)
	})

	t.Run("[]decimal.Decimal", func(t *testing.T) {
		expected := []decimal.Decimal{decimal.NewFromFloat(1.00), decimal.NewFromFloat(1.32)}
		got := []decimal.Decimal{}
		sliceTest[decimal.Decimal](t, expected, got)
	})

	t.Run("[]bool", func(t *testing.T) {
		expected := []bool{true, false}
		got := []bool{}
		sliceTest[bool](t, expected, got)
	})

	t.Run("[]big.Int", func(t *testing.T) {
		expected := []big.Int{*big.NewInt(1), *big.NewInt(2)}
		got := []big.Int{}
		sliceTest[big.Int](t, expected, got)
	})
}

func sliceTest[T any](t *testing.T, expected []T, got []T) {
	v, err := Wrap(expected)
	require.NoError(t, err)

	err = v.UnwrapTo(&got)
	require.NoError(t, err)

	assert.Equal(t, expected, got)

	gotn := (*[]T)(nil)
	err = v.UnwrapTo(gotn)
	assert.ErrorContains(t, err, "cannot unwrap to nil pointer")
}
