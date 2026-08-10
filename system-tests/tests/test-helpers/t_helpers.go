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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/solkey"
	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/cron/types"
	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/proof-of-reserve/cron-based/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
	stellarfeature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/stellar"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	crecrypto "github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
	consensus_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/consensus/config"
	evmread_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmread-negative/config"
	evmwrite_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/evmwrite-negative/config"
	logtrigger_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/evm/logtrigger-negative/config"
	http_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/http/config"
	httpaction_negative_config "github.com/smartcontractkit/chainlink/system-tests/tests/regression/cre/httpaction-negative/config"
	aptoswrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite/config"
	aptoswriteroundtrip_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswriteroundtrip/config"
	evmread_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/evmread/config"
	logtrigger_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/evm/logtrigger/config"
	httpaction_smoke_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/httpaction/config"
	sollogtrigger_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/config"
	solread_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solread/config"
	solwrite_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	vaultsecret_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/vaultsecret/config"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const WorkflowEngineInitErrorLog = "Workflow Engine initialization failed"
const maxWorkflowNameLen = 64

// defaultDeployMaxParallel is used when CRE_TEST_DEPLOY_MAX_PARALLEL is unset or invalid.
const defaultDeployMaxParallel = 20

var deleteWorkflowsMu sync.Mutex

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
	for _, bcOutput := range testEnv.CreEnvironment.Blockchains {
		for _, don := range testEnv.Dons.List() {
			if flags.RequiresForwarderContract(don.Flags, bcOutput.ChainID()) {
				if !slices.Contains(writeableChains, bcOutput.ChainID()) {
					writeableChains = append(writeableChains, bcOutput.ChainID())
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
		enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(cre.EVMCapability)
		require.NoError(t, err, "failed to get enabled chain IDs for EVM capability")

		for _, chainID := range enabledChainIDs {
			strChainID := strconv.FormatUint(chainID, 10)
			enabledChains[strChainID] = struct{}{}
		}
	}
	require.NotEmpty(t, enabledChains, "No chains have EVM capability enabled in any node set")
	return enabledChains
}

/*
Starts ChipIngressStack (Kafka listener for workflow messages from the local Chip Ingress stack).
Recommendation: Use it in tests that need to listen for Chip Ingress stack messages.
*/
func StartChipIngressStack(t *testing.T, testLogger zerolog.Logger, testEnv *ttypes.TestEnvironment) (context.Context, <-chan proto.Message, <-chan error) {
	t.Helper()

	stack, err := NewChipIngressStack(framework.L, testEnv.TestConfig)
	require.NoError(t, err, "failed to create chip ingress stack instance")

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
	testLogger.Info().Dur("timeout", timeout).Msg("Starting Chip Ingress stack listener...")
	listenerCtx, cancelListener := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(func() {
		cancelListener()
		testLogger.Info().Msg("Chip Ingress stack listener stopped.")
	})

	msgChan, errChan := stack.SubscribeToChipIngressStackMessages(listenerCtx, messageTypes)

	// Fail fast if there is an error from the heartbeat validation subscription
	select {
	case err := <-errChan:
		require.NoError(t, err, "Chip Ingress stack subscription failed during initialization")
	default:
		// No immediate error, proceed
	}

	testLogger.Info().Msg("Chip Ingress stack listener ready")
	return listenerCtx, msgChan, errChan
}

/*
Asserts that a specific log message is received from the Chip Ingress stack within a timeout period.
Returns an error if found in error channel or timeouts if a log message is not received.
*/
func AssertChipIngressStackMessage(ctx context.Context, t *testing.T, expectedLog string, testLogger zerolog.Logger, messageChan <-chan proto.Message, kafkaErrChan <-chan error, timeout time.Duration) error {
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
					if strings.Contains(typedMsg.Msg, WorkflowEngineInitErrorLog) {
						foundErrorLog <- true
					}
				case *workflowevents.UserLogs:
					testLogger.Info().Msg("➡️ Chip Ingress stack message received in test. Asserting...")
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
		testLogger.Warn().Msg("chip ingress stack found engine initialization failure message! (may be expected in negative tests)")
		return errors.New("chip ingress stack message validation completed with error: found engine initialization failure message")
	case <-time.After(timeout):
		testLogger.Error().Str("expected_log", expectedLog).Msg("Timed out waiting for expected user log message")
		if receivedUserLogs > 0 {
			testLogger.Warn().Int("received_user_logs", receivedUserLogs).Msg("Received some UserLogs messages, but none matched expected log")
		} else {
			testLogger.Warn().Msg("Did not receive any UserLogs messages")
		}
		require.Failf(t, "Timed out waiting for the expected user log message (or error)", "Expected user log message: '%s' not found after %s", expectedLog, timeout.String())
	case err := <-kafkaErrChan:
		testLogger.Error().Err(err).Msg("Kafka listener encountered an error during execution. Ensure Chip Ingress stack is running and accessible.")
		require.Fail(t, "Kafka listener failed", err.Error())
	}
	return nil
}

//////////////////////////////
//      CRYPTO HELPERS      //
//////////////////////////////

// CreateAndFundAddressesEVM - creates and funds a specified number of new Ethereum addresses on a given chain.
func CreateAndFundAddressesEVM(t *testing.T, testLogger zerolog.Logger, numberOfAddressesToCreate int, amountToFund *big.Int, bcOutput blockchains.Blockchain) ([]common.Address, error) {
	t.Helper()
	return createAndFundAddresses(t, testLogger, numberOfAddressesToCreate, amountToFund, bcOutput, func() (common.Address, error) {
		addr, _, err := crecrypto.GenerateNewKeyPair()
		return addr, err
	})
}

// CreateAndFundAddressesSolana - creates and funds a specified number of new Solana addresses on a given chain.
func CreateAndFundAddressesSolana(t *testing.T, testLogger zerolog.Logger, numberOfAddressesToCreate int, amountToFund *big.Int, bcOutput blockchains.Blockchain) ([]solana.PublicKey, error) {
	t.Helper()
	return createAndFundAddresses(t, testLogger, numberOfAddressesToCreate, amountToFund, bcOutput, func() (solana.PublicKey, error) {
		key, err := solkey.New()
		return key.PublicKey(), err
	})
}

func createAndFundAddresses[T interface{ String() string }](t *testing.T, testLogger zerolog.Logger, numberOfAddressesToCreate int, amountToFund *big.Int, bcOutput blockchains.Blockchain, generateKey func() (T, error)) ([]T, error) {
	t.Helper()

	testLogger.Info().Msgf("Creating and funding %d addresses...", numberOfAddressesToCreate)
	var addrs []T

	for i := range numberOfAddressesToCreate {
		addr, addrErr := generateKey()
		require.NoError(t, addrErr, "failed to generate address")
		orderNum := i + 1
		testLogger.Info().Msgf("Generated address #%d: %s", orderNum, addr.String())

		testLogger.Info().Msgf("Funding address '%s' with amount of '%s' wei", addr.String(), amountToFund.String())
		if err := bcOutput.Fund(t.Context(), addr.String(), amountToFund.Uint64()); err != nil {
			return nil, err
		}

		addrs = append(addrs, addr)
	}

	return addrs, nil
}

//////////////////////////////
// WORKFLOW-RELATED HELPERS //
//////////////////////////////

// Generic WorkflowConfig interface for creation of different workflow configurations
// Register your workflow configuration types here
type WorkflowConfig interface {
	None |
		portypes.WorkflowConfig |
		AptosReadWorkflowConfig |
		StellarReadWorkflowConfig |
		StellarWriteWorkflowConfig |
		aptoswrite_config.Config |
		aptoswriteroundtrip_config.Config |
		crontypes.WorkflowConfig |
		HTTPWorkflowConfig |
		consensus_negative_config.Config |
		evmread_config.Config |
		logtrigger_config.Config |
		evmread_negative_config.Config |
		evmwrite_negative_config.Config |
		logtrigger_negative_config.Config |
		http_config.Config |
		httpaction_smoke_config.Config |
		httpaction_negative_config.Config |
		solwrite_config.Config |
		sollogtrigger_config.Config |
		vaultsecret_config.Config |
		solread_config.Config
}

// None represents an empty workflow configuration
// It is used to satisfy the workflowConfigFactory, avoiding workflow config creation
type None struct{}

type HTTPWorkflowConfig struct {
	AuthorizedKey common.Address `json:"authorizedKey"`
	URL           string         `json:"url"`
}

type AptosReadWorkflowConfig struct {
	ChainSelector    uint64 `yaml:"chainSelector"`
	WorkflowName     string `yaml:"workflowName"`
	ExpectedCoinName string `yaml:"expectedCoinName"`
}

// StellarReadWorkflowConfig mirrors the fields of the stellarread workflow's
// config.Config (system-tests/tests/smoke/cre/stellar/stellarread/config).
// StellarReadKind mirrors config.ReadKind (stellarread workflow) — the iota values must
// match, since the config is serialized to the workflow as an int over YAML.
type StellarReadKind int

const (
	StellarReadKindLatestLedger StellarReadKind = iota
	StellarReadKindReadContract
)

type StellarReadWorkflowConfig struct {
	ChainSelector     uint64 `yaml:"chainSelector"`
	WorkflowName      string `yaml:"workflowName"`
	MinLedgerSequence uint64 `yaml:"minLedgerSequence"`
	CronSchedule      string `yaml:"cronSchedule"`

	// --- ReadContract path ---
	ReadKind      StellarReadKind `yaml:"readKind"`
	SourceAccount string          `yaml:"sourceAccount"`
	// Cases is the batch of ReadContract invocations run and asserted in a single trigger
	// when ReadKind == StellarReadKindReadContract.
	Cases []StellarReadContractStep `yaml:"cases"`
}

// StellarReadContractStep mirrors the stellarread workflow's config.ReadContractStep
// (system-tests/tests/smoke/cre/stellar/stellarread/config) — one ReadContract invocation
// asserted within a batch run.
type StellarReadContractStep struct {
	Name                      string   `yaml:"name"`
	ContractID                string   `yaml:"contractID"`
	Function                  string   `yaml:"function"`
	ArgMode                   string   `yaml:"argMode"`
	ArgBool                   bool     `yaml:"argBool"`
	ArgU32A                   uint32   `yaml:"argU32A"`
	ArgU32B                   uint32   `yaml:"argU32B"`
	ArgString                 string   `yaml:"argString"`
	ArgBytesHex               string   `yaml:"argBytesHex"`
	ArgI64                    int64    `yaml:"argI64"`
	ArgSymbol                 string   `yaml:"argSymbol"`
	ArgVecU32                 []uint32 `yaml:"argVecU32"`
	ArgReceiverContractIDHex  string   `yaml:"argReceiverContractIDHex"`
	ArgWorkflowExecutionIDHex string   `yaml:"argWorkflowExecutionIDHex"`
	ArgReportIDHex            string   `yaml:"argReportIDHex"`
	ExpectedResult            string   `yaml:"expectedResult"`
	ExpectError               bool     `yaml:"expectError"`
}

// StellarWriteWorkflowConfig mirrors the fields of the stellarwrite workflow's
// config.Config (system-tests/tests/smoke/cre/stellar/stellarwrite/config).
type StellarWriteWorkflowConfig struct {
	ChainSelector      uint64 `yaml:"chainSelector"`
	WorkflowName       string `yaml:"workflowName"`
	ReceiverContractID string `yaml:"receiverContractID"`
	ReportPayloadHex   string `yaml:"reportPayloadHex"`
	RequiredSignatures int    `yaml:"requiredSignatures"`
	ExpectFailure      bool   `yaml:"expectFailure"`
}

// WorkflowRegistrationConfig holds configuration for workflow registration
type WorkflowRegistrationConfig struct {
	WorkflowName            string
	WorkflowLocation        string
	WorkflowTag             string
	ConfigFilePath          string
	CompressedWasmPath      string
	SecretsURL              string
	WorkflowRegistryAddr    common.Address
	WorkflowRegistryVersion *semver.Version
	ChainID                 uint64
	DonID                   uint64
	DonFamily               string
	ContainerTargetDir      string
	SethClient              *seth.Client
	Attributes              []byte
}

/*
Creates the necessary workflow artifacts based on WorkflowConfig:
 1. Configuration for a workflow (or no config if typed nil is passed for workflowConfig);
 2. Compiled and compressed workflow WASM file;
 3. Copies the workflow artifacts to the Docker containers.

It returns the paths to:
 1. the compressed WASM file;
 2. the workflow config file.
*/
func createWorkflowArtifacts[T WorkflowConfig](t *testing.T, testLogger zerolog.Logger, workflowName string, workflowDONs []*cre.Don, workflowConfig *T, workflowFileLocation, artifactDir string) (string, string) {
	t.Helper()

	workflowConfigFilePath := workflowConfigFactory(t, testLogger, workflowName, workflowConfig, artifactDir)
	compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflowToDir(t.Context(), workflowFileLocation, workflowName, artifactDir)
	require.NoError(t, compileErr, "failed to compile workflow '%s'", workflowFileLocation)
	testLogger.Info().Msg("Workflow compiled successfully.")

	// Copy workflow artifacts to Docker containers to use blockchain client running inside for workflow registration
	testLogger.Info().Msg("Copying workflow artifacts to Docker containers.")
	for _, don := range workflowDONs {
		copyErr := creworkflow.CopyArtifactsToDockerContainers(creworkflow.DefaultWorkflowTargetDir, ns.NodeNamePrefix(don.Name), compressedWorkflowWasmPath, workflowConfigFilePath)
		require.NoError(t, copyErr, "failed to copy workflow artifacts to docker containers")
	}
	testLogger.Info().Msg("Workflow artifacts successfully copied to the Docker containers.")

	return compressedWorkflowWasmPath, workflowConfigFilePath
}

/*
Creates the necessary workflow configuration based on a type registered in the WorkflowConfig interface
Pass `nil` to skip workflow config file creation.

Returns the path to the workflow config file.
*/
func workflowConfigFactory[T WorkflowConfig](t *testing.T, testLogger zerolog.Logger, workflowName string, workflowConfig *T, outputDir string) (filePath string) {
	t.Helper()

	var workflowConfigFilePath string

	// nil is an acceptable argument that allows skipping config file creation when it is not necessary
	if workflowConfig != nil {
		switch cfg := any(workflowConfig).(type) {
		case *None:
			workflowConfigFilePath = ""
			testLogger.Info().Msg("Workflow config file is not requested and will not be created.")

		case *portypes.WorkflowConfig:
			// Validate and format the feed ID (truncate to 16 bytes / 32 hex chars)
			cleanID := strings.TrimPrefix(cfg.FeedID, "0x")
			if len(cleanID) < 32 {
				require.NoError(t, fmt.Errorf("PoR feed ID must be at least 32 hex characters, got %d", len(cleanID)))
			}
			if len(cleanID) > 32 {
				cleanID = cleanID[:32]
			}
			cfg.FeedID = "0x" + cleanID
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create PoR workflow config file")
			testLogger.Info().Msg("PoR workflow config file created.")

		case *AptosReadWorkflowConfig:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create aptos read workflow config file")
			testLogger.Info().Msg("Aptos read workflow config file created.")

		case *StellarReadWorkflowConfig:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create stellar read workflow config file")
			testLogger.Info().Msg("Stellar read workflow config file created.")

		case *StellarWriteWorkflowConfig:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create stellar write workflow config file")
			testLogger.Info().Msg("Stellar write workflow config file created.")

		case *aptoswrite_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create aptos write workflow config file")
			testLogger.Info().Msg("Aptos write workflow config file created.")

		case *aptoswriteroundtrip_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create aptos write roundtrip workflow config file")
			testLogger.Info().Msg("Aptos write roundtrip workflow config file created.")

		case *HTTPWorkflowConfig:
			workflowCfgFilePath, configErr := createHTTPWorkflowConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create HTTP workflow config file")
			testLogger.Info().Msg("HTTP workflow config file created.")

		case *crontypes.WorkflowConfig:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create Cron workflow config file")
			testLogger.Info().Msg("Cron workflow config file created.")

		case *consensus_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create consensus workflow config file")
			testLogger.Info().Msg("Consensus workflow config file created.")

		case *evmread_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmread workflow config file")
			testLogger.Info().Msg("EVM Read workflow config file created.")

		case *logtrigger_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create logtrigger workflow config file")
			testLogger.Info().Msg("EVM LogTrigger workflow config file created.")

		case *evmread_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmread-negative workflow config file")
			testLogger.Info().Msg("EVM Read negative workflow config file created.")

		case *evmwrite_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create evmwrite-negative workflow config file")
			testLogger.Info().Msg("EVM Write negative workflow config file created.")

		case *logtrigger_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create logtrigger-negative workflow config file")
			testLogger.Info().Msg("EVM LogTrigger negative workflow config file created.")

		case *http_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create http-negative workflow config file")
			testLogger.Info().Msg("HTTP negative workflow config file created.")

		case *httpaction_smoke_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create httpaction smoke workflow config file")
			testLogger.Info().Msg("HTTP Action smoke workflow config file created.")

		case *httpaction_negative_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create httpaction negative workflow config file")
			testLogger.Info().Msg("HTTP Action negative workflow config file created.")
		case *solwrite_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create solwrite workflow config file")
			testLogger.Info().Msg("Solana write workflow config file created.")
		case *sollogtrigger_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create solana logtrigger workflow config file")
			testLogger.Info().Msg("Solana log trigger workflow config file created.")
		case *solread_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create solana read workflow config file")
			testLogger.Info().Msg("Solana read workflow config file created.")
		case *vaultsecret_config.Config:
			workflowCfgFilePath, configErr := CreateWorkflowYamlConfigFile(workflowName, cfg, outputDir)
			workflowConfigFilePath = workflowCfgFilePath
			require.NoError(t, configErr, "failed to create vaultsecret workflow config file")
			testLogger.Info().Msg("Vault secret workflow config file created.")
		default:
			require.NoError(t, fmt.Errorf("unsupported workflow config type: %T", cfg))
		}
	}
	return workflowConfigFilePath
}

func createHTTPWorkflowConfigFile(workflowName string, cfg *HTTPWorkflowConfig, outputDir string) (string, error) {
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
	configPath := filepath.Join(outputDir, configFileName)

	writeErr := os.WriteFile(configPath, configBytes, 0o644) //nolint:gosec // this is a test file
	if writeErr != nil {
		return "", errors.Wrap(writeErr, "failed to write HTTP workflow config file")
	}

	return configPath, nil
}

/*
Creates .yaml workflow configuration file and returns the absolute path to the created config file.
*/
func CreateWorkflowYamlConfigFile(workflowName string, workflowConfig any, outputDir string) (string, error) {
	// Write workflow config to a .yaml file
	configMarshalled, err := yaml.Marshal(workflowConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal workflow config")
	}
	if mkErr := os.MkdirAll(outputDir, 0o755); mkErr != nil {
		return "", errors.Wrap(mkErr, "failed to create output directory")
	}

	workflowConfigFile, tempErr := os.CreateTemp(outputDir, workflowName+"-*_config.yaml")
	if tempErr != nil {
		return "", errors.Wrap(tempErr, "failed to create workflow config file")
	}
	workflowConfigOutputFile := workflowConfigFile.Name()
	if closeErr := workflowConfigFile.Close(); closeErr != nil {
		return "", errors.Wrap(closeErr, "failed to close workflow config file")
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
) string {
	t.Helper()

	t.Cleanup(func() {
		deleteWorkflows(t, wfRegCfg.WorkflowName, wfRegCfg.ConfigFilePath,
			wfRegCfg.CompressedWasmPath,
			wfRegCfg.WorkflowRegistryAddr, wfRegCfg.WorkflowRegistryVersion, wfRegCfg.SethClient,
		)
	})

	workflowID, registerErr := registerWorkflowErr(ctx, wfRegCfg, sethClient)
	require.NoError(t, registerErr, "failed to register workflow '%s'", wfRegCfg.WorkflowName)
	testLogger.Info().Msgf("Workflow registered successfully: '%s'", workflowID)
	return workflowID
}

// registerWorkflowErr registers a workflow on-chain without test cleanup or require.
func registerWorkflowErr(ctx context.Context, wfRegCfg *WorkflowRegistrationConfig, sethClient *seth.Client) (string, error) {
	var configURL *string
	if wfRegCfg.ConfigFilePath != "" {
		u := "file://" + wfRegCfg.ConfigFilePath
		configURL = &u
	}
	binaryURL := "file://" + wfRegCfg.CompressedWasmPath
	containerTargetDir := &wfRegCfg.ContainerTargetDir

	return creworkflow.RegisterWithContract(
		ctx,
		sethClient,
		wfRegCfg.WorkflowRegistryAddr,
		wfRegCfg.WorkflowRegistryVersion,
		wfRegCfg.DonID,
		wfRegCfg.DonFamily,
		wfRegCfg.WorkflowName,
		wfRegCfg.WorkflowTag,
		binaryURL,
		configURL,
		nil, // no secrets yet
		wfRegCfg.Attributes,
		containerTargetDir,
	)
}

const ReentrancySentryOOGError = "ReentrancySentryOOG"

/*
Deletes workflows from:
 1. Local environment
 2. Workflow Registry

Recommendation:
Use it at the end of your test to `t.Cleanup()` the env after test run
*/
func deleteWorkflows(
	t *testing.T,
	uniqueWorkflowName string,
	workflowConfigFilePath string,
	compressedWorkflowWasmPath string,
	workflowRegistryAddress common.Address,
	version *semver.Version,
	sethClient *seth.Client,
) {
	t.Helper()

	testLogger := framework.L
	testLogger.Info().Msgf("Deleting workflow artifacts (%s) after test.", uniqueWorkflowName)
	localEnvErr := creworkflow.RemoveWorkflowArtifactsFromLocalEnv(workflowConfigFilePath, compressedWorkflowWasmPath)
	require.NoError(t, localEnvErr, "failed to remove workflow artifacts from local environment")

	deleteWorkflowsMu.Lock()
	defer deleteWorkflowsMu.Unlock()

	retryErr := retry.Do(func() error {
		return creworkflow.DeleteWithContract(t.Context(), sethClient, workflowRegistryAddress, version, uniqueWorkflowName)
	}, retry.Attempts(3), retry.Delay(1*time.Second), retry.DelayType(retry.BackOffDelay), retry.RetryIf(func(err error) bool {
		/**
		 * NOTE: ReentrancySentryOOG occurs if the EVM hits a gas cliff during the
		 * ReentrancyGuard cleanup phase (the 63/64ths rule).
		 * Intermittent failure in simulation is usually due to Cold vs Warm storage
		 * slot costs. Retrying might work, because the subsequent attempt might benefit from
		 * warmed storage/access lists, saving ~2,000 gas.
		 */
		return strings.Contains(err.Error(), ReentrancySentryOOGError)
	}), retry.OnRetry(func(n uint, err error) {
		testLogger.Error().Msgf("Error deleting workflow '%s': %s", uniqueWorkflowName, err.Error())
	}))

	if retryErr != nil && isWorkflowNotFound(retryErr) {
		// Eager delete (see DeleteWorkflowFromRegistry) may have already removed it; treat as success.
		testLogger.Info().Msgf("Workflow '%s' already absent from the registry (treating as deleted).", uniqueWorkflowName)
		return
	}
	require.NoError(t, retryErr, "failed to delete workflow '%s'", uniqueWorkflowName)
	testLogger.Info().Msgf("Workflow '%s' deleted successfully from the registry.", uniqueWorkflowName)
}

// isWorkflowNotFound reports whether err indicates the workflow is already absent from
// the registry, so a delete can be treated as an idempotent success.
func isWorkflowNotFound(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "not found")
}

func CompileAndDeployWorkflow[T WorkflowConfig](t *testing.T,
	testEnv *ttypes.TestEnvironment, testLogger zerolog.Logger, workflowName string,
	workflowConfig *T, workflowFileLocation string,
	opts ...CompileAndDeployWorkflowOpt,
) string {
	t.Helper()
	cfg := compileAndDeployWorkflowCfg{
		artifactCopyDONTypes: []cre.CapabilityFlag{cre.WorkflowDON},
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	testLogger.Info().
		Str("workflow_name", workflowName).
		Str("workflow_file_location", workflowFileLocation).
		Msgf("compiling and registering workflow '%s'", workflowName)
	artifactDir := workflowArtifactsDir(t, testEnv)
	registryChainSelector := testEnv.CreEnvironment.Blockchains[0].ChainSelector()

	workflowDONs := selectArtifactTargetDONs(testEnv, cfg.artifactCopyDONTypes)

	compressedWorkflowWasmPath, workflowConfigPath := createWorkflowArtifacts(t, testLogger, workflowName, workflowDONs, workflowConfig, workflowFileLocation, artifactDir)
	require.NotEmpty(t, compressedWorkflowWasmPath, "failed to find workflow DON in the topology")

	workflowRegistryAddress := crecontracts.MustGetAddressRefFromDataStore(testEnv.CreEnvironment.CldfEnvironment.DataStore, testEnv.CreEnvironment.Blockchains[0].ChainSelector(), keystone_changeset.WorkflowRegistry.String(), testEnv.CreEnvironment.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")

	workflowRegConfig := &WorkflowRegistrationConfig{
		WorkflowName:            workflowName,
		WorkflowLocation:        workflowFileLocation,
		WorkflowTag:             creworkflow.DefaultWorkflowTag,
		ConfigFilePath:          workflowConfigPath,
		CompressedWasmPath:      compressedWorkflowWasmPath,
		WorkflowRegistryAddr:    common.HexToAddress(workflowRegistryAddress.Address),
		WorkflowRegistryVersion: workflowRegistryAddress.Version,
		ChainID:                 registryChainSelector,
		DonID:                   testEnv.Dons.MustWorkflowDON().ID,
		DonFamily:               testEnv.Dons.MustWorkflowDON().DonFamily,
		ContainerTargetDir:      creworkflow.DefaultWorkflowTargetDir,
		SethClient:              testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient,
		Attributes:              cfg.attributes,
	}
	require.IsType(t, &evm.Blockchain{}, testEnv.CreEnvironment.Blockchains[0], "expected EVM blockchain type")
	workflowID := registerWorkflow(t.Context(), t, workflowRegConfig, testEnv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient, testLogger)
	return workflowID
}

func envVarOrDefault(envVar string, defaultValue int) int {
	v := strings.TrimSpace(os.Getenv(envVar))
	if v == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultValue
	}
	return n
}

type compileAndDeployWorkflowCfg struct {
	artifactCopyDONTypes []cre.CapabilityFlag
	attributes           []byte
}

// CompileAndDeployWorkflowOpt customizes workflow compilation/deployment behavior.
type CompileAndDeployWorkflowOpt func(*compileAndDeployWorkflowCfg)

// WithArtifactCopyDONTypes sets DON types where workflow artifacts should be copied.
func WithArtifactCopyDONTypes(donTypes ...cre.CapabilityFlag) CompileAndDeployWorkflowOpt {
	return func(cfg *compileAndDeployWorkflowCfg) {
		if len(donTypes) == 0 {
			return
		}
		cfg.artifactCopyDONTypes = append([]cre.CapabilityFlag{}, donTypes...)
	}
}

// WithAttributes sets the workflow attributes byte blob (JSON) written to the
// WorkflowRegistry contract on upsert. The CRE syncer reads this to decide
// routing (e.g. confidential execution via ConfidentialModule). The input is
// cloned so later caller mutations don't affect stored config.
func WithAttributes(attributes []byte) CompileAndDeployWorkflowOpt {
	return func(cfg *compileAndDeployWorkflowCfg) {
		cfg.attributes = slices.Clone(attributes)
	}
}

func selectArtifactTargetDONs(testEnv *ttypes.TestEnvironment, donTypes []cre.CapabilityFlag) []*cre.Don {
	if len(donTypes) == 0 {
		donTypes = []cre.CapabilityFlag{cre.WorkflowDON}
	}
	allow := make(map[cre.CapabilityFlag]struct{}, len(donTypes))
	for _, donType := range donTypes {
		allow[donType] = struct{}{}
	}

	targetDONs := make([]*cre.Don, 0)
	for _, don := range testEnv.Dons.List() {
		for donType := range allow {
			if don.HasFlag(donType) {
				targetDONs = append(targetDONs, don)
				break
			}
		}
	}
	return targetDONs
}

func workflowArtifactsDir(t *testing.T, testEnv *ttypes.TestEnvironment) string {
	t.Helper()
	if testEnv.Execution == nil || testEnv.Execution.TestID == "" {
		dir, err := os.MkdirTemp("", "cre-workflow-artifacts-*")
		require.NoError(t, err, "failed to create artifacts directory")
		return dir
	}

	dir := filepath.Join(os.TempDir(), "cre-workflow-artifacts", testEnv.Execution.TestID)
	require.NoError(t, os.MkdirAll(dir, 0o755), "failed to create artifacts directory")
	return dir
}

func UniqueWorkflowName(testEnv *ttypes.TestEnvironment, baseName string) string {
	testID := ""
	if testEnv != nil && testEnv.Execution != nil {
		testID = testEnv.Execution.TestID
	}
	if testID == "" {
		return truncateWorkflowName(baseName, baseName)
	}
	return truncateWorkflowName(fmt.Sprintf("%s-%s", baseName, testID), fmt.Sprintf("%s:%s", baseName, testID))
}

func truncateWorkflowName(name, uniquenessSeed string) string {
	if len(name) <= maxWorkflowNameLen {
		return name
	}

	sum := sha256.Sum256([]byte(uniquenessSeed))
	suffix := hex.EncodeToString(sum[:])[:8]
	prefixLen := maxWorkflowNameLen - len(suffix) - 1 // include hyphen
	if prefixLen < 1 {
		return suffix[:maxWorkflowNameLen]
	}
	if prefixLen > len(name) {
		prefixLen = len(name)
	}
	return fmt.Sprintf("%s-%s", name[:prefixLen], suffix)
}

func ParallelEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("CRE_TEST_PARALLEL_ENABLED")))
	return v == "1" || v == "true" || v == "yes"
}

