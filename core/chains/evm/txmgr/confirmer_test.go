package txmgr_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-framework/chains/fees"
	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	gasmocks "github.com/smartcontractkit/chainlink-evm/pkg/gas/mocks"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func newBroadcastLegacyEthTxAttempt(t testing.TB, etxID int64, gasPrice ...int64) txmgr.TxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(gp)}
	}
	return attempt
}

func newTxReceipt(hash gethCommon.Hash, blockNumber int, txIndex uint) evmtypes.Receipt {
	return evmtypes.Receipt{
		TxHash:           hash,
		BlockHash:        testutils.NewHash(),
		BlockNumber:      big.NewInt(int64(blockNumber)),
		TransactionIndex: txIndex,
		Status:           uint64(1),
	}
}

func newInProgressLegacyEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) txmgr.TxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgrtypes.TxAttemptInProgress
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(gp)}
	}
	return attempt
}

func mustInsertInProgressEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress gethCommon.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.State = txmgrcommon.TxInProgress
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	require.NoError(t, txStore.InsertTx(t.Context(), &etx))

	return etx
}

func mustInsertConfirmedEthTx(t testing.TB, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress gethCommon.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.State = txmgrcommon.TxConfirmed
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	now := time.Now()
	etx.BroadcastAt = &now
	etx.InitialBroadcastAt = &now
	require.NoError(t, txStore.InsertTx(t.Context(), &etx))

	return etx
}

func TestEthConfirmer_Lifecycle(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	config := configtest.NewChainScopedConfig(t, nil)
	txStore := newTxStore(t, db)

	memKS := keystest.NewMemoryChainStore()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())

	// Add some fromAddresses
	memKS.MustCreate(t)
	memKS.MustCreate(t)
	estimator := gasmocks.NewEvmEstimator(t)
	newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), config.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), confirmerConfig{}, ethKeyStore, txBuilder, lggr, stuckTxDetector, metrics)
	ctx := t.Context()

	// Can't close unstarted instance
	err = ec.Close()
	require.Error(t, err)

	// Can successfully start once
	err = ec.Start(ctx)
	require.NoError(t, err)

	// Can't start an already started instance
	err = ec.Start(ctx)
	require.Error(t, err)

	latestFinalizedHead := &evmtypes.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	// We are guaranteed to receive a latestFinalizedHead.
	latestFinalizedHead.IsFinalized.Store(true)

	h9 := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	h9.Parent.Store(latestFinalizedHead)
	head := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}
	head.Parent.Store(h9)

	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

	err = ec.ProcessHead(ctx, head)
	require.NoError(t, err)
	// Can successfully close once
	err = ec.Close()
	require.NoError(t, err)

	// Can't start more than once (Confirmer uses services.StateMachine)
	err = ec.Start(ctx)
	require.Error(t, err)
	// Can't close more than once (Confirmer use services.StateMachine)
	err = ec.Close()
	require.Error(t, err)

	// Can't closeInternal unstarted instance
	require.Error(t, ec.XXXTestCloseInternal())

	// Can successfully startInternal a previously closed instance
	require.NoError(t, ec.XXXTestStartInternal())
	// Can't startInternal already started instance
	require.Error(t, ec.XXXTestStartInternal())
	// Can successfully closeInternal again
	require.NoError(t, ec.XXXTestCloseInternal())
}

