package environment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// workflowDeployTargets are the on-chain and gateway parameters passed to deployWorkflow.
//
// workflow_don_resolver.go decides *which* workflow DON is targeted; this file decides
// *what values* that deploy uses. See workflow.go for container copy vs DON identity flags.
//
// Shard index (--shard-index) and nodesets.name (--workflow-don-name) live on
// workflowDONSelector and are consumed during ResolveWorkflowDONMetadata; they are not
// UpsertWorkflow inputs. After resolution, donID identifies the shard on chain.
type workflowDeployTargets struct {
	donID      uint32 // Capabilities Registry DON ID written into UpsertWorkflow
	donFamily  string // Workflow registry family; nodes sync workflows registered under this key
	gatewayURL string // gateway URL for vault secret encryption (empty when no secrets)
}

// resolveWorkflowDeployTargets merges CLI flags with the topology saved by the last env start.
//
// Resolution order (later steps only fill fields the user did not override):
//  1. Start from --don-id, --don-family, --gateway-url flag values.
//  2. Resolve the workflow DON via selector → ResolveWorkflowDONMetadata (workflow_don_resolver.go).
//  3. If both --workflow-don-name and --don-family were explicitly set, validate they agree with state.
//  4. Default donID / donFamily from the selected DON's metadata when those flags were omitted.
//  5. finalizeWorkflowDonFamily — require don_family from state/flags when local CRE state exists.
//  6. Default gatewayURL from the gateway paired to that don_family (needed for --secrets-file-path).
//
// Multi-DON topologies must identify the workflow DON (--workflow-don-name, or --don-family
// with optional --shard-index when unambiguous). Container-name-pattern is copy-only.
//
// When resolver is nil (no local CRE state on disk), --don-family must be passed explicitly;
// donID and gatewayURL stay at whatever the caller passed.
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
		family, err := finalizeWorkflowDonFamily(donFamilyFlag)
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

	if validateErr := validateWorkflowDeployFlags(cmd, donMeta, selector, donFamilyFlag); validateErr != nil {
		return targets, validateErr
	}

	if !cmd.Flags().Changed("don-id") {
		targets.donID = libc.MustSafeUint32FromUint64(donMeta.ID)
	}
	// Infer registry family from the selected DON when --don-family was omitted.
	if !cmd.Flags().Changed("don-family") && donMeta.DonFamily != "" {
		targets.donFamily = donMeta.DonFamily
	}

	family, err := finalizeWorkflowDonFamily(targets.donFamily)
	if err != nil {
		return targets, err
	}
	targets.donFamily = family

	// Family-scoped gateway for vault secret encryption; skip when operator passed --gateway-url.
	if !cmd.Flags().Changed("gateway-url") {
		url, gwErr := resolver.GatewayURLForDonFamily(family)
		if gwErr != nil {
			return targets, fmt.Errorf("❌ failed to resolve gateway URL for don_family %q: %w", family, gwErr)
		}
		targets.gatewayURL = url
	}

	return targets, nil
}

// finalizeWorkflowDonFamily returns the don_family string written to the workflow registry on deploy.
// Must match nodesets.don_family from env start or be passed explicitly via --don-family.
func finalizeWorkflowDonFamily(donFamily string) (string, error) {
	if donFamily != "" {
		return donFamily, nil
	}
	return "", errors.New("❌ --don-family is required (set nodesets.don_family in the topology or pass --don-family / --workflow-don-name with local CRE state)")
}

// validateWorkflowDeployFlags is a deploy-time guard when the operator sets both identity flags explicitly.
//
// We only cross-check when --workflow-don-name and --don-family were both passed on the CLI. If only one
// is set, the other is inferred from local CRE state and there is nothing to disagree about. When both
// are set, a typo (e.g. --workflow-don-name feeds-zone-a with --don-family feeds-zone-b) would register
// the workflow under the wrong family while nodes and gateway connectors expect the topology value — fail
// before UpsertWorkflow instead of leaving a silent mismatch.
func validateWorkflowDeployFlags(cmd *cobra.Command, donMeta *cre.DonMetadata, selector workflowDONSelector, donFamilyFlag string) error {
	if !cmd.Flags().Changed("workflow-don-name") || !cmd.Flags().Changed("don-family") {
		return nil
	}
	if donMeta.DonFamily == "" {
		return nil
	}
	// CLI don_family must match the selected DON's topology value.
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
