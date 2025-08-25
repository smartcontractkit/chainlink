package cre

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/vault"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crevault "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/vault"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
)

type TestEnvironment struct {
	Config                   *envconfig.Config
	EnvArtifact              environment.EnvArtifact
	Logger                   zerolog.Logger
	FullCldEnvOutput         *cre.FullCLDEnvironmentOutput
	WrappedBlockchainOutputs []*cre.WrappedBlockchainOutput
}

type WorkflowTestConfig struct {
	WorkflowName     string
	WorkflowLocation string
	FeedIDs          []string
	Timeout          time.Duration
}

/*
To execute on local start the local CRE first with following command:
# inside core/scripts/cre/environment directory
go run . env start
*/
func Test_CRE_Workflow_Don(t *testing.T) {
	testEnv := setupTestEnvironment(t)

	// currently we can't run these tests in parallel, because each test rebuilds environment structs and that includes
	// logging into CL node with GraphQL API, which allows only 1 session per user at a time.
	t.Run("cron-based PoR workflow", func(t *testing.T) {
		executePoRTest(t, testEnv)
	})

	t.Run("vault DON test", func(t *testing.T) {
		executeVaultTest(t, testEnv)
	})

	t.Run("http trigger and action test", func(t *testing.T) {
		executeHTTPTriggerActionTest(t, testEnv)
	})
}

// WorkflowRegistrationConfig holds configuration for workflow registration
type WorkflowRegistrationConfig struct {
	WorkflowName         string
	WorkflowLocation     string
	ConfigFilePath       string
	CompressedWasmPath   string
	WorkflowRegistryAddr common.Address
	DonID                uint64
	ContainerTargetDir   string
}

// WorkflowCleanupConfig holds configuration for cleanup operations
type WorkflowCleanupConfig struct {
	WasmPath          string
	ConfigPath        string
	WorkflowName      string
	RegistryAddress   common.Address
	SethClient        *seth.Client
	Logger            zerolog.Logger
	TestEnv           *TestEnvironment
	FullCldEnvOutput  *cre.FullCLDEnvironmentOutput
	BlockchainOutputs []*cre.WrappedBlockchainOutput
	FeedIDs           []string
}

// setupTestEnvironment initializes the common test environment
func setupTestEnvironment(t *testing.T) *TestEnvironment {
	confErr := setConfigurationIfMissing(DefaultConfigPath, DefaultTopology)
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
	require.NoError(t, err, "failed to read artifact file")
	err = json.Unmarshal(artFile, &envArtifact)
	require.NoError(t, err, "failed to unmarshal artifact file")

	fullCldEnvOutput, wrappedBlockchainOutputs, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, err, "failed to load environment")

	return &TestEnvironment{
		Config:                   in,
		EnvArtifact:              envArtifact,
		Logger:                   framework.L,
		FullCldEnvOutput:         fullCldEnvOutput,
		WrappedBlockchainOutputs: wrappedBlockchainOutputs,
	}
}

// copyWorkflowFilesToContainers copies workflow files to Docker containers
func copyWorkflowFilesToContainers(t *testing.T, wasmPath, configPath, containerTargetDir string) {
	workflowCopyErr := creworkflow.CopyWorkflowToDockerContainers(wasmPath, WorkflowNodePrefix, containerTargetDir)
	require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")

	configCopyErr := creworkflow.CopyWorkflowToDockerContainers(configPath, WorkflowNodePrefix, containerTargetDir)
	require.NoError(t, configCopyErr, "failed to copy workflow config to docker containers")
}

// registerWorkflow registers a workflow with the contract
func registerWorkflow(ctx context.Context, t *testing.T, config *WorkflowRegistrationConfig, sethClient *seth.Client) {
	registerErr := creworkflow.RegisterWithContract(
		ctx,
		sethClient,
		config.WorkflowRegistryAddr,
		config.DonID,
		config.WorkflowName,
		"file://"+config.CompressedWasmPath,
		ptr.Ptr("file://"+config.ConfigFilePath),
		nil,
		&config.ContainerTargetDir,
	)
	require.NoError(t, registerErr, "failed to register workflow '%s'", config.WorkflowName)
}

