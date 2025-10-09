package cre

import (
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	tron_df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	tron_keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/tron"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const PoRWFV1Location = "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
const PoRWFV2Location = "../../../../core/scripts/cre/environment/examples/workflows/v2/proof-of-reserve/cron-based/main.go"

type WorkflowTestConfig struct {
	WorkflowName         string
	WorkflowFileLocation string
	FeedIDs              []string
}

func beforePoRTest(t *testing.T, testEnv *ttypes.TestEnvironment, workflowName, workflowLocation string) (PriceProvider, WorkflowTestConfig) {
	porWfCfg := WorkflowTestConfig{
		FeedIDs:              []string{"018e16c38e000320000000000000000000000000000000000000000000000000", "018e16c39e000320000000000000000000000000000000000000000000000000"},
		WorkflowName:         workflowName,
		WorkflowFileLocation: workflowLocation,
	}
	// AuthorizationKeySecretName := "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	// It is needed for FakePriceProvider

	testLogger := framework.L
	AuthorizationKey := "" // required by FakePriceProvider
	priceProvider, err := NewFakePriceProvider(testLogger, testEnv.Config.Fake, AuthorizationKey, porWfCfg.FeedIDs)
	require.NoError(t, err, "failed to create fake price provider")

	return priceProvider, porWfCfg
}

func ExecutePoRTest(t *testing.T, testEnv *ttypes.TestEnvironment, priceProvider PriceProvider, cfg WorkflowTestConfig, withBilling bool) {
	testLogger := framework.L
	blockchainOutputs := testEnv.WrappedBlockchainOutputs

	var billingState billingAssertionState
	if withBilling {
		billingState = getBillingAssertionState(t, testEnv.TestConfig.RelativePathToRepoRoot) // establish a baseline
	}

	writeableChains := t_helpers.GetWritableChainsFromSavedEnvironmentState(t, testEnv)
	require.Len(t, cfg.FeedIDs, len(writeableChains), "a number of writeable chains must match the number of feed IDs (check what chains 'evm' and 'write-evm' capabilities are enabled for)")

	/*
		DEPLOY DATA FEEDS CACHE + READ BALANCES CONTRACTS ON ALL CHAINS (except read-only ones)
		Workflow will write price data to the data feeds cache contract

		REGISTER ONE WORKFLOW PER CHAIN (except read-only ones)
	*/

	// amountToFund is moved to the outer scope to correctly count the final amount sent
	// to the requested number of new addresses used to read balances from in the PoR workflow.
	// This amount is added to the prices from the (http) PriceProvider,
	// forming the final PoR "expected" total price written on-chain.
	var amountToFund *big.Int
	numberOfAddressesToCreate := 2
	var workflowOwner common.Address
	for idx, bcOutput := range blockchainOutputs {
		chainFamily := bcOutput.BlockchainOutput.Family
		chainID := bcOutput.ChainID
		chainSelector := bcOutput.ChainSelector
		chainType := bcOutput.BlockchainOutput.Type
		perChainSethClient := bcOutput.SethClient
		creEnvironment := testEnv.CreEnvironment
		feedID := cfg.FeedIDs[idx]

		if chainType == blockchain.FamilySolana {
			continue
		}

		// Deploy Data Feeds Cache contract only on chains that are writable
		if !slices.Contains(writeableChains, chainID) {
			continue
		}

		var dataFeedsCacheAddress common.Address
		var readBalancesAddress common.Address

		uniqueWorkflowName := cfg.WorkflowName + "-" + bcOutput.BlockchainOutput.ChainID + "-" + uuid.New().String()[0:4]                                                                       // e.g. 'por-workflow-1337-5f37_config'
		forwarderAddress, _, forwarderErr := crecontracts.FindAddressesForChain(creEnvironment.CldfEnvironment.ExistingAddresses, chainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		require.NoError(t, forwarderErr, "failed to find Forwarder address for chain %d", chainSelector)

		switch chainFamily {
		case blockchain.FamilyTron:
			dataFeedsCacheAddress, readBalancesAddress = deployAndConfigureTronContracts(t, testLogger, chainSelector, creEnvironment, workflowOwner, uniqueWorkflowName, feedID, forwarderAddress)
			chainFamily = blockchain.FamilyEVM
		default:
			workflowOwner = bcOutput.SethClient.MustGetRootKeyAddress()
			dataFeedsCacheAddress, readBalancesAddress = deployAndConfigureEVMContracts(t, testLogger, chainSelector, chainID, creEnvironment, workflowOwner, uniqueWorkflowName, feedID, forwarderAddress)
		}

		// reset to avoid incrementing on each iteration
		amountToFund = big.NewInt(0).SetUint64(10) // 10 wei
		addressesToRead, addrErr := t_helpers.CreateAndFundAddresses(t, testLogger, numberOfAddressesToCreate, amountToFund, perChainSethClient, bcOutput, creEnvironment)
		require.NoError(t, addrErr, "failed to create and fund addresses to read")

		testLogger.Info().Msg("Creating PoR workflow configuration file...")
		writeTargetName := corevm.GenerateWriteTargetName(chainID)
		testLogger.Info().Msgf("Generated WriteTargetName for chain %d (%s): %s", chainID, chainFamily, writeTargetName)

		workflowConfig := portypes.WorkflowConfig{
			ChainFamily:   chainFamily,
			ChainID:       strconv.FormatUint(chainID, 10),
			ChainSelector: chainSelector,
			BalanceReaderConfig: portypes.BalanceReaderConfig{
				BalanceReaderAddress: readBalancesAddress.Hex(),
				AddressesToRead:      addressesToRead,
			},
			ComputeConfig: portypes.ComputeConfig{
				FeedID:                feedID,
				URL:                   priceProvider.URL(),
				DataFeedsCacheAddress: dataFeedsCacheAddress.Hex(),
				WriteTargetName:       writeTargetName,
			},
		}
		testLogger.Info().Msgf("Workflow config for chain %d: WriteTarget=%s, DataFeedsCache=%s, FeedID: %s", chainID, writeTargetName, dataFeedsCacheAddress.Hex(), feedID)
		workflowFileLocation := cfg.WorkflowFileLocation

		t_helpers.CompileAndDeployWorkflow(t, testEnv, testLogger, uniqueWorkflowName, &workflowConfig, workflowFileLocation)
	}
	/*
		START THE VALIDATION PHASE
		Check whether each feed has been updated with the expected prices, which workflow fetches from the price provider
	*/
	// final expected total = amount to fund * the number of addresses to create
	amountToFund.Mul(amountToFund, big.NewInt(int64(numberOfAddressesToCreate)))
	validatePoRPrices(t, testEnv, priceProvider, &cfg, *amountToFund)

	if withBilling {
		expectedMinChange := float64(49)
		assertBillingStateChanged(t, billingState, 2*time.Minute, expectedMinChange)
	}
}

func deployAndConfigureEVMContracts(t *testing.T, testLogger zerolog.Logger, chainSelector uint64, chainID uint64, creEnvironment *cre.Environment, workflowOwner common.Address, uniqueWorkflowName string, feedID string, forwarderAddress common.Address) (common.Address, common.Address) {
	testLogger.Info().Msgf("Deploying additional contracts to chain %d (%d)", chainID, chainSelector)
	dfAddress, dfOutput, dfErr := crecontracts.DeployDataFeedsCacheContract(testLogger, chainSelector, creEnvironment)
	require.NoError(t, dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)

	rbAddress, rbOutput, rbErr := crecontracts.DeployReadBalancesContract(testLogger, chainSelector, creEnvironment)

	require.NoError(t, rbErr, "failed to deploy Read Balances contract on chain %d", chainSelector)

	crecontracts.MergeAllDataStores(creEnvironment, dfOutput, rbOutput)

	testLogger.Info().Msgf("Configuring Data Feeds Cache contract...")

	configInput := &cre.ConfigureDataFeedsCacheInput{
		CldEnv:                creEnvironment.CldfEnvironment,
		ChainSelector:         chainSelector,
		FeedIDs:               []string{feedID},
		Descriptions:          []string{"PoR test feed"},
		DataFeedsCacheAddress: dfAddress,
		AdminAddress:          workflowOwner,
		AllowedSenders:        []common.Address{forwarderAddress},
		AllowedWorkflowNames:  []string{uniqueWorkflowName},
		AllowedWorkflowOwners: []common.Address{workflowOwner},
	}
	_, dfConfigErr := crecontracts.ConfigureDataFeedsCache(testLogger, configInput)
	require.NoError(t, dfConfigErr, "failed to configure Data Feeds Cache contract")
	testLogger.Info().Msg("Data Feeds Cache contract configured successfully.")

	return dfAddress, rbAddress
}

func deployAndConfigureTronContracts(t *testing.T, testLogger zerolog.Logger, chainSelector uint64, creEnvironment *cre.Environment, workflowOwner common.Address, uniqueWorkflowName string, feedID string, forwarderAddress common.Address) (common.Address, common.Address) {
	// Use Tron-specific changeset with deploy options
	deployOptions := cldf_tron.DefaultDeployOptions()
	deployOptions.FeeLimit = 1_000_000_000

	tronDeployConfig := df_changeset_types.DeployTronConfig{
		ChainsToDeploy: []uint64{chainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
		DeployOptions:  deployOptions,
	}

	dfOutput, dfErr := changeset.RunChangeset(tron_df_changeset.DeployCacheChangeset, *creEnvironment.CldfEnvironment, tronDeployConfig)
	require.NoError(t, dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)

	rbOutput, rbErr := changeset.RunChangeset(tron_keystone_changeset.DeployReadBalanceChangeset, *creEnvironment.CldfEnvironment, tronDeployConfig)
	require.NoError(t, rbErr, "failed to deploy Read Balances contract on chain %d", chainSelector)

	crecontracts.MergeAllDataStores(creEnvironment, dfOutput, rbOutput)

	// Get DataFeedsCache address from merged DataStore
	dfAddressRefs := creEnvironment.CldfEnvironment.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(df_changeset.DataFeedsCache),
	)
	require.Len(t, dfAddressRefs, 1, "DataFeedsCache address not found in merged DataStore for chain %d", chainSelector)
	dataFeedsCacheAddress := common.HexToAddress(dfAddressRefs[0].Address)

	// Get BalanceReader address from merged DataStore
	rbAddressRefs := creEnvironment.CldfEnvironment.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType("BalanceReader"),
	)
	require.Len(t, rbAddressRefs, 1, "BalanceReader address not found in merged DataStore for chain %d", chainSelector)
	readBalancesAddress := common.HexToAddress(rbAddressRefs[0].Address)

	testLogger.Info().Msgf("Tron DataFeedsCache address: %s", dataFeedsCacheAddress.Hex())
	testLogger.Info().Msgf("Tron BalanceReader address: %s", readBalancesAddress.Hex())

	tronChains := creEnvironment.CldfEnvironment.BlockChains.TronChains()
	tronChain, exists := tronChains[chainSelector]
	require.True(t, exists, "Tron chain %d not found in environment", chainSelector)

	triggerOptions := cldf_tron.DefaultTriggerOptions()
	triggerOptions.FeeLimit = 1_000_000_000

	setDeployerAdminConfig := df_changeset_types.SetFeedAdminTronConfig{
		ChainSelector:  chainSelector,
		CacheAddress:   address.EVMAddressToAddress(dataFeedsCacheAddress),
		AdminAddress:   tronChain.Address, // Deployer address (equivalent to MustGetRootKeyAddress)
		IsAdmin:        true,
		TriggerOptions: triggerOptions,
	}

	_, setDeployerAdminErr := changeset.RunChangeset(tron_df_changeset.SetFeedAdminChangeset, *creEnvironment.CldfEnvironment, setDeployerAdminConfig)
	require.NoError(t, setDeployerAdminErr, "failed to set deployer as admin for Tron chain")

	workflowNameBytes := df_changeset.HashedWorkflowName(uniqueWorkflowName)

	workflowMetadata := []df_changeset_types.DataFeedsCacheTronWorkflowMetadata{
		{
			AllowedSender:        address.EVMAddressToAddress(forwarderAddress),
			AllowedWorkflowOwner: address.EVMAddressToAddress(workflowOwner), // Use home chain's deployer address for consistency
			AllowedWorkflowName:  workflowNameBytes,
		},
	}

	feedIDTruncated := feedID
	feedIDTruncated = strings.TrimPrefix(feedIDTruncated, "0x")
	if len(feedIDTruncated) > 32 {
		feedIDTruncated = feedIDTruncated[:32]
	}

	setFeedConfigConfig := df_changeset_types.SetFeedDecimalTronConfig{
		ChainSelector:    chainSelector,
		CacheAddress:     address.EVMAddressToAddress(dataFeedsCacheAddress),
		DataIDs:          []string{feedIDTruncated},
		Descriptions:     []string{"PoR test feed"},
		WorkflowMetadata: workflowMetadata,
		TriggerOptions:   triggerOptions,
	}

	_, setConfigErr := changeset.RunChangeset(tron_df_changeset.SetFeedConfigChangeset, *creEnvironment.CldfEnvironment, setFeedConfigConfig)
	require.NoError(t, setConfigErr, "failed to set feed config for Tron chain")

	testLogger.Info().Msgf("Successfully configured Tron data feeds cache for chain %d", chainSelector)

	return dataFeedsCacheAddress, readBalancesAddress
}

func validateTronPrices(t *testing.T, testEnv *ttypes.TestEnvironment, bcOutput *cre.WrappedBlockchainOutput, feedID string, priceProvider PriceProvider, startTime time.Time, waitFor time.Duration, tick time.Duration) error {
	dfAddressRefs := testEnv.CreEnvironment.CldfEnvironment.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(bcOutput.ChainSelector),
		datastore.AddressRefByType(df_changeset.DataFeedsCache),
	)

	if len(dfAddressRefs) == 0 {
		return fmt.Errorf("DataFeedsCache address not found in DataStore for chain %d", bcOutput.ChainSelector)
	}

	dataFeedsCacheAddresses := common.HexToAddress(dfAddressRefs[0].Address)

	tronChains := testEnv.CreEnvironment.CldfEnvironment.BlockChains.TronChains()
	tronChain, exists := tronChains[bcOutput.ChainSelector]
	if !exists {
		return fmt.Errorf("Tron chain %d not found in environment", bcOutput.ChainSelector)
	}

	cacheAddr := address.EVMAddressToAddress(dataFeedsCacheAddresses)
	testEnv.Logger.Info().Msgf("Tron chain %d: Contract address conversion - EVM: %s -> Tron: %s", bcOutput.ChainSelector, dataFeedsCacheAddresses.Hex(), cacheAddr.String())

	require.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)

		accountInfo, accountErr := tronChain.Client.GetAccount(cacheAddr)
		if accountErr != nil {
			testEnv.Logger.Error().Err(accountErr).Msgf("Tron chain %d: Failed to get account info for contract %s", bcOutput.ChainSelector, cacheAddr.String())
			return false
		}

		if accountInfo == nil || len(accountInfo.Address) == 0 {
			testEnv.Logger.Error().Msgf("Tron chain %d: Contract %s does not exist or is not deployed", bcOutput.ChainSelector, cacheAddr.String())
			return false
		}

		testEnv.Logger.Info().Msgf("Tron chain %d: Calling getLatestAnswer for feed %s on contract %s", bcOutput.ChainSelector, feedID, cacheAddr.String())

		result, err := tronChain.Client.TriggerConstantContract(
			tronChain.Address,          // caller address
			cacheAddr,                  // contract address
			"getLatestAnswer(bytes16)", // function signature
			[]any{"bytes16", [16]byte(common.Hex2Bytes(feedID))}, // parameters
		)
		if err != nil {
			testEnv.Logger.Error().Err(err).Msgf("FAILED to call getLatestAnswer on Tron chain %d", bcOutput.ChainSelector)
			return false
		}

		testEnv.Logger.Info().Msgf("Tron chain %d: Got result from contract call: %+v", bcOutput.ChainSelector, result)

		if len(result.ConstantResult) == 0 {
			testEnv.Logger.Error().Msgf("NO RESULT from getLatestAnswer on Tron chain %d", bcOutput.ChainSelector)
			return false
		}

		priceBytes := result.ConstantResult[0]
		if len(priceBytes) == 0 {
			testEnv.Logger.Error().Msgf("EMPTY price result from Tron chain %d", bcOutput.ChainSelector)
			return false
		}

		testEnv.Logger.Info().Msgf("Tron chain %d: Raw price bytes: %s", bcOutput.ChainSelector, priceBytes)

		price := new(big.Int)
		if len(priceBytes) >= 2 && priceBytes[:2] == "0x" {
			price.SetString(priceBytes[2:], 16)
		} else {
			price.SetString(priceBytes, 16)
		}

		testEnv.Logger.Info().Msgf("Tron chain %d: Parsed price %s for feed %s", bcOutput.ChainSelector, price.String(), feedID)

		return !priceProvider.NextPrice(feedID, price, elapsed)
	}, waitFor, tick, "feed %s did not update, timeout after: %s", feedID, waitFor.String())

	return nil
}

