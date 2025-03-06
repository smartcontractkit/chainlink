package txmgr_test

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func BenchmarkCreateTransactionTxStore(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	txStore := newTxStore(b, db)
	kst := cltest.NewKeyStore(b, db)

	_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	ethClient := clienttest.NewClientWithDefaultChainID(b)

	// queue under capacity inserts eth_tx
	subject := uuid.New()
	strategy := newMockTxStrategy(b)
	strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txCount := n+1
		b.StartTimer()
		etx, err := txStore.CreateTransaction(tests.Context(b), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, ethClient.ConfiguredChainID())
		b.StopTimer()

		assert.NoError(b, err)
		assert.Greater(b, etx.ID, int64(0))
		assert.Equal(b, etx.State, txmgrcommon.TxUnstarted)
		assert.Equal(b, gasLimit, etx.FeeLimit)
		assert.Equal(b, fromAddress, etx.FromAddress)
		assert.Equal(b, toAddress, etx.ToAddress)
		assert.Equal(b, payload, etx.EncodedPayload)
		assert.Equal(b, big.Int(assets.NewEthValue(0)), etx.Value)
		assert.Equal(b, subject, etx.Subject.UUID)

		cltest.AssertCount(b, db, "evm.txes", int64(txCount))

		var dbEthTx txmgr.DbEthTx
		require.NoError(b, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(b, dbEthTx.State, txmgrcommon.TxUnstarted)
		assert.Equal(b, gasLimit, dbEthTx.GasLimit)
		assert.Equal(b, fromAddress, dbEthTx.FromAddress)
		assert.Equal(b, toAddress, dbEthTx.ToAddress)
		assert.Equal(b, payload, dbEthTx.EncodedPayload)
		assert.Equal(b, assets.NewEthValue(0), dbEthTx.Value)
		assert.Equal(b, subject, dbEthTx.Subject.UUID)
	}
}

func NewTestTxStore(t testing.TB, db *sqlx.DB) txmgr.TestEvmTxStore {
	t.Helper()
	return txmgr.NewTxStore(db, logger.Test(t))
}

func BenchmarkFindAttemptsRequiringReceiptFetch(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	txStore := NewTestTxStore(b, db)
	ctx := tests.Context(b)

	blockNum := int64(100)

	kst := cltest.NewKeyStore(b, db)
	_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())

	// Transactions whose attempts should not be picked up for receipt fetch
	mustInsertFatalErrorEthTx(b, txStore, fromAddress)
	mustInsertUnstartedTx(b, txStore, fromAddress)
	mustInsertInProgressEthTxWithAttempt(b, txStore, 4, fromAddress)
	mustInsertUnconfirmedEthTxWithAttemptState(b, txStore, 3, fromAddress, txmgrtypes.TxAttemptBroadcast)
	mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, 2, blockNum)
	// Terminally stuck transaction with receipt should NOT be picked up for receipt fetch
	stuckTx := mustInsertTerminallyStuckTxWithAttempt(b, txStore, fromAddress, 1, blockNum)
	mustInsertEthReceipt(b, txStore, blockNum, utils.NewHash(), stuckTx.TxAttempts[0].Hash)
	// Fatal transactions with nil nonce and stored attempts should NOT be picked up for receipt fetch
	fatalTxWithAttempt := mustInsertFatalErrorEthTx(b, txStore, fromAddress)
	attempt := newBroadcastLegacyEthTxAttempt(b, fatalTxWithAttempt.ID)
	err := txStore.InsertTxAttempt(ctx, &attempt)
	require.NoError(b, err)

	// Confirmed transaction without receipt should be picked up for receipt fetch
	confirmedTx := mustInsertConfirmedEthTx(b, txStore, 0, fromAddress)
	attempt = newBroadcastLegacyEthTxAttempt(b, confirmedTx.ID)
	err = txStore.InsertTxAttempt(ctx, &attempt)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		attempts, err := txStore.FindAttemptsRequiringReceiptFetch(ctx, testutils.FixtureChainID)
		b.StopTimer()

		require.NoError(b, err)
		require.Len(b, attempts, 1)
		require.Equal(b, attempt.Hash.String(), attempts[0].Hash.String())
	}
}

func BenchmarkFindTxesByIDs(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	txStore := NewTestTxStore(b, db)
	ctx := tests.Context(b)
	ethKeyStore := cltest.NewKeyStore(b, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(b, ethKeyStore)

	etx1 := mustInsertInProgressEthTxWithAttempt(b, txStore, 3, fromAddress)
	etx2 := mustInsertUnconfirmedEthTxWithAttemptState(b, txStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
	etx3 := mustInsertTerminallyStuckTxWithAttempt(b, txStore, fromAddress, 1, 100)
	etx4 := mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, 0, 100)
	etxIDs := []int64{etx1.ID, etx2.ID, etx3.ID, etx4.ID}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		oldTxs, err := txStore.FindTxesByIDs(ctx, etxIDs, testutils.FixtureChainID)
		b.StopTimer()

		require.NoError(b, err)
		require.Len(b, oldTxs, 4)
	}
}

func BenchmarkFindConfirmedTxesReceipts(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	txStore := NewTestTxStore(b, db)
	finalizedBlockNum := int64(100)
	kst := cltest.NewKeyStore(b, db)

	_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())
	for i := 0; i < 100; i++ {
		mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, int64(i), finalizedBlockNum)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		receipts, err := txStore.FindConfirmedTxesReceipts(tests.Context(b), finalizedBlockNum, &cltest.FixtureChainID)
		b.StopTimer()
		require.NoError(b, err)
		require.Len(b, receipts, 100)
	}
}
