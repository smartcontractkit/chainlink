package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

const flag = cre.ConsensusCapability
const consensusLabelledName = "consensus"

type Consensus struct{}

func (c *Consensus) Flag() cre.CapabilityFlag {
	return flag
}

func (c *Consensus) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   consensusLabelledName,
			Version:        "1.0.0-alpha",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		Config: &capabilitiespb.CapabilityConfig{
			LocalOnly: don.HasOnlyLocalCapabilities(),
		},
		UseCapRegOCRConfig: true,
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config: map[string]*ocr3.OracleConfig{
			consensusLabelledName: contracts.DefaultOCR3Config(),
		},
		CapabilityToExtraSignerFamilies: cre.CapabilityToExtraSignerFamilies(
			cre.OCRExtraSignerFamilies(creEnv.Blockchains),
			consensusLabelledName,
		),
	}, nil
}

const ContractQualifier = "capability_consensus"

// The consensus binary runs as a capability runner rather than as a standard
// capability: it hosts no node of its own, so its rage networking, the keys it
// signs with and the OCR configuration it runs under all come from the crecore
// process beside it. Those two addresses are the ones the node launches crecore
// with, so they must match the [Capabilities.Proxy] block in the topology config
// (GRPCPort), which is what enables crecore in the first place.
const (
	// creCoreGRPCTarget is crecore's single gRPC address: the OCR proxy (rage
	// networking plus signing) and the capabilities registry are both served on it.
	creCoreGRPCTarget = "localhost:50051"

	// runnerHTTPPort serves the runner's /metrics, /debug/pprof, health endpoints
	// and - the part the node needs - the settings reload endpoint it is notified
	// on. Loopback inside the node's container, so one port suffices for every
	// node; it only has to avoid crecore's 50051/50052.
	runnerHTTPPort = 50053
)

// consensusFlagPrefix is the namespace the consensus binary registers its own
// settings under, so a value from capability_defaults.toml becomes a flag.
const consensusFlagPrefix = "--consensus."

// consensusValueFlags maps the keys accepted under
// [capability_configs.consensus.values] to the binary's flag names. Unknown keys
// are rejected rather than dropped: a value set in the TOML and silently ignored
// looks like the binary disagreeing with its configuration.
var consensusValueFlags = map[string]string{
	"RequestBatchSize":             "request-batch-size",
	"MaxRequestSizeBytes":          "max-request-size-bytes",
	"MaxRequestOutcomeSize":        "max-request-outcome-size",
	"KeyBundleIDForValueConsensus": "key-bundle-id-for-value-consensus",
}

func (c *Consensus) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return fmt.Errorf("failed to create consensus jobs: %w", jobsErr)
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
		return fmt.Errorf("failed to get command for consensus capability: %w", commandErr)
	}

	configFlags, configErr := buildConfigFlags(capabilityConfig)
	if configErr != nil {
		return configErr
	}

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	specs := make(cre.DonJobs, 0, len(don.Nodes))
	for _, node := range don.Nodes {
		// Only worker nodes run the capability; the bootstrap node has no capability
		// DON membership and nothing to reach consensus over.
		if !node.HasRole(cre.RoleWorker) {
			continue
		}
		spec, specErr := capabilityRunnerJobSpec(node, don, creEnv, command, formatBootstrapPeer(bootstrap), configFlags)
		if specErr != nil {
			return specErr
		}
		specs = append(specs, spec)
	}
	if len(specs) == 0 {
		return fmt.Errorf("no worker nodes found in %s DON to run the consensus capability", don.GetName())
	}

	if err := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to create consensus jobs: %w", err)
	}

	return nil
}

