package adapters

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	// DataFormatBytes instructs the EthTx Adapter to treat the input value as a
	// bytes string, rather than a hexadecimal encoded bytes32
	DataFormatBytes = "bytes"
)

// EthTx holds the Address to send the result to and the FunctionSelector
// to execute.
type EthTx struct {
	ToAddress common.Address `json:"address"`
	// NOTE: FromAddress is deprecated and kept for backwards compatibility, new job specs should use fromAddresses
	FromAddress      common.Address          `json:"fromAddress,omitempty"`
	FromAddresses    []common.Address        `json:"fromAddresses,omitempty"`
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

func (e *EthTx) checkForConfirmation(trtx models.EthTaskRunTx,
	input models.RunInput, store *strpkg.Store) models.RunOutput {
	switch trtx.EthTx.State {
	case models.EthTxConfirmed:
		return e.checkEthTxForReceipt(trtx.EthTx.ID, input, store)
	case models.EthTxFatalError:
		return models.NewRunOutputError(trtx.EthTx.GetError())
	default:
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}
}

func (e *EthTx) pickFromAddress(input models.RunInput, store *strpkg.Store) (common.Address, error) {
	if len(e.FromAddresses) > 0 {
		if e.FromAddress != utils.ZeroAddress {
			logger.Warnf("task spec for task run %s specified both fromAddress and fromAddresses."+
				" fromAddress is deprecated, it will be ignored and fromAddresses used instead. "+
				"Specifying both of these keys in a job spec may result in an error in future versions of Chainlink", input.TaskRunID())
		}
		return store.GetRoundRobinAddress(e.FromAddresses...)
	}
	if e.FromAddress == utils.ZeroAddress {
		return store.GetRoundRobinAddress()
	}
	logger.Warnf(`DEPRECATION WARNING: task spec for task run %s specified a fromAddress of %s. fromAddress has been deprecated and will be removed in a future version of Chainlink. Please use fromAddresses instead. You can pin a job to one address simply by using only one element, like so:
{
	"type": "EthTx",
	"fromAddresses": ["%s"],
}
`, input.TaskRunID(), e.FromAddress.Hex(), e.FromAddress.Hex())
	return e.FromAddress, nil
}

func (e *EthTx) insertEthTx(input models.RunInput, store *strpkg.Store) models.RunOutput {
	txData, err := getTxData(e, input)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	taskRunID := input.TaskRunID()
	toAddress := e.ToAddress
	fromAddress, err := e.pickFromAddress(input, store)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed to pickFromAddress")
		logger.Error(err)
		return models.NewRunOutputError(err)
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

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
}

func (e *EthTx) checkEthTxForReceipt(ethTxID int64, input models.RunInput, s *strpkg.Store) models.RunOutput {
	var minRequiredOutgoingConfirmations uint64
	if e.MinRequiredOutgoingConfirmations == 0 {
		minRequiredOutgoingConfirmations = s.Config.MinRequiredOutgoingConfirmations()
	} else {
		minRequiredOutgoingConfirmations = e.MinRequiredOutgoingConfirmations
	}

	hash, err := getConfirmedTxHash(ethTxID, s.DB, minRequiredOutgoingConfirmations)

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

// getTxData returns the data to save against the callback encoded according to
// the dataFormat parameter in the job spec
func getTxData(e *EthTx, input models.RunInput) ([]byte, error) {
	result := input.Result()
	if e.DataFormat == "" {
		return common.HexToHash(result.Str).Bytes(), nil
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
