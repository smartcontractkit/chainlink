package gateway

import (
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var GatewayJobSpecFactoryFn = func(extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			extraAllowedPorts,
			extraAllowedIPs,
			extraAllowedIPsCIDR,
			input.DonTopology.GatewayConnectorOutput,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string, gatewayConnectorOutput *cre.GatewayConnectorOutput) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(cre.DonsToJobSpecs)

	// if we don't have a gateway connector outputs, we don't need to create any job specs
	if gatewayConnectorOutput == nil || len(gatewayConnectorOutput.Configurations) == 0 {
		return donToJobSpecs, nil
	}

	// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
	// This map will be used to configure the gateway job on the node that runs it. Currently, we support only a single gateway connector, even if CRE supports multiple
	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// if it's a workflow DON or it has custom compute capability, it needs access to gateway connector
		if !flags.HasFlag(donWithMetadata.Flags, cre.WorkflowDON) && !don.NodeNeedsGateway(donWithMetadata.Flags) {
			continue
		}

		workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		ethAddresses := make([]string, len(workflowNodeSet))
		var ethAddressErr error
		for i, n := range workflowNodeSet {
			ethAddresses[i], ethAddressErr = node.FindLabelValue(n, node.AddressKeyFromSelector(donTopology.HomeChainSelector))
			if ethAddressErr != nil {
				return nil, errors.Wrap(ethAddressErr, "failed to get eth address from labels")
			}
		}

		for idx := range gatewayConnectorOutput.Configurations {
			gatewayConnectorOutput.Configurations[idx].Dons = append(gatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donWithMetadata.ID,
			})
		}
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// create job specs for the gateway node or vault DON
		if !flags.HasFlag(donWithMetadata.Flags, cre.GatewayDON) && !flags.HasFlag(donWithMetadata.Flags, cre.VaultCapability) {
			continue
		}

		gatewayNode, nodeErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: node.ExtraRolesKey, Value: cre.GatewayNode}, node.LabelContains)
		if nodeErr != nil {
			return nil, errors.Wrap(nodeErr, "failed to find bootstrap node")
		}

		gatewayNodeID, gatewayErr := node.FindLabelValue(gatewayNode, node.NodeIDKey)
		if gatewayErr != nil {
			return nil, errors.Wrap(gatewayErr, "failed to get gateway node id from labels")
		}

		homeChainID, homeChainErr := chainselectors.ChainIdFromSelector(donTopology.HomeChainSelector)
		if homeChainErr != nil {
			return nil, errors.Wrap(homeChainErr, "failed to get home chain id from selector")
		}

		if flags.HasFlag(donWithMetadata.Flags, cre.GatewayDON) {
			gatewayConfigurations := don.GatewayConfigurationsForHandler(coregateway.WebAPICapabilitiesType, gatewayConnectorOutput)
			if len(gatewayConfigurations) == 0 {
				return nil, errors.New("no gateway connector configurations found for handler type " + string(coregateway.WebAPICapabilitiesType))
			}

			handlerConfig := `
			[gatewayConfig.Dons.HandlerConfig]
			MaxAllowedMessageAgeSec = 1_000
			[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
			GlobalBurst = 10
			GlobalRPS = 50
			PerSenderBurst = 10
			PerSenderRPS = 10
			`

			for _, gatewayConfiguration := range gatewayConfigurations {
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.AnyGateway(gatewayNodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConnectorOutput.DonID, handlerConfig, gatewayConfiguration))
			}
		}

		if flags.HasFlag(donWithMetadata.Flags, cre.VaultCapability) {
			gatewayConfigurations := don.GatewayConfigurationsForHandler(coregateway.VaultHandlerType, gatewayConnectorOutput)
			if len(gatewayConfigurations) == 0 {
				return nil, errors.New("no gateway connector configurations found for handler type " + string(coregateway.VaultHandlerType))
			}

			// for some reason vault expects different field names than web API
			handlerConfig := `
			[gatewayConfig.Dons.HandlerConfig]
			request_timeout_sec = 30
			[gatewayConfig.Dons.HandlerConfig.node_rate_limiter]
			globalRPS = 100
			globalBurst = 100
			perSenderRPS = 10
			perSenderBurst = 10
			`

			for _, gatewayConfiguration := range gatewayConfigurations {
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.AnyGateway(gatewayNodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, gatewayConnectorOutput.DonID, handlerConfig, gatewayConfiguration))
			}
		}
	}

	return donToJobSpecs, nil
}
