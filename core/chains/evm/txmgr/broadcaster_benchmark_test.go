package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/configtest"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

// happy path
func BenchmarkEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.B) {
	db := testutils.NewSqlxDB(t)
	ctx := tests.Context(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

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

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	evmcfg := configtest.NewChainScopedConfig(t, nil)
	checkerFactory := &txmgr.CheckerFactory{Client: ethClient}
	lggr := logger.Test(t)
	nonceTracker := txmgr.NewNonceTracker(lggr, txStore, txmgr.NewEvmTxmClient(ethClient, nil))

	ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(0), nil)

	eb := NewTestEthBroadcaster(t, txStore, ethClient, ethKeyStore, dbListenerCfg, evmcfg.EVM(), checkerFactory, false, nonceTracker)

	ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(multinode.Successful, nil)

	// Insertion order deliberately reversed to test ordering
	require.NoError(t, txStore.InsertTx(ctx, &expensiveEthTx))

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		eb.ProcessUnstartedTxs(tests.Context(t), fromAddress)
	}
}
