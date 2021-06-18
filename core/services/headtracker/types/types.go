package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type Tracker interface {
	HighestSeenHeadFromDB() (*models.Head, error)
	Start() error
	Stop() error
	SetLogger(logger *logger.Logger)
}

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//go:generate mockery --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable interface {
	Connect(head *models.Head) error
	OnNewLongestChain(ctx context.Context, head models.Head)
}

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (unsubscribe func())
}

// HeadBroadcaster is the external interface of headBroadcaster
//go:generate mockery --name HeadBroadcaster --output ../mocks/ --case=underscore
type HeadBroadcaster interface {
	service.Service
	HeadTrackable
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
