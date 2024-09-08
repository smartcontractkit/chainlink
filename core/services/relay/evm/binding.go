package evm

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

type readBinding interface {
	Bind(ctx context.Context, binding commontypes.BoundContract) error
	SetCodec(codec commontypes.RemoteCodec)
	Register(ctx context.Context) error
	Unregister(ctx context.Context) error
	GetLatestValue(ctx context.Context, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error
	QueryKey(ctx context.Context, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error)
}
