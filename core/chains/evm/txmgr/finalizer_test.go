package txmgr_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func TestFinalizer_MarkTxFinalized(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	feeLimit := uint64(10_000)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	rpcBatchSize := uint32(1)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	h99 := &types.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	h99.IsFinalized.Store(true)
	head := &types.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(h99)

	t.Run("returns not finalized for tx with receipt newer than finalized block", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := types.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		attemptHash := insertTxAndAttemptWithIdempotencyKey(t, txStore, tx, idempotencyKey)
		// Insert receipt for unfinalized block num
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attemptHash)
		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(head.Parent.Load(), nil).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, tx.State)
	})

	t.Run("returns not finalized for tx with receipt re-org'd out and deletes stale receipt", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := types.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		attemptHash := insertTxAndAttemptWithIdempotencyKey(t, txStore, tx, idempotencyKey)
		// Insert receipt for finalized block num
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, utils.NewHash(), attemptHash)
		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(head.Parent.Load(), nil).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, tx.State)
		require.Len(t, tx.TxAttempts, 1)
		require.Empty(t, tx.TxAttempts[0].Receipts)
	})

	t.Run("returns finalized for tx with receipt in a finalized block", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := types.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		attemptHash := insertTxAndAttemptWithIdempotencyKey(t, txStore, tx, idempotencyKey)
		// Insert receipt for finalized block num
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, head.Parent.Load().Hash, attemptHash)
		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(head.Parent.Load(), nil).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFinalized, tx.State)
	})

	t.Run("returns finalized for tx with receipt older than block history depth", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := types.Nonce(0)
		broadcast := time.Now()
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		attemptHash := insertTxAndAttemptWithIdempotencyKey(t, txStore, tx, idempotencyKey)
		// Insert receipt for finalized block num
		receiptBlockHash1 := utils.NewHash()
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number-2, receiptBlockHash1, attemptHash)
		idempotencyKey = uuid.New().String()
		nonce = types.Nonce(1)
		tx = &txmgr.Tx{
			Sequence:           &nonce,
			IdempotencyKey:     &idempotencyKey,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			FeeLimit:           feeLimit,
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		attemptHash = insertTxAndAttemptWithIdempotencyKey(t, txStore, tx, idempotencyKey)
		// Insert receipt for finalized block num
		receiptBlockHash2 := utils.NewHash()
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number-1, receiptBlockHash2, attemptHash)
		// Separate batch calls will be made for each tx due to RPC batch size set to 1 when finalizer initialized above
		ethClient.On("BatchCallContext", mock.Anything, mock.IsType([]rpc.BatchElem{})).Run(func(args mock.Arguments) {
			rpcElements := args.Get(1).([]rpc.BatchElem)
			require.Len(t, rpcElements, 1)

			require.Equal(t, "eth_getBlockByNumber", rpcElements[0].Method)
			require.Equal(t, false, rpcElements[0].Args[1])

			reqBlockNum := rpcElements[0].Args[0].(string)
			req1BlockNum := hexutil.EncodeBig(big.NewInt(head.Parent.Load().Number - 2))
			req2BlockNum := hexutil.EncodeBig(big.NewInt(head.Parent.Load().Number - 1))
			var headResult types.Head
			if req1BlockNum == reqBlockNum {
				headResult = types.Head{Number: head.Parent.Load().Number - 2, Hash: receiptBlockHash1}
			} else if req2BlockNum == reqBlockNum {
				headResult = types.Head{Number: head.Parent.Load().Number - 1, Hash: receiptBlockHash2}
			} else {
				require.Fail(t, "unrecognized block hash")
			}
			rpcElements[0].Result = &headResult
		}).Return(nil).Twice()
		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(head.Parent.Load(), nil).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.NoError(t, err)
		tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFinalized, tx.State)
	})

	t.Run("returns error if failed to retrieve latest head in headtracker", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(nil, errors.New("failed to get latest head")).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.Error(t, err)
	})

	t.Run("returns error if failed to calculate latest finalized head in headtracker", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		servicetest.Run(t, finalizer)

		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(nil, errors.New("failed to calculate latest finalized head")).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.Error(t, err)
	})
}

func insertTxAndAttemptWithIdempotencyKey(tb testing.TB, txStore txmgr.TestEvmTxStore, tx *txmgr.Tx, idempotencyKey string) common.Hash {
	ctx := tests.Context(tb)
	err := txStore.InsertTx(ctx, tx)
	require.NoError(tb, err)
	tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
	require.NoError(tb, err)
	attempt := cltest.NewLegacyEthTxAttempt(tb, tx.ID)
	err = txStore.InsertTxAttempt(ctx, &attempt)
	require.NoError(tb, err)
	return attempt.Hash
}

