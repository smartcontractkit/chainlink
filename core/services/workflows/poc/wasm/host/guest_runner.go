package host

import "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"

// GuestRunner probably should be a service so you get start and stop once
type GuestRunner interface {
	// Run would we want to send config here?
	Run() (*workflow.Spec, error)
}
