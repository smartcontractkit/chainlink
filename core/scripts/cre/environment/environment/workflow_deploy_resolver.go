package environment

import (
	"fmt"

	"github.com/spf13/cobra"

	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

type workflowDeployTargets struct {
	donID      uint32
	donFamily  string
	gatewayURL string
}

// resolveWorkflowDeployTargets fills donID, donFamily, and gatewayURL from CLI flags and,
// when available, the local CRE state file. Container name patterns (e.g. feeds-zone-a-node)
// take precedence over single-workflow-DON defaults.
func resolveWorkflowDeployTargets(
	cmd *cobra.Command,
	resolver *LocalCREStateResolver,
	containerPattern string,
	donIDFlag uint32,
	donFamilyFlag string,
	gatewayURLFlag string,
) (workflowDeployTargets, error) {
	targets := workflowDeployTargets{
		donID:      donIDFlag,
		donFamily:  donFamilyFlag,
		gatewayURL: gatewayURLFlag,
	}

	if resolver == nil {
		family, err := finalizeWorkflowDonFamily("", false)
		if err != nil {
			return targets, err
		}
		targets.donFamily = family
		return targets, nil
	}

	if !cmd.Flags().Changed("don-id") {
		if id, ok := resolver.resolveWorkflowDONID(containerPattern); ok {
			targets.donID = id
		}
	}

	if !cmd.Flags().Changed("don-family") {
		if family, ok := resolver.resolveWorkflowDonFamily(containerPattern); ok {
			targets.donFamily = family
		}
	}

	family, err := finalizeWorkflowDonFamily(targets.donFamily, resolver.DonFamilyGatewayPairingEnabled())
	if err != nil {
		return targets, err
	}
	targets.donFamily = family

	if !cmd.Flags().Changed("gateway-url") {
		if url, ok := resolver.resolveGatewayURL(family); ok {
			targets.gatewayURL = url
		}
	}

	return targets, nil
}

func finalizeWorkflowDonFamily(donFamily string, pairingEnabled bool) (string, error) {
	if donFamily != "" {
		return donFamily, nil
	}
	if pairingEnabled {
		return "", fmt.Errorf("❌ --don-family is required when topology uses don_family gateway pairing (or set nodesets.don_family in the local CRE state file)")
	}
	return envconfig.DefaultDONFamily, nil
}

func (r *LocalCREStateResolver) resolveWorkflowDONID(containerPattern string) (uint32, bool) {
	if containerPattern != "" {
		if id, err := r.WorkflowDONIDForContainerPattern(containerPattern); err == nil {
			return id, true
		}
	}
	id, err := r.WorkflowDONID()
	return id, err == nil
}

func (r *LocalCREStateResolver) resolveWorkflowDonFamily(containerPattern string) (string, bool) {
	if containerPattern != "" {
		if family, err := r.WorkflowDONFamilyForContainerPattern(containerPattern); err == nil {
			return family, true
		}
	}
	family, err := r.WorkflowDONFamily()
	return family, err == nil && family != ""
}

func (r *LocalCREStateResolver) resolveGatewayURL(donFamily string) (string, bool) {
	if url, err := r.GatewayURLForDonFamily(donFamily); err == nil {
		return url, true
	}
	url, err := r.GatewayURL()
	return url, err == nil
}
