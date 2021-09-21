package types

import (
	"context"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

type Tracker interface {
	HighestSeenHeadFromDB() (*eth.Head, error)
	Start() error
	Stop() error
	SetLogger(logger logger.Logger)
	Ready() error
	Healthy() error
}

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//go:generate mockery --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable interface {
	OnNewLongestChain(ctx context.Context, head eth.Head)
}

type SubscribeFunc func(callback HeadTrackable) (unsubscribe func())

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (currentLongestChain *eth.Head, unsubscribe func())
}

// HeadBroadcaster is the external interface of headBroadcaster
//go:generate mockery --name HeadBroadcaster --output ../mocks/ --case=underscore
type HeadBroadcaster interface {
	service.Service
	HeadTrackable
	Subscribe(callback HeadTrackable) (currentLongestChain *eth.Head, unsubscribe func())
}
