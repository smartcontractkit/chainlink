package cre

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type Chaos struct {
	Mode                        string   `toml:"mode"`
	Latency                     string   `toml:"latency"`
	Jitter                      string   `toml:"jitter"`
	DashboardUIDs               []string `toml:"dashboard_uids"`
	WaitBeforeStart             string   `toml:"wait_before_start"`
	ExperimentFullInterval      string   `toml:"experiment_full_interval"`
	ExperimentInjectionInterval string   `toml:"experiment_injection_interval"`
}

type TestConfigLoadTest struct {
	Duration                      string                               `toml:"duration"`
	Blockchains                   []*blockchain.Input                  `toml:"blockchains" validate:"required"`
	NodeSets                      []*ns.Input                          `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                            `toml:"jd" validate:"required"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *libtypes.InfraInput                 `toml:"infra" validate:"required"`
	WorkflowDONLoad               *WorkflowLoad                        `toml:"workflow_load"`
	MockCapabilities              []*MockCapabilities                  `toml:"mock_capabilities"`
	BinariesConfig                *BinariesConfig                      `toml:"binaries_config"`
	Chaos                         *Chaos                               `toml:"chaos"`
}

type BinariesConfig struct {
	MockCapabilityBinaryPath string `toml:"mock_capability_binary_path"`
}

type MockCapabilities struct {
	Name        string `toml:"name"`
	Version     string `toml:"version"`
	Type        string `toml:"type"`
	Description string `toml:"description"`
}

type WorkflowLoad struct {
	Streams       int32 `toml:"streams" validate:"required"`
	Jobs          int32 `toml:"jobs" validate:"required"`
	FeedAddresses [][]string
}

type FeedWithStreamID struct {
	Feed     string `json:"feed"`
	StreamID int32  `json:"streamID"`
}

type loadTestSetupOutput struct {
	dataFeedsCacheAddress common.Address
	forwarderAddress      common.Address
	blockchainOutput      []*creenv.BlockchainOutput
	donTopology           *keystonetypes.DonTopology
	nodeOutput            []*keystonetypes.WrappedNodeOutput
}

func setupLoadTestEnvironment(
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

	forwarderAddress, forwarderErr := crecontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainSelector)

	return &loadTestSetupOutput{
		forwarderAddress: forwarderAddress,
		blockchainOutput: universalSetupOutput.BlockchainOutput,
		donTopology:      universalSetupOutput.DonTopology,
		nodeOutput:       universalSetupOutput.NodeOutput,
	}
}

