package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

var ZeroAddress = common.Address{}

// CreateOCRKeeperJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRKeeperJobs(chainlinkNodes []*client.Chainlink, mockserver *ctfClient.MockserverClient, registryAddr string) {

	bootstrapNode := chainlinkNodes[0]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
	bootstrapSpec := &client.OCRBootstrapJobSpec{
		Name:            fmt.Sprintf("automation-bootstrap-%s", uuid.NewV4().String()),
		ContractAddress: registryAddr, //registry addr
		P2PPeerID:       bootstrapP2PId,
		IsBootstrapPeer: true,
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bootstrap job on bootstrap node")

	for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
		//nodeP2PIds, err := chainlinkNodes[nodeIndex].MustReadP2PKeys()
		//Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from OCR node %d", nodeIndex+1)
		//nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
		nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
		nodeOCRKeys, err := chainlinkNodes[nodeIndex].MustReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
		nodeOCRKeyId := nodeOCRKeys.Data[0].ID

		shortNodeAddr := nodeTransmitterAddress[2:12]
		shortOCRAddr := registryAddr[2:12]
		nodeContractPairID := strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr))
		//nodeContractPairID := BuildNodeContractPairID(chainlinkNodes[nodeIndex], ocrInstance)
		Expect(err).ShouldNot(HaveOccurred())
		bta := client.BridgeTypeAttributes{
			Name: nodeContractPairID,
			URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, nodeContractPairID),
		}

		err = chainlinkNodes[nodeIndex].MustCreateBridge(&bta)
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bridge in OCR node %d", nodeIndex+1)

		//ocrSpec := &client.OCRTaskJobSpec{
		//	ContractAddress:    ocrInstance.Address(),
		//	P2PPeerID:          nodeP2PId,
		//	P2PBootstrapPeers:  []*client.Chainlink{bootstrapNode},
		//	KeyBundleID:        nodeOCRKeyId,
		//	TransmitterAddress: nodeTransmitterAddress,
		//	ObservationSource:  client.ObservationSourceSpecBridge(bta),
		//}
		_, err = chainlinkNodes[nodeIndex].MustCreateJob(&client.KeeperOCRJobSpec{
			ContractID:         registryAddr,           // registryAddr
			OCRKeyBundleID:     nodeOCRKeyId,           // get node ocr2config.ID
			TransmitterID:      nodeTransmitterAddress, // node addr
			P2Pv2Bootstrappers: "",                     // bootstrap node key and address <p2p-key>@bootstrap:8000
			ChainID:            0,
		})
		Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
	}
}

func CreateKeeperJobs(chainlinkNodes []*client.Chainlink, keeperRegistry contracts.KeeperRegistry) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddresses(chainlinkNodes)
	Expect(err).ShouldNot(HaveOccurred(), "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving chainlink node address")
		_, err = chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress,
			MinIncomingConfirmations: 1,
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	}
}

func CreateKeeperJobsWithKeyIndex(chainlinkNodes []*client.Chainlink, keeperRegistry contracts.KeeperRegistry, keyIndex int) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddresses, err := primaryNode.EthAddresses()
	Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddressesAtIndex(chainlinkNodes, keyIndex)
	Expect(err).ShouldNot(HaveOccurred(), "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddresses[keyIndex])
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees)
	Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.EthAddresses()
		Expect(err).ShouldNot(HaveOccurred(), "Error retrieving chainlink node address")
		_, err = chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress[keyIndex],
			MinIncomingConfirmations: 1,
		})
		Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
	}
}

func DeleteKeeperJobsWithId(chainlinkNodes []*client.Chainlink, id int) {
	for _, chainlinkNode := range chainlinkNodes {
		err := chainlinkNode.MustDeleteJob(strconv.Itoa(id))
		Expect(err).ShouldNot(HaveOccurred(), "Deleting KeeperV2 Job shouldn't fail")
	}
}

