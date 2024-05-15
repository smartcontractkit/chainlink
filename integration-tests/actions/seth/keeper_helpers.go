package actions_seth

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	ctf_concurrency "github.com/smartcontractkit/chainlink-testing-framework/concurrency"
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

	err = DeployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfContracts, linkFundsForEachUpkeep)
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

	err = DeployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfContracts, linkFundsForEachUpkeep)
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

type upkeepRegistrationResult struct {
	upkeepID UpkeepId
}

func (r upkeepRegistrationResult) GetResult() *big.Int {
	return r.upkeepID
}

type upkeepConfig struct {
	address string
	data    []byte
}

type UpkeepId = *big.Int

func RegisterUpkeepContractsWithCheckData(t *testing.T, client *seth.Client, linkToken contracts.LinkToken, linkFunds *big.Int, upkeepGasLimit uint32, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, numberOfContracts int, upkeepAddresses []string, checkData [][]byte, isLogTrigger bool, isMercury bool) []*big.Int {
	l := logging.GetTestLogger(t)

	concurrency, err := GetAndAssertCorrectConcurrency(client, 1)
	require.NoError(t, err, "Insufficient concurrency to execute action")

	executor := ctf_concurrency.NewConcurrentExecutor[UpkeepId, upkeepRegistrationResult, upkeepConfig](l)

	configs := make([]upkeepConfig, 0)
	for i := 0; i < len(upkeepAddresses); i++ {
		configs = append(configs, upkeepConfig{address: upkeepAddresses[i], data: checkData[i]})
	}

	var registerUpkeepFn = func(resultCh chan upkeepRegistrationResult, errorCh chan error, executorNum int, config upkeepConfig) {
		uuid := uuid.New().String()
		keyNum := executorNum + 1 // key 0 is the root key

		req, err := registrar.EncodeRegisterRequest(
			fmt.Sprintf("upkeep_%s", uuid),
			[]byte("test@mail.com"),
			config.address,
			upkeepGasLimit,
			client.MustGetRootKeyAddress().Hex(), // upkeep Admin
			config.data,
			linkFunds,
			0,
			client.Addresses[keyNum].Hex(),
			isLogTrigger,
			isMercury,
		)

		if err != nil {
			errorCh <- errors.Wrapf(err, "[uuid: %s] Failed to encode register request for upkeep at %s", uuid, config.address)
			return
		}

		balance, err := linkToken.BalanceOf(context.Background(), client.Addresses[keyNum].Hex())
		if err != nil {
			errorCh <- errors.Wrapf(err, "[uuid: %s]Failed to get LINK balance of %s", uuid, client.Addresses[keyNum].Hex())
			return
		}

		// not stricly necessary, but helps us to avoid an errorless revert if there is not enough LINK
		if balance.Cmp(linkFunds) < 0 {
			errorCh <- fmt.Errorf("[uuid: %s] Not enough LINK balance for %s. Has: %s. Needs: %s", uuid, client.Addresses[keyNum].Hex(), balance.String(), linkFunds.String())
			return
		}

		tx, err := linkToken.TransferAndCallFromKey(registrar.Address(), linkFunds, req, keyNum)
		if err != nil {
			errorCh <- errors.Wrapf(err, "[uuid: %s] Failed to register upkeep at %s", uuid, config.address)
			return
		}

		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			errorCh <- errors.Wrapf(err, "[uuid: %s] Failed to get receipt for upkeep at %s and tx hash %s", uuid, config.address, tx.Hash())
			return
		}

		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}

		if upkeepId == nil {
			errorCh <- errors.Wrapf(err, "[uuid: %s] Failed find upkeep ID for upkeep at %s in logs of tx with hash %s", uuid, config.address, tx.Hash())
			return
		}

		l.Debug().
			Str("TxHash", tx.Hash().String()).
			Str("Upkeep ID", upkeepId.String()).
			Msg("Found upkeepId in tx hash")

		resultCh <- upkeepRegistrationResult{upkeepID: upkeepId}
	}

	upkeepIds, err := executor.Execute(concurrency, configs, registerUpkeepFn)
	require.NoError(t, err, "Failed to register upkeeps using executor")

	require.Equal(t, numberOfContracts, len(upkeepIds), "Incorrect number of Keeper Consumer Contracts registered")
	l.Info().Msg("Successfully registered all Keeper Consumer Contracts")

	return upkeepIds
}

type keeperConsumerResult struct {
	contract contracts.KeeperConsumer
}

func (k keeperConsumerResult) GetResult() contracts.KeeperConsumer {
	return k.contract
}

