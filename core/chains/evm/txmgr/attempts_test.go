package txmgr_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
)

func TestTxm_NewDynamicFeeTx(t *testing.T) {
	addr := testutils.NewAddress()
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst := ksmocks.NewEth(t)
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	var n int64
	lggr := logger.TestLogger(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		gcfg := configtest.NewGeneralConfig(t, nil)
		cfg := evmtest.NewChainScopedConfig(t, gcfg)
		cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)
		dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
		a, _, err := cks.NewAttemptWithType(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x2, lggr)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.Nil(t, a.GasPrice)
		assert.NotNil(t, a.GasTipCap)
		assert.Equal(t, assets.GWei(100).String(), a.GasTipCap.String())
		assert.NotNil(t, a.GasFeeCap)
		assert.Equal(t, assets.GWei(200).String(), a.GasFeeCap.String())
	})

	t.Run("verifies gas tip and fees", func(t *testing.T) {
		tests := []struct {
			name        string
			tipcap      *assets.Wei
			feecap      *assets.Wei
			setCfg      func(*chainlink.Config, *chainlink.Secrets)
			expectError string
		}{
			{"gas tip = fee cap", assets.GWei(5), assets.GWei(5), nil, ""},
			{"gas tip < fee cap", assets.GWei(4), assets.GWei(5), nil, ""},
			{"gas tip > fee cap", assets.GWei(6), assets.GWei(5), nil, "gas fee cap must be greater than or equal to gas tip cap (fee cap: 5 gwei, tip cap: 6 gwei)"},
			{"fee cap exceeds max allowed", assets.GWei(5), assets.GWei(5), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMax = (*assets.Wei)(assets.GWei(4))
			}, "specified gas fee cap of 5 gwei would exceed max configured gas price of 4 gwei"},
			{"ignores global min gas price", assets.GWei(5), assets.GWei(5), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMin = (*assets.Wei)(assets.GWei(6))
			}, ""},
			{"tip cap below min allowed", assets.GWei(5), assets.GWei(5), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.TipCapMin = (*assets.Wei)(assets.GWei(6))
			}, "specified gas tip cap of 5 gwei is below min configured gas tip of 6 gwei"},
		}

		for _, tt := range tests {
			test := tt
			t.Run(test.name, func(t *testing.T) {
				gcfg := configtest.NewGeneralConfig(t, test.setCfg)
				cfg := evmtest.NewChainScopedConfig(t, gcfg)
				cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)
				dynamicFee := gas.DynamicFee{TipCap: test.tipcap, FeeCap: test.feecap}
				_, _, err := cks.NewAttemptWithType(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x2, lggr)
				if test.expectError == "" {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, test.expectError)
				}
			})
		}
	})
}

func TestTxm_NewLegacyAttempt(t *testing.T) {
	addr := testutils.NewAddress()
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(50)
		c.EVM[0].GasEstimator.PriceMin = assets.NewWeiI(10)
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := ksmocks.NewEth(t)
	tx := types.NewTx(&types.LegacyTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)
	lggr := logger.TestLogger(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, _, err := cks.NewAttemptWithType(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(25)}, 100, 0x0, lggr)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.NotNil(t, a.GasPrice)
		assert.Equal(t, "25 wei", a.GasPrice.String())
		assert.Nil(t, a.GasTipCap)
		assert.Nil(t, a.GasFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, _, err := cks.NewAttemptWithType(txmgr.EthTx{FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(100)}, 100, 0x0, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 wei would exceed max configured gas price of 50 wei for key %s", addr.Hex()))
	})
}

func TestTxm_NewAttempt_NonRetryableErrors(t *testing.T) {
	t.Parallel()

	cfg := mocks.NewConfig(t)
	kst := ksmocks.NewEth(t)
	lggr := logger.TestLogger(t)
	cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)

	dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
	legacyFee := assets.NewWeiI(100)

	t.Run("NewAttempt: mismatch fee + type", func(t *testing.T) {
		t.Run("dynamic fee with legacy fee configured chain", func(t *testing.T) {
			cfg.On("EvmEIP1559DynamicFees").Return(false).Once()
			_, err := cks.NewAttempt(txmgr.EthTx{}, gas.EvmFee{Dynamic: &dynamicFee}, 100, lggr)
			require.Error(t, err)
		})
		t.Run("legacy fee with dynamic fee configured chain", func(t *testing.T) {
			cfg.On("EvmEIP1559DynamicFees").Return(true).Once()
			_, err := cks.NewAttempt(txmgr.EthTx{}, gas.EvmFee{Legacy: legacyFee}, 100, lggr)
			require.Error(t, err)
		})
	})

	t.Run("NewAttemptWithType: mismatch fee + type", func(t *testing.T) {
		t.Run("dynamic fee with legacy tx type", func(t *testing.T) {
			_, retryable, err := cks.NewAttemptWithType(txmgr.EthTx{}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x0, lggr)
			require.Error(t, err)
			assert.False(t, retryable)
		})
		t.Run("legacy fee with dynamic tx type", func(t *testing.T) {
			_, retryable, err := cks.NewAttemptWithType(txmgr.EthTx{}, gas.EvmFee{Legacy: legacyFee}, 100, 0x2, lggr)
			require.Error(t, err)
			assert.False(t, retryable)
		})
	})

	t.Run("NewAttemptWithType: invalid type", func(t *testing.T) {
		_, retryable, err := cks.NewAttemptWithType(txmgr.EthTx{}, gas.EvmFee{}, 100, 0xA, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})
}
