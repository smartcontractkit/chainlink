package txmgr_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_UpdateTxForRebroadcast(t *testing.T) {
	t.Parallel()

	t.Run("delete all receipts for transaction", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
		persistentStore := cltest.NewTestTxStore(t, db)
		kst := cltest.NewKeyStore(t, db, dbcfg)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		lggr := logger.TestSugared(t)
		chainID := ethClient.ConfiguredChainID()
		ctx := testutils.Context(t)

		inMemoryStore, err := commontxmgr.NewInMemoryStore[
			*big.Int,
			common.Address, common.Hash, common.Hash,
			*evmtypes.Receipt,
			evmtypes.Nonce,
			evmgas.EvmFee,
		](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
		require.NoError(t, err)

		// Insert a transaction into persistent store
		inTx := mustInsertConfirmedEthTxWithReceipt(t, persistentStore, fromAddress, 777, 1)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		txAttempt := inTx.TxAttempts[0]
		err = inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, commontxmgr.TxUnconfirmed, actTx.State)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, actTx.TxAttempts[0].State)
		assert.Nil(t, actTx.TxAttempts[0].BroadcastBeforeBlockNum)
		assert.Equal(t, 0, len(actTx.TxAttempts[0].Receipts))
	})

	t.Run("error parity for in-memory vs persistent store", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
		persistentStore := cltest.NewTestTxStore(t, db)
		kst := cltest.NewKeyStore(t, db, dbcfg)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		lggr := logger.TestSugared(t)
		chainID := ethClient.ConfiguredChainID()
		ctx := testutils.Context(t)

		inMemoryStore, err := commontxmgr.NewInMemoryStore[
			*big.Int,
			common.Address, common.Hash, common.Hash,
			*evmtypes.Receipt,
			evmtypes.Nonce,
			evmgas.EvmFee,
		](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
		require.NoError(t, err)

		// Insert a transaction into persistent store
		inTx := mustInsertConfirmedEthTxWithReceipt(t, persistentStore, fromAddress, 777, 1)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		txAttempt := inTx.TxAttempts[0]

		t.Run("error when attempt is not in Broadcast state", func(t *testing.T) {
			txAttempt.State = txmgrtypes.TxAttemptInProgress
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Error(t, expErr)
			assert.Error(t, actErr)
			txAttempt.State = txmgrtypes.TxAttemptBroadcast
		})
		t.Run("error when transaction is not in confirmed state", func(t *testing.T) {
			inTx.State = commontxmgr.TxUnconfirmed
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Error(t, expErr)
			assert.Error(t, actErr)
			inTx.State = commontxmgr.TxConfirmed
		})
		t.Run("wrong fromAddress has no error", func(t *testing.T) {
			inTx.FromAddress = common.Address{}
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Equal(t, expErr, actErr)
			assert.Nil(t, actErr)
			inTx.FromAddress = fromAddress
		})

	})
}

// assertTxEqual asserts that two transactions are equal
func assertTxEqual(t *testing.T, exp, act evmtxmgr.Tx) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.IdempotencyKey, act.IdempotencyKey)
	assert.Equal(t, exp.Sequence, act.Sequence)
	assert.Equal(t, exp.FromAddress, act.FromAddress)
	assert.Equal(t, exp.ToAddress, act.ToAddress)
	assert.Equal(t, exp.EncodedPayload, act.EncodedPayload)
	assert.Equal(t, exp.Value, act.Value)
	assert.Equal(t, exp.FeeLimit, act.FeeLimit)
	assert.Equal(t, exp.Error, act.Error)
	assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	assert.Equal(t, exp.InitialBroadcastAt, act.InitialBroadcastAt)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.Meta, act.Meta)
	assert.Equal(t, exp.Subject, act.Subject)
	assert.Equal(t, exp.ChainID, act.ChainID)
	assert.Equal(t, exp.PipelineTaskRunID, act.PipelineTaskRunID)
	assert.Equal(t, exp.MinConfirmations, act.MinConfirmations)
	assert.Equal(t, exp.TransmitChecker, act.TransmitChecker)
	assert.Equal(t, exp.SignalCallback, act.SignalCallback)
	assert.Equal(t, exp.CallbackCompleted, act.CallbackCompleted)

	require.Equal(t, len(exp.TxAttempts), len(act.TxAttempts))
	for i := 0; i < len(exp.TxAttempts); i++ {
		assertTxAttemptEqual(t, exp.TxAttempts[i], act.TxAttempts[i])
	}
}

func assertTxAttemptEqual(t *testing.T, exp, act evmtxmgr.TxAttempt) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.TxID, act.TxID)
	assert.Equal(t, exp.Tx, act.Tx)
	assert.Equal(t, exp.TxFee, act.TxFee)
	assert.Equal(t, exp.ChainSpecificFeeLimit, act.ChainSpecificFeeLimit)
	assert.Equal(t, exp.SignedRawTx, act.SignedRawTx)
	assert.Equal(t, exp.Hash, act.Hash)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.BroadcastBeforeBlockNum, act.BroadcastBeforeBlockNum)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.TxType, act.TxType)

	require.Equal(t, len(exp.Receipts), len(act.Receipts))
	for i := 0; i < len(exp.Receipts); i++ {
		assertChainReceiptEqual(t, exp.Receipts[i], act.Receipts[i])
	}
}

func assertChainReceiptEqual(t *testing.T, exp, act evmtxmgr.ChainReceipt) {
	assert.Equal(t, exp.GetStatus(), act.GetStatus())
	assert.Equal(t, exp.GetTxHash(), act.GetTxHash())
	assert.Equal(t, exp.GetBlockNumber(), act.GetBlockNumber())
	assert.Equal(t, exp.IsZero(), act.IsZero())
	assert.Equal(t, exp.IsUnmined(), act.IsUnmined())
	assert.Equal(t, exp.GetFeeUsed(), act.GetFeeUsed())
	assert.Equal(t, exp.GetTransactionIndex(), act.GetTransactionIndex())
	assert.Equal(t, exp.GetBlockHash(), act.GetBlockHash())
}