// DeployKeeperConsumers concurrently deploys keeper consumer contracts. It requires at least 1 ephemeral key to be present in Seth config.
func DeployKeeperConsumers(t *testing.T, client *seth.Client, numberOfContracts int, isLogTrigger bool, isMercury bool) []contracts.KeeperConsumer {
	l := logging.GetTestLogger(t)

	concurrency, err := GetAndAssertCorrectConcurrency(client, 1)
	require.NoError(t, err, "Insufficient concurrency to execute action")

	executor := ctf_concurrency.NewConcurrentExecutor[contracts.KeeperConsumer, keeperConsumerResult, ctf_concurrency.NoTaskType](l)

	var deployContractFn = func(channel chan keeperConsumerResult, errorCh chan error, executorNum int) {
		keyNum := executorNum + 1 // key 0 is the root key
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

		channel <- keeperConsumerResult{contract: keeperConsumerInstance}
	}

	results, err := executor.ExecuteSimple(concurrency, numberOfContracts, deployContractFn)
	require.NoError(t, err, "Failed to deploy keeper consumers")

	// require.Equal(t, 0, len(deplymentErrors), "Error deploying consumer contracts")
	require.Equal(t, numberOfContracts, len(results), "Incorrect number of Keeper Consumer Contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return results
}

// DeployKeeperConsumersPerformance sequentially deploys keeper performance consumer contracts.
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

	require.Equal(t, numberOfContracts, len(upkeeps), "Incorrect number of consumers contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeeps
}

// DeployPerformDataChecker sequentially deploys keeper perform data checker contracts.
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
	require.Equal(t, numberOfContracts, len(upkeeps), "Incorrect number of PerformDataChecker contracts deployed")
	l.Info().Msg("Successfully deployed all PerformDataChecker Contracts")

	return upkeeps
}

// DeployUpkeepCounters sequentially deploys a set amount of upkeep counter contracts.
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
	require.Equal(t, numberOfContracts, len(upkeepCounters), "Incorrect number of Keeper Consumer contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

// DeployUpkeepPerformCounter sequentially deploys a set amount of upkeep perform counter restrictive contracts.
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
	require.Equal(t, numberOfContracts, len(upkeepCounters), "Incorrect number of Keeper Consumer contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return upkeepCounters
}

// RegisterNewUpkeeps concurrently registers the given amount of new upkeeps, using the registry and registrar,
// which are passed as parameters. It returns the newly deployed contracts (consumers), as well as their upkeep IDs.
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

	concurrency, err := GetAndAssertCorrectConcurrency(chainClient, 1)
	require.NoError(t, err, "Insufficient concurrency to execute action")

	operationsPerAddress := numberOfNewUpkeeps / concurrency

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	require.NoError(t, err, "Error deploying multicall contract")

	linkFundsForEachUpkeep := big.NewInt(9e18)

	err = SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfNewUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	newUpkeepIDs := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps, false, false)

	return newlyDeployedUpkeeps, newUpkeepIDs
}

var INSUFFICIENT_EPHEMERAL_KEYS = `
Error: Insufficient Ephemeral Addresses for Simulated Network

To operate on a simulated network, you must configure at least one ephemeral address. Currently, %d ephemeral address(es) are set. Please update your TOML configuration file as follows to meet this requirement:
[Seth] ephemeral_addresses_number = 1

This adjustment ensures that your setup is minimaly viable. Although it is highly recommended to use at least 20 ephemeral addresses.
`

var INSUFFICIENT_STATIC_KEYS = `
Error: Insufficient Private Keys for Live Network

To run this test on a live network, you must either:
1. Set at least two private keys in the '[Network.WalletKeys]' section of your TOML configuration file. Example format:
   [Network.WalletKeys]
   NETWORK_NAME=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]
2. Set at least two private keys in the '[Network.EVMNetworks.NETWORK_NAME] section of your TOML configuration file. Example format:
   evm_keys=["PRIVATE_KEY_1", "PRIVATE_KEY_2"]

Currently, only %d private key/s is/are set.

Recommended Action:
Distribute your funds across multiple private keys and update your configuration accordingly. Even though 1 private key is sufficient for testing, it is highly recommended to use at least 10 private keys.
`

// GetAndAssertCorrectConcurrency checks Seth configuration for the number of ephemeral keys or static keys (depending on Seth configuration) and makes sure that
// the number is at least minConcurrency. If the number is less than minConcurrency, it returns an error. The root key is always excluded from the count.
func GetAndAssertCorrectConcurrency(client *seth.Client, minConcurrency int) (int, error) {
	concurrency := client.Cfg.GetMaxConcurrency()

	var msg string
	if client.Cfg.IsSimulatedNetwork() {
		msg = fmt.Sprintf(INSUFFICIENT_EPHEMERAL_KEYS, concurrency)
	} else {
		msg = fmt.Sprintf(INSUFFICIENT_STATIC_KEYS, concurrency)
	}

	if concurrency < minConcurrency {
		return 0, fmt.Errorf(msg)
	}

	return concurrency, nil
}
