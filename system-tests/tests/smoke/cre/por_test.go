package cre

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	writeevmcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	portypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v1/proof-of-reserve/cron-based/types"
)

var (
	SinglePoRDonCapabilitiesFlags = []string{cre.CronCapability, cre.OCR3Capability, cre.CustomComputeCapability, cre.WriteEVMCapability}
)

type CustomAnvilMiner struct {
	BlockSpeedSeconds int `toml:"block_speed_seconds"`
}

type TestConfig struct {
	Blockchains                   []*cre.WrappedBlockchainInput `toml:"blockchains" validate:"required"`
	CustomAnvilMiner              *CustomAnvilMiner             `toml:"custom_anvil_miner"`
	NodeSets                      []*ns.Input                   `toml:"nodesets" validate:"required"`
	WorkflowConfigs               []WorkflowConfig              `toml:"workflow_configs" validate:"required"`
	JD                            *jd.Input                     `toml:"jd" validate:"required"`
	Fake                          *fake.Input                   `toml:"fake"`
	WorkflowRegistryConfiguration *cre.WorkflowRegistryInput    `toml:"workflow_registry_configuration"`
	Infra                         *infra.Input                  `toml:"infra" validate:"required"`
	DependenciesConfig            *DependenciesConfig           `toml:"dependencies" validate:"required"`
}

type WorkflowConfig struct {
	// Tells the test where the workflow file is located
	WorkflowFileLocation string `toml:"workflow_file_location" validate:"required"`
	FeedID               string `toml:"feed_id" validate:"required,startsnotwith=0x"`
}

// Defines the location of the binary files that are required to run the test
// When test runs in CI hardcoded versions will be downloaded before the test starts
// Command that downloads them is part of "test_cmd" in .github/e2e-tests.yml file
type DependenciesConfig struct {
	CronCapabilityBinaryPath string `toml:"cron_capability_binary_path"`
}

const (
	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""
)

func validateEnvVars(t *testing.T) {
	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorImageEnvVarName), "missing env var: "+creenv.E2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorVersionEnvVarName), "missing env var: "+creenv.E2eJobDistributorVersionEnvVarName)
	}
}