// validatePoRPrices validates that all feeds receive the expected prices from the price provider
func validatePoRPrices(t *testing.T, testEnv *ttypes.TestEnvironment, priceProvider PriceProvider, config *WorkflowTestConfig, additionalPrice big.Int) {
	t.Helper()
	eg := &errgroup.Group{}

	for idx, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}

		eg.Go(func() error {
			feedID := config.FeedIDs[idx]
			testEnv.Logger.Info().Msgf("Waiting for feed %s to update...", feedID)

			startTime := time.Now()
			waitFor := 5 * time.Minute
			tick := 5 * time.Second

			switch bcOutput.BlockchainOutput.Family {
			case blockchain.FamilyTron:
				if err := validateTronPrices(t, testEnv, bcOutput, feedID, priceProvider, startTime, waitFor, tick); err != nil {
					return err
				}
			case blockchain.FamilyEVM:
				if err := validateEVMPrices(t, testEnv, bcOutput, feedID, priceProvider, startTime, waitFor, tick); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported blockchain family: %s", bcOutput.BlockchainOutput.Family)
			}

			ppExpectedPrices := priceProvider.ExpectedPrices(feedID)
			expected := totalPoRExpectedPrices(ppExpectedPrices, &additionalPrice)
			actual := priceProvider.ActualPrices(feedID)

			testEnv.Logger.Info().Msgf("Feed %s - expected): %v", feedID, expected)
			testEnv.Logger.Info().Msgf("Feed %s - actual: %v", feedID, actual)

			if len(expected) != len(actual) {
				return fmt.Errorf("expected %d prices, got %d", len(expected), len(actual))
			}

			for i := range expected {
				if expected[i].Cmp(actual[i]) != 0 {
					return fmt.Errorf("expected price %d, got %d", expected[i], actual[i])
				}
			}

			testEnv.Logger.Info().Msgf("All prices were found in the feed %s", feedID)
			return nil
		})
	}

	err := eg.Wait()
	require.NoError(t, err, "price validation failed")

	testEnv.Logger.Info().Msgf("All prices were found for all feeds")
}

