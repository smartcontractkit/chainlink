package types

import "context"

// HeadTrackable represents a generic type, implemented by the core txm,
// to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
//go:generate mockery --quiet --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable[HEAD any] interface {
	OnNewLongestChain(ctx context.Context, head HeadView[HEAD])
}
