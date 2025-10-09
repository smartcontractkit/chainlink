// helpers_test.go
//
// This file contains reusable test helper functions that encapsulate common,
// logically grouped test-specific steps. They hide and abstract away
// the complexities of the test setup and execution.
//
// All helpers here are intentionally unexported functions (lowercase)
// so they do not leak outside this package.
//
// By keeping repeated setup and execution logic in one place,
// we make individual tests shorter, clearer, and easier to maintain.
//
// Recommendations:
// 1. Keep naming action-oriented: mustStartDB, withEnv, seedUsers.
// 2. Ensure proper cleanup after steps, where necessary, to avoid side effects.
package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	consensus_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/consensus/config"
	evmread_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmread-negative/config"
	evmwrite_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmwrite-negative/config"
	http_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
	evmread_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	environment "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	crefunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
)

/////////////////////////
// ENVIRONMENT HELPERS //
/////////////////////////

/*
Parse through chain configs and extract "writable" chain IDs.
If a chain requires a Forwarder contract, it is considered a "writable" chain.

Recommendation: Use it to determine on which chains to deploy certain contracts and register workflows.
See an example in a test using PoR workflow.
*/
func GetWritableChainsFromSavedEnvironmentState(t *testing.T, testEnv *ttypes.TestEnvironment) []uint64 {
	t.Helper()

	testLogger := framework.L
	testLogger.Info().Msg("Getting writable chains from saved environment state.")
	writeableChains := []uint64{}
	for _, bcOutput := range testEnv.WrappedBlockchainOutputs {
		for _, don := range testEnv.CreEnvironment.DonTopology.Dons.List() {
			if flags.RequiresForwarderContract(don.Flags, bcOutput.ChainID) {
				if !slices.Contains(writeableChains, bcOutput.ChainID) {
					writeableChains = append(writeableChains, bcOutput.ChainID)
				}
			}
		}
	}
	testLogger.Info().Msgf("Writable chains: '%v'", writeableChains)
	return writeableChains
}

func GetEVMEnabledChains(t *testing.T, testEnv *ttypes.TestEnvironment) map[string]struct{} {
	t.Helper()

	enabledChains := map[string]struct{}{}
	for _, nodeSet := range testEnv.Config.NodeSets {
		require.NoError(t, nodeSet.ParseChainCapabilities())
		if nodeSet.ChainCapabilities == nil || nodeSet.ChainCapabilities[cre.EVMCapability] == nil {
			continue
		}

		for _, chainID := range nodeSet.ChainCapabilities[cre.EVMCapability].EnabledChains {
			strChainID := strconv.FormatUint(chainID, 10)
			enabledChains[strChainID] = struct{}{}
		}
	}
	require.NotEmpty(t, enabledChains, "No chains have EVM capability enabled in any node set")
	return enabledChains
}

/*
Starts Beholder
Recommendation: Use it in tests that need to listen for Beholder messages.
*/
func StartBeholder(t *testing.T, testLogger zerolog.Logger, testEnv *ttypes.TestEnvironment) (context.Context, <-chan proto.Message, <-chan error) {
	t.Helper()

	beholder, err := NewBeholder(framework.L, testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath)
	require.NoError(t, err, "failed to create beholder instance")

	// We are interested in UserLogs (successful execution)
	// or BaseMessage with specific error message (engine initialization failure)
	messageTypes := map[string]func() proto.Message{
		"workflows.v1.UserLogs": func() proto.Message {
			return &workflowevents.UserLogs{}
		},
		"BaseMessage": func() proto.Message {
			return &commonevents.BaseMessage{}
		},
	}

	timeout := 5 * time.Minute
	testLogger.Info().Dur("timeout", timeout).Msg("Starting Beholder listener...")
	listenerCtx, cancelListener := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(func() {
		cancelListener()
		testLogger.Info().Msg("Beholder listener stopped.")
	})

	beholderMsgChan, beholderErrChan := beholder.SubscribeToBeholderMessages(listenerCtx, messageTypes)

	// Fail fast if there is an error from the heartbeat validation subscription
	select {
	case err := <-beholderErrChan:
		require.NoError(t, err, "Beholder subscription failed during initialization")
	default:
		// No immediate error, proceed
	}

	testLogger.Info().Msg("Beholder listener ready")
	return listenerCtx, beholderMsgChan, beholderErrChan
}

