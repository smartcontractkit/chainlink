package metering

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestBalanceStore(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		// 1 of resourceA is worth 2 credits
		// 2 credits is worth 1 of resourceA
		rate := decimal.NewFromInt(2)
		balanceStore := NewBalanceStore(10, map[string]decimal.Decimal{"resourceA": rate}, logger.TestLogger(t))

		balance := balanceStore.Get()
		require.Equal(t, int64(10), balance)

		balanceAs := balanceStore.GetAs("resourceA")
		require.Equal(t, int64(5), balanceAs)

		err := balanceStore.Add(1)
		require.NoError(t, err)
		balance = balanceStore.Get()
		require.Equal(t, int64(11), balance)

		err = balanceStore.Minus(2)
		require.NoError(t, err)
		balance = balanceStore.Get()
		require.Equal(t, int64(9), balance)

		err = balanceStore.AddAs("resourceA", 1)
		require.NoError(t, err)
		balance = balanceStore.Get()
		require.Equal(t, int64(11), balance)

		err = balanceStore.MinusAs("resourceA", 2)
		require.NoError(t, err)
		balance = balanceStore.Get()
		require.Equal(t, int64(7), balance)
	})

	t.Run("handles unknown resources as 1:1", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(10, map[string]decimal.Decimal{}, logger.TestLogger(t))

		balanceAs := balanceStore.GetAs("")
		require.Equal(t, int64(10), balanceAs)

		err := balanceStore.MinusAs("", 1)
		require.NoError(t, err)
		balance := balanceStore.Get()
		require.Equal(t, int64(9), balance)
	})

	t.Run("throws out negative conversion rates", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(10, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(-1)}, logger.TestLogger(t))

		balanceAs := balanceStore.GetAs("resourceA")
		require.Equal(t, int64(10), balanceAs)

		err := balanceStore.MinusAs("resourceA", 1)
		require.NoError(t, err)
		balance := balanceStore.Get()
		require.Equal(t, int64(9), balance)
	})

	t.Run("cannot go negative by default", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(1)}, logger.TestLogger(t))

		err := balanceStore.Minus(1)
		require.ErrorIs(t, ErrInsufficientBalance, err)
	})

	t.Run("can go negative after allowNegative", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(1)}, logger.TestLogger(t))

		balanceStore.AllowNegative()
		err := balanceStore.Minus(1)
		require.NoError(t, err)
		balance := balanceStore.Get()
		require.Equal(t, int64(-1), balance)
	})

	t.Run("returns negative balances as 0 when converted to a resource", func(t *testing.T) {
		t.Parallel()

		balanceStore := NewBalanceStore(0, map[string]decimal.Decimal{"resourceA": decimal.NewFromInt(10)}, logger.TestLogger(t))

		balanceStore.AllowNegative()
		err := balanceStore.Minus(1)
		require.NoError(t, err)
		balance := balanceStore.Get()
		require.Equal(t, int64(-1), balance)

		balanceAs := balanceStore.GetAs("resourceA")
		require.Equal(t, int64(0), balanceAs)
	})

	t.Run("handles decimal rates", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.1 credits
		rate, err := decimal.NewFromString("0.1")
		require.NoError(t, err)
		balanceStore := NewBalanceStore(10, map[string]decimal.Decimal{"resourceA": rate}, logger.TestLogger(t))

		balance := balanceStore.Get()
		require.Equal(t, int64(10), balance)

		balanceAs := balanceStore.GetAs("resourceA")
		require.Equal(t, int64(100), balanceAs)
	})

	t.Run("rounds up decimals due to conversion to nearest integer", func(t *testing.T) {
		t.Parallel()

		// 1 of resource A is worth 0.2 credits
		rate, err := decimal.NewFromString("0.2")
		require.NoError(t, err)
		balanceStore := NewBalanceStore(2, map[string]decimal.Decimal{"resourceA": rate}, logger.TestLogger(t))

		balance := balanceStore.Get()
		require.Equal(t, int64(2), balance)

		err = balanceStore.MinusAs("resourceA", 1)
		require.NoError(t, err)
		balance = balanceStore.Get()
		require.Equal(t, int64(1), balance)
	})
}