// setupWorkflowCleanup sets up cleanup for workflow resources
func setupWorkflowCleanup(t *testing.T, config *WorkflowCleanupConfig) {
	t.Cleanup(func() {
		if wasmErr := os.Remove(config.WasmPath); wasmErr != nil {
			config.Logger.Warn().Msgf("failed to remove workflow wasm file %s: %s", config.WasmPath, wasmErr.Error())
		}
		if configErr := os.Remove(config.ConfigPath); configErr != nil {
			config.Logger.Warn().Msgf("failed to remove workflow config file %s: %s", config.ConfigPath, configErr.Error())
		}
		if deleteErr := creworkflow.DeleteWithContract(t.Context(), config.SethClient, config.RegistryAddress, config.WorkflowName); deleteErr != nil {
			config.Logger.Warn().Msgf("failed to delete workflow %s: %s. Please delete it manually.", config.WorkflowName, deleteErr.Error())
		}
	})

	t.Run("DON Time test", func(t *testing.T) {
		// TODO: Implement smoke test - https://smartcontract-it.atlassian.net/browse/CAPPL-1028
		t.Skip()
	})
}

type HTTPTestConfig struct {
	WorkflowName     string
	MockServer       *httptest.Server
	Recorder         *MockServerRecorder
	SigningKey       *ecdsa.PrivateKey
	ConfigPath       string
	WorkflowLocation string
}

// setupHTTPWorkflowTest sets up the HTTP workflow test infrastructure
func setupHTTPWorkflowTest(t *testing.T, testEnv *TestEnvironment) *HTTPTestConfig {
	mockServer, recorder := startMockHTTPServerOnPort(t, testEnv.Config.Fake.Port)
	workflowName := "http-trigger-action-test-" + uuid.New().String()[0:8]
	configPath, signingKey := createTestWorkflowConfig(t, workflowName, mockServer.URL)

	return &HTTPTestConfig{
		WorkflowName:     workflowName,
		MockServer:       mockServer,
		Recorder:         recorder,
		SigningKey:       signingKey,
		ConfigPath:       configPath,
		WorkflowLocation: HTTPWorkflowLocation,
	}
}

// executeHTTPTriggerRequest executes an HTTP trigger request and waits for successful response
func executeHTTPTriggerRequest(t *testing.T, testEnv *TestEnvironment, gatewayURL *url.URL, httpConfig *HTTPTestConfig, workflowOwnerAddress string) {
	var finalResponse jsonrpc.Response[json.RawMessage]
	var triggerRequest jsonrpc.Request[json.RawMessage]

	require.Eventually(t, func() bool {
		triggerRequest = createHTTPTriggerRequestWithKey(t, httpConfig.WorkflowName, workflowOwnerAddress, httpConfig.SigningKey)
		triggerRequestBody, err := json.Marshal(triggerRequest)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to marshal trigger request: %v", err)
			return false
		}

		testEnv.Logger.Info().Msgf("Gateway URL: %s", gatewayURL.String())
		testEnv.Logger.Info().Msg("Executing HTTP trigger request with retries until workflow is loaded...")

		req, err := http.NewRequestWithContext(t.Context(), "POST", gatewayURL.String(), bytes.NewBuffer(triggerRequestBody))
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to create request: %v", err)
			return false
		}
		req.Header.Set("Content-Type", "application/jsonrpc")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to execute request: %v", err)
			return false
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to read response body: %v", err)
			return false
		}

		testEnv.Logger.Info().Msgf("HTTP trigger response (status %d): %s", resp.StatusCode, string(body))

		if resp.StatusCode != http.StatusOK {
			testEnv.Logger.Warn().Msgf("Gateway returned status %d, retrying...", resp.StatusCode)
			return false
		}

		err = json.Unmarshal(body, &finalResponse)
		if err != nil {
			testEnv.Logger.Warn().Msgf("Failed to unmarshal response: %v", err)
			return false
		}

		if finalResponse.Error != nil {
			testEnv.Logger.Warn().Msgf("JSON-RPC error in response: %v", finalResponse.Error)
			return false
		}

		testEnv.Logger.Info().Msg("Successfully received 200 OK response from gateway")
		return true
	}, tests.WaitTimeout(t), RetryInterval, "gateway should respond with 200 OK and valid response once workflow is loaded")

	require.Equal(t, jsonrpc.JsonRpcVersion, finalResponse.Version, "expected JSON-RPC version %s, got %s", jsonrpc.JsonRpcVersion, finalResponse.Version)
	require.Equal(t, triggerRequest.ID, finalResponse.ID, "expected response ID %s, got %s", triggerRequest.ID, finalResponse.ID)
	require.Nil(t, finalResponse.Error, "unexpected error in response: %v", finalResponse.Error)
}