func TestLoad_Workflow_Streams_MockCapabilities(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfigLoadTest](t)
	require.NoError(t, err, "couldn't load test config")
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.OCR3Capability},
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.MockCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: 0,
			},
		}
	}

	feedsAddresses := make([][]FeedWithStreamID, in.WorkflowDONLoad.Jobs)
	for i := range in.WorkflowDONLoad.Jobs {
		feedsAddresses[i] = make([]FeedWithStreamID, 0)
		for streamID := range in.WorkflowDONLoad.Streams {
			_, id := NewFeedID(t)
			feedsAddresses[i] = append(feedsAddresses[i], FeedWithStreamID{
				Feed:     id,
				StreamID: (in.WorkflowDONLoad.Streams * i) + streamID + 1,
			})
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
				if flags.HasFlag(donWithMetadata.Flags, keystonetypes.WorkflowDON) {
					for i := range feedsAddresses {
						feedConfig := make([]FeedConfig, 0)

						for _, feed := range feedsAddresses[i] {
							feedID, err2 := datastreams.NewFeedID(feed.Feed)
							if err2 != nil {
								return nil, err2
							}
							feedBytes := feedID.Bytes()
							feedConfig = append(feedConfig, FeedConfig{
								FeedIDsIndex: feed.StreamID,
								Deviation:    "0.001",
								Heartbeat:    3600,
								RemappedID:   "0x" + hex.EncodeToString(feedBytes[:]),
							})
						}

						jobSpecs = append(jobSpecs, WorkflowsJob(nodeID, fmt.Sprintf("load_%d", i), feedConfig))
					}
				}

				if flags.HasFlag(donWithMetadata.Flags, keystonetypes.MockCapability) && in.MockCapabilities != nil {
					jobSpecs = append(jobSpecs, MockCapabilitiesJob(nodeID, filepath.Join(containerPath, filepath.Base(in.BinariesConfig.MockCapabilityBinaryPath)), in.MockCapabilities))
				}
			}

			donTojobSpecs[donWithMetadata.ID] = jobSpecs
		}

		return donTojobSpecs, nil
	}

	WorkflowDONLoadTestCapabilitiesFactoryFn := func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig {
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

		if flags.HasFlag(donFlags, keystonetypes.CustomComputeCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "custom-compute",
					Version:        "1.0.0",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

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

		return capabilities
	}

	homeChain := in.Blockchains[0]
	homeChainIDUint64, homeChainErr := strconv.ParseUint(homeChain.ChainID, 10, 64)
	require.NoError(t, homeChainErr, "failed to convert chain ID to int")

	setupOutput := setupLoadTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig{WorkflowDONLoadTestCapabilitiesFactoryFn, libcontracts.ChainWriterCapabilityFactory(homeChainIDUint64)},
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

	require.NoError(t, saveFeedAddresses(feedsAddresses), "could not save feeds")

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
		for i := range setupOutput.nodeOutput[1].CLNodes {
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

	g, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5, // Give enough time for the workflow to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 120*time.Minute),
			),
			Gun:                   NewStreamsGun(mocksClient, kb, feedsAddresses, "streams-trigger@2.0.0", receiveChannel, 500, 1),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(false)
	require.NoError(t, err, "wasp load test did not finish successfully")

	runChaosSuite(t, in)
	g.Wait()
}

// TestWithReconnect Re-runs the load test against an existing DON deployment. It expects feeds, OCR2 keys, and
// mock addresses to be cached from a previous test run. This is useful for tweaking load patterns or debugging
// workflow execution without redeploying the entire test environment.
func TestWithReconnect(t *testing.T) {
	testLogger := framework.L
	ctx := t.Context()

	kb, err := loadKeyBundlesFromCache()
	require.NoError(t, err, "could not load OCR2 keys")

	feedAddresses, err := loadFeedAddressesFromCache()
	require.NoError(t, err, "could not load feed addresses")
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

	sg := NewStreamsGun(mocksClient, kb, feedAddresses, "streams-trigger@2.0.0", receiveChannel, 600, 2)
	time.Sleep(time.Second * 5) // Give time for the report to be generated
	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5, // Give enough time for the workflow to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 15*time.Minute),
			),
			Gun:                   sg,
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err, "wasp load test did not finish successfully")
}

var _ wasp.Gun = (*StreamsGun)(nil)

type StreamsGun struct {
	capProxy    *mock_capability.Controller
	keyBundles  []ocr2key.KeyBundle
	feeds       [][]FeedWithStreamID
	triggerID   string
	waitChans   map[int64]chan interface{}
	receiveChan <-chan capabilities.CapabilityRequest
	mu          sync.Mutex
	feedLimit   int
	jobLimit    int
	event       *capabilities.OCRTriggerEvent
	eventID     string
	timestamp   time.Time
}

func NewStreamsGun(capProxy *mock_capability.Controller, keyBundles []ocr2key.KeyBundle, feeds [][]FeedWithStreamID, triggerID string, ch <-chan capabilities.CapabilityRequest, feedLimit int, jobLimit int) *StreamsGun {
	sg := &StreamsGun{
		capProxy:    capProxy,
		keyBundles:  keyBundles,
		feeds:       feeds,
		triggerID:   triggerID,
		receiveChan: ch,
		feedLimit:   feedLimit,
		jobLimit:    jobLimit,
	}
	go sg.precomputeReports()
	go sg.waitHOOKloop()
	return sg
}

