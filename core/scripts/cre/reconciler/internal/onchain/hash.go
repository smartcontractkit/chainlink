package onchain

import (
	"sort"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// Phase-hash keys stored in StateFile.PhaseHashes. PreEnvStartup and Configure
// CapReg share one key: PreEnvStartup's output (capability configs) is consumed
// in-memory, directly, by Configure CapReg, and can't be persisted cleanly (it
// carries a protobuf message with an unexported oneof field) — so the two are
// treated as one hash-gated unit, always run or skipped together.
const (
	phaseKeyPreEnvStartupCapReg  = "pre_env_startup_capreg"
	phaseKeyConfigureWorkflowReg = "configure_workflowreg"
	phaseKeyJobs                 = "jobs"
)

// donHashInput is the canonical, hashable shape of a DON's desired-state
// declaration plus its discovered member facts. Slices are explicitly sorted
// so reordering entries in desired.toml (or map iteration order) never
// changes the hash; encoding/json already sorts map keys on its own.
type donHashInput struct {
	Name                   string
	DONTypes               []string
	Capabilities           []string
	CapabilityConfigs      map[string]domain.CapabilityConfig
	RegistryBasedAllowlist []string
	ExposesRemoteCaps      bool
	Members                []memberHashInput
}

// memberHashInput is one DON member's discovered identity, sorted into
// donHashInput.Members by P2PID.
type memberHashInput struct {
	P2PID         string
	CSAKey        string
	EVMAddress    map[string]string
	OCR2BundleIDs map[string]string
}

func sortedCopy(s []string) []string {
	out := append([]string(nil), s...)
	sort.Strings(out)
	return out
}

// buildDonHashInput builds the canonical hash-input for one DON, combining its
// built topology metadata (for member/role data, already resolved) with its
// desired-state declaration (capabilities/config/allowlist). don may be nil
// (e.g. a synthetic gateway nodeset with no matching desired.DON entry).
func buildDonHashInput(donMeta *cre.DonMetadata, don *domain.DON, runtime map[string]domain.NodeRuntimeInfo) donHashInput {
	input := donHashInput{Name: donMeta.Name}
	if don != nil {
		input.DONTypes = sortedCopy(don.DONTypes)
		input.Capabilities = sortedCopy(don.Capabilities)
		input.CapabilityConfigs = don.CapabilityConfigs
		input.RegistryBasedAllowlist = sortedCopy(don.RegistryBasedAllowlist)
		input.ExposesRemoteCaps = don.ExposesRemoteCaps
	}

	workers, err := donMeta.Workers()
	if err != nil {
		// No worker-role nodes on this DON (bootstrap-only/gateway-only) — not
		// an error, just nothing to add to Members.
		return input
	}
	for _, worker := range workers {
		if worker == nil || worker.Keys == nil {
			continue
		}
		nodeName := chartNodeNameForWorker(worker, runtime)
		nodeRuntime := runtime[nodeName]
		input.Members = append(input.Members, memberHashInput{
			P2PID:         worker.Keys.PeerID(),
			CSAKey:        nodeRuntime.CSAKey,
			EVMAddress:    nodeRuntime.EVMAddress,
			OCR2BundleIDs: nodeRuntime.OCR2BundleIDs,
		})
	}
	sort.Slice(input.Members, func(i, j int) bool { return input.Members[i].P2PID < input.Members[j].P2PID })

	return input
}

// donHashInputs walks every DON in topology, applying exclude (if non-nil) to
// skip DONs that shouldn't factor into a given phase's hash, and returns them
// sorted by name.
func donHashInputs(
	desired *domain.DesiredState,
	topology *cre.Topology,
	runtime map[string]domain.NodeRuntimeInfo,
	exclude func(donMeta *cre.DonMetadata) bool,
) []donHashInput {
	if topology == nil || topology.DonsMetadata == nil {
		return nil
	}

	var out []donHashInput
	for _, donMeta := range topology.DonsMetadata.List() {
		if exclude != nil && exclude(donMeta) {
			continue
		}
		don := desired.DONByName(donMeta.Name)
		out = append(out, buildDonHashInput(donMeta, don, runtime))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// capRegExcludedDON mirrors capRegWorkerNodes' bootstrap/gateway-only
// exclusion, so the hash reflects exactly the DON scope Configure CapReg (and
// PreEnvStartup, its permanently-coupled predecessor) actually acts on.
func capRegExcludedDON(donMeta *cre.DonMetadata) bool {
	return flags.HasNoOtherFlags(donMeta.Flags, []string{cre.GatewayDON, cre.BootstrapDON})
}

// preEnvStartupCapRegHashInput is the canonical input for the combined
// PreEnvStartup + Configure CapReg unit.
type preEnvStartupCapRegHashInput struct {
	RegistryChainID         uint64
	GlobalCapabilityConfigs map[string]domain.CapabilityConfig
	DONs                    []donHashInput
}

func hashPreEnvStartupCapReg(
	desired *domain.DesiredState,
	topology *cre.Topology,
	runtime map[string]domain.NodeRuntimeInfo,
) (string, error) {
	registryChain, _ := desired.RegistryChain()
	input := preEnvStartupCapRegHashInput{
		RegistryChainID:         registryChain.ChainID,
		GlobalCapabilityConfigs: desired.CapabilityConfigs,
		DONs:                    donHashInputs(desired, topology, runtime, capRegExcludedDON),
	}
	return domain.CanonicalHash(input)
}

// configureWorkflowRegHashInput is the canonical input for Configure
// WorkflowReg: which DONs are workflow DONs, and the owner address that will
// be registered (derived from the deployer key, not the raw key itself).
type configureWorkflowRegHashInput struct {
	WorkflowDONNames []string
	WorkflowOwner    string
}

func hashConfigureWorkflowReg(desired *domain.DesiredState, workflowOwner string) (string, error) {
	var names []string
	for _, don := range desired.DONs {
		if don.IsWorkflowDon() {
			names = append(names, don.Name)
		}
	}
	sort.Strings(names)

	input := configureWorkflowRegHashInput{
		WorkflowDONNames: names,
		WorkflowOwner:    workflowOwner,
	}
	return domain.CanonicalHash(input)
}

// jobsHashInput is the canonical input for the Jobs (job sync) phase: every
// DON regardless of type (jobs are created for bootstrap/gateway DONs too,
// unlike CapReg) plus the persisted gateway service configs.
type jobsHashInput struct {
	DONs                  []donHashInput
	GatewayServiceConfigs []domain.GatewayServiceConfigState
}

func hashJobs(
	desired *domain.DesiredState,
	topology *cre.Topology,
	runtime map[string]domain.NodeRuntimeInfo,
	gatewayServiceConfigs []domain.GatewayServiceConfigState,
) (string, error) {
	configs := append([]domain.GatewayServiceConfigState(nil), gatewayServiceConfigs...)
	sort.Slice(configs, func(i, j int) bool { return configs[i].ServiceName < configs[j].ServiceName })

	input := jobsHashInput{
		DONs:                  donHashInputs(desired, topology, runtime, nil),
		GatewayServiceConfigs: configs,
	}
	return domain.CanonicalHash(input)
}
