package gateway

import (
	"github.com/pkg/errors"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	chainselectors "github.com/smartcontractkit/chain-selectors"

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
	// This map will be used to configure the gateway job on the node that runs it.
	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// if it's a workflow DON or it has custom compute capability or it has vault capability, it needs access to gateway connector
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
			donID := donWithMetadata.Name
			if flags.HasFlag(donWithMetadata.Flags, cre.VaultCapability) {
				donID = cre.VaultGatewayDonID
			}

			gatewayConnectorOutput.Configurations[idx].Dons = append(gatewayConnectorOutput.Configurations[idx].Dons, cre.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donID,
			})
		}
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
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

		homeChainID, homeChainErr := chainselectors.ChainIdFromSelector(donTopology.HomeChainSelector)
		if homeChainErr != nil {
			return nil, errors.Wrap(homeChainErr, "failed to get home chain id from selector")
		}

		handlers := make(map[string]string)

		nodeRateLimiterConfig := `
		[gatewayConfig.Dons.Handlers.Config.NodeRateLimiter]
		globalBurst = 10
		globalRPS = 50
		perSenderBurst = 10
		perSenderRPS = 10
		`

		if flags.HasFlag(donWithMetadata.Flags, cre.GatewayDON) {
			handlerConfig := `
			[gatewayConfig.Dons.Handlers.Config]
			maxAllowedMessageAgeSec = 1_000
			` + nodeRateLimiterConfig

			handlers[coregateway.WebAPICapabilitiesType] = handlerConfig
		}

		var donMetadata []*cre.DonMetadata
		for _, don := range donTopology.DonsWithMetadata {
			donMetadata = append(donMetadata, don.DonMetadata)
		}

		// if any of the DONs have http action or http trigger capability, we need to add a http handler to the jobspec for the gateway node
		if don.AnyDonHasCapability(donMetadata, cre.HTTPActionCapability) || don.AnyDonHasCapability(donMetadata, cre.HTTPTriggerCapability) {
			handlerConfig := `
			[gatewayConfig.Dons.Handlers.Config]
			maxTriggerRequestDurationMs = 5_000
			` + nodeRateLimiterConfig + `
			[gatewayConfig.Dons.Handlers.Config.UserRateLimiter]
			globalBurst = 10
			globalRPS = 50
			perSenderBurst = 10
			perSenderRPS = 10`

			handlers[coregateway.HTTPCapabilityType] = handlerConfig
		}

		// if any of the DONs have vault capability, we need to add a vault handler to the jobspec for the gateway node
		if don.AnyDonHasCapability(donMetadata, cre.VaultCapability) {
			handlerConfig := `
			[gatewayConfig.Dons.Handlers.Config]
			requestTimeoutSec = 30
			` + nodeRateLimiterConfig

			handlers[coregateway.VaultHandlerType] = handlerConfig
		}

		for _, gatewayConfiguration := range gatewayConnectorOutput.Configurations {
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.AnyGateway(gatewayNodeID, homeChainID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, handlers, gatewayConfiguration))
		}
	}

	return donToJobSpecs, nil
}
