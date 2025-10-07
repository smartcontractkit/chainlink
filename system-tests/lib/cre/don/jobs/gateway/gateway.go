package gateway

import (
	"fmt"
	"maps"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func JobSpec(extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string) cre.JobSpecFn {
	return func(input *cre.JobSpecInput) (cre.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}

		donToJobSpecs := make(cre.DonsToJobSpecs)

		// if we don't have a gateway connector outputs, we don't need to create any job specs
		// GatewayConnectorOutput is added by `system-tests/lib/cre/don/don.go`.BuildTopology() function, which builds gateway configuration
		// based on DON flags (cre.GatewayDON) and `gateway_node_index` and adds it to the topology.
		// `system-tests/lib/cre/don/don.go`.ValidateTopology() makes sure that if any DON needs gateway connector, there is at least one nodeSet with a gateway node.
		if input.DonTopology.GatewayConnectorOutput == nil || len(input.DonTopology.GatewayConnectorOutput.Configurations) == 0 {
			return donToJobSpecs, nil
		}

		// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
		// This map will be used to configure the gateway job on the node that runs it.
		for _, don := range input.DonTopology.Dons.List() {
			// if it's a workflow DON or has capabilities that require to connect to external resources, it needs access to gateway connector
			if !don.HasFlag(cre.WorkflowDON) && !don.RequiresGateway() {
				continue
			}

			workerNode, wErr := don.Workers()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			ethAddresses := make([]string, len(workerNode))
			chainID, err := chainselectors.ChainIdFromSelector(input.DonTopology.HomeChainSelector)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get home chain id from selector")
			}
			for i, workerNode := range workerNode {
				evmKey, ok := workerNode.Keys.EVM[chainID]
				if !ok {
					return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
				}
				ethAddresses[i] = evmKey.PublicAddress.Hex()
			}

			handlers := map[string]string{}
			if don.HasFlag(cre.WorkflowDON) || don.RequiresWebAPI() {
				handlerConfig := `
				[gatewayConfig.Dons.Handlers.Config]
				maxAllowedMessageAgeSec = 1_000
				[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
				globalBurst = 10
				globalRPS = 50
				perSenderBurst = 10
				perSenderRPS = 10
				`
				handlers[coregateway.WebAPICapabilitiesType] = handlerConfig
			}

			for _, capability := range input.Capabilities {
				if capability.GatewayJobHandlerConfigFn() == nil {
					continue
				}

				handlerConfig, handlerConfigErr := capability.GatewayJobHandlerConfigFn()(don)
				if handlerConfigErr != nil {
					return nil, errors.Wrap(handlerConfigErr, "failed to get handler config")
				}
				maps.Copy(handlers, handlerConfig)
			}

			for idx := range input.DonTopology.GatewayConnectorOutput.Configurations {
				// determine here what handlers we want to build.
				input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons = append(input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
					MembersEthAddresses: ethAddresses,
					ID:                  don.Name,
					Handlers:            handlers,
				})
			}
		}

		// we know that at least one DON must be the gateway DON, because topology.validate() checks that
		hasGateway := false
		for _, don := range input.DonTopology.Dons.List() {
			var gatewayNode *cre.Node
			gatewayNode, hasGateway = don.Gateway()
			if !hasGateway {
				continue
			}

			homeChainID, homeChainErr := chainselectors.ChainIdFromSelector(input.DonTopology.HomeChainSelector)
			if homeChainErr != nil {
				return nil, errors.Wrap(homeChainErr, "failed to get home chain id from selector")
			}

			for _, gatewayConfiguration := range input.DonTopology.GatewayConnectorOutput.Configurations {
				donToJobSpecs[don.ID] = append(donToJobSpecs[don.ID], jobs.AnyGateway(gatewayNode.JobDistributorDetails.NodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConfiguration))
			}
		}

		if !hasGateway {
			return nil, errors.New("no gateway node found in any DON, but at least one is required")
		}

		return donToJobSpecs, nil
	}
}

