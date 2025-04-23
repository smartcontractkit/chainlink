package cre

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	pb2 "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type TestConfigLoadTestWriter struct {
	BlockchainA                   *blockchain.Input                        `toml:"blockchain_a" validate:"required"`
	NodeSets                      []*ns.Input                              `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                                `toml:"jd" validate:"required"`
	KeystoneContracts             *keystonetypes.KeystoneContractsInput    `toml:"keystone_contracts"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput     `toml:"workflow_registry_configuration"`
	DataFeedsCache                *keystonetypes.DeployDataFeedsCacheInput `toml:"feed_consumer"`
	Infra                         *libtypes.InfraInput                     `toml:"infra" validate:"required"`
	MockCapabilities              []*MockCapabilities                      `toml:"mock_capabilities"`
	BinariesConfig                *BinariesConfig                          `toml:"binaries_config"`
}

func setupLoadTestWriterEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfigLoadTest,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
	jobSpecFactoryFns []keystonetypes.JobSpecFactoryFn,
) *loadTestSetupOutput {
	absMockCapabilityBinaryPath, err := filepath.Abs(in.BinariesConfig.MockCapabilityBinaryPath)
	require.NoError(t, err, "failed to get absolute path for mock capability binary")

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     *in.BlockchainA,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  map[string]string{keystonetypes.MockCapability: absMockCapabilityBinaryPath},
		JobSpecFactoryFunctions:              jobSpecFactoryFns,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	// Set inputs in the test config, so that they can be saved
	in.KeystoneContracts = &keystonetypes.KeystoneContractsInput{}
	in.KeystoneContracts.Out = universalSetupOutput.KeystoneContractsOutput
	in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{}
	in.WorkflowRegistryConfiguration.Out = universalSetupOutput.WorkflowRegistryConfigurationOutput

	return &loadTestSetupOutput{
		// feedsConsumerAddress: deployFeedsConsumerOutput.FeedConsumerAddress,
		forwarderAddress: universalSetupOutput.KeystoneContractsOutput.ForwarderAddress,
		blockchainOutput: universalSetupOutput.BlockchainOutput.BlockchainOutput,
		donTopology:      universalSetupOutput.DonTopology,
		nodeOutput:       universalSetupOutput.NodeOutput,
	}
}

