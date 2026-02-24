package changeset

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// RenounceTimelockDeployerChainsChangeset renounces the deployer/KMS key's ADMIN role
// on the RBAC Timelock for each given chain so that only the Timelock is admin of itself.
// Run after the transfer_mcms_ownership_to_timelock proposal has been executed.
var RenounceTimelockDeployerChainsChangeset cldf.ChangeSetV2[types.RenounceTimelockDeployerChainsConfig] = renounceTimelockDeployerChainsChangeset{}

type renounceTimelockDeployerChainsChangeset struct{}

func (r renounceTimelockDeployerChainsChangeset) VerifyPreconditions(e cldf.Environment, cfg types.RenounceTimelockDeployerChainsConfig) error {
	return ValidateRenounceTimelockDeployerChainsConfig(e, cfg)
}

func (r renounceTimelockDeployerChainsChangeset) Apply(e cldf.Environment, cfg types.RenounceTimelockDeployerChainsConfig) (cldf.ChangesetOutput, error) {
	for _, chainSelector := range cfg.ChainSelectors {
		output, err := commonchangeset.RenounceTimelockDeployer(e, commonchangeset.RenounceTimelockDeployerConfig{
			ChainSel: chainSelector,
		})
		if err != nil {
			return output, fmt.Errorf("chain %d: %w", chainSelector, err)
		}
	}
	return cldf.ChangesetOutput{}, nil
}
