package writesolana

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	ks_solana "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func GetGenerateConfig() func(cre.GenerateConfigsInput) (cre.NodeIndexToConfigOverride, error) {
	return func(input cre.GenerateConfigsInput) (cre.NodeIndexToConfigOverride, error) {
		configOverrides := make(cre.NodeIndexToConfigOverride)
		if flags.HasFlag(input.Flags, cre.WriteSolanaCapability) {
			workerSolanaInputs := make([]*config.WorkerSolanaInput, 0)
			for chainSelector, bcOut := range input.BlockchainOutput {
				if bcOut.SolChain == nil {
					continue
				}

				chainID, err := bcOut.SolClient.GetGenesisHash(context.Background())
				if err != nil {
					return nil, errors.Wrap(err, "failed to get chainID from solana")
				}

				forwarder, err := input.Datastore.Addresses().Get(datastore.NewAddressRefKey(
					bcOut.SolChain.ChainSelector,
					ks_solana.ForwarderContract,
					semver.MustParse("1.0.0"),
					"test-forwarder",
				))
				if err != nil {
					return nil, errors.Wrap(err, "failed to get test-forwarder address")
				}
				forwarderState, err := input.Datastore.Addresses().Get(datastore.NewAddressRefKey(
					bcOut.SolChain.ChainSelector,
					ks_solana.ForwarderState,
					semver.MustParse("1.0.0"),
					"test-forwarder",
				))
				if err != nil {
					return nil, errors.Wrap(err, "failed to get test-forwarder state address")
				}

				workerSolanaInputs = append(workerSolanaInputs, &config.WorkerSolanaInput{
					Name:             fmt.Sprintf("node-%d", chainSelector),
					ChainID:          chainID.String(),
					NodeURL:          bcOut.BlockchainOutput.Nodes[0].InternalHTTPUrl,
					ForwarderAddress: forwarder.Address,
					ForwarderState:   forwarderState.Address,
				})
			}

			//workflowNodeSet, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.WorkerNode}, node.EqualLabels)
			//if err != nil {
			//	return nil, errors.Wrap(err, "failed to find worker nodes")
			//}

			//for i := range workflowNodeSet {
			//	var nodeIndex int
			//	for _, label := range workflowNodeSet[i].Labels {
			//		if label.Key == node.IndexKey {
			//			nodeIndex, err = strconv.Atoi(label.Value)
			//			if err != nil {
			////				return nil, errors.Wrap(err, "failed to convert node index to int")
			//			}
			//			break
			//		}
			//	}

			//configOverrides[nodeIndex] = config.WorkerSolana(workerSolanaInputs)
			//}

			bootstrapNodes, err := node.FindManyWithLabel(input.DonMetadata.NodesMetadata, &cre.Label{Key: node.NodeTypeKey, Value: cre.BootstrapNode}, node.EqualLabels)
			if err != nil {
				return nil, errors.Wrap(err, "failed to find bootstrap nodes")
			}
			if len(bootstrapNodes) > 0 {
				bootstrapNode := bootstrapNodes[0]
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
				fmt.Println("bootstrap node idx", nodeIndex)
				configOverrides[nodeIndex] = config.WorkerSolana(workerSolanaInputs)
			}
		}
		return configOverrides, nil
	}
}
