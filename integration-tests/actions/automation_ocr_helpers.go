package actions

//revive:disable:dot-imports
import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	tt "github.com/smartcontractkit/chainlink/integration-tests/types"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
)

// DeployAutoOCRRegistryAndRegistrar registry and registrar
func DeployAutoOCRRegistryAndRegistrar(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	linkToken contracts.LinkToken,
	wethToken contracts.WETHToken,
	ethUSDFeed contracts.MockETHUSDFeed,
) (contracts.KeeperRegistry, contracts.KeeperRegistrar) {
	registry := deployRegistry(t, client, registryVersion, registrySettings, linkToken, wethToken, ethUSDFeed)
	registrar := deployRegistrar(t, client, registryVersion, registry, linkToken, wethToken)

	return registry, registrar
}

// DeployLegacyConsumers deploys and registers keeper consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkpeep will be registered with root key address as the admin. Which means
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deploys multicall contract and sends link funds to each deployment address.
func DeployLegacyConsumers(t *testing.T, chainClient *seth.Client, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken) ([]contracts.KeeperConsumer, []*big.Int) {
	// Fund deployers with LINK, no need to do this for Native token
	if !isBillingTokenNative {
		err := DeployMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep)
		require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")
	}

	upkeeps := DeployKeeperConsumers(t, chainClient, numberOfUpkeeps, isLogTrigger, isMercury)
	require.Equal(t, numberOfUpkeeps, len(upkeeps), "Number of upkeeps should match")
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(
		t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, isLogTrigger, isMercury, isBillingTokenNative, wethToken,
	)
	require.Equal(t, numberOfUpkeeps, len(upkeepIds), "Number of upkeepIds should match")
	return upkeeps, upkeepIds
}

// DeployConsumers deploys and registers keeper consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkpeep will be registered with root key address as the admin. Which means
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deploys multicall contract and sends link funds to each deployment address.
func DeployConsumers(t *testing.T, chainClient *seth.Client, registry contracts.KeeperRegistry, registrar contracts.KeeperRegistrar, linkToken contracts.LinkToken, numberOfUpkeeps int, linkFundsForEachUpkeep *big.Int, upkeepGasLimit uint32, isLogTrigger bool, isMercury bool, isBillingTokenNative bool, wethToken contracts.WETHToken, config tt.AutomationTestConfig) ([]contracts.KeeperConsumer, []*big.Int) {
	// Fund deployers with LINK, no need to do this for Native token
	if !isBillingTokenNative {
		err := SetupMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep, config)
		require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")
	}

	upkeeps := SetupKeeperConsumers(t, chainClient, numberOfUpkeeps, isLogTrigger, isMercury, config)
	require.Equal(t, numberOfUpkeeps, len(upkeeps), "Number of upkeeps should match")
	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(
		t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, isLogTrigger, isMercury, isBillingTokenNative, wethToken,
	)
	require.Equal(t, numberOfUpkeeps, len(upkeepIds), "Number of upkeepIds should match")
	return upkeeps, upkeepIds
}

// DeployPerformanceConsumers deploys and registers keeper performance consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkeep will be registered with root key address as the admin.
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deploys multicall contract and sends link funds to each deployment address.
func DeployPerformanceConsumers(
	t *testing.T,
	chainClient *seth.Client,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	upkeepGasLimit uint32,
	blockRange, // How many blocks to run the test for
	blockInterval, // Interval of blocks that upkeeps are expected to be performed
	checkGasToBurn, // How much gas should be burned on checkUpkeep() calls
	performGasToBurn int64, // How much gas should be burned on performUpkeep() calls
	config tt.AutomationTestConfig,
) ([]contracts.KeeperConsumerPerformance, []*big.Int) {
	upkeeps := DeployKeeperConsumersPerformance(
		t, chainClient, numberOfUpkeeps, blockRange, blockInterval, checkGasToBurn, performGasToBurn,
	)

	err := SetupMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep, config)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, false, false, false, nil)
	return upkeeps, upkeepIds
}

// DeployPerformDataCheckerConsumers deploys and registers keeper performance data checkers consumers. If ephemeral addresses are enabled, it will deploy and register the consumers from ephemeral addresses, but each upkpeep will be registered with root key address as the admin.
// that functions like setting upkeep configuration, pausing, unpausing, etc. will be done by the root key address. It deployes multicall contract and sends link funds to each deployment address.
func DeployPerformDataCheckerConsumers(
	t *testing.T,
	chainClient *seth.Client,
	registry contracts.KeeperRegistry,
	registrar contracts.KeeperRegistrar,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	upkeepGasLimit uint32,
	expectedData []byte,
	config tt.AutomationTestConfig,
) ([]contracts.KeeperPerformDataChecker, []*big.Int) {
	upkeeps := DeployPerformDataChecker(t, chainClient, numberOfUpkeeps, expectedData)

	err := SetupMultiCallAndFundDeploymentAddresses(chainClient, linkToken, numberOfUpkeeps, linkFundsForEachUpkeep, config)
	require.NoError(t, err, "Sending link funds to deployment addresses shouldn't fail")

	var upkeepsAddresses []string
	for _, upkeep := range upkeeps {
		upkeepsAddresses = append(upkeepsAddresses, upkeep.Address())
	}
	upkeepIds := RegisterUpkeepContracts(t, chainClient, linkToken, linkFundsForEachUpkeep, upkeepGasLimit, registry, registrar, numberOfUpkeeps, upkeepsAddresses, false, false, false, nil)
	return upkeeps, upkeepIds
}