func TestEthConfirmer_CheckForConfirmation(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	ctx := t.Context()
	blockNum := int64(100)
	head := evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: blockNum,
	}
	head.IsFinalized.Store(true)

	t.Run("does nothing if no re-org'd or included transactions found", func(t *testing.T) {
		fromAddress := cltest.MustGenerateRandomKey(t).Address
		ethKeyStore := &keystest.FakeChainStore{Addresses: keystest.Addresses{fromAddress}}
		etx1 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)
		etx2 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 4, fromAddress, 1, blockNum, assets.NewWeiI(1))
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(1), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx1.State)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
	})

	t.Run("marks re-org'd confirmed transaction as unconfirmed, marks latest attempt as in-progress, deletes receipt", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		// Insert confirmed transaction that stays confirmed
		etx := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(0), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)
		require.Empty(t, attempt.Receipts)
	})

	t.Run("marks re-org'd terminally stuck transaction as unconfirmed, marks latest attempt as in-progress, deletes receipt, removed error", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		// Insert terminally stuck transaction that stays fatal error
		etx := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 0, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx.TxAttempts[0].Hash)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(0), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		require.Equal(t, "", etx.Error.String)
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)
		require.Empty(t, attempt.Receipts)
	})

	t.Run("handles multiple re-org transactions at a time", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		// Insert confirmed transaction that stays confirmed
		etx1 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)
		// Insert terminally stuck transaction that stays fatal error
		etx2 := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx2.TxAttempts[0].Hash)
		// Insert confirmed transaction that gets re-org'd
		etx3 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 2, blockNum)
		// Insert terminally stuck transaction that gets re-org'd
		etx4 := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 3, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx4.TxAttempts[0].Hash)
		// Insert unconfirmed transaction that is untouched
		etx5 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 4, fromAddress, 1, blockNum, assets.NewWeiI(1))
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(2), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx1.State)
		attempt1 := etx1.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1.State)
		require.Len(t, attempt1.Receipts, 1)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx2.State)
		require.Equal(t, client.TerminallyStuckMsg, etx2.Error.String)
		attempt2 := etx2.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2.State)
		require.Len(t, attempt2.Receipts, 1)

		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)
		attempt3 := etx3.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt3.State)
		require.Empty(t, attempt3.Receipts)

		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)
		require.Equal(t, "", etx4.Error.String)
		attempt4 := etx4.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt4.State)
		require.True(t, attempt4.IsPurgeAttempt)
		require.Empty(t, attempt4.Receipts)

		etx5, err = txStore.FindTxWithAttempts(ctx, etx5.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx5.State)
		attempt5 := etx5.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt5.State)
	})

	t.Run("marks valid transaction as confirmed if nonce less than mined tx count", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, blockNum, assets.NewWeiI(1))
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(1), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
	})

	t.Run("marks purge transaction as terminally stuck if nonce less than mined tx count", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 0, fromAddress)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(1), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx.State)
		require.Equal(t, client.TerminallyStuckMsg, etx.Error.String)
	})

	t.Run("handles multiple confirmed transactions at a time", func(t *testing.T) {
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		// Insert valid confirmed transaction that is untouched
		etx1 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)
		// Insert terminally stuck transaction that is untouched
		etx2 := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx2.TxAttempts[0].Hash)
		// Insert valid unconfirmed transaction that is confirmed
		etx3 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 2, fromAddress, 1, blockNum, assets.NewWeiI(1))
		// Insert unconfirmed purge transaction that is confirmed and marked as terminally stuck
		etx4 := mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 3, fromAddress)
		// Insert unconfirmed transact that is not confirmed and left untouched
		etx5 := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 4, fromAddress, 1, blockNum, assets.NewWeiI(1))
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("NonceAt", mock.Anything, fromAddress, mock.Anything).Return(uint64(4), nil).Maybe()
		require.NoError(t, ec.CheckForConfirmation(ctx, &head))

		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx1.State)
		attempt1 := etx1.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1.State)
		require.Len(t, attempt1.Receipts, 1)

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx2.State)
		require.Equal(t, client.TerminallyStuckMsg, etx2.Error.String)
		attempt2 := etx2.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2.State)
		require.Len(t, attempt2.Receipts, 1)

		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
		attempt3 := etx3.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt3.State)
		require.Empty(t, attempt3.Receipts)

		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx4.State)
		require.Equal(t, client.TerminallyStuckMsg, etx4.Error.String)
		attempt4 := etx4.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt4.State)
		require.True(t, attempt4.IsPurgeAttempt)
		require.Empty(t, attempt4.Receipts)

		etx5, err = txStore.FindTxWithAttempts(ctx, etx5.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx5.State)
		attempt5 := etx5.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt5.State)
		require.Empty(t, attempt3.Receipts)
	})
}

