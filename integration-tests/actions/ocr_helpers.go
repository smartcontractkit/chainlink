package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// This actions file often returns functions, rather than just values. These are used as common test helpers, and are
// handy to have returning as functions so that Ginkgo can use them in an aesthetically pleasing way.

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRContracts(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []*client.Chainlink,
	client blockchain.EVMClient,
) []contracts.OffchainAggregator {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying OCR instance %d shouldn't fail", contractCount+1)
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for OCR Contract deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contract deployments")

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, node := range chainlinkNodes[1:] {
		addr, err := node.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Error getting node's primary ETH address")
		transmitters = append(transmitters, addr)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		Expect(err).ShouldNot(HaveOccurred(), "Error setting OCR payees")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting OCR payees")
		}
	}
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to set payees and transmitters")

	// Set Config
	transmitterAddresses, err := ChainlinkNodeAddresses(chainlinkNodes[1:])
	Expect(err).ShouldNot(HaveOccurred(), "Getting node common addresses should not fail")
	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			chainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
			transmitterAddresses,
		)
		Expect(err).ShouldNot(HaveOccurred())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting OCR config")
		}
	}
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to set config")
	return ocrInstances
}

// DeployOCRContractsForwarderFlow deploys and funds a certain number of offchain
// aggregator contracts with forwarders as effectiveTransmitters
func DeployOCRContractsForwarderFlow(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []*client.Chainlink,
	forwarderAddresses []common.Address,
	client blockchain.EVMClient,
) []contracts.OffchainAggregator {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying OCR instance %d shouldn't fail", contractCount+1)
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for OCR Contract deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contract deployments")

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, forwarderCommonAddress := range forwarderAddresses {
		forwarderAddress := forwarderCommonAddress.Hex()
		transmitters = append(transmitters, forwarderAddress)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		Expect(err).ShouldNot(HaveOccurred(), "Error setting OCR payees")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting OCR payees")
		}
	}
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to set payees and transmitters")

	// Set Config
	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			chainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
			forwarderAddresses,
		)
		Expect(err).ShouldNot(HaveOccurred())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for setting OCR config")
		}
	}
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to set config")
	return ocrInstances
}

// CreateOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRJobs(
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) func() {
	return func() {
		for _, ocrInstance := range ocrInstances {
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from bootstrap node")
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bootstrap job on bootstrap node")

			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				nodeP2PIds, err := chainlinkNodes[nodeIndex].MustReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from OCR node %d", nodeIndex+1)
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
				nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				nodeContractPairID := BuildNodeContractPairID(chainlinkNodes[nodeIndex], ocrInstance)
				Expect(err).ShouldNot(HaveOccurred())
				bta := client.BridgeTypeAttributes{
					Name: nodeContractPairID,
					URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, nodeContractPairID),
				}

				// This sets a default value for all node and ocr instances in order to avoid 404 issues
				SetAllAdapterResponsesToTheSameValue(0, ocrInstances, chainlinkNodes, mockserver)

				err = chainlinkNodes[nodeIndex].MustCreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bridge in OCR node %d", nodeIndex+1)

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []*client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = chainlinkNodes[nodeIndex].MustCreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
			}
		}
	}
}

// CreateOCRJobsWithForwarder bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRJobsWithForwarder(
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) func() {
	return func() {
		for _, ocrInstance := range ocrInstances {
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from bootstrap node")
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bootstrap job on bootstrap node")

			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				nodeP2PIds, err := chainlinkNodes[nodeIndex].MustReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from OCR node %d", nodeIndex+1)
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
				nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				nodeContractPairID := BuildNodeContractPairID(chainlinkNodes[nodeIndex], ocrInstance)
				Expect(err).ShouldNot(HaveOccurred())
				bta := client.BridgeTypeAttributes{
					Name: nodeContractPairID,
					URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, nodeContractPairID),
				}

				// This sets default value for all node and ocr instances in order to avoid 404 issues
				SetAllAdapterResponsesToTheSameValue(0, ocrInstances, chainlinkNodes, mockserver)

				err = chainlinkNodes[nodeIndex].MustCreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bridge in OCR node %d", nodeIndex+1)

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []*client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
					ForwardingAllowed:  true,
				}
				_, err = chainlinkNodes[nodeIndex].MustCreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
			}
		}
	}
}

// SetAdapterResponse sets a single adapter response that correlates with an ocr contract and a chainlink node
func SetAdapterResponse(
	response int,
	ocrInstance contracts.OffchainAggregator,
	chainlinkNode *client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) func() {
	return func() {
		nodeContractPairID := BuildNodeContractPairID(chainlinkNode, ocrInstance)
		path := fmt.Sprintf("/%s", nodeContractPairID)
		err := mockserver.SetValuePath(path, response)
		Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")
	}
}

// SetAllAdapterResponsesToTheSameValue sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to the same response
func SetAllAdapterResponsesToTheSameValue(
	response int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) func() {
	return func() {
		var adapterVals sync.WaitGroup
		for _, o := range ocrInstances {
			ocrInstance := o
			for _, n := range chainlinkNodes {
				node := n
				adapterVals.Add(1)
				go func() {
					defer adapterVals.Done()
					SetAdapterResponse(response, ocrInstance, node, mockserver)()
				}()

			}
		}
		adapterVals.Wait()
	}
}

// SetAllAdapterResponsesToDifferentValues sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to different responses
func SetAllAdapterResponsesToDifferentValues(
	responses []int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) func() {
	return func() {
		Expect(len(responses)).Should(BeNumerically("==", len(chainlinkNodes[1:])))
		for _, ocrInstance := range ocrInstances {
			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				SetAdapterResponse(responses[nodeIndex-1], ocrInstance, chainlinkNodes[nodeIndex], mockserver)()
			}
		}
	}
}

// StartNewRound requests a new round from the ocr contracts and waits for confirmation
func StartNewRound(
	roundNr int64,
	ocrInstances []contracts.OffchainAggregator,
	client blockchain.EVMClient,
) func() {
	return func() {
		roundTimeout := time.Minute * 2
		for i := 0; i < len(ocrInstances); i++ {
			err := ocrInstances[i].RequestNewRound()
			Expect(err).ShouldNot(HaveOccurred(), "Requesting new round in OCR instance %d shouldn't fail", i+1)
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstances[i], big.NewInt(roundNr), roundTimeout, nil)
			client.AddHeaderEventSubscription(ocrInstances[i].Address(), ocrRound)
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for Event subscriptions of OCR instance %d shouldn't fail", i+1)
		}
	}
}

// BuildNodeContractPairID builds a UUID based on a related pair of a Chainlink node and OCR contract
func BuildNodeContractPairID(node *client.Chainlink, ocrInstance contracts.OffchainAggregator) string {
	Expect(node).ShouldNot(BeNil())
	Expect(ocrInstance).ShouldNot(BeNil())
	nodeAddress, err := node.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Getting chainlink node's primary ETH address shouldn't fail")
	shortNodeAddr := nodeAddress[2:12]
	shortOCRAddr := ocrInstance.Address()[2:12]
	return strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr))
}
