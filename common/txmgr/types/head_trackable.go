package types

import "context"

// HeadTrackable is implemented by the core txm,
// to be able to receive head events from any chain.
// Chain implementations should notify head events to the core txm via this interface.
//
//go:generate mockery --quiet --name HeadTrackable --output ./mocks/ --case=underscore
type HeadTrackable interface {
	OnNewLongestChain(ctx context.Context, head Head)
}