func GatewayConfig(
	cldEnvironment *cldf.Environment,
	blockchainOutput *blockchain.Output,
	donTopology *cre.DonTopology,
	infraInput infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string,
) (map[uint64][]config.GatewayConfig, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}

	// donToJobSpecs := make(cre.DonsToJobSpecs)

	// if we don't have a gateway connector outputs, we don't need to create any job specs
	// GatewayConnectorOutput is added by `system-tests/lib/cre/don/don.go`.BuildTopology() function, which builds gateway configuration
	// based on DON flags (cre.GatewayDON) and `gateway_node_index` and adds it to the topology.
	// `system-tests/lib/cre/don/don.go`.ValidateTopology() makes sure that if any DON needs gateway connector, there is at least one nodeSet with a gateway node.
	if donTopology.GatewayConnectorOutput == nil || len(donTopology.GatewayConnectorOutput.Configurations) == 0 {
		return nil, nil
	}

	// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
	// This map will be used to configure the gateway job on the node that runs it.
	for _, don := range donTopology.Dons.List() {
		// if it's a workflow DON or has capabilities that require to connect to external resources, it needs access to gateway connector
		if !don.HasFlag(cre.WorkflowDON) && !don.RequiresGateway() {
			continue
		}

		workerNode, wErr := don.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		ethAddresses := make([]string, len(workerNode))
		chainID, err := chainselectors.ChainIdFromSelector(donTopology.HomeChainSelector)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get home chain id from selector")
		}
		for i, workerNode := range workerNode {
			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			ethAddresses[i] = evmKey.PublicAddress.Hex()
		}

		handlers := map[string]string{}
		if don.HasFlag(cre.WorkflowDON) || don.RequiresWebAPI() {
			handlerConfig := `
				[gatewayConfig.Dons.Handlers.Config]
				maxAllowedMessageAgeSec = 1_000
				[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
				globalBurst = 10
				globalRPS = 50
				perSenderBurst = 10
				perSenderRPS = 10
				`
			handlers[coregateway.WebAPICapabilitiesType] = handlerConfig
		}

		// for _, capability := range input.Capabilities {
		// 	if capability.GatewayJobHandlerConfigFn() == nil {
		// 		continue
		// 	}

		// 	handlerConfig, handlerConfigErr := capability.GatewayJobHandlerConfigFn()(don)
		// 	if handlerConfigErr != nil {
		// 		return nil, errors.Wrap(handlerConfigErr, "failed to get handler config")
		// 	}
		// 	maps.Copy(handlers, handlerConfig)
		// }

		for idx := range donTopology.GatewayConnectorOutput.Configurations {
			donTopology.GatewayConnectorOutput.Configurations[idx].Dons = append(donTopology.GatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  don.Name,
				Handlers:            handlers,
			})
		}
	}

	result := make(map[uint64][]config.GatewayConfig)

	// we know that at least one DON must be the gateway DON, because topology.validate() checks that
	hasGateway := false
	for _, don := range donTopology.Dons.List() {
		_, hasGateway = don.Gateway()
		if !hasGateway {
			continue
		}

		for _, gc := range donTopology.GatewayConnectorOutput.Configurations {
			config, cErr := gatewayConfiguration(extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gc)
			if cErr != nil {
				return nil, errors.Wrap(cErr, "failed to get gateway configuration")
			}

			result[don.ID] = append(result[don.ID], config)
		}
	}

	if !hasGateway {
		return nil, errors.New("no gateway node found in any DON, but at least one is required")
	}

	return result, nil
}

type GatewayHandler struct {
	Name   string
	Config string
}

var (
	DefaultAllowedPorts = []int{80, 443}
)