func (s *StreamsGun) Call(l *wasp.Generator) *wasp.Response {
	workingTimestamp := s.timestamp.Unix()
	err := s.prepareWaitHOOK(workingTimestamp)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	payload, err := s.event.ToMap()
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	payloadBytes, err := mock_capability.MapToBytes(payload)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	err = s.capProxy.SendTrigger(l.ResponsesCtx, s.triggerID, s.eventID, payloadBytes)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	go s.precomputeReports()
	// Wait for the DON to execute on the write target
	err = s.waitForHOOK(workingTimestamp)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	return &wasp.Response{}
}

func (s *StreamsGun) waitHOOKloop() {
	for {
		m, ok := <-s.receiveChan
		if !ok {
			framework.L.Error().Msg("channel closed")
			return
		}

		inputs, err := decodeTargetInput(m.Inputs)
		if err != nil {
			framework.L.Error().Msg("error decoding inputs")
			return
		}

		// To get the timestamp we look at the last 64 chars of the hex encoded report
		hexReport := hex.EncodeToString(inputs.Inputs.SignedReport.Report)
		timestampInHex := hexReport[len(hexReport)-64:]
		timestamp, err := strconv.ParseInt(timestampInHex, 16, 64)
		if err != nil {
			framework.L.Error().Msg("error parsing timestamp")
			return
		}

		s.mu.Lock()
		// Check if exist
		if ch, exist := s.waitChans[timestamp]; exist {
			s.mu.Unlock()
			ch <- m // This is blocking
		} else {
			s.mu.Unlock()
		}
	}
}

func (s *StreamsGun) prepareWaitHOOK(reportTimestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitChans == nil {
		s.waitChans = make(map[int64]chan interface{})
	}
	if _, exists := s.waitChans[reportTimestamp]; exists {
		return fmt.Errorf("cannot prepare for HOOK, timestamp  %d already exits", reportTimestamp)
	}
	s.waitChans[reportTimestamp] = make(chan interface{})
	return nil
}

func (s *StreamsGun) waitForHOOK(timestamp int64) error {
	s.mu.Lock()
	ch, exists := s.waitChans[timestamp]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("cannot wait for HOOK, timestamp  %q does not exist", timestamp)
	}
	s.mu.Unlock()
	<-ch
	delete(s.waitChans, timestamp)
	return nil
}

func (s *StreamsGun) precomputeReports() {
	timestamp := time.Now()
	start := time.Now()

	price := decimal.NewFromInt(int64(rand.IntN(100)))

	feeds := make([]FeedWithStreamID, 0)
	for jobNr := range s.feeds {
		if jobNr >= s.jobLimit {
			break
		}

		for feedNr, feed := range s.feeds[jobNr] {
			if feedNr >= s.feedLimit {
				break
			}
			feeds = append(feeds, feed)
		}
	}

	event, eventID, err := createFeedReport(logger.NullLogger, price, uint64(timestamp.UnixNano()), feeds, s.keyBundles) //nolint:gosec // G115 don't care in test code
	if err != nil {
		panic(err)
	}

	s.event = event
	s.eventID = eventID
	s.timestamp = timestamp

	framework.L.Info().Msgf("precomputeReports took %s", time.Since(start))
}

