package evm

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
)

type readBinding interface {
	GetLatestValue(ctx context.Context, address common.Address, params, returnVal any) error
	QueryOne(ctx context.Context, address common.Address, queryFilter query.Filter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error)
	Bind(ctx context.Context, address common.Address) error
	UnBind(ctx context.Context, address common.Address) error
	SetCodec(codec commontypes.RemoteCodec)
	Register(ctx context.Context, address common.Address) error
	Unregister(ctx context.Context, address common.Address) error
	UnregisterAll(ctx context.Context) error
}
