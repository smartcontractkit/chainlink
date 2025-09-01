package cre

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

type WorkflowTestConfig struct {
	WorkflowName         string
	WorkflowFileLocation string
	FeedIDs              []string
}

func ExecutePoRTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	// AuthorizationKeySecretName := "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	// It is needed for FakePriceProvider
	AuthorizationKey := "" // required by FakePriceProvider
	PoRWorkflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
	blockchainOutputs := testEnv.WrappedBlockchainOutputs
	baseWorkflowName := "por-workflow"
	feedIDs := []string{"018e16c39e000320000000000000000000000000000000000000000000000000", "018e16c38e000320000000000000000000000000000000000000000000000000"}
	baseWorkflowTestConfig := &WorkflowTestConfig{
		WorkflowName:         baseWorkflowName,
		WorkflowFileLocation: PoRWorkflowFileLocation,
		FeedIDs:              feedIDs,
	}

	writeableChains := getWritableChainsFromSavedEnvironmentState(t, testEnv)
	require.Len(t, baseWorkflowTestConfig.FeedIDs, len(writeableChains), "a number of writeable chains must match the number of feed IDs (check what chains 'evm' and 'write-evm' capabilities are enabled for)")

	priceProvider, err := NewFakePriceProvider(testLogger, testEnv.Config.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, err, "failed to create fake price provider")

	/*
		DEPLOY DATA FEEDS CACHE + READ BALANCES CONTRACTS ON ALL CHAINS (except read-only ones)
		Workflow will write price data to the data feeds cache contract

		REGISTER ONE WORKFLOW PER CHAIN (except read-only ones)
	*/
	for idx, bcOutput := range blockchainOutputs {
		chainFamily := bcOutput.BlockchainOutput.Family
		chainID := bcOutput.ChainID
		chainSelector := bcOutput.ChainSelector
		chainType := bcOutput.BlockchainOutput.Type
		fullCldEnvOutput := testEnv.FullCldEnvOutput
		feedID := baseWorkflowTestConfig.FeedIDs[idx]

		if chainType == blockchain.FamilySolana {
			continue
		}

		// Deploy Data Feeds Cache contract only on chains that are writable
		if !slices.Contains(writeableChains, chainID) {
			continue
		}

		testLogger.Info().Msgf("Deploying additional contracts to chain %d (%d)", chainID, chainSelector)
		dataFeedsCacheAddress, dfOutput, dfErr := crecontracts.DeployDataFeedsCacheContract(testLogger, chainSelector, fullCldEnvOutput)
		require.NoError(t, dfErr, "failed to deploy Data Feeds Cache contract on chain %d", chainSelector)

		readBalancesAddress, rbOutput, rbErr := crecontracts.DeployReadBalancesContract(testLogger, chainSelector, fullCldEnvOutput)
		require.NoError(t, rbErr, "failed to deploy Read Balances contract on chain %d", chainSelector)
		crecontracts.MergeAllDataStores(fullCldEnvOutput, dfOutput, rbOutput)

		testLogger.Info().Msgf("Configuring Data Feeds Cache contract...")
		forwarderAddress, forwarderErr := crecontracts.FindAddressesForChain(fullCldEnvOutput.Environment.ExistingAddresses, chainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
		require.NoError(t, forwarderErr, "failed to find Forwarder address for chain %d", chainSelector)

		uniqueWorkflowName := baseWorkflowTestConfig.WorkflowName + "-" + bcOutput.BlockchainOutput.ChainID + "-" + uuid.New().String()[0:4] // e.g. 'por-workflow-1337-5f37_config'
		configInput := &cre.ConfigureDataFeedsCacheInput{
			CldEnv:                fullCldEnvOutput.Environment,
			ChainSelector:         chainSelector,
			FeedIDs:               []string{feedID},
			Descriptions:          []string{"PoR test feed"},
			DataFeedsCacheAddress: dataFeedsCacheAddress,
			AdminAddress:          bcOutput.SethClient.MustGetRootKeyAddress(),
			AllowedSenders:        []common.Address{forwarderAddress},
			AllowedWorkflowNames:  []string{uniqueWorkflowName},
			AllowedWorkflowOwners: []common.Address{bcOutput.SethClient.MustGetRootKeyAddress()},
		}
		_, dfConfigErr := crecontracts.ConfigureDataFeedsCache(testLogger, configInput)
		require.NoError(t, dfConfigErr, "failed to configure Data Feeds Cache contract")
		testLogger.Info().Msg("Data Feeds Cache contract configured successfully.")

		testLogger.Info().Msg("Creating PoR workflow configuration file...")
		workflowConfig := portypes.WorkflowConfig{
			ChainFamily: chainFamily,
			ChainID:     strconv.FormatUint(chainID, 10),
			BalanceReaderConfig: portypes.BalanceReaderConfig{
				BalanceReaderAddress: readBalancesAddress.Hex(),
			},
			ComputeConfig: portypes.ComputeConfig{
				FeedID:                feedID,
				URL:                   priceProvider.URL(),
				DataFeedsCacheAddress: dataFeedsCacheAddress.Hex(),
				WriteTargetName:       corevm.GenerateWriteTargetName(chainID),
			},
		}
		workflowFileLocation := baseWorkflowTestConfig.WorkflowFileLocation

		compileAndDeployWorkflow(t, testEnv, testLogger, uniqueWorkflowName, &workflowConfig, workflowFileLocation)
	}
	/*
		START THE VALIDATION PHASE
		Check whether each feed has been updated with the expected prices, which workflow fetches from the price provider
	*/
	validatePoRPrices(t, testEnv, priceProvider, baseWorkflowTestConfig)
}

