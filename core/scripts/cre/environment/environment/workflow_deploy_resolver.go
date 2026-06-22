package environment

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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
//
// don_family and gateway URL must match the selected workflow DON so registry registration,
// vault secrets (when used), and node-side workflow filtering stay aligned with topology.
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

	if err := validateWorkflowDeployFlags(cmd, donMeta, selector, donFamilyFlag); err != nil {
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
		// Vault secrets flow (FetchVaultPublicKey / ExecuteSecrets) needs the gateway for
		// this don_family, not the first gateway in the topology.
		if url, err := resolver.resolveGatewayURL(family); err != nil {
			return targets, fmt.Errorf("❌ failed to resolve gateway URL for don_family %q: %w", family, err)
		} else {
			targets.gatewayURL = url
		}
	}

	return targets, nil
}

// finalizeWorkflowDonFamily applies the legacy default when pairing is off. When pairing
// is on, don_family must be explicit — it is the routing key shared by cap registry,
// workflow registry limits, gateway connectors, and node workflow sync.
func finalizeWorkflowDonFamily(donFamily string, pairingEnabled bool) (string, error) {
	if donFamily != "" {
		return donFamily, nil
	}
	if pairingEnabled {
		return "", fmt.Errorf("❌ --don-family is required when topology uses don_family gateway pairing (or set nodesets.don_family in the local CRE state file)")
	}
	return envconfig.DefaultDONFamily, nil
}

// validateWorkflowDeployFlags rejects inconsistent explicit deploy flags against local CRE state.
// Only runs when both --workflow-don-name and --don-family are set; catches typos before
// on-chain UpsertWorkflow writes the wrong family.
func validateWorkflowDeployFlags(cmd *cobra.Command, donMeta *cre.DonMetadata, selector workflowDONSelector, donFamilyFlag string) error {
	if !cmd.Flags().Changed("workflow-don-name") || !cmd.Flags().Changed("don-family") {
		return nil
	}
	if donMeta.DonFamily == "" {
		return nil
	}
	if donFamilyFlag != donMeta.DonFamily {
		return fmt.Errorf(
			"❌ --don-family %q does not match don_family %q for workflow DON %q in local CRE state",
			donFamilyFlag,
			donMeta.DonFamily,
			donMeta.Name,
		)
	}
	if name := strings.TrimSpace(selector.ExplicitName); name != "" && name != donMeta.Name {
		return fmt.Errorf(
			"❌ --workflow-don-name %q does not match workflow DON %q resolved from local CRE state",
			name,
			donMeta.Name,
		)
	}
	return nil
}

// resolveGatewayURL picks the gateway for donFamily when pairing is enabled, otherwise
// falls back to the first connector (legacy single-gateway topologies).
func (r *LocalCREStateResolver) resolveGatewayURL(donFamily string) (string, error) {
	if url, err := r.GatewayURLForDonFamily(donFamily); err == nil {
		return url, nil
	}
	return r.GatewayURL()
}