// DeployKeeperContracts deploys keeper registry and a number of basic upkeep contracts with an update interval of 5.
// It returns the freshly deployed registry, registrar, consumers and the IDs of the upkeeps.
func DeployKeeperContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	numberOfUpkeeps int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	linkFundsForEachUpkeep *big.Int,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperConsumer, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	// Deploy the transcoder here, and then set it to the registry
	transcoder := DeployUpkeepTranscoder(contractDeployer, client)
	registry := DeployKeeperRegistry(contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  transcoder.Address(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployKeeperConsumers(contractDeployer, client, numberOfUpkeeps)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registrySettings *contracts.KeeperRegistrySettings,
	linkFundsForEachUpkeep *big.Int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperConsumerPerformance, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  ZeroAddress.Hex(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        *registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployKeeperConsumersPerformance(contractDeployer, client, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformDataCheckerContracts deploys a set amount of keeper perform data checker contracts registered to a single registry
func DeployPerformDataCheckerContracts(
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registrySettings *contracts.KeeperRegistrySettings,
	linkFundsForEachUpkeep *big.Int,
	expectedData []byte,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperPerformDataChecker, []*big.Int) {
	ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(contractDeployer, client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion: registryVersion,
			LinkAddr:        linkToken.Address(),
			ETHFeedAddr:     ef.Address(),
			GasFeedAddr:     gf.Address(),
			TranscoderAddr:  ZeroAddress.Hex(),
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        *registrySettings,
		},
	)

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployPerformDataChecker(contractDeployer, client, numberOfContracts, expectedData)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses)

	return registry, registrar, upkeeps, upkeepIds
}

func DeployKeeperRegistry(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registryOpts *contracts.KeeperRegistryOpts,
) contracts.KeeperRegistry {
	registry, err := contractDeployer.DeployKeeperRegistry(
		registryOpts,
	)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for keeper registry to deploy")

	return registry
}

func DeployKeeperRegistrar(
	linkToken contracts.LinkToken,
	registrarSettings contracts.KeeperRegistrarSettings,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registry contracts.KeeperRegistry,
) contracts.KeeperRegistrar {
	registrar, err := contractDeployer.DeployKeeperRegistrar(linkToken.Address(), registrarSettings)

	Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperRegistrar contract shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registrar to deploy")
	err = registry.SetRegistrar(registrar.Address())
	Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for registry to set registrar")

	return registrar
}

func DeployUpkeepTranscoder(contractDeployer contracts.ContractDeployer, client blockchain.EVMClient) contracts.UpkeepTranscoder {
	transcoder, err := contractDeployer.DeployUpkeepTranscoder()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepTranscoder contract shouldn't fail")
	err = client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for transcoder to deploy")

	return transcoder
}

func RegisterUpkeepContracts(
	linkToken contracts.LinkToken,
	linkFunds *big.Int,
	client blockchain.EVMClient,
	upkeepGasLimit uint32,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	numberOfContracts int,
	upkeepAddresses []string,
) []*big.Int {
	registrationTxHashes := make([]common.Hash, 0)
	upkeepIds := make([]*big.Int, 0)
	for contractCount, upkeepAddress := range upkeepAddresses {
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("0x1234"),
			upkeepAddress,
			upkeepGasLimit,
			client.GetDefaultWallet().Address(), // upkeep Admin
			[]byte("0x"),
			linkFunds,
			0,
			client.GetDefaultWallet().Address(),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
		tx, err := linkToken.TransferAndCall(registrar.Address(), linkFunds, req)
		Expect(err).ShouldNot(HaveOccurred(), "Error registering the upkeep consumer to the registrar")
		log.Debug().
			Str("Contract Address", upkeepAddress).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Str("TxHash", tx.Hash().String()).
			Msg("Registered Keeper Consumer Contract")
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait after registering upkeep consumers")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed while waiting for all consumer contracts to be registered to registrar")

	// Fetch the upkeep IDs
	for _, txHash := range registrationTxHashes {
		receipt, err := client.GetTxReceipt(txHash)
		Expect(err).ShouldNot(HaveOccurred(), "Registration tx should be completed")
		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		Expect(upkeepId).ShouldNot(BeNil(), "Upkeep ID should be found after registration")
		log.Debug().
			Str("TxHash", txHash.String()).
			Str("Upkeep ID", upkeepId.String()).
			Msg("Found upkeepId in tx hash")
		upkeepIds = append(upkeepIds, upkeepId)
	}
	log.Info().Msg("Successfully registered all Keeper Consumer Contracts")
	return upkeepIds
}

func DeployKeeperConsumers(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
) []contracts.KeeperConsumer {
	keeperConsumerContracts := make([]contracts.KeeperConsumer, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumer(big.NewInt(5))
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		keeperConsumerContracts = append(keeperConsumerContracts, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerContracts
}

func DeployKeeperConsumersPerformance(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) []contracts.KeeperConsumerPerformance {
	upkeeps := make([]contracts.KeeperConsumerPerformance, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, keeperConsumerInstance)
		log.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Performance Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumerPerformance deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

func DeployPerformDataChecker(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	expectedData []byte,
) []contracts.KeeperPerformDataChecker {
	upkeeps := make([]contracts.KeeperPerformDataChecker, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		performDataCheckerInstance, err := contractDeployer.DeployKeeperPerformDataChecker(expectedData)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperPerformDataChecker instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, performDataCheckerInstance)
		log.Debug().
			Str("Contract Address", performDataCheckerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed PerformDataChecker Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 {
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for PerformDataChecker deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper perform data checker contracts")
	log.Info().Msg("Successfully deployed all PerformDataChecker Contracts")

	return upkeeps
}

func DeployUpkeepCounters(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	testRange *big.Int,
	interval *big.Int,
) []contracts.UpkeepCounter {
	upkeepCounters := make([]contracts.UpkeepCounter, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepCounter(testRange, interval)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		log.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

func DeployUpkeepPerformCounterRestrictive(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	testRange *big.Int,
	averageEligibilityCadence *big.Int,
) []contracts.UpkeepPerformCounterRestrictive {
	upkeepCounters := make([]contracts.UpkeepPerformCounterRestrictive, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepPerformCounterRestrictive(testRange, averageEligibilityCadence)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		log.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Failed waiting for to deploy all keeper consumer contracts")
	log.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

// RegisterNewUpkeeps registers the given amount of new upkeeps, using the registry and registrar
// which are passed as parameters.
// It returns the newly deployed contracts (consumers), as well as their upkeep IDs.
func RegisterNewUpkeeps(
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	linkToken contracts.LinkToken,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	upkeepGasLimit uint32,
	numberOfNewUpkeeps int,
) ([]contracts.KeeperConsumer, []*big.Int) {
	newlyDeployedUpkeeps := DeployKeeperConsumers(contractDeployer, client, numberOfNewUpkeeps)

	var addressesOfNewUpkeeps []string
	for _, upkeep := range newlyDeployedUpkeeps {
		addressesOfNewUpkeeps = append(addressesOfNewUpkeeps, upkeep.Address())
	}

	newUpkeepIDs := RegisterUpkeepContracts(linkToken, big.NewInt(9e18), client, upkeepGasLimit,
		registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps)

	return newlyDeployedUpkeeps, newUpkeepIDs
}
