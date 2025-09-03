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
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
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

	var testLogger = framework.L
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

//////////////////////////////
// WORKFLOW-RELATED HELPERS //
//////////////////////////////

// Generic WorkflowConfig interface for creation of different workflow configurations
// Register your workflow configuration types here
type WorkflowConfig interface {
	None | portypes.WorkflowConfig | HTTPWorkflowConfig
}

// None represents an empty workflow configuration
// It is used to satisfy the workflowConfigFactory, avoiding workflow config creation
type None struct{}

// WorkflowRegistrationConfig holds configuration for workflow registration
type WorkflowRegistrationConfig struct {
	WorkflowName         string
	WorkflowLocation     string
	ConfigFilePath       string
	CompressedWasmPath   string
	SecretsURL           string
	WorkflowRegistryAddr common.Address
	DonID                uint64
	ContainerTargetDir   string
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
func createWorkflowArtifacts[T WorkflowConfig](t *testing.T, testLogger zerolog.Logger, workflowName string, workflowConfig *T, workflowFileLocation string) (string, string) {
	t.Helper()

	workflowConfigFilePath := workflowConfigFactory(t, testLogger, workflowName, workflowConfig)
	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(workflowFileLocation, workflowName)
	require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)
	testLogger.Info().Msg("Workflow compiled successfully.")

	// Copy workflow artifacts to Docker containers to use blockchain client running inside for workflow registration
	testLogger.Info().Msg("Copying workflow artifacts to Docker containers.")
	copyErr := creworkflow.CopyArtifactsToDockerContainers(creworkflow.DefaultWorkflowTargetDir, creworkflow.DefaultWorkflowNodePattern, compressedWorkflowWasmPath, workflowConfigFilePath)
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

		case *HTTPWorkflowConfig:
			workflowCfgFilePath, configErr := createHTTPWorkflowConfigFile(workflowName, cfg)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create HTTP workflow config file")
			testLogger.Info().Msg("HTTP Workflow config file created.")

		default:
			require.NoError(t, fmt.Errorf("unsupported workflow config type: %T", cfg))
		}
	}
	return workflowConfigFilePath
}

/*
Registers a workflow with the specified configuration.
*/
func registerWorkflow(ctx context.Context, t *testing.T, workflowConfig *WorkflowRegistrationConfig, sethClient *seth.Client, testLogger zerolog.Logger) {
	t.Helper()

	workflowRegistryAddress := workflowConfig.WorkflowRegistryAddr
	donID := workflowConfig.DonID
	workflowName := workflowConfig.WorkflowName
	binaryURL := "file://" + workflowConfig.CompressedWasmPath
	configURL := ptr.Ptr("file://" + workflowConfig.ConfigFilePath)
	containerTargetDir := &workflowConfig.ContainerTargetDir

	if workflowConfig.ConfigFilePath == "" {
		configURL = nil
	}

	workflowID, registerErr := creworkflow.RegisterWithContract(
		ctx,
		sethClient,
		workflowRegistryAddress,
		donID,
		workflowName,
		binaryURL,
		configURL,
		nil, // no secrets yet
		containerTargetDir,
	)
	require.NoError(t, registerErr, "failed to register workflow '%s'", workflowConfig.WorkflowName)
	testLogger.Info().Msgf("Workflow registered successfully: '%s'", workflowID)
}

/*
Deletes workflows from:
 1. Local environment
 2. Workflow Registry

Recommendation:
Use it at the end of your test to `t.Cleanup()` the env after test run
*/
func deleteWorkflows(t *testing.T, uniqueWorkflowName string, workflowConfigFilePath string, compressedWorkflowWasmPath string, blockchainOutputs []*cre.WrappedBlockchainOutput, workflowRegistryAddress common.Address) {
	t.Helper()

	var testLogger = framework.L
	testLogger.Info().Msgf("Deleting workflow artifacts (%s) after test.\n", uniqueWorkflowName)
	localEnvErr := creworkflow.RemoveWorkflowArtifactsFromLocalEnv(workflowConfigFilePath, compressedWorkflowWasmPath)
	require.NoError(t, localEnvErr, "failed to remove workflow artifacts from local environment")

	deleteErr := creworkflow.DeleteWithContract(t.Context(), blockchainOutputs[0].SethClient, workflowRegistryAddress, uniqueWorkflowName)
	require.NoError(t, deleteErr, "failed to delete workflow '%s'. Please delete/unregister it manually.", uniqueWorkflowName)
}

func compileAndDeployWorkflow[T WorkflowConfig](t *testing.T, testEnv *TestEnvironment, testLogger zerolog.Logger, workflowName string, workflowConfig *T, workflowFileLocation string) {
	homeChainSelector := testEnv.WrappedBlockchainOutputs[0].ChainSelector

	compressedWorkflowWasmPath, workflowConfigPath := createWorkflowArtifacts(t, testLogger, workflowName, workflowConfig, workflowFileLocation)

	// Ignoring the deprecation warning as the suggest solution is not working in CI
	//lint:ignore SA1019 ignoring deprecation warning for this usage
	workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
		testEnv.FullCldEnvOutput.Environment.ExistingAddresses, //lint:ignore SA1019 ignoring deprecation warning for this usage
		homeChainSelector, keystone_changeset.WorkflowRegistry.String())
	require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", testEnv.WrappedBlockchainOutputs[0].ChainID)

	t.Cleanup(func() {
		deleteWorkflows(t, workflowName, workflowConfigPath, compressedWorkflowWasmPath, testEnv.WrappedBlockchainOutputs, workflowRegistryAddress)
	})

	workflowRegConfig := &WorkflowRegistrationConfig{
		WorkflowName:         workflowName,
		WorkflowLocation:     workflowFileLocation,
		ConfigFilePath:       workflowConfigPath,
		CompressedWasmPath:   compressedWorkflowWasmPath,
		WorkflowRegistryAddr: workflowRegistryAddress,
		DonID:                testEnv.FullCldEnvOutput.DonTopology.DonsWithMetadata[0].ID,
		ContainerTargetDir:   creworkflow.DefaultWorkflowTargetDir,
	}
	registerWorkflow(t.Context(), t, workflowRegConfig, testEnv.WrappedBlockchainOutputs[0].SethClient, testLogger)
}