// capabilityRunnerJobSpec builds the capabilityrunner job that launches the
// consensus binary on one node.
//
// The node supervises the process over the empty LOOP and notifies it of CRE
// settings updates on --http.port; everything else the binary needs it is told
// here, since a process that hosts no node of its own cannot look any of it up:
// the peer whose identity crecore holds on its behalf, the registry that says
// which DON it is and what OCR configuration that DON runs under, and the peers
// to dial before it has heard of anyone.
//
// The registry is one address serving both: --capabilities.proxy-url is where the
// capabilities and the OCR configuration are read from, and --ocr.proxy-address is
// where this node's peer and keys live. crecore serves both today, so they are the
// same value.
func capabilityRunnerJobSpec(node *cre.Node, don *cre.Don, creEnv *cre.Environment, command, bootstrapPeer string, configFlags []string) (*jobv1.ProposeJobRequest, error) {
	if node.JobDistributorDetails == nil {
		return nil, fmt.Errorf("node %s has no job distributor details", node.Name)
	}
	peerID := strings.TrimPrefix(node.PeerID(), "p2p_")
	if peerID == "" {
		return nil, fmt.Errorf("node %s has no P2P peer ID", node.Name)
	}
	transmitAccount, err := transmitAccountFor(node, creEnv)
	if err != nil {
		return nil, err
	}

	args := []string{
		// One instance, configured as it is here. The binary also has an "embed"
		// subcommand for running several in one process, which is for local runs
		// rather than for the node.
		"run",
		fmt.Sprintf("--http.port=%d", runnerHTTPPort),
		"--ocr.proxy-address=" + creCoreGRPCTarget,
		"--ocr.peer-id=" + peerID,
		"--ocr.transmit-account=" + transmitAccount,
		"--capabilities.proxy-url=" + creCoreGRPCTarget,
		fmt.Sprintf("--capabilities.capability-don-id=%d", don.ID),
		// Where the DON can first be reached belongs to the networking rather than to the
		// capability, so it is an ocr.* setting beside the proxy that does the dialling.
		"--ocr.bootstrappers=" + bootstrapPeer,
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
name = "consensus-worker"
command = "%s"
args = [%s]
`,
			uuid.NewString(),
			command,
			strings.Join(quoted, ", "),
		),
	}, nil
}

// buildConfigFlags turns [capability_configs.consensus.values] into
// --consensus.* flags. An empty set leaves every field at the binary's own
// default, which is what the commented-out capability_defaults.toml entries mean.
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
		name, ok := consensusValueFlags[k]
		if !ok {
			return nil, fmt.Errorf("unknown consensus capability config value %q; the consensus binary accepts %s",
				k, strings.Join(sortedKeys(consensusValueFlags), ", "))
		}
		flags = append(flags, fmt.Sprintf("%s%s=%v", consensusFlagPrefix, name, capConfig.Values[k]))
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

// transmitAccountFor returns the account the node is registered to transmit from:
// its EVM address on the registry chain, as lowercase hex with no 0x prefix.
//
// The presentation is the point, not just the address. libocr compares this to the
// config as an account string, and the config's copy is produced by whoever decodes
// it: for a capability read out of the CapabilitiesRegistry that is
// capabilitiespb.OCR3ConfigFromProto, which renders each on-chain transmitter as
// hex.EncodeToString of its 20 bytes. So this renders it the same way rather than
// reading it back from JD, which reports the EIP-55 form and would never match.
func transmitAccountFor(node *cre.Node, creEnv *cre.Environment) (string, error) {
	chainID, err := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to resolve the registry chain ID for node %s: %w", node.Name, err)
	}
	if node.Keys == nil {
		return "", fmt.Errorf("node %s has no keys", node.Name)
	}
	key, ok := node.Keys.EVM[chainID]
	if !ok || key == nil {
		return "", fmt.Errorf("node %s has no EVM key for the registry chain %d, so it has no account to transmit from", node.Name, chainID)
	}
	return hex.EncodeToString(key.PublicAddress.Bytes()), nil
}

func formatBootstrapPeer(bootstrap *cre.Node) string {
	return fmt.Sprintf("%s@%s:%d",
		strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"),
		bootstrap.Host,
		cre.OCRPeeringPort)
}
