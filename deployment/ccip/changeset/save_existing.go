package changeset

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

var (
	_ cldf.ChangeSet[ExistingContractsConfig] = SaveExistingContractsChangeset
)

type Contract struct {
	Address        string
	TypeAndVersion cldf.TypeAndVersion
	ChainSelector  uint64
	// Qualifier is the semantic name used to distinguish multiple contracts with the same
	// chain, type, and version in the datastore. Leave it empty only when this contract is a
	// singleton for that key; never derive it from Address.
	Qualifier string
}

type ExistingContractsConfig struct {
	ExistingContracts []Contract
}

type existingContractKey struct {
	chainSelector uint64
	typeName      cldf.ContractType
	version       string
	qualifier     string
}

type existingContractClaim struct {
	index    int
	contract Contract
}

func validateExistingContract(ec Contract) error {
	if err := cldf.IsValidChainSelector(ec.ChainSelector); err != nil {
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
		a := common.HexToAddress(ec.Address)
		if a == (common.Address{}) {
			return fmt.Errorf("invalid address: %s", ec.Address)
		}
	default:
		return fmt.Errorf("unsupported chain family: %s", family)
	}
	if ec.TypeAndVersion.Type == "" {
		return errors.New("type must be set")
	}
	if val, err := ec.TypeAndVersion.Version.Value(); err != nil || val == "" {
		return errors.New("version must be set")
	}
	if ec.Qualifier != "" && strings.Contains(strings.ToLower(ec.Qualifier), strings.ToLower(ec.Address)) {
		return fmt.Errorf("qualifier %q contains the address being imported: name the instance in domain terms instead", ec.Qualifier)
	}
	return nil
}

func (cfg ExistingContractsConfig) Validate() error {
	claims := make(map[existingContractKey]existingContractClaim, len(cfg.ExistingContracts))
	for i, ec := range cfg.ExistingContracts {
		if err := validateExistingContract(ec); err != nil {
			return fmt.Errorf("existing contract at index %d: %w", i, err)
		}

		key := existingContractKey{
			chainSelector: ec.ChainSelector,
			typeName:      ec.TypeAndVersion.Type,
			version:       ec.TypeAndVersion.Version.String(),
			qualifier:     ec.Qualifier,
		}
		first, exists := claims[key]
		if !exists {
			claims[key] = existingContractClaim{index: i, contract: ec}
			continue
		}

		// Repeating the same imported contract is safe and makes rerunning an import
		// idempotent. A second address under the same key is different: the datastore
		// cannot retain both contracts, so fail here with the missing operator decision
		// instead of exposing the lower-level store conflict.
		if !shared.AddressesEqual(ec.ChainSelector, first.contract.Address, ec.Address) {
			return fmt.Errorf(
				"existing contracts at indexes %d and %d use the same datastore key (chain=%d, type=%q, version=%s, qualifier=%q) but have different addresses (%q and %q); provide distinct semantic qualifiers",
				first.index,
				i,
				ec.ChainSelector,
				ec.TypeAndVersion.Type,
				ec.TypeAndVersion.Version.String(),
				ec.Qualifier,
				first.contract.Address,
				ec.Address,
			)
		}

		if !first.contract.TypeAndVersion.Equal(ec.TypeAndVersion) {
			return fmt.Errorf(
				"existing contracts at indexes %d and %d use the same address and datastore key (chain=%d, type=%q, version=%s, qualifier=%q) but have different type/version labels; keep only one declaration",
				first.index,
				i,
				ec.ChainSelector,
				ec.TypeAndVersion.Type,
				ec.TypeAndVersion.Version.String(),
				ec.Qualifier,
			)
		}
	}
	return nil
}

// SaveExistingContractsChangeset imports existing contracts into both the legacy AddressBook
// and the datastore. The AddressBook retains the address and type/version, while the datastore
// also retains the caller-provided semantic qualifier needed to distinguish multiple instances.
// The qualifier is never derived from the address, and may be omitted for chain singletons.
func SaveExistingContractsChangeset(env cldf.Environment, cfg ExistingContractsConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("%w: %w", err, cldf.ErrInvalidConfig)
	}
	ab := cldf.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	for _, ec := range cfg.ExistingContracts {
		if err := shared.RecordAddress(ab, ds, ec.ChainSelector, ec.Address, ec.TypeAndVersion, ec.Qualifier); err != nil {
			env.Logger.Errorw("Failed to record existing contract", "err", err, "address", ec.Address)
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to record existing contract: %w", err)
		}
	}
	return cldf.ChangesetOutput{
		AddressBook: ab, //nolint:staticcheck // SA1019 AddressBook is deprecated
		DataStore:   ds,
	}, nil
}