// validateHTTPWorkflowRequest validates that the workflow made the expected HTTP request
func validateHTTPWorkflowRequest(t *testing.T, testEnv *TestEnvironment, recorder *MockServerRecorder) {
	require.Eventually(t, func() bool {
		return len(recorder.GetRequests()) > 0
	}, tests.WaitTimeout(t), RetryInterval, "workflow should have made at least one HTTP request to mock server")

	recordedRequest := recorder.GetRequests()[0]
	testEnv.Logger.Info().Msgf("Recorded request: %+v", recordedRequest)

	require.Equal(t, "POST", recordedRequest.Method, "expected POST method")
	require.Equal(t, "/orders", recordedRequest.URL, "expected /orders endpoint")
	require.Equal(t, "application/json", recordedRequest.Headers["Content-Type"], "expected JSON content type")

	var workflowRequestBody map[string]interface{}
	err := json.Unmarshal([]byte(recordedRequest.Body), &workflowRequestBody)
	require.NoError(t, err, "request body should be valid JSON")

	require.Equal(t, "test-customer", workflowRequestBody["customer"], "expected customer field")
	require.Equal(t, "large", workflowRequestBody["size"], "expected size field")
	require.Contains(t, workflowRequestBody, "toppings", "expected toppings field")
}

// validatePoRPrices validates that all feeds receive the expected prices from the price provider
func validatePoRPrices(t *testing.T, testEnv *TestEnvironment, priceProvider PriceProvider, config *WorkflowTestConfig) {
	eg := &errgroup.Group{}

	for idx, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}
		eg.Go(func() error {
			feedID := config.FeedIDs[idx]
			testEnv.Logger.Info().Msgf("Waiting for feed %s to update...", feedID)

			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
				testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
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
			require.Eventually(t, func() bool {
				elapsed := time.Since(startTime).Round(time.Second)
				price, err := dataFeedsCacheInstance.GetLatestAnswer(bcOutput.SethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(feedID)))
				if err != nil {
					testEnv.Logger.Error().Err(err).Msg("failed to get price from Data Feeds Cache contract")
					return false
				}

				// if there are no more prices to be found, we can stop waiting
				return !priceProvider.NextPrice(feedID, price, elapsed)
			}, config.Timeout, ValidationInterval, "feed %s did not update, timeout after: %s", feedID, config.Timeout)

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

func executePoRTest(t *testing.T, testEnv *TestEnvironment) {
	feedIDs := []string{"018e16c39e000320000000000000000000000000000000000000000000000000", "018e16c38e000320000000000000000000000000000000000000000000000000"}
	priceProvider, err := NewFakePriceProvider(testEnv.Logger, testEnv.Config.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, err, "failed to create fake price provider")
	defer priceProvider.Close(t.Context())

	config := &WorkflowTestConfig{
		WorkflowName:     "por-workflow",
		WorkflowLocation: PoRWorkflowLocation,
		FeedIDs:          feedIDs,
		Timeout:          DefaultVerificationTimeout,
	}

	executePoRWorkflowTest(t, testEnv, priceProvider, config)
}

// executePoRWorkflowTest handles the main PoR workflow test logic
func executePoRWorkflowTest(t *testing.T, testEnv *TestEnvironment, priceProvider PriceProvider, config *WorkflowTestConfig) {
	homeChainSelector := testEnv.WrappedBlockchainOutputs[0].ChainSelector
	writeableChains := []uint64{}
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		for _, donMetadata := range testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata {
			if flags.RequiresForwarderContract(donMetadata.Flags, bcOutput.ChainID) {
				if !slices.Contains(writeableChains, bcOutput.ChainID) {
					writeableChains = append(writeableChains, bcOutput.ChainID)
				}
			}
		}
	}
	require.Len(t, config.FeedIDs, len(writeableChains), "number of writeable chains must match number of feed IDs (check what chains 'evm' and 'write-evm' capabilities are enabled for)")

	/*
		DEPLOY DATA FEEDS CACHE CONTRACTS ON ALL CHAINS (except read-only ones)
		Workflow will write price data to the data feeds cache contract

		REGISTER ONE WORKFLOW PER CHAIN (except read-only ones)
	*/
	for idx, bcOutput := range testEnv.WrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}
		// deploy data feeds cache contract only on chains that require a forwarder contract. It's required for the PoR workflow to work and we treat it as a proxy
		// for deciding whether need to deploy the data feeds cache contract.
		hasForwarderContract := false
		for _, donMetadata := range testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata {
			if flags.RequiresForwarderContract(donMetadata.Flags, bcOutput.ChainID) {
				hasForwarderContract = true
				break
			}
		}

		if !hasForwarderContract {
			continue
		}

		deployConfig := df_changeset_types.DeployConfig{
			ChainsToDeploy: []uint64{bcOutput.ChainSelector},
			Labels:         []string{"data-feeds"}, // label required by the changeset
		}

		dfOutput, dfErr := changeset.RunChangeset(df_changeset.DeployCacheChangeset, *testEnv.FullCldEnvOutput.Environment, deployConfig)
		require.NoError(t, dfErr, "failed to deploy data feed cache contract")

		mergeErr := testEnv.FullCldEnvOutput.Environment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr, "failed to merge address book")
		testEnv.FullCldEnvOutput.Environment.DataStore = dfOutput.DataStore.Seal()

		workflowName := "por-workflow-" + bcOutput.BlockchainOutput.ChainID + "-" + uuid.New().String()[0:4]

		dfConfigInput := &configureDataFeedsCacheInput{
			chainSelector:      bcOutput.ChainSelector,
			fullCldEnvironment: testEnv.FullCldEnvOutput.Environment,
			workflowName:       workflowName,
			feedID:             config.FeedIDs[idx],
			sethClient:         bcOutput.SethClient,
			blockchain:         bcOutput.BlockchainOutput,
		}
		dfConfigErr := configureDataFeedsCacheContract(testEnv.Logger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		testEnv.Logger.Info().Msg("Proceeding to register PoR workflow...")

		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
			testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			homeChainSelector,
			keystone_changeset.WorkflowRegistry.String(),
		)
		require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", bcOutput.ChainID)

		dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
			testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			bcOutput.ChainSelector,
			df_changeset.DataFeedsCache.String(),
		)
		require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", bcOutput.ChainID)

		workflowConfigFilePath, configErr := createConfigFile(dataFeedsCacheAddress, workflowName, config.FeedIDs[idx], priceProvider.URL(), corevm.GenerateWriteTargetName(bcOutput.ChainID))
		require.NoError(t, configErr, "failed to create workflow config file")

		compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(config.WorkflowLocation, workflowName)

		require.NoError(t, compileErr, "failed to compile workflow")

		cleanupConfig := &WorkflowCleanupConfig{
			WasmPath:          compressedWorkflowWasmPath,
			ConfigPath:        workflowConfigFilePath,
			WorkflowName:      workflowName,
			RegistryAddress:   workflowRegistryAddress,
			SethClient:        testEnv.WrappedBlockchainOutputs[0].SethClient,
			Logger:            testEnv.Logger,
			TestEnv:           testEnv,
			FullCldEnvOutput:  testEnv.FullCldEnvOutput,
			BlockchainOutputs: testEnv.WrappedBlockchainOutputs,
			FeedIDs:           config.FeedIDs,
		}
		setupWorkflowCleanup(t, cleanupConfig)
		debugPoRTest(t, testEnv.Logger, testEnv.Config, testEnv.FullCldEnvOutput, testEnv.WrappedBlockchainOutputs, config.FeedIDs)

		copyWorkflowFilesToContainers(t, compressedWorkflowWasmPath, workflowConfigFilePath, ContainerTargetDir)

		regConfig := &WorkflowRegistrationConfig{
			WorkflowName:         workflowName,
			WorkflowLocation:     config.WorkflowLocation,
			ConfigFilePath:       workflowConfigFilePath,
			CompressedWasmPath:   compressedWorkflowWasmPath,
			WorkflowRegistryAddr: workflowRegistryAddress,
			DonID:                testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
			ContainerTargetDir:   ContainerTargetDir,
		}
		registerWorkflow(t.Context(), t, regConfig, testEnv.WrappedBlockchainOutputs[0].SethClient)
	}

	validatePoRPrices(t, testEnv, priceProvider, config)
}

