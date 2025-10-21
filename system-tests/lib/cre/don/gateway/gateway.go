package gateway

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

type WhitelistConfig struct {
	ExtraAllowedPorts                    []int
	ExtraAllowedIPs, ExtraAllowedIPsCIDR []string
}

func JobConfigs(
	registryChainOutput *blockchain.Output,
	topology *cre.Topology,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	whitelistConfig WhitelistConfig,
) (map[cre.NodeUUID]*config.GatewayConfig, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	chainID, chErr := strconv.ParseUint(registryChainOutput.ChainID, 10, 64)
	if chErr != nil {
		return nil, errors.Wrap(chErr, "failed to parse chain ID")
	}

	// if we don't have a gateway connector outputs, it means that this topology does not require a gateway
	// so we can skip the rest of the setup
	if topology.GatewayConnectors == nil || len(topology.GatewayConnectors.Configurations) == 0 {
		return nil, nil
	}

	// for each gateway node prepare GatewayConfig, which will be later used in a job spec
	// by default we add only add web-api handler to the workflow DON (so that it can download workflows)
	// all other handlers should be added by capabilities/features that require them
	// gateway configurations that contain networking data are added, when a new topology is created
	result := make(map[string]*config.GatewayConfig)
	for _, donMetadata := range topology.DonsMetadata.List() {
		gateway, hasGateway := donMetadata.Gateway()
		if !hasGateway {
			continue
		}

		configuration, cErr := topology.GatewayConnectors.FindByNodeUUID(gateway.UUID)
		if cErr != nil {
			return nil, errors.Wrapf(cErr, "failed to find gateway configuration for node UUID %s", gateway.UUID)
		}

		c := config.GatewayConfig{
			ConnectionManagerConfig: config.ConnectionManagerConfig{
				AuthGatewayId:             configuration.AuthGatewayID,
				AuthChallengeLen:          10,
				AuthTimestampToleranceSec: 5,
				HeartbeatIntervalSec:      20,
			},
			NodeServerConfig: gw_net.WebSocketServerConfig{
				HandshakeTimeoutMillis: 1000,
				HTTPServerConfig: gw_net.HTTPServerConfig{
					MaxRequestBytes:      100_000,
					ReadTimeoutMillis:    1_000,
					RequestTimeoutMillis: 10_000,
					WriteTimeoutMillis:   1_000,
					Path:                 configuration.Outgoing.Path,
					Port:                 uint16(configuration.Outgoing.Port), //nolint:gosec //should never happen unless someone uses an incorrect negative port
				},
			},
			UserServerConfig: gw_net.HTTPServerConfig{
				ContentTypeHeader:    "application/jsonrpc",
				MaxRequestBytes:      100_000,
				ReadTimeoutMillis:    80_000,
				RequestTimeoutMillis: 80_000,
				WriteTimeoutMillis:   80_000,
				CORSEnabled:          false,
				CORSAllowedOrigins:   []string{},
				Path:                 configuration.Incoming.Path,
				Port:                 uint16(configuration.Incoming.InternalPort), //nolint:gosec //should never happen unless someone uses an incorrect negative port
			},
			HTTPClientConfig: gw_net.HTTPClientConfig{
				MaxResponseBytes: 100_000_000,
				AllowedPorts:     append(whitelistConfig.ExtraAllowedPorts, DefaultAllowedPorts...),
				AllowedIPs:       whitelistConfig.ExtraAllowedIPs,
				AllowedIPsCIDR:   whitelistConfig.ExtraAllowedIPsCIDR,
			},
		}

		workflowDON, donErr := topology.DonsMetadata.WorkflowDON()
		if donErr != nil {
			return nil, errors.Wrap(donErr, "failed to find workflow DON")
		}

		workerNodes, wErr := workflowDON.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		donConfig := config.DONConfig{
			DonId:   workflowDON.Name,
			F:       1,
			Members: make([]config.NodeConfig, len(workerNodes)),
		}

		for i, workerNode := range workerNodes {
			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			donConfig.Members[i] = config.NodeConfig{
				Address: evmKey.PublicAddress.Hex(),
				Name:    fmt.Sprintf("%s-node-%d", workflowDON.Name, i),
			}
		}

		handlerConfig, hErr := HandlerConfig(coregateway.WebAPICapabilitiesType)
		if hErr != nil {
			return nil, errors.Wrap(hErr, "failed to get web-api capability handler config")
		}

		donConfig.Handlers = []config.Handler{handlerConfig}
		c.Dons = append(c.Dons, donConfig)
		result[gateway.UUID] = &c
	}

	if len(result) == 0 {
		return nil, errors.New("no gateway configurations were created, although at least one is expected")
	}

	return result, nil
}

