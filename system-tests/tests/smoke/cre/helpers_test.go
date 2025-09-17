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
package cre

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"slices"
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
	evmread_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evmread/config"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	crefunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
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
func getWritableChainsFromSavedEnvironmentState(t *testing.T, testEnv *TestEnvironment) []uint64 {
	t.Helper()

	testLogger := framework.L
	testLogger.Info().Msg("Getting writable chains from saved environment state.")
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
	testLogger.Info().Msgf("Writable chains: '%v'", writeableChains)
	return writeableChains
}

/*
Starts Beholder
1. Starts Beholder if it is not running already
2. Loads Beholder stack cache to get Kafka connection details
3. Starts a Kafka listener for Beholder messages

Returns:
1. Context for the listener (with timeout)
2. Channel to receive messages
3. Channel to receive errors from the listener

Recommendation: Use it in tests that need to listen for Beholder messages.
*/
func startBeholder(t *testing.T, testLogger zerolog.Logger, testEnv *TestEnvironment) (context.Context, <-chan proto.Message, <-chan error) {
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
	return listenerCtx, beholderMsgChan, beholderErrChan
}

// Logs all messages received from Beholder until the context is done
func logBeholderMessages(ctx context.Context, t *testing.T, testLogger zerolog.Logger, testEnv *TestEnvironment, messageChan <-chan proto.Message, errChan <-chan error) {
	t.Helper()

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errChan:
			require.FailNowf(t, "Kafka error received from Kafka %s", err.Error())
		case msg := <-messageChan:
			switch typedMsg := msg.(type) {
			case *commonevents.BaseMessage:
				testLogger.Info().Msgf("Received BaseMessage from Beholder: %s", typedMsg.Msg)
			case *workflowevents.UserLogs:
				for _, logLine := range typedMsg.LogLines {
					testLogger.Info().Msgf("Received workflow msg: %s", logLine.Message)
				}
			default:
				testLogger.Info().Msgf("Received unknown message of type '%T'", msg)
			}
		}
	}
}

/*
Asserts that a specific log message is received from a Beholder within a timeout period.
Returns an error if found in error channel or timeouts if a log message is not received.
*/
func assertBeholderMessage(ctx context.Context, t *testing.T, expectedLog string, testLogger zerolog.Logger, messageChan <-chan proto.Message, kafkaErrChan <-chan error, timeout time.Duration) error {
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
					testLogger.Info().Msg("🎉 Received UserLogs message in test")
					receivedUserLogs++

					for _, logLine := range typedMsg.LogLines {
						if strings.Contains(logLine.Message, expectedLog) {
							testLogger.Info().
								Str("expected_log", expectedLog).
								Str("found_message", strings.TrimSpace(logLine.Message)).
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
							Msg("Received UserLogs message, but it does not match expected log")
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
		testLogger.Error().Msg("Timed out waiting for expected user log message")
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
func createAndFundAddresses(t *testing.T, testLogger zerolog.Logger, numberOfAddressesToCreate int, amountToFund *big.Int, sethClient *seth.Client) ([]common.Address, error) {
	t.Helper()

	testLogger.Info().Msgf("Creating and funding %d addresses...", numberOfAddressesToCreate)
	var addressesToRead []common.Address

	for i := 0; i < numberOfAddressesToCreate; i++ {
		addressToRead, _, addrErr := crecrypto.GenerateNewKeyPair()
		require.NoError(t, addrErr, "failed to generate address to read")
		orderNum := i + 1
		testLogger.Info().Msgf("Generated address #%d: %s", orderNum, addressToRead.Hex())

		testLogger.Info().Msgf("Funding address '%s' with amount of '%s' wei", addressToRead.Hex(), amountToFund.String())
		receipt, funErr := crefunding.SendFunds(t.Context(), testLogger, sethClient, crefunding.FundsToSend{
			ToAddress:  addressToRead,
			Amount:     amountToFund,
			PrivateKey: sethClient.MustGetRootPrivateKey(),
		})
		require.NoError(t, funErr, "failed to send funds")
		testLogger.Info().Msgf("Funds sent successfully to address '%s': txHash='%s'", addressToRead.Hex(), receipt.TxHash)

		addressesToRead = append(addressesToRead, addressToRead)
	}

	return addressesToRead, nil
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
		evmread_config.Config
}

// None represents an empty workflow configuration
// It is used to satisfy the workflowConfigFactory, avoiding workflow config creation
type None struct{}

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
			testLogger.Info().Msg("PoR Workflow config file created.")

		case *crontypes.WorkflowConfig:
			workflowCfgFilePath, configErr := createWorkflowYamlConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create Cron workflow config file")
			testLogger.Info().Msg("Cron Workflow config file created.")

		case *HTTPWorkflowConfig:
			workflowCfgFilePath, configErr := createHTTPWorkflowConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create HTTP workflow config file")
			testLogger.Info().Msg("HTTP Workflow config file created.")

		case *evmread_config.Config:
			var configErr error
			workflowConfigFilePath, configErr = createWorkflowYamlConfigFile(workflowName, cfg)
			require.NoError(t, configErr, "failed to create evmread workflow config file")
			testLogger.Info().Msg("EVM Read Workflow config file created.")
		default:
			require.NoError(t, fmt.Errorf("unsupported workflow config type: %T", cfg))
		}
	}
	return workflowConfigFilePath
}

/*
Creates .yaml workflow configuration file and returns the absolute path to the created config file.
*/
func createWorkflowYamlConfigFile(workflowName string, workflowConfig any) (string, error) {
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

	switch tv.Version.Major() {
	case 2:
		// TODO(CRE-876): delete with workflowID
		return
	default:
	}
	deleteErr := creworkflow.DeleteWithContract(t.Context(), blockchainOutputs[0].SethClient, workflowRegistryAddress, tv, uniqueWorkflowName)
	require.NoError(t, deleteErr, "failed to delete workflow '%s'. Please delete/unregister it manually.", uniqueWorkflowName)
}

func compileAndDeployWorkflow[T WorkflowConfig](t *testing.T,
	testEnv *TestEnvironment, testLogger zerolog.Logger, workflowName string,
	workflowConfig *T, workflowFileLocation string,
) {
	t.Helper()

	testLogger.Info().Msgf("compiling and registering workflow '%s'", workflowName)
	homeChainSelector := testEnv.WrappedBlockchainOutputs[0].ChainSelector

	workflowDON, donErr := flags.OneDonMetadataWithFlag(testEnv.FullCldEnvOutput.DonTopology.ToDonMetadata(), cre.WorkflowDON)
	require.NoError(t, donErr, "failed to get find workflow DON in the topology")
	compressedWorkflowWasmPath, workflowConfigPath := createWorkflowArtifacts(t, testLogger, workflowName, workflowDON.Name, workflowConfig, workflowFileLocation)

	// Ignoring the deprecation warning as the suggest solution is not working in CI
	//lint:ignore SA1019 ignoring deprecation warning for this usage
	workflowRegistryAddress, tv, workflowRegistryErr := crecontracts.FindAddressesForChain(
		testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //nolint:staticcheck // SA1019 ignoring deprecation warning for this usage
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
		DonID:                       testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
		ContainerTargetDir:          creworkflow.DefaultWorkflowTargetDir,
		WrappedBlockchainOutputs:    testEnv.WrappedBlockchainOutputs,
	}
	registerWorkflow(t.Context(), t, workflowRegConfig, testEnv.WrappedBlockchainOutputs[0].SethClient, testLogger)
}
