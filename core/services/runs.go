package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func markInProgressIfSufficientIncomingConfirmations(run *models.JobRun, taskRun *models.TaskRun, currentHeight *utils.Big, ethClient eth.Client) {
	updateTaskRunObservedIncomingConfirmations(currentHeight, run, taskRun)

	if !meetsMinRequiredIncomingConfirmations(run, taskRun, run.ObservedHeight) {
		logger.Debugw("Pausing run pending confirmations",
			run.ForLogger("required_height", taskRun.MinRequiredIncomingConfirmations)...,
		)

		taskRun.Status = models.RunStatusPendingIncomingConfirmations
		run.SetStatus(models.RunStatusPendingIncomingConfirmations)

	} else if err := validateOnMainChain(run, taskRun, ethClient); err != nil {
		logger.Warnw("Failure while trying to validate chain",
			run.ForLogger("error", err)...,
		)

		taskRun.SetError(err)
		run.SetError(err)

	} else {
		run.SetStatus(models.RunStatusInProgress)
	}
}

func validateOnMainChain(run *models.JobRun, taskRun *models.TaskRun, ethClient eth.Client) error {
	txhash := run.RunRequest.TxHash
	if txhash == nil || !taskRun.MinRequiredIncomingConfirmations.Valid || taskRun.MinRequiredIncomingConfirmations.Uint32 == 0 {
		return nil
	}

	receipt, err := ethClient.TransactionReceipt(context.TODO(), *txhash)
	if err != nil {
		return err
	}
	if invalidRequest(run.RunRequest, receipt) {
		return fmt.Errorf(
			"TxHash %s initiating run %s not on main chain; presumably has been uncled",
			txhash.Hex(),
			run.ID.String(),
		)
	}
	return nil
}

func updateTaskRunObservedIncomingConfirmations(currentHeight *utils.Big, jr *models.JobRun, taskRun *models.TaskRun) {
	if !taskRun.MinRequiredIncomingConfirmations.Valid || jr.CreationHeight == nil || currentHeight == nil {
		return
	}

	confs := blockConfirmations(currentHeight, jr.CreationHeight)
	diff := utils.MinBigs(confs, big.NewInt(int64(taskRun.MinRequiredIncomingConfirmations.Uint32)))

	// diff's ceiling is guaranteed to be MaxUint32 since MinRequiredIncomingConfirmations
	// ceiling is MaxUint32.
	taskRun.ObservedIncomingConfirmations = clnull.Uint32From(uint32(diff.Int64()))
}

func invalidRequest(request models.RunRequest, receipt *types.Receipt) bool {
	return models.ReceiptIsUnconfirmed(receipt) ||
		(request.BlockHash != nil && *request.BlockHash != receipt.BlockHash)
}

func meetsMinRequiredIncomingConfirmations(
	run *models.JobRun,
	taskRun *models.TaskRun,
	currentHeight *utils.Big) bool {

	if !taskRun.MinRequiredIncomingConfirmations.Valid || run.CreationHeight == nil || currentHeight == nil {
		return true
	}

	diff := blockConfirmations(currentHeight, run.CreationHeight)
	return diff.Cmp(big.NewInt(int64(taskRun.MinRequiredIncomingConfirmations.Uint32))) >= 0
}

func blockConfirmations(currentHeight, creationHeight *utils.Big) *big.Int {
	bigDiff := new(big.Int).Sub(currentHeight.ToInt(), creationHeight.ToInt())
	confs := bigDiff.Add(bigDiff, big.NewInt(1)) // creation of runlog alone warrants 1 confirmation
	if confs.Cmp(big.NewInt(0)) < 0 {            // negative, so floor at 0
		confs.SetUint64(0)
	}
	return confs
}
