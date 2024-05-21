package headtracker

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

// Config represents a subset of options needed by head tracker
type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	FinalityDepth() uint32
	FinalityTagEnabled() bool
}

type HeadTrackerConfig interface {
	config.HeadTracker
}
