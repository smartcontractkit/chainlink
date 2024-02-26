package txmgr_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	persistentStore := cltest.NewTestTxStore(t, db, cfg.Database())
	kst := cltest.NewKeyStore(t, db, cfg.Database())

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := context.Background()

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore)
	require.NoError(t, err)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		actTx, err := inMemoryStore.CreateTransaction(testutils.Context(t), evmtxmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, chainID)
		require.NoError(t, err)

		assert.Greater(t, actTx.ID, int64(0))
		assert.Equal(t, commontxmgr.TxUnstarted, actTx.State)
		assert.Equal(t, gasLimit, actTx.FeeLimit)
		assert.Equal(t, fromAddress, actTx.FromAddress)
		assert.Equal(t, toAddress, actTx.ToAddress)
		assert.Equal(t, payload, actTx.EncodedPayload)
		assert.Equal(t, big.Int(evmassets.NewEthValue(0)), actTx.Value)
		assert.Equal(t, subject, actTx.Subject.UUID)

		cltest.AssertCount(t, db, "evm.txes", 1)

		var dbEthTx evmtxmgr.DbEthTx
		require.NoError(t, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, commontxmgr.TxUnstarted, dbEthTx.State)
		assert.Equal(t, gasLimit, dbEthTx.GasLimit)
		assert.Equal(t, fromAddress, dbEthTx.FromAddress)
		assert.Equal(t, toAddress, dbEthTx.ToAddress)
		assert.Equal(t, payload, dbEthTx.EncodedPayload)
		assert.Equal(t, evmassets.NewEthValue(0), dbEthTx.Value)
		assert.Equal(t, subject, dbEthTx.Subject.UUID)

		var expTx evmtxmgr.Tx
		dbEthTx.ToTx(&expTx)
		assertTxEqual(t, expTx, actTx)
	})
}

// assertTxEqual asserts that two transactions are equal
func assertTxEqual(t *testing.T, exp, act evmtxmgr.Tx) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.IdempotencyKey, act.IdempotencyKey)
	assert.Equal(t, exp.Sequence, act.Sequence)
	assert.Equal(t, exp.FromAddress, act.FromAddress)
	assert.Equal(t, exp.ToAddress, act.ToAddress)
	assert.Equal(t, exp.EncodedPayload, act.EncodedPayload)
	assert.Equal(t, exp.Value, act.Value)
	assert.Equal(t, exp.FeeLimit, act.FeeLimit)
	assert.Equal(t, exp.Error, act.Error)
	assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	assert.Equal(t, exp.InitialBroadcastAt, act.InitialBroadcastAt)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.Meta, act.Meta)
	assert.Equal(t, exp.Subject, act.Subject)
	assert.Equal(t, exp.ChainID, act.ChainID)
	assert.Equal(t, exp.PipelineTaskRunID, act.PipelineTaskRunID)
	assert.Equal(t, exp.MinConfirmations, act.MinConfirmations)
	assert.Equal(t, exp.TransmitChecker, act.TransmitChecker)
	assert.Equal(t, exp.SignalCallback, act.SignalCallback)
	assert.Equal(t, exp.CallbackCompleted, act.CallbackCompleted)

	require.Len(t, exp.TxAttempts, len(act.TxAttempts))
	for i := 0; i < len(exp.TxAttempts); i++ {
		assertTxAttemptEqual(t, exp.TxAttempts[i], act.TxAttempts[i])
	}
}

func assertTxAttemptEqual(t *testing.T, exp, act evmtxmgr.TxAttempt) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.TxID, act.TxID)
	assert.Equal(t, exp.Tx, act.Tx)
	assert.Equal(t, exp.TxFee, act.TxFee)
	assert.Equal(t, exp.ChainSpecificFeeLimit, act.ChainSpecificFeeLimit)
	assert.Equal(t, exp.SignedRawTx, act.SignedRawTx)
	assert.Equal(t, exp.Hash, act.Hash)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.BroadcastBeforeBlockNum, act.BroadcastBeforeBlockNum)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.TxType, act.TxType)

	require.Equal(t, len(exp.Receipts), len(act.Receipts))
	for i := 0; i < len(exp.Receipts); i++ {
		assertChainReceiptEqual(t, exp.Receipts[i], act.Receipts[i])
	}
}

func assertChainReceiptEqual(t *testing.T, exp, act evmtxmgr.ChainReceipt) {
	assert.Equal(t, exp.GetStatus(), act.GetStatus())
	assert.Equal(t, exp.GetTxHash(), act.GetTxHash())
	assert.Equal(t, exp.GetBlockNumber(), act.GetBlockNumber())
	assert.Equal(t, exp.IsZero(), act.IsZero())
	assert.Equal(t, exp.IsUnmined(), act.IsUnmined())
	assert.Equal(t, exp.GetFeeUsed(), act.GetFeeUsed())
	assert.Equal(t, exp.GetTransactionIndex(), act.GetTransactionIndex())
	assert.Equal(t, exp.GetBlockHash(), act.GetBlockHash())
}
