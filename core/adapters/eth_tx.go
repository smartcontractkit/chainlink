package adapters

import (
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	ToAddress gethCommon.Address `json:"address"`
	// NOTE: FromAddress is deprecated and kept for backwards compatibility, new job specs should use fromAddresses
	FromAddress      gethCommon.Address      `json:"fromAddress,omitempty"`
	FromAddresses    []gethCommon.Address    `json:"fromAddresses,omitempty"`
	FunctionSelector models.FunctionSelector `json:"functionSelector"`
	DataPrefix       hexutil.Bytes           `json:"dataPrefix"`
	DataFormat       string                  `json:"format"`
	GasLimit         uint64                  `json:"gasLimit,omitempty"`

	// GasPrice only needed for legacy tx manager
	GasPrice *utils.Big `json:"gasPrice" gorm:"type:numeric"`

	// MinRequiredOutgoingConfirmations only works with bulletprooftxmanager
	MinRequiredOutgoingConfirmations uint64 `json:"minRequiredOutgoingConfirmations,omitempty"`
}

func (e *EthTx) GetToAddress() gethCommon.Address {
	return e.ToAddress
}

func (e *EthTx) GetFromAddresses() []gethCommon.Address {
	if len(e.FromAddresses) > 0 {
		if e.FromAddress != utils.ZeroAddress {
			logger.Warn("task spec for task run specified both fromAddress and fromAddresses." +
				" fromAddress is deprecated, it will be ignored and fromAddresses used instead. " +
				"Specifying both of these keys in a job spec may result in an error in future versions of Chainlink")
		}
		return e.FromAddresses
	}

	if e.FromAddress == utils.ZeroAddress {
		return []gethCommon.Address{}
	}
	logger.Warnf(`DEPRECATION WARNING: task spec specified a fromAddress of %s. fromAddress has been deprecated and will be removed in a future version of Chainlink. Please use fromAddresses instead. You can pin a job to one address simply by using only one element, like so:
{
	"type": "EthTx",
	"fromAddresses": ["%s"],
} 
`, e.FromAddress.Hex(), e.FromAddress.Hex())
	return []gethCommon.Address{e.FromAddress}
}

func (e *EthTx) GetGasLimit() uint64 {
	return e.GasLimit
}

func (e *EthTx) GetGasPrice() *utils.Big {
	return e.GasPrice
}

func (e *EthTx) GetMinRequiredOutgoingConfirmations() uint64 {
	return e.MinRequiredOutgoingConfirmations
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
	return findOrInsertEthTx(e, input, store, func() ([]byte, error) {
		txData, err := getTxData(e, input)
		if err != nil {
			return nil, err
		}
		return utils.ConcatBytes(e.FunctionSelector.Bytes(), e.DataPrefix, txData), nil
	})
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
		return gethCommon.HexToHash(result.Str).Bytes(), nil
	}

	output, err := utils.EVMTranscodeJSONWithFormat(result, e.DataFormat)
	if err != nil {
		return []byte{}, err
	}
	if e.DataFormat == DataFormatBytes || len(e.DataPrefix) > 0 {
		payloadOffset := utils.EVMWordUint64(utils.EVMWordByteLen)
		if len(e.DataPrefix) > 0 {
			payloadOffset = utils.EVMWordUint64(utils.EVMWordByteLen * 2)
			return utils.ConcatBytes(payloadOffset, output), nil
		}
		return utils.ConcatBytes(payloadOffset, output), nil
	}
	return utils.ConcatBytes(output), nil
}

func createTxRunResult(
	address gethCommon.Address,
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

	hash := gethCommon.HexToHash(val)
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
