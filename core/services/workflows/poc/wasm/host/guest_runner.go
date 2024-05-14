package host

import "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"

// GuestRunner probably should be a service so you get start and stop once
// Also, if this lives in a different directory, then we can allow ourselves to easily add say a V8 runner etc.
type GuestRunner interface {
	// Run would we want to send config here?
	Run() (*workflow.Spec, error)
}