func SetupMultiCallAddress(chainClient *seth.Client, testConfig tt.AutomationTestConfig) (common.Address, error) {
	if testConfig.GetAutomationConfig().UseExistingMultiCallContract() {
		multiCallAddress, err := testConfig.GetAutomationConfig().MultiCallContractAddress()
		if err != nil {
			return common.Address{}, errors.Wrap(err, "Error getting existing multicall contract address")
		}
		return multiCallAddress, nil
	}

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	if err != nil {
		return common.Address{}, errors.Wrap(err, "Error deploying multicall contract")
	}
	return multicallAddress, nil
}

// SetupMultiCallAndFundDeploymentAddresses setups multicall contract and sends link funds to each deployment address
func SetupMultiCallAndFundDeploymentAddresses(
	chainClient *seth.Client,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
	testConfig tt.AutomationTestConfig,
) error {
	concurrency, err := GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return err
	}

	operationsPerAddress := numberOfUpkeeps / concurrency

	multicallAddress, err := SetupMultiCallAddress(chainClient, testConfig)
	if err != nil {
		return errors.Wrap(err, "Error deploying multicall contract")
	}

	return SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
}

// DeployMultiCallAndFundDeploymentAddresses deploys multicall contract and sends link funds to each deployment address
func DeployMultiCallAndFundDeploymentAddresses(
	chainClient *seth.Client,
	linkToken contracts.LinkToken,
	numberOfUpkeeps int,
	linkFundsForEachUpkeep *big.Int,
) error {
	concurrency, err := GetAndAssertCorrectConcurrency(chainClient, 1)
	if err != nil {
		return err
	}

	operationsPerAddress := numberOfUpkeeps / concurrency

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	if err != nil {
		return errors.Wrap(err, "Error deploying multicall contract")
	}

	return SendLinkFundsToDeploymentAddresses(chainClient, concurrency, numberOfUpkeeps, operationsPerAddress, multicallAddress, linkFundsForEachUpkeep, linkToken)
}

func deployRegistrar(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	registry contracts.KeeperRegistry,
	linkToken contracts.LinkToken,
	wethToken contracts.WETHToken,
) contracts.KeeperRegistrar {
	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: 2,
		AutoApproveMaxAllowed: math.MaxUint16,
		RegistryAddr:          registry.Address(),
		MinLinkJuels:          big.NewInt(0),
		WETHTokenAddr:         wethToken.Address(),
	}
	registrar, err := contracts.DeployKeeperRegistrar(client, registryVersion, linkToken.Address(), registrarSettings)
	require.NoError(t, err, "Deploying KeeperRegistrar contract shouldn't fail")
	return registrar
}

func deployRegistry(
	t *testing.T,
	client *seth.Client,
	registryVersion ethereum.KeeperRegistryVersion,
	registrySettings contracts.KeeperRegistrySettings,
	linkToken contracts.LinkToken,
	wethToken contracts.WETHToken,
	ethUSDFeed contracts.MockETHUSDFeed,
) contracts.KeeperRegistry {
	ef, err := contracts.DeployMockLINKETHFeed(client, big.NewInt(2e18))
	require.NoError(t, err, "Deploying mock ETH-Link feed shouldn't fail")
	gf, err := contracts.DeployMockGASFeed(client, big.NewInt(2e11))
	require.NoError(t, err, "Deploying mock gas feed shouldn't fail")

	// Deploy the transcoder here, and then set it to the registry
	transcoder, err := contracts.DeployUpkeepTranscoder(client)
	require.NoError(t, err, "Deploying upkeep transcoder shouldn't fail")

	registry, err := contracts.DeployKeeperRegistry(
		client,
		&contracts.KeeperRegistryOpts{
			RegistryVersion:   registryVersion,
			LinkAddr:          linkToken.Address(),
			ETHFeedAddr:       ef.Address(),
			GasFeedAddr:       gf.Address(),
			TranscoderAddr:    transcoder.Address(),
			RegistrarAddr:     ZeroAddress.Hex(),
			Settings:          registrySettings,
			LinkUSDFeedAddr:   ethUSDFeed.Address(),
			NativeUSDFeedAddr: ethUSDFeed.Address(),
			WrappedNativeAddr: wethToken.Address(),
		},
	)
	require.NoError(t, err, "Deploying KeeperRegistry contract shouldn't fail")
	return registry
}
