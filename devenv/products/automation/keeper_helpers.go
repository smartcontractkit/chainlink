package automation

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	ctf_concurrency "github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
)

var ZeroAddress = common.Address{}

func RegisterUpkeepContracts(t *testing.T, client *seth.Client, linkToken contracts.LinkToken, fundsForEachUpkeep *big.Int, upkeepGasLimit uint32, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, numberOfContracts int, upkeepAddresses []string, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken) []*big.Int {
	checkData := make([][]byte, 0)
	for range numberOfContracts {
		checkData = append(checkData, []byte("0"))
	}
	return RegisterUpkeepContractsWithCheckData(
		t, client, linkToken, fundsForEachUpkeep, upkeepGasLimit, registry, registrar,
		numberOfContracts, upkeepAddresses, checkData, isLogTrigger, isMercury, isBillingTokenNative, wethToken)
}

type UpkeepID = *big.Int

type upkeepRegistrationResult struct {
	upkeepID UpkeepID
}

func (r upkeepRegistrationResult) GetResult() *big.Int {
	return r.upkeepID
}

type upkeepConfig struct {
	address string
	data    []byte
}

func RegisterUpkeepContractsWithCheckData(t *testing.T, client *seth.Client, linkToken contracts.LinkToken, fundsForEachUpkeep *big.Int, upkeepGasLimit uint32, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, numberOfContracts int, upkeepAddresses []string, checkData [][]byte, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken) []*big.Int {
	l := framework.L

	concurrency, err := GetAndAssertCorrectConcurrency(client, 1)
	require.NoError(t, err, "Insufficient concurrency to execute action")

	executor := ctf_concurrency.NewConcurrentExecutor[UpkeepID, upkeepRegistrationResult, upkeepConfig](l)

	configs := make([]upkeepConfig, 0)
	for i := range upkeepAddresses {
		configs = append(configs, upkeepConfig{address: upkeepAddresses[i], data: checkData[i]})
	}

	var registerUpkeepFn = func(resultCh chan upkeepRegistrationResult, errorCh chan error, executorNum int, config upkeepConfig) {
		id := uuid.New().String()
		keyNum := executorNum + 1 // key 0 is the root key
		var tx *types.Transaction

		if isBillingTokenNative {
			// register upkeep with native token
			tx, err = registrar.RegisterUpkeepFromKey(
				keyNum,
				"upkeep_"+id,
				[]byte("test@mail.com"),
				config.address,
				upkeepGasLimit,
				client.MustGetRootKeyAddress().Hex(), // upkeep Admin
				config.data,
				fundsForEachUpkeep,
				wethToken.Address(),
				isLogTrigger,
				isMercury,
			)
			if err != nil {
				errorCh <- errors.Wrapf(err, "[id: %s] Failed to register upkeep at %s", id, config.address)
				return
			}
		} else {
			// register upkeep with LINK
			req, err := registrar.EncodeRegisterRequest(
				"upkeep_"+id,
				[]byte("test@mail.com"),
				config.address,
				upkeepGasLimit,
				client.MustGetRootKeyAddress().Hex(), // upkeep Admin
				config.data,
				fundsForEachUpkeep,
				0,
				client.Addresses[keyNum].Hex(),
				isLogTrigger,
				isMercury,
				linkToken.Address(),
			)

			if err != nil {
				errorCh <- errors.Wrapf(err, "[id: %s] Failed to encode register request for upkeep at %s", id, config.address)
				return
			}

			balance, err := linkToken.BalanceOf(context.Background(), client.Addresses[keyNum].Hex())
			if err != nil {
				errorCh <- errors.Wrapf(err, "[id: %s]Failed to get LINK balance of %s", id, client.Addresses[keyNum].Hex())
				return
			}

			// not strictly necessary, but helps us to avoid an errorless revert if there is not enough LINK
			if balance.Cmp(fundsForEachUpkeep) < 0 {
				errorCh <- fmt.Errorf("[id: %s] Not enough LINK balance for %s. Has: %s. Needs: %s", id, client.Addresses[keyNum].Hex(), balance.String(), fundsForEachUpkeep.String())
				return
			}

			decodedTx, err := client.Decode(linkToken.TransferAndCallFromKey(registrar.Address(), fundsForEachUpkeep, req, keyNum))
			if err != nil {
				errorCh <- errors.Wrapf(err, "[id: %s] Failed to register upkeep at %s", id, config.address)
				return
			}
			tx = decodedTx.Transaction
		}

		// parse txn to get upkeep ID
		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			errorCh <- errors.Wrapf(err, "[id: %s] Failed to get receipt for upkeep at %s and tx hash %s", id, config.address, tx.Hash())
			return
		}

		var upkeepID *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepID, err := registry.ParseUpkeepIDFromRegisteredLog(rawLog)
			if err == nil {
				upkeepID = parsedUpkeepID
				break
			}
		}

		if upkeepID == nil {
			errorCh <- errors.Wrapf(err, "[id: %s] Failed find upkeep ID for upkeep at %s in logs of tx with hash %s", id, config.address, tx.Hash())
			return
		}

		l.Debug().
			Str("TxHash", tx.Hash().String()).
			Str("Upkeep ID", upkeepID.String()).
			Msg("Found upkeepId in tx hash")

		resultCh <- upkeepRegistrationResult{upkeepID: upkeepID}
	}

	upkeepIDs, err := executor.Execute(concurrency, configs, registerUpkeepFn)
	require.NoError(t, err, "Failed to register upkeeps using executor")

	require.Len(t, upkeepIDs, numberOfContracts, "Incorrect number of Keeper Consumer Contracts registered")
	l.Info().Msg("Successfully registered all Keeper Consumer Contracts")

	return upkeepIDs
}

