package keystone

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

// AddCapabilities adds the capabilities to the registry
// it tries to add all capabilities in one go, if that fails, it falls back to adding them one by one
func AddCapabilities(lggr logger.Logger, registry *kcr.CapabilitiesRegistry, chain deployment.Chain, capabilities []kcr.CapabilitiesRegistryCapability) error {
	if len(capabilities) == 0 {
		return nil
	}
	// dedup capabilities
	var deduped []kcr.CapabilitiesRegistryCapability
	seen := make(map[string]struct{})
	for _, cap := range capabilities {
		if _, ok := seen[CapabilityID(cap)]; !ok {
			seen[CapabilityID(cap)] = struct{}{}
			deduped = append(deduped, cap)
		}
	}

	tx, err := registry.AddCapabilities(chain.DeployerKey, deduped)
	if err != nil {
		err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
		// no typed errors in the abi, so we have to do string matching
		// try to add all capabilities in one go, if that fails, fall back to 1-by-1
		if !strings.Contains(err.Error(), "CapabilityAlreadyExists") {
			return fmt.Errorf("failed to call AddCapabilities: %w", err)
		}
		lggr.Warnw("capabilities already exist, falling back to 1-by-1", "capabilities", deduped)
		for _, cap := range deduped {
			tx, err = registry.AddCapabilities(chain.DeployerKey, []kcr.CapabilitiesRegistryCapability{cap})
			if err != nil {
				err = DecodeErr(kcr.CapabilitiesRegistryABI, err)
				if strings.Contains(err.Error(), "CapabilityAlreadyExists") {
					lggr.Warnw("capability already exists, skipping", "capability", cap)
					continue
				}
				return fmt.Errorf("failed to call AddCapabilities for capability %v: %w", cap, err)
			}
			// 1-by-1 tx is pending and we need to wait for it to be mined
			_, err = chain.Confirm(tx)
			if err != nil {
				return fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
			}
			lggr.Debugw("registered capability", "capability", cap)

		}
	} else {
		// the bulk add tx is pending and we need to wait for it to be mined
		_, err = chain.Confirm(tx)
		if err != nil {
			return fmt.Errorf("failed to confirm AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
		}
		lggr.Info("registered capabilities", "capabilities", deduped)
	}
	return nil
}

// CapabilityID returns a unique id for the capability
// TODO: mv to chainlink-common? ref https://github.com/smartcontractkit/chainlink/blob/4fb06b4525f03c169c121a68defa9b13677f5f20/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L170
func CapabilityID(c kcr.CapabilitiesRegistryCapability) string {
	return fmt.Sprintf("%s@%s", c.LabelledName, c.Version)
}