func executeVaultTest(t *testing.T, testEnv *TestEnvironment) {
	// Skip till we figure out and fix the issues with environment startup on this test
	t.Skip()
	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	framework.L.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")

	framework.L.Info().Msgf("Gateway URL: %s", gatewayURL.String())

	framework.L.Info().Msgf("Sleeping 2 minutes to allow the Vault DON to start...")
	time.Sleep(2 * time.Minute)
	framework.L.Info().Msgf("Sleep over. Executing test now...")

	secretID := strconv.Itoa(rand.Intn(10000)) // generate a random secret ID for testing
	owner := "Owner1"
	secretValue := "Secret Value to be stored"

	executeVaultSecretsCreateTest(t, secretValue, secretID, owner, gatewayURL.String())

	framework.L.Info().Msg("------------------------------------------------------")
	framework.L.Info().Msg("------------------------------------------------------")
	framework.L.Info().Msg("------------------------------------------------------")
	framework.L.Info().Msg("------------------------------------------------------")
	framework.L.Info().Msg("------------------------------------------------------")

	executeVaultSecretsGetTest(t, secretValue, secretID, owner, gatewayURL.String())
}

func executeVaultSecretsCreateTest(t *testing.T, secretValue, secretID, owner, gatewayURL string) {
	framework.L.Info().Msg("Creating secret...")
	uniqueRequestID := uuid.New().String()

	secretsCreateRequest := jsonrpc.Request[vault.SecretsCreateRequest]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  vault.MethodSecretsCreate,
		Params: &vault.SecretsCreateRequest{
			ID:    secretID,
			Value: encryptSecret(t, secretValue),
			Owner: owner,
		},
		ID: uniqueRequestID,
	}
	requestBody, err := json.Marshal(secretsCreateRequest)
	require.NoError(t, err, "failed to marshal secrets request")

	framework.L.Info().Msgf("Request Body: %s", string(requestBody))
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read createResponse body")
	framework.L.Info().Msgf("Response Body: %s", string(body))

	framework.L.Info().Msg("Checking createResponse status...")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking createResponse structure...")
	var createResponse jsonrpc.Response[json.RawMessage]
	err = json.Unmarshal(body, &createResponse)
	require.NoError(t, err, "failed to unmarshal createResponse")
	framework.L.Info().Msgf("createResponse Body: %v", createResponse)
	if createResponse.Error != nil {
		require.Empty(t, createResponse.Error.Error())
	}
	result := *createResponse.Result
	framework.L.Info().Msgf("SecretsCreateResponse: %s", string(result))

	require.Equal(t, jsonrpc.JsonRpcVersion, createResponse.Version)

	framework.L.Info().Msg("Secret created successfully")
}

