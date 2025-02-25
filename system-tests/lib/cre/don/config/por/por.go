package por

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	keystoneflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateConfigs(input types.GeneratePoRConfigsInput) (types.NodeIndexToConfigOverrides, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}
	configOverrides := make(types.NodeIndexToConfigOverrides)

	chainIDInt, err := strconv.Atoi(input.BlockchainOutput.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert chain ID to int")
	}
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	// find bootstrap node
	bootstrapNode, err := node.FindOneWithLabel(input.Don, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.BootstrapNode)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	donBootstrapNodePeerID, err := node.ToP2PID(*bootstrapNode, node.KeyExtractingTransformFn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get bootstrap node peer ID")
	}

	var donBootstrapNodeHost string
	for _, label := range bootstrapNode.Labels() {
		if label.Key == node.HostLabelKey {
			donBootstrapNodeHost = *label.Value
			break
		}
	}

	if donBootstrapNodeHost == "" {
		return nil, errors.New("failed to get bootstrap node host from labels")
	}

	var nodeIndex int
	for _, label := range bootstrapNode.Labels() {
		if label.Key == node.NodeIndexKey {
			nodeIndex, err = strconv.Atoi(*label.Value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert node index to int")
			}
		}
	}

	// generat configuration for the bootstrap node
	configOverrides[nodeIndex] = config.BootstrapEVM(donBootstrapNodePeerID, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)

	if keystoneflags.HasFlag(input.Flags, types.WorkflowDON) {
		configOverrides[nodeIndex] += config.BoostrapDon2DonPeering(input.PeeringData)

		if input.GatewayConnectorOutput == nil {
			return nil, errors.New("GatewayConnectorOutput is required for Workflow DON")
		}
		input.GatewayConnectorOutput.Host = donBootstrapNodeHost
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.Don, &ptypes.Label{Key: node.RoleLabelKey, Value: ptr.Ptr(types.WorkerNode)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	for i := range workflowNodeSet {
		var nodeIndex int
		for _, label := range workflowNodeSet[i].Labels() {
			if label.Key == node.NodeIndexKey {
				nodeIndex, err = strconv.Atoi(*label.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to convert node index to int")
				}
			}
		}

		configOverrides[nodeIndex] = config.WorkerEVM(donBootstrapNodePeerID, donBootstrapNodeHost, input.PeeringData, chainIDUint64, input.CapabilitiesRegistryAddress, input.BlockchainOutput.Nodes[0].DockerInternalHTTPUrl, input.BlockchainOutput.Nodes[0].DockerInternalWSUrl)
		nodeEthAddr := common.HexToAddress(workflowNodeSet[i].AccountAddr[chainIDUint64])

		if keystoneflags.HasFlag(input.Flags, types.WriteEVMCapability) {
			configOverrides[nodeIndex] += config.WorkerWriteEMV(
				nodeEthAddr,
				input.ForwarderAddress,
			)
		}

		// if it's workflow DON configure workflow registry
		if keystoneflags.HasFlag(input.Flags, types.WorkflowDON) {
			configOverrides[nodeIndex] += config.WorkerWorkflowRegistry(
				input.WorkflowRegistryAddress, chainIDUint64)
		}

		// workflow DON nodes always needs gateway connector, otherwise they won't be able to fetch the workflow
		// it's also required by custom compute, which can only run on workflow DON nodes
		if keystoneflags.HasFlag(input.Flags, types.WorkflowDON) || keystoneflags.HasFlag(input.Flags, types.CustomComputeCapability) {
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
