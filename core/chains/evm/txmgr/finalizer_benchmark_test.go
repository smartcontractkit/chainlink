package txmgr_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink-integrations/evm/heads/headstest"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func BenchmarkFinalizer(b *testing.B) {
	ctx := tests.Context(b)
	db := testutils.NewSqlxDB(b)
	txStore := cltest.NewTestTxStore(b, db)
	ethKeyStore := cltest.NewKeyStore(b, db).Eth()
	feeLimit := uint64(10_000)
	ethClient := clienttest.NewClientWithDefaultChainID(b)
	txmClient := txmgr.NewEvmTxmClient(ethClient, nil)
	rpcBatchSize := uint32(1)
	ht := headstest.NewSimulatedHeadTracker(ethClient, true, 0)

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
	finalizer := txmgr.NewEvmFinalizer(logger.Test(b), testutils.FixtureChainID, rpcBatchSize, false, txStore, txmClient, ht)
	servicetest.Run(b, finalizer)

	_, fromAddress := cltest.MustInsertRandomKey(b, ethKeyStore)
	
	broadcast := time.Now()
	ethClient.On("HeadByNumber", mock.Anything, mock.Anything).Return(head, nil)
	ethClient.On("LatestFinalizedBlock", mock.Anything).Return(head.Parent.Load(), nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		idempotencyKey := uuid.New().String()
		nonce := types.Nonce(i)
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
		attemptHash := insertTxAndAttemptWithIdempotencyKey(b, txStore, tx, idempotencyKey)
		// Insert receipt for finalized block num
		mustInsertEthReceipt(b, txStore, head.Parent.Load().Number, head.Parent.Load().Hash, attemptHash)
		b.StartTimer()
		err := finalizer.ProcessHead(ctx, head)
		b.StopTimer()
		require.NoError(b, err)
		deleteTx(ctx, b, tx, db)
	}
}