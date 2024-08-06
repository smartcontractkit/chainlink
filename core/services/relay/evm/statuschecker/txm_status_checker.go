package statuschecker

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// CCIPTransactionStatusChecker is an interface that defines the method for checking the status of a transaction.
// CheckMessageStatus checks the status of a transaction for a given message ID.
// It returns a list of transaction statuses, the retry counter, and an error if any occurred during the process.
//

type CCIPTransactionStatusChecker interface {
	CheckMessageStatus(ctx context.Context, msgID string) (transactionStatuses []types.TransactionStatus, retryCounter int, err error)
}

type TxmStatusChecker struct {
	getTransactionStatus func(ctx context.Context, transactionID string) (types.TransactionStatus, error)
}

func NewTxmStatusChecker(getTransactionStatus func(ctx context.Context, transactionID string) (types.TransactionStatus, error)) *TxmStatusChecker {
	return &TxmStatusChecker{getTransactionStatus: getTransactionStatus}
}

// CheckMessageStatus checks the status of a message by checking the status of all transactions associated with the message ID.
// It returns a slice of all statuses and the number of transactions found (-1 if none).
// The key will follow the format: <msgID>-<counter>. TXM will be queried for each key until a NotFound error is returned.
// The goal is to find all transactions associated with a message ID and snooze messages if they are fatal in the Execution Plugin.
func (tsc *TxmStatusChecker) CheckMessageStatus(ctx context.Context, msgID string) ([]types.TransactionStatus, int, error) {
	var counter int
	const maxStatuses = 1000 // Cap the number of statuses to avoid infinite loop

	allStatuses := make([]types.TransactionStatus, 0)

	for {
		transactionID := fmt.Sprintf("%s-%d", msgID, counter)
		status, err := tsc.getTransactionStatus(ctx, transactionID)
		if err != nil && status == types.Unknown {
			// If the status is unknown and err not nil, it means the transaction was not found
			break
		}
		allStatuses = append(allStatuses, status)
		counter++

		// Break the loop if the cap is reached
		if counter >= maxStatuses {
			return allStatuses, counter - 1, fmt.Errorf("maximum number of statuses reached, possible infinite loop")
		}
	}

	return allStatuses, counter - 1, nil
}
