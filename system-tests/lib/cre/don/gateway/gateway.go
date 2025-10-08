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

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

var (
	DefaultAllowedPorts = []int{80, 443}
)

func GatewayConfig(
	cldEnvironment *cldf.Environment,
	blockchainOutput *blockchain.Output,
	topology *cre.Topology,
	infraInput infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string,
) (map[cre.NodeUUID]*config.GatewayConfig, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	chainID, chErr := strconv.ParseUint(blockchainOutput.ChainID, 10, 64)
	if chErr != nil {
		return nil, errors.Wrap(chErr, "failed to parse chain ID")
	}

	// ALWAYS REQUIRE A GATEWAY
	// if we don't have a gateway connector outputs, we don't need to create any job specs
	// if topology.GatewayConnectorOutput == nil || len(topology.GatewayConnectorOutput.Configurations) == 0 {
	// 	return nil, nil
	// }

	// for each gateway node prepare GatewayConfig, which will be later used in a job spec
	// by default we add only add web-api handler to the workflow DON (so that it can download workflows)
	// all other handlers should be added by capabilities/features that require them
	result := make(map[string]*config.GatewayConfig)
	for _, donMetadata := range topology.DonsMetadata.List() {
		gateway, hasGateway := donMetadata.Gateway()
		if !hasGateway {
			continue
		}

		configuration, cErr := topology.GatewayConnectorOutput.FindByNodeUUID(gateway.UUID)
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
					Port:                 uint16(configuration.Outgoing.Port),
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
				Port:                 uint16(configuration.Incoming.InternalPort),
			},
			HTTPClientConfig: gw_net.HTTPClientConfig{
				MaxResponseBytes: 100_000_000,
				AllowedPorts:     append(extraAllowedPorts, DefaultAllowedPorts...),
				AllowedIPs:       extraAllowedIPs,
				AllowedIPsCIDR:   extraAllowedIPsCIDR,
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

		donConfig.Handlers = []config.Handler{
			{
				Name: coregateway.WebAPICapabilitiesType,
				Config: []byte(`
maxAllowedMessageAgeSec = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`),
			},
		}
		c.Dons = append(c.Dons, donConfig)
		result[gateway.UUID] = &c
	}

	if len(result) == 0 {
		return nil, errors.New("no gateway configurations were created, although at least one is expected")
	}

	return result, nil
}

func CreateJobs(ctx context.Context, jd *jd.JobDistributor, donTopology *cre.DonTopology, gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig) error {
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

		gatewayNode, found := donTopology.Dons.NodeWithUUID(nodeUUID)
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

	return jobs.Create(ctx, jd, donTopology, jobSpecs)
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
