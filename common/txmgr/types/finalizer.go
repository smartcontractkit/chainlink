package types

import (
	"context"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Finalizer[BLOCK_HASH types.Hashable, HEAD types.Head[BLOCK_HASH]] interface {
	// interfaces for running the underlying estimator
	services.Service
	DeliverLatestHead(head HEAD) bool
	SetResumeCallback(callback func(ctx context.Context, id uuid.UUID, result interface{}, err error) error)
}
