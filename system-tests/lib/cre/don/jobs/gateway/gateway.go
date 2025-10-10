package gateway

import (
	"fmt"
	"maps"

	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
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
