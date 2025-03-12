package por

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateConfigs(input cretypes.GeneratePoRConfigsInput) (cretypes.NodeIndexToConfigOverride, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	configOverrides := make(cretypes.NodeIndexToConfigOverride)

	// if it's only a gateway DON, we don't need to generate any extra configuration, the default one will do
	if keystoneflags.HasFlag(input.Flags, cretypes.GatewayDON) && (!keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) && !keystoneflags.HasFlag(input.Flags, cretypes.CapabilitiesDON)) {
		return configOverrides, nil
	}

	chainIDInt, err := strconv.Atoi(input.BlockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	// find bootstrap node for the Don
	var donBootstrapNodeHost string
	var donBootstrapNodePeerID string

	bootstrapNode, err := node.FindOneWithLabel(input.DonMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.BootstrapNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err = node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	for _, label := range bootstrapNode.Labels {
		if label.Key == node.HostLabelKey {
			donBootstrapNodeHost = label.Value
			break
		}
	}

	if donBootstrapNodeHost == "" {
		return nil, errors.New("failed to get bootstrap node host from labels")
	}

	var nodeIndex int
	for _, label := range bootstrapNode.Labels {
		if label.Key == node.IndexKey {
			nodeIndex, err = strconv.Atoi(label.Value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert node index to int")
			}
			break
		}
	}

	// generate configuration for the bootstrap node
	configOverrides[nodeIndex] = config.BootstrapEVM(donBootstrapNodePeerID, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)

	if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) {
		configOverrides[nodeIndex] += config.BoostrapDon2DonPeering(input.PeeringData)
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.WorkerNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.IndexKey {
				nodeIndex, err = strconv.Atoi(label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
			}
		}

		// for now we just assume that every worker node is connected to one EVM chain
		configOverrides[nodeIndex] = config.WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.PeeringData, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)
		var nodeEthAddr common.Address
		for _, label := range workflowNodeSet[i].Labels {
			if label.Key == node.EthAddressKey {
				if label.Value == "" {
					return nil, errors.New("eth address label value is empty")
				}
				nodeEthAddr = common.HexToAddress(label.Value)
				break
			}
		}

		if keystoneflags.HasFlag(input.Flags, cretypes.WriteEVMCapability) {
			configOverrides[nodeIndex] += config.WorkerWriteEMV(
				nodeEthAddr,
				input.ForwarderAddress,
			)
		}

		// if it's workflow DON configure workflow registry
		if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) {
			configOverrides[nodeIndex] += config.WorkerWorkflowRegistry(
				input.WorkflowRegistryAddress, chainIDUint64)
		}

		// workflow DON nodes always needs gateway connector, otherwise they won't be able to fetch the workflow
		// it's also required by custom compute, which can only run on workflow DON nodes
		if keystoneflags.HasFlag(input.Flags, cretypes.WorkflowDON) || keystoneflags.HasFlag(input.Flags, cretypes.CustomComputeCapability) {
			configOverrides[nodeIndex] += config.WorkerGateway(
				nodeEthAddr,
				chainIDUint64,
				input.DonID,
				*input.GatewayConnectorOutput,
			)
		}
	}

	return configOverrides, nil
}
