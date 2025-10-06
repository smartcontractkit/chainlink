package gateway

import (
	"fmt"
	"maps"

	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
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
		for _, donMetadata := range input.DonTopology.ToDonMetadata() {
			// if it's a workflow DON or it has custom compute capability or it has vault capability, it needs access to gateway connector
			if !flags.HasFlag(donMetadata.Flags, cre.WorkflowDON) && !don.NodeNeedsAnyGateway(donMetadata.Flags) {
				continue
			}

			workerNode, wErr := donMetadata.WorkerNodes()
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
			if flags.HasFlag(donMetadata.Flags, cre.WorkflowDON) || don.NodeNeedsWebAPIGateway(donMetadata.Flags) {
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

				handlerConfig, handlerConfigErr := capability.GatewayJobHandlerConfigFn()(donMetadata)
				if handlerConfigErr != nil {
					return nil, errors.Wrap(handlerConfigErr, "failed to get handler config")
				}
				maps.Copy(handlers, handlerConfig)
			}

			for idx := range input.DonTopology.GatewayConnectorOutput.Configurations {
				// determine here what handlers we want to build.
				input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons = append(input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
					MembersEthAddresses: ethAddresses,
					ID:                  donMetadata.Name,
					Handlers:            handlers,
				})
			}
		}

		for _, donMetadata := range input.DonTopology.ToDonMetadata() {
			// create job specs only for the gateway node
			if !flags.HasFlag(donMetadata.Flags, cre.GatewayDON) {
				continue
			}

			gatewayNode, nodeErr := donMetadata.GatewayNode()
			if nodeErr != nil {
				return nil, errors.Wrap(nodeErr, "failed to find gateway node")
			}

			gatewayNodeID, gatewayErr := node.FindLabelValue(gatewayNode, node.NodeIDKey)
			if gatewayErr != nil {
				return nil, errors.Wrap(gatewayErr, "failed to get gateway node id from labels")
			}

			homeChainID, homeChainErr := chainselectors.ChainIdFromSelector(input.DonTopology.HomeChainSelector)
			if homeChainErr != nil {
				return nil, errors.Wrap(homeChainErr, "failed to get home chain id from selector")
			}

			for _, gatewayConfiguration := range input.DonTopology.GatewayConnectorOutput.Configurations {
				donToJobSpecs[donMetadata.ID] = append(donToJobSpecs[donMetadata.ID], jobs.AnyGateway(gatewayNodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConfiguration))
			}
		}

		return donToJobSpecs, nil
	}
}