func executeVaultSecretsGetTest(t *testing.T, secretValue, secretID, owner, gatewayURL string) {
	uniqueRequestID := uuid.New().String()
	framework.L.Info().Msg("Getting secret...")
	secretsGetRequest := jsonrpc.Request[vault.SecretsGetRequest]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  vault.MethodSecretsGet,
		Params: &vault.SecretsGetRequest{
			ID:    secretID,
			Owner: owner,
		},
		ID: uniqueRequestID,
	}
	requestBody, err := json.Marshal(secretsGetRequest)
	require.NoError(t, err, "failed to marshal secrets request")
	framework.L.Info().Msgf("Get Request Body: %s", string(requestBody))
	framework.L.Info().Msg("Executing Get request...")
	req, err := http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "failed to create request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "failed to execute request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read getResponse body")
	framework.L.Info().Msgf("getResponse Body: %s", string(body))

	framework.L.Info().Msg("Checking getResponse status...")
	require.Equal(t, http.StatusOK, resp.StatusCode, "Gateway endpoint should respond with 200 OK")

	framework.L.Info().Msg("Checking getResponse structure...")
	var getResponse jsonrpc.Response[vault.SecretsGetResponse]
	err = json.Unmarshal(body, &getResponse)
	require.NoError(t, err, "failed to unmarshal getResponse")
	framework.L.Info().Msgf("getResponse Body: %v", getResponse)
	if getResponse.Error != nil {
		require.Empty(t, getResponse.Error.Error())
	}

	result := getResponse.Result
	framework.L.Info().Msgf("SecretsGetResponse: %s", result)

	require.Equal(t, jsonrpc.JsonRpcVersion, getResponse.Version)

	require.Empty(t, result.Error)
	require.Equal(t, secretID, result.SecretID.Key)
	require.Equal(t, owner, result.SecretID.Owner)

	framework.L.Info().Msg("Secret get successful")
}