// CreateJobs creates gateway job spec for each gateway node in the DON topology and sends it to JD for creation and approval
func CreateJobs(ctx context.Context, jd *jd.JobDistributor, dons *cre.Dons, gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig) error {
	jobSpecs := make(cre.DonJobs, 0)

	header := `
type = "gateway"
schemaVersion = 1
externalJobID = "%s"
name = "cre-gateway"
forwardingAllowed = false
`

	for nodeUUID, gc := range gatewayConfigs {
		jobSpec := fmt.Sprintf(header, uuid.NewString())

		type wrapper struct {
			GC config.GatewayConfig `json:"gatewayConfig"  toml:"gatewayConfig"`
		}

		gatewayNode, found := dons.NodeWithUUID(nodeUUID)
		if !found {
			return fmt.Errorf("could not find gateway node with UUID %s in DON topology", nodeUUID)
		}

		tomlStr, mErr := toml.Marshal(wrapper{GC: *gc})
		if mErr != nil {
			return fmt.Errorf("failed to marshal gateway config to toml: %w", mErr)
		}

		// hack for json.RawMessage that otherwise outputs a byte array instead of JSON string in toml, which cannot be parsed by gateway
		replaced, rErr := expandConfigByteArray(string(tomlStr), []string{"gatewayConfig", "Dons", "Handlers", "Config"})
		if rErr != nil {
			return fmt.Errorf("failed to expand config byte arrays: %w", rErr)
		}

		jobSpec += "\n" + replaced
		jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
			NodeId: gatewayNode.JobDistributorDetails.NodeID,
			Spec:   jobSpec,
		})
	}

	return jobs.Create(ctx, jd, dons, jobSpecs)
}

// AddHandlers adds the given handler configurations to the gateway job config of the given DON. It only adds handlers, if they are not already present.
func AddHandlers(donMetadata cre.DonMetadata, registryChainID uint64, gatewayJobConfigs map[cre.NodeUUID]*config.GatewayConfig, handlerConfigs []config.Handler) error {
	workers, wErr := donMetadata.Workers()
	if wErr != nil {
		return wErr
	}
	evmKey, ok := workers[0].Keys.EVM[registryChainID]
	if !ok {
		return fmt.Errorf("worker node at index %d does not have EVM key for chainID %d", workers[0].Index, registryChainID)
	}

	// for each DON, we need to add a handler config specific for this capability
	for _, gc := range gatewayJobConfigs {
		donFound := false
		for donIdx, maybeDON := range gc.Dons {
			// first we try to find DON configuration that matches current don, because it might be already present
			for _, member := range maybeDON.Members {
				// if any of the member's address matches the EVM key of the worker node, we found the right DON
				if member.Address == evmKey.PublicAddress.Hex() {
					donFound = true
					break
				}
			}

			if donFound {
				for _, newHandler := range handlerConfigs {
					alreadyPresent := false
					for _, existingHandlers := range maybeDON.Handlers {
						if strings.EqualFold(existingHandlers.Name, newHandler.Name) {
							alreadyPresent = true
							break
						}
					}
					if !alreadyPresent {
						gc.Dons[donIdx].Handlers = append(gc.Dons[donIdx].Handlers, newHandler)
					}
				}
				break
			}
		}

		// if we did not find the DON in the gateway config, we need to add it
		if !donFound {
			members := make([]config.NodeConfig, len(workers))
			for i, worker := range workers {
				evmKey, ok := worker.Keys.EVM[registryChainID]
				if !ok {
					return fmt.Errorf("worker node at index %d does not have EVM key for chain ID %d", worker.Index, registryChainID)
				}

				members[i] = config.NodeConfig{
					Address: evmKey.PublicAddress.Hex(),
					Name:    fmt.Sprintf("%s-node-%d", donMetadata.Name, worker.Index),
				}
			}

			gc.Dons = append(gc.Dons, config.DONConfig{
				DonId:    donMetadata.Name,
				F:        1,
				Members:  members,
				Handlers: handlerConfigs,
			})
		}
	}

	return nil
}

