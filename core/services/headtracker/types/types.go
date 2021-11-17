package types

import (
	"context"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

type Tracker interface {
	HighestSeenHeadFromDB(context.Context) (*eth.Head, error)
	Start() error
	Stop() error
	SetLogLevel(lvl zapcore.Level)
	Ready() error
	Healthy() error
	LatestChain() *eth.Head
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