/*
Asserts that a specific log message is received from a Beholder within a timeout period.
Returns an error if found in error channel or timeouts if a log message is not received.
*/
func AssertBeholderMessage(ctx context.Context, t *testing.T, expectedLog string, testLogger zerolog.Logger, messageChan <-chan proto.Message, kafkaErrChan <-chan error, timeout time.Duration) error {
	foundExpectedLog := make(chan bool, 1) // Channel to signal when expected log is found
	foundErrorLog := make(chan bool, 1)    // Channel to signal when engine initialization failure is detected
	receivedUserLogs := 0

	// Start message processor goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messageChan:
				// Process received messages
				switch typedMsg := msg.(type) {
				case *commonevents.BaseMessage:
					if strings.Contains(typedMsg.Msg, "Workflow Engine initialization failed") {
						foundErrorLog <- true
					}
				case *workflowevents.UserLogs:
					testLogger.Info().Msg("➡️ Beholder message received in test. Asserting...")
					receivedUserLogs++

					for _, logLine := range typedMsg.LogLines {
						if strings.Contains(logLine.Message, expectedLog) {
							testLogger.Info().
								Str("expected_log", expectedLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
								Str("workflow_id", typedMsg.M.WorkflowExecutionID).
								Msg("🎯 Found expected user log message!")

							select {
							case foundExpectedLog <- true:
							default: // Channel might already have a value
							}
							return // Exit the processor goroutine
						}

						testLogger.Warn().
							Str("expected_log", expectedLog).
							Str("found_message", strings.TrimSpace(logLine.Message)).
							Msg("[soft assertion] Received UserLogs message, but it does not match expected log")
					}
				default:
					// ignore other message types
				}
			}
		}
	}()

	testLogger.Info().
		Str("expected_log", expectedLog).
		Dur("timeout", timeout).
		Msg("Waiting for expected user log message or timeout")

	// Wait for either the expected log to be found, or engine initialization failure to be detected
	select {
	case <-foundExpectedLog:
		testLogger.Info().Str("expected_log", expectedLog).Msg("✅ Test completed successfully - found expected user log message!")
		return nil
	case <-foundErrorLog:
		testLogger.Warn().Msg("beholder found engine initialization failure message! (may be expected in negative tests)")
		return errors.New("beholder message validation completed with error: found engine initialization failure message")
	case <-time.After(timeout):
		testLogger.Error().Str("expected_log", expectedLog).Msg("Timed out waiting for expected user log message")
		if receivedUserLogs > 0 {
			testLogger.Warn().Int("received_user_logs", receivedUserLogs).Msg("Received some UserLogs messages, but none matched expected log")
		} else {
			testLogger.Warn().Msg("Did not receive any UserLogs messages")
		}
		require.Failf(t, "Timed out waiting for the expected user log message (or error)", "Expected user log message: '%s' not found after %s", expectedLog, timeout.String())
	case err := <-kafkaErrChan:
		testLogger.Error().Err(err).Msg("Kafka listener encountered an error during execution. Ensure Beholder is running and accessible.")
		require.Fail(t, "Kafka listener failed", err.Error())
	}
	return nil
}

//////////////////////////////
//      CRYPTO HELPERS      //
//////////////////////////////

// Creates and funds a specified number of new Ethereum addresses on a given chain.
func CreateAndFundAddresses(t *testing.T, testLogger zerolog.Logger, numberOfAddressesToCreate int, amountToFund *big.Int, sethClient *seth.Client, bcOutput *cre.WrappedBlockchainOutput, fullCldEnvOutput *cre.Environment) ([]common.Address, error) {
	t.Helper()

	testLogger.Info().Msgf("Creating and funding %d addresses...", numberOfAddressesToCreate)
	var addressesToRead []common.Address

	for i := range numberOfAddressesToCreate {
		addressToRead, _, addrErr := crecrypto.GenerateNewKeyPair()
		require.NoError(t, addrErr, "failed to generate address to read")
		orderNum := i + 1
		testLogger.Info().Msgf("Generated address #%d: %s", orderNum, addressToRead.Hex())

		testLogger.Info().Msgf("Funding address '%s' with amount of '%s' wei", addressToRead.Hex(), amountToFund.String())

		switch bcOutput.BlockchainOutput.Family {
		case blockchain.FamilyTron:
			if err := environment.FundTronAddress(t.Context(), testLogger, addressToRead, amountToFund.Uint64(), bcOutput, fullCldEnvOutput.CldfEnvironment); err != nil {
				return nil, err
			}
		default:
			if err := fundEthAddress(t, testLogger, addressToRead, amountToFund, sethClient); err != nil {
				return nil, err
			}
		}

		addressesToRead = append(addressesToRead, addressToRead)
	}

	return addressesToRead, nil
}

