package v1_6

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[InitChainE2EConfig] = InitChainE2E

// InitChainE2EConfig is a configuration struct for InitChainE2E.
type InitChainE2EConfig struct {
	ChainSelector uint64
}

// InitChainE2E sequences all changesets required to add a chain to CCIP from scratch.
func InitChainE2E(e deployment.Environment, config InitChainE2EConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()

	output, err := changeset.DeployPrerequisitesChangeset(e, changeset.DeployPrerequisiteConfig{
		Configs: []changeset.DeployPrerequisiteConfigPerChain{
			changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: config.ChainSelector,
			},
		},
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run DeployPrerequisitesChangeset: %w", err)
	}
	newAddresses.Merge(output.AddressBook)

	// ... and so on

	return deployment.ChangesetOutput{}, nil
}
