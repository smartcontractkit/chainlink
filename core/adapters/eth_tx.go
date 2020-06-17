package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

const (
	// DataFormatBytes instructs the EthTx Adapter to treat the input value as a
	// bytes string, rather than a hexadecimal encoded bytes32
	DataFormatBytes = "bytes"
)

// EthTx holds the Address to send the result to and the FunctionSelector
// to execute.
type EthTx struct {
	ToAddress        common.Address          `json:"address"`
	FromAddress      common.Address          `json:"fromAddress,omitempty"`
	FunctionSelector models.FunctionSelector `json:"functionSelector"`
	DataPrefix       hexutil.Bytes           `json:"dataPrefix"`
	DataFormat       string                  `json:"format"`
	GasLimit         uint64                  `json:"gasLimit,omitempty"`

	// GasPrice only needed for legacy tx manager
	GasPrice *utils.Big `json:"gasPrice" gorm:"type:numeric"`

	// MinRequiredOutgoingConfirmations only works with bulletprooftxmanager
	MinRequiredOutgoingConfirmations uint64 `json:"minRequiredOutgoingConfirmations,omitempty"`
}

// TaskType returns the type of Adapter.
func (e *EthTx) TaskType() models.TaskType {
	return TaskTypeEthTx
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (e *EthTx) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	if store.Config.EnableBulletproofTxManager() {
		return e.perform(input, store)
	}
	return e.legacyPerform(input, store)
}