func encryptSecret(t *testing.T, secret string) string {
	masterPublicKey := tdh2easy.PublicKey{}
	masterPublicKeyBytes, err := hex.DecodeString(crevault.MasterPublicKeyStr)
	require.NoError(t, err)
	err = masterPublicKey.Unmarshal(masterPublicKeyBytes)
	require.NoError(t, err)
	cipher, err := tdh2easy.Encrypt(&masterPublicKey, []byte(secret))
	require.NoError(t, err)
	cipherBytes, err := cipher.Marshal()
	require.NoError(t, err)
	return hex.EncodeToString(cipherBytes)
}

func executeHTTPTriggerActionTest(t *testing.T, testEnv *TestEnvironment) {
	homeChainSelector := testEnv.WrappedBlockchainOutputs[0].ChainSelector
	testEnv.Logger.Info().Msg("Starting HTTP trigger and action test...")

	httpConfig := setupHTTPWorkflowTest(t, testEnv)
	defer httpConfig.MockServer.Close()

	compressedWorkflowWasmPath, err := creworkflow.CompileWorkflow(httpConfig.WorkflowLocation, httpConfig.WorkflowName)
	require.NoError(t, err, "failed to compile workflow")

	workflowRegistryAddress, err := crecontracts.FindAddressesForChain(
		testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		homeChainSelector,
		keystone_changeset.WorkflowRegistry.String(),
	)
	require.NoError(t, err, "failed to find workflow registry address for chain %d", homeChainSelector)

	cleanupConfig := &WorkflowCleanupConfig{
		WasmPath:          compressedWorkflowWasmPath,
		ConfigPath:        httpConfig.ConfigPath,
		WorkflowName:      httpConfig.WorkflowName,
		RegistryAddress:   workflowRegistryAddress,
		SethClient:        testEnv.WrappedBlockchainOutputs[0].SethClient,
		Logger:            testEnv.Logger,
		TestEnv:           testEnv,
		FullCldEnvOutput:  testEnv.FullCldEnvOutput,
		BlockchainOutputs: testEnv.WrappedBlockchainOutputs,
		FeedIDs:           []string{}, // No feed IDs for HTTP test
	}
	setupWorkflowCleanup(t, cleanupConfig)

	copyWorkflowFilesToContainers(t, compressedWorkflowWasmPath, httpConfig.ConfigPath, ContainerTargetDir)

	regConfig := &WorkflowRegistrationConfig{
		WorkflowName:         httpConfig.WorkflowName,
		WorkflowLocation:     httpConfig.WorkflowLocation,
		ConfigFilePath:       httpConfig.ConfigPath,
		CompressedWasmPath:   compressedWorkflowWasmPath,
		WorkflowRegistryAddr: workflowRegistryAddress,
		DonID:                testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
		ContainerTargetDir:   ContainerTargetDir,
	}
	registerWorkflow(t.Context(), t, regConfig, testEnv.WrappedBlockchainOutputs[0].SethClient)

	testEnv.Logger.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + testEnv.FullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
	require.NoError(t, err, "failed to parse gateway URL")

	workflowOwner, err := crypto.HexToECDSA(testEnv.WrappedBlockchainOutputs[0].DeployerPrivateKey)
	require.NoError(t, err, "failed to convert private key to ECDSA")
	workflowOwnerAddress := strings.ToLower(crypto.PubkeyToAddress(workflowOwner.PublicKey).Hex())

	testEnv.Logger.Info().Msgf("Workflow owner address: %s", workflowOwnerAddress)
	testEnv.Logger.Info().Msgf("Workflow name: %s", httpConfig.WorkflowName)

	executeHTTPTriggerRequest(t, testEnv, gatewayURL, httpConfig, workflowOwnerAddress)

	validateHTTPWorkflowRequest(t, testEnv, httpConfig.Recorder)

	testEnv.Logger.Info().Msg("HTTP trigger and action test completed successfully")
}

type MockServerRecorder struct {
	mu       sync.Mutex
	requests []RecordedRequest
}

type RecordedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (r *MockServerRecorder) RecordRequest(method, url string, headers http.Header, body []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	headerMap := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			headerMap[key] = values[0]
		}
	}

	r.requests = append(r.requests, RecordedRequest{
		Method:  method,
		URL:     url,
		Headers: headerMap,
		Body:    string(body),
	})
}

