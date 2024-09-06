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

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

/*
	These methods should be cleaned merged after we decouple ChainlinkClient and ChainlinkK8sClient
	Please, use them while refactoring other tests to local docker env
*/

func ChainlinkNodeAddressesLocal(nodes []*client.ChainlinkClient) ([]common.Address, error) {
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
	bootstrapNode *client.ChainlinkClient,
	workerNodes []*client.ChainlinkClient,
	mockValue int,
	mockAdapter *test_env.Killgrave,
	evmChainID *big.Int,
) error {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		if err != nil {
			return fmt.Errorf("reading P2P keys from bootstrap node have failed: %w", err)
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
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
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			if err != nil {
				return err
			}
			bta := &client.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockAdapter.InternalEndpoint, strings.TrimPrefix(nodeContractPairID, "/")),
			}
			err = SetAdapterResponseLocal(mockValue, ocrInstance, node, mockAdapter)
			if err != nil {
				return fmt.Errorf("setting adapter response for OCR node failed: %w", err)
			}
			err = node.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("creating bridge on CL node failed: %w", err)
			}

			bootstrapPeers := []*client.ChainlinkClient{bootstrapNode}
			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				EVMChainID:         evmChainID.String(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  bootstrapPeers,
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
			}
			_, err = node.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR job on OCR node failed: %w", err)
			}
		}
	}
	return nil
}

func SetAdapterResponseLocal(
	response int,
	ocrInstance contracts.OffchainAggregator,
	chainlinkNode *client.ChainlinkClient,
	mockAdapter *test_env.Killgrave,
) error {
	nodeContractPairID, err := BuildNodeContractPairID(chainlinkNode, ocrInstance)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/%s", nodeContractPairID)
	err = mockAdapter.SetAdapterBasedIntValuePath(path, []string{http.MethodGet, http.MethodPost}, response)
	if err != nil {
		return fmt.Errorf("setting mock adapter value path failed: %w", err)
	}
	return nil
}

func SetAllAdapterResponsesToTheSameValueLocal(
	response int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.ChainlinkClient,
	mockAdapter *test_env.Killgrave,
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
	bootstrapNode *client.ChainlinkClient,
	workerNodes []*client.ChainlinkClient,
	mockValue int,
	mockAdapter *test_env.Killgrave,
	evmChainID string,
) error {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		if err != nil {
			return err
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
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
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			if err != nil {
				return err
			}
			bta := &client.BridgeTypeAttributes{
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

			bootstrapPeers := []*client.ChainlinkClient{bootstrapNode}
			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				EVMChainID:         evmChainID,
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  bootstrapPeers,
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
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