func fundEthAddress(t *testing.T, testLogger zerolog.Logger, addressToRead common.Address, amountToFund *big.Int, sethClient *seth.Client) error {
	receipt, funErr := crefunding.SendFunds(t.Context(), testLogger, sethClient, crefunding.FundsToSend{
		ToAddress:  addressToRead,
		Amount:     amountToFund,
		PrivateKey: sethClient.MustGetRootPrivateKey(),
	})
	require.NoError(t, funErr, "failed to send funds")
	testLogger.Info().Msgf("Funds sent successfully to address '%s': txHash='%s'", addressToRead.Hex(), receipt.TxHash)
	return nil
}

//////////////////////////////
// WORKFLOW-RELATED HELPERS //
//////////////////////////////

// Generic WorkflowConfig interface for creation of different workflow configurations
// Register your workflow configuration types here
type WorkflowConfig interface {
	None |
		portypes.WorkflowConfig |
		crontypes.WorkflowConfig |
		HTTPWorkflowConfig |
		consensus_negative_config.Config |
		evmread_config.Config |
		evmread_negative_config.Config |
		evmwrite_negative_config.Config |
		http_negative_config.Config
}

// None represents an empty workflow configuration
// It is used to satisfy the workflowConfigFactory, avoiding workflow config creation
type None struct{}

type HTTPWorkflowConfig struct {
	AuthorizedKey common.Address `json:"authorizedKey"`
	URL           string         `json:"url"`
}

// WorkflowRegistrationConfig holds configuration for workflow registration
type WorkflowRegistrationConfig struct {
	WorkflowName                string
	WorkflowLocation            string
	ConfigFilePath              string
	CompressedWasmPath          string
	SecretsURL                  string
	WorkflowRegistryAddr        common.Address
	WorkflowRegistryTypeVersion deployment.TypeAndVersion
	ChainID                     uint64
	DonID                       uint64
	ContainerTargetDir          string
	WrappedBlockchainOutputs    []*cre.WrappedBlockchainOutput
}

/*
Creates the necessary workflow artifacts based on WorkflowConfig:
 1. Configuration for a workflow (or no config if typed nil is passed for workflowConfig);
 2. Compiled and compressed workflow WASM file;
 3. Copies the workflow artifacts to the Docker containers

It returns the paths to:
 1. the compressed WASM file;
 2. the workflow config file.
*/
func createWorkflowArtifacts[T WorkflowConfig](t *testing.T, testLogger zerolog.Logger, workflowName, workflowDONName string, workflowConfig *T, workflowFileLocation string) (string, string) {
	t.Helper()

	workflowConfigFilePath := workflowConfigFactory(t, testLogger, workflowName, workflowConfig)
	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
	require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)
	testLogger.Info().Msg("Workflow compiled successfully.")

	// Copy workflow artifacts to Docker containers to use blockchain client running inside for workflow registration
	testLogger.Info().Msg("Copying workflow artifacts to Docker containers.")
	copyErr := creworkflow.CopyArtifactsToDockerContainers(creworkflow.DefaultWorkflowTargetDir, ns.NodeNamePrefix(workflowDONName), compressedWorkflowWasmPath, workflowConfigFilePath)
	require.NoError(t, copyErr, "failed to copy workflow artifacts to docker containers")
	testLogger.Info().Msg("Workflow artifacts successfully copied to the Docker containers.")

	return compressedWorkflowWasmPath, workflowConfigFilePath
}