/*
Creates .yaml workflow configuration file.
It stores the values used by a workflow (main.go),
(i.e. feedID, read/write contract addresses)

The values are written to types.WorkflowConfig.
The method returns the absolute path to the created config file.
*/
func createPoRWorkflowConfigFile(workflowName string, workflowConfig *portypes.WorkflowConfig) (string, error) {
	feedIDToUse, fIDerr := validateAndFormatFeedID(workflowConfig)
	if fIDerr != nil {
		return "", errors.Wrap(fIDerr, "failed to validate and format feed ID")
	}
	workflowConfig.FeedID = feedIDToUse

	// Write workflow config to a .yaml file
	configMarshalled, err := yaml.Marshal(workflowConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal workflow config")
	}
	workflowSuffix := "_config.yaml"
	workflowConfigOutputFile := workflowName + workflowSuffix

	// remove the duplicate if it already exists
	_, statErr := os.Stat(workflowConfigOutputFile)
	if statErr == nil {
		if err := os.Remove(workflowConfigOutputFile); err != nil {
			return "", errors.Wrap(err, "failed to remove existing output file")
		}
	}

	if err := os.WriteFile(workflowConfigOutputFile, configMarshalled, 0644); err != nil { //nolint:gosec // G306: we want it to be readable by everyone
		return "", errors.Wrap(err, "failed to write output file")
	}

	outputFileAbsPath, outputFileAbsPathErr := filepath.Abs(workflowConfigOutputFile)
	if outputFileAbsPathErr != nil {
		return "", errors.Wrap(outputFileAbsPathErr, "failed to get absolute path of the config file")
	}

	return outputFileAbsPath, nil
}

func validateAndFormatFeedID(workflowConfig *portypes.WorkflowConfig) (string, error) {
	feedID := workflowConfig.FeedID

	// validate and format feed ID to fit 32 bytes
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedIDLength := len(cleanFeedID)
	if feedIDLength < 32 {
		return "", errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedIDLength)
	}

	if feedIDLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	// override feed ID in workflow config to ensure it is exactly 32 bytes
	feedIDToUse := "0x" + cleanFeedID
	return feedIDToUse, nil
}

// validatePoRPrices validates that all feeds receive the expected prices from the price provider
func validatePoRPrices(t *testing.T, testEnv *TestEnvironment, priceProvider PriceProvider, config *WorkflowTestConfig) {
	t.Helper()
	eg := &errgroup.Group{}

	for idx, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}
		eg.Go(func() error {
			feedID := config.FeedIDs[idx]
			testEnv.Logger.Info().Msgf("Waiting for feed %s to update...", feedID)

			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
				testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck,nolintlint // SA1019: deprecated but we don't want to migrate now
				bcOutput.ChainSelector,
				df_changeset.DataFeedsCache.String(),
			)
			if dataFeedsCacheErr != nil {
				return fmt.Errorf("failed to find data feeds cache address for chain %d: %w", bcOutput.ChainID, dataFeedsCacheErr)
			}

			dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, bcOutput.SethClient.Client)
			if instanceErr != nil {
				return fmt.Errorf("failed to create data feeds cache instance: %w", instanceErr)
			}

			startTime := time.Now()
			waitFor := 5 * time.Minute
			tick := 5 * time.Second
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

			expected := priceProvider.ExpectedPrices(feedID)
			actual := priceProvider.ActualPrices(feedID)

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
