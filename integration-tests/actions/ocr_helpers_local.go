package actions

import (
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

/*
	These methods should be cleaned merged after we decouple ChainlinkClient and ChainlinkK8sClient
	Please, use them while refactoring other tests to local docker env
*/

// FundChainlinkNodesLocal will fund all the provided Chainlink nodes with a set amount of native currency
func FundChainlinkNodesLocal(
	nodes []*client.ChainlinkClient,
	client blockchain.EVMClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		toAddr := common.HexToAddress(toAddress)
		gasEstimates, err := client.EstimateGas(ethereum.CallMsg{
			To: &toAddr,
		})
		if err != nil {
			return err
		}
		err = client.Fund(toAddress, amount, gasEstimates)
		if err != nil {
			return err
		}
	}
	return client.WaitForEvents()
}

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

func DeployOCRContractsLocal(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	workerNodes []*client.ChainlinkClient,
	client blockchain.EVMClient,
) ([]contracts.OffchainAggregator, error) {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("OCR instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, ocrInstance)
		err = client.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("failed to wait for OCR contract deployments: %w", err)
		}
	}
	err := client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contract deployments: %w", err)
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, node := range workerNodes {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, fmt.Errorf("error getting node's primary ETH address: %w", err)
		}
		transmitters = append(transmitters, addr)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for _, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		err = client.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("failed to wait for setting OCR payees: %w", err)
		}
	}

	// Set Config
	transmitterAddresses, err := ChainlinkNodeAddressesLocal(workerNodes)
	if err != nil {
		return nil, fmt.Errorf("getting node common addresses should not fail: %w", err)
	}

	for _, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes),
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			transmitterAddresses,
		)
		if err != nil {
			return nil, fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		err = client.WaitForEvents()
		if err != nil {
			return nil, fmt.Errorf("failed to wait for setting OCR config: %w", err)
		}
	}
	err = client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contracts to set config: %w", err)
	}
	return ocrInstances, nil
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

func TrackForwarderLocal(
	chainClient blockchain.EVMClient,
	authorizedForwarder common.Address,
	node *client.ChainlinkClient,
	logger zerolog.Logger,
) error {
	chainID := chainClient.GetChainID()
	_, _, err := node.TrackForwarder(chainID, authorizedForwarder)
	if err != nil {
		return fmt.Errorf("failed to track forwarder, err: %w", err)
	}
	logger.Info().Str("NodeURL", node.Config.URL).
		Str("ForwarderAddress", authorizedForwarder.Hex()).
		Str("ChaindID", chainID.String()).
		Msg("Forwarder tracked")
	return nil
}

func DeployOCRContractsForwarderFlowLocal(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	workerNodes []*client.ChainlinkClient,
	forwarderAddresses []common.Address,
	client blockchain.EVMClient,
) ([]contracts.OffchainAggregator, error) {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy offchain aggregator, err: %w", err)
		}
		ocrInstances = append(ocrInstances, ocrInstance)
		err = client.WaitForEvents()
		if err != nil {
			return nil, err
		}
	}
	if err := client.WaitForEvents(); err != nil {
		return nil, err
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, forwarderCommonAddress := range forwarderAddresses {
		forwarderAddress := forwarderCommonAddress.Hex()
		transmitters = append(transmitters, forwarderAddress)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for _, ocrInstance := range ocrInstances {
		err := ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("failed to set OCR payees, err: %w", err)
		}
		if err := client.WaitForEvents(); err != nil {
			return nil, err
		}
	}
	if err := client.WaitForEvents(); err != nil {
		return nil, err
	}

	// Set Config
	for _, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err := ocrInstance.SetConfig(
			contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes),
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			forwarderAddresses,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to set on-chain config, err: %w", err)
		}
		if err = client.WaitForEvents(); err != nil {
			return nil, err
		}
	}
	return ocrInstances, client.WaitForEvents()
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