func gatewayConfiguration(extraAllowedPorts []int, extraAllowedIps, extrAallowedIPsCIDR []string, gatewayConfiguration *cre.DonGatewayConfiguration) (config.GatewayConfig, error) {
	var gatewayDons string

	for _, don := range gatewayConfiguration.Dons {
		var gatewayMembers string

		for i := 0; i < len(don.MembersEthAddresses); i++ {
			gatewayMembers += fmt.Sprintf(`
	[[gatewayConfig.Dons.Members]]
	Address = "%s"
	Name = "Node %d"`,
				don.MembersEthAddresses[i],
				i+1,
			)
		}

		var handlersConfig string
		for name, config := range don.Handlers {
			handlersConfig += fmt.Sprintf(`
	[[gatewayConfig.Dons.Handlers]]
	Name = "%s"
	%s
		`, name, config)
		}

		gatewayDons += fmt.Sprintf(`
	[[gatewayConfig.Dons]]
	DonId = "%s"
	F = 1
	%s
	%s
		`, don.ID, gatewayMembers, handlersConfig)
	}

	gatewayConfig := fmt.Sprintf(`
	forwardingAllowed = false
	[gatewayConfig.ConnectionManagerConfig]
	AuthChallengeLen = 10
	AuthGatewayId = "%s"
	AuthTimestampToleranceSec = 5
	HeartbeatIntervalSec = 20
	%s
	[gatewayConfig.NodeServerConfig]
	HandshakeTimeoutMillis = 1_000
	MaxRequestBytes = 100_000
	# this is the path other nodes will use to connect to the gateway
	Path = "%s"
	# this is the port other nodes will use to connect to the gateway
	Port = %d
	ReadTimeoutMillis = 1_000
	RequestTimeoutMillis = 10_000
	WriteTimeoutMillis = 1_000
	[gatewayConfig.UserServerConfig]
	ContentTypeHeader = "application/jsonrpc"
	MaxRequestBytes = 100_000
	Path = "%s"
	Port = %d
	ReadTimeoutMillis = 80_000
	RequestTimeoutMillis = 80_000
	WriteTimeoutMillis = 80_000
	CORSEnabled = false
	CORSAllowedOrigins = []
	[gatewayConfig.HTTPClientConfig]
	MaxResponseBytes = 100_000_000
`,
		gatewayConfiguration.AuthGatewayID,
		gatewayDons,
		gatewayConfiguration.Outgoing.Path,
		gatewayConfiguration.Outgoing.Port,
		gatewayConfiguration.Incoming.Path,
		gatewayConfiguration.Incoming.InternalPort,
	)

	if len(extraAllowedPorts) != 0 {
		var allowedPorts string
		allPorts := make([]int, 0, len(DefaultAllowedPorts)+len(extraAllowedPorts))
		allPorts = append(allPorts, append(extraAllowedPorts, DefaultAllowedPorts...)...)
		for _, port := range allPorts {
			allowedPorts += fmt.Sprintf("%d, ", port)
		}

		// when we pass custom allowed IPs, defaults are not used and we need to
		// pass HTTP and HTTPS explicitly
		gatewayConfig += fmt.Sprintf(`
	AllowedPorts = [%s]
`,
			allowedPorts,
		)
	}

	if len(extraAllowedIps) != 0 {
		allowedIPs := strings.Join(extraAllowedIps, `", "`)

		gatewayConfig += fmt.Sprintf(`
	AllowedIps = ["%s"]
`,
			allowedIPs,
		)
	}

	if len(extrAallowedIPsCIDR) != 0 {
		allowedIPsCIDR := strings.Join(extrAallowedIPsCIDR, `", "`)

		gatewayConfig += fmt.Sprintf(`
	AllowedIPsCIDR = ["%s"]
`,
			allowedIPsCIDR,
		)
	}

	var gc config.GatewayConfig
	err := toml.Unmarshal([]byte(gatewayConfig), &gc)
	if err != nil {
		return config.GatewayConfig{}, errors.Wrap(err, "failed to unmarshal gateway config")
	}

	return gc, nil
}
