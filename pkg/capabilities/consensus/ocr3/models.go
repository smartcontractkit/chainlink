package consensus

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type request struct {
	Observations values.Value `mapstructure:"-"`
	ExpiresAt    time.Time

	CallbackCh chan<- capabilities.CapabilityResponse
	RequestCtx context.Context
	RequestID  string
}

//nolint:unused
type response struct {
	Value     values.Value
	Err       error
	RequestID string
}
