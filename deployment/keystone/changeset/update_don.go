package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var _ deployment.ChangeSet[*UpdateDonRequest] = UpdateDon

// CapabilityConfig is a struct that holds a capability and its configuration
type CapabilityConfig = internal.CapabilityConfig

type UpdateDonRequest = internal.UpdateDonRequest

type UpdateDonResponse struct {
	DonInfo kcr.CapabilitiesRegistryDONInfo
}

// UpdateDon updates the capabilities of a Don
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateDon(env deployment.Environment, req *UpdateDonRequest) (deployment.ChangesetOutput, error) {
	_, err := internal.UpdateDon(env.Logger, req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