func TestFinalizer_ResumePendingRuns(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	rpcBatchSize := uint32(1)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	grandParentHead := &types.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	parentHead := &types.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	parentHead.Parent.Store(grandParentHead)
	head := types.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}
	head.Parent.Store(parentHead)

	minConfirmations := int64(2)

	testutils.MustExec(t, db, `SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	testutils.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)

	t.Run("doesn't process task runs that are not suspended (possibly already previously resumed)", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(context.Context, uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})
		servicetest.Run(t, finalizer)

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)
		// Setting both signal_callback and callback_completed to TRUE to simulate a completed pipeline task
		// It would only be in a state past suspended if the resume callback was called and callback_completed was set to TRUE
		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE, callback_completed = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := finalizer.ResumePendingTaskRuns(ctx, head.BlockNumber(), 0)
		require.NoError(t, err)
	})

	t.Run("doesn't process task runs where the receipt is younger than minConfirmations", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(context.Context, uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})
		servicetest.Run(t, finalizer)

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 2, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, etx.TxAttempts[0].Hash)

		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := finalizer.ResumePendingTaskRuns(ctx, head.BlockNumber(), 0)
		require.NoError(t, err)
	})

	t.Run("processes transactions with receipts older than minConfirmations", func(t *testing.T) {
		ch := make(chan interface{})
		nonce := types.Nonce(3)
		var err error
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(ctx context.Context, id uuid.UUID, value interface{}, thisErr error) error {
			err = thisErr
			ch <- value
			return nil
		})
		servicetest.Run(t, finalizer)

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
		testutils.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		testutils.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)
		receipt := mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)

		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		done := make(chan struct{})
		t.Cleanup(func() { <-done })
		go func() {
			defer close(done)
			err2 := finalizer.ResumePendingTaskRuns(ctx, head.BlockNumber(), 0)
			assert.NoError(t, err2)

			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(ctx, fromAddress, nonce)
			assert.NoError(t, err3)
			assert.True(t, updateTx.CallbackCompleted)
		}()

		select {
		case data := <-ch:
			require.NoError(t, err)

			require.IsType(t, &types.Receipt{}, data)
			r := data.(*types.Receipt)
			require.Equal(t, receipt.TxHash, r.TxHash)

		case <-time.After(time.Second):
			t.Fatal("no value received")
		}
	})

	testutils.MustExec(t, db, `DELETE FROM pipeline_runs`)

	t.Run("processes transactions with receipt older than minConfirmations that reverted", func(t *testing.T) {
		type data struct {
			value any
			error
		}
		ch := make(chan data)
		nonce := types.Nonce(4)
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(ctx context.Context, id uuid.UUID, value interface{}, err error) error {
			ch <- data{value, err}
			return nil
		})
		servicetest.Run(t, finalizer)

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
		testutils.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		testutils.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)

		// receipt is not passed through as a value since it reverted and caused an error
		mustInsertRevertedEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)

		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		done := make(chan struct{})
		t.Cleanup(func() { <-done })
		go func() {
			defer close(done)
			err2 := finalizer.ResumePendingTaskRuns(ctx, head.BlockNumber(), 0)
			assert.NoError(t, err2)

			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(ctx, fromAddress, nonce)
			assert.NoError(t, err3)
			assert.True(t, updateTx.CallbackCompleted)
		}()

		select {
		case data := <-ch:
			require.Error(t, data.error)

			require.EqualError(t, data.error, fmt.Sprintf("transaction %s reverted on-chain", etx.TxAttempts[0].Hash.String()))

			require.Nil(t, data.value)

		case <-time.After(tests.WaitTimeout(t)):
			t.Fatal("no value received")
		}
	})

	t.Run("does not mark callback complete if callback fails", func(t *testing.T) {
		nonce := types.Nonce(5)
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(ctx context.Context, id uuid.UUID, value interface{}, err error) error {
			return errors.New("error")
		})
		servicetest.Run(t, finalizer)

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)
		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := finalizer.ResumePendingTaskRuns(ctx, head.BlockNumber(), 0)
		require.Error(t, err)

		// Retrieve Tx to check if callback completed flag was left unchanged
		updateTx, err := txStore.FindTxWithSequence(ctx, fromAddress, nonce)
		require.NoError(t, err)
		require.False(t, updateTx.CallbackCompleted)
	})
}

func TestFinalizer_FetchAndStoreReceipts(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	config := configtest.NewChainScopedConfig(t, nil)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	rpcBatchSize := config.EVM().RPCDefaultBatchSize()
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	latestFinalizedHead := &types.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	latestFinalizedHead.IsFinalized.Store(true)
	head := &types.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(latestFinalizedHead)

	t.Run("does nothing if no confirmed transactions without receipts found", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, config.EVM().RPCDefaultBatchSize(), false, txStore, txmClient, ht, metrics)

		mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		mustInsertInProgressEthTx(t, txStore, 0, fromAddress)
		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 2, fromAddress)
		mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())
		// Insert confirmed transactions with receipt and multiple attempts to ensure none of the attempts are picked up
		etx := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 3, head.Number)
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))
	})

	t.Run("fetches receipt for confirmed transaction without a receipt", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		// Transaction not confirmed yet, receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &types.Receipt{}
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		require.NoError(t, err)
		require.Empty(t, attempt.Receipts)
	})

	t.Run("saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		txmReceipt := types.Receipt{
			TxHash:           testutils.NewHash(),
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = txmReceipt
		}).Once()

		// No error because it is merely logged
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		require.Empty(t, etx.TxAttempts[0].Receipts)
	})

	t.Run("saves nothing if query returns error", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		txmReceipt := types.Receipt{
			TxHash:           etx.TxAttempts[0].Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		// Batch receipt call fails
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = txmReceipt
			elems[0].Error = errors.New("foo")
		}).Once()

		// No error because it is merely logged
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		require.Empty(t, etx.TxAttempts[0].Receipts)
	})

	t.Run("saves valid receipt returned by client", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		// Insert confirmed transaction without receipt
		etx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, head.Number, fromAddress)
		txmReceipt := types.Receipt{
			TxHash:           etx1.TxAttempts[0].Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], etx1.TxAttempts[0].Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], etx2.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			*(elems[0].Result.(*types.Receipt)) = txmReceipt
			// Second transaction still unconfirmed
			elems[1].Result = &types.Receipt{}
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Check that the receipt was saved
		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)

		require.Equal(t, txmgrcommon.TxConfirmed, etx1.State)
		require.Len(t, etx1.TxAttempts, 1)
		attempt := etx1.TxAttempts[0]
		require.Len(t, attempt.Receipts, 1)
		receipt := attempt.Receipts[0]
		require.Equal(t, txmReceipt.TxHash, receipt.GetTxHash())
		require.Equal(t, txmReceipt.BlockHash, receipt.GetBlockHash())
		require.Equal(t, txmReceipt.BlockNumber.Int64(), receipt.GetBlockNumber().Int64())
		require.Equal(t, txmReceipt.TransactionIndex, receipt.GetTransactionIndex())

		receiptJSON, err := json.Marshal(txmReceipt)
		require.NoError(t, err)

		storedReceiptJSON, err := json.Marshal(receipt)
		require.NoError(t, err)
		require.JSONEq(t, string(receiptJSON), string(storedReceiptJSON))
	})

	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		attempt1 := etx.TxAttempts[0]
		attempt2 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		attempt3 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 3)

		// Insert order deliberately reversed to test sorting by gas price
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2))

		txmReceipt := types.Receipt{
			TxHash:           attempt2.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				cltest.BatchElemMatchesParams(b[2], attempt1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[0], attempt3.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Most expensive attempt still unconfirmed
			elems[2].Result = &types.Receipt{}
			// Second most expensive attempt is confirmed
			*(elems[1].Result.(*types.Receipt)) = txmReceipt
			// Cheapest attempt still unconfirmed
			elems[0].Result = &types.Receipt{}
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Check that the receipt was stored
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 3)
		require.Empty(t, etx.TxAttempts[0].Receipts)
		require.Len(t, etx.TxAttempts[1].Receipts, 1)
		require.Empty(t, etx.TxAttempts[2].Receipts)
	})

	t.Run("ignores receipt missing BlockHash that comes from querying parity too early", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		receipt := types.Receipt{
			TxHash: etx.TxAttempts[0].Hash,
			Status: uint64(1),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = receipt
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// No receipt, but no error either
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		require.Empty(t, attempt.Receipts)
	})

	t.Run("does not panic if receipt has BlockHash but is missing some other fields somehow", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		// NOTE: This should never happen, but we shouldn't panic regardless
		receipt := types.Receipt{
			TxHash:    etx.TxAttempts[0].Hash,
			BlockHash: testutils.NewHash(),
			Status:    uint64(1),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = receipt
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// No receipt, but no error either
		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]
		require.Empty(t, attempt.Receipts)
	})

	t.Run("simulate on revert", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)
		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		attempt := etx.TxAttempts[0]
		txmReceipt := types.Receipt{
			TxHash:           attempt.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(0),
		}

		// First attempt is confirmed and reverted
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First attempt still unconfirmed
			*(elems[0].Result.(*types.Receipt)) = txmReceipt
		}).Once()
		data, err := utils.ABIEncode(`[{"type":"uint256"}]`, big.NewInt(10))
		require.NoError(t, err)
		sig := utils.Keccak256Fixed([]byte(`MyError(uint256)`))
		ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, &client.JsonError{
			Code:    1,
			Message: "reverted",
			Data:    utils.ConcatBytes(sig[:4], data),
		}).Once()

		// Do the thing
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Check that the state was updated
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		attempt = etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		require.NotNil(t, attempt.BroadcastBeforeBlockNum)
		// Check receipts
		require.Len(t, attempt.Receipts, 1)
	})

	t.Run("find receipt for old transaction, avoid marking as fatal", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, true, txStore, txmClient, ht, metrics)

		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, latestFinalizedHead.Number, fromAddress)

		txmReceipt := types.Receipt{
			TxHash:           etx.TxAttempts[0].Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		// Transaction receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = txmReceipt
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Check that transaction was picked up as old and marked as fatal
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
	})

	t.Run("old transaction failed to find receipt, marked as fatal", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, true, txStore, txmClient, ht, metrics)

		// Insert confirmed transaction without receipt
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, latestFinalizedHead.Number, fromAddress)

		// Transaction receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], etx.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &types.Receipt{}
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Check that transaction was picked up as old and marked as fatal
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx.State)
		require.Equal(t, txmgr.ErrCouldNotGetReceipt, etx.Error.String)
	})

	t.Run("attempts requiring receipt fetch is not fetched from TxStore every head", func(t *testing.T) {
		txStore := mocks.NewEvmTxStore(t)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)

		// Mock finalizer txstore calls that are not needed
		txStore.On("SaveFetchedReceipts", mock.Anything, mock.Anything).Return(nil).Maybe()
		txStore.On("FindTxesPendingCallback", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
		txStore.On("UpdateTxCallbackCompleted", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
		txStore.On("FindConfirmedTxesReceipts", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()
		txStore.On("FindTxesByIDs", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()

		// RPC returns nil receipt for attempt
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Maybe()

		// Should fetch attempts list from txstore
		attempt := cltest.NewLegacyEthTxAttempt(t, 0)
		txStore.On("FindAttemptsRequiringReceiptFetch", mock.Anything, mock.Anything).Return([]txmgr.TxAttempt{attempt}, nil).Once()
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))
		// Should use the attempts cache for receipt fetch
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))
		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))
	})
}

func TestFinalizer_FetchAndStoreReceipts_batching(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	latestFinalizedHead := &types.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	latestFinalizedHead.IsFinalized.Store(true)
	head := &types.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(latestFinalizedHead)

	t.Run("fetch and store receipts from multiple batch calls", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		rpcBatchSize := uint32(2)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)

		// Insert confirmed transaction without receipt
		etx := mustInsertConfirmedEthTx(t, txStore, 0, fromAddress)

		var attempts []txmgr.TxAttempt
		// Total of 5 attempts should lead to 3 batched fetches (2, 2, 1)v
		for i := 0; i < 5; i++ {
			attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(i+2))
			attempt.BroadcastBeforeBlockNum = &head.Number
			require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
			attempts = append(attempts, attempt)
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempts[4].Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempts[3].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &types.Receipt{}
			elems[1].Result = &types.Receipt{}
		}).Once()
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempts[2].Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempts[1].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &types.Receipt{}
			elems[1].Result = &types.Receipt{}
		}).Once()
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				cltest.BatchElemMatchesParams(b[0], attempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &types.Receipt{}
		}).Once()

		require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))
	})

	t.Run("continue to fetch and store receipts after batch call error", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		rpcBatchSize := uint32(1)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht, metrics)

		// Insert confirmed transactions without receipts
		etx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, head.Number, fromAddress)
		etx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, head.Number, fromAddress)

		txmReceipt := types.Receipt{
			TxHash:           etx2.TxAttempts[0].Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				cltest.BatchElemMatchesParams(b[0], etx1.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(errors.New("batch call failed")).Once()
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				cltest.BatchElemMatchesParams(b[0], etx2.TxAttempts[0].Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*types.Receipt)) = txmReceipt // confirmed
		}).Once()

		// Returns error due to batch call failure
		require.Error(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

		// Still fetches and stores receipt for later batch call that succeeds
		var err error
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Len(t, etx2.TxAttempts, 1)
		attempt := etx2.TxAttempts[0]
		require.Len(t, attempt.Receipts, 1)
	})
}

func TestFinalizer_FetchAndStoreReceipts_HandlesNonFwdTxsWithForwardingEnabled(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	latestFinalizedHead := &types.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	latestFinalizedHead.IsFinalized.Store(true)
	head := &types.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(latestFinalizedHead)

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, 1, true, txStore, txmClient, ht, metrics)

	// tx is not forwarded and doesn't have meta set. Confirmer should handle nil meta values
	etx := mustInsertConfirmedEthTx(t, txStore, 0, fromAddress)
	attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
	attempt.Tx.Meta = nil
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	dbtx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	require.Empty(t, dbtx.TxAttempts[0].Receipts)

	txmReceipt := types.Receipt{
		TxHash:           attempt.Hash,
		BlockHash:        testutils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
		Status:           uint64(1),
	}

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesParams(b[0], attempt.Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		*(elems[0].Result.(*types.Receipt)) = txmReceipt // confirmed
	}).Once()

	require.NoError(t, finalizer.FetchAndStoreReceipts(ctx, head, latestFinalizedHead))

	// Check receipt is inserted correctly.
	dbtx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	require.Len(t, dbtx.TxAttempts[0].Receipts, 1)
}

func TestFinalizer_ProcessOldTxsWithoutReceipts(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)

	latestFinalizedHead := &types.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	latestFinalizedHead.IsFinalized.Store(true)
	head := &types.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(latestFinalizedHead)

	t.Run("does nothing if no old transactions found", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, 1, true, txStore, txmClient, ht, metrics)
		require.NoError(t, finalizer.ProcessOldTxsWithoutReceipts(ctx, []int64{}, head, latestFinalizedHead))
	})

	t.Run("marks multiple old transactions as fatal", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, 1, true, txStore, txmClient, ht, metrics)

		// Insert confirmed transaction without receipt
		etx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, latestFinalizedHead.Number, fromAddress)
		etx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, latestFinalizedHead.Number, fromAddress)

		etxIDs := []int64{etx1.ID, etx2.ID}
		require.NoError(t, finalizer.ProcessOldTxsWithoutReceipts(ctx, etxIDs, head, latestFinalizedHead))

		// Check transactions marked as fatal
		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx1.State)
		require.Equal(t, txmgr.ErrCouldNotGetReceipt, etx1.Error.String)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx2.State)
		require.Equal(t, txmgr.ErrCouldNotGetReceipt, etx2.Error.String)
	})

	t.Run("marks old transaction as fatal, resumes pending task as failed", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, 1, true, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(context.Context, uuid.UUID, interface{}, error) error { return nil })

		// Insert confirmed transaction with pending task run
		etx := cltest.NewEthTx(fromAddress)
		etx.State = txmgrcommon.TxConfirmed
		n := types.Nonce(0)
		etx.Sequence = &n
		now := time.Now()
		etx.BroadcastAt = &now
		etx.InitialBroadcastAt = &now
		etx.SignalCallback = true
		etx.PipelineTaskRunID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		require.NoError(t, txStore.InsertTx(t.Context(), &etx))

		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 0)
		attempt.BroadcastBeforeBlockNum = &latestFinalizedHead.Number // set broadcast time to finalized block num
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

		require.NoError(t, finalizer.ProcessOldTxsWithoutReceipts(ctx, []int64{etx.ID}, head, latestFinalizedHead))

		// Check transaction marked as fatal
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx.State)
		require.Equal(t, txmgr.ErrCouldNotGetReceipt, etx.Error.String)
		require.True(t, etx.CallbackCompleted)
	})

	t.Run("transaction stays confirmed if failure to resume pending task", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, 1, true, txStore, txmClient, ht, metrics)
		finalizer.SetResumeCallback(func(context.Context, uuid.UUID, interface{}, error) error { return errors.New("failure") })

		// Insert confirmed transaction with pending task run
		etx := cltest.NewEthTx(fromAddress)
		etx.State = txmgrcommon.TxConfirmed
		n := types.Nonce(0)
		etx.Sequence = &n
		now := time.Now()
		etx.BroadcastAt = &now
		etx.InitialBroadcastAt = &now
		etx.SignalCallback = true
		etx.PipelineTaskRunID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		require.NoError(t, txStore.InsertTx(t.Context(), &etx))

		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 0)
		attempt.BroadcastBeforeBlockNum = &latestFinalizedHead.Number // set broadcast time to finalized block num
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

		// Expect error since resuming pending task failed
		require.Error(t, finalizer.ProcessOldTxsWithoutReceipts(ctx, []int64{etx.ID}, head, latestFinalizedHead))

		// Check transaction marked as fatal
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.False(t, etx.CallbackCompleted)
	})
}
