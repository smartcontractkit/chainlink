package cre

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
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
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

type TestConfigLoadTestWriter struct {
	Blockchains                   []*blockchain.Input                  `toml:"blockchains" validate:"required"`
	NodeSets                      []*ns.Input                          `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                            `toml:"jd" validate:"required"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *libtypes.InfraInput                 `toml:"infra" validate:"required"`
	MockCapabilities              []*MockCapabilities                  `toml:"mock_capabilities"`
	BinariesConfig                *BinariesConfig                      `toml:"binaries_config"`
}

func setupLoadTestWriterEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfigLoadTestWriter,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
	jobSpecFactoryFns []keystonetypes.JobSpecFactoryFn,
) *loadTestSetupOutput {
	absMockCapabilityBinaryPath, err := filepath.Abs(in.BinariesConfig.MockCapabilityBinaryPath)
	require.NoError(t, err, "failed to get absolute path for mock capability binary")

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  map[string]string{keystonetypes.MockCapability: absMockCapabilityBinaryPath},
		JobSpecFactoryFunctions:              jobSpecFactoryFns,
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")
	// Set inputs in the test config, so that they can be saved
	in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{}
	in.WorkflowRegistryConfiguration.Out = universalSetupOutput.WorkflowRegistryConfigurationOutput

	forwarderAddress, forwarderErr := libcontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainSelector)

	return &loadTestSetupOutput{
		//feedsConsumerAddress: deployFeedsConsumerOutput.FeedConsumerAddress,
		forwarderAddress: forwarderAddress,
		blockchainOutput: universalSetupOutput.BlockchainOutput,
		donTopology:      universalSetupOutput.DonTopology,
		nodeOutput:       universalSetupOutput.NodeOutput,
	}
}

func TestLoad_Writer_MockCapabilities(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfigLoadTestWriter](t)
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

		if flags.HasFlag(donFlags, keystonetypes.OCR3Capability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "offchain_reporting",
					Version:        "1.0.0",
					CapabilityType: 2, // CONSENSUS
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

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

	homeChain := in.Blockchains[0]
	homeChainIDUint64, homeChainErr := strconv.ParseUint(homeChain.ChainID, 10, 64)
	require.NoError(t, homeChainErr, "failed to convert chain ID to int")

	setupOutput := setupLoadTestWriterEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig{WriterDONLoadTestCapabilitiesFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(homeChainIDUint64)))},
		[]keystonetypes.JobSpecFactoryFn{loadTestJobSpecsFactoryFn, consensus.ConsensusJobSpecFactoryFn(homeChainIDUint64)},
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
				BlockchainOutput: setupOutput.blockchainOutput[0].BlockchainOutput,
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

	require.NoError(t, saveForwarderAddress(setupOutput.dataFeedsCacheAddress.String()), "could not save forwarder address")

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
			mockClientsAddress = append(mockClientsAddress, fmt.Sprintf("%s-%s-%d-mock.main.stage.cldev.sh:443", in.Infra.CRIB.Namespace, setupOutput.nodeOutput[0].NodeSetName, i-1))
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

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5, // Give enough time for the workflow to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 120*time.Minute),
			),
			Gun:                   NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", receiveChannel, setupOutput.dataFeedsCacheAddress.String(), logger.TestLogger(t)),
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

	forwarderAddress, err := loadForwarderAddressFromCache()
	require.NoError(t, err, "could not load forwarder address from cache")

	kb, err := loadKeyBundlesFromCache()
	require.NoError(t, err, "could not load OCR2 keys")

	testLogger.Info().Msg("Connecting to mock capabilities...")
	var mocksClient *mock_capability.Controller

	mocksClient, err = mock_capability.NewMockCapabilityControllerFromCache(testLogger, false)
	require.NoError(t, err, "could not create mock controller")

	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel), "could not hook into executable")

	labels := map[string]string{
		"go_test_name": "Workflow DON Load Test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	sg := NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", receiveChannel, forwarderAddress, logger.TestLogger(t))
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
	capProxy              *mock_capability.Controller
	keyBundles            []ocr2key.KeyBundle
	triggerID             string
	waitChans             map[int64]chan interface{}
	receiveChan           <-chan capabilities.CapabilityRequest
	mu                    sync.Mutex
	event                 *capabilities.OCRTriggerEvent
	eventID               string
	dataFeedsCacheAddress string
	lggr                  logger.Logger
}

