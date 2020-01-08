package adapters

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"

	"chainlink/core/eth"
	"chainlink/core/logger"
	"chainlink/core/store"
	strpkg "chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
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
	Address          common.Address       `json:"address"`
	FunctionSelector eth.FunctionSelector `json:"functionSelector"`
	DataPrefix       hexutil.Bytes        `json:"dataPrefix"`
	DataFormat       string               `json:"format"`
	GasPrice         *utils.Big           `json:"gasPrice" gorm:"type:numeric"`
	GasLimit         uint64               `json:"gasLimit"`
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (etx *EthTx) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	logger.Info("ethTX.Perform",
		"runInput.ID", input.JobRunID().String(),
		"runInput.Data", input.Data().String(),
		"runInput.Status", input.Status(),
	)

	if !store.TxManager.Connected() {
		return pendingConfirmationsOrConnection(input)
	}

	if input.Status().PendingConfirmations() {
		return ensureTxRunResult(input, store)
	}

	value, err := getTxData(etx, input)
	if err != nil {
		err = errors.Wrap(err, "while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	data := utils.ConcatBytes(etx.FunctionSelector.Bytes(), etx.DataPrefix, value)
	ret := createTxRunResult(etx.Address, etx.GasPrice, etx.GasLimit, data, input, store)
	logger.Info("ethTX.Perform",
		"runOutput.Data", ret.Data().String(),
		"runOutput.Status", ret.Status(),
		"runOutput.Error", ret.Error(),
	)
	return ret
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
	if IsClientRetriable(err) {
		return models.NewRunOutputPendingConnection()
	} else if err != nil {
		return models.NewRunOutputError(err)
	}

	output, err := models.JSON{}.Add("result", tx.Hash.String())
	if err != nil {
		return models.NewRunOutputError(err)
	}

	txAttempt := tx.Attempts[0]
	receipt, state, err := store.TxManager.CheckAttempt(txAttempt, tx.SentAt)
	if IsClientRetriable(err) {
		return models.NewRunOutputPendingConnectionWithData(output)
	} else if IsClientEmptyError(err) {
		return models.NewRunOutputPendingConfirmationsWithData(output)
	} else if err != nil {
		return models.NewRunOutputError(err)
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
		return addReceiptToResult(receipt, input, output)
	}

	return models.NewRunOutputPendingConfirmationsWithData(output)
}

func ensureTxRunResult(input models.RunInput, str *strpkg.Store) models.RunOutput {
	val, err := input.ResultString()
	if err != nil {
		return models.NewRunOutputError(err)
	}

	hash := common.HexToHash(val)
	receipt, state, err := str.TxManager.BumpGasUntilSafe(hash)
	if err != nil {
		if IsClientEmptyError(err) {
			return models.NewRunOutputPendingConfirmations()
		} else if state == strpkg.Unknown {
			return models.NewRunOutputError(err)
		}

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
		return addReceiptToResult(receipt, input, output)
	}

	return models.NewRunOutputPendingConfirmationsWithData(output)
}

func addReceiptToResult(
	receipt *eth.TxReceipt,
	input models.RunInput,
	data models.JSON,
) models.RunOutput {
	receipts := []eth.TxReceipt{}

	ethereumReceipts := input.Data().Get("ethereumReceipts").String()
	if ethereumReceipts != "" {
		if err := json.Unmarshal([]byte(ethereumReceipts), &receipts); err != nil {
			logger.Errorw("Error unmarshaling ethereum Receipts", "error", err)
		}
	}

	receipts = append(receipts, *receipt)
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

// IsClientRetriable does its best effort to see if an error indicates one that
// might have a different outcome if we retried the operation
func IsClientRetriable(err error) bool {
	if err == nil {
		return false
	}

	if err, ok := err.(net.Error); ok {
		return err.Timeout() || err.Temporary()
	} else if errors.Cause(err) == store.ErrPendingConnection {
		return true
	}

	return false
}

var (
	parityEmptyResponseRegex = regexp.MustCompile("Error cause was EmptyResponse")
)

// Parity light clients can return an EmptyResponse error when they don't have
// access to the transaction in the mempool. If we wait long enough it should
// eventually return a transaction receipt.
func IsClientEmptyError(err error) bool {
	return err != nil && parityEmptyResponseRegex.MatchString(err.Error())
}

func pendingConfirmationsOrConnection(input models.RunInput) models.RunOutput {
	// If the input is not pending confirmations next time it may, then it may
	// submit a new transaction.
	if input.Status().PendingConfirmations() {
		return models.NewRunOutputPendingConfirmations()
	}
	return models.NewRunOutputPendingConnection()
}