// AddConnectors adds gateway connector configuration to the node TOML config of each node in the given DON. It only adds connectors, if they are not already present.
func AddConnectors(donMetadata *cre.DonMetadata, registryChainID uint64, connectors cre.GatewayConnectors) error {
	workers, wErr := donMetadata.Workers()
	if wErr != nil {
		return wErr
	}

	for _, workerNode := range workers {
		currentConfig := donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

		var typedConfig corechainlink.Config
		unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
		if unmarshallErr != nil {
			return errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
		}

		evmKey, ok := workerNode.Keys.EVM[registryChainID]
		if !ok {
			return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", registryChainID, workerNode.Index)
		}

		// if no gateways are configured, then gateway connector config is most probably also not configured
		if len(typedConfig.Capabilities.GatewayConnector.Gateways) == 0 {
			typedConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
				DonID:             ptr.Ptr(donMetadata.Name),
				ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(registryChainID, 10)),
				NodeAddress:       ptr.Ptr(evmKey.PublicAddress.Hex()),
			}
		}

		// make sure that all other gateways are also present in the config
		for _, gatewayConnector := range connectors.Configurations {
			alreadyPresent := false
			for _, existingGateway := range typedConfig.Capabilities.GatewayConnector.Gateways {
				if gatewayConnector.AuthGatewayID == *existingGateway.ID {
					alreadyPresent = true
					continue
				}
			}

			if !alreadyPresent {
				typedConfig.Capabilities.GatewayConnector.Gateways = append(typedConfig.Capabilities.GatewayConnector.Gateways, coretoml.ConnectorGateway{
					ID: ptr.Ptr(gatewayConnector.AuthGatewayID),
					URL: ptr.Ptr(fmt.Sprintf("ws://%s:%d%s",
						gatewayConnector.Outgoing.Host,
						gatewayConnector.Outgoing.Port,
						gatewayConnector.Outgoing.Path)),
				})
			}
		}

		stringifiedConfig, mErr := toml.Marshal(typedConfig)
		if mErr != nil {
			return errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
		}

		donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
	}

	return nil
}

func HandlerConfig(handler string) (config.Handler, error) {
	switch handler {
	case coregateway.HTTPCapabilityType:
		return config.Handler{
			Name:        coregateway.HTTPCapabilityType,
			ServiceName: "workflows",
			Config: []byte(`
maxTriggerRequestDurationMs = 5_000
metadataPullIntervalMs = 1_000
metadataAggregationIntervalMs = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`),
		}, nil
	case coregateway.WebAPICapabilitiesType:
		return config.Handler{
			Name: coregateway.WebAPICapabilitiesType,
			Config: []byte(`
maxAllowedMessageAgeSec = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`),
		}, nil
	case coregateway.VaultHandlerType:
		return config.Handler{
			Name:        coregateway.VaultHandlerType,
			ServiceName: "vault",
			Config: []byte(`
requestTimeoutSec = 70
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`)}, nil
	default:
		return config.Handler{}, fmt.Errorf("unknown handler type: %s", handler)
	}
}

// ExpandConfigByteArray finds lines like `Config = [10, 109, ...]` and replaces them
// with TOML tables under the given path, using the bytes as TOML text.
// Example path: []string{"gatewayConfig","Dons","Handlers","Config"}
func expandConfigByteArray(tomlDoc string, path []string) (string, error) {
	re := regexp.MustCompile(`(?m)^(\s*)Config\s*=\s*\[([0-9,\s]+)\]\s*$`)
	return re.ReplaceAllStringFunc(tomlDoc, func(line string) string {
		m := re.FindStringSubmatch(line)
		if m == nil {
			return line
		}
		indent := m[1]
		nums := m[2]

		// parse the byte array back to text
		bs, err := parseByteArray(nums)
		if err != nil {
			// if parsing fails, keep original line to avoid breaking output
			return line
		}
		snippet := string(bs)

		// rewrite snippet under full path
		expanded := embedUnderPath(snippet, path)

		// keep the same indentation before each emitted line
		var out strings.Builder
		for _, l := range strings.Split(expanded, "\n") {
			if len(strings.TrimSpace(l)) == 0 {
				out.WriteString("\n")
				continue
			}
			out.WriteString(indent)
			out.WriteString(l)
			out.WriteString("\n")
		}
		return strings.TrimRight(out.String(), "\n")
	}), nil
}

func parseByteArray(s string) ([]byte, error) {
	fields := strings.Split(s, ",")
	buf := bytes.NewBuffer(nil)
	for _, f := range fields {
		t := strings.TrimSpace(f)
		if t == "" {
			continue
		}
		n, err := strconv.Atoi(t)
		if err != nil || n < 0 || n > 255 {
			return nil, fmt.Errorf("invalid byte: %q", t)
		}
		buf.WriteByte(byte(n))
	}
	return buf.Bytes(), nil
}

// embedUnderPath prefixes any table headers in snippet with path,
// and adds a `[path...]` header for the root keys.
func embedUnderPath(snippet string, path []string) string {
	base := "[" + strings.Join(path, ".") + "]"
	var out strings.Builder
	out.WriteString(base)
	out.WriteString("\n")

	for _, raw := range strings.Split(strings.ReplaceAll(snippet, "\r\n", "\n"), "\n") {
		line := strings.TrimRight(raw, " \t")
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}
		// table header inside snippet? e.g. [NodeRateLimiter]
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inner := strings.TrimSuffix(strings.TrimPrefix(trim, "["), "]")
			out.WriteString("[" + strings.Join(path, ".") + "." + inner + "]\n")
			continue
		}
		// regular key/value line
		out.WriteString(line)
		out.WriteString("\n")
	}
	return strings.TrimRight(out.String(), "\n")
}
