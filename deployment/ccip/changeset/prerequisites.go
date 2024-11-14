package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
)

var (
	_                        deployment.ChangeSet[PrerequisiteConfig] = InitializePrerequisites
	mapPrerequisiteContracts                                          = map[deployment.ContractType]struct{}{
		ccipdeployment.LinkToken:          {},
		ccipdeployment.WETH9:              {},
		ccipdeployment.TokenAdminRegistry: {},
		ccipdeployment.RegistryModule:     {},
		ccipdeployment.Router:             {},
	}
)

// InitializePrerequisites loads the existing contracts into the address book.
// This is required for contracts which can be reused from previous versions of CCIP
// If PrerequisiteConfig.Deploy is true, it will deploy the prerequisite contracts except Router
// Router is deployed as part of DeployChainContracts Changeset due to its dependency on RMNProxy
func InitializePrerequisites(env deployment.Environment, cfg PrerequisiteConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	if cfg.Deploy {
		err = ccipdeployment.DeployPrerequisiteContracts(env, ab, env.Chains[cfg.ChainSelector])
		if err != nil {
			env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
		}
		return deployment.ChangesetOutput{
			Proposals:   []timelock.MCMSWithTimelockProposal{},
			AddressBook: ab,
			JobSpecs:    nil,
		}, nil
	}
	for _, ec := range cfg.ExistingContracts {
		err = ab.Save(cfg.ChainSelector, ec.Address.String(), ec.TypeAndVersion)
		if err != nil {
			env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save existing contract: %w", err)
		}
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}

type ContractConfig struct {
	Address        common.Address
	TypeAndVersion deployment.TypeAndVersion
}

type PrerequisiteConfig struct {
	ChainSelector     uint64
	ExistingContracts map[deployment.ContractType]ContractConfig
	Deploy            bool
}

func (c PrerequisiteConfig) Validate() error {
	if c.ChainSelector == 0 {
		return fmt.Errorf("chain selector must be set")
	}
	_, err := chain_selectors.ChainIdFromSelector(c.ChainSelector)
	if err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	if c.Deploy {
		return nil
	}
	for ct := range mapPrerequisiteContracts {
		if _, ok := c.ExistingContracts[ct]; !ok {
			return fmt.Errorf("missing existing contract: %s", ct)
		}
	}
	for _, ec := range c.ExistingContracts {
		if ec.Address == (common.Address{}) {
			return fmt.Errorf("address must be set")
		}
		if ec.TypeAndVersion.Type == "" {
			return fmt.Errorf("type must be set")
		}
		if val, err := ec.TypeAndVersion.Version.Value(); err != nil || val == "" {
			return fmt.Errorf("version must be set")
		}
	}
	return nil
}
