package adapters

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	Address          common.Address          `json:"address"`
	FunctionSelector models.FunctionSelector `json:"functionSelector"`
	DataPrefix       hexutil.Bytes           `json:"dataPrefix"`
	DataFormat       string                  `json:"format"`
	GasPrice         *models.Big             `json:"gasPrice" gorm:"type:numeric"`
	GasLimit         uint64                  `json:"gasLimit"`
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (etx *EthTx) Perform(input models.JSON, result models.RunResult, store *strpkg.Store) models.RunResult {
	if !store.TxManager.Connected() {
		result.MarkPendingConnection()
		return result
	}

	if !result.Status.PendingConfirmations() {
		createTxRunResult(etx, input, &result, store)
		return result
	}
	ensureTxRunResult(&result, store)
	return result
}

// getTxData returns the data to save against the callback encoded according to
// the dataFormat parameter in the job spec
func getTxData(e *EthTx, input models.JSON) ([]byte, error) {
	result := input.Get("result")
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
	e *EthTx,
	input models.JSON,
	result *models.RunResult,
	store *strpkg.Store,
) {
	value, err := getTxData(e, input)
	if err != nil {
		result.SetError(err)
		return
	}

	data := utils.ConcatBytes(e.FunctionSelector.Bytes(), e.DataPrefix, value)

	tx, err := store.TxManager.CreateTxWithGas(
		null.StringFrom(result.CachedJobRunID),
		e.Address,
		data,
		e.GasPrice.ToInt(),
		e.GasLimit,
	)
	if IsClientRetriable(err) {
		result.MarkPendingConnection()
		return
	} else if err != nil {
		result.SetError(err)
		return
	}

	result.ApplyResult(tx.Hash.String())

	txAttempt := tx.Attempts[0]
	logger.Debugw(
		fmt.Sprintf("Tx #0 checking on-chain state"),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
	)

	receipt, state, err := store.TxManager.CheckAttempt(txAttempt, tx.SentAt)
	if IsClientRetriable(err) {
		result.MarkPendingConnection()
		return
	} else if IsClientEmptyError(err) {
		result.MarkPendingConfirmations()
		return
	} else if err != nil {
		result.SetError(err)
		return
	}

	logger.Debugw(
		fmt.Sprintf("Tx #0 is %s", state),
		"txHash", txAttempt.Hash.String(),
		"txID", txAttempt.TxID,
		"receiptBlockNumber", receipt.BlockNumber.ToInt(),
		"currentBlockNumber", tx.SentAt,
		"receiptHash", receipt.Hash.Hex(),
	)

	if state != strpkg.Safe {
		result.MarkPendingConfirmations()
		return
	}

	addReceiptToResult(receipt, result)
}

func ensureTxRunResult(result *models.RunResult, str *strpkg.Store) {
	val, err := result.ResultString()
	if err != nil {
		result.SetError(err)
		return
	}

	hash := common.HexToHash(val)
	if err != nil {
		result.SetError(err)
		return
	}

	receipt, state, err := str.TxManager.BumpGasUntilSafe(hash)
	if err != nil {
		if IsClientEmptyError(err) {
			result.MarkPendingConfirmations()
			return
		} else if state == strpkg.Unknown {
			result.SetError(err)
			return
		}

		// We failed to get one of the TxAttempt receipts, so we won't mark this
		// run as errored in order to try again
		logger.Warn("EthTx Adapter Perform Resuming: ", err)
	}

	recordLatestTxHash(receipt, result)
	if state != strpkg.Safe {
		result.MarkPendingConfirmations()
		return
	}

	addReceiptToResult(receipt, result)
}

var zero = common.Hash{}

// recordLatestTxHash adds the current tx hash to the run result
func recordLatestTxHash(receipt *models.TxReceipt, result *models.RunResult) {
	if receipt == nil || receipt.Unconfirmed() {
		return
	}
	hex := receipt.Hash.String()
	result.ApplyResult(hex)
	result.Add("latestOutgoingTxHash", hex)
}

func addReceiptToResult(receipt *models.TxReceipt, result *models.RunResult) {
	receipts := []models.TxReceipt{}

	if !result.Get("ethereumReceipts").IsArray() {
		result.Add("ethereumReceipts", receipts)
	}

	if err := json.Unmarshal([]byte(result.Get("ethereumReceipts").String()), &receipts); err != nil {
		logger.Error(fmt.Errorf("EthTx Adapter unmarshaling ethereum Receipts: %v", err))
	}

	receipts = append(receipts, *receipt)
	result.Add("ethereumReceipts", receipts)
	result.CompleteWithResult(receipt.Hash.String())
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
