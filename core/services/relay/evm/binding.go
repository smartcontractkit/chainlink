package evm

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

type readBinding interface {
	GetLatestValue(ctx context.Context, params, returnVal any) error
	QueryKey(ctx context.Context, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error)
	Bind(ctx context.Context, binding commontypes.BoundContract) error
	SetCodec(codec commontypes.RemoteCodec)
	Register(ctx context.Context) error
	Unregister(ctx context.Context) error
}