type keeperConsumerResult struct {
	contract contracts.KeeperConsumer
}

func (k keeperConsumerResult) GetResult() contracts.KeeperConsumer {
	return k.contract
}

// DeployKeeperConsumers concurrently deploys keeper consumer contracts. It requires at least 1 ephemeral key to be present in Seth config.
func DeployKeeperConsumers(t *testing.T, client *seth.Client, numberOfContracts int, isLogTrigger bool, isMercury bool) []contracts.KeeperConsumer {
	l := framework.L

	concurrency, err := GetAndAssertCorrectConcurrency(client, 1)
	require.NoError(t, err, "Insufficient concurrency to execute action")

	executor := ctf_concurrency.NewConcurrentExecutor[contracts.KeeperConsumer, keeperConsumerResult, ctf_concurrency.NoTaskType](l)

	var deployContractFn = func(channel chan keeperConsumerResult, errorCh chan error, executorNum int) {
		keyNum := executorNum + 1 // key 0 is the root key
		var keeperConsumerInstance contracts.KeeperConsumer
		var err error

		switch {
		case isMercury && isLogTrigger:
			// v2.1 only: Log triggered based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggeredStreamsLookupUpkeepConsumerFromKey(client, keyNum)
		case isMercury:
			// v2.1 only: Conditional based contract with Mercury enabled
			keeperConsumerInstance, err = contracts.DeployAutomationStreamsLookupUpkeepConsumerFromKey(client, keyNum, big.NewInt(1000), big.NewInt(5), false, true, false) // 1000 block test range
		case isLogTrigger:
			// v2.1+: Log triggered based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployAutomationLogTriggerConsumerFromKey(client, keyNum, big.NewInt(1000)) // 1000 block test range
		default:
			// v2.0+: Conditional based contract without Mercury
			keeperConsumerInstance, err = contracts.DeployUpkeepCounterFromKey(client, keyNum, big.NewInt(999999), big.NewInt(5))
		}

		if err != nil {
			errorCh <- errors.Wrapf(err, "Failed to deploy keeper consumer contract")
			return
		}

		channel <- keeperConsumerResult{contract: keeperConsumerInstance}
	}

	results, err := executor.ExecuteSimple(concurrency, numberOfContracts, deployContractFn)
	require.NoError(t, err, "Failed to deploy keeper consumers")

	require.Len(t, results, numberOfContracts, "Incorrect number of Keeper Consumer Contracts deployed")
	l.Info().Msg("Successfully deployed all Keeper Consumer Contracts")

	return results
}

// SetupKeeperConsumers concurrently loads or deploys keeper consumer contracts. It requires at least 1 ephemeral key to be present in Seth config.
func SetupKeeperConsumers(t *testing.T, client *seth.Client, numberOfContracts int, isLogTrigger bool, isMercury bool, config Automation) []contracts.KeeperConsumer {
	l := framework.L

	results := []contracts.KeeperConsumer{}

	if len(config.DeployedContracts.Upkeeps) == 0 {
		// Deploy new contracts
		return DeployKeeperConsumers(t, client, numberOfContracts, isLogTrigger, isMercury)
	}

	require.Len(t, len(config.DeployedContracts.Upkeeps), numberOfContracts, "Incorrect number of Keeper Consumer Contracts loaded")
	l.Info().Int("Number of Contracts", numberOfContracts).Msg("Loading upkeep contracts from config")
	// Load existing contracts
	for i := range numberOfContracts {
		contract, err := contracts.LoadKeeperConsumer(client, common.HexToAddress(config.DeployedContracts.Upkeeps[i]))
		require.NoError(t, err, "Failed to load keeper consumer contract")
		l.Info().Str("Contract Address", contract.Address()).Int("Number", i+1).Int("Out Of", numberOfContracts).Msg("Loaded Keeper Consumer Contract")
		results = append(results, contract)
	}

	return results
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

	addressesOfNewUpkeeps := []string{}
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

	newUpkeepIDs := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfNewUpkeeps, addressesOfNewUpkeeps, false, false, false, nil)

	return newlyDeployedUpkeeps, newUpkeepIDs
}

var InsufficientStaticKeys = `
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

	if concurrency < minConcurrency {
		return 0, fmt.Errorf(InsufficientStaticKeys, concurrency)
	}

	return concurrency, nil
}