/*
Creates the necessary workflow configuration based on a type registered in the WorkflowConfig interface
Pass `nil` to skip workflow config file creation.

Returns the path to the workflow config file.
*/
func workflowConfigFactory[T WorkflowConfig](t *testing.T, testLogger zerolog.Logger, workflowName string, workflowConfig *T) (filePath string) {
	t.Helper()

	var workflowConfigFilePath string

	// nil is an acceptable argument that allows skipping config file creation when it is not necessary
	if workflowConfig != nil {
		switch cfg := any(workflowConfig).(type) {
		case *None:
			workflowConfigFilePath = ""
			testLogger.Info().Msg("Workflow config file is not requested and will not be created.")

		case *portypes.WorkflowConfig:
			workflowCfgFilePath, configErr := createPoRWorkflowConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create PoR workflow config file")
			testLogger.Info().Msg("PoR workflow config file created.")

		case *HTTPWorkflowConfig:
			workflowCfgFilePath, configErr := createHTTPWorkflowConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create HTTP workflow config file")
			testLogger.Info().Msg("HTTP workflow config file created.")

		case *crontypes.WorkflowConfig:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create Cron workflow config file")
			testLogger.Info().Msg("Cron workflow config file created.")

		case *consensus_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create consensus workflow config file")
			testLogger.Info().Msg("Consensus workflow config file created.")

		case *evmread_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmread workflow config file")
			testLogger.Info().Msg("EVM Read workflow config file created.")

		case *evmread_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmread-negative workflow config file")
			testLogger.Info().Msg("EVM Read negative workflow config file created.")

		case *evmwrite_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmwrite-negative workflow config file")
			testLogger.Info().Msg("EVM Write negative workflow config file created.")

		case *http_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create http-negative workflow config file")
			testLogger.Info().Msg("HTTP negative workflow config file created.")

		default:
			require.NoError(t, fmt.Errorf("unsupported workflow config type: %T", cfg))
		}
	}
	return workflowConfigFilePath
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

	return CreateWorkflowYamlConfigFile(workflowName, workflowConfig)
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

func createHTTPWorkflowConfigFile(workflowName string, cfg *HTTPWorkflowConfig) (string, error) {
	testLogger := framework.L
	mockServerURL := cfg.URL
	parsedURL, urlErr := url.Parse(mockServerURL)
	if urlErr != nil {
		return "", errors.Wrap(urlErr, "failed to parse HTTP mock server URL")
	}

	url := fmt.Sprintf("%s:%s", framework.HostDockerInternal(), parsedURL.Port())
	testLogger.Info().Msgf("Mock server URL transformed from '%s' to '%s' for Docker access", mockServerURL, url)

	// override values in the initial workflow configuration
	cfg.URL = url + "/orders"

	configBytes, marshalErr := json.Marshal(cfg)
	if marshalErr != nil {
		return "", errors.Wrap(marshalErr, "failed to marshal HTTP workflow config")
	}

	configFileName := fmt.Sprintf("test_http_workflow_config_%s.json", workflowName)
	configPath := filepath.Join(os.TempDir(), configFileName)

	writeErr := os.WriteFile(configPath, configBytes, 0o644) //nolint:gosec // this is a test file
	if writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write HTTP workflow config file")
	}

	return configPath, nil
}

/*
Creates .yaml workflow configuration file and returns the absolute path to the created config file.
*/
func CreateWorkflowYamlConfigFile(workflowName string, workflowConfig any) (string, error) {
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

	if err := os.WriteFile(workflowConfigOutputFile, configMarshalled, 0o644); err != nil { //nolint:gosec // G306: we want it to be readable by everyone
		return "", errors.Wrap(err, "failed to write output file")
	}

	outputFileAbsPath, outputFileAbsPathErr := filepath.Abs(workflowConfigOutputFile)
	if outputFileAbsPathErr != nil {
		return "", errors.Wrap(outputFileAbsPathErr, "failed to get absolute path of the config file")
	}

	return outputFileAbsPath, nil
}

/*
Registers a workflow with the specified configuration.
*/
func registerWorkflow(ctx context.Context, t *testing.T,
	wfRegCfg *WorkflowRegistrationConfig, sethClient *seth.Client,
	testLogger zerolog.Logger,
) {
	t.Helper()

	t.Cleanup(func() {
		deleteWorkflows(t, wfRegCfg.WorkflowName, wfRegCfg.ConfigFilePath,
			wfRegCfg.CompressedWasmPath, wfRegCfg.WrappedBlockchainOutputs,
			wfRegCfg.WorkflowRegistryAddr, wfRegCfg.WorkflowRegistryTypeVersion,
		)
	})

	donID := wfRegCfg.DonID
	workflowName := wfRegCfg.WorkflowName
	binaryURL := "file://" + wfRegCfg.CompressedWasmPath
	configURL := ptr.Ptr("file://" + wfRegCfg.ConfigFilePath)
	containerTargetDir := &wfRegCfg.ContainerTargetDir

	if wfRegCfg.ConfigFilePath == "" {
		configURL = nil
	}

	workflowID, registerErr := creworkflow.RegisterWithContract(
		ctx,
		sethClient,
		wfRegCfg.WorkflowRegistryAddr,
		wfRegCfg.WorkflowRegistryTypeVersion,
		donID,
		workflowName,
		binaryURL,
		configURL,
		nil, // no secrets yet
		containerTargetDir,
	)
	require.NoError(t, registerErr, "failed to register workflow '%s'", wfRegCfg.WorkflowName)
	testLogger.Info().Msgf("Workflow registered successfully: '%s'", workflowID)
}

