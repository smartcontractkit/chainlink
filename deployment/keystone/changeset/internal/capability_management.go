package internal

import (
	"fmt"
	"math/big"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// AddCapabilities adds the capabilities to the registry
//
// It is idempotent. It deduplicates the input capabilities.
func AddCapabilities(lggr logger.Logger, registry *kcr.CapabilitiesRegistry, chain deployment.Chain, capabilities []kcr.CapabilitiesRegistryCapability, useMCMS bool) (*mcmstypes.BatchOperation, error) {
	if len(capabilities) == 0 {
		return nil, nil
	}
	deduped, err := dedupCapabilities(registry, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to dedup capabilities: %w", err)
	}

	if useMCMS {
		return addCapabilitiesMCMSProposal(registry, deduped, chain)
	}

	tx, err := registry.AddCapabilities(chain.DeployerKey, deduped)
	if err != nil {
		err = deployment.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
	}
	lggr.Info("registered capabilities", "capabilities", deduped)

	return nil, nil
}

func addCapabilitiesMCMSProposal(registry *kcr.CapabilitiesRegistry, caps []kcr.CapabilitiesRegistryCapability, regChain deployment.Chain) (*mcmstypes.BatchOperation, error) {
	tx, err := registry.AddCapabilities(deployment.SimTransactOpts(), caps)
	if err != nil {
		err = deployment.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call AddNodeOperators: %w", err)
	}

	ops, err := proposalutils.BatchOperationForChain(regChain.Selector, registry.Address().Hex(), tx.Data(), big.NewInt(0), string(CapabilitiesRegistry), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch operation: %w", err)
	}

	return &ops, nil
}

// CapabilityID returns a unique id for the capability
// TODO: mv to chainlink-common? ref https://github.com/smartcontractkit/chainlink/blob/4fb06b4525f03c169c121a68defa9b13677f5f20/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L170
func CapabilityID(c kcr.CapabilitiesRegistryCapability) string {
	return fmt.Sprintf("%s@%s", c.LabelledName, c.Version)
}

// dedupCapabilities deduplicates the capabilities with respect to the registry
//
// the contract reverts on adding the same capability twice and that would cause the whole transaction to revert
// this is particularly important when using MCMS, because it would cause the whole batch to revert
func dedupCapabilities(registry *kcr.CapabilitiesRegistry, capabilities []kcr.CapabilitiesRegistryCapability) ([]kcr.CapabilitiesRegistryCapability, error) {
	var out []kcr.CapabilitiesRegistryCapability
	existing, err := registry.GetCapabilities(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetCapabilities: %w", err)
	}
	existingByID := make(map[[32]byte]struct{})
	for _, cap := range existing {
		existingByID[cap.HashedId] = struct{}{}
	}
	seen := make(map[string]struct{})
	for _, candidate := range capabilities {
		h, err := registry.GetHashedCapabilityId(nil, candidate.LabelledName, candidate.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to call GetHashedCapabilityId: %w", err)
		}
		// dedup input capabilities
		if _, exists := seen[CapabilityID(candidate)]; exists {
			continue
		}
		seen[CapabilityID(candidate)] = struct{}{}
		// dedup with respect to the registry
		if _, exists := existingByID[h]; !exists {
			out = append(out, candidate)
		}
	}
	return out, nil
}