// TODO(sam): https://www.pivotaltracker.com/story/show/173280188
func (e *EthTx) perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	trtx, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
	if err != nil {
		err = errors.Wrap(err, "FindEthTaskRunTxByTaskRunID failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}
	if trtx != nil {
		return e.checkForConfirmation(*trtx, input, store)
	}
	return e.insertEthTx(input, store)
}

func (e *EthTx) checkForConfirmation(trtx models.EthTaskRunTx, input models.RunInput, store *store.Store) models.RunOutput {
	switch trtx.EthTx.State {
	case models.EthTxConfirmed:
		return e.checkEthTxForReceipt(trtx.EthTx.ID, input, store)
	case models.EthTxFatalError:
		return models.NewRunOutputError(trtx.EthTx.GetError())
	default:
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}
}

func (e *EthTx) insertEthTx(input models.RunInput, store *store.Store) models.RunOutput {
	txData, err := getTxData(e, input)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	taskRunID := input.TaskRunID()
	toAddress := e.ToAddress
	var fromAddress common.Address
	if e.FromAddress == utils.ZeroAddress {
		fromAddress, err = store.GetDefaultAddress()
		if err != nil {
			err = errors.Wrap(err, "insertEthTx failed to GetDefaultAddress")
			logger.Error(err)
			return models.NewRunOutputError(err)
		}
	} else {
		fromAddress = e.FromAddress
	}
	encodedPayload := utils.ConcatBytes(e.FunctionSelector.Bytes(), e.DataPrefix, txData)

	var gasLimit uint64
	if e.GasLimit == 0 {
		gasLimit = store.Config.EthGasLimitDefault()
	} else {
		gasLimit = e.GasLimit
	}

	if err := store.IdempotentInsertEthTaskRunTx(taskRunID, fromAddress, toAddress, encodedPayload, gasLimit); err != nil {
		err = errors.Wrap(err, "insertEthTx failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	store.NotifyNewEthTx.Trigger()

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
}

func (e *EthTx) checkEthTxForReceipt(ethTxID int64, input models.RunInput, s *store.Store) models.RunOutput {
	var minRequiredOutgoingConfirmations uint64
	if e.MinRequiredOutgoingConfirmations == 0 {
		minRequiredOutgoingConfirmations = s.Config.MinRequiredOutgoingConfirmations()
	} else {
		minRequiredOutgoingConfirmations = e.MinRequiredOutgoingConfirmations
	}

	hash, err := getConfirmedTxHash(ethTxID, s.GetRawDB(), minRequiredOutgoingConfirmations)

	if err != nil {
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	if hash == nil {
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}

	hexHash := (*hash).Hex()

	output := input.Data()
	output, err = output.MultiAdd(models.KV{
		"result": hexHash,
		// HACK: latestOutgoingTxHash is used for backwards compatibility with the stats pusher
		"latestOutgoingTxHash": hexHash,
	})
	if err != nil {
		err = errors.Wrap(err, "checkEthTxForReceipt failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}
	return models.NewRunOutputComplete(output)
}

func getConfirmedTxHash(ethTxID int64, db *gorm.DB, minRequiredOutgoingConfirmations uint64) (*common.Hash, error) {
	receipt := models.EthReceipt{}
	err := db.
		Joins("INNER JOIN eth_tx_attempts ON eth_tx_attempts.hash = eth_receipts.tx_hash AND eth_tx_attempts.eth_tx_id = ?", ethTxID).
		Joins("INNER JOIN eth_txes ON eth_txes.id = eth_tx_attempts.eth_tx_id AND eth_txes.state = 'confirmed'").
		Where("eth_receipts.block_number <= (SELECT max(number) - ? FROM heads)", minRequiredOutgoingConfirmations).
		First(&receipt).
		Error

	if err == nil {
		return &receipt.TxHash, nil
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return nil, errors.Wrap(err, "getConfirmedTxHash failed")

}

func (e *EthTx) legacyPerform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	if !store.TxManager.Connected() {
		return pendingOutgoingConfirmationsOrConnection(input)
	}

	if input.Status().PendingOutgoingConfirmations() {
		return ensureTxRunResult(input, store)
	}

	value, err := getTxData(e, input)
	if err != nil {
		err = errors.Wrap(err, "while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	data := utils.ConcatBytes(e.FunctionSelector.Bytes(), e.DataPrefix, value)
	return createTxRunResult(e.ToAddress, e.GasPrice, e.GasLimit, data, input, store)
}

// getTxData returns the data to save against the callback encoded according to
// the dataFormat parameter in the job spec
func getTxData(e *EthTx, input models.RunInput) ([]byte, error) {
	result := input.Result()
	if e.DataFormat == "" {
		return common.HexToHash(result.Str).Bytes(), nil
	}

	payloadOffset := utils.EVMWordUint64(utils.EVMWordByteLen)
	if len(e.DataPrefix) > 0 {
		payloadOffset = utils.EVMWordUint64(utils.EVMWordByteLen * 2)
	}
	output, err := utils.EVMTranscodeJSONWithFormat(result, e.DataFormat)
	if err != nil {
		return []byte{}, err
	}
	return utils.ConcatBytes(payloadOffset, output), nil
}

func createTxRunResult(
	address common.Address,
	gasPrice *utils.Big,
	gasLimit uint64,
	data []byte,
	input models.RunInput,
	store *strpkg.Store,
) models.RunOutput {
	tx, err := store.TxManager.CreateTxWithGas(
		null.StringFrom(input.JobRunID().String()),
		address,
		data,
		gasPrice.ToInt(),
		gasLimit,
	)
	if err != nil {
		logger.Error(errors.Wrap(err, "createTxRunResult failed"))
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}

	output, err := models.JSON{}.Add("result", tx.Hash.String())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	txAttempt := tx.Attempts[0]
	receipt, state, err := store.TxManager.CheckAttempt(txAttempt, tx.SentAt)
	if err != nil {
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(output)
	}

	logger.Debugw(
		fmt.Sprintf("Tx #0 is %s", state),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
		"receiptBlockNumber", receipt.BlockNumber.ToInt(),
		"currentBlockNumber", tx.SentAt,
		"receiptHash", receipt.Hash.Hex(),
	)

	if state == strpkg.Safe {
		// I don't see how the receipt could possibly be nil here, but handle it just in case
		if receipt == nil {
			err := errors.New("missing receipt for transaction")
			return models.NewRunOutputError(err)
		}
		return addReceiptToResult(*receipt, input, output)
	}

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(output)
}

func ensureTxRunResult(input models.RunInput, str *strpkg.Store) models.RunOutput {
	val, err := input.ResultString()
	if err != nil {
		return models.NewRunOutputError(errors.Wrapf(err, "while processing ethtx input %#+v", input))
	}

	hash := common.HexToHash(val)
	receipt, state, err := str.TxManager.BumpGasUntilSafe(hash)
	if err != nil {
		// We failed to get one of the TxAttempt receipts, so we won't mark this
		// run as errored in order to try again
		logger.Warn("EthTx Adapter Perform Resuming: ", err)
	}

	var output models.JSON

	if receipt != nil && !receipt.Unconfirmed() {
		// If the tx has been confirmed, record the hash in the output
		hex := receipt.Hash.String()
		output, err = output.Add("result", hex)
		if err != nil {
			return models.NewRunOutputError(err)
		}
		output, err = output.Add("latestOutgoingTxHash", hex)
		if err != nil {
			return models.NewRunOutputError(err)
		}
	} else {
		// If the tx is still unconfirmed, just copy over the original tx hash.
		output, err = output.Add("result", hash)
		if err != nil {
			return models.NewRunOutputError(err)
		}
	}

	if state == strpkg.Safe {
		// FIXME: Receipt can definitely be nil here, although I don't really know how
		// it can be "Safe" without a receipt... maybe we should just keep
		// waiting for confirmations instead?
		if receipt == nil {
			err := errors.New("missing receipt for transaction")
			return models.NewRunOutputError(err)
		}

		return addReceiptToResult(*receipt, input, output)
	}

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(output)
}

func addReceiptToResult(
	receipt models.TxReceipt,
	input models.RunInput,
	data models.JSON,
) models.RunOutput {
	receipts := []models.TxReceipt{}

	ethereumReceipts := input.Data().Get("ethereumReceipts").String()
	if ethereumReceipts != "" {
		if err := json.Unmarshal([]byte(ethereumReceipts), &receipts); err != nil {
			logger.Errorw("Error unmarshaling ethereum Receipts", "error", err)
		}
	}

	receipts = append(receipts, receipt)
	var err error
	data, err = data.Add("ethereumReceipts", receipts)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	data, err = data.Add("result", receipt.Hash.String())
	if err != nil {
		return models.NewRunOutputError(err)
	}
	return models.NewRunOutputComplete(data)
}

func pendingOutgoingConfirmationsOrConnection(input models.RunInput) models.RunOutput {
	// If the input is not pending outgoing confirmations next time
	// then it may submit a new transaction.
	if input.Status().PendingOutgoingConfirmations() {
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}
	return models.NewRunOutputPendingConnection()
}
