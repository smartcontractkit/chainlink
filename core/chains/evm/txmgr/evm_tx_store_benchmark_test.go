package txmgr_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func NewTestTxStore(t testing.TB, db *sqlx.DB) txmgr.TestEvmTxStore {
	t.Helper()
	return txmgr.NewTxStore(db, logger.Test(t))
}

var benchmarkSizes = []struct {
	name string
	size int
}{
	{"Factor_1", 1},
	{"Factor_10", 10},
	{"Factor_100", 100},
	{"Factor_1000", 1000},
}

func BenchmarkTxStoreCreateTransaction(b *testing.B) {
	for _, bs := range benchmarkSizes {
		b.Run(bs.name, func(b *testing.B) {
			db := testutils.NewSqlxDB(b)
			txStore := newTxStore(b, db)
			kst := cltest.NewKeyStore(b, db)
			_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())
			gasLimit := uint64(1000)
			payload := []byte{1, 2, 3}
			ethClient := clienttest.NewClientWithDefaultChainID(b)

			subject := uuid.New()
			strategy := newMockTxStrategy(b)
			strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})

			for i := 0; i < bs.size; i++ {
				toAddress := testutils.NewAddress()
				_, err := txStore.CreateTransaction(tests.Context(b), txmgr.TxRequest{
					FromAddress:    fromAddress,
					ToAddress:      toAddress,
					EncodedPayload: payload,
					FeeLimit:       gasLimit,
					Strategy:       strategy,
				}, ethClient.ConfiguredChainID())
				assert.NoError(b, err)
			}

			b.StopTimer()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				toAddress := testutils.NewAddress()
				b.StartTimer()
				_, err := txStore.CreateTransaction(tests.Context(b), txmgr.TxRequest{
					FromAddress:    fromAddress,
					ToAddress:      toAddress,
					EncodedPayload: payload,
					FeeLimit:       gasLimit,
					Strategy:       strategy,
				}, ethClient.ConfiguredChainID())
				b.StopTimer()
				assert.NoError(b, err)
			}
		})
	}
}

func BenchmarkTxStoreFindAttemptsRequiringReceiptFetch(b *testing.B) {
	for _, bs := range benchmarkSizes {
		b.Run(bs.name, func(b *testing.B) {
			db := testutils.NewSqlxDB(b)
			txStore := NewTestTxStore(b, db)
			ctx := tests.Context(b)
			blockNum := int64(100)
			kst := cltest.NewKeyStore(b, db)
			_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())

			var nonce = evmtypes.Nonce(0)
			for i := 0; i < bs.size; i++ {
				// Transactions whose attempts should not be picked up for receipt fetch
				mustInsertFatalErrorEthTx(b, txStore, fromAddress)
				mustInsertUnstartedTx(b, txStore, fromAddress)
				mustInsertUnconfirmedEthTxWithAttemptState(b, txStore, nonce.Int64(), fromAddress, txmgrtypes.TxAttemptBroadcast)
				nonce++
				mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, nonce.Int64(), blockNum)
				nonce++
				// Terminally stuck transaction with receipt should NOT be picked up for receipt fetch
				stuckTx := mustInsertTerminallyStuckTxWithAttempt(b, txStore, fromAddress, nonce.Int64(), blockNum)
				nonce++
				mustInsertEthReceipt(b, txStore, blockNum, utils.NewHash(), stuckTx.TxAttempts[0].Hash)
				// Fatal transactions with nil nonce and stored attempts should NOT be picked up for receipt fetch
				fatalTx := mustInsertFatalErrorEthTx(b, txStore, fromAddress)
				attempt := newBroadcastLegacyEthTxAttempt(b, fatalTx.ID)
				require.NoError(b, txStore.InsertTxAttempt(ctx, &attempt))
				// Confirmed transaction without receipt should be picked up for receipt fetch
				confirmedTx := mustInsertConfirmedEthTx(b, txStore, nonce.Int64(), fromAddress)
				nonce++
				attempt = newBroadcastLegacyEthTxAttempt(b, confirmedTx.ID)
				require.NoError(b, txStore.InsertTxAttempt(ctx, &attempt))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				attempts, err := txStore.FindAttemptsRequiringReceiptFetch(ctx, testutils.FixtureChainID)
				b.StopTimer()
				require.NoError(b, err)
				require.Equal(b, bs.size, len(attempts))
			}
		})
	}
}

func BenchmarkFindTxesByIDs(b *testing.B) {
	for _, bs := range benchmarkSizes {
		b.Run(bs.name, func(b *testing.B) {
			db := testutils.NewSqlxDB(b)
			txStore := NewTestTxStore(b, db)
			ctx := tests.Context(b)
			ethKeyStore := cltest.NewKeyStore(b, db).Eth()
			_, fromAddress := cltest.MustInsertRandomKeyReturningState(b, ethKeyStore)

			var etxIDs []int64
			for i := 0; i < bs.size; i++ {
				etx := mustInsertUnconfirmedEthTxWithAttemptState(b, txStore, int64(i), fromAddress, txmgrtypes.TxAttemptBroadcast)
				etxIDs = append(etxIDs, etx.ID)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				txs, err := txStore.FindTxesByIDs(ctx, etxIDs, testutils.FixtureChainID)
				b.StopTimer()
				require.NoError(b, err)
				require.Len(b, txs, len(etxIDs))
			}
		})
	}
}

func BenchmarkFindConfirmedTxesReceipts(b *testing.B) {
	for _, bs := range benchmarkSizes {
		b.Run(bs.name, func(b *testing.B) {
			db := testutils.NewSqlxDB(b)
			txStore := NewTestTxStore(b, db)
			finalizedBlockNum := int64(100)
			kst := cltest.NewKeyStore(b, db)
			_, fromAddress := cltest.MustInsertRandomKey(b, kst.Eth())

			for i := 0; i < bs.size; i++ {
				mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, int64(i), finalizedBlockNum)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				receipts, err := txStore.FindConfirmedTxesReceipts(tests.Context(b), finalizedBlockNum, &cltest.FixtureChainID)
				b.StopTimer()
				require.NoError(b, err)
				require.Len(b, receipts, bs.size)
			}
		})
	}
}
