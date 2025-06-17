package actions

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

/*
	These methods should be cleaned merged after we decouple ChainlinkClient and ChainlinkK8sClient
	Please, use them while refactoring other tests to local docker env
*/

func ChainlinkNodeAddressesLocal(nodes []*nodeclient.ChainlinkClient) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

func CreateOCRJobsLocal(
	ocrInstances []contracts.OffchainAggregator,
	bootstrapNode *nodeclient.ChainlinkClient,
	workerNodes []*nodeclient.ChainlinkClient,
	mockValue int,
	mockAdapter *test_env.Parrot,
	evmChainID *big.Int,
) error {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		if err != nil {
			return fmt.Errorf("reading P2P keys from bootstrap node have failed: %w", err)
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &nodeclient.OCRBootstrapJobSpec{
			Name:            fmt.Sprintf("bootstrap-%s", uuid.New().String()),
			ContractAddress: ocrInstance.Address(),
			EVMChainID:      evmChainID.String(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
		if err != nil {
			return fmt.Errorf("creating bootstrap job have failed: %w", err)
		}

		for _, node := range workerNodes {
			nodeP2PIds, err := node.MustReadP2PKeys()
			if err != nil {
				return fmt.Errorf("reading P2P keys from OCR node have failed: %w", err)
			}
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := node.PrimaryEthAddress()
			if err != nil {
				return fmt.Errorf("getting primary ETH address from OCR node have failed: %w", err)
			}
			nodeOCRKeys, err := node.MustReadOCRKeys()
			if err != nil {
				return fmt.Errorf("getting OCR keys from OCR node have failed: %w", err)
			}
			nodeOCRKeyID := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			if err != nil {
				return err
			}
			bta := &nodeclient.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockAdapter.InternalEndpoint, strings.TrimPrefix(nodeContractPairID, "/")),
			}
			err = SetAdapterResponseLocal(mockValue, ocrInstance, node, mockAdapter)
			if err != nil {
				return fmt.Errorf("setting adapter response for OCR node failed: %w", err)
			}
			err = node.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("creating bridge to %s CL node failed: %w", bta.URL, err)
			}

			bootstrapPeers := []*nodeclient.ChainlinkClient{bootstrapNode}
			ocrSpec := &nodeclient.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				EVMChainID:         evmChainID.String(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  bootstrapPeers,
				KeyBundleID:        nodeOCRKeyID,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  nodeclient.ObservationSourceSpecBridge(bta),
			}
			_, err = node.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR job on OCR node failed: %w", err)
			}
		}
	}
	return nil
}

// SetAdapterResponseLocal sets the response of the mock adapter for a given OCR instance and Chainlink node
func SetAdapterResponseLocal(
	response int,
	ocrInstance contracts.OffchainAggregator,
	chainlinkNode *nodeclient.ChainlinkClient,
	mockAdapter *test_env.Parrot,
) error {
	nodeContractPairID, err := BuildNodeContractPairID(chainlinkNode, ocrInstance)
	if err != nil {
		return err
	}
	route := &parrot.Route{
		Method:             parrot.MethodAny,
		Path:               "/" + nodeContractPairID,
		ResponseBody:       response,
		ResponseStatusCode: http.StatusOK,
	}
	err = mockAdapter.SetAdapterRoute(route)
	if err != nil {
		return fmt.Errorf("setting mock adapter value path failed: %w", err)
	}
	return nil
}

// SetAllAdapterResponsesToTheSameValueLocal sets the response of the mock adapter for all OCR instances and Chainlink nodes to the same value
func SetAllAdapterResponsesToTheSameValueLocal(
	response int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*nodeclient.ChainlinkClient,
	mockAdapter *test_env.Parrot,
) error {
	eg := &errgroup.Group{}
	for _, o := range ocrInstances {
		ocrInstance := o
		for _, n := range chainlinkNodes {
			node := n
			eg.Go(func() error {
				return SetAdapterResponseLocal(response, ocrInstance, node, mockAdapter)
			})
		}
	}
	return eg.Wait()
}

func CreateOCRJobsWithForwarderLocal(
	ocrInstances []contracts.OffchainAggregator,
	bootstrapNode *nodeclient.ChainlinkClient,
	workerNodes []*nodeclient.ChainlinkClient,
	mockValue int,
	mockAdapter *test_env.Parrot,
	evmChainID string,
) error {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		if err != nil {
			return err
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &nodeclient.OCRBootstrapJobSpec{
			Name:            fmt.Sprintf("bootstrap-%s", uuid.New().String()),
			ContractAddress: ocrInstance.Address(),
			EVMChainID:      evmChainID,
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
		if err != nil {
			return err
		}

		for _, node := range workerNodes {
			nodeP2PIds, err := node.MustReadP2PKeys()
			if err != nil {
				return err
			}
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := node.PrimaryEthAddress()
			if err != nil {
				return err
			}
			nodeOCRKeys, err := node.MustReadOCRKeys()
			if err != nil {
				return err
			}
			nodeOCRKeyID := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			if err != nil {
				return err
			}
			bta := &nodeclient.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockAdapter.InternalEndpoint, strings.TrimPrefix(nodeContractPairID, "/")),
			}
			err = SetAdapterResponseLocal(mockValue, ocrInstance, node, mockAdapter)
			if err != nil {
				return err
			}
			err = node.MustCreateBridge(bta)
			if err != nil {
				return err
			}

			bootstrapPeers := []*nodeclient.ChainlinkClient{bootstrapNode}
			ocrSpec := &nodeclient.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				EVMChainID:         evmChainID,
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  bootstrapPeers,
				KeyBundleID:        nodeOCRKeyID,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  nodeclient.ObservationSourceSpecBridge(bta),
				ForwardingAllowed:  true,
			}
			_, err = node.MustCreateJob(ocrSpec)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
