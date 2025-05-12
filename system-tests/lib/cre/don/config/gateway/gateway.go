package gateway

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"

	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateConfig(input types.GenerateConfigsInput) (types.NodeIndexToConfigOverride, error) {
	configOverrides := make(types.NodeIndexToConfigOverride)

	if input.GatewayConnectorOutput == nil {
		return configOverrides, errors.New("gateway connector output is not set")
	}

	// find worker nodes
	workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &types.Label{Key: node.NodeTypeKey, Value: types.WorkerNode}, node.EqualLabels)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find worker nodes")
	}

	homeChainID, homeErr := chain_selectors.ChainIdFromSelector(input.HomeChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	workflowRegistryAddress, workErr := crecontracts.FindAddressesForChain(input.AddressBook, input.HomeChainSelector, keystone_changeset.WorkflowRegistry.String())
	if workErr != nil {
		return nil, errors.Wrap(workErr, "failed to find WorkflowRegistry address")
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

		// we need to configure workflow registry
		if flags.HasFlag(input.Flags, types.WorkflowDON) {
			configOverrides[nodeIndex] += config.WorkerWorkflowRegistry(
				workflowRegistryAddress, homeChainID)
		}

		// workflow DON nodes might need gateway connector to download WASM workflow binaries,
		// but if the workflowDON is using only workflow jobs, we don't need to set the gateway connector
		// gateway is also required by various capabilities
		if flags.HasFlag(input.Flags, types.WorkflowDON) || don.NodeNeedsGateway(input.Flags) {
			var nodeEthAddr common.Address
			expectedAddressKey := node.AddressKeyFromSelector(input.HomeChainSelector)
			for _, label := range workflowNodeSet[i].Labels {
				if label.Key == expectedAddressKey {
					if label.Value == "" {
						return nil, errors.Errorf("%s label value is empty", expectedAddressKey)
					}
					nodeEthAddr = common.HexToAddress(label.Value)
					break
				}
			}

			configOverrides[nodeIndex] += config.WorkerGateway(
				nodeEthAddr,
				homeChainID,
				input.DonMetadata.ID,
				*input.GatewayConnectorOutput,
			)
		}
	}

	return configOverrides, nil
}