func TestEthConfirmer_FindTxsRequiringRebroadcast(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := t.Context()

	ethClient := clienttest.NewClientWithDefaultChainID(t)

	evmcfg := configtest.NewChainScopedConfig(t, nil)

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
	evmFromAddress := fromAddress
	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	mustInsertConfirmedEthTx(t, txStore, nonce, fromAddress)
	nonce++

	evmOtherAddress := memKS.MustCreate(t)

	lggr := logger.Test(t)

	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Empty(t, etxs)
	})

	mustInsertInProgressEthTx(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Empty(t, etxs)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Empty(t, etxs)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx1.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt1_1.ID))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Empty(t, etxs)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt2_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt2_1.ID))

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Empty(t, etxs)
	})

	etxWithoutAttempts := cltest.NewEthTx(fromAddress)
	{
		n := evmtypes.Nonce(nonce)
		etxWithoutAttempts.Sequence = &n
	}
	now := time.Now()
	etxWithoutAttempts.BroadcastAt = &now
	etxWithoutAttempts.InitialBroadcastAt = &now
	etxWithoutAttempts.State = txmgrcommon.TxUnconfirmed
	require.NoError(t, txStore.InsertTx(ctx, &etxWithoutAttempts))
	nonce++

	t.Run("does nothing if the transaction is from a different address than the one given", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmOtherAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Empty(t, etxs)
	})

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	t.Run("returns nothing for different chain id", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, big.NewInt(42))
		require.NoError(t, err)

		require.Empty(t, etxs)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt3_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_1.ID))

	// NOTE: It should ignore qualifying eth_txes from a different address
	etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, evmOtherAddress)
	attemptOther1 := etxOther.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attemptOther1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attemptOther1.ID))

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	t.Run("returns nothing if threshold is zero", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, 0, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Empty(t, etxs)
	})

	t.Run("does not return more transactions for gas bumping than gasBumpThreshold", func(t *testing.T) {
		// Unconfirmed txes in DB are:
		// (unnamed) (nonce 2)
		// etx1 (nonce 3)
		// etx2 (nonce 4)
		// etxWithoutAttempts (nonce 5)
		// etx3 (nonce 6) - ready for bump
		// etx4 (nonce 7) - ready for bump
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 4, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1) // returns etxWithoutAttempts only - eligible for gas bumping because it technically doesn't have any attempts within gasBumpThreshold blocks
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)

		etxs, err = ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 5, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)

		// Zero limit disables it
		etxs, err = ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 0, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt4_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt4_1.ID))

	t.Run("ignores pending transactions for another key", func(t *testing.T) {
		// Re-use etx3 nonce for another key, it should not affect the results for this key
		etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, (*etx3.Sequence).Int64(), evmOtherAddress)
		aOther := etxOther.TxAttempts[0]
		dbAttempt = txmgr.DbEthTxAttempt{}
		dbAttempt.FromTxAttempt(&aOther)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, aOther.ID))

		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 6, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_2))

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_3 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_3))

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		require.Equal(t, evmtypes.Nonce(5), *etxWithoutAttempts.Sequence)
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		require.Equal(t, evmtypes.Nonce(7), *etx4.Sequence)
	})

	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etxWithoutAttempts.ID)
	attempt0_1.State = txmgrtypes.TxAttemptInsufficientFunds
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_1))

	// This attempt has insufficient_eth, but there is also another attempt4_1
	// which is old enough, so this will be caught by both queries and should
	// not be duplicated
	attempt4_2 := cltest.NewLegacyEthTxAttempt(t, etx4.ID)
	attempt4_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt4_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_2))

	etx5 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	nonce++

	// This etx has one attempt that is too new, which would exclude it from
	// the gas bumping query, but it should still be caught by the insufficient
	// eth query
	etx6 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	attempt6_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt6_2.BroadcastBeforeBlockNum = &tooNew
	attempt6_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(30001)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt6_2))

	t.Run("returns unique attempts requiring resubmission due to insufficient eth, ordered by nonce asc", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 4)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		assert.Equal(t, etx5.ID, etxs[2].ID)
		assert.Equal(t, *etx5.Sequence, *(etxs[2].Sequence))
		assert.Equal(t, etx6.ID, etxs[3].ID)
		assert.Equal(t, *etx6.Sequence, *(etxs[3].Sequence))
	})

	t.Run("applies limit", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(t.Context(), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 2, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_WithConnectivityCheck(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	db := testutils.NewSqlxDB(t)
	ethClient := clienttest.NewClientWithDefaultChainID(t)

	t.Run("should retry previous attempt if connectivity check failed for legacy transactions", func(t *testing.T) {
		ccfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.EIP1559DynamicFees = ptr(false)
			c.GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})

		ctx := t.Context()
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

		estimator := gasmocks.NewEvmEstimator(t)
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		estimator.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, uint64(0), pkgerrors.Wrapf(fees.ErrConnectivity, "transaction..."))
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
		kst := keystest.Addresses{fromAddress}
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, keystest.TxSigner(nil), feeEstimator)
		stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), ccfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
		metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
		require.NoError(t, err)
		// Create confirmer with necessary state
		ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmFeeConfig(ccfg.EVM().GasEstimator()), ccfg.EVM().Transactions(), confirmerConfig{}, kst, txBuilder, lggr, stuckTxDetector, metrics)
		servicetest.Run(t, ec)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(multinode.Successful, nil).Once()

		err = ec.RebroadcastWhereNecessary(t.Context(), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("should retry previous attempt if connectivity check failed for dynamic transactions", func(t *testing.T) {
		ccfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.EIP1559DynamicFees = ptr(true)
			c.GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})

		ctx := t.Context()
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

		estimator := gasmocks.NewEvmEstimator(t)
		estimator.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.DynamicFee{}, pkgerrors.Wrapf(fees.ErrConnectivity, "transaction..."))
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		// Create confirmer with necessary state
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
		kst := keystest.Addresses{fromAddress}
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, keystest.TxSigner(nil), feeEstimator)
		stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), ccfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
		metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
		require.NoError(t, err)
		ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmFeeConfig(ccfg.EVM().GasEstimator()), ccfg.EVM().Transactions(), confirmerConfig{}, kst, txBuilder, lggr, stuckTxDetector, metrics)
		servicetest.Run(t, ec)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(multinode.Successful, nil).Once()

		err = ec.RebroadcastWhereNecessary(t.Context(), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_MaxFeeScenario(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := t.Context()

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	kst := &keystest.FakeChainStore{Addresses: keystest.Addresses{fromAddress}}
	// Use a mock keystore for this test
	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, kst, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

	t.Run("treats an exceeds max fee attempt as a success", func(t *testing.T) {
		// Once for the bumped attempt which exceeds limit
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.GasPrice().Int64() == int64(20000000000) && tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.ExceedsMaxFee, errors.New("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// Check that the attempt is saved
		require.Len(t, etx.TxAttempts, 2)

		// broadcast_at did change
		require.Greater(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		require.Equal(t, etx.InitialBroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
		c.GasEstimator.BumpMin = assets.NewWeiI(0)
	})
	currentHead := int64(30)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())

		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
	})

	t.Run("re-sends previous transaction on keystore error", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		fromAddress := keystest.NewMemoryChainStore().MustCreate(t)
		kst := &keystest.FakeChainStore{
			Addresses: keystest.Addresses{fromAddress},
			TxSigner: func(ctx context.Context, from common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if from == fromAddress {
					return nil, errors.New("signing error")
				}
				return tx, nil
			},
		}

		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, 25, assets.NewWeiI(100))

		// Use a mock keystore for this test
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, kst, nil)

		err := ec.RebroadcastWhereNecessary(ctx, currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("does nothing and continues on fatal error", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, 25, assets.NewWeiI(100))
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Fatal, errors.New("exceeds block gas limit")).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		latestGasPrice := assets.GWei(20)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, 25, latestGasPrice)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)

		// Got the new attempt
		bumpAttempt := etx.TxAttempts[0]
		expectedBumpedGas := latestGasPrice.AddPercentage(evmcfg.EVM().GasEstimator().BumpPercent())
		require.Equal(t, expectedBumpedGas.Int64(), bumpAttempt.TxFee.GasPrice.Int64())
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, bumpAttempt.State)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 0, fromAddress, txmgrtypes.TxAttemptBroadcast)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		latestGasPrice := assets.GWei(20)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, 25, latestGasPrice)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, fmt.Errorf("known transaction: %s", etx.TxAttempts[0].Hash.Hex())).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)

		// Got the new attempt
		bumpAttempt := etx.TxAttempts[0]
		expectedBumpedGas := latestGasPrice.AddPercentage(evmcfg.EVM().GasEstimator().BumpPercent())
		require.Equal(t, expectedBumpedGas.Int64(), bumpAttempt.TxFee.GasPrice.Int64())
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, bumpAttempt.State)
	})

	t.Run("saves new attempt even for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		latestGasPrice := assets.GWei(20)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, 25, latestGasPrice)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)

		// Got the new attempt
		// Got the new attempt
		bumpedAttempt := etx.TxAttempts[0]
		expectedBumpedGas := latestGasPrice.AddPercentage(evmcfg.EVM().GasEstimator().BumpPercent())
		require.Equal(t, expectedBumpedGas.Int64(), bumpedAttempt.TxFee.GasPrice.Int64())

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[1].State)
	})

	t.Run("saves in-progress attempt on temporary error and returns error", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		latestGasPrice := assets.GWei(20)
		broadcastBlockNum := int64(25)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, broadcastBlockNum, latestGasPrice)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Unknown, errors.New("some network error")).Once()

		err := ec.RebroadcastWhereNecessary(ctx, currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// Old attempt is untouched
		require.Len(t, etx.TxAttempts, 2)
		originalAttempt := etx.TxAttempts[1]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, originalAttempt.State)
		require.Equal(t, broadcastBlockNum, *originalAttempt.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		bumpedAttempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, bumpedAttempt.State)
		require.Nil(t, bumpedAttempt.BroadcastBeforeBlockNum)

		// Try again and move the attempt into "broadcast"
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// New in_progress attempt saved and marked "broadcast"
		require.Len(t, etx.TxAttempts, 2)
		bumpedAttempt = etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, bumpedAttempt.State)
		require.Nil(t, bumpedAttempt.BroadcastBeforeBlockNum)
	})

	t.Run("re-bumps attempt if initial bump is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		latestGasPrice := assets.GWei(20)
		broadcastBlockNum := int64(25)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, broadcastBlockNum, latestGasPrice)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Underpriced, errors.New("replacement transaction underpriced")).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		require.Len(t, etx.TxAttempts, 2)
		bumpedAttempt := etx.TxAttempts[0]
		expectedBumpedGas := latestGasPrice.AddPercentage(evmcfg.EVM().GasEstimator().BumpPercent())
		expectedBumpedGas = expectedBumpedGas.AddPercentage(evmcfg.EVM().GasEstimator().BumpPercent())
		require.Equal(t, expectedBumpedGas.Int64(), bumpedAttempt.TxFee.GasPrice.Int64())
	})

	t.Run("resubmits at the old price and does not create a new attempt if one of the bumped transactions would exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		priceMax := assets.GWei(30)
		newCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.PriceMax = priceMax
		})
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		broadcastBlockNum := int64(25)
		currentAttemptPrice := priceMax.Sub(assets.GWei(1))
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, broadcastBlockNum, currentAttemptPrice)
		ec := newEthConfirmer(t, txStore, ethClient, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Underpriced, errors.New("underpriced")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.Error(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// No new tx attempts
		require.Len(t, etx.TxAttempts, 1)
		bumpedAttempt := etx.TxAttempts[0]
		require.Equal(t, currentAttemptPrice.Int64(), bumpedAttempt.TxFee.GasPrice.Int64())
	})

	t.Run("resubmits at the old price and does not create a new attempt if the current price is exactly EVM.GasEstimator.PriceMax", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		priceMax := assets.GWei(30)
		newCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.PriceMax = priceMax
		})
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		broadcastBlockNum := int64(25)
		etx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, 0, fromAddress, 1, broadcastBlockNum, priceMax)
		ec := newEthConfirmer(t, txStore, ethClient, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Underpriced, errors.New("underpriced")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.Error(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// No new tx attempts
		require.Len(t, etx.TxAttempts, 1)
		bumpedAttempt := etx.TxAttempts[0]
		require.Equal(t, priceMax.Int64(), bumpedAttempt.TxFee.GasPrice.Int64())
	})

	t.Run("EIP-1559: bumps using EIP-1559 rules when existing attempts are of type 0x2", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		newCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.BumpMin = assets.GWei(1)
		})
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, 0, fromAddress)
		err := txStore.UpdateTxAttemptBroadcastBeforeBlockNum(ctx, etx.ID, uint(25))
		require.NoError(t, err)
		ec := newEthConfirmer(t, txStore, ethClient, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()
		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// A new, bumped attempt
		require.Len(t, etx.TxAttempts, 2)
		bumpAttempt := etx.TxAttempts[0]
		require.Nil(t, bumpAttempt.TxFee.GasPrice)
		bumpedGas := assets.NewWeiI(1).Add(newCfg.EVM().GasEstimator().BumpMin())
		require.Equal(t, bumpedGas.Int64(), bumpAttempt.TxFee.GasTipCap.Int64())
		require.Equal(t, bumpedGas.Int64(), bumpAttempt.TxFee.GasFeeCap.Int64())
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, bumpAttempt.State)
	})

	t.Run("EIP-1559: resubmits at the old price and does not create a new attempt if one of the bumped EIP-1559 transactions would have its tip cap exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		newCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.PriceMax = assets.NewWeiI(1)
		})
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, 0, fromAddress)
		err := txStore.UpdateTxAttemptBroadcastBeforeBlockNum(ctx, etx.ID, uint(25))
		require.NoError(t, err)
		ec := newEthConfirmer(t, txStore, ethClient, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Underpriced, errors.New("underpriced")).Once()

		require.Error(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		// No new tx attempts
		require.Len(t, etx.TxAttempts, 1)
		bumpedAttempt := etx.TxAttempts[0]
		assert.Equal(t, assets.NewWeiI(1).Int64(), bumpedAttempt.TxFee.GasTipCap.Int64())
		assert.Equal(t, assets.NewWeiI(1).Int64(), bumpedAttempt.TxFee.GasFeeCap.Int64())
	})

	t.Run("EIP-1559: re-bumps attempt if initial bump is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		newCfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.BumpMin = assets.GWei(1)
		})
		// NOTE: This test case was empirically impossible when I tried it on eth mainnet (any EIP1559 transaction with a higher tip cap is accepted even if it's only 1 wei more) but appears to be possible on Polygon/Matic, probably due to poor design that applies the 10% minimum to the overall value (base fee + tip cap)
		memKS := keystest.NewMemoryChainStore()
		fromAddress := memKS.MustCreate(t)
		ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())
		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, 0, fromAddress)
		err := txStore.UpdateTxAttemptBroadcastBeforeBlockNum(ctx, etx.ID, uint(25))
		require.NoError(t, err)
		ec := newEthConfirmer(t, txStore, ethClient, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Underpriced, errors.New("replacement transaction underpriced")).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()

		// Do it
		require.NoError(t, ec.RebroadcastWhereNecessary(ctx, currentHead))
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		require.Len(t, etx.TxAttempts, 2)
		bumpAttempt := etx.TxAttempts[0]
		bumpedGas := assets.NewWeiI(1).Add(newCfg.EVM().GasEstimator().BumpMin())
		bumpedGas = bumpedGas.Add(newCfg.EVM().GasEstimator().BumpMin())
		assert.Equal(t, bumpedGas.Int64(), bumpAttempt.TxFee.GasTipCap.Int64())
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_TerminallyUnderpriced_ThenGoesThrough(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	kst := keys.NewChainStore(memKS, big.NewInt(0))

	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	currentHead := int64(30)
	oldEnough := 5
	nonce := int64(0)

	t.Run("terminally underpriced transaction with in_progress attempt is retried with more gas", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, kst, nil)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, nonce, fromAddress, txmgrtypes.TxAttemptInProgress, originalBroadcastAt)
		require.Equal(t, originalBroadcastAt, *etx.BroadcastAt)
		nonce++

		// Fail the first time with terminally underpriced.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Once()
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in legacy mode", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, kst, nil)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
		nonce++
		legacyAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&legacyAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, legacyAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in EIP-1559 mode", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, kst, nil)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress)
		nonce++
		dxFeeAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&dxFeeAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, dxFeeAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Underpriced, errors.New("transaction underpriced")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			multinode.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_WhenOutOfEth(t *testing.T) {
	t.Parallel()
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := t.Context()

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())

	config := configtest.NewChainScopedConfig(t, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))
	var attempt1_2 txmgr.TxAttempt

	insufficientEthError := errors.New("insufficient funds for gas * price + value")

	t.Run("saves attempt with state 'insufficient_eth' if eth node returns this error", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.GasPrice.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(multinode.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
		assert.Nil(t, attempt1_2.BroadcastBeforeBlockNum)
	})

	t.Run("does not bump gas when previous error was 'out of eth', instead resubmits existing transaction", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.GasPrice.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(multinode.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// The attempt is still "out of eth"
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
	})

	t.Run("saves the attempt as broadcast after node wallet has been topped up with sufficient balance", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.GasPrice.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(multinode.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))

		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// Attempt is now 'broadcast'
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	t.Run("resubmitting due to insufficient eth is not limited by EVM.GasEstimator.BumpTxDepth", func(t *testing.T) {
		depth := 2
		etxCount := 4

		evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.GasEstimator.BumpTxDepth = ptr(uint32(depth))
		})
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)

		for i := 0; i < etxCount; i++ {
			n := nonce
			mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(n)
			}), fromAddress).Return(multinode.Successful, nil).Once()

			nonce++
		}

		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))

		var dbAttempts []txmgr.DbEthTxAttempt

		require.NoError(t, db.Select(&dbAttempts, "SELECT * FROM evm.tx_attempts WHERE state = 'insufficient_eth'"))
		require.Empty(t, dbAttempts)
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_TerminallyStuckError(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := t.Context()

	ethClient := clienttest.NewClientWithDefaultChainID(t)
	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	ethKeyStore := keys.NewChainStore(memKS, ethClient.ConfiguredChainID())

	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.PriceMax = assets.GWei(500)
	})

	// Use a mock keystore for this test
	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)
	terminallyStuckError := "failed to add tx to the pool: not enough step counters to continue the execution"

	t.Run("terminally stuck transaction replaced with purge attempt", func(t *testing.T) {
		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		nonce++
		attempt1_1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

		// Return terminally stuck error on first rebroadcast
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.TerminallyStuck, errors.New(terminallyStuckError)).Once()
		// Return successful for purge attempt
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) //nolint:gosec // disable G115
		}), fromAddress).Return(multinode.Successful, nil).Once()

		// Start processing transactions for rebroadcast
		require.NoError(t, ec.RebroadcastWhereNecessary(t.Context(), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		purgeAttempt := etx.TxAttempts[0]
		require.True(t, purgeAttempt.IsPurgeAttempt)
	})
}

