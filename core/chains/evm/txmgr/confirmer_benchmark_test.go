package txmgr_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func BenchmarkEthConfirmer(b *testing.B) {
	db := testutils.NewSqlxDB(b)
	txStore := cltest.NewTestTxStore(b, db)
	ethClient := clienttest.NewClientWithDefaultChainID(b)
	evmcfg := configtest.NewChainScopedConfig(b, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	blockNum := int64(100)
	head := evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: blockNum,
	}
	head.IsFinalized.Store(true)

	memKeystore := keystest.NewMemoryChainStore()
	ethKeyStore := keys.NewChainStore(memKeystore, ethClient.ConfiguredChainID())
	fromAddress := memKeystore.MustCreate(b)
	ec := newEthConfirmer(b, txStore, ethClient, evmcfg, ethKeyStore, nil)
	ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(1), nil).Maybe()
	ctx := tests.Context(b)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		etx1 := mustInsertConfirmedEthTxWithReceipt(b, txStore, fromAddress, 0, blockNum)
		etx2 := mustInsertUnconfirmedTxWithBroadcastAttempts(b, txStore, 4, fromAddress, 1, blockNum, assets.NewWeiI(1))

		var err error
		b.StartTimer()
		err = ec.CheckForConfirmation(ctx, &head)
		b.StopTimer()
		require.NoError(b, err)

		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(b, err)
		require.Equal(b, txmgrcommon.TxConfirmed, etx1.State)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(b, err)
		require.Equal(b, txmgrcommon.TxUnconfirmed, etx2.State)

		deleteTx(ctx, b, &etx1, db)
		deleteTx(ctx, b, &etx2, db)
	}
}
