package cron

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const flag = cre.CronCapability

type Cron struct{}

func (c *Cron) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Cron) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "cron-trigger",
			Version:        "1.0.0",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

// The cron binary runs as a capability runner rather than as a standard
// capability launched from the registry: it hosts no node of its own, so the
// registry it announces itself to, and the CRE settings its schedule limits
// resolve against, both come from the crecore process beside it. That address is
// the one the node launches crecore with, so it must match the
// [Capabilities.Proxy] block in the topology config (GRPCPort), which is what
// enables crecore in the first place.
const (
	// creCoreGRPCTarget is crecore's single gRPC address, which is what serves the
	// capabilities registry this binary announces itself to.
	creCoreGRPCTarget = "localhost:50051"

	// runnerHTTPPort serves the runner's /metrics, /debug/pprof, health endpoints
	// and - the part the node needs - the settings reload endpoint it is notified
	// on. Loopback inside the node's container, so one port suffices for every
	// node; it only has to avoid crecore's 50051/50052 and whatever another
	// capability runner on the same node listens on (consensus takes 50053).
	runnerHTTPPort = 50054
)

// cronFlagPrefix is the namespace the cron binary registers its own settings
// under, so a value from capability_defaults.toml becomes a flag.
const cronFlagPrefix = "--cron."

// cronValueFlags maps the keys accepted under [capability_configs.cron.values]
// to the binary's flag names. Unknown keys are rejected rather than dropped: a
// value set in the TOML and silently ignored looks like the binary disagreeing
// with its configuration.
var cronValueFlags = map[string]string{
	"FastestScheduleIntervalSeconds": "fastest-schedule-interval-seconds",
}

func (c *Cron) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	if err := createJobs(ctx, don, dons, creEnv); err != nil {
		return fmt.Errorf("failed to create cron jobs: %w", err)
	}

	return nil
}

func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	capabilityConfig, ok := don.GetCapabilityConfig(flag)
	if !ok {
		return fmt.Errorf("config for '%s' capability not found for %s DON", flag, don.GetName())
	}

	command, commandErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if commandErr != nil {
		return fmt.Errorf("failed to get command for cron capability: %w", commandErr)
	}

	configFlags, configErr := buildConfigFlags(capabilityConfig)
	if configErr != nil {
		return configErr
	}

	specs := make(cre.DonJobs, 0, len(don.Nodes))
	for _, node := range don.Nodes {
		// Only worker nodes run the capability; the bootstrap node has no
		// capability DON membership and no workflow to fire a trigger for.
		if !node.HasRole(cre.RoleWorker) {
			continue
		}
		spec, specErr := capabilityRunnerJobSpec(node, don, command, configFlags)
		if specErr != nil {
			return specErr
		}
		specs = append(specs, spec)
	}
	if len(specs) == 0 {
		return fmt.Errorf("no worker nodes found in %s DON to run the cron capability", don.GetName())
	}

	if err := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to create cron jobs: %w", err)
	}

	return nil
}

// capabilityRunnerJobSpec builds the capabilityrunner job that launches the cron
// binary on one node.
//
// The node supervises the process over the empty LOOP and notifies it of CRE
// settings updates on --http.port; everything else the binary needs it is told
// here, since a process that hosts no node of its own cannot look any of it up:
// the registry it announces the capability to, and which DON it is announcing
// itself as part of.
func capabilityRunnerJobSpec(node *cre.Node, don *cre.Don, command string, configFlags []string) (*jobv1.ProposeJobRequest, error) {
	if node.JobDistributorDetails == nil {
		return nil, fmt.Errorf("node %s has no job distributor details", node.Name)
	}

	args := []string{
		// One instance, configured as it is here. The binary also has an "embed"
		// subcommand for running several in one process, which is for local runs
		// rather than for the node.
		"run",
		fmt.Sprintf("--http.port=%d", runnerHTTPPort),
		"--capabilities.proxy-url=" + creCoreGRPCTarget,
		fmt.Sprintf("--capabilities.capability-don-id=%d", don.ID),
	}
	args = append(args, configFlags...)

	quoted := make([]string, 0, len(args))
	for _, a := range args {
		quoted = append(quoted, fmt.Sprintf("%q", a))
	}

	return &jobv1.ProposeJobRequest{
		NodeId: node.JobDistributorDetails.NodeID,
		Spec: fmt.Sprintf(`
type = "capabilityrunner"
schemaVersion = 1
externalJobID = "%s"
name = "cron-worker"
command = "%s"
args = [%s]
`,
			uuid.NewString(),
			command,
			strings.Join(quoted, ", "),
		),
	}, nil
}

// buildConfigFlags turns [capability_configs.cron.values] into --cron.* flags. An
// empty set leaves every field at the binary's own default, which is what the
// commented-out capability_defaults.toml entries mean.
func buildConfigFlags(capConfig cre.CapabilityConfig) ([]string, error) {
	keys := make([]string, 0, len(capConfig.Values))
	for k := range capConfig.Values {
		keys = append(keys, k)
	}
	// Sorted so the args - and so the job spec a node is proposed - are the same
	// on every run rather than following map iteration order.
	sort.Strings(keys)

	flags := make([]string, 0, len(keys))
	for _, k := range keys {
		name, ok := cronValueFlags[k]
		if !ok {
			return nil, fmt.Errorf("unknown cron capability config value %q; the cron binary accepts %s",
				k, strings.Join(sortedKeys(cronValueFlags), ", "))
		}
		flags = append(flags, fmt.Sprintf("%s%s=%v", cronFlagPrefix, name, capConfig.Values[k]))
	}
	return flags, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
