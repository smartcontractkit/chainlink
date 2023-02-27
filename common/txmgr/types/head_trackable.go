package types

import "context"

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//
//go:generate mockery --quiet --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable interface {
	OnNewLongestChain(ctx context.Context, head *HeadView)
}
