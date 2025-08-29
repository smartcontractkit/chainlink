package cre

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp/benchspy"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
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
	Duration                      string                          `toml:"duration"`
	Blockchains                   []blockchain.Input              `toml:"blockchains" validate:"required"`
	NodeSets                      []*ns.Input                     `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                       `toml:"jd" validate:"required"`
	WorkflowRegistryConfiguration *cretypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *infra.Input                    `toml:"infra" validate:"required"`
	WorkflowDONLoad               *WorkflowLoad                   `toml:"workflow_load"`
	MockCapabilities              []*MockCapabilities             `toml:"mock_capabilities"`
	Chaos                         *Chaos                          `toml:"chaos"`
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
	blockchainOutput      []*cretypes.WrappedBlockchainOutput
	donTopology           *cretypes.DonTopology
	nodeOutput            []*cretypes.WrappedNodeOutput
}

func setupLoadTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfigLoadTest,
	mustSetCapabilitiesFn func(input []*ns.Input) []*cretypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []cretypes.CapabilityRegistryConfigFn,
	jobSpecFactoryFns []cretypes.JobSpecFn,
	workflowJobsFn cretypes.JobSpecFn,
) *loadTestSetupOutput {
	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		JobSpecFactoryFunctions:              jobSpecFactoryFns,
	}

	singleFileLogger := cldlogger.NewSingleFileLogger(t)
	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, singleFileLogger, universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	// Set inputs in the test config, so that they can be saved
	in.WorkflowRegistryConfiguration = &cretypes.WorkflowRegistryInput{}
	in.WorkflowRegistryConfiguration.Out = universalSetupOutput.WorkflowRegistryConfigurationOutput

	forwarderAddress, forwarderErr := crecontracts.FindAddressesForChain(
		universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // will not migrate now
		universalSetupOutput.BlockchainOutput[0].ChainSelector,
		keystone_changeset.KeystoneForwarder.String(),
	)
	require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainSelector)

	// Create workflow jobs only after capability registry configuration is complete to avoid initialization failures
	createJobsInput := creenv.CreateJobsWithJdOpInput{}
	createJobsDeps := creenv.CreateJobsWithJdOpDeps{
		Logger:                    testLogger,
		SingleFileLogger:          singleFileLogger,
		HomeChainBlockchainOutput: universalSetupOutput.BlockchainOutput[0].BlockchainOutput,
		AddressBook:               universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // will not migrate now
		JobSpecFactoryFunctions:   []cretypes.JobSpecFn{workflowJobsFn},
		FullCLDEnvOutput: &cretypes.FullCLDEnvironmentOutput{
			Environment: universalSetupOutput.CldEnvironment,
			DonTopology: universalSetupOutput.DonTopology,
		},
	}

	_, createJobsErr := operations.ExecuteOperation(universalSetupOutput.CldEnvironment.OperationsBundle, creenv.CreateJobsWithJdOpFactory("load-test-jobs", "1.0.0"), createJobsDeps, createJobsInput)
	require.NoError(t, createJobsErr, "failed to create jobs with Job Distributor")

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
	require.NotEmpty(t, os.Getenv("PROMETHEUS_URL"), "PROMETHEUS_URL must be set")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*cretypes.CapabilitiesAwareNodeSet {
		return []*cretypes.CapabilitiesAwareNodeSet{
			{
				Input:        input[0],
				Capabilities: []string{cretypes.ConsensusCapability},
				// TODO quick hack, this needs to be removed after the migration to TOML
				ComputedCapabilities: []string{cretypes.ConsensusCapability},
				DONTypes:             []string{cretypes.WorkflowDON},
				BootstrapNodeIndex:   0,
			},
			{
				Input:        input[1],
				Capabilities: []string{cretypes.MockCapability},
				// TODO quick hack, this needs to be removed after the migration to TOML
				ComputedCapabilities: []string{cretypes.MockCapability, cretypes.EVMCapability + "-1337"},
				DONTypes:             []string{cretypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex:   -1,
			},
		}
	}

	feedsAddresses := make([][]FeedWithStreamID, in.WorkflowDONLoad.Jobs)
	for i := range in.WorkflowDONLoad.Jobs {
		feedsAddresses[i] = make([]FeedWithStreamID, 0)
		for streamID := range in.WorkflowDONLoad.Streams {
			_, id := NewFeedIDDF2(t)
			feedsAddresses[i] = append(feedsAddresses[i], FeedWithStreamID{
				Feed:     id,
				StreamID: (in.WorkflowDONLoad.Streams * i) + streamID + 1,
			})
		}
	}

	mockJobSpecsFactoryFn := func(input *cretypes.JobSpecInput) (cretypes.DonsToJobSpecs, error) {
		donTojobSpecs := make(cretypes.DonsToJobSpecs, 0)

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			jobSpecs := make(cretypes.DonJobs, 0)
			workflowNodeSet, err2 := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.WorkerNode}, node.EqualLabels)
			if err2 != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err2, "failed to find worker nodes")
			}
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				if flags.HasFlag(donWithMetadata.Flags, cretypes.MockCapability) && in.MockCapabilities != nil {
					jobSpecs = append(jobSpecs, MockCapabilitiesJob(nodeID, "mock", in.MockCapabilities))
				}
			}

			donTojobSpecs[donWithMetadata.ID] = jobSpecs
		}

		return donTojobSpecs, nil
	}

	loadTestJobSpecsFactoryFn := func(input *cretypes.JobSpecInput) (cretypes.DonsToJobSpecs, error) {
		donTojobSpecs := make(cretypes.DonsToJobSpecs, 0)

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			jobSpecs := make(cretypes.DonJobs, 0)
			workflowNodeSet, err2 := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.WorkerNode}, node.EqualLabels)
			if err2 != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err2, "failed to find worker nodes")
			}
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}
				if flags.HasFlag(donWithMetadata.Flags, cretypes.WorkflowDON) {
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
			}

			donTojobSpecs[donWithMetadata.ID] = jobSpecs
		}

		return donTojobSpecs, nil
	}

	WorkflowDONLoadTestCapabilitiesFactoryFn := func(donFlags []string, _ *cretypes.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
		var capabilities []keystone_changeset.DONCapabilityWithConfig

		if flags.HasFlag(donFlags, cretypes.MockCapability) {
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

		if flags.HasFlag(donFlags, cretypes.CustomComputeCapability) {
			capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "custom-compute",
					Version:        "1.0.0",
					CapabilityType: 1, // ACTION
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}

		if flags.HasFlag(donFlags, cretypes.ConsensusCapability) {
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

		return capabilities, nil
	}

	homeChain := in.Blockchains[0]
	homeChainIDUint64, homeChainErr := strconv.ParseUint(homeChain.ChainID, 10, 64)
	require.NoError(t, homeChainErr, "failed to convert chain ID to int")

	setupOutput := setupLoadTestEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		[]cretypes.CapabilityRegistryConfigFn{WorkflowDONLoadTestCapabilitiesFactoryFn, registerEVMWithV1},
		[]cretypes.JobSpecFn{mockJobSpecsFactoryFn, consensusJobSpec(homeChainIDUint64)},
		loadTestJobSpecsFactoryFn,
	)

	ctx := t.Context()
	// Get OCR2 keys needed to sign the reports
	kb := make([]ocr2key.KeyBundle, 0)
	for _, don := range setupOutput.donTopology.DonsWithMetadata {
		if flags.HasFlag(don.Flags, cretypes.MockCapability) {
			for _, n := range don.DON.Nodes {
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

	// If were not running in CI then save the feeds and OCR2 keys to a file so we can reuse them later
	cacheClients := false
	if os.Getenv("CI") != "true" {
		cacheClients = true
		require.NoError(t, saveFeedAddresses(feedsAddresses), "could not save feeds")

		// Export key bundles so we can import them later in another test, used when crib cluster is already setup and we just want to connect to mocks for a different test
		require.NoError(t, saveKeyBundles(kb), "could not save OCR2 Keys")
	}
	testLogger.Info().Msg("Connecting to mock capabilities...")

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.Type == infra.Docker {
		for _, nodeSet := range in.NodeSets {
			if nodeSet.Name == "capabilities" {
				for _, n := range nodeSet.NodeSpecs {
					if len(n.Node.CustomPorts) == 0 {
						panic("no custom port specified, mock capability running in kind must have a custom port in order to connect")
					}
					ports := strings.Split(n.Node.CustomPorts[0], ":")
					mockClientsAddress = append(mockClientsAddress, "127.0.0.1:"+ports[0])
				}
			}
		}
	} else {
		for i := range setupOutput.nodeOutput[1].CLNodes {
			mockClientsAddress = append(mockClientsAddress, fmt.Sprintf("%s-%s-%d-mock.main.stage.cldev.sh:443", in.Infra.CRIB.Namespace, setupOutput.nodeOutput[1].NodeSetName, i-1))
		}
	}

	require.NotEmpty(t, mockClientsAddress, "Could not create mock capability client addresses")

	// Use insecure gRPC connection for local Docker containers. For AWS, use TLS credentials
	// due to ingress requirements, as grpc.insecure.NewCredentials() doesn't work properly with AWS ingress
	useInsecure := in.Infra.Type == infra.Docker

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, useInsecure, cacheClients), "could not connect to mock capabilities")

	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel), "could not hook into mock executable")

	// Wait for the remote capability to be exposed, we check if the streams-trigger has subscribers
	require.NoError(t, mocksClient.WaitForTriggerSubscribers(ctx, "streams-trigger@2.0.0", time.Minute*5), "error while waiting for trigger subscribers")

	labels := map[string]string{
		"go_test_name": "workflow-don-load-test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	generator, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		CallTimeout: time.Minute * 2, // Give enough time for the workflow to execute
		LoadType:    wasp.RPS,
		Schedule: wasp.Combine(
			wasp.Plain(4, 10*time.Minute),
		),
		Gun:                   NewStreamsGun(mocksClient, kb, feedsAddresses, "streams-trigger@2.0.0", receiveChannel, int(in.WorkflowDONLoad.Streams), int(in.WorkflowDONLoad.Jobs)),
		Labels:                labels,
		RateLimitUnitDuration: time.Minute,
	})
	require.NoError(t, err, "could not create generator")
	// run the load
	generator.Run(true)

	tag := "local-test-" + time.Now().Format("20060102150405")
	if os.Getenv("CI") == "true" {
		// When running in CI, use the GitHub commit SHA
		commitSHA := os.Getenv("GITHUB_SHA")
		if commitSHA != "" {
			tag = commitSHA + time.Now().Format("20060102150405")
		}
	} else if gitSHA := os.Getenv("GITHUB_SHA"); gitSHA != "" {
		// For local runs with manually set GITHUB_SHA
		tag = gitSHA
	}

	promConfig := benchspy.NewPrometheusConfig()

	prometheusExecutor, err := benchspy.NewPrometheusQueryExecutor(
		map[string]string{
			"cpu_percent":          `avg (rate(container_cpu_usage_seconds_total{name=~"workflow-node[1-9][0-9]*"}[10m]) * 100)`,
			"mem_peak":             `avg (max_over_time(container_memory_working_set_bytes{name=~"workflow-node[1-9][0-9]*"}[10m]))`,
			"mem_avg":              `avg (avg_over_time(container_memory_working_set_bytes{name=~"workflow-node[1-9][0-9]*"}[10m]))`,
			"disk_io_time_seconds": `avg (container_fs_io_time_seconds_total{name=~"workflow-node[1-9][0-9]*"})`,
			"network_tx":           `avg (container_network_transmit_bytes_total{name=~"workflow-node[1-9][0-9]*"})`,
			"network_rx":           `avg (container_network_receive_bytes_total{name=~"workflow-node[1-9][0-9]*"})`,
		},
		promConfig,
	)
	require.NoError(t, err)

	benchmarkReport, baselineReport, err := benchspy.FetchNewStandardReportAndLoadLatestPrevious(
		ctx,
		tag,
		benchspy.WithStandardQueries(benchspy.StandardQueryExecutor_Direct),
		benchspy.WithQueryExecutors(prometheusExecutor),
		benchspy.WithGenerators(generator),
	)
	require.NoError(t, err, "failed to create benchmark report")

	fetchErr := benchmarkReport.FetchData(ctx)
	require.NoError(t, fetchErr, "failed to fetch data for benchmark report")

	path, storeErr := benchmarkReport.Store()
	require.NoError(t, storeErr, "failed to store benchmark report", path)
	require.NoError(t, err, "workflow load test did not finish successfully")

	// Compare benchmark with baseline if available
	if baselineReport != nil {
		compareBenchmarkReports(t, benchmarkReport, baselineReport)
	}
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
				wasp.Plain(4, 5*time.Minute),
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
	go sg.waitLoop()
	return sg
}

func (s *StreamsGun) Call(l *wasp.Generator) *wasp.Response {
	event, eventID, timestamp, err := s.createReport()
	if err != nil {
		return &wasp.Response{Failed: true, Error: err.Error()}
	}
	err = s.createWaitChannelForTimestamp(timestamp.Unix())
	if err != nil {
		return &wasp.Response{Failed: true, Error: err.Error()}
	}

	outputs, err := event.ToMap()
	if err != nil {
		return &wasp.Response{Failed: true, Error: err.Error()}
	}

	outputsBytes, err := mock_capability.MapToBytes(outputs)
	if err != nil {
		return &wasp.Response{Failed: true, Error: err.Error()}
	}

	message := pb.SendTriggerEventRequest{
		TriggerID: s.triggerID,
		ID:        eventID,
		Outputs:   outputsBytes,
	}

	err = s.capProxy.SendTrigger(context.Background(), &message)
	if err != nil {
		framework.L.Error().Msgf("error sending trigger: %s", err.Error())
		return &wasp.Response{Failed: true, Error: err.Error()}
	}

	// Wait for the DON to execute on the write target
	err = s.waitForReportWithTimestamp(timestamp.Unix())
	if err != nil {
		return &wasp.Response{Failed: true, Error: err.Error()}
	}
	return &wasp.Response{}
}

func (s *StreamsGun) waitLoop() {
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

func (s *StreamsGun) createWaitChannelForTimestamp(reportTimestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitChans == nil {
		s.waitChans = make(map[int64]chan interface{})
	}
	if _, exists := s.waitChans[reportTimestamp]; exists {
		return fmt.Errorf("cannot create wait channel, timestamp  %d already exits", reportTimestamp)
	}
	s.waitChans[reportTimestamp] = make(chan interface{})
	return nil
}

func (s *StreamsGun) waitForReportWithTimestamp(timestamp int64) error {
	s.mu.Lock()
	ch, exists := s.waitChans[timestamp]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("cannot wait, timestamp  %q does not exist", timestamp)
	}
	s.mu.Unlock()
	<-ch
	delete(s.waitChans, timestamp)
	framework.L.Info().Msgf("ACK report with timestamp %d", timestamp)
	return nil
}

func (s *StreamsGun) createReport() (capabilities.OCRTriggerEvent, string, time.Time, error) {
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

	event, eventID, err := createFeedReport(logger.Nop(), price, uint64(timestamp.UnixNano()), feeds, s.keyBundles) //nolint:gosec // G115 don't care in test code
	if err != nil {
		return capabilities.OCRTriggerEvent{}, "", time.Time{}, err
	}

	framework.L.Info().Msgf("create report with timestamp %d containing %d feeds in %s", timestamp.Unix(), len(feeds), time.Since(start))
	return *event, eventID, timestamp, nil
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
			return nil, "", err2
		}
		event.Sigs = append(event.Sigs, capabilities.OCRAttributedOnchainSignature{
			Signer:    uint32(i), //nolint:gosec // G115 don't care in test code
			Signature: sig,
		})
	}

	return event, eventID, nil
}

func decodeTargetInput(inputs *values.Map) (evm.TargetRequest, error) {
	var r evm.TargetRequest
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
		framework.L.Info().Msgf("Saving OCR2 Key ID: %s, OnChainPublicKey: %s", kb.ID(), kb.OnChainPublicKey())
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

// NewFeedIDDF2 creates a random Data Feeds 2.0 format https://docs.google.com/document/d/13ciwTx8lSUfyz1IdETwpxlIVSn1lwYzGtzOBBTpl5Vg/edit?tab=t.0#heading=h.dxx2wwn1dmoz
func NewFeedIDDF2(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := crand.Read(buf[:])
	require.NoError(t, err, "cannot create feedID")
	buf[0] = 0x01 // format byte
	buf[5] = 0x00 // attribute
	buf[6] = 0x03 // attribute
	buf[7] = 0x00 // data type byte

	for i := 8; i < 16; i++ {
		buf[i] = 0x00
	}

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
   - id: write_ethereum_mock@1.0.0
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

func compareBenchmarkReports(t *testing.T, baselineReport, currentReport *benchspy.StandardReport) {
	// Define threshold percentages for different metrics
	thresholds := map[string]float64{
		"cpu_percent":             10.0, // 10% increase
		"mem_peak":                10.0,
		"mem_avg":                 10.0,
		"network_tx":              10.0,
		"network_rx":              10.0,
		"95th_percentile_latency": 10.0,
		"99th_percentile_latency": 10.0,
		"median_latency":          10.0,
		"error_rate":              10.0,
		"max_latency":             10.0,
	}

	// Fetch all metrics
	require.Len(t, baselineReport.QueryExecutors, 2, "expected two query executors in baseline report")
	require.Len(t, currentReport.QueryExecutors, 2, "expected two query executors in benchmark report")

	baselineReportMetrics := make(map[string]float64)
	for _, qe := range baselineReport.QueryExecutors {
		for metricName, metricValue := range qe.Results() {
			// Check if the metricValue is a slice
			if sliceVal, ok := metricValue.([]float64); ok && len(sliceVal) > 0 {
				// If it's a slice of float64, get the last element
				baselineReportMetrics[metricName] = sliceVal[len(sliceVal)-1]
			} else if floatVal, ok := metricValue.(float64); ok {
				// If it's a single float64, use it directly
				baselineReportMetrics[metricName] = floatVal
			} else if vector, ok := metricValue.(model.Vector); ok {
				if len(vector) > 0 {
					// Use the most recent sample's value from the vector
					baselineReportMetrics[metricName] = float64(vector[len(vector)-1].Value)
				} else {
					// Log the case where vector is empty
					framework.L.Warn().Msgf("Metric %s has empty vector value", metricName)
				}
			} else {
				// Log the case where the value is not a float64 or slice of float64
				framework.L.Warn().Msgf("Metric %s has unsupported value type: %T", metricName, metricValue)
			}
		}
	}

	currentReportMetrics := make(map[string]float64)
	for _, qe := range currentReport.QueryExecutors {
		for metricName, metricValue := range qe.Results() {
			// Check if the metricValue is a slice
			if sliceVal, ok := metricValue.([]float64); ok && len(sliceVal) > 0 {
				// If it's a slice of float64, get the last element
				currentReportMetrics[metricName] = sliceVal[len(sliceVal)-1]
			} else if floatVal, ok := metricValue.(float64); ok {
				// If it's a single float64, use it directly
				currentReportMetrics[metricName] = floatVal
			} else if vector, ok := metricValue.(model.Vector); ok {
				if len(vector) > 0 {
					// Use the most recent sample's value from the vector
					currentReportMetrics[metricName] = float64(vector[len(vector)-1].Value)
				} else {
					// Log the case where vector is empty
					framework.L.Warn().Msgf("Metric %s has empty vector value", metricName)
				}
			} else {
				// Log the case where the value is not a float64 or slice of float64
				framework.L.Warn().Msgf("Metric %s has unsupported value type: %T", metricName, metricValue)
			}
		}
	}

	// 	// Compare metrics
	var warnings []string
	for metric, threshold := range thresholds {
		if baselineReportMetrics[metric] > 0 {
			percentIncrease := ((currentReportMetrics[metric] - baselineReportMetrics[metric]) / baselineReportMetrics[metric]) * 100
			if percentIncrease > threshold {
				warnings = append(warnings, fmt.Sprintf(
					"PERFORMANCE REGRESSION: %s increased by %.2f%% (baseline: %.2f, current: %.2f, threshold: %.2f%%)",
					metric, percentIncrease, baselineReportMetrics[metric], currentReportMetrics[metric], threshold,
				))
			}
		}
	}

	// Log any warnings
	if len(warnings) > 0 {
		framework.L.Warn().Msgf("Performance regression detected compared to baseline %s", baselineReport.CommitOrTag)
		for _, warning := range warnings {
			framework.L.Warn().Msg(warning)
		}
	} else {
		framework.L.Info().Msgf("No significant performance regressions detected compared to baseline %s", baselineReport.CommitOrTag)
	}
}

// Deprecated: remove this once load tests have been migrated
func registerEVMWithV1(_ []string, nodeSetInput *cretypes.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	capabilities := make([]keystone_changeset.DONCapabilityWithConfig, 0)

	if nodeSetInput == nil {
		return nil, errors.New("node set input is nil")
	}

	// it's fine if there are no chain capabilities
	if nodeSetInput.ChainCapabilities == nil {
		return nil, nil
	}

	if _, ok := nodeSetInput.ChainCapabilities[cretypes.WriteEVMCapability]; !ok {
		return nil, nil
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[cretypes.WriteEVMCapability].EnabledChains {
		fullName := evm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   splitName[0],
				Version:        splitName[1],
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}

// Deprecated: remove this once load tests have been migrated
func consensusJobSpec(chainID uint64) cretypes.JobSpecFn {
	return func(input *cretypes.JobSpecInput) (cretypes.DonsToJobSpecs, error) {
		if input.DonTopology == nil {
			return nil, errors.New("topology is nil")
		}
		donToJobSpecs := make(cretypes.DonsToJobSpecs)

		ocr3Key := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"capability_ocr3",
		)
		ocr3CapabilityAddress, err := input.CldEnvironment.DataStore.Addresses().Get(ocr3Key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get Vault capability address")
		}

		donTimeKey := datastore.NewAddressRefKey(
			input.DonTopology.HomeChainSelector,
			datastore.ContractType(keystone_changeset.OCR3Capability.String()),
			semver.MustParse("1.0.0"),
			"DONTime",
		)
		donTimeAddress, err := input.CldEnvironment.DataStore.Addresses().Get(donTimeKey)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get DON Time address")
		}

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			if !flags.HasFlag(donWithMetadata.Flags, cretypes.ConsensusCapability) {
				continue
			}

			// create job specs for the worker nodes
			workflowNodeSet, err := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.WorkerNode}, node.EqualLabels)
			if err != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, errors.Wrap(err, "failed to find worker nodes")
			}

			// look for boostrap node and then for required values in its labels
			bootstrapNode, bootErr := node.FindOneWithLabel(donWithMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.BootstrapNode}, node.EqualLabels)
			if bootErr != nil {
				return nil, errors.Wrap(bootErr, "failed to find bootstrap node")
			}

			donBootstrapNodePeerID, pIDErr := node.ToP2PID(bootstrapNode, node.KeyExtractingTransformFn)
			if pIDErr != nil {
				return nil, errors.Wrap(pIDErr, "failed to get bootstrap node peer ID")
			}

			donBootstrapNodeHost, hostErr := node.FindLabelValue(bootstrapNode, node.HostLabelKey)
			if hostErr != nil {
				return nil, errors.Wrap(hostErr, "failed to get bootstrap node host from labels")
			}

			bootstrapNodeID, nodeIDErr := node.FindLabelValue(bootstrapNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get bootstrap node id from labels")
			}

			// create job specs for the bootstrap node
			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.BootstrapOCR3(bootstrapNodeID, "ocr3-capability", ocr3CapabilityAddress.Address, chainID))

			ocrPeeringData := cretypes.OCRPeeringData{
				OCRBootstraperPeerID: donBootstrapNodePeerID,
				OCRBootstraperHost:   donBootstrapNodeHost,
				Port:                 don.OCRPeeringPort,
			}

			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
				}

				nodeEthAddr, ethErr := node.FindLabelValue(workerNode, node.AddressKeyFromSelector(input.DonTopology.HomeChainSelector))
				if ethErr != nil {
					return nil, errors.Wrap(ethErr, "failed to get eth address from labels")
				}

				ocr2KeyBundleID, ocr2Err := node.FindLabelValue(workerNode, node.NodeOCR2KeyBundleIDKey)
				if ocr2Err != nil {
					return nil, errors.Wrap(ocr2Err, "failed to get ocr2 key bundle id from labels")
				}
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerOCR3(nodeID, ocr3CapabilityAddress.Address, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainID))
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.DonTimeJob(nodeID, donTimeAddress.Address, nodeEthAddr, ocr2KeyBundleID, ocrPeeringData, chainID))
			}
		}

		return donToJobSpecs, nil
	}
}
