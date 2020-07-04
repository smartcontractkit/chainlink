package adapters

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// EthTxCommon represents a common interface for all EthTx adapters
type EthTxCommon interface {
	GetToAddress() gethCommon.Address
	GetFromAddresses() []gethCommon.Address
	GetGasLimit() uint64
	// GasPrice only needed for legacy tx manager
	GetGasPrice() *utils.Big
	// MinRequiredOutgoingConfirmations only works with bulletprooftxmanager
	GetMinRequiredOutgoingConfirmations() uint64

	GetEncodedPayload(input models.RunInput) ([]byte, error)
}

// TODO(sam): https://www.pivotaltracker.com/story/show/173280188
func findOrInsertEthTx(e EthTxCommon, input models.RunInput, store *strpkg.Store) models.RunOutput {
	trtx, err := store.FindEthTaskRunTxByTaskRunID(input.TaskRunID().UUID())
	if err != nil {
		err = errors.Wrap(err, "FindEthTaskRunTxByTaskRunID failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}
	if trtx != nil {
		return checkForConfirmation(e, *trtx, input, store)
	}
	return insertEthTx(e, input, store)
}

func checkForConfirmation(e EthTxCommon, trtx models.EthTaskRunTx, input models.RunInput, store *strpkg.Store) models.RunOutput {
	switch trtx.EthTx.State {
	case models.EthTxConfirmed:
		return checkEthTxForReceipt(e, trtx.EthTx.ID, input, store)
	case models.EthTxFatalError:
		return models.NewRunOutputError(trtx.EthTx.GetError())
	default:
		return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
	}
}

func insertEthTx(e EthTxCommon, input models.RunInput, store *strpkg.Store) models.RunOutput {
	encodedPayload, err := e.GetEncodedPayload(input)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed while constructing EthTx data")
		return models.NewRunOutputError(err)
	}

	taskRunID := input.TaskRunID()
	toAddress := e.GetToAddress()
	fromAddress, err := store.GetRoundRobinAddress(e.GetFromAddresses()...)
	if err != nil {
		err = errors.Wrap(err, "insertEthTx failed to GetRoundRobinAddress")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	var gasLimit uint64
	if e.GetGasLimit() == 0 {
		gasLimit = store.Config.EthGasLimitDefault()
	} else {
		gasLimit = e.GetGasLimit()
	}

	if err := store.IdempotentInsertEthTaskRunTx(taskRunID, fromAddress, toAddress, encodedPayload, gasLimit); err != nil {
		err = errors.Wrap(err, "insertEthTx failed")
		logger.Error(err)
		return models.NewRunOutputError(err)
	}

	store.NotifyNewEthTx.Trigger()

	return models.NewRunOutputPendingOutgoingConfirmationsWithData(input.Data())
}

func checkEthTxForReceipt(e EthTxCommon, ethTxID int64, input models.RunInput, s *strpkg.Store) models.RunOutput {
	var minRequiredOutgoingConfirmations uint64
	if e.GetMinRequiredOutgoingConfirmations() == 0 {
		minRequiredOutgoingConfirmations = s.Config.MinRequiredOutgoingConfirmations()
	} else {
		minRequiredOutgoingConfirmations = e.GetMinRequiredOutgoingConfirmations()
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

func getConfirmedTxHash(ethTxID int64, db *gorm.DB, minRequiredOutgoingConfirmations uint64) (*gethCommon.Hash, error) {
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
