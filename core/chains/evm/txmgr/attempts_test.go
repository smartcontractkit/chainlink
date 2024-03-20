package txmgr_test

import (
	"fmt"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func NewEvmAddress() gethcommon.Address {
	return testutils.NewAddress()
}

type feeConfig struct {
	eip1559DynamicFees bool
	tipCapMin          *assets.Wei
	priceMin           *assets.Wei
	priceMax           *assets.Wei
}

func newFeeConfig() *feeConfig {
	return &feeConfig{
		tipCapMin: assets.NewWeiI(0),
		priceMin:  assets.NewWeiI(0),
		priceMax:  assets.NewWeiI(0),
	}
}

func (g *feeConfig) EIP1559DynamicFees() bool                        { return g.eip1559DynamicFees }
func (g *feeConfig) TipCapMin() *assets.Wei                          { return g.tipCapMin }
func (g *feeConfig) PriceMin() *assets.Wei                           { return g.priceMin }
func (g *feeConfig) PriceMaxKey(addr gethcommon.Address) *assets.Wei { return g.priceMax }

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
		kst := ksmocks.NewEth(t)
		kst.On("SignTx", mock.Anything, to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, newFeeConfig(), kst, nil)
		hash, rawBytes, err := cks.SignTx(testutils.Context(t), addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.String())
	})
	// okex used to have a custom hash but now this just verifies that is it the same
	t.Run("returns correct hash for okex chains", func(t *testing.T) {
		chainID := big.NewInt(1)
		kst := ksmocks.NewEth(t)
		kst.On("SignTx", mock.Anything, to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, newFeeConfig(), kst, nil)
		hash, rawBytes, err := cks.SignTx(testutils.Context(t), addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xdd68f554373fdea7ec6713a6e437e7646465d553a6aa0b43233093366cc87ef0", hash.String())
	})
	t.Run("can properly encoded and decode raw transaction for LegacyTx", func(t *testing.T) {
		chainID := big.NewInt(1)
		kst := ksmocks.NewEth(t)
		kst.On("SignTx", mock.Anything, to, tx, chainID).Return(tx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, newFeeConfig(), kst, nil)

		_, rawBytes, err := cks.SignTx(testutils.Context(t), addr, tx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xe42a82015681f294b921f7763960b296b9cbad586ff066a18d749724818e83010203808080", hexutil.Encode(rawBytes))

		var decodedTx *gethtypes.Transaction
		decodedTx, err = txmgr.GetGethSignedTx(rawBytes)
		require.NoError(t, err)
		require.Equal(t, tx.Hash(), decodedTx.Hash())
	})
	t.Run("can properly encoded and decode raw transaction for DynamicFeeTx", func(t *testing.T) {
		chainID := big.NewInt(1)
		kst := ksmocks.NewEth(t)
		typedTx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
			Nonce: 42,
			To:    &to,
			Value: big.NewInt(142),
			Gas:   242,
			Data:  []byte{1, 2, 3},
		})
		kst.On("SignTx", mock.Anything, to, typedTx, chainID).Return(typedTx, nil).Once()
		cks := txmgr.NewEvmTxAttemptBuilder(*chainID, newFeeConfig(), kst, nil)
		_, rawBytes, err := cks.SignTx(testutils.Context(t), addr, typedTx)
		require.NoError(t, err)
		require.NotNil(t, rawBytes)
		require.Equal(t, "0xa702e5802a808081f294b921f7763960b296b9cbad586ff066a18d749724818e83010203c0808080", hexutil.Encode(rawBytes))

		var decodedTx *gethtypes.Transaction
		decodedTx, err = txmgr.GetGethSignedTx(rawBytes)
		require.NoError(t, err)
		require.Equal(t, typedTx.Hash(), decodedTx.Hash())
	})
}

