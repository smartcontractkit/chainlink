// Package http sets up the HTTP capabilities and the gateway they talk to.
//
// It replaces the two features that set them up separately. They were separate
// because they were separate binaries in separate LOOPs; they are one binary now,
// sharing one connection to a gateway, so what starts them is one feature.
//
// The gateway is no longer part of a node. It is its own binary, started on the
// gateway DON's node the same way a capability is - as a capabilityrunner job -
// and reached over HTTP rather than a websocket. What a customer sends it is
// unchanged.
package http

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
)

// flag is the capability this feature is enabled by. One flag for both
// capabilities: they are one binary, and a DON that wants the trigger without the
// action is not a shape this serves.
const flag = cre.HTTPTriggerCapability

// The addresses the binaries use.
//
// They are not 5002 and 5003, which is where a gateway lives in this topology,
// because the node's own gateway is still there serving what has not moved -
// vault, the confidential relay, the DAG-era web API - and two gateways cannot
// share a port. So this one sits beside it until there is nothing left in the
// old one.
//
// The health port is what the job that launches this polls; it only has to avoid
// crecore's 50051/50052 and the capability runner's 50053.
const (
	gatewayUserPort   = 5012
	gatewayNodePort   = 5013
	gatewayHealthPort = 50054
	capabilityPort    = 50055

	// creCoreGRPCTarget is crecore: the registry these capabilities announce
	// themselves to, and the keystore they sign as this node through.
	creCoreGRPCTarget = "localhost:50051"

	// databaseSchema is where this capability keeps its own tables in the node's
	// database. Created by the node's migrations, so this binary needs no right to
	// make one.
	databaseSchema = "http_capability"

	// gatewayID is what the nodes authenticate to. It matches what the topology
	// names the gateway node, so that a node's header names the gateway it reached.
	gatewayID = "gateway-node-0"
)

// httpFlagPrefix is the namespace the HTTP binary registers its own settings
// under, so a value from capability_defaults.toml becomes a flag.
const httpFlagPrefix = "--"

// valueFlags maps the keys accepted under [capability_configs.http.values] to the
// binary's flag names. Unknown keys are rejected rather than dropped: a value set
// in the TOML and silently ignored looks like the binary disagreeing with its
// configuration.
var valueFlags = map[string]string{
	"SendChannelBufferSize":        "trigger.send-channel-buffer-size",
	"RequestCacheTTL":              "trigger.request-cache-ttl",
	"MaxAuthorizedKeysPerWorkflow": "trigger.max-authorized-keys-per-workflow",
	"ProxyMode":                    "action.proxy-mode",
}

type HTTP struct{}

func (o *HTTP) Flag() cre.CapabilityFlag { return flag }

// PreEnvStartup registers both capabilities with the registry.
//
// No gateway handler is added to any node: the gateway is a binary of its own
// now, and the node that used to serve it does not.
func (o *HTTP) PreEnvStartup(
	_ context.Context,
	_ zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	_ *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: []keystone_changeset.DONCapabilityWithConfig{
			{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "http-trigger",
					Version:        "1.0.0-alpha",
					CapabilityType: 0, // TRIGGER
				},
				Config: &capabilitiespb.CapabilityConfig{LocalOnly: don.HasOnlyLocalCapabilities()},
			},
			{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "http-actions",
					Version:        "1.0.0-alpha",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{LocalOnly: don.HasOnlyLocalCapabilities()},
			},
		},
	}, nil
}

// PostEnvStartup starts the gateway, then the capabilities that talk to it.
func (o *HTTP) PostEnvStartup(
	ctx context.Context,
	_ zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	capabilityConfig, ok := don.GetCapabilityConfig(flag)
	if !ok {
		return fmt.Errorf("config for '%s' capability not found for %s DON", flag, don.GetName())
	}

	command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get the command for the HTTP capabilities")
	}

	workers, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	nodes, aErr := accounts(workers, creEnv)
	if aErr != nil {
		return aErr
	}

	gatewayNode, found := dons.Gateway()
	if !found {
		return errors.New("no gateway node in the topology, so the HTTP capabilities have nothing to connect to")
	}

	configFlags, fErr := buildConfigFlags(capabilityConfig)
	if fErr != nil {
		return fErr
	}

	specs := make(cre.DonJobs, 0, len(workers)+1)

	// The gateway first: the capabilities connect to it, and one that is not there
	// yet is one they retry until it is.
	gateway, gwErr := gatewayJob(don, gatewayNode, creEnv, nodes)
	if gwErr != nil {
		return gwErr
	}
	specs = append(specs, gateway)

	for _, node := range workers {
		spec, sErr := capabilityJob(node, don, gatewayNode, creEnv, command, configFlags)
		if sErr != nil {
			return sErr
		}
		specs = append(specs, spec)
	}

	if err := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to create HTTP capability jobs: %w", err)
	}
	return nil
}

