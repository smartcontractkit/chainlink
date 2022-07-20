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
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
)

func TestTxm_NewDynamicFeeTx(t *testing.T) {
	addr := testutils.NewAddress()
	gcfg := cltest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := new(ksmocks.Eth)
	kst.Test(t)
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	var n int64

	t.Run("creates attempt with fields", func(t *testing.T) {
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
			tipcap      *big.Int
			feecap      *big.Int
			setCfg      func(cfg *configtest.TestGeneralConfig)
			expectError string
		}{
			{"gas tip = fee cap", assets.GWei(5), assets.GWei(5), nil, ""},
			{"gas tip < fee cap", assets.GWei(4), assets.GWei(5), nil, ""},
			{"gas tip > fee cap", assets.GWei(6), assets.GWei(5), nil, "gas fee cap must be greater than or equal to gas tip cap (fee cap: 5000000000, tip cap: 6000000000)"},
			{"fee cap exceeds max allowed", assets.GWei(5), assets.GWei(5), func(cfg *configtest.TestGeneralConfig) {
				cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(4)
			}, "specified gas fee cap of 5000000000 would exceed max configured gas price of 4000000000"},
			{"ignores global min gas price", assets.GWei(5), assets.GWei(5), func(cfg *configtest.TestGeneralConfig) {
				cfg.Overrides.GlobalEvmMinGasPriceWei = assets.GWei(6)
			}, ""},
			{"tip cap below min allowed", assets.GWei(5), assets.GWei(5), func(cfg *configtest.TestGeneralConfig) {
				cfg.Overrides.GlobalEvmGasTipCapMinimum = assets.GWei(6)
			}, "specified gas tip cap of 5000000000 is below min configured gas tip of 6000000000"},
		}

		for _, tt := range tests {
			test := tt
			t.Run(test.name, func(t *testing.T) {
				gcfg := configtest.NewTestGeneralConfig(t)
				if test.setCfg != nil {
					test.setCfg(gcfg)
				}
				cfg := evmtest.NewChainScopedConfig(t, gcfg)
				cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)
				_, err := cks.NewDynamicFeeAttempt(txmgr.EthTx{Nonce: &n, FromAddress: addr}, gas.DynamicFee{TipCap: test.tipcap, FeeCap: test.feecap}, 100)
				if test.expectError == "" {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					assert.Contains(t, err.Error(), test.expectError)
				}
			})
		}
	})
}

func TestTxm_NewLegacyAttempt(t *testing.T) {
	addr := testutils.NewAddress()
	gcfg := cltest.NewTestGeneralConfig(t)
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	gcfg.Overrides.GlobalEvmMaxGasPriceWei = big.NewInt(50)
	gcfg.Overrides.GlobalEvmMinGasPriceWei = big.NewInt(10)
	kst := new(ksmocks.Eth)
	kst.Test(t)
	tx := types.NewTx(&types.LegacyTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	cks := txmgr.NewChainKeyStore(*big.NewInt(1), cfg, kst)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, err := cks.NewLegacyAttempt(txmgr.EthTx{Nonce: &n, FromAddress: addr}, big.NewInt(25), 100)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.NotNil(t, a.GasPrice)
		assert.Equal(t, big.NewInt(25).String(), a.GasPrice.String())
		assert.Nil(t, a.GasTipCap)
		assert.Nil(t, a.GasFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, err := cks.NewLegacyAttempt(txmgr.EthTx{FromAddress: addr}, big.NewInt(100), 100)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 would exceed max configured gas price of 50 for key %s", addr.Hex()))
	})
}