func TestLoad_Writer_MockCapabilities(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfigLoadTest](t)
	require.NoError(t, err, "couldn't load test config")
	require.Len(t, in.NodeSets, 1, "expected 1 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.WriteEVMCapability, keystonetypes.MockCapability, keystonetypes.OCR3Capability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON, keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
		}
	}

	containerPath, pathErr := crecapabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")

	loadTestJobSpecsFactoryFn := func(input *keystonetypes.JobSpecFactoryInput) (keystonetypes.DonsToJobSpecs, error) {
		donTojobSpecs := make(keystonetypes.DonsToJobSpecs, 0)

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			jobSpecs := make(keystonetypes.DonJobs, 0)
			workflowNodeSet, err2 := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &keystonetypes.Label{Key: node.NodeTypeKey, Value: keystonetypes.WorkerNode}, node.EqualLabels)
			if err2 != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err2, "failed to find worker nodes")
			}
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				if flags.HasFlag(donWithMetadata.Flags, keystonetypes.MockCapability) && in.MockCapabilities != nil {
					jobSpecs = append(jobSpecs, MockCapabilitiesJob(nodeID, filepath.Join(containerPath, filepath.Base(in.BinariesConfig.MockCapabilityBinaryPath)), in.MockCapabilities))
				}
			}

			donTojobSpecs[donWithMetadata.ID] = jobSpecs
		}

		return donTojobSpecs, nil
	}

	WriterDONLoadTestCapabilitiesFactoryFn := func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		if flags.HasFlag(donFlags, keystonetypes.MockCapability) {
			for _, m := range in.MockCapabilities {
				capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
					Capability: kcr.CapabilitiesRegistryCapability{
						LabelledName:   m.Name,
						Version:        m.Version,
						CapabilityType: capTypeToInt(m.Type),
					},
					Config: &capabilitiespb.CapabilityConfig{},
				})
			}
		}

		return capabilities
	}

	chainIDInt, chainErr := strconv.ParseUint(in.BlockchainA.ChainID, 10, 64)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupLoadTestWriterEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig{WriterDONLoadTestCapabilitiesFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))},
		[]keystonetypes.JobSpecFactoryFn{loadTestJobSpecsFactoryFn, consensus.ConsensusJobSpecFactoryFn(chainIDInt)},
	)

	ctx := t.Context()
	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, "n/a", "n/a", setupOutput.dataFeedsCacheAddress.Hex(), setupOutput.forwarderAddress.Hex())

			logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

			removeErr := os.RemoveAll(logDir)
			if removeErr != nil {
				testLogger.Error().Err(removeErr).Msg("failed to remove log directory")
				return
			}

			_, saveErr := framework.SaveContainerLogs(logDir)
			if saveErr != nil {
				testLogger.Error().Err(saveErr).Msg("failed to save container logs")
				return
			}

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
			for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range setupOutput.nodeOutput[i].Output.CLNodes {
					containerNames = append(containerNames, output.Node.ContainerName)
				}
				debugDons = append(debugDons, &keystonetypes.DebugDon{
					NodesMetadata:  donWithMetadata.NodesMetadata,
					Flags:          donWithMetadata.Flags,
					ContainerNames: containerNames,
				})
			}

			debugInput := keystonetypes.DebugInput{
				DebugDons:        debugDons,
				BlockchainOutput: setupOutput.blockchainOutput,
				InfraInput:       in.Infra,
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	// Get OCR2 keys needed to sign the reports
	kb := make([]ocr2key.KeyBundle, 0)
	for _, don := range setupOutput.donTopology.DonsWithMetadata {
		if flags.HasFlag(don.Flags, keystonetypes.MockCapability) {
			for i, n := range don.DON.Nodes {
				if i == 0 {
					continue // Skip bootstrap nodes
				}

				key, err2 := n.ExportOCR2Keys(n.Ocr2KeyBundleID)
				if err2 == nil {
					b, err3 := json.Marshal(key)
					require.NoError(t, err3, "could not marshal OCR2 key")
					kk, err3 := ocr2key.FromEncryptedJSON(b, nodeclient.ChainlinkKeyPassword)
					require.NoError(t, err3, "could not decrypt OCR2 key json")
					kb = append(kb, kk)
				} else {
					testLogger.Error().Msgf("Could not export OCR2 key: %s", err2)
				}
			}
		}
	}

	// Export key bundles so we can import them later in another test, used when crib cluster is already setup and we just want to connect to mocks for a different test
	require.NoError(t, saveKeyBundles(kb), "could not save OCR2 Keys")

	testLogger.Info().Msg("Connecting to mock capabilities...")

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.InfraType == "docker" {
		// TODO: For CTFv2 we should get the ports from the .toml
		// Need to add addresses manually
		mockClientsAddress = []string{"127.0.0.1:13401", "127.0.0.1:13402", "127.0.0.1:13403", "127.0.0.1:13404"}
	} else {
		for i := range setupOutput.nodeOutput[0].CLNodes {
			// TODO: This is brittle, switch to checking the node label
			if i == 0 { // Skip bootstrap node
				continue
			}
			mockClientsAddress = append(mockClientsAddress, fmt.Sprintf("%s-%s-%d-mock.main.stage.cldev.sh:443", in.Infra.CRIB.Namespace, setupOutput.nodeOutput[1].NodeSetName, i-1))
		}
	}

	// Use insecure gRPC connection for local Docker containers. For AWS, use TLS credentials
	// due to ingress requirements, as grpc.insecure.NewCredentials() doesn't work properly with AWS ingress
	useInsecure := false
	if in.Infra.InfraType == "docker" {
		useInsecure = true
	}

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, useInsecure, true), "could not connect to mock capabilities")

	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel), "could not hook into mock executable")

	labels := map[string]string{
		"go_test_name": "test1",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	//TODO: @george-dorin wait for cap to be exposed
	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5, // Give enough time for the workflow to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 120*time.Minute),
			),
			Gun:                   NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", receiveChannel),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "wasp load test did not finish successfully")
}