// gatewayJob launches the gateway binary on the gateway DON's node.
//
// It runs there because that is where the addresses already point - the topology
// gives the gateway node 5002 for customers and 5003 for the DON - and because a
// gateway that moved would be a gateway every workflow's configuration has to be
// told about again.
func gatewayJob(don *cre.Don, gatewayNode *cre.Node, creEnv *cre.Environment, nodes []string) (*jobv1.ProposeJobRequest, error) {
	command, cErr := standardcapability.GetCommand("http-gateway")
	if cErr != nil {
		return nil, errors.Wrap(cErr, "failed to get the command for the gateway")
	}

	args := []string{
		"--gateway.gateway-id=" + gatewayID,
		"--gateway.don-id=" + don.Name,
		"--gateway.nodes=" + strings.Join(nodes, ","),
		// F as the DON is configured: a customer is answered when F+1 of its nodes have
		// said the same thing, which is what makes the answer the DON's rather than one
		// node's.
		fmt.Sprintf("--gateway.f=%d", don.F),
		fmt.Sprintf("--gateway.user-address=:%d", gatewayUserPort),
		fmt.Sprintf("--gateway.node-address=:%d", gatewayNodePort),
		fmt.Sprintf("--http.port=%d", gatewayHealthPort),
	}

	return capabilityRunner(gatewayNode, creEnv, "http-gateway", command, args)
}

// capabilityJob launches the HTTP capabilities on one worker node.
func capabilityJob(node *cre.Node, don *cre.Don, gatewayNode *cre.Node, creEnv *cre.Environment, command string, configFlags []string) (*jobv1.ProposeJobRequest, error) {
	account, aErr := account(node, creEnv)
	if aErr != nil {
		return nil, aErr
	}

	args := []string{
		// One instance, configured as it is here. The binary also has an "embed"
		// subcommand that runs several with a gateway of their own, which is for local
		// runs rather than for a node.
		"run",
		fmt.Sprintf("--http.port=%d", capabilityPort),
		"--capabilities.proxy-url=" + creCoreGRPCTarget,
		fmt.Sprintf("--capabilities.capability-don-id=%d", don.ID),
		// The keys are the node's: this binary signs as it without holding them, which
		// is what the gateway checks its handshake against.
		"--keystore.proxy-address=" + creCoreGRPCTarget,
		// Its own schema in the node's database, created by the node's migrations. What
		// it keeps there is the requests it has answered, so that a customer's retry is
		// answered rather than run a second time - which is what the node's own key-value
		// table did for this capability before it was a binary of its own.
		//
		// The URL is not passed: CL_DATABASE_URL is in this process's environment
		// already, and the binary reads it the same way it reads any other setting.
		"--database.schema=" + databaseSchema,
		"--gateway.node-address=" + account,
		"--gateway.don-id=" + don.Name,
		"--gateway.gateways=" + gatewayID + "=" + gatewayURL(gatewayNode),
	}
	args = append(args, configFlags...)

	return capabilityRunner(node, creEnv, "http", command, args)
}

// capabilityRunner is the job that launches a binary beside a node.
func capabilityRunner(node *cre.Node, creEnv *cre.Environment, name, command string, args []string) (*jobv1.ProposeJobRequest, error) {
	quoted := make([]string, 0, len(args))
	for _, a := range args {
		quoted = append(quoted, strconv.Quote(a))
	}

	externalJobID := uuid.NewString()
	if !creEnv.FreshExternalJobIDs {
		externalJobID = uuid.NewString()
	}

	return &jobv1.ProposeJobRequest{
		NodeId: node.JobDistributorDetails.NodeID,
		Spec: fmt.Sprintf(`
type = "capabilityrunner"
schemaVersion = 1
externalJobID = "%s"
name = "%s-%s"
command = "%s"
args = [%s]
`,
			externalJobID,
			name,
			node.Name,
			command,
			strings.Join(quoted, ", "),
		),
	}, nil
}

// gatewayURL is where the nodes reach the gateway from inside the environment:
// the gateway node's own hostname, which is what the topology hands a node for
// every other purpose.
func gatewayURL(gateway *cre.Node) string {
	return fmt.Sprintf("http://%s:%d", gateway.Host, gatewayNodePort)
}

// accounts are the DON's members, as the gateway knows them.
func accounts(nodes []*cre.Node, creEnv *cre.Environment) ([]string, error) {
	addresses := make([]string, 0, len(nodes))
	for _, node := range nodes {
		address, err := account(node, creEnv)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

// account is the node's EVM address on the registry chain: the identity the DON's
// membership is recorded under, and so the one it authenticates to a gateway with.
func account(node *cre.Node, creEnv *cre.Environment) (string, error) {
	chainID, err := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to resolve the registry chain ID for node %s: %w", node.Name, err)
	}
	if node.Keys == nil {
		return "", fmt.Errorf("node %s has no keys", node.Name)
	}
	key, ok := node.Keys.EVM[chainID]
	if !ok || key == nil {
		return "", fmt.Errorf("node %s has no EVM key for the registry chain %d, so it has no identity to authenticate with", node.Name, chainID)
	}
	// Checksummed, which is how a node's keystore stores it. The gateway compares
	// addresses case-insensitively, so this is about the keystore rather than about
	// the gateway.
	return key.PublicAddress.Hex(), nil
}

// buildConfigFlags turns [capability_configs.http.values] into flags.
func buildConfigFlags(capConfig cre.CapabilityConfig) ([]string, error) {
	keys := make([]string, 0, len(capConfig.Values))
	for k := range capConfig.Values {
		keys = append(keys, k)
	}
	// Sorted, so the args - and the job spec a node is proposed - are the same on
	// every run rather than following map iteration order.
	sort.Strings(keys)

	flags := make([]string, 0, len(keys))
	for _, k := range keys {
		name, ok := valueFlags[k]
		if !ok {
			return nil, fmt.Errorf("unknown HTTP capability config value %q; the binary accepts %s",
				k, strings.Join(sortedKeys(valueFlags), ", "))
		}
		flags = append(flags, fmt.Sprintf("%s%s=%v", httpFlagPrefix, name, capConfig.Values[k]))
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
