package changeset

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitUnwrapAmount(t *testing.T) {
	t.Parallel()

	t.Run("balance twice the test amount", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Mul(TEST_AMOUNT_WEI, big.NewInt(2))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, TEST_AMOUNT_WEI, test)
		require.Equal(t, TEST_AMOUNT_WEI, remainder)
	})

	t.Run("balance just above the test amount", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Add(TEST_AMOUNT_WEI, big.NewInt(1))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, TEST_AMOUNT_WEI, test)
		require.Equal(t, big.NewInt(1), remainder)
	})

	t.Run("balance equal to the test amount", func(t *testing.T) {
		t.Parallel()

		_, _, err := splitUnwrapAmount(new(big.Int).Set(TEST_AMOUNT_WEI))
		require.Error(t, err)
	})

	t.Run("balance below the test amount", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Sub(TEST_AMOUNT_WEI, big.NewInt(1))
		_, _, err := splitUnwrapAmount(balance)
		require.Error(t, err)
	})

	t.Run("zero balance", func(t *testing.T) {
		t.Parallel()

		_, _, err := splitUnwrapAmount(big.NewInt(0))
		require.Error(t, err)
	})

	t.Run("nil balance", func(t *testing.T) {
		t.Parallel()

		_, _, err := splitUnwrapAmount(nil)
		require.Error(t, err)
	})

	t.Run("does not mutate the input balance", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Add(TEST_AMOUNT_WEI, big.NewInt(42))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, TEST_AMOUNT_WEI, test)
		require.Equal(t, big.NewInt(42), remainder)
		require.Equal(t, new(big.Int).Add(TEST_AMOUNT_WEI, big.NewInt(42)), balance)
	})
}
