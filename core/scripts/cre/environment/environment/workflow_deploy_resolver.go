package environment

import (
	"fmt"

	"github.com/spf13/cobra"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

type workflowDeployTargets struct {
	donID      uint32
	donFamily  string
	gatewayURL string
}

// resolveWorkflowDeployTargets fills donID, donFamily, and gatewayURL from CLI flags and
// the local CRE state file. Multi-DON topologies require --workflow-don-name or an
// unambiguous container-name-pattern.
func resolveWorkflowDeployTargets(
	cmd *cobra.Command,
	resolver *LocalCREStateResolver,
	selector workflowDONSelector,
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
		family, err := finalizeWorkflowDonFamily(donFamilyFlag, false)
		if err != nil {
			return targets, err
		}
		targets.donFamily = family
		return targets, nil
	}

	donMeta, err := resolver.ResolveWorkflowDONMetadata(selector)
	if err != nil {
		return targets, err
	}

	if !cmd.Flags().Changed("don-id") {
		targets.donID = libc.MustSafeUint32FromUint64(donMeta.ID)
	}
	if !cmd.Flags().Changed("don-family") && donMeta.DonFamily != "" {
		targets.donFamily = donMeta.DonFamily
	}

	family, err := finalizeWorkflowDonFamily(targets.donFamily, resolver.DonFamilyGatewayPairingEnabled())
	if err != nil {
		return targets, err
	}
	targets.donFamily = family

	if !cmd.Flags().Changed("gateway-url") {
		if url, err := resolver.resolveGatewayURL(family); err != nil {
			return targets, fmt.Errorf("❌ failed to resolve gateway URL for don_family %q: %w", family, err)
		} else {
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

func (r *LocalCREStateResolver) resolveGatewayURL(donFamily string) (string, error) {
	if url, err := r.GatewayURLForDonFamily(donFamily); err == nil {
		return url, nil
	}
	return r.GatewayURL()
}
