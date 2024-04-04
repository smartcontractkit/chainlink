package actions_seth

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

var ZeroAddress = common.Address{}

// DeployKeeperContracts deploys keeper registry and a number of basic upkeep contracts with an update interval of 5.
// It returns the freshly deployed registry, registrar, consumers and the IDs of the upkeeps.
func DeployKeeperContracts(
	t *testing.T,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	numberOfUpkeeps int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	client *seth.Client,
	linkFundsForEachUpkeep *big.Int,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperConsumer, []*big.Int) {
	ef, err := contracts.DeployMockETHLINKFeed(client, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(client, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	// Deploy the transcoder here, and then set it to the registry
	transcoder, err := contracts.DeployUpkeepTranscoder(client)
	require.NoError(t, err, "Deploying UpkeepTranscoder contract shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		client,
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
	require.NoError(t, err, "Deploying KeeperRegistry shouldn't fail")

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfUpkeeps))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}

	registrar := DeployKeeperRegistrar(t, client, registryVersion, linkToken, registrarSettings, registry)
	upkeeps := DeployKeeperConsumers(t, client, numberOfUpkeeps, false, false)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(t, client, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, false, false)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	registrySettings *contracts.KeeperRegistrySettings,
	linkFundsForEachUpkeep *big.Int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperConsumerPerformance, []*big.Int) {
	ef, err := contracts.DeployMockETHLINKFeed(client, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(client, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		client,
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
	require.NoError(t, err, "Deploying KeeperRegistry shouldn't fail")

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}
	registrar := DeployKeeperRegistrar(t, client, registryVersion, linkToken, registrarSettings, registry)

	upkeeps := DeployKeeperConsumersPerformance(
		t, client, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn,
	)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(t, client, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses, false, false)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformDataCheckerContracts deploys a set amount of keeper perform data checker contracts registered to a single registry
func DeployPerformDataCheckerContracts(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	registrySettings *contracts.KeeperRegistrySettings,
	linkFundsForEachUpkeep *big.Int,
	expectedData []byte,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperPerformDataChecker, []*big.Int) {
	ef, err := contracts.DeployMockETHLINKFeed(client, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(client, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		client,
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
	require.NoError(t, err, "Deploying KeeperRegistry shouldn't fail")

	// Fund the registry with 1 LINK * amount of KeeperConsumerPerformance contracts
	err = linkToken.Transfer(registry.Address(), big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(numberOfContracts))))
	require.NoError(t, err, "Funding keeper registry contract shouldn't fail")

	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
	}

	registrar := DeployKeeperRegistrar(t, client, registryVersion, linkToken, registrarSettings, registry)
	upkeeps := DeployPerformDataChecker(t, client, numberOfContracts, expectedData)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(t, client, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses, false, false)

	return registry, registrar, upkeeps, upkeepIds
}

func DeployKeeperRegistrar(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	linkToken contracts.LinkToken,
	registrarSettings contracts.KeeperRegistrarSettings,
	registry contracts.KeeperRegistry,
) contracts.KeeperRegistrar {
	registrar, err := contracts.DeployKeeperRegistrar(client, registryVersion, linkToken.Address(), registrarSettings)
	require.NoError(t, err, "Failed waiting for registrar to deploy")
	if registryVersion != ethereum.RegistryVersion_2_0 {
		err = registry.SetRegistrar(registrar.Address())
		require.NoError(t, err, "Registering the registrar address on the registry shouldn't fail")
	}

	return registrar
}

func RegisterUpkeepContracts(t *testing.T, client *seth.Client, linkToken contracts.LinkToken, linkFunds *big.Int, upkeepGasLimit uint32, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, numberOfContracts int, upkeepAddresses []string, isLogTrigger bool, isMercury bool) []*big.Int {
	checkData := make([][]byte, 0)
	for i := 0; i < numberOfContracts; i++ {
		checkData = append(checkData, []byte("0"))
	}
	return RegisterUpkeepContractsWithCheckData(
		t, client, linkToken, linkFunds, upkeepGasLimit, registry, registrar,
		numberOfContracts, upkeepAddresses, checkData, isLogTrigger, isMercury)
}

func RegisterUpkeepContractsWithCheckData(t *testing.T, client *seth.Client, linkToken contracts.LinkToken, linkFunds *big.Int, upkeepGasLimit uint32, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, numberOfContracts int, upkeepAddresses []string, checkData [][]byte, isLogTrigger bool, isMercury bool) []*big.Int {
	l := logging.GetTestLogger(t)
	registrationTxHashes := make([]common.Hash, 0)
	upkeepIds := make([]*big.Int, 0)
	for contractCount, upkeepAddress := range upkeepAddresses {
		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", contractCount+1),
			[]byte("test@mail.com"),
			upkeepAddress,
			upkeepGasLimit,
			client.Addresses[0].Hex(), // upkeep Admin
			checkData[contractCount],
			linkFunds,
			0,
			client.Addresses[0].Hex(),
			isLogTrigger,
			isMercury,
		)
		require.NoError(t, err, "Encoding the register request shouldn't fail")
		tx, err := linkToken.TransferAndCall(registrar.Address(), linkFunds, req)
		require.NoError(t, err, "Error registering the upkeep consumer to the registrar")
		l.Debug().
			Str("Contract Address", upkeepAddress).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Str("TxHash", tx.Hash().String()).
			Str("Check Data", hexutil.Encode(checkData[contractCount])).
			Msg("Registered Keeper Consumer Contract")
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
	}

	// Fetch the upkeep IDs
	for _, txHash := range registrationTxHashes {
		receipt, err := client.Client.TransactionReceipt(context.Background(), txHash)
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

func DeployKeeperConsumers(t *testing.T, client *seth.Client, numberOfContracts int, isLogTrigger bool, isMercury bool) []contracts.KeeperConsumer {
	l := logging.GetTestLogger(t)
	keeperConsumerContracts := make([]contracts.KeeperConsumer, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		var keeperConsumerInstance contracts.KeeperConsumer
		var err error

		if isMercury && isLogTrigger {
			// v2.1 only: Log triggered based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggeredStreamsLookupUpkeepConsumer(client)
		} else if isMercury {
			// v2.1 only: Conditional based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationStreamsLookupUpkeepConsumer(client, big.NewInt(1000), big.NewInt(5), false, true, false) // 1000 block test range
		} else if isLogTrigger {
			// v2.1 only: Log triggered based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggerConsumer(client, big.NewInt(1000)) // 1000 block test range
		} else {
			// v2.0 and v2.1: Conditional based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployUpkeepCounter(client, big.NewInt(999999), big.NewInt(5))
		}

		require.NoError(t, err, "Deploying Consumer instance %d shouldn't fail", contractCount+1)
		keeperConsumerContracts = append(keeperConsumerContracts, keeperConsumerInstance)
		l.Debug().
			Str("Contract Address", keeperConsumerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
	}
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return keeperConsumerContracts
}

func DeployKeeperConsumersPerformance(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
) []contracts.KeeperConsumerPerformance {
	l := logging.GetTestLogger(t)
	upkeeps := make([]contracts.KeeperConsumerPerformance, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		keeperConsumerInstance, err := contracts.DeployKeeperConsumerPerformance(
			client,
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
	}
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

func DeployPerformDataChecker(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	expectedData []byte,
) []contracts.KeeperPerformDataChecker {
	l := logging.GetTestLogger(t)
	upkeeps := make([]contracts.KeeperPerformDataChecker, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		performDataCheckerInstance, err := contracts.DeployKeeperPerformDataChecker(client, expectedData)
		require.NoError(t, err, "Deploying KeeperPerformDataChecker instance %d shouldn't fail", contractCount+1)
		upkeeps = append(upkeeps, performDataCheckerInstance)
		l.Debug().
			Str("Contract Address", performDataCheckerInstance.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed PerformDataChecker Contract")
	}
	l.Info().Msg("Successfully deployed all PerformDataChecker Contracts")

	return upkeeps
}

func DeployUpkeepCounters(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	testRange *big.Int,
	interval *big.Int,
) []contracts.UpkeepCounter {
	l := logging.GetTestLogger(t)
	upkeepCounters := make([]contracts.UpkeepCounter, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contracts.DeployUpkeepCounter(client, testRange, interval)
		require.NoError(t, err, "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		l.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
	}
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

func DeployUpkeepPerformCounterRestrictive(
	t *testing.T,
	client *seth.Client,
	numberOfContracts int,
	testRange *big.Int,
	averageEligibilityCadence *big.Int,
) []contracts.UpkeepPerformCounterRestrictive {
	l := logging.GetTestLogger(t)
	upkeepCounters := make([]contracts.UpkeepPerformCounterRestrictive, 0)

	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		// Deploy consumer
		upkeepCounter, err := contracts.DeployUpkeepPerformCounterRestrictive(client, testRange, averageEligibilityCadence)
		require.NoError(t, err, "Deploying KeeperConsumer instance %d shouldn't fail", contractCount+1)
		upkeepCounters = append(upkeepCounters, upkeepCounter)
		l.Debug().
			Str("Contract Address", upkeepCounter.Address()).
			Int("Number", contractCount+1).
			Int("Out Of", numberOfContracts).
			Msg("Deployed Keeper Consumer Contract")
	}
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

// RegisterNewUpkeeps registers the given amount of new upkeeps, using the registry and registrar
// which are passed as parameters.
// It returns the newly deployed contracts (consumers), as well as their upkeep IDs.
func RegisterNewUpkeeps(
	t *testing.T,
	client *seth.Client,
	linkToken contracts.LinkToken,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	upkeepGasLimit uint32,
	numberOfNewUpkeeps int,
) ([]contracts.KeeperConsumer, []*big.Int) {
	newlyDeployedUpkeeps := DeployKeeperConsumers(t, client, numberOfNewUpkeeps, false, false)

	var addressesOfNewUpkeeps []string
	for _, upkeep := range newlyDeployedUpkeeps {
		addressesOfNewUpkeeps = append(addressesOfNewUpkeeps, upkeep.Address())
	}

	newUpkeepIDs := RegisterUpkeepContracts(t, client, linkToken, big.NewInt(9e18), upkeepGasLimit, registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps, false, false)

	return newlyDeployedUpkeeps, newUpkeepIDs
}
