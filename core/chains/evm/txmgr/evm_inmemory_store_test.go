package txmgr_test

import (
	"context"
	"math/big"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_PruneUnstartedTxQueue(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db, dbcfg)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := context.Background()

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("doesnt prune unstarted transactions if under maxQueueSize", func(t *testing.T) {
		maxQueueSize := uint32(5)
		nTxs := 3
		subject := uuid.NullUUID{UUID: uuid.New(), Valid: true}
		strat := commontxmgr.NewDropOldestStrategy(subject.UUID, maxQueueSize, dbcfg.DefaultQueryTimeout())
		for i := 0; i < nTxs; i++ {
			inTx := cltest.NewEthTx(fromAddress)
			inTx.Subject = subject
			// insert the transaction into the persistent store
			require.NoError(t, persistentStore.InsertTx(&inTx))
			// insert the transaction into the in-memory store
			require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
		}

		ids, err := strat.PruneQueue(ctx, inMemoryStore)
		require.NoError(t, err)
		assert.Equal(t, 0, len(ids))

		AssertCountPerSubject(t, persistentStore, int64(nTxs), subject.UUID)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		states := []txmgrtypes.TxState{commontxmgr.TxUnstarted}
		actTxs := inMemoryStore.XXXTestFindTxs(states, fn)
		expTxs, err := persistentStore.FindTxesByFromAddressAndState(ctx, fromAddress, "unstarted")
		require.NoError(t, err)
		require.Equal(t, len(expTxs), len(actTxs))

		// sort by ID to ensure the order is the same for comparison
		sort.SliceStable(actTxs, func(i, j int) bool {
			return actTxs[i].ID < actTxs[j].ID
		})
		sort.SliceStable(expTxs, func(i, j int) bool {
			return expTxs[i].ID < expTxs[j].ID
		})
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], actTxs[i])
		}
	})
	t.Run("prunes unstarted transactions", func(t *testing.T) {
		maxQueueSize := uint32(5)
		nTxs := 5
		subject := uuid.NullUUID{UUID: uuid.New(), Valid: true}
		strat := commontxmgr.NewDropOldestStrategy(subject.UUID, maxQueueSize, dbcfg.DefaultQueryTimeout())
		for i := 0; i < nTxs; i++ {
			inTx := cltest.NewEthTx(fromAddress)
			inTx.Subject = subject
			// insert the transaction into the persistent store
			require.NoError(t, persistentStore.InsertTx(&inTx))
			// insert the transaction into the in-memory store
			require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
		}

		ids, err := strat.PruneQueue(ctx, inMemoryStore)
		require.NoError(t, err)
		assert.Equal(t, int(nTxs)-int(maxQueueSize-1), len(ids))

		AssertCountPerSubject(t, persistentStore, int64(maxQueueSize-1), subject.UUID)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		states := []txmgrtypes.TxState{commontxmgr.TxUnstarted}
		actTxs := inMemoryStore.XXXTestFindTxs(states, fn)
		expTxs, err := persistentStore.FindTxesByFromAddressAndState(ctx, fromAddress, "unstarted")
		require.NoError(t, err)
		require.Equal(t, len(expTxs), len(actTxs))

		// sort by ID to ensure the order is the same for comparison
		sort.SliceStable(actTxs, func(i, j int) bool {
			return actTxs[i].ID < actTxs[j].ID
		})
		sort.SliceStable(expTxs, func(i, j int) bool {
			return expTxs[i].ID < expTxs[j].ID
		})
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], actTxs[i])
		}
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

	require.Len(t, exp.TxAttempts, len(act.TxAttempts))
	for i := 0; i < len(exp.TxAttempts); i++ {
		assertTxAttemptEqual(t, exp.TxAttempts[i], act.TxAttempts[i])
	}
}

func assertTxAttemptEqual(t *testing.T, exp, act evmtxmgr.TxAttempt) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.TxID, act.TxID)
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
