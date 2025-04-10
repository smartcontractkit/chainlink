package txmgr_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

// happy path
func BenchmarkEthBroadcaster_ProcessUnstartedEthTxs_Success(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	ctx := tests.Context(b)

	ethClient := clienttest.NewClientWithDefaultChainID(b)
	txStore := cltest.NewTestTxStore(b, db)
	memKeystore := keystest.NewMemoryChainStore()
	ethKeyStore := keys.NewChainStore(memKeystore, ethClient.ConfiguredChainID())
	fromAddress := memKeystore.MustCreate(b)

	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")

	gasLimit := uint64(242)

	// Higher value
	expensiveEthTx := txmgr.Tx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: []byte{42, 42, 0},
		Value:          big.Int(assets.NewEthValue(242)),
		FeeLimit:       gasLimit,
		CreatedAt:      time.Unix(0, 0),
		State:          txmgrcommon.TxUnstarted,
	}

	evmcfg := configtest.NewChainScopedConfig(b, nil)
	checkerFactory := &txmgr.CheckerFactory{Client: ethClient}
	lggr := logger.Test(b)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))

	ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(0), nil)

	eb := NewTestEthBroadcaster(b, txStore, ethClient, ethKeyStore, dbListenerCfg, evmcfg.EVM(), checkerFactory, false, nonceTracker)

	ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(multinode.Successful, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Insertion order deliberately reversed to test ordering
		require.NoError(b, txStore.InsertTx(ctx, &expensiveEthTx))
		b.StartTimer()

		eb.ProcessUnstartedTxs(tests.Context(b), fromAddress)
		b.StopTimer()
		deleteTx(ctx, b, &expensiveEthTx, db)
		b.StartTimer()
	}
}

func deleteTx(ctx context.Context, b *testing.B, etx *txmgr.Tx, db *sqlx.DB) {
	var dbTx txmgr.DbEthTx
	dbTx.FromTx(etx)
	txID := dbTx.ID
	require.NotNil(b, txID)
	deleteTxQuery := fmt.Sprintf("DELETE FROM evm.txes WHERE id = %d", dbTx.ID)
	_ = deleteTxQuery
	db.Exec(deleteTxQuery)
}
