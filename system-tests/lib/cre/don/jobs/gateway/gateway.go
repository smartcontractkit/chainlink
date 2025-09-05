package gateway

import (
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
		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			// if it's a workflow DON or it has custom compute capability or it has vault capability, it needs access to gateway connector
			if !flags.HasFlag(donWithMetadata.Flags, cre.WorkflowDON) && !don.NodeNeedsAnyGateway(donWithMetadata.Flags) {
				continue
			}

			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			ethAddresses := make([]string, len(workflowNodeSet))
			var ethAddressErr error
			for i, n := range workflowNodeSet {
				ethAddresses[i], ethAddressErr = node.FindLabelValue(n, node.AddressKeyFromSelector(input.DonTopology.HomeChainSelector))
				if ethAddressErr != nil {
					return nil, errors.Wrap(ethAddressErr, "failed to get eth address from labels")
				}
			}

			handlers := map[string]string{}
			if flags.HasFlag(donWithMetadata.Flags, cre.WorkflowDON) || don.NodeNeedsWebAPIGateway(donWithMetadata.Flags) {
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

				handlerConfig, handlerConfigErr := capability.GatewayJobHandlerConfigFn()(donWithMetadata.DonMetadata)
				if handlerConfigErr != nil {
					return nil, errors.Wrap(handlerConfigErr, "failed to get handler config")
				}
				maps.Copy(handlers, handlerConfig)
			}

			for idx := range input.DonTopology.GatewayConnectorOutput.Configurations {
				// determine here what handlers we want to build.
				input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons = append(input.DonTopology.GatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
					MembersEthAddresses: ethAddresses,
					ID:                  donWithMetadata.Name,
					Handlers:            handlers,
				})
			}
		}

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			// create job specs only for the gateway node
			if !flags.HasFlag(donWithMetadata.Flags, cre.GatewayDON) {
				continue
			}

			gatewayNode, nodeErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.ExtraRolesKey, Value: cre.GatewayNode}, node.LabelContains)
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
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.AnyGateway(gatewayNodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConfiguration))
			}
		}

		return donToJobSpecs, nil
	}
}
