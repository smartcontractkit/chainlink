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
		ccipdeployment.ARMProxy:           {},
	}
)

func InitializePrerequisites(env deployment.Environment, cfg PrerequisiteConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	for _, ec := range cfg.ExistingContracts {
		err = ab.Save(ec.ChainSelector, ec.Address.String(), ec.TypeAndVersion)
		if err != nil {
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
	ChainSelector  uint64
	Address        common.Address
	TypeAndVersion deployment.TypeAndVersion
}

type PrerequisiteConfig struct {
	ExistingContracts   map[deployment.ContractType]ContractConfig
	DeployPrerequisites bool
}

func (c PrerequisiteConfig) Validate() error {
	if !c.DeployPrerequisites && len(c.ExistingContracts) == 0 {
		return fmt.Errorf("either deploy prerequisites or provide existing contracts")
	}
	for ct := range mapPrerequisiteContracts {
		if _, ok := c.ExistingContracts[ct]; !ok {
			return fmt.Errorf("missing existing contract: %s", ct)
		}
	}
	for _, ec := range c.ExistingContracts {
		if ec.ChainSelector == 0 {
			return fmt.Errorf("chain selector must be set")
		}
		_, err := chain_selectors.ChainIdFromSelector(ec.ChainSelector)
		if err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", ec.ChainSelector, err)
		}
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
