package txmgr_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestFinalizer_MarkTxFinalized(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	feeLimit := uint64(10_000)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	rpcBatchSize := uint32(1)
	ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)

	h99 := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 99,
	}
	h99.IsFinalized.Store(true)
	head := &evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 100,
	}
	head.Parent.Store(h99)

	t.Run("returns not finalized for tx with receipt newer than finalized block", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := evmtypes.Nonce(0)
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

	t.Run("returns not finalized for tx with receipt re-org'd out", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := evmtypes.Nonce(0)
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
	})

	t.Run("returns finalized for tx with receipt in a finalized block", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := evmtypes.Nonce(0)
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
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		idempotencyKey := uuid.New().String()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
		nonce := evmtypes.Nonce(0)
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
		nonce = evmtypes.Nonce(1)
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
			require.Equal(t, 1, len(rpcElements))

			require.Equal(t, "eth_getBlockByNumber", rpcElements[0].Method)
			require.Equal(t, false, rpcElements[0].Args[1])

			reqBlockNum := rpcElements[0].Args[0].(string)
			req1BlockNum := hexutil.EncodeBig(big.NewInt(head.Parent.Load().Number - 2))
			req2BlockNum := hexutil.EncodeBig(big.NewInt(head.Parent.Load().Number - 1))
			var headResult evmtypes.Head
			if req1BlockNum == reqBlockNum {
				headResult = evmtypes.Head{Number: head.Parent.Load().Number - 2, Hash: receiptBlockHash1}
			} else if req2BlockNum == reqBlockNum {
				headResult = evmtypes.Head{Number: head.Parent.Load().Number - 1, Hash: receiptBlockHash2}
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
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(nil, errors.New("failed to get latest head")).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.Error(t, err)
	})

	t.Run("returns error if failed to calculate latest finalized head in headtracker", func(t *testing.T) {
		finalizer := txmgr.NewEvmFinalizer(logger.Test(t), testutils.FixtureChainID, rpcBatchSize, txStore, ethClient, ht)
		servicetest.Run(t, finalizer)

		ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(nil, errors.New("failed to calculate latest finalized head")).Once()
		err := finalizer.ProcessHead(ctx, head)
		require.Error(t, err)
	})
}

func insertTxAndAttemptWithIdempotencyKey(t *testing.T, txStore txmgr.TestEvmTxStore, tx *txmgr.Tx, idempotencyKey string) common.Hash {
	ctx := tests.Context(t)
	err := txStore.InsertTx(ctx, tx)
	require.NoError(t, err)
	tx, err = txStore.FindTxWithIdempotencyKey(ctx, idempotencyKey, testutils.FixtureChainID)
	require.NoError(t, err)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx.ID)
	err = txStore.InsertTxAttempt(ctx, &attempt)
	require.NoError(t, err)
	return attempt.Hash
}