func validateEVMPrices(t *testing.T, testEnv *ttypes.TestEnvironment, bcOutput *cre.WrappedBlockchainOutput, feedID string, priceProvider PriceProvider, startTime time.Time, waitFor time.Duration, tick time.Duration) error {
	dataFeedsCacheAddresses, _, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
		testEnv.CreEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		bcOutput.ChainSelector,
		df_changeset.DataFeedsCache.String(),
	)
	if dataFeedsCacheErr != nil {
		return fmt.Errorf("failed to find Data Feeds Cache address for chain %d: %w", bcOutput.ChainID, dataFeedsCacheErr)
	}

	dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, bcOutput.SethClient.Client)
	if instanceErr != nil {
		return fmt.Errorf("failed to create Data Feeds Cache instance: %w", instanceErr)
	}

	require.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, err := dataFeedsCacheInstance.GetLatestAnswer(bcOutput.SethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(feedID)))
		if err != nil {
			testEnv.Logger.Error().Err(err).Msg("failed to get price from Data Feeds Cache contract")
			return false
		}

		// if there are no more prices to be found, we can stop waiting
		return !priceProvider.NextPrice(feedID, price, elapsed)
	}, waitFor, tick, "feed %s did not update, timeout after: %s", feedID, waitFor.String())

	return nil
}

// Adds the additional price (if any) to each expected price since it's included in actual prices
func totalPoRExpectedPrices(ppExpectedPrices []*big.Int, additionalPrice *big.Int) []*big.Int {
	expected := make([]*big.Int, len(ppExpectedPrices))
	for i, price := range ppExpectedPrices {
		expected[i] = new(big.Int).Add(price, additionalPrice)
	}
	return expected
}
