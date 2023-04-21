package txmgr_test

import (
	"fmt"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	txmgrmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/types/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	txmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/services/keystore/mocks"
)

func NewEvmAddress() gethcommon.Address {
	return testutils.NewAddress()
}

func TestTxm_SignTx(t *testing.T) {
	t.Parallel()

	addr := gethcommon.HexToAddress("0xb921F7763960b296B9cbAD586ff066A18D749724")
	to := gethcommon.HexToAddress("0xb921F7763960b296B9cbAD586ff066A18D749724")
	tx := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    42,
		To:       &to,
		Value:    big.NewInt(142),
		Gas:      242,
		GasPrice: big.NewInt(342),
		Data:     []byte{1, 2, 3},
	})

	t.Run("returns correct hash for non-okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		cfg := txmmocks.NewConfig(t)
		kst := ksmocks.NewEth(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, cfg, kst, nil)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.String())
	})
	// okex used to have a custom hash but now this just verifies that is it the same
	t.Run("returns correct hash for okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		cfg := txmmocks.NewConfig(t)
		kst := ksmocks.NewEth(t)
		kst.On("SignTx", to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, cfg, kst, nil)
		hash, rawBytes, err := cks.SignTx(addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.String())
	})
}

func TestTxm_NewDynamicFeeTx(t *testing.T) {
	addr := NewEvmAddress()
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst := ksmocks.NewEth(t)
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	var n int64
	lggr := logger.TestLogger(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		gcfg := configtest.NewGeneralConfig(t, nil)
		cfg := evmtest.NewChainScopedConfig(t, gcfg)
		cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg, kst, nil)
		dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
		a, _, err := cks.NewCustomTxAttempt(txmgr.EvmTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x2, lggr)
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
				cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg, kst, nil)
				dynamicFee := gas.DynamicFee{TipCap: test.tipcap, FeeCap: test.feecap}
				_, _, err := cks.NewCustomTxAttempt(txmgr.EvmTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x2, lggr)
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
	addr := NewEvmAddress()
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(50)
		c.EVM[0].GasEstimator.PriceMin = assets.NewWeiI(10)
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)
	kst := ksmocks.NewEth(t)
	tx := types.NewTx(&types.LegacyTx{})
	kst.On("SignTx", addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg, kst, nil)
	lggr := logger.TestLogger(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n int64
		a, _, err := cks.NewCustomTxAttempt(txmgr.EvmTx{Nonce: &n, FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(25)}, 100, 0x0, lggr)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificGasLimit))
		assert.NotNil(t, a.GasPrice)
		assert.Equal(t, "25 wei", a.GasPrice.String())
		assert.Nil(t, a.GasTipCap)
		assert.Nil(t, a.GasFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, _, err := cks.NewCustomTxAttempt(txmgr.EvmTx{FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(100)}, 100, 0x0, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 wei would exceed max configured gas price of 50 wei for key %s", addr.String()))
	})
}

func TestTxm_NewCustomTxAttempt_NonRetryableErrors(t *testing.T) {
	t.Parallel()

	cfg := txmmocks.NewConfig(t)
	kst := ksmocks.NewEth(t)
	lggr := logger.TestLogger(t)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg, kst, nil)

	dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
	legacyFee := assets.NewWeiI(100)

	t.Run("dynamic fee with legacy tx type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(txmgr.EvmTx{}, gas.EvmFee{Dynamic: &dynamicFee}, 100, 0x0, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})
	t.Run("legacy fee with dynamic tx type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(txmgr.EvmTx{}, gas.EvmFee{Legacy: legacyFee}, 100, 0x2, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(txmgr.EvmTx{}, gas.EvmFee{}, 100, 0xA, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})
}

func TestTxm_EvmTxAttemptBuilder_RetryableEstimatorError(t *testing.T) {
	est := txmgrmocks.NewFeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, gethcommon.Hash](t)
	est.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{}, uint32(0), errors.New("fail"))
	est.On("BumpFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{}, uint32(0), errors.New("fail"))

	cfg := txmmocks.NewConfig(t)
	cfg.On("EvmEIP1559DynamicFees").Return(true)
	cfg.On("KeySpecificMaxGasPriceWei", mock.Anything).Return(assets.NewWeiI(100))

	kst := ksmocks.NewEth(t)
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg, kst, est)

	t.Run("NewAttempt", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewTxAttempt(ctx, txmgr.EvmTx{}, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get fee")
		assert.True(t, retryable)
	})
	t.Run("NewAttemptWithType", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewTxAttemptWithType(ctx, txmgr.EvmTx{}, lggr, 0x0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get fee")
		assert.True(t, retryable)
	})
	t.Run("NewBumpAttempt", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewBumpTxAttempt(ctx, txmgr.EvmTx{}, txmgr.EvmTxAttempt{}, nil, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to bump fee")
		assert.True(t, retryable)
	})
}
