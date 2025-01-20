package txm

import (
	"math/big"
	"testing"

	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

func TestAttemptBuilder_newLegacyAttempt(t *testing.T) {
	ks := mocks.NewEth(t)
	ab := NewAttemptBuilder(testutils.FixtureChainID, nil, nil, ks)
	address := testutils.NewAddress()
	toAddress := testutils.NewAddress()
	lggr := logger.Test(t)
	var gasLimit uint64 = 100

	t.Run("fails if GasPrice is nil", func(t *testing.T) {
		tx := &types.Transaction{ID: 10, FromAddress: address}
		_, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{DynamicFee: gas.DynamicFee{GasTipCap: assets.NewWeiI(1), GasFeeCap: assets.NewWeiI(2)}}, gasLimit, evmtypes.LegacyTxType, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "estimator did not return legacy fee")
	})

	t.Run("fails if tx doesn't have a nonce", func(t *testing.T) {
		tx := &types.Transaction{ID: 10, FromAddress: address}
		_, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{GasPrice: assets.NewWeiI(25)}, gasLimit, evmtypes.LegacyTxType, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonce empty")
	})

	t.Run("creates attempt with fields", func(t *testing.T) {
		var nonce uint64 = 77
		tx := &types.Transaction{ID: 10, FromAddress: address, Nonce: &nonce}
		legacyTx := evmtypes.NewTx(&evmtypes.LegacyTx{Nonce: nonce, To: &toAddress, Gas: gasLimit, GasPrice: big.NewInt(25)})
		ks.On("SignTx", mock.Anything, mock.Anything, mock.Anything, testutils.FixtureChainID).Return(legacyTx, nil).Once()
		a, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{GasPrice: assets.NewWeiI(25)}, gasLimit, evmtypes.LegacyTxType, lggr)
		require.NoError(t, err)
		assert.Equal(t, tx.ID, a.TxID)
		assert.Equal(t, evmtypes.LegacyTxType, int(a.Type))
		assert.NotNil(t, a.Fee.GasPrice)
		assert.Equal(t, "25 wei", a.Fee.GasPrice.String())
		assert.Nil(t, a.Fee.GasTipCap)
		assert.Nil(t, a.Fee.GasFeeCap)
		assert.Equal(t, gasLimit, a.GasLimit)
	})
}

func TestAttemptBuilder_newDynamicFeeAttempt(t *testing.T) {
	ks := mocks.NewEth(t)
	ab := NewAttemptBuilder(testutils.FixtureChainID, nil, nil, ks)
	address := testutils.NewAddress()
	toAddress := testutils.NewAddress()
	lggr := logger.Test(t)
	var gasLimit uint64 = 100

	t.Run("fails if DynamicFee is invalid", func(t *testing.T) {
		tx := &types.Transaction{ID: 10, FromAddress: address}
		_, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{GasPrice: assets.NewWeiI(1)}, gasLimit, evmtypes.DynamicFeeTxType, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "estimator did not return dynamic fee")
	})

	t.Run("fails if tx doesn't have a nonce", func(t *testing.T) {
		tx := &types.Transaction{ID: 10, FromAddress: address}
		_, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{DynamicFee: gas.DynamicFee{GasTipCap: assets.NewWeiI(1), GasFeeCap: assets.NewWeiI(2)}}, gasLimit, evmtypes.DynamicFeeTxType, lggr)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonce empty")
	})

	t.Run("creates attempt with fields", func(t *testing.T) {
		var nonce uint64 = 77
		tx := &types.Transaction{ID: 10, FromAddress: address, Nonce: &nonce}
		legacyTx := evmtypes.NewTx(&evmtypes.LegacyTx{Nonce: nonce, To: &toAddress, Gas: gasLimit, GasPrice: big.NewInt(25)})
		ks.On("SignTx", mock.Anything, mock.Anything, mock.Anything, testutils.FixtureChainID).Return(legacyTx, nil).Once()
		a, err := ab.newCustomAttempt(tests.Context(t), tx, gas.EvmFee{DynamicFee: gas.DynamicFee{GasTipCap: assets.NewWeiI(1), GasFeeCap: assets.NewWeiI(2)}}, gasLimit, evmtypes.DynamicFeeTxType, lggr)
		require.NoError(t, err)
		assert.Equal(t, tx.ID, a.TxID)
		assert.Equal(t, evmtypes.DynamicFeeTxType, int(a.Type))
		assert.Equal(t, "1 wei", a.Fee.DynamicFee.GasTipCap.String())
		assert.Equal(t, "2 wei", a.Fee.DynamicFee.GasFeeCap.String())
		assert.Nil(t, a.Fee.GasPrice)
		assert.Equal(t, gasLimit, a.GasLimit)
	})
}
