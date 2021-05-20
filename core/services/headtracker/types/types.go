package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being attached to HeadTracker.
//go:generate mockery --name HeadTrackable --output ../internal/mocks/ --case=underscore
type HeadTrackable interface {
	Connect(head *models.Head) error
	OnNewLongestChain(ctx context.Context, head models.Head)
}

// HeadBroadcastable defines the interface for listeners
type HeadBroadcastable interface {
	Connect(head *models.Head) error
	OnNewLongestChain(ctx context.Context, head models.Head)
}

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadBroadcastable) (unsubscribe func())
}
