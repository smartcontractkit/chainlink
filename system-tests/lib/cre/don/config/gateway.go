package config

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func GatewayConfig(
	cldEnvironment *cldf.Environment,
	blockchainOutput *blockchain.Output,
	topology *cre.Topology,
	infraInput infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string,
) (map[cre.NodeUUID]config.GatewayConfig, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	chainID, chErr := strconv.ParseUint(blockchainOutput.ChainID, 10, 64)
	if chErr != nil {
		return nil, errors.Wrap(chErr, "failed to parse chain ID")
	}

	// donToJobSpecs := make(cre.DonsToJobSpecs)

	// if we don't have a gateway connector outputs, we don't need to create any job specs
	// GatewayConnectorOutput is added by `system-tests/lib/cre/don/don.go`.BuildTopology() function, which builds gateway configuration
	// based on DON flags (cre.GatewayDON) and `gateway_node_index` and adds it to the topology.
	// `system-tests/lib/cre/don/don.go`.ValidateTopology() makes sure that if any DON needs gateway connector, there is at least one nodeSet with a gateway node.
	if topology.GatewayConnectorOutput == nil || len(topology.GatewayConnectorOutput.Configurations) == 0 {
		return nil, nil
	}

	// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
	// This map will be used to configure the gateway job on the node that runs it.
	for _, donMetadata := range topology.DonsMetadata.List() {
		// if it's a workflow DON or has capabilities that require to connect to external resources, it needs access to gateway connector
		if !donMetadata.HasFlag(cre.WorkflowDON) && !donMetadata.RequiresGateway() {
			continue
		}

		workerNode, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		ethAddresses := make([]string, len(workerNode))
		// chainID, err := chainselectors.ChainIdFromSelector(topology.HomeChainSelector)
		// if err != nil {
		// 	return nil, errors.Wrap(err, "failed to get home chain id from selector")
		// }
		for i, workerNode := range workerNode {
			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}
			ethAddresses[i] = evmKey.PublicAddress.Hex()
		}

		// handlers := map[string]string{}
		// hanldersC := config.Handler{}
		// if donMetadata.HasFlag(cre.WorkflowDON) || donMetadata.RequiresWebAPI() {
		// 	handlerConfig := `
		// 		[Dons.Handlers.Config]
		// 		maxAllowedMessageAgeSec = 1_000
		// 		[Dons.Handlers.Config.NodeRateLimiter]
		// 		globalBurst = 10
		// 		globalRPS = 50
		// 		perSenderBurst = 10
		// 		perSenderRPS = 10
		// 		`
		// 	handlers[coregateway.WebAPICapabilitiesType] = handlerConfig
		// }

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

		handlers := []config.Handler{}

		if donMetadata.HasFlag(cre.WorkflowDON) || donMetadata.RequiresWebAPI() {
			webApiConfig := config.Handler{
				Name: coregateway.WebAPICapabilitiesType,
				Config: []byte(`
maxAllowedMessageAgeSec = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10
`),
			}
			// handlerConfig := `
			// 	[Dons.Handlers.Config]
			// 	maxAllowedMessageAgeSec = 1_000
			// 	[Dons.Handlers.Config.NodeRateLimiter]
			// 	globalBurst = 10
			// 	globalRPS = 50
			// 	perSenderBurst = 10
			// 	perSenderRPS = 10
			// 	`
			// // handlers[coregateway.WebAPICapabilitiesType] = handlerConfig
			handlers = append(handlers, webApiConfig)
		}

		for idx := range topology.GatewayConnectorOutput.Configurations {
			topology.GatewayConnectorOutput.Configurations[idx].Dons = append(topology.GatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donMetadata.Name,
				HandlersC:           handlers,
				// Handlers:            handlers,
			})
		}
	}

	result := make(map[string]config.GatewayConfig)

	// we know that at least one DON must have a gateway node, because topology.validate() checks that
	hasGateway := false
	for _, donMetadata := range topology.DonsMetadata.List() {
		var gateway *cre.NodeMetadata
		gateway, hasGateway = donMetadata.Gateway()
		if !hasGateway {
			continue
		}

		for _, gc := range topology.GatewayConnectorOutput.Configurations {
			if gc.NodeUUID != gateway.UUID {
				continue
			}

			// [UserServerConfig]
			// ContentTypeHeader = "application/jsonrpc"
			// MaxRequestBytes = 100_000
			// Path = "%s"
			// Port = %d
			// ReadTimeoutMillis = 80_000
			// RequestTimeoutMillis = 80_000
			// WriteTimeoutMillis = 80_000
			// CORSEnabled = false
			// CORSAllowedOrigins = []

			c := config.GatewayConfig{
				ConnectionManagerConfig: config.ConnectionManagerConfig{
					AuthGatewayId:             gc.AuthGatewayID,
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
						Path:                 gc.Outgoing.Path,
						Port:                 uint16(gc.Outgoing.Port),
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
					Path:                 gc.Incoming.Path,
					Port:                 uint16(gc.Incoming.InternalPort),
				},
				HTTPClientConfig: gw_net.HTTPClientConfig{
					MaxResponseBytes: 100_000_000,
					AllowedPorts:     append(extraAllowedPorts, DefaultAllowedPorts...),
					AllowedIPs:       extraAllowedIPs,
					AllowedIPsCIDR:   extraAllowedIPsCIDR,
					// DefaultTimeout:   time.Duration(60 * time.Second),
				},
			}

			for _, don := range gc.Dons {
				donConfig := config.DONConfig{
					DonId:   don.ID,
					F:       1,
					Members: make([]config.NodeConfig, len(don.MembersEthAddresses)),
				}

				for i, memberAddress := range don.MembersEthAddresses {
					donConfig.Members[i] = config.NodeConfig{
						Address: memberAddress,
						Name:    fmt.Sprintf("%s-node-%d", don.ID, i),
					}
				}

				for _, handler := range don.HandlersC {
					donConfig.Handlers = append(donConfig.Handlers, handler)
				}

				c.Dons = append(c.Dons, donConfig)
			}

			// config, cErr := gatewayConfiguration(extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gc)
			// if cErr != nil {
			// 	return nil, errors.Wrap(cErr, "failed to get gateway configuration")
			// }

			result[gateway.UUID] = c
		}
	}

	if !hasGateway {
		return nil, errors.New("no gateway node found in any DON, but at least one is required")
	}

	return result, nil
}