// TestWriteWithReconnect Re-runs the load test against an existing DON deployment. It expects feeds, OCR2 keys, and
// mock addresses to be cached from a previous test run. This is useful for tweaking load patterns or debugging
// workflow execution without redeploying the entire test environment.
func TestWriteWithReconnect(t *testing.T) {
	testLogger := framework.L
	ctx := t.Context()

	kb, err := loadKeyBundlesFromCache()
	require.NoError(t, err, "could not load OCR2 keys")

	testLogger.Info().Msg("Connecting to mock capabilities...")
	var mocksClient *mock_capability.Controller

	mocksClient, err = mock_capability.NewMockCapabilityControllerFromCache(testLogger, false)
	require.NoError(t, err, "could not create mock controller")

	caps, err := mocksClient.List(ctx)
	require.NoError(t, err, "error getting the nodes capabilities")
	fmt.Println(caps)
	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel), "could not hook into executable")

	labels := map[string]string{
		"go_test_name": "Workflow DON Load Test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	sg := NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", receiveChannel)
	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5, // Give enough time for the workflow to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(10, 15*time.Minute),
			),
			Gun:                   sg,
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "wasp load test did not finish successfully")
}

var _ wasp.Gun = (*WriterGun)(nil)

type WriterGun struct {
	capProxy    *mock_capability.Controller
	keyBundles  []ocr2key.KeyBundle
	triggerID   string
	waitChans   map[int64]chan interface{}
	receiveChan <-chan capabilities.CapabilityRequest
	mu          sync.Mutex
	event       *capabilities.OCRTriggerEvent
	eventID     string
}

func NewWriterGun(capProxy *mock_capability.Controller, keyBundles []ocr2key.KeyBundle, triggerID string, ch <-chan capabilities.CapabilityRequest) *WriterGun {
	sg := &WriterGun{
		capProxy:    capProxy,
		keyBundles:  keyBundles,
		triggerID:   triggerID,
		receiveChan: ch,
	}
	return sg
}

func (s *WriterGun) Call(l *wasp.Generator) *wasp.Response {
	reportID := [2]byte{0x00, 0x01}
	reportMetadata := evm.ReportV1Metadata{
		Version:             1,
		WorkflowExecutionID: [32]byte{},
		Timestamp:           0,
		DonID:               0,
		DonConfigVersion:    0,
		WorkflowCID:         [32]byte{},
		WorkflowName:        [10]byte{},
		WorkflowOwner:       [20]byte{},
		ReportID:            reportID,
	}

	reportMetadataBytes, err := reportMetadata.Encode()
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	signatures := [][]byte{}

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     reportMetadataBytes,
			"signatures": signatures,
			"context":    []byte{4, 5},
			"id":         reportID[:],
		},
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	validMetadata := capabilities.RequestMetadata{
		WorkflowID:          hex.EncodeToString(reportMetadata.WorkflowCID[:]),
		WorkflowOwner:       hex.EncodeToString(reportMetadata.WorkflowOwner[:]),
		WorkflowName:        hex.EncodeToString(reportMetadata.WorkflowName[:]),
		WorkflowExecutionID: hex.EncodeToString(reportMetadata.WorkflowExecutionID[:]),
	}

	validConfig, err := values.NewMap(map[string]any{
		"Address": testutils.NewAddress().String(),
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	//gasLimitDefault := uint64(400_000)

	configBytes, err := mock_capability.MapToBytes(validConfig)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	inputBytes, err := mock_capability.MapToBytes(validInputs)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	req := &pb2.ExecutableRequest{
		ID:             s.triggerID,
		CapabilityType: 4,
		RequestMetadata: &pb2.Metadata{
			WorkflowID:               validMetadata.WorkflowID,
			WorkflowOwner:            validMetadata.WorkflowOwner,
			WorkflowExecutionID:      validMetadata.WorkflowExecutionID,
			WorkflowName:             validMetadata.WorkflowName,
			WorkflowDonID:            validMetadata.WorkflowDonID,
			WorkflowDonConfigVersion: validMetadata.WorkflowDonConfigVersion,
			ReferenceID:              validMetadata.ReferenceID,
			DecodedWorkflowName:      validMetadata.DecodedWorkflowName,
		},
		Config: configBytes,
		Inputs: inputBytes,
	}

	err = s.capProxy.Execute(context.TODO(), req)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	return &wasp.Response{}
}
