package bulletprooftxmanager_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/gas"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
)

func TestBulletproofTxManager_NewDynamicFeeTx(t *testing.T) {
	addr := cltest.NewAddress()
	gcfg := cltest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := new(ksmocks.Eth)
	kst.Test(t)
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, err := bulletprooftxmanager.NewDynamicFeeAttempt(cfg, kst, big.NewInt(1), bulletprooftxmanager.EthTx{Nonce: &n, FromAddress: addr}, gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}, 100)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.Nil(t, a.GasPrice)
		assert.NotNil(t, a.GasTipCap)
		assert.Equal(t, assets.GWei(100), a.GasTipCap)
		assert.NotNil(t, a.GasFeeCap)
		assert.Equal(t, assets.GWei(200), a.GasFeeCap)
	})
}

func TestBulletproofTxManager_NewLegacyAttempt(t *testing.T) {
	addr := cltest.NewAddress()
	gcfg := cltest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	gcfg.Overrides.GlobalEvmMaxGasPriceWei = big.NewInt(50)
	gcfg.Overrides.GlobalEvmMinGasPriceWei = big.NewInt(10)
	kst := new(ksmocks.Eth)
	kst.Test(t)
	tx := types.NewTx(&types.LegacyTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, err := bulletprooftxmanager.NewLegacyAttempt(cfg, kst, big.NewInt(1), bulletprooftxmanager.EthTx{Nonce: &n, FromAddress: addr}, big.NewInt(25), 100)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.NotNil(t, a.GasPrice)
		assert.Equal(t, big.NewInt(25), a.GasPrice)
		assert.Nil(t, a.GasTipCap)
		assert.Nil(t, a.GasFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, err := bulletprooftxmanager.NewLegacyAttempt(cfg, nil, big.NewInt(1), bulletprooftxmanager.EthTx{FromAddress: addr}, big.NewInt(100), 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 would exceed max configured gas price of 50 for key %s", addr.Hex()))
	})
}