type configureDataFeedsCacheInput struct {
	chainSelector      uint64
	fullCldEnvironment *cldf.Environment
	workflowName       string
	feedID             string
	sethClient         *seth.Client
	blockchain         *blockchain.Output
	deployerPrivateKey string
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

func createConfigFile(feedsConsumerAddress common.Address, workflowName, feedID, dataURL, writeTargetName string) (string, error) {
	cleanFeedID := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedID)

	if feedLength < 32 {
		return "", errors.Errorf("feed ID must be at least 32 characters long, but was %d", feedLength)
	}

	if feedLength > 32 {
		cleanFeedID = cleanFeedID[:32]
	}

	feedIDToUse := "0x" + cleanFeedID

	workflowConfig := portypes.WorkflowConfig{
		ComputeConfig: portypes.ComputeConfig{
			FeedID:                feedIDToUse,
			URL:                   dataURL,
			DataFeedsCacheAddress: feedsConsumerAddress.Hex(),
			WriteTargetName:       writeTargetName,
		},
	}

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

type porSetupOutput struct {
	priceProvider                   PriceProvider
	addressBook                     cldf.AddressBook
	chainSelectorToSethClient       map[uint64]*seth.Client
	chainSelectorToBlockchainOutput map[uint64]*blockchain.Output
	donTopology                     *cre.DonTopology
	nodeOutput                      []*cre.WrappedNodeOutput
	chainSelectorToWorkflowConfig   map[uint64]WorkflowConfig
}

func setupPoRTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	priceProvider PriceProvider,
	mustSetCapabilitiesFn func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *porSetupOutput {
	extraAllowedGatewayPorts := []int{}
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedGatewayPorts = append(extraAllowedGatewayPorts, in.Fake.Port)
	}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.Type)
	require.NoError(t, pathErr, "failed to get default container directory")
	var cronBinaryPathInTheContainer string
	if in.DependenciesConfig.CronCapabilityBinaryPath != "" {
		// where cron binary is located in the container
		cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.CronCapabilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[cre.CronCapability] = in.DependenciesConfig.CronCapabilityBinaryPath
	} else {
		// assume that if cron binary is already in the image it is in the default location and has default name
		cronBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
	}

	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  customBinariesPaths,
		JobSpecFactoryFunctions: []cre.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
			cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
			crecompute.ComputeJobSpecFactoryFn,
		},
		ConfigFactoryFunctions: []cre.ConfigFactoryFn{
			gatewayconfig.GenerateConfig,
		},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(t.Context(), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")
	homeChainOutput := universalSetupOutput.BlockchainOutput[0]

	if in.CustomAnvilMiner != nil {
		for _, bi := range universalSetupInput.BlockchainsInput {
			if bi.Type == blockchain.TypeAnvil {
				require.NotContains(t, bi.DockerCmdParamsOverrides, "-b", "custom_anvil_miner was specified but Anvil has '-b' key set, remove that parameter from 'docker_cmd_params' to run deployments instantly or remove custom_anvil_miner key from TOML config")
			}
		}
		for _, bo := range universalSetupOutput.BlockchainOutput {
			if bo.BlockchainOutput.Type == blockchain.TypeAnvil {
				miner := rpc.NewRemoteAnvilMiner(bo.BlockchainOutput.Nodes[0].ExternalHTTPUrl, nil)
				miner.MinePeriodically(time.Duration(in.CustomAnvilMiner.BlockSpeedSeconds) * time.Second)
			}
		}
	}

	chainSelectorToWorkflowConfig := make(map[uint64]WorkflowConfig)
	chainSelectorToSethClient := make(map[uint64]*seth.Client)
	chainSelectorToBlockchainOutput := make(map[uint64]*blockchain.Output)

	for idx, bo := range universalSetupOutput.BlockchainOutput {
		if bo.ReadOnly {
			continue
		}
		chainSelectorToWorkflowConfig[bo.ChainSelector] = in.WorkflowConfigs[idx]
		chainSelectorToSethClient[bo.ChainSelector] = bo.SethClient
		chainSelectorToBlockchainOutput[bo.ChainSelector] = bo.BlockchainOutput

		deployConfig := df_changeset_types.DeployConfig{
			ChainsToDeploy: []uint64{bo.ChainSelector},
			Labels:         []string{"data-feeds"}, // label required by the changeset
		}

		dfOutput, dfErr := changeset.RunChangeset(df_changeset.DeployCacheChangeset, *universalSetupOutput.CldEnvironment, deployConfig)
		require.NoError(t, dfErr, "failed to deploy data feed cache contract")

		mergeErr := universalSetupOutput.CldEnvironment.ExistingAddresses.Merge(dfOutput.AddressBook) //nolint:staticcheck // won't migrate now
		require.NoError(t, mergeErr, "failed to merge address book")

		workflowName := "por-workflow-" + strconv.FormatUint(bo.ChainID, 10)

		dfConfigInput := &configureDataFeedsCacheInput{
			chainSelector:      bo.ChainSelector,
			fullCldEnvironment: universalSetupOutput.CldEnvironment,
			workflowName:       workflowName,
			feedID:             in.WorkflowConfigs[idx].FeedID,
			sethClient:         bo.SethClient,
			blockchain:         bo.BlockchainOutput,
			deployerPrivateKey: bo.DeployerPrivateKey,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		testLogger.Info().Msg("Proceeding to register PoR workflow...")

		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(
			universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			homeChainOutput.ChainSelector,
			keystone_changeset.WorkflowRegistry.String(),
		)
		require.NoError(t, workflowRegistryErr, "failed to find workflow registry address for chain %d", bo.ChainSelector)

		dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(
			universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			bo.ChainSelector,
			df_changeset.DataFeedsCache.String(),
		)
		require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", bo.ChainSelector)

		workflowConfigFilePath, configErr := createConfigFile(dataFeedsCacheAddress, workflowName, in.WorkflowConfigs[idx].FeedID, priceProvider.URL(), corevm.GenerateWriteTargetName(bo.ChainID))
		require.NoError(t, configErr, "failed to create workflow config file")

		compressedWorkflowWasmPath, compileErr := creworkflow.CompileWorkflow(in.WorkflowConfigs[idx].WorkflowFileLocation, workflowName)
		require.NoError(t, compileErr, "failed to compile workflow '%s'", in.WorkflowConfigs[idx].WorkflowFileLocation)

		t.Cleanup(func() {
			_ = os.Remove(compressedWorkflowWasmPath)
			_ = os.Remove(workflowConfigFilePath)
		})

		containerTargetDir := "/home/chainlink/workflows"
		workflowCopyErr := creworkflow.CopyWorkflowToDockerContainers(compressedWorkflowWasmPath, "workflow-node", containerTargetDir)
		require.NoError(t, workflowCopyErr, "failed to copy workflow to docker containers")

		configCopyErr := creworkflow.CopyWorkflowToDockerContainers(workflowConfigFilePath, "workflow-node", containerTargetDir)
		require.NoError(t, configCopyErr, "failed to copy workflow config to docker containers")

		registerErr := creworkflow.RegisterWithContract(
			t.Context(),
			homeChainOutput.SethClient,
			workflowRegistryAddress,
			universalSetupOutput.DonTopology.WorkflowDonID,
			workflowName,
			"file://"+compressedWorkflowWasmPath,
			ptr.Ptr("file://"+workflowConfigFilePath),
			nil,
			&containerTargetDir,
		)
		require.NoError(t, registerErr, "failed to register PoR workflow")
	}
	// Workflow-specific configuration -- END

	return &porSetupOutput{
		priceProvider:                   priceProvider,
		chainSelectorToSethClient:       chainSelectorToSethClient,
		chainSelectorToBlockchainOutput: chainSelectorToBlockchainOutput,
		donTopology:                     universalSetupOutput.DonTopology,
		nodeOutput:                      universalSetupOutput.NodeOutput,
		addressBook:                     universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
		chainSelectorToWorkflowConfig:   chainSelectorToWorkflowConfig,
	}
}

func TestCRE_OCR3_PoR_Workflow_SingleDon_MultipleWriters_MockedPrice(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-one-don-multichain.toml")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{cre.WorkflowDON, cre.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	capabilityFactoryFns := []cre.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
	}

	for _, bc := range in.Blockchains {
		chainID, chainErr := strconv.Atoi(bc.ChainID)
		require.NoError(t, chainErr, "failed to convert chain ID to int")
		capabilityFactoryFns = append(capabilityFactoryFns, writeevmcap.WriteEVMCapabilityFactory(libc.MustSafeUint64(int64(chainID))))
	}

	setupOutput := setupPoRTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		capabilityFactoryFns,
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

func TestCRE_OCR3_PoR_Workflow_GatewayDon_MockedPrice(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-gateway-don.toml")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{cre.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{},
				DONTypes:           []string{cre.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                       // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, []string{in.WorkflowConfigs[0].FeedID})
	require.NoError(t, priceErr, "failed to create fake price provider")

	firstBlockchain := in.Blockchains[0]
	chainIDInt, chainErr := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []cre.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
		writeevmcap.WriteEVMCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt))),
	})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

