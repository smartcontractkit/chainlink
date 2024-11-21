package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	_ deployment.ChangeSet[ExistingContractsConfig] = SaveExistingContracts
)

type Contract struct {
	Address        common.Address
	TypeAndVersion deployment.TypeAndVersion
	ChainSelector  uint64
}

type ExistingContractsConfig struct {
	ExistingContracts []Contract
}

func (cfg ExistingContractsConfig) Validate() error {
	for _, ec := range cfg.ExistingContracts {
		if err := deployment.IsValidChainSelector(ec.ChainSelector); err != nil {
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

func SaveExistingContracts(env deployment.Environment, cfg ExistingContractsConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	for _, ec := range cfg.ExistingContracts {
		err = ab.Save(ec.ChainSelector, ec.Address.String(), ec.TypeAndVersion)
		if err != nil {
			env.Logger.Errorw("Failed to save existing contract", "err", err, "addressBook", ab)
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save existing contract: %w", err)
		}
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}