func TestEthConfirmer_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)

	memKS := keystest.NewMemoryChainStore()
	fromAddress := memKS.MustCreate(t)
	config := configtest.NewChainScopedConfig(t, nil)
	ethKeyStore := keys.NewChainStore(memKS, big.NewInt(0))

	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())
	mustInsertInProgressEthTx(t, txStore, 0, fromAddress)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	gasPriceWei := gas.EvmFee{GasPrice: assets.GWei(52)}
	overrideGasLimit := uint64(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() &&
				tx.Gas() == overrideGasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(t.Context(), []evmtypes.Nonce{1}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() &&
				tx.Gas() == etx1.FeeLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(t.Context(), []evmtypes.Nonce{(1)}, gasPriceWei, fromAddress, 0))
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) && tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() && tx.Gas() == overrideGasLimit
		}), mock.Anything).Return(multinode.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Sequence) && tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() && tx.Gas() == overrideGasLimit
		}), mock.Anything).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(t.Context(), []evmtypes.Nonce{(1), (2)}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(1)
		}), mock.Anything).Return(multinode.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(2)
		}), mock.Anything).Return(multinode.Successful, nil).Once()
		for i := 3; i <= 5; i++ {
			nonce := i
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(nonce) &&
					tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() &&
					tx.Gas() == overrideGasLimit &&
					*tx.To() == fromAddress &&
					tx.Value().Cmp(big.NewInt(0)) == 0 &&
					len(tx.Data()) == 0
			}), mock.Anything).Return(multinode.Successful, nil).Once()
		}
		nonces := []evmtypes.Nonce{(1), (2), (3), (4), (5)}

		require.NoError(t, ec.ForceRebroadcast(t.Context(), nonces, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := clienttest.NewClientWithDefaultChainID(t)
		ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && tx.GasPrice().Int64() == gasPriceWei.GasPrice.Int64() && tx.Gas() == config.EVM().GasEstimator().LimitDefault()
		}), mock.Anything).Return(multinode.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(t.Context(), []evmtypes.Nonce{(0)}, gasPriceWei, fromAddress, 0))
	})
}

