package models

import "context"

// Interfaces used by the whole of chainlink should be put here

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being attached to HeadTracker.
//go:generate mockery --name HeadTrackable --output ../../internal/mocks/ --case=underscore
type HeadTrackable interface {
	Connect(head *Head) error
	Disconnect()
	OnNewLongestChain(ctx context.Context, head Head)
}
