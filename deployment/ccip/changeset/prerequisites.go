package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var (
	_ deployment.ChangeSet[DeployPrerequisiteConfig] = DeployPrerequisites
)

// DeployPrerequisites deploys the pre-requisite contracts for CCIP
// pre-requisite contracts are the contracts which can be reused from previous versions of CCIP
func DeployPrerequisites(env deployment.Environment, cfg DeployPrerequisiteConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployPrerequisiteChainContracts(env, ab, cfg.ChainSelectors)
	if err != nil {
		env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
		return deployment.ChangesetOutput{
			AddressBook: ab,
		}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}

type DeployPrerequisiteConfig struct {
	ChainSelectors []uint64
	// TODO handle tokens and feeds in prerequisite config
	Tokens map[ccipdeployment.TokenSymbol]common.Address
	Feeds  map[ccipdeployment.TokenSymbol]common.Address
}

func (c DeployPrerequisiteConfig) Validate() error {
	for _, cs := range c.ChainSelectors {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}
