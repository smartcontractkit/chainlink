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
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
)

func TestTxm_NewDynamicFeeTx(t *testing.T) {
	addr := testutils.NewAddress()
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst := ksmocks.NewEth(t)
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	var n int64

	t.Run("creates attempt with fields", func(t *testing.T) {
		gcfg := configtest.NewGeneralConfig(t, nil)
		cfg := evmtest.NewChainScopedConfig(t, gcfg)
		cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)
		a, err := cks.NewDynamicFeeAttempt(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}, 100)
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
				_, err := cks.NewDynamicFeeAttempt(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.DynamicFee{TipCap: test.tipcap, FeeCap: test.feecap}, 100)
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

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, err := cks.NewLegacyAttempt(txmgr.EthTx{Nonce: &n, FromAddress: addr}, assets.NewWeiI(25), 100)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.NotNil(t, a.GasPrice)
		assert.Equal(t, "25 wei", a.GasPrice.String())
		assert.Nil(t, a.GasTipCap)
		assert.Nil(t, a.GasFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, err := cks.NewLegacyAttempt(txmgr.EthTx{FromAddress: addr}, assets.NewWeiI(100), 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 wei would exceed max configured gas price of 50 wei for key %s", addr.Hex()))
	})
}
