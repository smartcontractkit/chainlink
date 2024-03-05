package txmgr_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	//"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"

	//	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

// TODO: CLEANUP?

func TestInMemoryStore_FindTxWithIdempotencyKey(t *testing.T) {
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

	idempotencyKey := "777"
	inTx := cltest.NewEthTx(fromAddress)
	inTx.IdempotencyKey = &idempotencyKey
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(&inTx))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

	tcs := []struct {
		name             string
		inIdempotencyKey string
		inChainID        *big.Int

		hasErr bool
		hasTx  bool
	}{
		{"no idempotency key", "", chainID, true, true},
		{"wrong idempotency key", "wrong", chainID, true, true},
		{"finds tx with idempotency key", idempotencyKey, chainID, false, false},
		{"wrong chain", idempotencyKey, big.NewInt(999), true, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			actTx, actErr := inMemoryStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			expTx, expErr := persistentStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, actErr)
				require.Nil(t, expErr)
			}
			if tc.hasTx {
				require.NotNil(t, actTx)
				require.NotNil(t, expTx)
				assertTxEqual(t, *expTx, *actTx)
			} else {
				require.Nil(t, actTx)
				require.Nil(t, expTx)
			}
		})
	}
}

func TestInMemoryStore_CheckTxQueueCapacity(t *testing.T) {
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

	inTxs := []evmtxmgr.Tx{
		cltest.NewEthTx(fromAddress),
		cltest.NewEthTx(fromAddress),
	}
	for _, inTx := range inTxs {
		// insert the transaction into the persistent store
		require.NoError(t, persistentStore.InsertTx(&inTx))
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name           string
		inFromAddress  common.Address
		inMaxQueuedTxs uint64
		inChainID      *big.Int

		hasErr bool
	}{
		{"capacity reached", fromAddress, 2, chainID, true},
		{"above capacity", fromAddress, 1, chainID, true},
		{"below capacity", fromAddress, 3, chainID, false},
		{"wrong chain", fromAddress, 2, big.NewInt(999), false},
		{"wrong address", common.Address{}, 2, chainID, false},
		{"max queued txs is 0", fromAddress, 0, chainID, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			actErr := inMemoryStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			expErr := persistentStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
			//require.Equal(t, expErr, actErr)
		})
	}
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
