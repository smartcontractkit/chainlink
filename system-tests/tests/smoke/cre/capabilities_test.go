package cre

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	tron_df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/tron"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	cldf_tron "github.com/smartcontractkit/chainlink-deployments-framework/chain/tron"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
)

/*
To execute on local start the local CRE first with following command:
# inside core/scripts/cre/environment directory
go run . env start
*/
func Test_CRE_Workflow_Don(t *testing.T) {
	// TODO: remove this temp test env vars
	os.Setenv("CTF_CONFIGS", "../../../../core/scripts/cre/environment/configs/capability_defaults-cache.toml")
	os.Setenv("ENV_ARTIFACT_PATH", "../../../../core/scripts/cre/environment/env_artifact/env_artifact.json")
	confErr := setConfigurationIfMissing("../../../../core/scripts/cre/environment/configs/workflow-don-tron.toml", "workflow")
	require.NoError(t, confErr, "failed to set configuration")

	configurationFiles := os.Getenv("CTF_CONFIGS")
	require.NotEmpty(t, configurationFiles, "CTF_CONFIGS env var is not set")

	topology := os.Getenv("CRE_TOPOLOGY")
	require.NotEmpty(t, topology, "CRE_TOPOLOGY env var is not set")

	createErr := createEnvironmentIfNotExists(configurationFiles, "../../../../core/scripts/cre/environment", topology)
	require.NoError(t, createErr, "failed to create environment")

	/*
		LOAD ENVIRONMENT STATE
	*/
	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")

	var envArtifact environment.EnvArtifact
	artFile, err := os.ReadFile(os.Getenv("ENV_ARTIFACT_PATH"))
	fmt.Println("ENV_ARTIFACT_PATH", os.Getenv("ENV_ARTIFACT_PATH"))
	fmt.Println("CTF_CONFIGS", os.Getenv("CTF_CONFIGS"))
	require.NoError(t, err, "failed to read artifact file")
	err = json.Unmarshal(artFile, &envArtifact)
	require.NoError(t, err, "failed to unmarshal artifact file")

	// currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("cron-based PoR workflow", func(t *testing.T) {
		executePoRTest(t, in, envArtifact, 5*time.Minute)
	})

	t.Run("vault DON test", func(t *testing.T) {
		executeVaultTest(t, in, envArtifact)
	})

	t.Run("DON Time test", func(t *testing.T) {
		// TODO: Implement smoke test - https://smartcontract-it.atlassian.net/browse/CAPPL-1028
		t.Skip()
	})
}

func executePoRTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact, verificationTimeout time.Duration) {
	testLogger := framework.L
	cldLogger := cldlogger.NewSingleFileLogger(t)

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
	feedIDs := []string{"018e16c39e000320000000000000000000000000000000000000000000000000", "018e16c38e000320000000000000000000000000000000000000000000000000"}

	priceProvider, priceErr := NewFakePriceProvider(framework.L, in.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	fullCldEnvOutput, wrappedBlockchainOutputs, loadErr := environment.BuildFromSavedState(t.Context(), cldLogger, in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

	homeChainSelector := wrappedBlockchainOutputs[0].ChainSelector
	homeChainSethClient := wrappedBlockchainOutputs[0].SethClient // Get home chain's seth client for workflow owner
	writeableChains := []uint64{}
	for _, bcOutput := range wrappedBlockchainOutputs {
		for _, donMetadata := range fullCldEnvOutput.DonTopology.DonsWithMetadata {
			if flags.RequiresForwarderContract(donMetadata.Flags, bcOutput.ChainID) {
				if !slices.Contains(writeableChains, bcOutput.ChainID) {
					writeableChains = append(writeableChains, bcOutput.ChainID)
				}
			}
		}
	}
	require.Len(t, feedIDs, len(writeableChains), "number of writeable chains must match number of feed IDs (check what chains 'evm' and 'write-evm' capabilities are enabled for)")

	/*
		DEPLOY DATA FEEDS CACHE CONTRACTS ON ALL CHAINS (except read-only ones)
		Workflow will write price data to the data feeds cache contract

		REGISTER ONE WORKFLOW PER CHAIN (except read-only ones)
	*/
	for idx, bcOutput := range wrappedBlockchainOutputs {
		// deploy data feeds cache contract only on chains that require a forwarder contract. It's required for the PoR workflow to work and we treat it as a proxy
		// for deciding whether need to deploy the data feeds cache contract.
		hasForwarderContract := false
		for _, donMetadata := range fullCldEnvOutput.DonTopology.DonsWithMetadata {
			if flags.RequiresForwarderContract(donMetadata.Flags, bcOutput.ChainID) {
				hasForwarderContract = true
				break
			}
		}

		if !hasForwarderContract {
			continue
		}

		var dfOutput cldf.ChangesetOutput
		var dfErr error

		// Use different changesets for Tron vs EVM chains
		if bcOutput.BlockchainOutput.Family == "tron" {
			// Use Tron-specific changeset with deploy options
			deployOptions := cldf_tron.DefaultDeployOptions()
			deployOptions.FeeLimit = 1_000_000_000

			tronDeployConfig := df_changeset_types.DeployTronConfig{
				ChainsToDeploy: []uint64{bcOutput.ChainSelector},
				Labels:         []string{"data-feeds"}, // label required by the changeset
				DeployOptions:  deployOptions,
			}

			dfOutput, dfErr = changeset.RunChangeset(tron_df_changeset.DeployCacheChangeset, *fullCldEnvOutput.Environment, tronDeployConfig)
		} else {
			// Use standard EVM changeset
			deployConfig := df_changeset_types.DeployConfig{
				ChainsToDeploy: []uint64{bcOutput.ChainSelector},
				Labels:         []string{"data-feeds"}, // label required by the changeset
			}

			dfOutput, dfErr = changeset.RunChangeset(df_changeset.DeployCacheChangeset, *fullCldEnvOutput.Environment, deployConfig)
		}
		require.NoError(t, dfErr, "failed to deploy data feed cache contract")

		mergeErr := fullCldEnvOutput.Environment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr, "failed to merge address book")
		fullCldEnvOutput.Environment.DataStore = dfOutput.DataStore.Seal()

		workflowName := "por-workflow-" + bcOutput.BlockchainOutput.ChainID + "-" + uuid.New().String()[0:4]

		dfConfigInput := &configureDataFeedsCacheInput{
			chainSelector:       bcOutput.ChainSelector,
			fullCldEnvironment:  fullCldEnvOutput.Environment,
			workflowName:        workflowName,
			feedID:              feedIDs[idx],
			sethClient:          bcOutput.SethClient,
			homeChainSethClient: homeChainSethClient, // Pass home chain's seth client for workflow owner
			blockchain:          bcOutput.BlockchainOutput,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		testLogger.Info().Msg("Proceeding to register PoR workflow...")

		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			homeChainSelector,
			keystone_changeset.WorkflowRegistry.String(),
		)
		require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", bcOutput.ChainID)

		dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			bcOutput.ChainSelector,
			df_changeset.DataFeedsCache.String(),
		)
		require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", bcOutput.ChainID)

		workflowConfigFilePath, configErr := createConfigFile(dataFeedsCacheAddress, workflowName, feedIDs[idx], priceProvider.URL(), corevm.GenerateWriteTargetName(bcOutput.ChainID))
		require.NoError(t, configErr, "failed to create workflow config file")

		compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
		require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)

		t.Cleanup(func() {
			wasmErr := os.Remove(compressedWorkflowWasmPath)
			if wasmErr != nil {
				framework.L.Warn().Msgf("failed to remove workflow wasm file %s: %s", compressedWorkflowWasmPath, wasmErr.Error())
			}
			configErr := os.Remove(workflowConfigFilePath)
			if configErr != nil {
				framework.L.Warn().Msgf("failed to remove workflow config file %s: %s", workflowConfigFilePath, configErr.Error())
			}
			deleteErr := creworkflow.DeleteWithContract(t.Context(), wrappedBlockchainOutputs[0].SethClient, workflowRegistryAddress, workflowName)
			if deleteErr != nil {
				framework.L.Warn().Msgf("failed to delete workflow %s: %s. Please delete it manually.", workflowName, deleteErr.Error())
			}
			debugPoRTest(t, testLogger, in, fullCldEnvOutput, wrappedBlockchainOutputs, feedIDs)
		})

		containerTargetDir := "/home/chainlink/workflows"
		workflowCopyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
		require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")

		configCopyErr := creworkflow.CopyWorkflowToDockerContainers(workflowConfigFilePath, "workflow-node", containerTargetDir)
		require.NoError(t, configCopyErr, "failed to copy workflow config to docker containers")

		registerErr := creworkflow.RegisterWithContract(
			t.Context(),
			wrappedBlockchainOutputs[0].SethClient, // crucial to use Seth Client connected to home chain (first chain in the set)
			workflowRegistryAddress,
			fullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
			workflowName,
			"file://"+compressedWorkflowWasmPath,
			ptr.Ptr("file://"+workflowConfigFilePath),
			nil,
			&containerTargetDir,
		)
		require.NoError(t, registerErr, "failed to register PoR workflow")
	}

	/*
		START THE VALIDATION PHASE
		Check whether each feed has been updated with the expected prices, which workflow fetches from the price provider
	*/
	eg := &errgroup.Group{}
	for idx, bcOutput := range wrappedBlockchainOutputs {
		eg.Go(func() error {
			feedID := feedIDs[idx]
			testLogger.Info().Msgf("Waiting for feed %s to update...", feedID)

			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
				fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				bcOutput.ChainSelector,
				df_changeset.DataFeedsCache.String(),
			)
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", bcOutput.ChainID)

			startTime := time.Now()

			// Handle Tron chains differently from EVM chains for reading contract data
			if bcOutput.BlockchainOutput.Family == "tron" {
				// For Tron chains, use the Tron chain client to read contract data
				tronChains := fullCldEnvOutput.Environment.BlockChains.TronChains()
				tronChain, exists := tronChains[bcOutput.ChainSelector]
				if !exists {
					return errors.Errorf("Tron chain %d not found in environment", bcOutput.ChainSelector)
				}

				// Convert Ethereum address back to Tron address for contract calls
				cacheAddr := address.EVMAddressToAddress(dataFeedsCacheAddresses)
				testLogger.Info().Msgf("Tron chain %d: Contract address conversion - EVM: %s -> Tron: %s", bcOutput.ChainSelector, dataFeedsCacheAddresses.Hex(), cacheAddr.String())

				assert.Eventually(t, func() bool {
					elapsed := time.Since(startTime).Round(time.Second)

					// First verify the contract exists by checking if it's deployed
					accountInfo, accountErr := tronChain.Client.GetAccount(cacheAddr)
					if accountErr != nil {
						testLogger.Error().Err(accountErr).Msgf("Tron chain %d: Failed to get account info for contract %s", bcOutput.ChainSelector, cacheAddr.String())
						return false
					}

					if accountInfo == nil || len(accountInfo.Address) == 0 {
						testLogger.Error().Msgf("Tron chain %d: Contract %s does not exist or is not deployed", bcOutput.ChainSelector, cacheAddr.String())
						return false
					}

					// Call getLatestAnswer on the Tron DataFeedsCache contract
					// This mirrors the EVM implementation: dataFeedsCacheInstance.GetLatestAnswer()
					testLogger.Info().Msgf("Tron chain %d: Calling getLatestAnswer for feed %s on contract %s", bcOutput.ChainSelector, feedID, cacheAddr.String())

					result, err := tronChain.Client.TriggerConstantContract(
						tronChain.Address,          // caller address
						cacheAddr,                  // contract address
						"getLatestAnswer(bytes16)", // function signature
						[]interface{}{"bytes16", [16]byte(common.Hex2Bytes(feedID))}, // parameters
					)
					if err != nil {
						testLogger.Error().Err(err).Msgf("FAILED to call getLatestAnswer on Tron chain %d", bcOutput.ChainSelector)
						return false
					}

					testLogger.Info().Msgf("Tron chain %d: Got result from contract call: %+v", bcOutput.ChainSelector, result)

					if len(result.ConstantResult) == 0 {
						testLogger.Error().Msgf("NO RESULT from getLatestAnswer on Tron chain %d", bcOutput.ChainSelector)
						return false
					}

					// Parse the result as a big integer (price)
					// The result is typically returned as hex-encoded bytes
					priceBytes := result.ConstantResult[0]
					if len(priceBytes) == 0 {
						testLogger.Error().Msgf("EMPTY price result from Tron chain %d", bcOutput.ChainSelector)
						return false
					}

					testLogger.Info().Msgf("Tron chain %d: Raw price bytes: %s", bcOutput.ChainSelector, priceBytes)

					// Convert hex string result to big.Int
					price := new(big.Int)
					if len(priceBytes) >= 2 && priceBytes[:2] == "0x" {
						price.SetString(priceBytes[2:], 16)
					} else {
						price.SetString(priceBytes, 16)
					}

					testLogger.Info().Msgf("Tron chain %d: Parsed price %s for feed %s", bcOutput.ChainSelector, price.String(), feedID)

					// if there are no more prices to be found, we can stop waiting
					return !priceProvider.NextPrice(feedID, price, elapsed)
				}, verificationTimeout, 10*time.Second, "feed %s did not update, timeout after: %s", feedID, verificationTimeout)

				expected := priceProvider.ExpectedPrices(feedID)
				actual := priceProvider.ActualPrices(feedID)

				if len(expected) != len(actual) {
					return errors.Errorf("expected %d prices, got %d", len(expected), len(actual))
				}

				for i := range expected {
					if expected[i].Cmp(actual[i]) != 0 {
						return errors.Errorf("expected price %d, got %d", expected[i], actual[i])
					}
				}

				testLogger.Info().Msgf("All prices were found in the feed %s", feedID)
			} else {
				// Handle EVM chains with existing logic
				dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, bcOutput.SethClient.Client)
				require.NoError(t, instanceErr, "failed to create data feeds cache instance")

				assert.Eventually(t, func() bool {
					elapsed := time.Since(startTime).Round(time.Second)
					price, err := dataFeedsCacheInstance.GetLatestAnswer(bcOutput.SethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(feedID)))
					require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

					// if there are no more prices to be found, we can stop waiting
					return !priceProvider.NextPrice(feedID, price, elapsed)
				}, verificationTimeout, 10*time.Second, "feed %s did not update, timeout after: %s", feedID, verificationTimeout)

				expected := priceProvider.ExpectedPrices(feedID)
				actual := priceProvider.ActualPrices(feedID)

				if len(expected) != len(actual) {
					return errors.Errorf("expected %d prices, got %d", len(expected), len(actual))
				}

				for i := range expected {
					if expected[i].Cmp(actual[i]) != 0 {
						return errors.Errorf("expected price %d, got %d", expected[i], actual[i])
					}
				}

				testLogger.Info().Msgf("All prices were found in the feed %s", feedID)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	testLogger.Info().Msgf("All prices were found for all feeds")
}

func executeVaultTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact) {
	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	fullCldEnvOutput, _, loadErr := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

	/*
		CREATE NEW VAULT SECRET
	*/
	framework.L.Info().Msg("Creating secret...")
	secretsRequest := jsonrpc.Request[vault.SecretsCreateRequest]{
		ID:      "request-id",
		Version: jsonrpc.JsonRpcVersion,
		Method:  vault.MethodSecretsCreate,
		Params: &vault.SecretsCreateRequest{
			ID:    "test-secret",
			Value: "test-secret-value",
		},
	}
	requestBody, err := json.Marshal(secretsRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	framework.L.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")

	framework.L.Info().Msgf("Gateway URL: %s", gatewayURL.String())
	framework.L.Info().Msg("Executing request...")
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL.String(), bytes.NewBuffer(requestBody))
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/jsonrpc")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")
	framework.L.Debug().Msgf("Response Body: %s", string(body))

	framework.L.Info().Msg("Checking response status...")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking response structure...")
	var response jsonrpc.Response[vault.SecretsCreateResponse]
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "failed to unmarshal response")

	require.Equal(t, jsonrpc.JsonRpcVersion, response.Version)
	require.NotEmpty(t, response.ID)
	require.NoError(t, err, "failed to unmarshal response")
	require.True(t, response.Result.Success)
	require.Equal(t, "test-secret", response.Result.SecretID)
	require.Empty(t, response.Result.ErrorMessage)

	framework.L.Info().Msg("Secret created successfully")
}

