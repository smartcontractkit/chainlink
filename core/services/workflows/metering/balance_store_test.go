package metering

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalanceStore(t *testing.T) {
	t.Parallel()

	one := decimal.NewFromInt(1)
	two := decimal.NewFromInt(2)
	five := decimal.NewFromInt(5)
	seven := decimal.NewFromInt(7)
	nine := decimal.NewFromInt(9)
	ten := decimal.NewFromInt(10)
	eleven := decimal.NewFromInt(11)

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		// 1 of resourceA is worth 2 credits
		// 2 credits is worth 1 of resourceA
		rate := decimal.NewFromInt(2)
		balanceStore, err := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": rate})
		require.NoError(t, err)

		assert.True(t, balanceStore.Get().Equal(ten), "initialization should set balance")

		bal, err := balanceStore.GetAs("resourceA")
		require.NoError(t, err)
		assert.True(t, bal.Equal(five), "rate should apply to balance")

		require.NoError(t, balanceStore.Add(one))
		assert.True(t, balanceStore.Get().Equal(eleven), "addition should update the balance")

		require.NoError(t, balanceStore.Minus(two))
		require.ErrorIs(t, balanceStore.Minus(decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(nine), "subtraction should update the balance")

		require.NoError(t, balanceStore.AddAs("resourceA", one))
		require.ErrorIs(t, balanceStore.AddAs("resourceA", decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(eleven), "addition by rate should update balance")

		require.NoError(t, balanceStore.MinusAs("resourceA", two))
		require.ErrorIs(t, balanceStore.MinusAs("resourceA", decimal.NewFromInt(-1)), ErrInvalidAmount)
		assert.True(t, balanceStore.Get().Equal(seven), "subtraction by rate should update balance")

		require.ErrorIs(t, balanceStore.AddAs("unknown", one), ErrResourceTypeNotFound)
	})

	t.Run("returns error for unknown resource", func(t *testing.T) {
		t.Parallel()

		balanceStore, err := NewBalanceStore(ten, map[string]decimal.Decimal{})
		require.NoError(t, err)

		bal, err := balanceStore.GetAs("")
		require.ErrorIs(t, err, ErrResourceTypeNotFound)
		assert.True(t, bal.Equal(ten))
		require.ErrorIs(t, balanceStore.MinusAs("", one), ErrResourceTypeNotFound)
		assert.True(t, balanceStore.Get().Equal(ten))
	})

	t.Run("errors when given negative conversion rates", func(t *testing.T) {
		t.Parallel()

		_, err := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(-1)})
		require.ErrorContains(t, err, "conversion rate -1 must be a positive number for resource resourceA")
	})

	t.Run("cannot go negative", func(t *testing.T) {
		t.Parallel()

		balanceStore, err := NewBalanceStore(decimal.Zero, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(1)})
		require.NoError(t, err)

		require.ErrorIs(t, balanceStore.Minus(one), ErrInsufficientBalance)
		require.ErrorIs(t, balanceStore.MinusAs("resourceA", one), ErrInsufficientBalance)
	})

	t.Run("handles decimal rates", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.1 credits
		rate := decimal.NewFromFloat(0.1)
		balanceStore, err := NewBalanceStore(ten, map[string]decimal.Decimal{"resourceA": rate})
		require.NoError(t, err)

		assert.True(t, balanceStore.Get().Equal(ten))

		bal, err := balanceStore.GetAs("resourceA")
		require.NoError(t, err)
		assert.True(t, bal.Equal(decimal.NewFromInt(100)))
	})

	t.Run("applies no rounding to result", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.2 credits
		rate := decimal.NewFromFloat(0.2)
		balanceStore, err := NewBalanceStore(two, map[string]decimal.Decimal{"resourceA": rate})
		require.NoError(t, err)

		assert.True(t, balanceStore.Get().Equal(two))
		require.NoError(t, balanceStore.MinusAs("resourceA", one))
		assert.True(t, balanceStore.Get().Equal(decimal.NewFromFloat(1.8)), balanceStore.Get())
	})
}
