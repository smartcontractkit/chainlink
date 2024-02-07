package txm

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
)

func TestResender_resendUnconfirmed(t *testing.T) {
	lggr := logger.Test(t)
	chainID := big.NewInt(0)
	blockTime := 2 * time.Second

	cfg, db := heavyweight.FullTestDBV2(t, nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	rcfg := ResenderConfig{
		BumpAfterThreshold:  3 * blockTime,
		MaxBumpCycles:       1,
		MaxInFlight:         evmcfg.EVM().Transactions().MaxInFlight(),
		ResendInterval:      blockTime,
		RPCDefaultBatchSize: evmcfg.EVM().RPCDefaultBatchSize(),
	}

	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	client := evmtest.NewEthClientMockWithDefaultChain(t)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	estimator := gasmocks.NewEvmFeeEstimator(t)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*chainID, evmcfg.EVM().GasEstimator(), keyStore, estimator)

	r := NewResender(txBuilder, lggr, txStore, client, keyStore, rcfg)

	ctx := testutils.Context(t)
	t.Run("no enabled addresses", func(t *testing.T) {
		require.NoError(t, r.resendUnconfirmed())
	})

	t.Run("no txs at all for enabled address", func(t *testing.T) {
		key1, addr1 := cltest.MustInsertRandomKey(t, keyStore)
		client.On("SequenceAt", mock.Anything, addr1, mock.Anything).Return(evmtypes.Nonce(0), nil).Once()

		require.NoError(t, r.resendUnconfirmed())
		keyStore.Delete(key1.ID())
	})

	t.Run("marks all unconfirmed txs as confirmed before the on-chain mined nonce", func(t *testing.T) {
		key2, addr2 := cltest.MustInsertRandomKey(t, keyStore)

		nonce0 := evmtypes.Nonce(0)
		encodedPayload := []byte{1, 2, 3}
		value := big.Int(assets.NewEthValue(142))
		gasLimit := uint32(242)
		timeNow := time.Now()
		txUnconfirmed1 := txmgr.Tx{
			Sequence:           &nonce0,
			FromAddress:        addr2,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              TxUnconfirmed,
		}

		nonce1 := evmtypes.Nonce(1)
		txUnconfirmed2 := txmgr.Tx{
			Sequence:           &nonce1,
			FromAddress:        addr2,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              TxUnconfirmed,
		}

		require.NoError(t, txStore.InsertTx(&txUnconfirmed1))
		require.NoError(t, txStore.InsertTx(&txUnconfirmed2))

		client.On("SequenceAt", mock.Anything, addr2, mock.Anything).Return(evmtypes.Nonce(1), nil).Once()
		require.NoError(t, r.resendUnconfirmed())

		n, err := txStore.CountUnconfirmedTransactions(ctx, addr2, chainID)
		require.NoError(t, err)
		require.Equal(t, uint32(1), n)

		tx, err := txStore.GetTxByID(ctx, txUnconfirmed1.ID)
		require.NoError(t, err)
		require.Equal(t, TxConfirmed, tx.State)

		keyStore.Delete(key2.ID())
	})

	t.Run("batch sends transactions that require gas bumping", func(t *testing.T) {
		_, addr3 := cltest.MustInsertRandomKey(t, keyStore)

		nonce0 := evmtypes.Nonce(0)
		encodedPayload := []byte{1, 2, 3}
		value := big.Int(assets.NewEthValue(142))
		gasLimit := uint32(242)
		timeNow := time.Now().Add(-time.Hour)

		txConfirmed := txmgr.Tx{
			Sequence:           &nonce0,
			FromAddress:        addr3,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              TxConfirmed,
		}

		nonce1 := evmtypes.Nonce(1)
		txUnconfirmed := txmgr.Tx{
			Sequence:           &nonce1,
			FromAddress:        addr3,
			ToAddress:          utils.RandomAddress(),
			EncodedPayload:     encodedPayload,
			Value:              value,
			FeeLimit:           gasLimit,
			BroadcastAt:        &timeNow,
			InitialBroadcastAt: &timeNow,
			Error:              null.String{},
			State:              TxUnconfirmed,
		}

		require.NoError(t, txStore.InsertTx(&txConfirmed))
		require.NoError(t, txStore.InsertTx(&txUnconfirmed))

		client.On("SequenceAt", mock.Anything, addr3, mock.Anything).Return(evmtypes.Nonce(1), nil).Once()
		estimator.On("GetFee", mock.Anything, txConfirmed.EncodedPayload, txConfirmed.FeeLimit, mock.Anything).
			Return(gas.EvmFee{Legacy: assets.GWei(32)}, uint32(500), nil).Once()
		estimator.On("GetFee", mock.Anything, txUnconfirmed.EncodedPayload, txUnconfirmed.FeeLimit, mock.Anything).
			Return(gas.EvmFee{Legacy: assets.GWei(35)}, uint32(500), nil).Once()

		estimator.On("BumpFee", mock.Anything, mock.Anything, txConfirmed.FeeLimit, mock.Anything, mock.Anything).
			Return(gas.EvmFee{Legacy: assets.GWei(42)}, uint32(500), nil).Once()
		estimator.On("BumpFee", mock.Anything, mock.Anything, txUnconfirmed.FeeLimit, mock.Anything, mock.Anything).
			Return(gas.EvmFee{Legacy: assets.GWei(45)}, uint32(500), nil).Once()

		client.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(elems []rpc.BatchElem) bool {
			assert.Len(t, elems, 2)
			return true
		})).Run(func(args mock.Arguments) {}).Return(nil)
		require.NoError(t, r.resendUnconfirmed())

		tx1, err := txStore.GetTxByID(ctx, txConfirmed.ID)
		require.NoError(t, err)
		require.Equal(t, TxUnconfirmed, tx1.State)
		require.NotEqual(t, &timeNow, tx1.BroadcastAt) // BroadcastAt should be updated if a new attempt was sent

		tx2, err := txStore.GetTxByID(ctx, txUnconfirmed.ID)
		require.NoError(t, err)
		require.Equal(t, TxUnconfirmed, tx2.State)
		require.NotEqual(t, &timeNow, tx2.BroadcastAt)// BroadcastAt should be updated if a new attempt was sent

	})

}
