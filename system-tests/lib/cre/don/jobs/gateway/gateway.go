package gateway

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var GatewayJobSpecFactoryFn = func(chainID uint64, extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			chainID,
			extraAllowedPorts,
			extraAllowedIPs,
			extraAllowedIPsCIDR,
			input.DonTopology.GatewayConnectorOutput,
		)
	}
}

func GenerateJobSpecs(donTopology *types.DonTopology, chainID uint64, extraAllowedPorts []int, extraAllowedIPs, extraAllowedIPsCIDR []string, gatewayConnectorOutput *types.GatewayConnectorOutput) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}

	donToJobSpecs := make(types.DonsToJobSpecs)

	// if we don't have a gateway connector output, we don't need to create any job specs
	if gatewayConnectorOutput == nil {
		return donToJobSpecs, nil
	}

	// we need to iterate over all DONs to see which need gateway connector and create a map of Don IDs and ETH addresses (which identify nodes that can use the connector)
	// This map will be used to configure the gateway job on the node that runs it. Currently, we support only a single gateway connector, even if CRE supports multiple
	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// if it's a workflow DON or it has custom compute capability, it needs access to gateway connector
		if flags.HasFlag(donWithMetadata.Flags, types.WorkflowDON) || don.NodeNeedsGateway(donWithMetadata.Flags) {
			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			ethAddresses := make([]string, len(workflowNodeSet))
			var ethAddressErr error
			for i, n := range workflowNodeSet {
				ethAddresses[i], ethAddressErr = node.FindLabelValue(n, node.EthAddressKey)
				if ethAddressErr != nil {
					return nil, errors.Wrap(ethAddressErr, "failed to get eth address from labels")
				}
			}
			gatewayConnectorOutput.Dons = append(gatewayConnectorOutput.Dons, types.GatewayConnectorDons{
				MembersEthAddresses: ethAddresses,
				ID:                  donWithMetadata.DonMetadata.ID,
			})
		}
	}

	if len(gatewayConnectorOutput.Dons) == 0 {
		return donToJobSpecs, nil
	}

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		// create job specs for the gateway node
		if flags.HasFlag(donWithMetadata.Flags, types.GatewayDON) {
			gatewayNode, nodeErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: node.ExtraRolesKey, Value: types.GatewayNode}, node.LabelContains)
			if nodeErr != nil {
				return nil, errors.Wrap(nodeErr, "failed to find bootstrap node")
			}

			gatewayNodeID, gatewayErr := node.FindLabelValue(gatewayNode, node.NodeIDKey)
			if gatewayErr != nil {
				return nil, errors.Wrap(gatewayErr, "failed to get gateway node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.AnyGateway(gatewayNodeID, chainID, donWithMetadata.ID, extraAllowedPorts, extraAllowedIPs, extraAllowedIPsCIDR, *gatewayConnectorOutput))
		}
	}

	return donToJobSpecs, nil
}