// type GatewayHandler struct {
// 	Name   string
// 	Config string
// }

var (
	DefaultAllowedPorts = []int{80, 443}
)

// func gatewayConfiguration(extraAllowedPorts []int, extraAllowedIps, extrAallowedIPsCIDR []string, gatewayConfiguration *cre.DonGatewayConfiguration) (config.GatewayConfig, error) {
// 	var gatewayDons string

// 	for _, don := range gatewayConfiguration.Dons {
// 		var gatewayMembers string

// 		for i := 0; i < len(don.MembersEthAddresses); i++ {
// 			gatewayMembers += fmt.Sprintf(`
// 	[[Dons.Members]]
// 	Address = "%s"
// 	Name = "Node %d"`,
// 				don.MembersEthAddresses[i],
// 				i+1,
// 			)
// 		}

// 		var handlersConfig string
// 		for name, config := range don.Handlers {
// 			handlersConfig += fmt.Sprintf(`
// 	[[Dons.Handlers]]
// 	Name = "%s"
// 	%s
// 		`, name, config)
// 		}

// 		gatewayDons += fmt.Sprintf(`
// 	[[Dons]]
// 	DonId = "%s"
// 	F = 1
// 	%s
// 	%s
// 		`, don.ID, gatewayMembers, handlersConfig)
// 	}

// 	gatewayConfig := fmt.Sprintf(`
// 	[ConnectionManagerConfig]
// 	AuthChallengeLen = 10
// 	AuthGatewayId = "%s"
// 	AuthTimestampToleranceSec = 5
// 	HeartbeatIntervalSec = 20
// 	%s
// 	[NodeServerConfig]
// 	HandshakeTimeoutMillis = 1_000
// 	MaxRequestBytes = 100_000
// 	# this is the path other nodes will use to connect to the gateway
// 	Path = "%s"
// 	# this is the port other nodes will use to connect to the gateway
// 	Port = %d
// 	ReadTimeoutMillis = 1_000
// 	RequestTimeoutMillis = 10_000
// 	WriteTimeoutMillis = 1_000
// 	[UserServerConfig]
// 	ContentTypeHeader = "application/jsonrpc"
// 	MaxRequestBytes = 100_000
// 	Path = "%s"
// 	Port = %d
// 	ReadTimeoutMillis = 80_000
// 	RequestTimeoutMillis = 80_000
// 	WriteTimeoutMillis = 80_000
// 	CORSEnabled = false
// 	CORSAllowedOrigins = []
// 	[HTTPClientConfig]
// 	MaxResponseBytes = 100_000_000
// `,
// 		gatewayConfiguration.AuthGatewayID,
// 		gatewayDons,
// 		gatewayConfiguration.Outgoing.Path,
// 		gatewayConfiguration.Outgoing.Port,
// 		gatewayConfiguration.Incoming.Path,
// 		gatewayConfiguration.Incoming.InternalPort,
// 	)

// 	if len(extraAllowedPorts) != 0 {
// 		var allowedPorts string
// 		allPorts := make([]int, 0, len(DefaultAllowedPorts)+len(extraAllowedPorts))
// 		allPorts = append(allPorts, append(extraAllowedPorts, DefaultAllowedPorts...)...)
// 		for _, port := range allPorts {
// 			allowedPorts += fmt.Sprintf("%d, ", port)
// 		}

// 		// when we pass custom allowed IPs, defaults are not used and we need to
// 		// pass HTTP and HTTPS explicitly
// 		gatewayConfig += fmt.Sprintf(`
// 	AllowedPorts = [%s]
// `,
// 			allowedPorts,
// 		)
// 	}

// 	if len(extraAllowedIps) != 0 {
// 		allowedIPs := strings.Join(extraAllowedIps, `", "`)

// 		gatewayConfig += fmt.Sprintf(`
// 	AllowedIps = ["%s"]
// `,
// 			allowedIPs,
// 		)
// 	}

// 	if len(extrAallowedIPsCIDR) != 0 {
// 		allowedIPsCIDR := strings.Join(extrAallowedIPsCIDR, `", "`)

// 		gatewayConfig += fmt.Sprintf(`
// 	AllowedIPsCIDR = ["%s"]
// `,
// 			allowedIPsCIDR,
// 		)
// 	}

// 	var gc config.GatewayConfig
// 	err := toml.Unmarshal([]byte(gatewayConfig), &gc)
// 	if err != nil {
// 		return config.GatewayConfig{}, errors.Wrap(err, "failed to unmarshal gateway config")
// 	}

// 	return gc, nil
// }
