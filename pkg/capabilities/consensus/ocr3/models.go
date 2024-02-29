package ocr3

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

	WorkflowExecutionID string
	WorkflowID          string
}

type response struct {
	WorkflowExecutionID string
	capabilities.CapabilityResponse
}
