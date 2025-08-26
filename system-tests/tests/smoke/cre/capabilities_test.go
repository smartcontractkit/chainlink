package cre

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/rand"
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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	common_events "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflow_events "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
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

const (
	DefaultLocalCreDir     = "../../../../core/scripts/cre/environment"
	DefaultWorkflowDonFile = "../../../../core/scripts/cre/environment/configs/workflow-don.toml"
	DefaultEnvArtifactFile = "../../../../core/scripts/cre/environment/env_artifact/env_artifact.json"
)

/*
To execute on local start the local CRE first with following command:
# inside core/scripts/cre/environment directory
1. ensure necessary capabilities (cron, readcontract) are added (see README in core/scripts/cre/environment for [extra_capabilities])
2. `go run . env start && ctf obs up && ctf bs up`.
It will start env + observability + blockscout.
*/
func Test_CRE_Workflow_Don(t *testing.T) {
	confErr := setConfigurationIfMissing(DefaultWorkflowDonFile, DefaultEnvArtifactFile)
	require.NoError(t, confErr, "failed to set default configuration")

	createErr := createEnvironmentIfNotExists(DefaultLocalCreDir)
	require.NoError(t, createErr, "failed to create environment")

	// transform the config file to the cache file, so that we can use the cached environment
	cachedConfigFile, cacheErr := ctfConfigToCacheFile()
	require.NoError(t, cacheErr, "failed to get cached config file")

	setErr := os.Setenv("CTF_CONFIGS", cachedConfigFile)
	require.NoError(t, setErr, "failed to set CTF_CONFIGS env var")

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

	t.Run("Beholder test", func(t *testing.T) {
		executeBeholderTest(t, in, envArtifact)
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
		DEPLOY DATA FEEDS CACHE + READ BALANCES CONTRACTS ON ALL CHAINS (except read-only ones)
		Workflow will write price data to the data feeds cache contract

		REGISTER ONE WORKFLOW PER CHAIN (except read-only ones)
	*/
	for idx, bcOutput := range wrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}
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

		chainSelector := bcOutput.ChainSelector

		testLogger.Info().Msgf("Deploying additional contracts to %d", chainSelector)
		testLogger.Info().Msg("Deploying Data Feeds Cache contract...")
		deployDfConfig := df_changeset_types.DeployConfig{
			ChainsToDeploy: []uint64{chainSelector},
			Labels:         []string{"data-feeds"}, // label required by the changeset
		}
		dfOutput, dfErr := changeset.RunChangeset(df_changeset.DeployCacheChangeset, *fullCldEnvOutput.Environment, deployDfConfig)
		require.NoError(t, dfErr, "failed to deploy Data Feed Cache contract")
		mergeErr := fullCldEnvOutput.Environment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr, "failed to merge address book of Data Feeds Cache contract")
		testLogger.Info().Msgf("Data Feeds Cache contract deployed to %d", chainSelector)

		testLogger.Info().Msg("Deploying Read Balances contract...")
		deployReadBalanceRequest := &keystone_changeset.DeployRequestV2{ChainSel: chainSelector}
		rbOutput, rbErr := keystone_changeset.DeployBalanceReaderV2(*fullCldEnvOutput.Environment, deployReadBalanceRequest)
		require.NoError(t, rbErr, "failed to deploy Read Balances contract")
		mergeErr2 := fullCldEnvOutput.Environment.ExistingAddresses.Merge(rbOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr2, "failed to merge address book of Read Balances contract")
		testLogger.Info().Msgf("Read Balances contract deployed to %d", chainSelector)

		mergeErr3 := dfOutput.DataStore.Merge(rbOutput.DataStore.Seal())
		require.NoError(t, mergeErr3, "failed to merge data stores")
		fullCldEnvOutput.Environment.DataStore = dfOutput.DataStore.Seal()

		workflowName := "por-workflow-" + bcOutput.BlockchainOutput.ChainID + "-" + uuid.New().String()[0:4]

		dfConfigInput := &configureDataFeedsCacheInput{
			chainSelector:      chainSelector,
			fullCldEnvironment: fullCldEnvOutput.Environment,
			workflowName:       workflowName,
			feedID:             feedIDs[idx],
			sethClient:         bcOutput.SethClient,
			blockchain:         bcOutput.BlockchainOutput,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		chainID := bcOutput.ChainID
		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			homeChainSelector, // it should live only on one chain, it is not deployed to all chains
			keystone_changeset.WorkflowRegistry.String(),
		)
		require.NoError(t, workflowRegistryErr, "failed to find Workflow Registry address.")
		testLogger.Info().Msgf("Workflow Registry contract found at chain selector %d at %s", homeChainSelector, workflowRegistryAddress)

		testLogger.Info().Msgf("Registering PoR workflow on chain %d (%d)", chainID, chainSelector)
		testLogger.Info().Msg("Creating workflow config file.")
		readBalancesAddress, readContractErr := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			chainSelector,
			keystone_changeset.BalanceReader.String(),
		)
		require.NoError(t, readContractErr, "failed to find Read Balances contract address for chain %d", chainID)
		testLogger.Info().Msgf("Read Balances contract found on chain %d at %s", chainID, readBalancesAddress)

		dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
			fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			chainSelector,
			df_changeset.DataFeedsCache.String(),
		)
		require.NoError(t, dataFeedsCacheErr, "failed to find Data Feeds Cache address for chain %d", chainID)
		testLogger.Info().Msgf("Data Feeds Cache contract found on chain %d at %s", chainID, dataFeedsCacheAddress)

		workflowConfigFilePath, configErr := createWorkflowConfigFile(bcOutput, readBalancesAddress, dataFeedsCacheAddress, workflowName, feedIDs[idx], priceProvider.URL(), corevm.GenerateWriteTargetName(chainID))
		require.NoError(t, configErr, "failed to create workflow config file")
		testLogger.Info().Msgf("Workflow config file created.")

		compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
		require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)
		testLogger.Info().Msgf("Workflow compiled successfully.")

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

		workflowCopyErr := creworkflow.CopyArtifactToDockerContainers(compressedWorkflowWasmPath, creworkflow.DefaultWorkflowNodePattern, creworkflow.DefaultWorkflowTargetDir)
		require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")

		configCopyErr := creworkflow.CopyArtifactToDockerContainers(workflowConfigFilePath, creworkflow.DefaultWorkflowNodePattern, creworkflow.DefaultWorkflowTargetDir)
		require.NoError(t, configCopyErr, "failed to copy workflow config to docker containers")

		workflowID, registerErr := creworkflow.RegisterWithContract(
			t.Context(),
			wrappedBlockchainOutputs[0].SethClient, // crucial to use Seth Client connected to home chain (first chain in the set)
			workflowRegistryAddress,
			fullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
			workflowName,
			"file://"+compressedWorkflowWasmPath,
			ptr.Ptr("file://"+workflowConfigFilePath),
			nil,
			&creworkflow.DefaultWorkflowTargetDir,
		)
		require.NoError(t, registerErr, "failed to register PoR workflow")
		testLogger.Info().Msgf("Workflow registered successfully: '%s'", workflowID)
	}

	/*
		START THE VALIDATION PHASE
		Check whether each feed has been updated with the expected prices, which workflow fetches from the price provider
	*/
	eg := &errgroup.Group{}
	for idx, bcOutput := range wrappedBlockchainOutputs {
		if bcOutput.BlockchainOutput.Type == blockchain.FamilySolana {
			continue
		}
		eg.Go(func() error {
			feedID := feedIDs[idx]
			testLogger.Info().Msgf("Waiting for feed %s to update...", feedID)

			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
				fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				bcOutput.ChainSelector,
				df_changeset.DataFeedsCache.String(),
			)
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", bcOutput.ChainID)

			dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, bcOutput.SethClient.Client)
			require.NoError(t, instanceErr, "failed to create data feeds cache instance")

			startTime := time.Now()
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

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	testLogger.Info().Msgf("All prices were found for all feeds")
}

func executeVaultTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact) {
	// Skip till we figure out and fix the issues with environment startup on this test
	t.Skip()
	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	fullCldEnvOutput, _, loadErr := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

	framework.L.Info().Msg("Getting gateway configuration...")
	require.NotEmpty(t, fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations, "expected at least one gateway configuration")
	gatewayURL, err := url.Parse(fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Protocol + "://" + fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Host + ":" + strconv.Itoa(fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.ExternalPort) + fullCldEnvOutput.DonTopology.GatewayConnectorOutput.Configurations[0].Incoming.Path)
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

const (
	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""
)

func createEnvironmentIfNotExists(environmentDir string) error {
	cachedConfigFile, cacheErr := ctfConfigToCacheFile()
	if cacheErr != nil {
		return errors.Wrap(cacheErr, "failed to get cached config file")
	}

	if _, err := os.Stat(cachedConfigFile); os.IsNotExist(err) {
		framework.L.Info().Str("cached_config_file", cachedConfigFile).Msg("Cached config file does not exist, starting environment...")
		cmd := exec.Command("go", "run", ".", "env", "start")
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

// Creates workflow configuration file storing the necessary values used by a workflow (i.e. feedID, read/write contract addresses)
// The values are written to types.WorkflowConfig
func createWorkflowConfigFile(bcOutput *cre.WrappedBlockchainOutput, readContractAddress, feedsConsumerAddress common.Address, workflowName, feedID, dataURL, writeTargetName string) (string, error) {
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return "", errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID
	chainFamily := bcOutput.BlockchainOutput.Family
	chainID := bcOutput.BlockchainOutput.ChainID

	workflowConfig := portypes.WorkflowConfig{
		ChainFamily: chainFamily,
		ChainID:     chainID,
		BalanceReaderConfig: portypes.BalanceReaderConfig{
			BalanceReaderAddress: readContractAddress.Hex(),
		},
		ComputeConfig: portypes.ComputeConfig{
			FeedID:                feedIDToUse,
			URL:                   dataURL,
			DataFeedsCacheAddress: feedsConsumerAddress.Hex(),
			WriteTargetName:       writeTargetName,
		},
	}

	// Write workflow config to a file
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

func executeBeholderTest(t *testing.T, in *envconfig.Config, envArtifact environment.EnvArtifact) {
	testLogger := framework.L
	cldLogger := cldlogger.NewSingleFileLogger(t)

	/*
		BUILD ENVIRONMENT FROM SAVED STATE
	*/
	fullCldEnvOutput, wrappedBlockchainOutputs, loadErr := environment.BuildFromSavedState(t.Context(), cldLogger, in, envArtifact)
	require.NoError(t, loadErr, "failed to load environment")

	bErr := startBeholderStackIfIsNotRunning(DefaultBeholderStackCacheFile, DefaultLocalCreDir)
	require.NoError(t, bErr, "failed to start Beholder")

	chipConfig, chipErr := loadBeholderStackCache()
	require.NoError(t, chipErr, "failed to load chip ingress cache")
	require.NotNil(t, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, "kafka external url is not set in the cache")
	require.NotEmpty(t, chipConfig.Kafka.Topics, "kafka topics are not set in the cache")

	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbeholder"

	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
	require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)

	homeChainSelector := wrappedBlockchainOutputs[0].ChainSelector
	workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
		fullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		homeChainSelector,
		keystone_changeset.WorkflowRegistry.String(),
	)
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", wrappedBlockchainOutputs[0].ChainID)

	t.Cleanup(func() {
		wasmErr := os.Remove(compressedWorkflowWasmPath)
		if wasmErr != nil {
			framework.L.Warn().Msgf("failed to remove workflow wasm file %s: %s", compressedWorkflowWasmPath, wasmErr.Error())
		}

		deleteErr := creworkflow.DeleteWithContract(t.Context(), wrappedBlockchainOutputs[0].SethClient, workflowRegistryAddress, workflowName)
		if deleteErr != nil {
			framework.L.Warn().Msgf("failed to delete workflow %s: %s. Please delete it manually.", workflowName, deleteErr.Error())
		}
	})

	workflowCopyErr := creworkflow.CopyArtifactToDockerContainers(compressedWorkflowWasmPath, creworkflow.DefaultWorkflowNodePattern, creworkflow.DefaultWorkflowTargetDir)
	require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")

	_, registerErr := creworkflow.RegisterWithContract(
		t.Context(),
		wrappedBlockchainOutputs[0].SethClient, // crucial to use Seth Client connected to home chain (first chain in the set)
		workflowRegistryAddress,
		fullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
		workflowName,
		"file://"+compressedWorkflowWasmPath,
		nil,
		nil,
		&creworkflow.DefaultWorkflowTargetDir,
	)
	require.NoError(t, registerErr, "failed to register cron beholder workflow")

	listenerCtx, cancelListener := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(func() {
		cancelListener()
	})

	kafkaErrChan := make(chan error, 1)
	messageChan := make(chan proto.Message, 10)

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	messageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflow_events.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &common_events.BaseMessage{}
		},
	}

	// Start listening for messages in the background
	go func() {
		listenForKafkaMessages(listenerCtx, testLogger, chipConfig.ChipIngress.Output.RedPanda.KafkaExternalURL, chipConfig.Kafka.Topics[0], messageTypes, messageChan, kafkaErrChan)
	}()

	expectedUserLog := "Amazing workflow user log"

	foundExpectedLog := make(chan bool, 1) // Channel to signal when expected log is found
	foundErrorLog := make(chan bool, 1)    // Channel to signal when engine initialization failure is detected
	receivedUserLogs := 0
	// Start message processor goroutine
	go func() {
		for {
			select {
			case <-listenerCtx.Done():
				return
			case msg := <-messageChan:
				// Process received messages
				switch typedMsg := msg.(type) {
				case *common_events.BaseMessage:
					if strings.Contains(typedMsg.Msg, "Workflow Engine initialization failed") {
						foundErrorLog <- true
					}
				case *workflow_events.UserLogs:
					testLogger.Info().
						Msg("🎉 Received UserLogs message in test")
					receivedUserLogs++

					for _, logLine := range typedMsg.LogLines {
						if strings.Contains(logLine.Message, expectedUserLog) {
							testLogger.Info().
								Str("expected_log", expectedUserLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
								Msg("🎯 Found expected user log message!")

							select {
							case foundExpectedLog <- true:
							default: // Channel might already have a value
							}
							return // Exit the processor goroutine
						}
						testLogger.Warn().
							Str("expected_log", expectedUserLog).
							Str("found_message", strings.TrimSpace(logLine.Message)).
							Msg("Received UserLogs message, but it does not match expected log")
					}
				default:
					// ignore other message types
				}
			}
		}
	}()

	timeout := 2 * time.Minute

	testLogger.Info().
		Str("expected_log", expectedUserLog).
		Dur("timeout", timeout).
		Msg("Waiting for expected user log message or timeout")

	// Wait for either the expected log to be found, or engine initialization failure to be detected, or timeout (2 minutes)
	select {
	case <-foundExpectedLog:
		testLogger.Info().
			Str("expected_log", expectedUserLog).
			Msg("✅ Test completed successfully - found expected user log message!")
		return
	case <-foundErrorLog:
		require.Fail(t, "Test completed with error - found engine initialization failure message!")
	case <-time.After(timeout):
		testLogger.Error().Msg("Timed out waiting for expected user log message")
		if receivedUserLogs > 0 {
			testLogger.Warn().Int("received_user_logs", receivedUserLogs).Msg("Received some UserLogs messages, but none matched expected log")
		} else {
			testLogger.Warn().Msg("Did not receive any UserLogs messages")
		}
		require.Failf(t, "Timed out waiting for expected user log message", "Expected user log message: '%s' not found after %s", expectedUserLog, timeout.String())
	case err := <-kafkaErrChan:
		testLogger.Error().Err(err).Msg("Kafka listener encountered an error during execution")
		require.Fail(t, "Kafka listener failed", err.Error())
	}

	testLogger.Info().Msg("Beholder test completed")
}

func setConfigurationIfMissing(configName, envArtifactPath string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	if os.Getenv("ENV_ARTIFACT_PATH") == "" {
		err := os.Setenv("ENV_ARTIFACT_PATH", envArtifactPath)
		if err != nil {
			return errors.Wrap(err, "failed to set ENV_ARTIFACT_PATH env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}

func ctfConfigToCacheFile() (string, error) {
	configFile := os.Getenv("CTF_CONFIGS")
	if configFile == "" {
		return "", errors.New("CTF_CONFIGS env var is not set")
	}

	split := strings.Split(configFile, ",")
	return strings.ReplaceAll(split[0], ".toml", "") + "-cache.toml", nil
}
