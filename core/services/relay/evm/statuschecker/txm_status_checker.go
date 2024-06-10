package statuschecker

import (
	"context"
	"fmt"
	"strings"
)

//go:generate mockery --quiet --name TransactionStatusChecker --output ../mocks/ --case=underscore

// TODO replace with actual implementation coming from "github.com/smartcontractkit/chainlink-common/pkg/types"
type TransactionStatus int

const (
	Unknown TransactionStatus = iota
	Unconfirmed
	Finalized
	Failed
	Fatal
)

type TxManager interface {
	GetTransactionStatus(ctx context.Context, transactionID string) (TransactionStatus, error)
}

type TransactionStatusChecker interface {
	CheckMessageStatus(ctx context.Context, msgID string) ([]TransactionStatus, int, error)
}

type TxmStatusChecker struct {
	txManager TxManager
}

type NoOpTxManager struct{}

func (n *NoOpTxManager) GetTransactionStatus() error {
	return nil
}

func NewTransactionStatusChecker(txManager TxManager) *TxmStatusChecker {
	return &TxmStatusChecker{txManager: txManager}
}

// CheckMessageStatus checks the status of a message by checking the status of all transactions
// associated with the message ID.
// It returns a slice of all statuses and the number of transactions found (-1 if none).
func (tsc *TxmStatusChecker) CheckMessageStatus(ctx context.Context, msgID string) ([]TransactionStatus, int, error) {
	var allStatuses []TransactionStatus
	var counter int

	for {
		transactionID := fmt.Sprintf("%s-%d", msgID, counter)
		status, err := tsc.txManager.GetTransactionStatus(ctx, transactionID)
		if err != nil {
			if status == Unknown && strings.Contains(err.Error(), fmt.Sprintf("failed to find transaction with IdempotencyKey %s", transactionID)) {
				break
			}
			return nil, counter - 1, err
		}
		allStatuses = append(allStatuses, status)
		counter++
	}

	return allStatuses, counter - 1, nil
}
