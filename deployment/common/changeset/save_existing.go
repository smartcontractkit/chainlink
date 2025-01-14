package changeset

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/mr-tron/base58"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

var (
	_ deployment.ChangeSet[ExistingContractsConfig] = SaveExistingContracts
)

type Contract struct {
	Address        string
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
		if ec.Address == "" {
			return errors.New("address must be set")
		}
		family, err := chain_selectors.GetSelectorFamily(ec.ChainSelector)
		if err != nil {
			return err
		}
		switch family {
		case chain_selectors.FamilySolana:
			decoded, err := base58.Decode(ec.Address)
			if err != nil {
				return fmt.Errorf("address must be a valid Solana address (i.e. base58 encoded): %w", err)
			}
			if len(decoded) != 32 {
				return fmt.Errorf("address must be a valid Solana address, got %d bytes expected 32", len(decoded))
			}
		case chain_selectors.FamilyEVM:
			// EVM is the default case
			fallthrough
		default:
			// aggregator must be an ethereum address
			decoded, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(ec.Address, "0x")))
			if err != nil {
				return fmt.Errorf("address must be a valid ethereum address (i.e hex encoded 20 bytes): %w", err)
			}
			if len(decoded) != 20 {
				return fmt.Errorf("address must be a valid ethereum address, got %d bytes expected 20", len(decoded))
			}
		}
		if ec.TypeAndVersion.Type == "" {
			return errors.New("type must be set")
		}
		if val, err := ec.TypeAndVersion.Version.Value(); err != nil || val == "" {
			return errors.New("version must be set")
		}
	}
	return nil
}

// SaveExistingContracts saves the existing contracts to the address book.
// Caller should update the environment's address book with the returned addresses.
func SaveExistingContracts(env deployment.Environment, cfg ExistingContractsConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	for _, ec := range cfg.ExistingContracts {
		err = ab.Save(ec.ChainSelector, ec.Address, ec.TypeAndVersion)
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