func TestCRE_OCR3_PoR_Workflow_CapabilitiesDons_LivePrice(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-capabilities-don.toml")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 3, "expected 3 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{cre.OCR3Capability, cre.CustomComputeCapability, cre.CronCapability},
				DONTypes:           []string{cre.WorkflowDON},
				BootstrapNodeIndex: 0,
				SupportedChains:    []uint64{1337}, // workflow DON has to support only home chain
			},
			{
				Input:              input[1],
				Capabilities:       []string{cre.WriteEVMCapability},
				DONTypes:           []string{cre.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                            // <----- indicate that capabilities DON doesn't have a bootstrap node and will use the global bootstrap node
				SupportedChains:    []uint64{1337, 2337},          // capabilities DON has to support both chains, because we want to make sure that second workflow that writes to the second chain is run using a remote capability
			},
			{
				Input:              input[2],
				Capabilities:       []string{},
				DONTypes:           []string{cre.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                       // <----- it's crucial to indicate there's no bootstrap node for the gateway DON
				GatewayNodeIndex:   0,
			},
		}
	}

	// we want to register write EVM capability only for the second blockchain
	secondBlockchain := in.Blockchains[1]
	secondChainIDInt, secondChainErr := strconv.Atoi(secondBlockchain.ChainID)
	require.NoError(t, secondChainErr, "failed to convert chain ID to int")

	priceProvider := NewTrueUSDPriceProvider(testLogger, []string{in.WorkflowConfigs[0].FeedID})
	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []cre.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
		writeevmcap.WriteEVMCapabilityFactory(libc.MustSafeUint64(int64(secondChainIDInt))),
	})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		debugTest(t, testLogger, setupOutput, in)
	})

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
}

