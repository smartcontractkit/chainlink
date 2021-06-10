package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//go:generate mockery --name HeadTrackable --output ../internal/mocks/ --case=underscore
type HeadTrackable interface {
	Connect(head *models.Head) error
	OnNewLongestChain(ctx context.Context, head models.Head)
}

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (unsubscribe func())
}

// HeadTrackableCallback is a simple wrapper around an On Connect callback
type HeadTrackableCallback struct {
	OnConnect func() error
}

func (c *HeadTrackableCallback) Connect(*models.Head) error {
	return c.OnConnect()
}

func (c *HeadTrackableCallback) OnNewLongestChain(context.Context, models.Head) {}
