package adapters_test

import (
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthTxAdapter_Perform_BPTXM(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	toAddress := cltest.NewAddress()
	gasLimit := uint64(42)
	functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
	dataPrefix := hexutil.MustDecode("0x88888888")

	t.Run("with valid data and empty DataFormat writes to database and returns run output pending outgoing confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())

		etrt, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)

		assert.Equal(t, taskRunID.UUID(), etrt.TaskRunID)
		require.NotNil(t, etrt.EthTx)
		assert.Nil(t, etrt.EthTx.Nonce)
		assert.Equal(t, toAddress, etrt.EthTx.ToAddress)
		assert.Equal(t, "70a08231888888880000000000000000000000000000000000000000000000000000009786856756", hex.EncodeToString(etrt.EthTx.EncodedPayload))
		assert.Equal(t, gasLimit, etrt.EthTx.GasLimit)
		assert.Equal(t, models.EthTxUnstarted, etrt.EthTx.State)
	})

	t.Run("if FromAddresses is provided but no key matches, returns job error", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			FromAddresses:    []gethCommon.Address{cltest.NewAddress()},
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.EqualError(t, runOutput.Error(), "insertEthTx failed to pickFromAddress: no keys available")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())
	})

	t.Run("with bytes DataFormat writes correct encoded data to database", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
			DataFormat:       "bytes",
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "cönfirmed", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())

		expectedData := hexutil.MustDecode(
			functionSelector.String() +
				"88888888" + // dataPrefix
				"0000000000000000000000000000000000000000000000000000000000000040" + // offset
				"000000000000000000000000000000000000000000000000000000000000000a" + // length in bytes
				"63c3b66e6669726d656400000000000000000000000000000000000000000000") // encoded string left padded

		etrt, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)

		assert.Equal(t, expectedData, etrt.EthTx.EncodedPayload)
	})

	t.Run("with invalid data returns run output error and does not write to DB", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
			DataFormat:       "some old bollocks",
		}
		jobRunID := models.NewID()
		taskRunID := models.NewID()
		input := models.NewRunInputWithResult(jobRunID, *taskRunID, "0x9786856756", models.RunStatusUnstarted)
		runOutput := adapter.Perform(*input, store)
		assert.Contains(t, runOutput.Error().Error(), "while constructing EthTx data: unsupported format: some old bollocks")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())

		trtx, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
		require.NoError(t, err)
		require.Nil(t, trtx)
	})

	t.Run("with unconfirmed transaction returns output pending confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0)
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction returns pending outgoing confirmations if receipt is missing (invariant violation, should never happen)", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 1, 1)
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is younger than minRequiredOutgoingConfirmations, returns output pending_outgoing_confirmations", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 2, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 12,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, runOutput.Status())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is equal to minRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 3, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 13,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		data := cltest.JSONFromString(t, `{"foo": "bar", "result": "some old bollocks"}`)
		input := models.NewRunInput(jobRunID, taskRunID, data, models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
		// Does not clobber previously assigned data
		assert.Equal(t, "bar", runOutput.Get("foo").String())
		// Assigns latestOutgoingTxHash for legacy compatibility
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Get("latestOutgoingTxHash").String())
	})

	t.Run("with confirmed transaction with exactly one attempt with exactly one receipt that is older than minRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:                        toAddress,
			GasLimit:                         gasLimit,
			FunctionSelector:                 functionSelector,
			DataPrefix:                       dataPrefix,
			MinRequiredOutgoingConfirmations: 12,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 4, 1)

		confirmedAttemptHash := etx.EthTxAttempts[0].Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: 14,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
	})

	t.Run("with confirmed transaction with two attempts, one of which has exactly one receipt that is older than the default MinRequiredOutgoingConfirmations, returns output complete with transaction hash pulled from receipt", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 5, 1)
		attempt2 := cltest.MustInsertBroadcastEthTxAttempt(t, etx.ID, store, 2)

		confirmedAttemptHash := attempt2.Hash

		cltest.MustInsertEthReceipt(t, store, 1, cltest.NewHash(), confirmedAttemptHash)
		require.NoError(t, store.IdempotentInsertHead(models.Head{
			Hash:   cltest.NewHash(),
			Number: int64(store.Config.MinRequiredOutgoingConfirmations()) + 2,
		}))
		store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID)
		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.NoError(t, runOutput.Error())
		assert.Equal(t, models.RunStatusCompleted, runOutput.Status())
		assert.Equal(t, confirmedAttemptHash.Hex(), runOutput.Result().String())
	})

	t.Run("with transaction that ended up in fatal_error state returns job run error", func(t *testing.T) {
		adapter := adapters.EthTx{
			ToAddress:        toAddress,
			GasLimit:         gasLimit,
			FunctionSelector: functionSelector,
			DataPrefix:       dataPrefix,
		}
		jobRunID := models.NewID()
		taskRunID := cltest.MustInsertTaskRun(t, store)
		etx := cltest.MustInsertFatalErrorEthTx(t, store)
		require.NoError(t, store.DB.Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID.UUID(), etx.ID).Error)

		input := models.NewRunInputWithResult(jobRunID, taskRunID, "0x9786856756", models.RunStatusUnstarted)

		// Do the thing
		runOutput := adapter.Perform(*input, store)

		require.EqualError(t, runOutput.Error(), "something exploded")
		assert.Equal(t, models.RunStatusErrored, runOutput.Status())
		assert.Equal(t, "", runOutput.Result().String())
	})
}
