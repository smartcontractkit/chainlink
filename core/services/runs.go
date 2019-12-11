package services

import (
	"fmt"
	"math/big"

	"chainlink/core/eth"
	"chainlink/core/logger"
	clnull "chainlink/core/null"
	"chainlink/core/store"
	"chainlink/core/store/models"
	"chainlink/core/utils"
)

func validateMinimumConfirmations(run *models.JobRun, taskRun *models.TaskRun, currentHeight *utils.Big, txManager store.TxManager) {
	updateTaskRunConfirmations(currentHeight, run, taskRun)

	if !meetsMinimumConfirmations(run, taskRun, run.ObservedHeight) {
		logger.Debugw("Pausing run pending confirmations",
			run.ForLogger("required_height", taskRun.MinimumConfirmations)...,
		)

		taskRun.Status = models.RunStatusPendingConfirmations
		run.Status = models.RunStatusPendingConfirmations

	} else if err := validateOnMainChain(run, taskRun, txManager); err != nil {
		logger.Warnw("Failure while trying to validate chain",
			run.ForLogger("error", err)...,
		)

		taskRun.SetError(err)
		run.SetError(err)

	} else {
		run.Status = models.RunStatusInProgress
	}
}

func validateOnMainChain(run *models.JobRun, taskRun *models.TaskRun, txManager store.TxManager) error {
	txhash := run.RunRequest.TxHash
	if txhash == nil || !taskRun.MinimumConfirmations.Valid || taskRun.MinimumConfirmations.Uint32 == 0 {
		return nil
	}

	receipt, err := txManager.GetTxReceipt(*txhash)
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

func updateTaskRunConfirmations(currentHeight *utils.Big, jr *models.JobRun, taskRun *models.TaskRun) {
	if !taskRun.MinimumConfirmations.Valid || jr.CreationHeight == nil || currentHeight == nil {
		return
	}

	confs := blockConfirmations(currentHeight, jr.CreationHeight)
	diff := utils.MinBigs(confs, big.NewInt(int64(taskRun.MinimumConfirmations.Uint32)))

	// diff's ceiling is guaranteed to be MaxUint32 since MinimumConfirmations
	// ceiling is MaxUint32.
	taskRun.Confirmations = clnull.Uint32From(uint32(diff.Int64()))
}

func invalidRequest(request models.RunRequest, receipt *eth.TxReceipt) bool {
	return receipt.Unconfirmed() ||
		(request.BlockHash != nil && *request.BlockHash != *receipt.BlockHash)
}

func meetsMinimumConfirmations(
	run *models.JobRun,
	taskRun *models.TaskRun,
	currentHeight *utils.Big) bool {
	if !taskRun.MinimumConfirmations.Valid || run.CreationHeight == nil || currentHeight == nil {
		return true
	}

	diff := blockConfirmations(currentHeight, run.CreationHeight)
	return diff.Cmp(big.NewInt(int64(taskRun.MinimumConfirmations.Uint32))) >= 0
}

func blockConfirmations(currentHeight, creationHeight *utils.Big) *big.Int {
	bigDiff := new(big.Int).Sub(currentHeight.ToInt(), creationHeight.ToInt())
	confs := bigDiff.Add(bigDiff, big.NewInt(1)) // creation of runlog alone warrants 1 confirmation
	if confs.Cmp(big.NewInt(0)) < 0 {            // negative, so floor at 0
		confs.SetUint64(0)
	}
	return confs
}
