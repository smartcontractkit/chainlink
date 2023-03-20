package actions

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

var ZeroAddress = common.Address{}

func CreateKeeperJobs(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	keeperRegistry contracts.KeeperRegistry,
	ocrConfig contracts.OCRConfig,
) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	require.NoError(t, err, "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	require.NoError(t, err, "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		require.NoError(t, err, "Error retrieving chainlink node address")
		_, err = chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress,
			MinIncomingConfirmations: 1,
		})
		require.NoError(t, err, "Creating KeeperV2 Job shouldn't fail")
	}
}

func CreateKeeperJobsWithKeyIndex(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	keeperRegistry contracts.KeeperRegistry,
	keyIndex int,
	ocrConfig contracts.OCRConfig,
) {
	// Send keeper jobs to registry and chainlink nodes
	primaryNode := chainlinkNodes[0]
	primaryNodeAddresses, err := primaryNode.EthAddresses()
	require.NoError(t, err, "Reading ETH Keys from Chainlink Client shouldn't fail")
	nodeAddresses, err := ChainlinkNodeAddressesAtIndex(chainlinkNodes, keyIndex)
	require.NoError(t, err, "Retrieving on-chain wallet addresses for chainlink nodes shouldn't fail")
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddresses[keyIndex])
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	require.NoError(t, err, "Setting keepers in the registry shouldn't fail")

	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.EthAddresses()
		require.NoError(t, err, "Error retrieving chainlink node address")
		_, err = chainlinkNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              chainlinkNodeAddress[keyIndex],
			MinIncomingConfirmations: 1,
		})
		require.NoError(t, err, "Creating KeeperV2 Job shouldn't fail")
	}
}

func DeleteKeeperJobsWithId(t *testing.T, chainlinkNodes []*client.Chainlink, id int) {
	for _, chainlinkNode := range chainlinkNodes {
		err := chainlinkNode.MustDeleteJob(strconv.Itoa(id))
		require.NoError(t, err, "Deleting KeeperV2 Job shouldn't fail")
	}
}