func createFeedReport(lggr logger.Logger, price decimal.Decimal, timestamp uint64,
	feeds []FeedWithStreamID, keyBundles []ocr2key.KeyBundle) (*capabilities.OCRTriggerEvent, string, error) {
	values := make([]datastreamsllo.StreamValue, 0)

	priceBytes, err := price.MarshalBinary()
	if err != nil {
		return nil, "", err
	}
	streams := make([]llotypes.Stream, 0)

	for _, f := range feeds {
		dec := &datastreamsllo.Decimal{}
		err2 := dec.UnmarshalBinary(priceBytes)
		if err2 != nil {
			return nil, "", err2
		}
		values = append(values, dec)
		streams = append(streams, llotypes.Stream{
			StreamID: llotypes.StreamID(f.StreamID), //nolint:gosec // G115 don't care in test code
		})
	}

	reportCodec := cre.NewReportCodecCapabilityTrigger(lggr, 1)

	report := datastreamsllo.Report{
		ObservationTimestampNanoseconds: timestamp,
		Values:                          values,
	}

	reportBytes, err := reportCodec.Encode(report, llotypes.ChannelDefinition{
		Streams: streams,
	})
	if err != nil {
		return nil, "", err
	}
	eventID := reportCodec.EventID(report)

	event := &capabilities.OCRTriggerEvent{
		ConfigDigest: []byte{0: 1, 31: 2},
		SeqNr:        0,
		Report:       reportBytes,
		Sigs:         make([]capabilities.OCRAttributedOnchainSignature, 0, len(keyBundles)),
	}

	for i, key := range keyBundles {
		sig, err2 := key.Sign3(ocrTypes.ConfigDigest(event.ConfigDigest), event.SeqNr, reportBytes)
		if err2 != nil {
			return nil, "", err
		}
		event.Sigs = append(event.Sigs, capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(i), //nolint:gosec // G115 don't care in test code
			Signature: sig,
		})
	}

	return event, eventID, nil
}

func decodeTargetInput(inputs *values.Map) (targets.Request, error) {
	var r targets.Request
	const signedReportField = "signed_report"
	signedReport, ok := inputs.Underlying[signedReportField]
	if !ok {
		return r, fmt.Errorf("missing required field %s", signedReportField)
	}

	if err := signedReport.UnwrapTo(&r.Inputs.SignedReport); err != nil {
		return r, err
	}

	return r, nil
}

func saveKeyBundles(keyBundles []ocr2key.KeyBundle) error {
	cacheDir := "cache/keys"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	for i, kb := range keyBundles {
		bytes, err := kb.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal key bundle %d: %w", i, err)
		}

		filename := fmt.Sprintf("%s/key_bundle_%d.json", cacheDir, i)
		if err := os.WriteFile(filename, bytes, 0600); err != nil {
			return fmt.Errorf("failed to write key bundle %d to file: %w", i, err)
		}
	}
	return nil
}

func loadKeyBundlesFromCache() ([]ocr2key.KeyBundle, error) {
	cacheDir := "cache/keys"
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var keyBundles []ocr2key.KeyBundle
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "key_bundle_") {
			bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", cacheDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read key bundle file %s: %w", file.Name(), err)
			}

			kb, err := ocr2key.New(chaintype.EVM)
			if err != nil {
				return nil, fmt.Errorf("cannot create new key bundle from %s: %w", file.Name(), err)
			}
			if err := kb.Unmarshal(bytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal key bundle from %s: %w", file.Name(), err)
			}
			keyBundles = append(keyBundles, kb)
		}
	}

	if len(keyBundles) == 0 {
		return nil, errors.New("no key bundles found in cache directory")
	}
	return keyBundles, nil
}

func NewFeedID(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := crand.Read(buf[:])
	require.NoError(t, err, "cannot create feedID")
	return buf, "0x" + hex.EncodeToString(buf[:])
}

func saveFeedAddresses(feedsAddresses [][]FeedWithStreamID) error {
	cacheDir := "cache/feeds"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	filename := cacheDir + "/feed_addresses.json"
	bytes, err := json.Marshal(feedsAddresses)
	if err != nil {
		return fmt.Errorf("failed to marshal feed addresses: %w", err)
	}

	if err := os.WriteFile(filename, bytes, 0600); err != nil {
		return fmt.Errorf("failed to write feed addresses to file: %w", err)
	}

	return nil
}

