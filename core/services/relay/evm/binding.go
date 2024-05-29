package evm

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

type readBinding interface {
	GetLatestValue(ctx context.Context, params, returnVal any) error
	QueryKey(ctx context.Context, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error)
	Bind(binding commontypes.BoundContract)
	SetCodec(codec commontypes.RemoteCodec)
}
