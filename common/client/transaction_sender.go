package client

import "context"

type TransactionSender[TX any] interface {
	SendTransaction(ctx context.Context, tx TX) (SendTxReturnCode, error)
}

// TxErrorClassifier - defines interface of a function that transforms raw RPC error into the SendTxReturnCode enum
// (e.g. Successful, Fatal, Retryable, etc.)
type TxErrorClassifier[TX any] func(tx TX, err error) SendTxReturnCode

// SendTxRPCClient - defines interface of an RPC used by TransactionSender to broadcast transaction
type SendTxRPCClient[TX any] interface {
	SendTransaction(ctx context.Context, tx TX) error
}