// ─── Stellar ReadContract test helpers (shared by smoke + regression) ───
//
// These are the generic, chain-family-agnostic bits reused by both the smoke
// (smoke/cre) and regression (regression/cre) Stellar suites, which are separate
// Go packages. They live here — alongside StellarReadWorkflowConfig — for the same
// reason EVM keeps StartChipTestSink/GetEVMEnabledChains here: shared infra that a
// _test.go file can't export across packages. The per-suite case tables and the
// deploy-and-watch loop are inlined in each test file (matching evm/solana/aptos).

// StellarWorkflowTimeout bounds how long a Stellar read/read-contract case waits
// for its expected log.
const StellarWorkflowTimeout = 4 * time.Minute

// stellarCronScheduleEnvVar overrides the cron schedule for Stellar test workflows.
const stellarCronScheduleEnvVar = "CRE_STELLAR_CRON_SCHEDULE"

var stellarWorkflowNameSeq atomic.Uint64

// StellarTestCronSchedule returns the cron schedule Stellar test workflows should use.
// Empty (the default) keeps the workflow's safe "*/30 * * * * *". Set
// CRE_STELLAR_CRON_SCHEDULE (e.g. "*/5 * * * * *") to cut per-case idle wait — only
// where the Local CRE cron fastest-interval limit permits a faster schedule.
func StellarTestCronSchedule() string { return os.Getenv(stellarCronScheduleEnvVar) }