func TestTxm_NewDynamicFeeTx(t *testing.T) {
	addr := NewEvmAddress()
	tx := types.NewTx(&types.DynamicFeeTx{})
	kst := ksmocks.NewEth(t)
	kst.On("SignTx", mock.Anything, addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	var n evmtypes.Nonce
	lggr := logger.Test(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		feeCfg := newFeeConfig()
		feeCfg.priceMax = assets.GWei(200)
		cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), feeCfg, kst, nil)
		dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
		a, _, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{Sequence: &n, FromAddress: addr}, gas.EvmFee{
			DynamicTipCap: dynamicFee.TipCap,
			DynamicFeeCap: dynamicFee.FeeCap,
		}, 100, 0x2, lggr)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificFeeLimit))
		assert.Nil(t, a.TxFee.Legacy)
		assert.NotNil(t, a.TxFee.DynamicTipCap)
		assert.Equal(t, assets.GWei(100).String(), a.TxFee.DynamicTipCap.String())
		assert.NotNil(t, a.TxFee.DynamicFeeCap)
		assert.Equal(t, assets.GWei(200).String(), a.TxFee.DynamicFeeCap.String())
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
				c.EVM[0].GasEstimator.PriceMax = assets.GWei(4)
			}, "specified gas fee cap of 5 gwei would exceed max configured gas price of 4 gwei"},
			{"ignores global min gas price", assets.GWei(5), assets.GWei(5), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.PriceMin = assets.GWei(6)
			}, ""},
			{"tip cap below min allowed", assets.GWei(5), assets.GWei(5), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.TipCapMin = assets.GWei(6)
			}, "specified gas tip cap of 5 gwei is below min configured gas tip of 6 gwei"},
		}

		for _, tt := range tests {
			test := tt
			t.Run(test.name, func(t *testing.T) {
				gcfg := configtest.NewGeneralConfig(t, test.setCfg)
				cfg := evmtest.NewChainScopedConfig(t, gcfg)
				cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), cfg.EVM().GasEstimator(), kst, nil)
				dynamicFee := gas.DynamicFee{TipCap: test.tipcap, FeeCap: test.feecap}
				_, _, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{Sequence: &n, FromAddress: addr}, gas.EvmFee{
					DynamicTipCap: dynamicFee.TipCap,
					DynamicFeeCap: dynamicFee.FeeCap,
				}, 100, 0x2, lggr)
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
	kst := ksmocks.NewEth(t)
	tx := types.NewTx(&types.LegacyTx{})
	kst.On("SignTx", mock.Anything, addr, mock.Anything, big.NewInt(1)).Return(tx, nil)
	gc := newFeeConfig()
	gc.priceMin = assets.NewWeiI(10)
	gc.priceMax = assets.NewWeiI(50)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), gc, kst, nil)
	lggr := logger.Test(t)

	t.Run("creates attempt with fields", func(t *testing.T) {
		var n evmtypes.Nonce
		a, _, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{Sequence: &n, FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(25)}, 100, 0x0, lggr)
		require.NoError(t, err)
		assert.Equal(t, 100, int(a.ChainSpecificFeeLimit))
		assert.NotNil(t, a.TxFee.Legacy)
		assert.Equal(t, "25 wei", a.TxFee.Legacy.String())
		assert.Nil(t, a.TxFee.DynamicTipCap)
		assert.Nil(t, a.TxFee.DynamicFeeCap)
	})

	t.Run("verifies max gas price", func(t *testing.T) {
		_, _, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{FromAddress: addr}, gas.EvmFee{Legacy: assets.NewWeiI(100)}, 100, 0x0, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("specified gas price of 100 wei would exceed max configured gas price of 50 wei for key %s", addr.String()))
	})
}

func TestTxm_NewCustomTxAttempt_NonRetryableErrors(t *testing.T) {
	t.Parallel()

	kst := ksmocks.NewEth(t)
	lggr := logger.Test(t)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), newFeeConfig(), kst, nil)

	dynamicFee := gas.DynamicFee{TipCap: assets.GWei(100), FeeCap: assets.GWei(200)}
	legacyFee := assets.NewWeiI(100)

	t.Run("dynamic fee with legacy tx type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{}, gas.EvmFee{
			DynamicTipCap: dynamicFee.TipCap,
			DynamicFeeCap: dynamicFee.FeeCap,
		}, 100, 0x0, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})
	t.Run("legacy fee with dynamic tx type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{}, gas.EvmFee{Legacy: legacyFee}, 100, 0x2, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, retryable, err := cks.NewCustomTxAttempt(testutils.Context(t), txmgr.Tx{}, gas.EvmFee{}, 100, 0xA, lggr)
		require.Error(t, err)
		assert.False(t, retryable)
	})
}

func TestTxm_EvmTxAttemptBuilder_RetryableEstimatorError(t *testing.T) {
	est := gasmocks.NewEvmFeeEstimator(t)
	est.On("GetFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{}, uint64(0), pkgerrors.New("fail"))
	est.On("BumpFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.EvmFee{}, uint64(0), pkgerrors.New("fail"))

	kst := ksmocks.NewEth(t)
	lggr := logger.Test(t)
	ctx := testutils.Context(t)
	cks := txmgr.NewEvmTxAttemptBuilder(*big.NewInt(1), &feeConfig{eip1559DynamicFees: true}, kst, est)

	t.Run("NewAttempt", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewTxAttempt(ctx, txmgr.Tx{}, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get fee")
		assert.True(t, retryable)
	})
	t.Run("NewAttemptWithType", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewTxAttemptWithType(ctx, txmgr.Tx{}, lggr, 0x0)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get fee")
		assert.True(t, retryable)
	})
	t.Run("NewBumpAttempt", func(t *testing.T) {
		_, _, _, retryable, err := cks.NewBumpTxAttempt(ctx, txmgr.Tx{}, txmgr.TxAttempt{}, nil, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to bump fee")
		assert.True(t, retryable)
	})
}
