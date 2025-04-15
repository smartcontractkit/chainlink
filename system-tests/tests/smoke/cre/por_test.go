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
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink/deployment"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/por"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

var (
	SinglePoRDonCapabilitiesFlags = []string{"ocr3", "cron", "custom-compute", "write-evm"}
)

type CustomAnvilMiner struct {
	BlockSpeedSeconds int `toml:"block_speed_seconds"`
}

type TestConfig struct {
	BlockchainA                   *blockchain.Input                        `toml:"blockchain_a" validate:"required"`
	CustomAnvilMiner              *CustomAnvilMiner                        `toml:"custom_anvil_miner"`
	NodeSets                      []*ns.Input                              `toml:"nodesets" validate:"required"`
	WorkflowConfig                *WorkflowConfig                          `toml:"workflow_config" validate:"required"`
	JD                            *jd.Input                                `toml:"jd" validate:"required"`
	Fake                          *fake.Input                              `toml:"fake"`
	KeystoneContracts             *keystonetypes.KeystoneContractsInput    `toml:"keystone_contracts"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput     `toml:"workflow_registry_configuration"`
	DataFeedsCacheContract        *keystonetypes.DeployDataFeedsCacheInput `toml:"data_feeds_cache"`
	Infra                         *libtypes.InfraInput                     `toml:"infra" validate:"required"`
}

type WorkflowConfig struct {
	UseCRECLI bool `toml:"use_cre_cli"`
	/*
		These tests can be run in two modes:
		1. existing mode: it uses a workflow binary (and configuration) file that is already uploaded to Gist
		2. compile mode: it compiles a new workflow binary and uploads it to Gist

		For the "compile" mode to work, the `GIST_WRITE_TOKEN` env var must be set to a token that has `gist:read` and `gist:write` permissions, but this permissions
		are tied to account not to repository. Currently, we have no service account in the CI at all. And using a token that's tied to personal account of a developer
		is not a good idea. So, for now, we are only allowing the `existing` mode in CI.

		If you wish to use "compile" mode set `ShouldCompileNewWorkflow` to `true`, set `GIST_WRITE_TOKEN` env var and provide the path to the workflow folder.
	*/
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow" validate:"no_cre_no_compilation,disabled_in_ci"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string             `toml:"workflow_folder_location" validate:"required_if=ShouldCompileNewWorkflow true"`
	CompiledWorkflowConfig *CompiledConfig     `toml:"compiled_config" validate:"required_if=ShouldCompileNewWorkflow false"`
	DependenciesConfig     *DependenciesConfig `toml:"dependencies" validate:"required"`
	WorkflowName           string              `toml:"workflow_name" validate:"required" `
	FeedID                 string              `toml:"feed_id" validate:"required,startsnotwith=0x"`
}

// noCRENoCompilation is a custom validator for the tag "no_cre_no_compilation".
// It ensures that if UseCRECLI is false, then ShouldCompileNewWorkflow must also be false.
func noCRENoCompilation(fl validator.FieldLevel) bool {
	// Use Parent() to access the WorkflowConfig struct.
	wc, ok := fl.Parent().Interface().(WorkflowConfig)
	if !ok {
		return false
	}
	// If not using CRE CLI and ShouldCompileNewWorkflow is true, fail validation.
	if !wc.UseCRECLI && fl.Field().Bool() {
		return false
	}
	return true
}

func disabledInCI(fl validator.FieldLevel) bool {
	if os.Getenv("CI") == "true" {
		return !fl.Field().Bool()
	}

	return true
}

func registerNoCRENoCompilationTranslation(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("no_cre_no_compilation", trans, func(ut ut.Translator) error {
		return ut.Add("no_cre_no_compilation", "{0} must be false when UseCRECLI is false, it is not possible to compile a workflow without it", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("no_cre_no_compilation", fe.Field())
		return t
	})
}

func registerNoFolderLocationTranslation(v *validator.Validate, trans ut.Translator) {
	_ = v.RegisterTranslation("folder_required_if_compiling", trans, func(ut ut.Translator) error {
		return ut.Add("folder_required_if_compiling", "{0} must set, when compiling the workflow", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("folder_required_if_compiling", fe.Field())
		return t
	})
}

func init() {
	err := framework.Validator.RegisterValidation("no_cre_no_compilation", noCRENoCompilation)
	if err != nil {
		panic(errors.Wrap(err, "failed to register no_cre_no_compilation validator"))
	}
	err = framework.Validator.RegisterValidation("disabled_in_ci", disabledInCI)
	if err != nil {
		panic(errors.Wrap(err, "failed to register disabled_in_ci validator"))
	}

	if framework.ValidatorTranslator != nil {
		registerNoCRENoCompilationTranslation(framework.Validator, framework.ValidatorTranslator)
		registerNoFolderLocationTranslation(framework.Validator, framework.ValidatorTranslator)
	}
}

// Defines the location of the binary files that are required to run the test
// When test runs in CI hardcoded versions will be downloaded before the test starts
// Command that downloads them is part of "test_cmd" in .github/e2e-tests.yml file
type DependenciesConfig struct {
	CronCapabilityBinaryPath string `toml:"cron_capability_binary_path"`
	CRECLIBinaryPath         string `toml:"cre_cli_binary_path" validate:"required"`
}

const (
	CronBinaryVersion   = "v1.0.2-alpha"
	CRECLIBinaryVersion = "v0.1.5"
)

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL string `toml:"binary_url" validate:"required"`
	ConfigURL string `toml:"config_url" validate:"required"`
}

func validateEnvVars(t *testing.T, in *TestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")

	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorImageEnvVarName), "missing env var: "+creenv.E2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(creenv.E2eJobDistributorVersionEnvVarName), "missing env var: "+creenv.E2eJobDistributorVersionEnvVarName)
	}

	if in.WorkflowConfig.UseCRECLI {
		if in.WorkflowConfig.ShouldCompileNewWorkflow {
			gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
			require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
			err := os.Setenv("CRE_GITHUB_API_TOKEN", gistWriteToken)
			require.NoError(t, err, "failed to set CRE_GITHUB_API_TOKEN env var")
		}
	}
}

type registerPoRWorkflowInput struct {
	*WorkflowConfig
	chainSelector           uint64
	writeTargetName         string
	workflowDonID           uint32
	feedID                  string
	workflowRegistryAddress common.Address
	dataFeedsCacheAddress   common.Address
	priceProvider           PriceProvider
	sethClient              *seth.Client
	deployerPrivateKey      string
	creCLIAbsPath           string
	creCLIsettingsFile      *os.File
}

type configureDataFeedsCacheInput struct {
	useCRECLI             bool
	chainSelector         uint64
	fullCldEnvironment    *deployment.Environment
	forwarderAddress      common.Address
	dataFeedsCacheAddress common.Address
	workflowName          string
	feedID                string
	sethClient            *seth.Client
	blockchain            *blockchain.Output
	creCLIAbsPath         string
	settingsFile          *os.File
	deployerPrivateKey    string
}

func configureDataFeedsCacheContract(testLogger zerolog.Logger, input *configureDataFeedsCacheInput) error {
	chainIDInt, intErr := strconv.Atoi(input.blockchain.ChainID)
	if intErr != nil {
		return errors.Wrap(intErr, "failed to convert chain ID to int")
	}

	if input.useCRECLI {
		// These two env vars are required by the CRE CLI
		err := os.Setenv("CRE_ETH_PRIVATE_KEY", input.deployerPrivateKey)
		if err != nil {
			return errors.Wrap(err, "failed to set CRE_ETH_PRIVATE_KEY")
		}

		dfAdminErr := libcrecli.SetFeedAdmin(input.creCLIAbsPath, chainIDInt, input.sethClient.MustGetRootKeyAddress(), input.settingsFile)
		if dfAdminErr != nil {
			return errors.Wrap(dfAdminErr, "failed to set feed admin")
		}

		cleanFeedID := strings.TrimPrefix(input.feedID, "0x")

		// Ensure the feed ID is long enough
		if len(cleanFeedID) < 14 { // Need at least 7 bytes (14 hex chars)
			return fmt.Errorf("feed ID too short: %s", input.feedID)
		} else if len(cleanFeedID) > 32 {
			cleanFeedID = cleanFeedID[:32]
		}

		// Extract decimals from feed ID
		decimals, decimalsErr := df_changeset.GetDecimalsFromFeedID(cleanFeedID)
		if decimalsErr != nil {
			return errors.Wrapf(decimalsErr, "failed to get decimals from feed ID %s", input.feedID)
		}

		dfConfigErr := libcrecli.SetFeedConfig(input.creCLIAbsPath, input.feedID, strconv.Itoa(int(decimals)), "PoR test feed", chainIDInt, []common.Address{input.forwarderAddress}, []common.Address{input.sethClient.MustGetRootKeyAddress()}, []string{input.workflowName}, input.settingsFile)
		if dfConfigErr != nil {
			return errors.Wrap(dfConfigErr, "failed to set feed config")
		}

		return nil
	}

	configInput := &keystonetypes.ConfigureDataFeedsCacheInput{
		CldEnv:                input.fullCldEnvironment,
		ChainSelector:         input.chainSelector,
		FeedIDs:               []string{input.feedID},
		Descriptions:          []string{"PoR test feed"},
		DataFeedsCacheAddress: input.dataFeedsCacheAddress,
		AdminAddress:          input.sethClient.MustGetRootKeyAddress(),
		AllowedSenders:        []common.Address{input.forwarderAddress},
		AllowedWorkflowNames:  []string{input.workflowName},
		AllowedWorkflowOwners: []common.Address{input.sethClient.MustGetRootKeyAddress()},
	}

	_, configErr := libcontracts.ConfigureDataFeedsCache(testLogger, configInput)

	return configErr
}

func registerPoRWorkflow(input registerPoRWorkflowInput) error {
	// Register workflow directly using the provided binary URL and optionally config and secrets URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !input.WorkflowConfig.ShouldCompileNewWorkflow && !input.WorkflowConfig.UseCRECLI {
		err := libcontracts.RegisterWorkflow(
			input.sethClient,
			input.workflowRegistryAddress,
			input.workflowDonID,
			input.WorkflowConfig.WorkflowName,
			input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			&input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
			nil, // TODO pass secrets URL once support for them has been added
		)
		if err != nil {
			return errors.Wrap(err, "failed to register workflow")
		}

		return nil
	}

	// This env var is required by the CRE CLI
	err := os.Setenv("CRE_ETH_PRIVATE_KEY", input.deployerPrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to set CRE_ETH_PRIVATE_KEY")
	}

	// create workflow-specific config file
	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(input.dataFeedsCacheAddress, input.feedID, input.priceProvider.URL(), input.writeTargetName)
	if configErr != nil {
		return errors.Wrap(configErr, "failed to create workflow config file")
	}

	workflowConfigFilePath := workflowConfigFile.Name()

	registerWorkflowInput := keystonetypes.RegisterWorkflowWithCRECLIInput{
		ChainSelector:            input.chainSelector,
		WorkflowDonID:            input.workflowDonID,
		WorkflowRegistryAddress:  input.workflowRegistryAddress,
		WorkflowOwnerAddress:     input.sethClient.MustGetRootKeyAddress(),
		CRECLIPrivateKey:         input.deployerPrivateKey,
		CRECLIAbsPath:            input.creCLIAbsPath,
		CRESettingsFile:          input.creCLIsettingsFile,
		WorkflowName:             input.WorkflowConfig.WorkflowName,
		ShouldCompileNewWorkflow: input.WorkflowConfig.ShouldCompileNewWorkflow,
	}

	if input.WorkflowConfig.ShouldCompileNewWorkflow {
		registerWorkflowInput.NewWorkflow = &keystonetypes.NewWorkflow{
			FolderLocation: *input.WorkflowConfig.WorkflowFolderLocation,
			ConfigFilePath: &workflowConfigFilePath,
		}
	} else {
		registerWorkflowInput.ExistingWorkflow = &keystonetypes.ExistingWorkflow{
			BinaryURL: input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			ConfigURL: &input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
		}
	}

	registerErr := creworkflow.RegisterWithCRECLI(registerWorkflowInput)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow with CRE CLI")
	}

	return nil
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, dataFeedsCacheAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("DataFeedsCache address: %s", dataFeedsCacheAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

type porSetupOutput struct {
	priceProvider         PriceProvider
	dataFeedsCacheAddress common.Address
	forwarderAddress      common.Address
	sethClient            *seth.Client
	blockchainOutput      *blockchain.Output
	donTopology           *keystonetypes.DonTopology
	nodeOutput            []*keystonetypes.WrappedNodeOutput
}

func setupPoRTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	priceProvider PriceProvider,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *porSetupOutput {
	extraAllowedPorts := []int{}
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedPorts = append(extraAllowedPorts, in.Fake.Port)
	}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	var cronBinaryPathInTheContainer string
	if in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath != "" {
		require.NoError(t, pathErr, "failed to get default container directory")
		// where cron binary is located in the container
		cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[keystonetypes.CronCapability] = in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath
	} else {
		// assume that if cron binary is already in the image it is in the default location and has default name
		cronBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
	}

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:  mustSetCapabilitiesFn(in.NodeSets),
		CapabilityFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:           *in.BlockchainA,
		JdInput:                    *in.JD,
		InfraInput:                 *in.Infra,
		CustomBinariesPaths:        customBinariesPaths,
		ExtraAllowedPorts:          extraAllowedPorts,
		JobSpecFactoryFunctions:    []keystonetypes.JobSpecFactoryFn{keystonepor.PoRJobSpecFactoryFn(cronBinaryPathInTheContainer, extraAllowedPorts, []string{}, []string{"0.0.0.0/0"})},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")

	if in.CustomAnvilMiner != nil {
		require.NotContains(t, in.BlockchainA.DockerCmdParamsOverrides, "-b", "custom_anvil_miner was specified but Anvil has '-b' key set, remove that parameter from 'docker_cmd_params' to run deployments instantly or remove custom_anvil_miner key from TOML config")
		require.Equal(t, "anvil", in.BlockchainA.Type, "custom_anvil_miner was specified but blockchain type is not Anvil")
		miner := rpc.NewRemoteAnvilMiner(universalSetupOutput.BlockchainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl, nil)
		miner.MinePeriodically(time.Duration(in.CustomAnvilMiner.BlockSpeedSeconds) * time.Second)
	}

	deployDataFeedsInput := &keystonetypes.DeployDataFeedsCacheInput{
		ChainSelector: universalSetupOutput.BlockchainOutput.ChainSelector,
		CldEnv:        universalSetupOutput.CldEnvironment,
	}
	deployDataFeedsCacheOutput, dfErr := libcontracts.DeployDataFeedsCache(testLogger, deployDataFeedsInput)
	require.NoError(t, dfErr, "failed to deploy data feeds cache")

	var creCLIAbsPath string
	var creCLISettingsFile *os.File
	if in.WorkflowConfig.UseCRECLI {
		// make sure that path is indeed absolute
		var pathErr error
		creCLIAbsPath, pathErr = filepath.Abs(in.WorkflowConfig.DependenciesConfig.CRECLIBinaryPath)
		require.NoError(t, pathErr, "failed to get absolute path for CRE CLI")

		// create CRE CLI settings file
		var settingsErr error
		creCLISettingsFile, settingsErr = libcrecli.PrepareCRECLISettingsFile(
			universalSetupOutput.BlockchainOutput.SethClient.MustGetRootKeyAddress(),
			universalSetupOutput.KeystoneContractsOutput.CapabilitiesRegistryAddress,
			universalSetupOutput.KeystoneContractsOutput.WorkflowRegistryAddress,
			&deployDataFeedsCacheOutput.DataFeedsCacheAddress,
			universalSetupOutput.DonTopology.WorkflowDonID,
			universalSetupOutput.BlockchainOutput.ChainSelector,
			universalSetupOutput.BlockchainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl)
		require.NoError(t, settingsErr, "failed to create CRE CLI settings file")
	}

	dfConfigInput := &configureDataFeedsCacheInput{
		useCRECLI:             in.WorkflowConfig.UseCRECLI,
		chainSelector:         universalSetupOutput.BlockchainOutput.ChainSelector,
		fullCldEnvironment:    universalSetupOutput.CldEnvironment,
		forwarderAddress:      universalSetupOutput.KeystoneContractsOutput.ForwarderAddress,
		dataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		workflowName:          in.WorkflowConfig.WorkflowName,
		feedID:                in.WorkflowConfig.FeedID,
		sethClient:            universalSetupOutput.BlockchainOutput.SethClient,
		blockchain:            universalSetupOutput.BlockchainOutput.BlockchainOutput,
		creCLIAbsPath:         creCLIAbsPath,
		settingsFile:          creCLISettingsFile,
		deployerPrivateKey:    universalSetupOutput.BlockchainOutput.DeployerPrivateKey,
	}
	dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
	require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

	registerInput := registerPoRWorkflowInput{
		WorkflowConfig:          in.WorkflowConfig,
		chainSelector:           universalSetupOutput.BlockchainOutput.ChainSelector,
		workflowDonID:           universalSetupOutput.DonTopology.WorkflowDonID,
		feedID:                  in.WorkflowConfig.FeedID,
		workflowRegistryAddress: universalSetupOutput.KeystoneContractsOutput.WorkflowRegistryAddress,
		dataFeedsCacheAddress:   deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		priceProvider:           priceProvider,
		sethClient:              universalSetupOutput.BlockchainOutput.SethClient,
		deployerPrivateKey:      universalSetupOutput.BlockchainOutput.DeployerPrivateKey,
		creCLIAbsPath:           creCLIAbsPath,
		creCLIsettingsFile:      creCLISettingsFile,
		writeTargetName:         corevm.GenerateWriteTargetName(universalSetupOutput.BlockchainOutput.ChainID),
	}

	workflowErr := registerPoRWorkflow(registerInput)
	require.NoError(t, workflowErr, "failed to register PoR workflow")
	// Workflow-specific configuration -- END

	// Set inputs in the test config, so that they can be saved
	in.KeystoneContracts = &keystonetypes.KeystoneContractsInput{
		Out: universalSetupOutput.KeystoneContractsOutput,
	}
	in.DataFeedsCacheContract = &keystonetypes.DeployDataFeedsCacheInput{
		Out: &keystonetypes.DeployDataFeedsCacheOutput{
			DataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		},
	}
	in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{
		Out: universalSetupOutput.WorkflowRegistryConfigurationOutput,
	}

	return &porSetupOutput{
		priceProvider:         priceProvider,
		dataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
		forwarderAddress:      universalSetupOutput.KeystoneContractsOutput.ForwarderAddress,
		sethClient:            universalSetupOutput.BlockchainOutput.SethClient,
		blockchainOutput:      universalSetupOutput.BlockchainOutput.BlockchainOutput,
		donTopology:           universalSetupOutput.DonTopology,
		nodeOutput:            universalSetupOutput.NodeOutput,
	}
}

// config file to use: environment-one-don.toml
func TestCRE_OCR3_PoR_Workflow_SingleDon_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake)
	require.NoError(t, priceErr, "failed to create fake price provider")

	chainIDInt, chainErr := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupPoRTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))},
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.dataFeedsCacheAddress.Hex(), setupOutput.forwarderAddress.Hex())

			// log scanning is not supported for CRIB
			if in.Infra.InfraType == libtypes.CRIB {
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

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(setupOutput.dataFeedsCacheAddress, setupOutput.sethClient.Client)
	require.NoError(t, instanceErr, "failed to create data feeds cache instance")

	startTime := time.Now()
	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.sethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(in.WorkflowConfig.FeedID)))
		require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

		// if there are no more prices to be found, we can stop waiting
		return !setupOutput.priceProvider.NextPrice(price, elapsed)
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "prices do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}

// config file to use: environment-gateway-don.toml
func TestCRE_OCR3_PoR_Workflow_GatewayDon_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                 // <----- it's crucial to indicate there's no bootstrap node
				GatewayNodeIndex:   0,
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake)
	require.NoError(t, priceErr, "failed to create fake price provider")

	chainIDInt, chainErr := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.dataFeedsCacheAddress.Hex(), setupOutput.forwarderAddress.Hex())

			// log scanning is not supported for CRIB
			if in.Infra.InfraType == libtypes.CRIB {
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

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(setupOutput.dataFeedsCacheAddress, setupOutput.sethClient.Client)
	require.NoError(t, instanceErr, "failed to create data feeds cache instance")

	startTime := time.Now()
	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.sethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(in.WorkflowConfig.FeedID)))
		require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

		// if there are no more prices to be found, we can stop waiting
		return !setupOutput.priceProvider.NextPrice(price, elapsed)
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "pricesup do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}

// config file to use: environment-capabilities-don.toml
func TestCRE_OCR3_PoR_Workflow_CapabilitiesDons_LivePrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 3, "expected 3 node sets in the test config")

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.OCR3Capability, keystonetypes.CustomComputeCapability, keystonetypes.CronCapability},
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.WriteEVMCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                      // <----- indicate that capabilities DON doesn't have a bootstrap node and will use the global bootstrap node
			},
			{
				Input:              input[2],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: -1,                                 // <----- it's crucial to indicate there's no bootstrap node for the gateway DON
				GatewayNodeIndex:   0,
			},
		}
	}

	chainIDInt, chainErr := strconv.Atoi(in.BlockchainA.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	priceProvider := NewTrueUSDPriceProvider(testLogger)
	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{libcontracts.DefaultCapabilityFactoryFn, libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(chainIDInt)))})

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.dataFeedsCacheAddress.Hex(), setupOutput.forwarderAddress.Hex())

			// log scanning is not supported for CRIB
			if in.Infra.InfraType == libtypes.CRIB {
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

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(setupOutput.dataFeedsCacheAddress, setupOutput.sethClient.Client)
	require.NoError(t, instanceErr, "failed to create data feeds cache instance")

	startTime := time.Now()
	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.sethClient.NewCallOpts(), [16]byte(common.Hex2Bytes(in.WorkflowConfig.FeedID)))
		require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

		// if there are no more prices to be found, we can stop waiting
		return !setupOutput.priceProvider.NextPrice(price, elapsed)
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "prices do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}