// UniqueStellarWorkflowName returns a re-run-safe workflow name.
func UniqueStellarWorkflowName(base string) string {
	return fmt.Sprintf("%s-%d-%d", base, time.Now().UnixNano(), stellarWorkflowNameSeq.Add(1))
}

// MustStellarChainInEnv returns the first Stellar chain in the environment, failing the test if none is present.
func MustStellarChainInEnv(t *testing.T, tenv *ttypes.TestEnvironment) *stellchain.Blockchain {
	t.Helper()
	require.NotNil(t, tenv, "Stellar suite requires a test environment")
	require.NotNil(t, tenv.CreEnvironment, "Stellar suite requires a CRE environment")
	require.NotEmpty(t, tenv.CreEnvironment.Blockchains, "Stellar suite expects at least one blockchain in the environment")
	for _, bc := range tenv.CreEnvironment.Blockchains {
		if bc.IsFamily(blockchain.FamilyStellar) {
			concrete, ok := bc.(*stellchain.Blockchain)
			require.True(t, ok, "expected concrete *stellar.Blockchain, got %T", bc)
			return concrete
		}
	}
	require.FailNow(t, "Stellar suite expects a Stellar chain in the environment (use config workflow-gateway-don-stellar.toml)")
	return nil
}

// MustDeployStellarReadFixture deploys the read_fixture contract and returns its
// C-address, failing the test on error.
func MustDeployStellarReadFixture(t *testing.T, stellarChain *stellchain.Blockchain) string {
	t.Helper()
	fixtureID, err := stellarfeature.DeployStellarReadFixture(context.Background(), stellarChain)
	require.NoError(t, err, "failed to deploy Stellar CRE read fixture")
	framework.L.Info().Str("fixture", fixtureID).Msg("Deployed Stellar ReadContract fixture")
	return fixtureID
}

