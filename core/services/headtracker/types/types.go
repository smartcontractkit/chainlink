package types

import (
	"context"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

type Tracker interface {
	services.Service
	HighestSeenHeadFromDB(context.Context) (*eth.Head, error)
	SetLogLevel(lvl zapcore.Level)
}

// HeadTrackable represents any object that wishes to respond to ethereum events,
// after being subscribed to HeadBroadcaster
//go:generate mockery --name HeadTrackable --output ../mocks/ --case=underscore
type HeadTrackable interface {
	OnNewLongestChain(ctx context.Context, head *eth.Head)
}

type SubscribeFunc func(callback HeadTrackable) (unsubscribe func())

type HeadBroadcasterRegistry interface {
	Subscribe(callback HeadTrackable) (currentLongestChain *eth.Head, unsubscribe func())
}

// HeadBroadcaster relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
//go:generate mockery --name HeadBroadcaster --output ../mocks/ --case=underscore
type HeadBroadcaster interface {
	services.Service
	BroadcastNewLongestChain(head *eth.Head)
	HeadBroadcasterRegistry
}
