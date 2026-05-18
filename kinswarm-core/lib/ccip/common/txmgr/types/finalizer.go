package types

import (
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type Finalizer[BLOCK_HASH types.Hashable, HEAD types.Head[BLOCK_HASH]] interface {
	// interfaces for running the underlying estimator
	services.Service
	DeliverLatestHead(head HEAD) bool
}