/*
Deletes workflows from:
 1. Local environment
 2. Workflow Registry

Recommendation:
Use it at the end of your test to `t.Cleanup()` the env after test run
*/
func deleteWorkflows(t *testing.T, uniqueWorkflowName string,
	workflowConfigFilePath string, compressedWorkflowWasmPath string,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	workflowRegistryAddress common.Address,
	tv deployment.TypeAndVersion,
) {
	t.Helper()

	testLogger := framework.L
	testLogger.Info().Msgf("Deleting workflow artifacts (%s) after test.", uniqueWorkflowName)
	localEnvErr := creworkflow.RemoveWorkflowArtifactsFromLocalEnv(workflowConfigFilePath, compressedWorkflowWasmPath)
	require.NoError(t, localEnvErr, "failed to remove workflow artifacts from local environment")

	deleteErr := creworkflow.DeleteWithContract(t.Context(), blockchainOutputs[0].SethClient, workflowRegistryAddress, tv, uniqueWorkflowName)
	require.NoError(t, deleteErr, "failed to delete workflow '%s'. Please delete/unregister it manually.", uniqueWorkflowName)
	testLogger.Info().Msgf("Workflow '%s' deleted successfully from the registry.", uniqueWorkflowName)
}

func CompileAndDeployWorkflow[T WorkflowConfig](t *testing.T,
	testEnv *ttypes.TestEnvironment, testLogger zerolog.Logger, workflowName string,
	workflowConfig *T, workflowFileLocation string,
) {
	t.Helper()

	testLogger.Info().Msgf("compiling and registering workflow '%s'", workflowName)
	homeChainSelector := testEnv.WrappedBlockchainOutputs[0].ChainSelector

	workflowDOName := ""
	for _, don := range testEnv.CreEnvironment.DonTopology.Dons.List() {
		if don.ID == testEnv.CreEnvironment.DonTopology.WorkflowDonID {
			workflowDOName = don.Name
			break
		}
	}
	require.NotEmpty(t, workflowDOName, "failed to find workflow DON in the topology")

	compressedWorkflowWasmPath, workflowConfigPath := createWorkflowArtifacts(t, testLogger, workflowName, workflowDOName, workflowConfig, workflowFileLocation)

	// Ignoring the deprecation warning as the suggest solution is not working in CI
	//lint:ignore SA1019 ignoring deprecation warning for this usage
	workflowRegistryAddress, tv, workflowRegistryErr := crecontracts.FindAddressesForChain(
		testEnv.CreEnvironment.CldfEnvironment.ExistingAddresses, //nolint:staticcheck // SA1019 ignoring deprecation warning for this usage
		homeChainSelector, keystone_changeset.WorkflowRegistry.String())
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", testEnv.WrappedBlockchainOutputs[0].ChainID)

	workflowRegConfig := &WorkflowRegistrationConfig{
		WorkflowName:                workflowName,
		WorkflowLocation:            workflowFileLocation,
		ConfigFilePath:              workflowConfigPath,
		CompressedWasmPath:          compressedWorkflowWasmPath,
		WorkflowRegistryAddr:        workflowRegistryAddress,
		WorkflowRegistryTypeVersion: tv,
		ChainID:                     homeChainSelector,
		DonID:                       testEnv.CreEnvironment.DonTopology.Dons.List()[0].ID,
		ContainerTargetDir:          creworkflow.DefaultWorkflowTargetDir,
		WrappedBlockchainOutputs:    testEnv.WrappedBlockchainOutputs,
	}
	registerWorkflow(t.Context(), t, workflowRegConfig, testEnv.WrappedBlockchainOutputs[0].SethClient, testLogger)
}