func loadFeedAddressesFromCache() ([][]FeedWithStreamID, error) {
	bytes, err := os.ReadFile("cache/feeds/feed_addresses.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read feed addresses file: %w", err)
	}

	var feedsAddresses [][]FeedWithStreamID
	if err := json.Unmarshal(bytes, &feedsAddresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed addresses: %w", err)
	}

	return feedsAddresses, nil
}

type FeedConfig struct {
	FeedIDsIndex int32  `json:"feedIDsIndex"`
	Deviation    string `json:"deviation"`
	Heartbeat    int32  `json:"heartbeat"`
	RemappedID   string `json:"remappedID"`
}

// TODO shouldn't consumer address be configurable?
func WorkflowsJob(nodeID string, workflowName string, feeds []FeedConfig) *jobv1.ProposeJobRequest {
	const workflowTemplateLoad = `
 type = "workflow"
 schemaVersion = 1
 name = "{{ .WorkflowName }}"
 externalJobID = "{{ .JobID }}"
 workflow = """
 name: "{{ .WorkflowName }}"
 owner: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512'
 triggers:
  - id: streams-trigger@2.0.0
    config:
      feedIds:
 {{- range .Feeds }}
        - '{{ .FeedIDsIndex }}'
 {{- end }}
 consensus:
   - id: "offchain_reporting@1.0.0"
     ref: "evm_median"
     inputs:
       observations:
         - "$(trigger.outputs)"
     config:
       report_id: "0001"
       key_id: "evm"
       aggregation_method: "llo_streams"
       aggregation_config:
         streams:
{{- range .Feeds }}
           "{{ .FeedIDsIndex }}":
             deviation: "{{ .Deviation }}"
             heartbeat: {{ .Heartbeat }}
             remappedID: {{ .RemappedID }}
{{- end }}
       encoder: "EVM"
       encoder_config:
         abi: "(bytes32 RemappedID, uint224 Price, uint32 Timestamp)[] Reports"
 targets:
   - id: write_ethereum@1.0.0
     inputs:
       signed_report: "$(evm_median.outputs)"
     config:
       address: "0xEB739A9641938934D21A325A0C6b26126D48926A"
       params: ["$(report)"]
       abi: "receive(report bytes)"
       deltaStage: 2s
       schedule: allAtOnce
 """
 `

	tmpl, err := template.New("workflow").Parse(workflowTemplateLoad)

	if err != nil {
		panic(err)
	}
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"WorkflowName": workflowName,
		"Feeds":        feeds,
		"JobID":        uuid.NewString(),
	})
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String()}
}

func MockCapabilitiesJob(nodeID, binaryPath string, mocks []*MockCapabilities) *jobv1.ProposeJobRequest {
	jobTemplate := `type = "standardcapabilities"
			schemaVersion = 1
			externalJobID = "{{ .JobID }}"
			name = "mock-capability"
			forwardingAllowed = false
			command = "{{ .BinaryPath }}"
			config = """
				port=7777
		{{ range $index, $m := .Mocks }}
 		  [[DefaultMocks]]
				id="{{ $m.ID }}"
				description="{{ $m.Description }}"
				type="{{ $m.Type }}"
 		{{- end }}
			"""`
	tmpl, err := template.New("mock-job").Parse(jobTemplate)

	if err != nil {
		panic(err)
	}
	mockJobsData := make([]map[string]string, 0)
	for _, m := range mocks {
		mockJobsData = append(mockJobsData, map[string]string{
			"ID":          m.Name + "@" + m.Version,
			"Description": m.Description,
			"Type":        m.Type,
		})
	}

	jobUUID := uuid.NewString()
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"JobID":      jobUUID,
		"ShortID":    jobUUID[0:8],
		"BinaryPath": binaryPath,
		"Mocks":      mockJobsData,
	})
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String(),
	}
}

func capTypeToInt(capType string) uint8 {
	switch capType {
	case "trigger":
		return 0
	case "action":
		return 1
	case "consensus":
		return 2
	case "target":
		return 3
	default:
		panic("unknown capability type " + capType)
	}
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}
