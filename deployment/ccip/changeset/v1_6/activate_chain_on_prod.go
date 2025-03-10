package v1_6

import (
	"github.com/smartcontractkit/chainlink/deployment"
)

var _ deployment.ChangeSet[ActivateChainOnProdConfig] = ActivateChainOnProd

// ActivateChainOnProdConfig is a configuration struct for ActivateChainOnProd.
type ActivateChainOnProdConfig struct{}

// ActivateChainOnProd activates a new chain on production routers in CCIP.
// InitChainForTesting should be executed first, along with some E2E transfers using test routers.
func ActivateChainOnProd(e deployment.Environment, config ActivateChainOnProdConfig) (deployment.ChangesetOutput, error) {
	/*
		Throughout the course of this changeset, aggregate all MCMS proposals into a single MCMS proposal.
		Output this MCMS proposal and all new addresses at the end of execution.

		1. Run changesets against new chain WITHOUT MCMS
		- UpdateOffRampSourcesChangeset (TestRouter = false)
		- UpdateOnRampsDestsChangeset (TestRouter = false)
		- UpdateRouterRampsChangeset (for both inbound and outbound traffic, TestRouter = false)

		2. Run changesets against all requested source chains (including home chain if required) WITH MCMS
		- UpdateOffRampSourcesChangeset (TestRouter = false)
		- UpdateOnRampsDestsChangeset (TestRouter = false)
		- UpdateRouterRampsChangeset (for both inbound and outbound traffic, TestRouter = false)

		3. MCMS cleanup
		- TransferToMCMSWithTimelock on all new chain contracts WITH MCMS
		- Run RenounceTimelockDeployer on new chain to revoke admin rights from deployer key
	*/

	return deployment.ChangesetOutput{}, nil
}