// DeployKeeperContracts deploys keeper registry and a number of basic upkeep contracts with an update interval of 5.
// It returns the freshly deployed registry, registrar, consumers and the IDs of the upkeeps.
func DeployKeeperContracts(
	t *testing.T,
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
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for mock feeds to deploy")

	// Deploy the transcoder here, and then set it to the registry
	transcoder := DeployUpkeepTranscoder(t, contractDeployer, client)
	registry := DeployKeeperRegistry(t, contractDeployer, client,
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
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(t, registryVersion, linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployKeeperConsumers(t, contractDeployer, client, numberOfUpkeeps)
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses,
	)
	err = client.WaitForEvents()
	require.NoError(t, err, "Error waiting for events")

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	t *testing.T,
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
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(t, contractDeployer, client,
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
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(t, registryVersion, linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployKeeperConsumersPerformance(
		t, contractDeployer, client, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn,
	)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses,
	)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformDataCheckerContracts deploys a set amount of keeper perform data checker contracts registered to a single registry
func DeployPerformDataCheckerContracts(
	t *testing.T,
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
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for mock feeds to deploy")

	registry := DeployKeeperRegistry(t, contractDeployer, client,
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
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(t, registryVersion, linkToken, registrarSettings, contractDeployer, client, registry)

	upkeeps := DeployPerformDataChecker(t, contractDeployer, client, numberOfContracts, expectedData)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(
		t, linkToken, linkFundsForEachUpkeep, client, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses,
	)

	return registry, registrar, upkeeps, upkeepIds
}

func DeployKeeperRegistry(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registryOpts *contracts.KeeperRegistryOpts,
) contracts.KeeperRegistry {
	registry, err := contractDeployer.DeployKeeperRegistry(
		registryOpts,
	)
	require.NoError(t, err, "Deploying keeper registry shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for keeper registry to deploy")

	return registry
}

func DeployKeeperRegistrar(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	linkToken contracts.LinkToken,
	registrarSettings contracts.KeeperRegistrarSettings,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	registry contracts.KeeperRegistry,
) contracts.KeeperRegistrar {
	registrar, err := contractDeployer.DeployKeeperRegistrar(registryVersion, linkToken.Address(), registrarSettings)

	require.NoError(t, err, "Deploying KeeperRegistrar contract shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for registrar to deploy")
	if registryVersion != ethereum.RegistryVersion_2_0 {
		err = registry.SetRegistrar(registrar.Address())
		require.NoError(t, err, "Registering the registrar address on the registry shouldn't fail")
		err = client.WaitForEvents()
		require.NoError(t, err, "Failed waiting for registry to set registrar")
	}

	return registrar
}

func DeployUpkeepTranscoder(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
) contracts.UpkeepTranscoder {
	transcoder, err := contractDeployer.DeployUpkeepTranscoder()
	require.NoError(t, err, "Deploying UpkeepTranscoder contract shouldn't fail")
	err = client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for transcoder to deploy")

	return transcoder
}

func RegisterUpkeepContracts(
	t *testing.T,
	linkToken contracts.LinkToken,
	linkFunds *big.Int,
	client blockchain.EVMClient,
	upkeepGasLimit uint32,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	numberOfContracts int,
	upkeepAddresses []string,
) []*big.Int {
	l := utils.GetTestLogger(t)
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
		require.NoError(t, err, "Encoding the register request shouldn't fail")
		tx, err := linkToken.TransferAndCall(registrar.Address(), linkFunds, req)
		require.NoError(t, err, "Error registering the upkeep consumer to the registrar")
		l.Debug().
			Str("Contract Address", upkeepAddress).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Str("TxHash", tx.Hash().String()).
			Msg("Registered Keeper Consumer Contract")
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait after registering upkeep consumers")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed while waiting for all consumer contracts to be registered to registrar")

	// Fetch the upkeep IDs
	for _, txHash := range registrationTxHashes {
		receipt, err := client.GetTxReceipt(txHash)
		require.NoError(t, err, "Registration tx should be completed")
		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		require.NotNil(t, upkeepId, "Upkeep ID should be found after registration")
		l.Debug().
			Str("TxHash", txHash.String()).
			Str("Upkeep ID", upkeepId.String()).
			Msg("Found upkeepId in tx hash")
		upkeepIds = append(upkeepIds, upkeepId)
	}
	l.Info().Msg("Successfully registered all Keeper Consumer Contracts")
	return upkeepIds
}

func DeployKeeperConsumers(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
) []contracts.KeeperConsumer {
	l := utils.GetTestLogger(t)
	keeperConsumerContracts := make([]contracts.KeeperConsumer, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumer(big.NewInt(5))
		require.NoError(t, err, "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		keeperConsumerContracts = append(keeperConsumerContracts, keeperConsumerInstance)
		l.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper consumer contracts")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerContracts
}

func DeployKeeperConsumersPerformance(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) []contracts.KeeperConsumerPerformance {
	l := utils.GetTestLogger(t)
	upkeeps := make([]contracts.KeeperConsumerPerformance, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contractDeployer.DeployKeeperConsumerPerformance(
			big.NewInt(blockRange),
			big.NewInt(blockInterval),
			big.NewInt(checkGasToBurn),
			big.NewInt(performGasToBurn),
		)
		require.NoError(t, err, "Deploying KeeperConsumerPerformance instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, keeperConsumerInstance)
		l.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Performance Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for KeeperConsumerPerformance deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper consumer contracts")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

func DeployPerformDataChecker(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	expectedData []byte,
) []contracts.KeeperPerformDataChecker {
	l := utils.GetTestLogger(t)
	upkeeps := make([]contracts.KeeperPerformDataChecker, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		performDataCheckerInstance, err := contractDeployer.DeployKeeperPerformDataChecker(expectedData)
		require.NoError(t, err, "Deploying KeeperPerformDataChecker instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, performDataCheckerInstance)
		l.Debug().
			Str("Contract Address", performDataCheckerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed PerformDataChecker Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 {
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for PerformDataChecker deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper perform data checker contracts")
	l.Info().Msg("Successfully deployed all PerformDataChecker Contracts")

	return upkeeps
}

func DeployUpkeepCounters(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	testRange *big.Int,
	interval *big.Int,
) []contracts.UpkeepCounter {
	l := utils.GetTestLogger(t)
	upkeepCounters := make([]contracts.UpkeepCounter, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepCounter(testRange, interval)
		require.NoError(t, err, "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		l.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper consumer contracts")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

func DeployUpkeepPerformCounterRestrictive(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	numberOfContracts int,
	testRange *big.Int,
	averageEligibilityCadence *big.Int,
) []contracts.UpkeepPerformCounterRestrictive {
	l := utils.GetTestLogger(t)
	upkeepCounters := make([]contracts.UpkeepPerformCounterRestrictive, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contractDeployer.DeployUpkeepPerformCounterRestrictive(testRange, averageEligibilityCadence)
		require.NoError(t, err, "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		l.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for KeeperConsumer deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Failed waiting for to deploy all keeper consumer contracts")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

// RegisterNewUpkeeps registers the given amount of new upkeeps, using the registry and registrar
// which are passed as parameters.
// It returns the newly deployed contracts (consumers), as well as their upkeep IDs.
func RegisterNewUpkeeps(
	t *testing.T,
	contractDeployer contracts.ContractDeployer,
	client blockchain.EVMClient,
	linkToken contracts.LinkToken,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	upkeepGasLimit uint32,
	numberOfNewUpkeeps int,
) ([]contracts.KeeperConsumer, []*big.Int) {
	newlyDeployedUpkeeps := DeployKeeperConsumers(t, contractDeployer, client, numberOfNewUpkeeps)

	var addressesOfNewUpkeeps []string
	for _, upkeep := range newlyDeployedUpkeeps {
		addressesOfNewUpkeeps = append(addressesOfNewUpkeeps, upkeep.Address())
	}

	newUpkeepIDs := RegisterUpkeepContracts(t, linkToken, big.NewInt(9e18), client, upkeepGasLimit,
		registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps)

	return newlyDeployedUpkeeps, newUpkeepIDs
}
