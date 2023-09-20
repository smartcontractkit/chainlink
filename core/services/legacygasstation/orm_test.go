package legacygasstation_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/sqlx"
	"github.com/test-go/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation"
	"github.com/smartcontractkit/capital-markets-projects/lib/services/legacygasstation/types"
	txmgrstate "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	lgsservice "github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation"
)

func TestORM_Insert(t *testing.T) {
	orm, _, txStore, ethKeyStore := setup(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	etx := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

	tx := legacygasstation.LegacyGaslessTx(t, legacygasstation.TestLegacyGaslessTx{
		EthTxID: etx.GetID(),
	})
	err := orm.InsertLegacyGaslessTx(testutils.Context(t), tx)
	require.NoError(t, err)

	txs, err := orm.SelectBySourceChainIDAndStatus(testutils.Context(t), tx.SourceChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])

	txs, err = orm.SelectByDestChainIDAndStatus(testutils.Context(t), tx.DestinationChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])
}

func TestORM_MultipleInserts(t *testing.T) {
	orm, _, txStore, ethKeyStore := setup(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	etx := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

	nonce1, ok := new(big.Int).SetString("1", 10)
	require.True(t, ok)
	nonce2, ok := new(big.Int).SetString("2", 10)
	require.True(t, ok)

	tx1 := legacygasstation.LegacyGaslessTx(t, legacygasstation.TestLegacyGaslessTx{
		EthTxID: etx.GetID(),
		Nonce:   nonce1,
	})
	tx2 := legacygasstation.LegacyGaslessTx(t, legacygasstation.TestLegacyGaslessTx{
		EthTxID: etx.GetID(),
		Nonce:   nonce2,
	})
	err := orm.InsertLegacyGaslessTx(testutils.Context(t), tx1)
	require.NoError(t, err)
	err = orm.InsertLegacyGaslessTx(testutils.Context(t), tx2)
	require.NoError(t, err)

	txs, err := orm.SelectBySourceChainIDAndStatus(testutils.Context(t), tx1.SourceChainID, tx1.Status)
	require.NoError(t, err)
	require.Equal(t, 2, len(txs))

	txs, err = orm.SelectByDestChainIDAndStatus(testutils.Context(t), tx1.DestinationChainID, tx1.Status)
	require.NoError(t, err)
	require.Equal(t, 2, len(txs))
}

func TestORM_Update(t *testing.T) {
	orm, _, txStore, ethKeyStore := setup(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	etx := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

	tx := legacygasstation.LegacyGaslessTx(t, legacygasstation.TestLegacyGaslessTx{
		EthTxID: etx.GetID(),
	})
	err := orm.InsertLegacyGaslessTx(testutils.Context(t), tx)
	require.NoError(t, err)

	tx.Status = types.SourceFinalized
	ccipMessageID := common.HexToHash("1")
	tx.CCIPMessageID = &ccipMessageID

	err = orm.UpdateLegacyGaslessTx(testutils.Context(t), tx)
	require.NoError(t, err)

	txs, err := orm.SelectBySourceChainIDAndStatus(testutils.Context(t), tx.SourceChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])

	txs, err = orm.SelectByDestChainIDAndStatus(testutils.Context(t), tx.DestinationChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])

	tx.Status = types.Failure
	failureReason := "executionReverted"
	tx.FailureReason = &failureReason

	err = orm.UpdateLegacyGaslessTx(testutils.Context(t), tx)
	require.NoError(t, err)

	txs, err = orm.SelectBySourceChainIDAndStatus(testutils.Context(t), tx.SourceChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])

	txs, err = orm.SelectByDestChainIDAndStatus(testutils.Context(t), tx.DestinationChainID, tx.Status)
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	legacygasstation.AssertTxEquals(t, tx, txs[0])
}

func TestORM_FailedEthTx(t *testing.T) {
	orm, _, txStore, ethKeyStore := setup(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	etx := cltest.MustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)
	errorMsg := "execution reverted"
	etx.Error = null.StringFrom(errorMsg)
	err := txStore.UpdateTxFatalError(&etx)
	require.NoError(t, err)

	tx := legacygasstation.LegacyGaslessTx(t, legacygasstation.TestLegacyGaslessTx{
		EthTxID: etx.GetID(),
	})
	err = orm.InsertLegacyGaslessTx(testutils.Context(t), tx)
	require.NoError(t, err)

	txs, err := orm.SelectBySourceChainIDAndEthTxStates(testutils.Context(t), tx.SourceChainID, []legacygasstation.EtxStatus{legacygasstation.EtxInProgress})
	require.NoError(t, err)
	require.Equal(t, 0, len(txs))

	txs, err = orm.SelectBySourceChainIDAndEthTxStates(testutils.Context(t), tx.SourceChainID, []legacygasstation.EtxStatus{legacygasstation.EtxFatalError})
	require.NoError(t, err)
	require.Equal(t, 1, len(txs))
	require.Equal(t, txs[0].EthTxStatus, string(txmgrstate.TxFatalError))
	require.Equal(t, *txs[0].EthTxError, errorMsg)
}

func setup(t *testing.T) (legacygasstation.ORM, *sqlx.DB, txmgr.TestEvmTxStore, keystore.Eth) {
	cfg := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	evmtest.NewChainScopedConfig(t, cfg)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	orm := lgsservice.NewORM(db, logger.TestLogger(t), cfg.Database())
	return orm, db, txStore, ethKeyStore
}
