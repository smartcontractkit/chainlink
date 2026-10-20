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

		balance := new(big.Int).Mul(testAmountWei, big.NewInt(2))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, testAmountWei, test)
		require.Equal(t, testAmountWei, remainder)
	})

	t.Run("balance just above the test amount", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Add(testAmountWei, big.NewInt(1))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, testAmountWei, test)
		require.Equal(t, big.NewInt(1), remainder)
	})

	t.Run("balance equal to the test amount", func(t *testing.T) {
		t.Parallel()

		_, _, err := splitUnwrapAmount(new(big.Int).Set(testAmountWei))
		require.Error(t, err)
	})

	t.Run("balance below the test amount", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Sub(testAmountWei, big.NewInt(1))
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

		balance := new(big.Int).Add(testAmountWei, big.NewInt(42))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, testAmountWei, test)
		require.Equal(t, big.NewInt(42), remainder)
		require.Equal(t, new(big.Int).Add(testAmountWei, big.NewInt(42)), balance)
	})

	t.Run("test and remainder are independent copies", func(t *testing.T) {
		t.Parallel()

		// Mutating the returned values must not leak into TEST_AMOUNT_WEI or the input,
		// otherwise a later chain's split (or a later run) would use a corrupted test amount.
		balance := new(big.Int).Mul(testAmountWei, big.NewInt(3))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)

		test.Add(test, big.NewInt(1))
		remainder.Add(remainder, big.NewInt(1))

		require.Equal(t, testAmountWei, test.Sub(test, big.NewInt(1)))
		require.Equal(t, new(big.Int).Mul(testAmountWei, big.NewInt(2)), remainder.Sub(remainder, big.NewInt(1)))
	})

	t.Run("error mentions both balances", func(t *testing.T) {
		t.Parallel()

		_, _, err := splitUnwrapAmount(big.NewInt(1))
		require.ErrorContains(t, err, "1")
		require.ErrorContains(t, err, testAmountWei.String())
	})

	t.Run("test amount plus remainder equals the balance", func(t *testing.T) {
		t.Parallel()

		balance := new(big.Int).Mul(testAmountWei, big.NewInt(7))
		test, remainder, err := splitUnwrapAmount(balance)
		require.NoError(t, err)
		require.Equal(t, balance, new(big.Int).Add(test, remainder))
	})
}
