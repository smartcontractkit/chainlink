package actions_seth

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	test_utils "github.com/smartcontractkit/chainlink/integration-tests/utils"

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
	upkeeps, upkeepIds := DeployConsumers(t, client, registry, registrar, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep, upkeepGasLimit, false, false)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformanceKeeperContracts deploys a set amount of keeper performance contracts registered to a single registry
func DeployPerformanceKeeperContracts(
	t *testing.T,
	chainClient *seth.Client,
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
	ef, err := contracts.DeployMockETHLINKFeed(chainClient, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(chainClient, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		chainClient,
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
	registrar := DeployKeeperRegistrar(t, chainClient, registryVersion, linkToken, registrarSettings, registry)

	err = deployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfContracts, linkFundsForEachUpkeep)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	upkeeps := DeployKeeperConsumersPerformance(
		t, chainClient, numberOfContracts, blockRange, blockInterval, checkGasToBurn, performGasToBurn,
	)

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses, false, false)

	return registry, registrar, upkeeps, upkeepIds
}

// DeployPerformDataCheckerContracts deploys a set amount of keeper perform data checker contracts registered to a single registry
func DeployPerformDataCheckerContracts(
	t *testing.T,
	chainClient *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	numberOfContracts int,
	upkeepGasLimit uint32,
	linkToken contracts.LinkToken,
	registrySettings *contracts.KeeperRegistrySettings,
	linkFundsForEachUpkeep *big.Int,
	expectedData []byte,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar, []contracts.KeeperPerformDataChecker, []*big.Int) {
	ef, err := contracts.DeployMockETHLINKFeed(chainClient, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(chainClient, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		chainClient,
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

	registrar := DeployKeeperRegistrar(t, chainClient, registryVersion, linkToken, registrarSettings, registry)
	upkeeps := DeployPerformDataChecker(t, chainClient, numberOfContracts, expectedData)

	err = deployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfContracts, linkFundsForEachUpkeep)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}

	upkeepIds := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfContracts, upkeepsAddresses, false, false)

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

	concurrency := int(*client.Cfg.EphemeralAddrs)
	require.GreaterOrEqual(t, concurrency, 1, "You need at least 1 ephemeral address to deploy consumers. Please set them in TOML config: `[Seth] ephemeral_addresses_number = 10`")

	type config struct {
		address string
		data    []byte
	}

	require.Equal(t, len(upkeepAddresses), len(checkData), "Number of upkeep addresses and check data should be the same")
	configs := make([]config, 0)

	for i := 0; i < len(upkeepAddresses); i++ {
		configs = append(configs, config{address: upkeepAddresses[i], data: checkData[i]})
	}

	type result struct {
		tx  *types.Transaction
		err error
	}

	var wgProcesses sync.WaitGroup
	wgProcesses.Add(len(upkeepAddresses))

	deplymentErrors := []error{}
	deploymentCh := make(chan result, numberOfContracts)

	atomicCounter := atomic.Uint64{}

	var registerUpkeepFn = func(channel chan result, keyNum int, config config) {
		atomicCounter.Add(1)

		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%d", atomicCounter.Load()),
			[]byte("test@mail.com"),
			config.address,
			upkeepGasLimit,
			client.Addresses[0].Hex(), // upkeep Admin
			config.data,
			linkFunds,
			0,
			client.Addresses[keyNum].Hex(),
			isLogTrigger,
			isMercury,
		)

		if err != nil {
			channel <- result{err: err}
			return
		}

		balance, err := linkToken.BalanceOf(context.Background(), client.Addresses[keyNum].Hex())
		if err != nil {
			channel <- result{err: fmt.Errorf("Failed to get LINK balance of %s: %w", client.Addresses[keyNum].Hex(), err)}
			return
		}

		if balance.Cmp(linkFunds) < 0 {
			channel <- result{err: fmt.Errorf("Not enough LINK balance for %s. Has: %s. Needs: %s", client.Addresses[keyNum].Hex(), balance.String(), linkFunds.String())}
			return
		}

		tx, err := linkToken.TransferAndCallFromKey(registrar.Address(), linkFunds, req, keyNum)
		channel <- result{tx: tx, err: err}
	}

	go func() {
		defer l.Debug().Msg("Finished listening to results of registering upkeeps")
		for r := range deploymentCh {
			if r.err != nil {
				l.Error().Err(r.err).Msg("Failed to register upkeep")
				deplymentErrors = append(deplymentErrors, r.err)
				wgProcesses.Done()
				continue
			}

			registrationTxHashes = append(registrationTxHashes, r.tx.Hash())
			l.Trace().Msg("Pushed upkeep address to data array")
			wgProcesses.Done()
		}
	}()

	dividedConfigs := test_utils.DivideSlice(configs, concurrency)

	for clientNum := 1; clientNum <= concurrency; clientNum++ {
		go func(key int) {
			configs := dividedConfigs[key-1]

			l.Trace().
				Int("Key Number", key).
				Int("Upkeeps to register", len(configs)).
				Msg("Preparing to register upkeeps")

			for i := 0; i < len(configs); i++ {
				registerUpkeepFn(deploymentCh, key, configs[i])
				l.Trace().
					Int("Key Number", key).
					Str("Done/Total", fmt.Sprintf("%d/%d", (i+1), len(configs))).
					Msg("Registered upkeep")
			}

			l.Debug().
				Int("Key Number", key).
				Msg("Finished registering upkeeps")
		}(clientNum)
	}

	wgProcesses.Wait()
	close(deploymentCh)

	require.Equal(t, 0, len(deplymentErrors), "Failed to register some upkeeps")

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

	concurrency := int(*client.Cfg.EphemeralAddrs)
	require.GreaterOrEqual(t, concurrency, 1, "You need at least 1 ephemeral address to deploy consumers. Please set them in TOML config: `[Seth] ephemeral_addresses_number = 10`")

	type result struct {
		contract contracts.KeeperConsumer
		err      error
	}

	var deplymentErr error
	deploymentCh := make(chan result, numberOfContracts)
	stopCh := make(chan struct{})

	var deployContractFn = func(channel chan result, keyNum int) {
		var keeperConsumerInstance contracts.KeeperConsumer
		var err error

		if isMercury && isLogTrigger {
			// v2.1 only: Log triggered based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggeredStreamsLookupUpkeepConsumerFromKey(client, keyNum)
		} else if isMercury {
			// v2.1 only: Conditional based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationStreamsLookupUpkeepConsumerFromKey(client, keyNum, big.NewInt(1000), big.NewInt(5), false, true, false) // 1000 block test range
		} else if isLogTrigger {
			// v2.1 only: Log triggered based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggerConsumerFromKey(client, keyNum, big.NewInt(1000)) // 1000 block test range
		} else {
			// v2.0 and v2.1: Conditional based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployUpkeepCounterFromKey(client, keyNum, big.NewInt(999999), big.NewInt(5))
		}

		require.NoError(t, err, "Deploying Consumer shouldn't fail")

		channel <- result{contract: keeperConsumerInstance, err: nil}
	}

	var wgProcess sync.WaitGroup
	for i := 0; i < numberOfContracts; i++ {
		wgProcess.Add(1)
	}

	go func() {
		defer l.Debug().Msg("Finished listening to results of deploying consumer contracts")
		for contractData := range deploymentCh {
			if contractData.err != nil {
				l.Error().Err(contractData.err).Msg("Error deploying customer contract")
				deplymentErr = contractData.err
				close(stopCh)
				return
			}
			if contractData.contract != nil {
				keeperConsumerContracts = append(keeperConsumerContracts, contractData.contract)
				l.Debug().
					Str("Contract Address", contractData.contract.Address()).
					Int("Number", len(keeperConsumerContracts)).
					Int("Out Of", numberOfContracts).
					Msg("Deployed Keeper Consumer Contract")
			}
			wgProcess.Done()
		}
	}()

	operationsPerClient := numberOfContracts / concurrency
	extraOperations := numberOfContracts % concurrency

	for clientNum := 1; clientNum <= concurrency; clientNum++ {
		go func(key int) {
			numTasks := operationsPerClient
			if key <= extraOperations {
				numTasks++
			}
			for i := 0; i < numTasks; i++ {
				select {
				case <-stopCh:
					return
				default:
					deployContractFn(deploymentCh, key)
				}
			}
		}(clientNum)
	}

	wgProcess.Wait()
	close(deploymentCh)

	require.NoError(t, deplymentErr, "Error deploying consumer contracts")
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
	chainClient *seth.Client,
	linkToken contracts.LinkToken,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	upkeepGasLimit uint32,
	numberOfNewUpkeeps int,
) ([]contracts.KeeperConsumer, []*big.Int) {
	newlyDeployedUpkeeps := DeployKeeperConsumers(t, chainClient, numberOfNewUpkeeps, false, false)

	var addressesOfNewUpkeeps []string
	for _, upkeep := range newlyDeployedUpkeeps {
		addressesOfNewUpkeeps = append(addressesOfNewUpkeeps, upkeep.Address())
	}

	concurrency := int(*chainClient.Cfg.EphemeralAddrs)
	operationsPerAddress := numberOfNewUpkeeps / concurrency

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	require.NoError(t, err, "Error deploying multicall contract")

	linkFundsForEachUpkeep := big.NewInt(9e18)

	err = SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfNewUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	newUpkeepIDs := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps, false, false)

	return newlyDeployedUpkeeps, newUpkeepIDs
}
