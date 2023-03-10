package types

import "context"

// HeadTrackable represents a generic type, implemented by the core txm,
// to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
// The generic type HEAD here indicates the chain specific Head type.
//
//go:generate mockery --quiet --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable[HEAD any] interface {
	OnNewLongestChain(ctx context.Context, head Head[HEAD])
}
