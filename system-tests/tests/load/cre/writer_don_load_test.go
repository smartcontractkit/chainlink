package cre

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	consensustypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	changeset2 "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
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

type WriterTest struct {
	NrOfFeeds     int    `toml:"nr_of_feeds"`
	WorkflowName  string `toml:"workflow_name"`
	WorkflowOwner string `toml:"workflow_owner"`
	WorkflowID    string `toml:"workflow_id"`
}
type TestConfigLoadTestWriter struct {
	Blockchains                   []*blockchain.Input                  `toml:"blockchains" validate:"required"`
	NodeSets                      []*ns.Input                          `toml:"nodesets" validate:"required"`
	JD                            *jd.Input                            `toml:"jd" validate:"required"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *libtypes.InfraInput                 `toml:"infra" validate:"required"`
	MockCapabilities              []*MockCapabilities                  `toml:"mock_capabilities"`
	BinariesConfig                *BinariesConfig                      `toml:"binaries_config"`
	WriterTest                    *WriterTest                          `toml:"writer_test"`
}

func setupLoadTestWriterEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfigLoadTestWriter,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
	jobSpecFactoryFns []keystonetypes.JobSpecFactoryFn,
	feedIDs []string,
	workflowNames []string,
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

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")
	// Set inputs in the test config, so that they can be saved
	in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{}
	in.WorkflowRegistryConfiguration.Out = universalSetupOutput.WorkflowRegistryConfigurationOutput

	forwarderAddress, forwarderErr := libcontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, keystone_changeset.KeystoneForwarder.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainSelector)

	// DF cache start

	// Deploy
	deployConfig := df_changeset_types.DeployConfig{
		ChainsToDeploy: []uint64{universalSetupOutput.BlockchainOutput[0].ChainSelector},
		Labels:         []string{"data-feeds"}, // label required by the changeset
	}
	dfOutput, dfErr := changeset2.RunChangeset(changeset.DeployCacheChangeset, *universalSetupOutput.CldEnvironment, deployConfig)
	require.NoError(t, dfErr, "failed to deploy data feed cache contract")

	mergeErr := universalSetupOutput.CldEnvironment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
	require.NoError(t, mergeErr, "failed to merge address book")

	dfCacheAddress, dfCacheErr := libcontracts.FindAddressesForChain(universalSetupOutput.CldEnvironment.ExistingAddresses, universalSetupOutput.BlockchainOutput[0].ChainSelector, changeset.DataFeedsCache.String()) //nolint:staticcheck // won't migrate now
	require.NoError(t, dfCacheErr, "failed to find df cache address for chain %d", universalSetupOutput.BlockchainOutput[0].ChainSelector)
	// Config
	_, configErr := libcontracts.ConfigureDataFeedsCache(testLogger, &keystonetypes.ConfigureDataFeedsCacheInput{
		CldEnv:                universalSetupOutput.CldEnvironment,
		ChainSelector:         universalSetupOutput.BlockchainOutput[0].ChainSelector,
		FeedIDs:               feedIDs,
		Descriptions:          feedIDs,
		DataFeedsCacheAddress: dfCacheAddress,
		AdminAddress:          universalSetupOutput.BlockchainOutput[0].SethClient.MustGetRootKeyAddress(),
		AllowedSenders:        []common.Address{forwarderAddress},
		AllowedWorkflowOwners: []common.Address{common.HexToAddress(in.WriterTest.WorkflowOwner)},
		AllowedWorkflowNames:  workflowNames,
		Out:                   nil,
	})
	require.NoError(t, configErr, "failed to configure data feeds cache")

	return &loadTestSetupOutput{
		dataFeedsCacheAddress: dfCacheAddress,
		forwarderAddress:      forwarderAddress,
		blockchainOutput:      universalSetupOutput.BlockchainOutput,
		donTopology:           universalSetupOutput.DonTopology,
		nodeOutput:            universalSetupOutput.NodeOutput,
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

	registryChain := in.Blockchains[0]
	homeChainIDUint64, homeChainErr := strconv.ParseUint(registryChain.ChainID, 10, 64)
	require.NoError(t, homeChainErr, "failed to convert chain ID to int")

	feedIDs := make([]string, 0)
	for range in.WriterTest.NrOfFeeds {
		_, id := NewFeedIDDF2(t)
		feedIDs = append(feedIDs, strings.TrimPrefix(id, "0x"))
	}

	setupOutput := setupLoadTestWriterEnvironment(
		t,
		testLogger,
		in,
		mustSetCapabilitiesFn,
		//nolint:gosec // disable G115
		[]func(donFlags []string) []keystone_changeset.DONCapabilityWithConfig{WriterDONLoadTestCapabilitiesFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(homeChainIDUint64)))},
		[]keystonetypes.JobSpecFactoryFn{loadTestJobSpecsFactoryFn, consensus.ConsensusJobSpecFactoryFn(homeChainIDUint64)},
		feedIDs,
		[]string{in.WriterTest.WorkflowName},
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
			lidebug.PrintTestDebug(ctx, t.Name(), testLogger, debugInput)
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

	f := 0
	// Nr of signatures needs to be equal with f+1, compute f based on the nr of ocr3 worker nodes
	for _, donMetadata := range setupOutput.donTopology.DonsWithMetadata {
		if flags.HasFlag(donMetadata.Flags, keystonetypes.OCR3Capability) {
			workerNodes, workerNodesErr := node.FindManyWithLabel(donMetadata.NodesMetadata, &keystonetypes.Label{
				Key:   node.NodeTypeKey,
				Value: keystonetypes.WorkerNode,
			}, node.EqualLabels)
			require.NoError(t, workerNodesErr, "could not find any worker nodes for ocr3")

			f = (len(workerNodes) - 1) / 3
		}
	}
	require.NotEqual(t, 0, f, "f cannot be 0")
	tParams := testParams{
		workflowName:          in.WriterTest.WorkflowName,
		workflowOwner:         in.WriterTest.WorkflowOwner,
		workflowID:            in.WriterTest.WorkflowID,
		feeds:                 feedIDs,
		nrOfReports:           in.WriterTest.NrOfFeeds,
		dataFeedsCacheAddress: setupOutput.dataFeedsCacheAddress,
		forwarderAddress:      setupOutput.forwarderAddress,
		f:                     f,
	}

	require.NoError(t, exportTestParams(tParams), "could not save test params")

	// Export key bundles so we can import them later in another test, used when crib cluster is already setup and we just want to connect to mocks for a different test
	require.NoError(t, saveKeyBundles(kb), "could not save OCR2 Keys")

	require.NoError(t, saveClientURL(setupOutput.blockchainOutput[0].SethClient.URL), "could not save seth client url")

	testLogger.Info().Msg("Connecting to mock capabilities...")

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.InfraType == "docker" {
		for _, nodeSet := range in.NodeSets {
			if nodeSet.Name == "writer" {
				for i, n := range nodeSet.NodeSpecs {
					if i == 0 {
						continue
					}
					if len(n.Node.CustomPorts) == 0 {
						panic("no custom port specified, mock capability running in kind must have a custom port in order to connect")
					}
					ports := strings.Split(n.Node.CustomPorts[0], ":")
					mockClientsAddress = append(mockClientsAddress, "127.0.0.1:"+ports[0])
				}
			}
		}
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
			CallTimeout: time.Minute * 5, // Give enough time for write to execute
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(1, 10*time.Minute),
			),
			Gun:                   NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", setupOutput.blockchainOutput[0].SethClient, tParams),
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

	tParams, err := importTestParams()
	require.NoError(t, err, "could not import test params")

	pkey := os.Getenv("PRIVATE_KEY")
	require.NotEmpty(t, pkey, "PRIVATE_KEY env var must be set")

	clientURL, err := loadClientURL()
	require.NoError(t, err, "could load the client url")

	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(clientURL).
		WithPrivateKeys([]string{pkey}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()

	require.NoError(t, err, "failed to create seth client")

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

	sg := NewWriterGun(mocksClient, kb, "write_geth-testnet@1.0.0", sethClient, tParams)
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

type testParams struct {
	workflowName          string
	workflowOwner         string
	workflowID            string
	feeds                 []string
	nrOfReports           int
	dataFeedsCacheAddress common.Address
	forwarderAddress      common.Address
	f                     int
}

func exportTestParams(params testParams) error {
	// Create cache directory if it doesn't exist
	cacheDir := filepath.Join(os.TempDir(), "cache")
	if err := os.MkdirAll(cacheDir, 0600); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Marshal data to JSON
	data := map[string]interface{}{
		"workflowName":          params.workflowName,
		"workflowOwner":         params.workflowOwner,
		"workflowID":            params.workflowID,
		"feeds":                 params.feeds,
		"nrOfReports":           params.nrOfReports,
		"dataFeedsCacheAddress": params.dataFeedsCacheAddress.String(),
		"forwarderAddress":      params.forwarderAddress.String(),
		"f":                     params.f,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal test params: %w", err)
	}

	// Write to file in cache directory
	cacheFile := filepath.Join(cacheDir, "test_params.json")
	if err := os.WriteFile(cacheFile, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write test params file: %w", err)
	}

	return nil
}

func importTestParams() (testParams, error) {
	var params testParams

	data, err := os.ReadFile(filepath.Join(os.TempDir(), "cache") + "/test_params.json")
	if err != nil {
		return params, fmt.Errorf("failed to read test params file: %w", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return params, fmt.Errorf("failed to unmarshal test params: %w", err)
	}

	params.workflowName = rawData["workflowName"].(string)
	params.workflowOwner = rawData["workflowOwner"].(string)
	params.workflowID = rawData["workflowID"].(string)

	// Convert interface{} array to []string
	feedsRaw := rawData["feeds"].([]interface{})
	params.feeds = make([]string, len(feedsRaw))
	for i, v := range feedsRaw {
		params.feeds[i] = v.(string)
	}

	params.nrOfReports = int(rawData["nrOfReports"].(float64))
	params.f = int(rawData["f"].(float64))
	params.dataFeedsCacheAddress = common.HexToAddress(rawData["dataFeedsCacheAddress"].(string))
	params.forwarderAddress = common.HexToAddress(rawData["forwarderAddress"].(string))

	return params, nil
}

type WriterGun struct {
	capProxy   *mock_capability.Controller
	keyBundles []ocr2key.KeyBundle
	triggerID  string
	waitChans  map[uint8]chan interface{}
	mu         sync.Mutex
	reportID   uint8
	seqNr      uint32
	testParams testParams

	seth *seth.Client
}

func NewWriterGun(capProxy *mock_capability.Controller, keyBundles []ocr2key.KeyBundle, triggerID string, seth *seth.Client, params testParams) *WriterGun {
	sg := &WriterGun{
		capProxy:   capProxy,
		keyBundles: keyBundles,
		triggerID:  triggerID,
		reportID:   1,
		seqNr:      1,
		seth:       seth,
		testParams: params,
		waitChans:  make(map[uint8]chan interface{}),
	}

	go func() {
		ethClinet, _ := seth.Client.(*ethclient.Client)
		fwd, _ := forwarder.NewKeystoneForwarder(params.forwarderAddress, ethClinet)
		ch := make(chan *types2.Header)
		_, err := seth.Client.SubscribeNewHead(context.TODO(), ch)
		if err != nil {
			panic(err)
		}
		framework.L.Info().Msg("Subscribed to new heads, check every block for tx event")
		for {
			head := <-ch
			block, err := seth.Client.BlockByNumber(context.TODO(), head.Number)
			if err != nil {
				panic(err)
			}

			framework.L.Info().Msgf("Block %d tx count: %d", block.NumberU64(), block.Transactions().Len())
			for _, t := range block.Transactions() {
				receipt, err := seth.Client.TransactionReceipt(context.TODO(), t.Hash())
				if err != nil {
					framework.L.Error().Msgf("could not fetch tx receipt %v", err)
					continue
				}

				for _, l := range receipt.Logs {
					log, err := fwd.ParseReportProcessed(*l)
					if err != nil {
						continue
					}
					framework.L.Info().Msgf("Found ReportProcessed event ReportID %04x result %t", log.ReportId, log.Result)
					reportID := log.ReportId[1]
					sg.mu.Lock()
					if ch, exists := sg.waitChans[reportID]; exists {
						close(ch) // Signal completion
						delete(sg.waitChans, reportID)
					}
					sg.mu.Unlock()
				}
			}
		}
	}()

	return sg
}

func (s *WriterGun) Call(l *wasp.Generator) *wasp.Response {
	// Generate workflow execution metadata
	metadata, err := s.createWorkflowMetadata()
	if err != nil {
		framework.L.Error().Err(err)
		return &wasp.Response{Error: err.Error()}
	}

	// Create report context from sequence number and config digest
	repContext, err := s.createReportContext()
	if err != nil {
		framework.L.Error().Err(err)
		return &wasp.Response{Error: err.Error()}
	}

	// Create and encode report data
	encodedReport, err := s.createEncodedReport(metadata)
	if err != nil {
		framework.L.Error().Err(err)
		return &wasp.Response{Error: err.Error()}
	}

	// Generate signatures for the report
	sigs, err := s.generateSignatures(encodedReport, repContext)
	if err != nil {
		framework.L.Error().Err(err)
		return &wasp.Response{Error: err.Error()}
	}

	s.mu.Lock()
	ch := make(chan interface{})
	s.waitChans[s.reportID] = ch
	s.mu.Unlock()
	// Create executable request and execute
	err = s.executeRequest(metadata, encodedReport, sigs, repContext)
	if err != nil {
		framework.L.Error().Err(err)
		return &wasp.Response{Error: err.Error()}
	}

	s.reportID++
	s.seqNr++

	<-ch // Wait for confirmation
	return &wasp.Response{}
}

func (s *WriterGun) createWorkflowMetadata() (*pb2.Metadata, error) {
	executionID, err := workflows.EncodeExecutionID(s.testParams.workflowID, uuid.New().String())
	if err != nil {
		return nil, fmt.Errorf("failed to encode execution ID: %w", err)
	}

	return &pb2.Metadata{
		WorkflowID:               s.testParams.workflowID,
		WorkflowOwner:            strings.TrimPrefix(s.testParams.workflowOwner, "0x"),
		WorkflowExecutionID:      executionID,
		WorkflowName:             convertToHashedWorkflowName(s.testParams.workflowName),
		WorkflowDonID:            1,
		WorkflowDonConfigVersion: 1,
		ReferenceID:              "write_geth-testnet@1.0.0",
		DecodedWorkflowName:      s.testParams.workflowName,
	}, nil
}

func (s *WriterGun) createReportContext() ([]byte, error) {
	seqToEpoch := make([]byte, 32)
	binary.BigEndian.PutUint32(seqToEpoch[32-5:32-1], s.seqNr)
	zeros := make([]byte, 32)
	configDigest := [32]byte{1}
	return append(append(configDigest[:], seqToEpoch...), zeros...), nil
}

func (s *WriterGun) createEncodedReport(metadata *pb2.Metadata) ([]byte, error) {
	timestamp := time.Now()

	evmEncoder, err := s.createEVMEncoder()
	if err != nil {
		return nil, err
	}

	reports := s.generateReports(timestamp)
	fakeReport := s.createFakeReport(reports, metadata, timestamp)

	wrappedFakeReport, err := values.NewMap(fakeReport)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap fake report: %w", err)
	}

	return evmEncoder.Encode(context.TODO(), *wrappedFakeReport)
}

func (s *WriterGun) generateSignatures(encodedReport []byte, repContext []byte) ([][]byte, error) {
	sigs := make([][]byte, s.testParams.f+1)
	sigData := append(crypto.Keccak256(encodedReport), repContext...)
	fullHash := crypto.Keccak256(sigData)

	for i := range s.testParams.f + 1 {
		sig, err := s.keyBundles[i].SignBlob(fullHash)
		if err != nil {
			return nil, fmt.Errorf("failed to sign blob: %w", err)
		}
		sigs[i] = sig
	}

	return sigs, nil
}

func (s *WriterGun) executeRequest(metadata *pb2.Metadata, encodedReport []byte, sigs [][]byte, repContext []byte) error {
	validInputs, validConfig, err := s.createRequestInputs(encodedReport, sigs, repContext)
	if err != nil {
		return err
	}

	req := &pb2.ExecutableRequest{
		ID:              s.triggerID,
		CapabilityType:  4,
		RequestMetadata: metadata,
		Config:          validConfig,
		Inputs:          validInputs,
	}

	return s.capProxy.Execute(context.TODO(), req)
}
func (s *WriterGun) createEVMEncoder() (consensustypes.Encoder, error) {
	evmEncoderConfig := map[string]any{
		"abi": "(bytes32 FeedID, uint32 Timestamp, uint224 Price)[] Reports",
	}
	wrappedEVMEncoderConfig, err := values.NewMap(evmEncoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder config: %w", err)
	}
	return evm.NewEVMEncoder(wrappedEVMEncoderConfig)
}

// generateReports creates a slice of reports for each feed ID, where each report contains the necessary data to be EVM encoded conforming to `(bytes32 FeedID, uint32 Timestamp, uint224 Price)[] Reports`
func (s *WriterGun) generateReports(timestamp time.Time) []any {
	reports := make([]any, 0, s.testParams.nrOfReports)
	for i := range s.testParams.nrOfReports {
		feedIDBytes, err := stringTo32Byte(s.testParams.feeds[i])
		if err != nil {
			framework.L.Error().Msgf("failed to convert feed ID to bytes: %v", err)
			continue
		}
		reports = append(reports, struct {
			FeedID    [32]byte
			Price     string
			Timestamp uint32
		}{
			FeedID:    feedIDBytes,
			Price:     big.NewInt(int64(rand.Intn(1000000))).String(),
			Timestamp: uint32(timestamp.Unix()), //nolint:gosec // disable G115
		})
	}
	return reports
}

func (s *WriterGun) createFakeReport(reports []any, metadata *pb2.Metadata, timestamp time.Time) map[string]any {
	return map[string]any{
		"Reports": reports,
		consensustypes.MetadataFieldName: consensustypes.Metadata{
			Version:          1,
			ExecutionID:      metadata.WorkflowExecutionID,
			Timestamp:        uint32(timestamp.Unix()), //nolint:gosec // disable G115
			DONID:            1,
			DONConfigVersion: 1,
			WorkflowID:       metadata.WorkflowID,
			WorkflowName:     metadata.WorkflowName,
			WorkflowOwner:    metadata.WorkflowOwner,
			ReportID:         fmt.Sprintf("%04x", s.reportID),
		},
	}
}

func (s *WriterGun) createRequestInputs(encodedReport []byte, sigs [][]byte, repContext []byte) ([]byte, []byte, error) {
	validInputs, err := values.NewMap(map[string]any{
		"signed_report": map[string]any{
			"report":     encodedReport,
			"signatures": sigs,
			"context":    repContext,
			"id":         [2]byte{0, s.reportID},
		},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create inputs: %w", err)
	}

	validConfig, err := values.NewMap(map[string]any{
		"Address":  s.testParams.dataFeedsCacheAddress.String(),
		"gasLimit": 4000000,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create config: %w", err)
	}

	configBytes, err := mock_capability.MapToBytes(validConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert config to bytes: %w", err)
	}

	inputBytes, err := mock_capability.MapToBytes(validInputs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert inputs to bytes: %w", err)
	}

	return inputBytes, configBytes, nil
}
func stringTo32Byte(input string) ([32]byte, error) {
	var result [32]byte

	// Remove "0x" prefix if present
	input = strings.TrimPrefix(input, "0x")

	// Decode hex string to bytes
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return result, fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Check length
	if len(decoded) > 32 {
		return result, fmt.Errorf("input string too long: got %d bytes, want <= 32", len(decoded))
	}

	// Copy bytes to result, right-aligned
	copy(result[32-len(decoded):], decoded)
	return result, nil
}

func convertToHashedWorkflowName(input string) string {
	// Create SHA256 hash
	hash := sha256.New()
	hash.Write([]byte(input))

	// Get the hex string of the hash
	hashHex := hex.EncodeToString(hash.Sum(nil))

	return hex.EncodeToString([]byte(hashHex[:10]))
}

func saveClientURL(url string) error {
	// Create cache directory if it doesn't exist
	cacheDir := filepath.Join(os.TempDir(), "cache")
	if err := os.MkdirAll(cacheDir, 0600); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Save URL to file
	cacheFile := filepath.Join(cacheDir, "client_url.txt")
	if err := os.WriteFile(cacheFile, []byte(url), 0600); err != nil {
		return fmt.Errorf("failed to write client URL file: %w", err)
	}

	return nil
}

func loadClientURL() (string, error) {
	cacheFile := filepath.Join(os.TempDir(), "cache", "client_url.txt")
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return "", fmt.Errorf("failed to read client URL file: %w", err)
	}

	return string(data), nil
}
