package txmgr_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink-integrations/evm/client/clienttest"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/configtest"
	"github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func BenchmarkEthConfirmer(t *testing.B) {
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	blockNum := int64(100)
	head := evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: blockNum,
	}
	head.IsFinalized.Store(true)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
	ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(1), nil).Maybe()
	ctx := tests.Context(t)

	t.ResetTimer()
	for n := 0; n < t.N; n++ {
		etx1 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)
		etx2 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 4, fromAddress, 1, blockNum, assets.NewWeiI(1))

		var err error
		t.StartTimer()
		err = ec.CheckForConfirmation(ctx, &head)
		t.StopTimer()
		require.NoError(t, err)

		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx1.State)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)

		deleteTx(ctx, t, &etx1, db)
		deleteTx(ctx, t, &etx2, db)
	}
}
