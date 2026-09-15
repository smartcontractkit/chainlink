package cre

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// UpdateDONCapabilityConfig rewrites one capability's config on a DON, leaving
// the DON's other capabilities, nodes, and settings as they are.
//
// This exists for configuration only knowable once the environment is running.
// Capabilities that poll the registry pick the new config up on their next
// refresh, so publishing it does not require restarting the environment.
//
// capabilityName is matched without a version, since the registry keys
// capabilities as "name@version".
func UpdateDONCapabilityConfig(
	ctx context.Context,
	sethClient *seth.Client,
	capabilitiesRegistryAddr string,
	donName string,
	capabilityName string,
	config []byte,
) error {
	capReg, err := capabilities_registry_v2.NewCapabilitiesRegistry(
		common.HexToAddress(capabilitiesRegistryAddr), sethClient.Client,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create capabilities registry wrapper")
	}

	don, err := capReg.GetDONByName(&bind.CallOpts{Context: ctx}, donName)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch DON %q from capabilities registry", donName)
	}

	// updateDON replaces the whole set, so carry every capability the DON has.
	updated := make([]capabilities_registry_v2.CapabilitiesRegistryCapabilityConfiguration, len(don.CapabilityConfigurations))
	copy(updated, don.CapabilityConfigurations)

	var found bool
	for i := range updated {
		if strings.HasPrefix(updated[i].CapabilityId, capabilityName+"@") {
			updated[i].Config = config
			found = true
		}
	}
	if !found {
		return errors.Errorf("capability %q is not configured on DON %q", capabilityName, donName)
	}

	_, err = sethClient.Decode(capReg.UpdateDON(sethClient.NewTXOpts(), don.Id, capabilities_registry_v2.CapabilitiesRegistryUpdateDONParams{
		Name:                     don.Name,
		Config:                   don.Config,
		CapabilityConfigurations: updated,
		Nodes:                    don.NodeP2PIds,
		F:                        don.F,
		IsPublic:                 don.IsPublic,
	}))
	if err != nil {
		return errors.Wrapf(err, "failed to submit updateDON for DON %q", donName)
	}

	return nil
}
