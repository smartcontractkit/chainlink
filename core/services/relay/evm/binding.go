package evm

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type readBinding interface {
	GetLatestValue(ctx context.Context, params, returnVal any) error
	Bind(ctx context.Context, binding commontypes.BoundContract) error
	SetCodec(codec commontypes.RemoteCodec)
	Register(ctx context.Context) error
	Unregister(ctx context.Context) error
}
