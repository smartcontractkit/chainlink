package stateview

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmscontracts "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/contracts/mcms"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// ValidateSolanaTimelockConfig validates a Solana timelock config against the environment
// datastore: the RBACTimelock account and the action-specific MCMS account must exist for the
// chain and parse as "programID.PDASeed" addresses. The bundle qualifier comes from
// tc.TimelockQualifierPerChain when configured, else the default CCIP bundle.
//
// CCIP-owned replacement for cldfproposalutils.TimelockConfig.ValidateSolana (address-book
// based); the framework is not patched.
func ValidateSolanaTimelockConfig(e cldf.Environment, chainSelector uint64, tc *cldfproposalutils.TimelockConfig) error {
	if tc == nil {
		return fmt.Errorf("timelock config is nil")
	}
	// default in place, mirroring the framework's validateCommon: callers observe the config
	if tc.MCMSAction == "" {
		tc.MCMSAction = mcmstypes.TimelockActionSchedule
	}
	switch tc.MCMSAction {
	case mcmstypes.TimelockActionSchedule, mcmstypes.TimelockActionCancel, mcmstypes.TimelockActionBypass:
	default:
		return fmt.Errorf("invalid MCMS action %s", tc.MCMSAction)
	}

	qualifier := DefaultMCMSQualifier
	if q, ok := tc.TimelockQualifierPerChain[chainSelector]; ok && q != "" {
		qualifier = q
	}

	validateContract := func(contractType cldf.ContractType) error {
		address, err := dataStoreSolanaContractAddress(e, chainSelector, contractType, qualifier)
		if err != nil {
			return fmt.Errorf("%s not present on the chain %w", contractType, err)
		}
		// Format is: "programID.PDASeed"
		if _, _, parseErr := mcmssolanasdk.ParseContractAddress(address); parseErr != nil {
			return fmt.Errorf("failed to parse timelock address: %w", parseErr)
		}
		return nil
	}

	if err := validateContract(mcmscontracts.RBACTimelock); err != nil {
		return err
	}

	switch tc.MCMSAction {
	case mcmstypes.TimelockActionSchedule:
		return validateContract(mcmscontracts.ProposerManyChainMultisig)
	case mcmstypes.TimelockActionCancel:
		return validateContract(mcmscontracts.CancellerManyChainMultisig)
	case mcmstypes.TimelockActionBypass:
		return validateContract(mcmscontracts.BypasserManyChainMultisig)
	}
	return nil
}

// dataStoreSolanaContractAddress resolves one Solana contract account of the given type and
// qualifier from the environment datastore.
func dataStoreSolanaContractAddress(e cldf.Environment, chainSelector uint64, contractType cldf.ContractType, qualifier string) (string, error) {
	if e.DataStore == nil {
		return "", fmt.Errorf("datastore not available for chain %d", chainSelector)
	}
	refs := e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType(contractType)),
	)
	var active []datastore.AddressRef
	for _, ref := range refs {
		if isActiveRef(ref) {
			active = append(active, ref)
		}
	}
	refs = active
	if len(refs) == 0 {
		return "", fmt.Errorf("no %s ref for chain %d", contractType, chainSelector)
	}
	var matching []datastore.AddressRef
	for _, ref := range refs {
		if ref.Qualifier == qualifier {
			matching = append(matching, ref)
		}
	}
	if len(matching) == 0 {
		return "", fmt.Errorf("no %s ref for chain %d with qualifier %q", contractType, chainSelector, qualifier)
	}
	if len(matching) > 1 {
		sortAddressRefs(matching)
		return "", fmt.Errorf("ambiguous %s on chain %d: %d refs share qualifier %q", contractType, chainSelector, len(matching), qualifier)
	}
	return matching[0].Address, nil
}