func waitForFeedUpdate(t *testing.T, testLogger zerolog.Logger, priceProvider PriceProvider, setupOutput *porSetupOutput, timeout time.Duration) {
	eg := &errgroup.Group{}
	for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
		eg.Go(func() error {
			testLogger.Info().Msgf("Waiting for feed %s to update...", workflowConfig.FeedID)

			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

			dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(dataFeedsCacheAddresses, setupOutput.chainSelectorToSethClient[chainSelector].Client)
			require.NoError(t, instanceErr, "failed to create data feeds cache instance")

			startTime := time.Now()
			assert.Eventually(t, func() bool {
				elapsed := time.Since(startTime).Round(time.Second)
				price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.chainSelectorToSethClient[chainSelector].NewCallOpts(), [16]byte(common.Hex2Bytes(workflowConfig.FeedID)))
				require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

				// if there are no more prices to be found, we can stop waiting
				return !setupOutput.priceProvider.NextPrice(workflowConfig.FeedID, price, elapsed)
			}, timeout, 10*time.Second, "feed %s did not update, timeout after: %s", workflowConfig.FeedID, timeout)

			expected := priceProvider.ExpectedPrices(workflowConfig.FeedID)
			actual := priceProvider.ActualPrices(workflowConfig.FeedID)

			if len(expected) != len(actual) {
				return errors.Errorf("expected %d prices, got %d", len(expected), len(actual))
			}

			for i := range expected {
				if expected[i].Cmp(actual[i]) != 0 {
					return errors.Errorf("expected price %d, got %d", expected[i], actual[i])
				}
			}

			testLogger.Info().Msgf("All %d prices were found in the feed %s", len(expected), workflowConfig.FeedID)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	testLogger.Info().Msgf("All prices were found for all feeds")
}

func debugTest(t *testing.T, testLogger zerolog.Logger, setupOutput *porSetupOutput, in *TestConfig) {
	if t.Failed() {
		counter := 0
		for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

			forwarderAddresses, forwarderErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, keystone_changeset.KeystoneForwarder.String())
			require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", chainSelector)

			logTestInfo(testLogger, workflowConfig.FeedID, dataFeedsCacheAddresses.Hex(), forwarderAddresses.Hex())
			counter++
			// log scanning is not supported for CRIB
			if in.Infra.Type == infra.CRIB {
				return
			}

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

			debugDons := make([]*cre.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
			for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range setupOutput.nodeOutput[i].CLNodes {
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
				BlockchainOutput: setupOutput.chainSelectorToBlockchainOutput[chainSelector],
				InfraInput:       in.Infra,
			}
			lidebug.PrintTestDebug(t.Context(), t.Name(), testLogger, debugInput)
		}
	}
}

func setConfigurationIfMissing(configName string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	if os.Getenv("PRIVATE_KEY") == "" {
		err := os.Setenv("PRIVATE_KEY", blockchain.DefaultAnvilPrivateKey)
		if err != nil {
			return errors.Wrap(err, "failed to set PRIVATE_KEY env var")
		}
	}

	return nil
}
