package v1_6

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[InitChainForTestingConfig] = InitChainForTesting

// InitChainForTestingConfig is a configuration struct for InitChainForTesting.
type InitChainForTestingConfig struct {
	ChainSelector uint64
}

// InitChainForTesting sequences all changesets required to add a chain to CCIP from scratch.
// This particular changeset activates the new lanes on test routers.
func InitChainForTesting(e deployment.Environment, config InitChainForTestingConfig) (deployment.ChangesetOutput, error) {
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

	/*
		Throughout the course of this changeset, aggregate all MCMS proposals into a single MCMS proposal.
		Output this MCMS proposal and all new addresses at the end of execution.

		1. Run changesets against new chain WITHOUT MCMS
		- DeployPrerequisitesChangeset
		- DeployChainContractsChangeset
		- SetOCR3OffRampChangeset
		- UpdateFeeQuoterPricesChangeset
		- UpdateFeeQuoterDestsChangeset
		- UpdateOffRampSourcesChangeset (TestRouter = true)
		- UpdateOnRampsDestsChangeset (TestRouter = true)
		- UpdateRouterRampsChangeset (for both inbound and outbound traffic, TestRouter = true)

		2. Run changesets against home chain WITH MCMS
		(likely need to extract CS logic into functions to bypass validations)
		- UpdateChainConfigChangeset
		- AddDonAndSetCandidateChangeset (candidate for commit plugin)
		- SetCandidateChangeset (for exec plugin)
		- PromoteCandidateChangeset (for commit plugin)
		- PromoteCandidateChangeset (for exec plugin)

		3. Run changesets against all requested source chains (including home chain if required) WITH MCMS
		- UpdateFeeQuoterPricesChangeset
		- UpdateFeeQuoterDestsChangeset
		- UpdateOffRampSourcesChangeset (TestRouter = true)
		- UpdateOnRampsDestsChangeset (TestRouter = true)

		4. UpdateRouterRampsChangeset WITHOUT MCMS (for both inbound and outbound traffic, TestRouter = true)

		5. RMN steps? (possibly optional for now)
	*/

	return deployment.ChangesetOutput{}, nil
}

/*

CHANGESETS (For practical BIX usage)

1. AddDonChangeset (claims the )



*/