func NewWriterGun(capProxy *mock_capability.Controller, keyBundles []ocr2key.KeyBundle, triggerID string, ch <-chan capabilities.CapabilityRequest, dataFeedsCacheAddress string, lggr logger.Logger) *WriterGun {
	sg := &WriterGun{
		capProxy:              capProxy,
		keyBundles:            keyBundles,
		triggerID:             triggerID,
		receiveChan:           ch,
		dataFeedsCacheAddress: dataFeedsCacheAddress,
		lggr:                  lggr,
	}
	return sg
}

func createReports(kb []ocr2key.KeyBundle) ([]*datastreams.FeedReport, error) {
	//Create report
	feeds := make([]FeedWithStreamID, 0)
	for i := range 10 {
		buf := [32]byte{}
		_, err := crand.Read(buf[:])
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, FeedWithStreamID{
			Feed:     "0x" + hex.EncodeToString(buf[:]),
			StreamID: int32(i),
		})
	}

	reports := []*datastreams.FeedReport{}
	for _, f := range feeds {
		price := decimal.NewFromInt(int64(rand.IntN(100)))
		reports = append(reports, createFeedReportDataStreams(price.BigInt(), time.Now().Unix(), f.Feed, kb))
	}

	return reports, nil
}

func (s *WriterGun) Call(l *wasp.Generator) *wasp.Response {
	reports, err := createReports(s.keyBundles)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	reportID := [2]byte{0x00, 0x01}
	wrappedReports, err := wrapReports(reports, 12, datastreams.Metadata{})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	executionID, err := workflows.EncodeExecutionID("5dbe5f217ff07d6b1dddb43519fe7bf13ccb10b540578fafdbea86b508abbd71", "")
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	metadata := pb2.Metadata{
		WorkflowID:               "5dbe5f217ff07d6b1dddb43519fe7bf13ccb10b540578fafdbea86b508abbd71",
		WorkflowOwner:            "0100000000000000000000000000000000000001",
		WorkflowExecutionID:      executionID,
		WorkflowName:             "61626364656630313233",
		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              "write_geth-testnet@1.0.0",
		DecodedWorkflowName:      "abcdef0123",
	}

	rMetadata := evm.ReportV1Metadata{
		Version:             1,
		WorkflowExecutionID: stringTo32Byte(metadata.WorkflowExecutionID),
		Timestamp:           0,
		DonID:               metadata.WorkflowDonID,
		DonConfigVersion:    metadata.WorkflowDonConfigVersion,
		WorkflowCID:         stringTo32Byte(metadata.WorkflowID),
		WorkflowName:        stringTo10Byte(metadata.WorkflowName),
		WorkflowOwner:       stringTo20Byte(metadata.WorkflowOwner),
		ReportID:            [2]byte{0, 1},
	}

	rMetaBytes, err := rMetadata.Encode()
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	reportBytes, err := mock_capability.MapToBytes(wrappedReports)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     append(rMetaBytes, reportBytes...),
			"signatures": reports[0].Signatures,
			"context":    []byte{4, 5},
			"id":         reportID[:],
		},
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	validConfig, err := values.NewMap(map[string]any{
		"Address":    s.dataFeedsCacheAddress,
		"abi":        "(bytes32 FeedID, bytes RawReport)[] Reports",
		"schedule":   "oneAtATime",
		"params":     []string{"$(report)"},
		"gasLimit":   400000,
		"deltaStage": "2s",
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	configBytes, err := mock_capability.MapToBytes(validConfig)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	inputBytes, err := mock_capability.MapToBytes(validInputs)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	req := &pb2.ExecutableRequest{
		ID:              s.triggerID,
		CapabilityType:  4, //Target
		RequestMetadata: &metadata,
		Config:          configBytes,
		Inputs:          inputBytes,
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
func saveForwarderAddress(forwarderAddress string) error {
	cacheDir := "./cache"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	forwarderPath := filepath.Join(cacheDir, "forwarder_address.txt")
	if err := os.WriteFile(forwarderPath, []byte(forwarderAddress), 0644); err != nil {
		return fmt.Errorf("failed to save forwarder address: %w", err)
	}

	return nil
}
func loadForwarderAddressFromCache() (string, error) {
	forwarderPath := filepath.Join("./cache", "forwarder_address.txt")
	data, err := os.ReadFile(forwarderPath)
	if err != nil {
		return "", fmt.Errorf("failed to read df cache address from cache: %w", err)
	}
	return string(data), nil
}
func stringTo32Byte(input string) [32]byte {
	var result [32]byte
	decoded, err := hex.DecodeString(input)
	if err != nil {
		panic(fmt.Sprintf("invalid string for conversion to [32]byte: %v", err))
	}
	copy(result[:], decoded)
	return result
}

func stringTo10Byte(input string) [10]byte {
	var result [10]byte
	decoded, err := hex.DecodeString(input)
	if err != nil {
		panic(fmt.Sprintf("invalid string for conversion to [32]byte: %v", err))
	}
	copy(result[:], decoded)
	return result
}

func stringTo20Byte(input string) [20]byte {
	var result [20]byte
	decoded, err := hex.DecodeString(input)
	if err != nil {
		panic(fmt.Sprintf("invalid string for conversion to [32]byte: %v", err))
	}
	copy(result[:], decoded)
	return result
}

func createFeedReportDataStreams(price *big.Int, observationTimestamp int64,
	feedIDString string,
	keyBundles []ocr2key.KeyBundle) *datastreams.FeedReport {
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	bytes, err := hex.DecodeString(feedIDString[2:])
	if err != nil {
		panic(err)
	}
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)

	report := &datastreams.FeedReport{
		FeedID:               feedIDString,
		FullReport:           newReport(feedIDBytes, price, observationTimestamp),
		BenchmarkPrice:       price.Bytes(),
		ObservationTimestamp: observationTimestamp,
		Signatures:           [][]byte{},
		ReportContext:        rawCtx,
	}

	for _, key := range keyBundles {
		sig, err := key.Sign(reportCtx, report.FullReport)
		if err != nil {
			panic(err)
		}
		report.Signatures = append(report.Signatures, sig)
		if len(report.Signatures) == 2 {
			break
		}
	}

	return report
}
func newReport(feedID [32]byte, price *big.Int, timestamp int64) []byte {
	v3Codec := reportcodec.NewReportCodec(feedID, logger.NullLogger)
	raw, err := v3Codec.BuildReport(context.TODO(), v3.ReportFields{
		BenchmarkPrice: price,

		Timestamp: uint32(timestamp), //nolint:gosec // G115
		Bid:       big.NewInt(0),
		Ask:       big.NewInt(0),
		LinkFee:   big.NewInt(0),
		NativeFee: big.NewInt(0),
	})
	if err != nil {
		panic(err)
	}
	return raw
}
func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func wrapReports(reportList []*datastreams.FeedReport,
	timestamp int64, meta datastreams.Metadata) (*values.Map, error) {
	rl := make([]datastreams.FeedReport, 0, len(reportList))
	for _, r := range reportList {
		rl = append(rl, *r)
	}

	return values.WrapMap(datastreams.StreamsTriggerEvent{
		Payload:   rl,
		Metadata:  meta,
		Timestamp: timestamp,
	})
}
