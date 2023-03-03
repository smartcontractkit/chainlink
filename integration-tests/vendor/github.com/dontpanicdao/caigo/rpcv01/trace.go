package rpcv01

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

// not implemented for testing yet
func (provider *Provider) TransactionTrace(ctx context.Context, hash string) error {
	return ErrNotImplemented
}

// not implemented for testing yet
func (provider *Provider) TraceBlockTransactions(ctx context.Context, hash string) error {
	return ErrNotImplemented
}
