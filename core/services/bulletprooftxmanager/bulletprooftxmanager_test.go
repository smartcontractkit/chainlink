package bulletprooftxmanager_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/require"
)

func TestBulletproofTxManager_BumpGas(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name             string
		originalGasPrice *big.Int
		priceDefault     *big.Int
		bumpPercent      uint16
		bumpWei          *big.Int
		maxGasPriceWei   *big.Int
		expected         *big.Int
	}{
		{
			name:             "defaults",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.6e10"), // 36 GWei
		},
		{
			name:             "original + percentage wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      30,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.9e10"), // 39 GWei
		},
		{
			name:             "original + fixed wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("2e10"), // 20 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("8e9"),    // 0.8 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("3.8e10"), // 38 GWei
		},
		{
			name:             "default + percentage wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("4e10"), // 40 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("5e9"),    // 0.5 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("4.8e10"), // 48 GWei
		},
		{
			name:             "default + fixed wins",
			originalGasPrice: toBigInt("3e10"), // 30 GWei
			priceDefault:     toBigInt("4e10"), // 40 GWei
			bumpPercent:      20,
			bumpWei:          toBigInt("9e9"),    // 0.9 GWei
			maxGasPriceWei:   toBigInt("5e11"),   // 0.5 uEther
			expected:         toBigInt("4.9e10"), // 49 GWei
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := orm.NewConfig()
			config.Set("ETH_GAS_PRICE_DEFAULT", test.priceDefault)
			config.Set("ETH_GAS_BUMP_PERCENT", test.bumpPercent)
			config.Set("ETH_GAS_BUMP_WEI", test.bumpWei)
			config.Set("ETH_MAX_GAS_PRICE_WEI", test.maxGasPriceWei)
			actual, err := bulletprooftxmanager.BumpGas(config, test.originalGasPrice)
			require.NoError(t, err)
			if actual.Cmp(test.expected) != 0 {
				t.Fatalf("Expected %s but got %s", test.expected.String(), actual.String())
			}
		})
	}
}

func TestBulletproofTxManager_BumpGas_HitsMaxError(t *testing.T) {
	t.Parallel()
	config := orm.NewConfig()
	config.Set("ETH_GAS_BUMP_PERCENT", "50")
	config.Set("ETH_GAS_PRICE_DEFAULT", toBigInt("2e10")) // 20 GWei
	config.Set("ETH_GAS_BUMP_WEI", toBigInt("5e9"))       // 0.5 GWei
	config.Set("ETH_MAX_GAS_PRICE_WEI", toBigInt("4e10")) // 40 Gwei

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, err := bulletprooftxmanager.BumpGas(config, originalGasPrice)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 45000000000 would exceed configured max gas price of 40000000000 (original price was 30000000000)")
}

func TestBulletproofTxManager_BumpGas_NoBumpError(t *testing.T) {
	t.Parallel()
	config := orm.NewConfig()
	config.Set("ETH_GAS_BUMP_PERCENT", "0")
	config.Set("ETH_GAS_BUMP_WEI", "0")
	config.Set("ETH_MAX_GAS_PRICE_WEI", "40000000000")

	originalGasPrice := toBigInt("3e10") // 30 GWei
	_, err := bulletprooftxmanager.BumpGas(config, originalGasPrice)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 30000000000 is equal to original gas price of 30000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

	// Even if it's exactly the maximum
	originalGasPrice = toBigInt("4e10") // 40 GWei
	_, err = bulletprooftxmanager.BumpGas(config, originalGasPrice)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bumped gas price of 40000000000 is equal to original gas price of 40000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")
}

// Helpers

// toBigInt is used to convert scientific notation string to a *big.Int
func toBigInt(input string) *big.Int {
	flt, _, err := big.ParseFloat(input, 10, 0, big.ToNearestEven)
	if err != nil {
		panic(fmt.Sprintf("unable to parse '%s' into a big.Float: %v", input, err))
	}
	var i = new(big.Int)
	i, _ = flt.Int(i)
	return i
}

func TestBulletproofTxManager_SendEther_DoesNotSendToZero(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	from := utils.ZeroAddress
	to := utils.ZeroAddress
	value := assets.NewEth(1)

	_, err := bulletprooftxmanager.SendEther(store, from, to, *value)
	require.Error(t, err)
	require.EqualError(t, err, "cannot send ether to zero address")
}

func TestBulletproofTxManager_CheckOKToTransmit(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, store)

	ctx := context.Background()
	db := store.MustSQLDB()
	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, int64(i), otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	var n int64 = 0
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithAttempt(t, store, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, n, fromAddress)
		n++
	}

	t.Run("with fewer unconfirmed eth_txes than limit returns nil", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, n, fromAddress)
	n++

	t.Run("with equal unconfirmed eth_txes to limit returns nil", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.NoError(t, err)
	})

	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, n, fromAddress)

	t.Run("with more unconfirmed eth_txes than limit returns error", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, maxUnconfirmedTransactions)
		require.Error(t, err)

		require.EqualError(t, err, fmt.Sprintf("cannot transmit eth transaction; there are currently %v unconfirmed transactions in the queue which exceeds the configured maximum of %v", maxUnconfirmedTransactions+1, maxUnconfirmedTransactions))
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := utils.CheckOKToTransmit(ctx, db, fromAddress, 0)
		require.NoError(t, err)
	})
}
