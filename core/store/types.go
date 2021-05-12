package store

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being attached to HeadTracker.
//go:generate mockery --name HeadTrackable --output ../internal/mocks/ --case=underscore
type HeadTrackable interface {
	Connect(head *models.Head) error
	Disconnect()
	OnNewLongestChain(ctx context.Context, head models.Head)
}