func TestEthConfirmer_ProcessStuckTransactions(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	fromAddress := cltest.MustGenerateRandomKey(t)
	ethKeyStore := &keystest.FakeChainStore{Addresses: keystest.Addresses{fromAddress.Address}}
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress.Address).Return(multinode.Successful, nil).Once()
	lggr := logger.Test(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)

	// Return 10 gwei as market gas price
	marketGasPrice := tenGwei
	fee := gas.EvmFee{GasPrice: marketGasPrice}
	bumpedLegacy := assets.GWei(30)
	bumpedFee := gas.EvmFee{GasPrice: bumpedLegacy}
	feeEstimator.On("GetFee", mock.Anything, []byte{}, uint64(0), mock.Anything, mock.Anything, mock.Anything).Return(fee, uint64(0), nil)
	feeEstimator.On("BumpFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(bumpedFee, uint64(10_000), nil)
	autoPurgeThreshold := uint32(5)
	autoPurgeMinAttempts := uint32(3)
	limitDefault := uint64(100)
	evmcfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
		c.GasEstimator.LimitDefault = ptr(limitDefault)
		c.Transactions.AutoPurge.Enabled = ptr(true)
		c.Transactions.AutoPurge.Threshold = ptr(autoPurgeThreshold)
		c.Transactions.AutoPurge.MinAttempts = ptr(autoPurgeMinAttempts)
	})
	ge := evmcfg.EVM().GasEstimator()
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), evmcfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmFeeConfig(ge), evmcfg.EVM().Transactions(), confirmerConfig{}, ethKeyStore, txBuilder, lggr, stuckTxDetector, metrics)
	fn := func(ctx context.Context, id uuid.UUID, result interface{}, err error) error {
		require.ErrorContains(t, err, client.TerminallyStuckMsg)
		return nil
	}
	ec.SetResumeCallback(fn)
	servicetest.Run(t, ec)

	ctx := t.Context()
	blockNum := int64(100)

	t.Run("detects and processes stuck transactions", func(t *testing.T) {
		nonce := int64(0)
		// Create attempts so that the oldest broadcast attempt's block num is what meets the threshold check
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		tx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, nonce, fromAddress.Address, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold), marketGasPrice.Add(oneGwei))
		// Update tx to signal callback once it is identified as terminally stuck
		testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, signal_callback = TRUE WHERE id = $2`, uuid.New(), tx.ID)
		head := evmtypes.Head{
			Hash:   testutils.NewHash(),
			Number: blockNum,
		}
		head.IsFinalized.Store(true)

		// Mined tx count does not increment due to terminally stuck transaction
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil).Once()

		// First call to ProcessHead should:
		// 1. Detect a stuck transaction
		// 2. Create a purge attempt for it
		// 3. Save the purge attempt to the DB
		// 4. Send the purge attempt
		err := ec.ProcessHead(ctx, &head)
		require.NoError(t, err)

		// Check if the purge attempt was saved to the DB properly
		dbTx, err := txStore.FindTxWithAttempts(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, dbTx)
		latestAttempt := dbTx.TxAttempts[0]
		require.True(t, latestAttempt.IsPurgeAttempt)
		require.Equal(t, limitDefault, latestAttempt.ChainSpecificFeeLimit)
		require.Equal(t, bumpedFee.GasPrice, latestAttempt.TxFee.GasPrice)

		head = evmtypes.Head{
			Hash:   testutils.NewHash(),
			Number: blockNum + 1,
		}
		// Mined tx count incremented because of purge attempt
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(1), nil)

		// Second call to ProcessHead on next head should:
		// 1. Check for receipts for purged transaction
		// 2. When receipts are found for a purge attempt, the transaction is marked in the DB as fatal error with error message
		err = ec.ProcessHead(ctx, &head)
		require.NoError(t, err)
		dbTx, err = txStore.FindTxWithAttempts(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, dbTx)
		require.Equal(t, txmgrcommon.TxFatalError, dbTx.State)
		require.Equal(t, client.TerminallyStuckMsg, dbTx.Error.String)
		require.True(t, dbTx.CallbackCompleted)
	})
}

func ptr[T any](t T) *T { return &t }

func newEthConfirmer(t testing.TB, txStore txmgr.EvmTxStore, ethClient client.Client, config evmconfig.ChainScopedConfig, ks keys.ChainStore, fn txmgrcommon.ResumeCallback) *txmgr.Confirmer {
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	estimator := gas.NewEvmFeeEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(ge, nil, ge.BlockHistory(), lggr, nil)
	}, ge.EIP1559DynamicFees(), ge, ethClient)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ks, estimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), config.EVM().Transactions().AutoPurge(), estimator, txStore, ethClient)
	metrics, err := txmgr.NewEVMTxmMetrics(ethClient.ConfiguredChainID().String())
	require.NoError(t, err)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), confirmerConfig{}, ks, txBuilder, lggr, stuckTxDetector, metrics)
	ec.SetResumeCallback(fn)
	servicetest.Run(t, ec)
	return ec
}

var _ txmgrtypes.ConfirmerDatabaseConfig = confirmerConfig{}

type confirmerConfig struct{}

func (d confirmerConfig) DefaultQueryTimeout() time.Duration {
	return 10 * time.Second
}