const (
	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""
)

func createEnvironmentIfNotExists(stateFile, environmentDir, topology string) error {
	split := strings.Split(stateFile, ",")
	if _, err := os.Stat(split[0]); os.IsNotExist(err) {
		ctfConfigs := os.Getenv("CTF_CONFIGS")
		defer func() {
			setErr := os.Setenv("CTF_CONFIGS", ctfConfigs)
			if setErr != nil {
				framework.L.Error().Err(setErr).Msg("failed to set CTF_CONFIGS env var")
			}
		}()

		// unset the CTF_CONFIGS env var to avoid using the cached environment
		setErr := os.Setenv("CTF_CONFIGS", "")
		if setErr != nil {
			return errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
		}

		cmd := exec.Command("go", "run", ".", "env", "start", "--topology", topology)
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start environment")
		}
	}

	return nil
}

type configureDataFeedsCacheInput struct {
	chainSelector       uint64
	fullCldEnvironment  *cldf.Environment
	workflowName        string
	feedID              string
	sethClient          *seth.Client
	homeChainSethClient *seth.Client // Home chain's seth client for workflow owner
	blockchain          *blockchain.Output
}

func configureDataFeedsCacheContract(testLogger zerolog.Logger, input *configureDataFeedsCacheInput) error {
	forwarderAddress, forwarderErr := crecontracts.FindAddressesForChain(input.fullCldEnvironment.ExistingAddresses, input.chainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck // won't migrate now
	if forwarderErr != nil {
		return errors.Wrapf(forwarderErr, "failed to find forwarder address for chain %d", input.chainSelector)
	}

	dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(input.fullCldEnvironment.ExistingAddresses, input.chainSelector, df_changeset.DataFeedsCache.String()) //nolint:staticcheck // won't migrate now
	if dataFeedsCacheErr != nil {
		return errors.Wrapf(dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", input.chainSelector)
	}

	// Handle Tron chains differently from EVM chains
	if input.blockchain.Family == "tron" {
		// Set feed admin using Tron changeset
		triggerOptions := cldf_tron.DefaultTriggerOptions()
		triggerOptions.FeeLimit = 1_000_000_000

		// Convert Ethereum address to Tron address for cache address
		cacheAddr := address.EVMAddressToAddress(dataFeedsCacheAddress)

		// Convert forwarder address to Tron format
		forwarderTronAddr := address.EVMAddressToAddress(forwarderAddress)

		// Get the Tron chain for the deployer address
		tronChains := input.fullCldEnvironment.BlockChains.TronChains()
		tronChain, exists := tronChains[input.chainSelector]
		if !exists {
			return errors.Errorf("Tron chain %d not found in environment", input.chainSelector)
		}

		// Set the deployer as admin (equivalent to sethClient.MustGetRootKeyAddress() for EVM)
		// This matches the EVM pattern where only the deployer is the admin
		setDeployerAdminConfig := df_changeset_types.SetFeedAdminTronConfig{
			ChainSelector:  input.chainSelector,
			CacheAddress:   cacheAddr,
			AdminAddress:   tronChain.Address, // Deployer address (equivalent to MustGetRootKeyAddress)
			IsAdmin:        true,
			TriggerOptions: triggerOptions,
		}

		_, setDeployerAdminErr := changeset.RunChangeset(tron_df_changeset.SetFeedAdminChangeset, *input.fullCldEnvironment, setDeployerAdminConfig)
		if setDeployerAdminErr != nil {
			return errors.Wrap(setDeployerAdminErr, "failed to set deployer as admin for Tron chain")
		}

		// Set feed config using Tron changeset - equivalent to EVM ConfigureDataFeedsCache

		// Convert workflow name to [10]byte hash (same as EVM implementation)
		workflowNameBytes := df_changeset.HashedWorkflowName(input.workflowName)

		// Use home chain's deployer address as AllowedWorkflowOwner (matches EVM pattern using MustGetRootKeyAddress)
		// This ensures the workflow owner is consistent across all chains (EVM and Tron)
		homeChainWorkflowOwner := address.EVMAddressToAddress(input.homeChainSethClient.MustGetRootKeyAddress())
		workflowMetadata := []df_changeset_types.DataFeedsCacheTronWorkflowMetadata{
			{
				AllowedSender:        forwarderTronAddr,
				AllowedWorkflowOwner: homeChainWorkflowOwner, // Use home chain's deployer address for consistency
				AllowedWorkflowName:  workflowNameBytes,
			},
		}

		// Truncate feed ID to 32 characters (16 bytes) for Tron
		feedIDTruncated := input.feedID
		// Remove 0x prefix if present
		feedIDTruncated = strings.TrimPrefix(feedIDTruncated, "0x")
		// Truncate to 32 characters (16 bytes)
		if len(feedIDTruncated) > 32 {
			feedIDTruncated = feedIDTruncated[:32]
		}

		setFeedConfigConfig := df_changeset_types.SetFeedDecimalTronConfig{
			ChainSelector:    input.chainSelector,
			CacheAddress:     cacheAddr,
			DataIDs:          []string{feedIDTruncated},
			Descriptions:     []string{"PoR test feed"},
			WorkflowMetadata: workflowMetadata,
			TriggerOptions:   triggerOptions,
		}

		_, setConfigErr := changeset.RunChangeset(tron_df_changeset.SetFeedConfigChangeset, *input.fullCldEnvironment, setFeedConfigConfig)
		if setConfigErr != nil {
			return errors.Wrap(setConfigErr, "failed to set feed config for Tron chain")
		}

		testLogger.Info().Msgf("Successfully configured Tron data feeds cache for chain %d", input.chainSelector)
		return nil
	} else {
		// Handle EVM chains with existing logic
		configInput := &cre.ConfigureDataFeedsCacheInput{
			CldEnv:                input.fullCldEnvironment,
			ChainSelector:         input.chainSelector,
			FeedIDs:               []string{input.feedID},
			Descriptions:          []string{"PoR test feed"},
			DataFeedsCacheAddress: dataFeedsCacheAddress,
			AdminAddress:          input.sethClient.MustGetRootKeyAddress(),
			AllowedSenders:        []common.Address{forwarderAddress},
			AllowedWorkflowNames:  []string{input.workflowName},
			AllowedWorkflowOwners: []common.Address{input.sethClient.MustGetRootKeyAddress()},
		}

		_, configErr := crecontracts.ConfigureDataFeedsCache(testLogger, configInput)
		return configErr
	}
}

func logTestInfo(l zerolog.Logger, feedID, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func createConfigFile(feedsConsumerAddress common.Address, workflowName, feedID, dataURL, writeTargetName string) (string, error) {
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return "", errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID

	workflowConfig := portypes.WorkflowConfig{
		ComputeConfig: portypes.ComputeConfig{
			FeedID:                feedIDToUse,
			URL:                   dataURL,
			DataFeedsCacheAddress: feedsConsumerAddress.Hex(),
			WriteTargetName:       writeTargetName,
		},
	}

	configMarshalled, err := yaml.Marshal(workflowConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal workflow config")
	}
	outputFile := workflowName + "_config.yaml"

	// remove the file if it already exists
	_, statErr := os.Stat(outputFile)
	if statErr == nil {
		if err := os.Remove(outputFile); err != nil {
			return "", errors.Wrap(err, "failed to remove existing output file")
		}
	}

	if err := os.WriteFile(outputFile, configMarshalled, 0644); err != nil { //nolint:gosec // G306: we want it to be readable by everyone
		return "", errors.Wrap(err, "failed to write output file")
	}

	outputFileAbsPath, outputFileAbsPathErr := filepath.Abs(outputFile)
	if outputFileAbsPathErr != nil {
		return "", errors.Wrap(outputFileAbsPathErr, "failed to get absolute path of the config file")
	}

	return outputFileAbsPath, nil
}

func debugPoRTest(t *testing.T, testLogger zerolog.Logger, in *envconfig.Config, env *cre.FullCLDEnvironmentOutput, wrappedBlockchainOutputs []*cre.WrappedBlockchainOutput, feedIDs []string) {
	if t.Failed() {
		counter := 0
		for idx, feedID := range feedIDs {
			chainSelector := wrappedBlockchainOutputs[idx].ChainSelector
			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
				env.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				chainSelector,
				df_changeset.DataFeedsCache.String(),
			)
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

			forwarderAddresses, forwarderErr := crecontracts.FindAddressesForChain(
				env.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				chainSelector,
				keystone_changeset.KeystoneForwarder.String(),
			)
			require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", chainSelector)

			logTestInfo(testLogger, feedID, dataFeedsCacheAddresses.Hex(), forwarderAddresses.Hex())
			counter++
			// log scanning is not supported for CRIB
			if in.Infra.Type == infra.CRIB {
				return
			}

			_, saveErr := framework.SaveContainerLogs(os.TempDir())
			if saveErr != nil {
				testLogger.Error().Err(saveErr).Msg("failed to save container logs")
				return
			}

			debugDons := make([]*cre.DebugDon, 0, len(env.DonTopology.DonsWithMetadata))
			for i, donWithMetadata := range env.DonTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range in.NodeSets[i].Out.CLNodes {
					containerNames = append(containerNames, output.Node.ContainerName)
				}
				debugDons = append(debugDons, &cre.DebugDon{
					NodesMetadata:  donWithMetadata.NodesMetadata,
					Flags:          donWithMetadata.Flags,
					ContainerNames: containerNames,
				})
			}

			debugInput := cre.DebugInput{
				DebugDons:        debugDons,
				BlockchainOutput: wrappedBlockchainOutputs[idx].BlockchainOutput,
				InfraInput:       in.Infra,
			}
			credebug.PrintTestDebug(t.Context(), t.Name(), testLogger, debugInput)
		}
	}
}

func setConfigurationIfMissing(configName, topology string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	if os.Getenv("CRE_TOPOLOGY") == "" {
		err := os.Setenv("CRE_TOPOLOGY", topology)
		if err != nil {
			return errors.Wrap(err, "failed to set CRE_TOPOLOGY env var")
		}
	}

	if os.Getenv("ENV_ARTIFACT_PATH") == "" {
		err := os.Setenv("ENV_ARTIFACT_PATH", "../../../..//core/scripts/cre/environment/env_artifact/env_artifact.json")
		if err != nil {
			return errors.Wrap(err, "failed to set ENV_ARTIFACT_PATH env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}