// StartChipTestSinkWithLogging starts a Chip test sink that both writes every event to
// logFilePath (the ./logs dir is uploaded as a GH artifact) AND returns the user-log /
// base-message channels for the caller to consume via WatchWorkflowLogs / WaitForUserLog.
func StartChipTestSinkWithLogging(t *testing.T, logFilePath string) (<-chan *workflowevents.UserLogs, <-chan *commonevents.BaseMessage) {
	t.Helper()
	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := StartChipTestSink(t, GetLoggingPublishFn(framework.L, userLogsCh, baseMessageCh, logFilePath))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})
	return userLogsCh, baseMessageCh
}

// StartLoggingOnlyChipTestSink starts a Chip test sink that only writes events to logFilePath.
// It passes nil channels to the publish handler (safeSend* treats nil as a no-op), so events
// are never queued for consumption — there is nothing to drain and Publish never blocks.
// Use this for tests that assert on external state and only need the event dump for debugging.
func StartLoggingOnlyChipTestSink(t *testing.T, logFilePath string) {
	t.Helper()
	server := StartChipTestSink(t, GetLoggingPublishFn(framework.L, nil, nil, logFilePath))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ShutdownChipSinkWithDrain(ctx, server)
	})
}

func LogFilePath(prefix, suffix string) string {
	safe := SanitizeLogToken(suffix)
	if safe == "" {
		safe = "default"
	}
	return filepath.Join("logs", fmt.Sprintf("%s_%s.log", prefix, safe))
}

func SanitizeLogToken(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}

	return b.String()
}
