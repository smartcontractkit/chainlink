package types

import (
	"context"

	commontypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

// HeadTrackable is implemented by the core txm,
// to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
//go:generate mockery --quiet --name HeadTrackable --output ./mocks/ --case=underscore
type HeadTrackable[H commontypes.Head[HASH], HASH commontypes.Hashable] interface {
	OnNewLongestChain(ctx context.Context, head H)
}