func (r *MockServerRecorder) GetRequests() []RecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()

	requests := make([]RecordedRequest, len(r.requests))
	copy(requests, r.requests)
	return requests
}

func (r *MockServerRecorder) GetRequestCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.requests)
}

func startMockHTTPServerOnPort(t *testing.T, port int) (*httptest.Server, *MockServerRecorder) {
	recorder := &MockServerRecorder{}
	mux := http.NewServeMux()

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		recorder.RecordRequest(r.Method, r.URL.String(), r.Header, body)

		framework.L.Info().Msgf("Mock server received order request: %s", string(body))

		response := map[string]interface{}{
			"orderId": "test-order-" + uuid.New().String()[0:8],
			"status":  "success",
			"message": "Order processed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	testServer := httptest.NewUnstartedServer(mux)
	testServer.Listener.Close()

	var err error
	testServer.Listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	require.NoError(t, err, "failed to listen on port %d", port)

	testServer.Start()

	framework.L.Info().Msgf("Mock HTTP server started on port %d at: %s", port, testServer.URL)
	return testServer, recorder
}

func createTestWorkflowConfig(t *testing.T, workflowName, mockServerURL string) (string, *ecdsa.PrivateKey) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	publicKeyAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	parsedURL, err := url.Parse(mockServerURL)
	require.NoError(t, err, "failed to parse mock server URL")

	url := fmt.Sprintf("%s:%s", framework.HostDockerInternal(), parsedURL.Port())
	framework.L.Info().Msgf("Mock server URL transformed from '%s' to '%s' for Docker access", mockServerURL, url)

	config := map[string]interface{}{
		"authorizedKey": publicKeyAddr.Hex(),
		"url":           url + "/orders",
	}

	configBytes, err := json.Marshal(config)
	require.NoError(t, err, "failed to marshal config")

	configFileName := fmt.Sprintf("test_http_workflow_config_%s.json", workflowName)
	configPath := filepath.Join(os.TempDir(), configFileName)

	err = os.WriteFile(configPath, configBytes, 0644) //nolint:gosec // this is a test file
	require.NoError(t, err, "failed to write config file")

	return configPath, privateKey
}

func createHTTPTriggerRequestWithKey(t *testing.T, workflowName, workflowOwner string, privateKey *ecdsa.PrivateKey) jsonrpc.Request[json.RawMessage] {
	triggerPayload := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowOwner: workflowOwner,
			WorkflowName:  workflowName,
			WorkflowTag:   "TEMP_TAG",
		},
		Input: []byte(`{
			"customer": "test-customer",
			"size": "large",
			"toppings": ["cheese", "pepperoni"],
			"dedupe": false
		}`),
	}

	payloadBytes, err := json.Marshal(triggerPayload)
	require.NoError(t, err)
	rawPayload := json.RawMessage(payloadBytes)

	req := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawPayload,
		ID:      "http-trigger-test-" + uuid.New().String()[0:8],
	}

	token, err := utils.CreateRequestJWT(req)
	require.NoError(t, err)

	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	req.Auth = tokenString

	return req
}

const (
	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""

	// Test configuration constants
	DefaultVerificationTimeout = 5 * time.Minute
	ContainerTargetDir         = "/home/chainlink/workflows"
	WorkflowNodePrefix         = "workflow-node"
	DefaultConfigPath          = "../../../../core/scripts/cre/environment/configs/workflow-don-cache.toml"
	DefaultTopology            = "workflow"
	DefaultEnvironmentDir      = "../../../../core/scripts/cre/environment"
	PoRWorkflowLocation        = "../../../../core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/main.go"
	HTTPWorkflowLocation       = "../../../../core/scripts/cre/environment/examples/workflows/v2/http_simple/main.go"
	DefaultEnvArtifactPath     = "../../../../core/scripts/cre/environment/env_artifact/env_artifact.json"
	RetryInterval              = 2 * time.Second
	ValidationInterval         = 10 * time.Second
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
	chainSelector      uint64
	fullCldEnvironment *cldf.Environment
	workflowName       string
	feedID             string
	sethClient         *seth.Client
	blockchain         *blockchain.Output
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
		err := os.Setenv("ENV_ARTIFACT_PATH", "../../../../core/scripts/cre/environment/env_artifact/env_artifact.json")
		if err != nil {
			return errors.Wrap(err, "failed to set ENV_ARTIFACT_PATH env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}
