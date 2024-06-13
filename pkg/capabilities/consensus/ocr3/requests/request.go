package requests

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Request struct {
	Observations *values.List `mapstructure:"-"`
	ExpiresAt    time.Time

	// CallbackCh is a channel to send a response back to the requester
	// after the request has been processed or timed out.
	CallbackCh chan capabilities.CapabilityResponse
	StopCh     services.StopChan

	WorkflowExecutionID string
	WorkflowID          string
	WorkflowOwner       string
	WorkflowName        string
	WorkflowDonID       string
	ReportID            string
}

type Response struct {
	WorkflowExecutionID string
	capabilities.CapabilityResponse
}
