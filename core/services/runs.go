package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func validateOnMainChain(run *models.JobRun, taskRun *models.TaskRun, ethClient eth.Client) error {
	txhash := run.RunRequest.TxHash
	if txhash == nil {
		logger.Debug("validateOnMainChain not performing check, runRequest missing txHash")
		return nil
	}
	if !taskRun.MinRequiredIncomingConfirmations.Valid || taskRun.MinRequiredIncomingConfirmations.Uint32 == 0 {
		logger.Debug("validateOnMainChain not performing check, no minimum required confirmations")
		return nil
	}

	receipt, err := ethClient.TransactionReceipt(context.TODO(), *txhash)
	if err != nil {
		return err
	}
	if models.ReceiptIsUnconfirmed(receipt) {
		logger.Debug("validateOnMainChain performing check, transaction is yet to confirm")
		return nil
	}
	request := run.RunRequest
	if request.BlockHash != nil && *request.BlockHash != receipt.BlockHash {
		return fmt.Errorf(
			"TxHash %s initiating run %s not on main chain; presumably has been uncled",
			txhash.Hex(),
			run.ID.String(),
		)
	}
	return nil
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
