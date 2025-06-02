package cre

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	df_changeset_types "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	writeevmcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	gatewayconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/gateway"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	creworkflow "github.com/smartcontractkit/chainlink/system-tests/lib/cre/workflow"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

var (
	SinglePoRDonCapabilitiesFlags = []string{keystonetypes.CronCapability, keystonetypes.OCR3Capability, keystonetypes.CustomComputeCapability, keystonetypes.WriteEVMCapability}
)

type CustomAnvilMiner struct {
	BlockSpeedSeconds int `toml:"block_speed_seconds"`
}

type TestConfig struct {
	Blockchains                   []*blockchain.Input                  `toml:"blockchains" validate:"required"`
	CustomAnvilMiner              *CustomAnvilMiner                    `toml:"custom_anvil_miner"`
	NodeSets                      []*ns.Input                          `toml:"nodesets" validate:"required"`
	WorkflowConfigs               []WorkflowConfig                     `toml:"workflow_configs" validate:"required"`
	JD                            *jd.Input                            `toml:"jd" validate:"required"`
	Fake                          *fake.Input                          `toml:"fake"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput `toml:"workflow_registry_configuration"`
	Infra                         *libtypes.InfraInput                 `toml:"infra" validate:"required"`
	DependenciesConfig            *DependenciesConfig                  `toml:"dependencies" validate:"required"`
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
	WorkflowFolderLocation *string         `toml:"workflow_folder_location" validate:"required_if=ShouldCompileNewWorkflow true"`
	CompiledWorkflowConfig *CompiledConfig `toml:"compiled_config" validate:"required_if=ShouldCompileNewWorkflow false"`
	WorkflowName           string          `toml:"workflow_name" validate:"required" `
	FeedID                 string          `toml:"feed_id" validate:"required,startsnotwith=0x"`
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
	AuthorizationKeySecretName = "AUTH_KEY"
	// TODO: use once we can run these tests in CI (https://smartcontract-it.atlassian.net/browse/DX-589)
	// AuthorizationKey           = "12a-281j&@91.sj1:_}"
	AuthorizationKey = ""
)

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL  string `toml:"binary_url" validate:"required"`
	ConfigURL  string `toml:"config_url" validate:"required"`
	SecretsURL string `toml:"secrets_url"`
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

	for _, workflowConfig := range in.WorkflowConfigs {
		if workflowConfig.UseCRECLI {
			if workflowConfig.ShouldCompileNewWorkflow {
				gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
				require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
				err := os.Setenv("CRE_GITHUB_API_TOKEN", gistWriteToken)
				require.NoError(t, err, "failed to set CRE_GITHUB_API_TOKEN env var")

				// set it only for the first workflow config, since it will be used for all workflows
				break
			}
		}
	}
}

type managePoRWorkflowInput struct {
	WorkflowConfig
	chainSelector      uint64
	homeChainSelector  uint64
	writeTargetName    string
	workflowDonID      uint32
	feedID             string
	addressBook        cldf.AddressBook
	priceProvider      PriceProvider
	sethClient         *seth.Client
	deployerPrivateKey string
	creCLIAbsPath      string
	creCLIsettingsFile *os.File
	authKey            string
	creCLIProfile      string
}

type configureDataFeedsCacheInput struct {
	useCRECLI          bool
	chainSelector      uint64
	fullCldEnvironment *cldf.Environment
	workflowName       string
	feedID             string
	sethClient         *seth.Client
	blockchain         *blockchain.Output
	creCLIAbsPath      string
	settingsFile       *os.File
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

	configInput := &keystonetypes.ConfigureDataFeedsCacheInput{
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

	_, configErr := libcontracts.ConfigureDataFeedsCache(testLogger, configInput)

	return configErr
}

func buildManageWorkflowInput(input managePoRWorkflowInput) (keystonetypes.ManageWorkflowWithCRECLIInput, error) {
	return keystonetypes.ManageWorkflowWithCRECLIInput{
		ChainSelector:            input.chainSelector,
		WorkflowDonID:            input.workflowDonID,
		WorkflowOwnerAddress:     input.sethClient.MustGetRootKeyAddress(),
		CRECLIPrivateKey:         input.deployerPrivateKey,
		CRECLIAbsPath:            input.creCLIAbsPath,
		CRESettingsFile:          input.creCLIsettingsFile,
		WorkflowName:             input.WorkflowConfig.WorkflowName,
		ShouldCompileNewWorkflow: input.WorkflowConfig.ShouldCompileNewWorkflow,
		CRECLIProfile:            input.creCLIProfile,
	}, nil
}

func pausePoRWorkflow(input managePoRWorkflowInput) error {
	workflowInput, err := buildManageWorkflowInput(input)
	if err != nil {
		return err
	}

	pauseErr := creworkflow.PauseWithCRECLI(workflowInput)
	if pauseErr != nil {
		return errors.Wrap(pauseErr, "failed to pause workflow with CRE CLI")
	}

	return nil
}

func activatePoRWorkflow(input managePoRWorkflowInput) error {
	workflowInput, err := buildManageWorkflowInput(input)
	if err != nil {
		return err
	}

	activateErr := creworkflow.ActivateWithCRECLI(workflowInput)
	if activateErr != nil {
		return errors.Wrap(activateErr, "failed to activate workflow with CRE CLI")
	}

	return nil
}

func registerPoRWorkflow(ctx context.Context, input managePoRWorkflowInput) error {
	// Register workflow directly using the provided binary URL and optionally config and secrets URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !input.WorkflowConfig.ShouldCompileNewWorkflow && !input.WorkflowConfig.UseCRECLI {
		workflowRegistryAddress, workflowRegistryErr := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, keystone_changeset.WorkflowRegistry.String())
		if workflowRegistryErr != nil {
			return errors.Wrapf(workflowRegistryErr, "failed to find workflow registry address for chain %d", input.chainSelector)
		}

		err := creworkflow.RegisterWithContract(
			ctx,
			input.sethClient,
			workflowRegistryAddress,
			input.workflowDonID,
			input.WorkflowConfig.WorkflowName,
			input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			&input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
			&input.WorkflowConfig.CompiledWorkflowConfig.SecretsURL,
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
	var secretNameToUse *string
	if input.authKey != "" {
		secretNameToUse = ptr.Ptr(AuthorizationKeySecretName)
	}

	dataFeedsCacheAddress, dataFeedsCacheErr := crecontracts.FindAddressesForChain(input.addressBook, input.chainSelector, df_changeset.DataFeedsCache.String())
	if dataFeedsCacheErr != nil {
		return errors.Wrapf(dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", input.chainSelector)
	}

	// pass nil if no secrets are used, otherwise workflow will fail if it cannot find the secret
	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(dataFeedsCacheAddress, input.feedID, input.priceProvider.URL(), input.writeTargetName, secretNameToUse)
	if configErr != nil {
		return errors.Wrap(configErr, "failed to create workflow config file")
	}
	workflowConfigFilePath := workflowConfigFile.Name()

	// indicate to the CRE CLI that the secret will be shared between all nodes in the workflow by using specific suffix
	authKeyEnvVarName := AuthorizationKeySecretName + libcrecli.SharedSecretEnvVarSuffix

	var secretsFilePath *string
	if input.authKey != "" {
		// create workflow-specific secrets file using the CRE CLI, which contains a mapping of secret names to environment variables that hold them
		// secrets will be read from the environment variables by the CRE CLI and encoded using nodes' public keys and when workflow executes it will
		// be able to read all secrets, which after decoding will be set as environment variables with names specified in the secrets file
		secrets := map[string][]string{
			AuthorizationKeySecretName: {authKeyEnvVarName},
		}

		secretsFile, secretsErr := libcrecli.CreateSecretsFile(secrets)
		if secretsErr != nil {
			return errors.Wrap(secretsErr, "failed to create secrets file")
		}
		secretsFilePath = ptr.Ptr(secretsFile.Name())
	}

	registerWorkflowInput := keystonetypes.ManageWorkflowWithCRECLIInput{
		ChainSelector:            input.chainSelector,
		WorkflowDonID:            input.workflowDonID,
		WorkflowOwnerAddress:     input.sethClient.MustGetRootKeyAddress(),
		CRECLIPrivateKey:         input.deployerPrivateKey,
		CRECLIAbsPath:            input.creCLIAbsPath,
		CRESettingsFile:          input.creCLIsettingsFile,
		WorkflowName:             input.WorkflowConfig.WorkflowName,
		ShouldCompileNewWorkflow: input.WorkflowConfig.ShouldCompileNewWorkflow,
		CRECLIProfile:            input.creCLIProfile,
	}

	if input.WorkflowConfig.ShouldCompileNewWorkflow {
		registerWorkflowInput.NewWorkflow = &keystonetypes.NewWorkflow{
			FolderLocation:   *input.WorkflowConfig.WorkflowFolderLocation,
			WorkflowFileName: "main.go",
			ConfigFilePath:   &workflowConfigFilePath,
			SecretsFilePath:  secretsFilePath,
			Secrets: map[string]string{
				authKeyEnvVarName: input.authKey,
			},
		}
	} else {
		registerWorkflowInput.ExistingWorkflow = &keystonetypes.ExistingWorkflow{
			BinaryURL:  input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL,
			ConfigURL:  &input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL,
			SecretsURL: &input.WorkflowConfig.CompiledWorkflowConfig.SecretsURL,
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
	priceProvider                   PriceProvider
	addressBook                     cldf.AddressBook
	chainSelectorToSethClient       map[uint64]*seth.Client
	chainSelectorToBlockchainOutput map[uint64]*blockchain.Output
	donTopology                     *keystonetypes.DonTopology
	nodeOutput                      []*keystonetypes.WrappedNodeOutput
	chainSelectorToWorkflowConfig   map[uint64]WorkflowConfig
}

func setupPoRTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	priceProvider PriceProvider,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *porSetupOutput {
	extraAllowedGatewayPorts := []int{}
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedGatewayPorts = append(extraAllowedGatewayPorts, in.Fake.Port)
	}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")
	var cronBinaryPathInTheContainer string
	if in.DependenciesConfig.CronCapabilityBinaryPath != "" {
		// where cron binary is located in the container
		cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.DependenciesConfig.CronCapabilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[keystonetypes.CronCapability] = in.DependenciesConfig.CronCapabilityBinaryPath
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
		JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
			cregateway.GatewayJobSpecFactoryFn(extraAllowedGatewayPorts, []string{}, []string{"0.0.0.0/0"}),
			crecompute.ComputeJobSpecFactoryFn,
		},
		ConfigFactoryFunctions: []keystonetypes.ConfigFactoryFn{
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

		var creCLIAbsPath string
		var creCLISettingsFile *os.File
		if in.WorkflowConfigs[idx].UseCRECLI {
			// make sure that path is indeed absolute
			var pathErr error
			creCLIAbsPath, pathErr = filepath.Abs(in.DependenciesConfig.CRECLIBinaryPath)
			require.NoError(t, pathErr, "failed to get absolute path for CRE CLI")

			rpcs := map[uint64]string{}
			for _, bcOut := range universalSetupOutput.BlockchainOutput {
				rpcs[bcOut.ChainSelector] = bcOut.BlockchainOutput.Nodes[0].ExternalHTTPUrl
			}

			// create CRE CLI settings file
			var settingsErr error
			creCLISettingsFile, settingsErr = libcrecli.PrepareCRECLISettingsFile(
				libcrecli.CRECLIProfile,
				bo.SethClient.MustGetRootKeyAddress(),
				universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
				universalSetupOutput.DonTopology.WorkflowDonID,
				homeChainOutput.ChainSelector,
				rpcs,
			)
			require.NoError(t, settingsErr, "failed to create CRE CLI settings file")
		}

		dfConfigInput := &configureDataFeedsCacheInput{
			useCRECLI:          in.WorkflowConfigs[idx].UseCRECLI,
			chainSelector:      bo.ChainSelector,
			fullCldEnvironment: universalSetupOutput.CldEnvironment,
			workflowName:       in.WorkflowConfigs[idx].WorkflowName,
			feedID:             in.WorkflowConfigs[idx].FeedID,
			sethClient:         bo.SethClient,
			blockchain:         bo.BlockchainOutput,
			creCLIAbsPath:      creCLIAbsPath,
			settingsFile:       creCLISettingsFile,
			deployerPrivateKey: bo.DeployerPrivateKey,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		testLogger.Info().Msg("Waiting for RegistrySyncer health checks...")
		syncerErr := waitForWorkflowRegistrySyncer(universalSetupOutput.NodeOutput, universalSetupOutput.DonTopology)
		require.NoError(t, syncerErr, "failed to wait for workflow registry syncer")
		testLogger.Info().Msg("Proceeding to register PoR workflow...")

		workflowInput := managePoRWorkflowInput{
			WorkflowConfig:     in.WorkflowConfigs[idx],
			homeChainSelector:  homeChainOutput.ChainSelector,
			chainSelector:      bo.ChainSelector,
			workflowDonID:      universalSetupOutput.DonTopology.WorkflowDonID,
			feedID:             in.WorkflowConfigs[idx].FeedID,
			addressBook:        universalSetupOutput.CldEnvironment.ExistingAddresses, //nolint:staticcheck // won't migrate now
			priceProvider:      priceProvider,
			sethClient:         bo.SethClient,
			deployerPrivateKey: bo.DeployerPrivateKey,
			creCLIAbsPath:      creCLIAbsPath,
			creCLIsettingsFile: creCLISettingsFile,
			writeTargetName:    corevm.GenerateWriteTargetName(bo.ChainID),
			creCLIProfile:      libcrecli.CRECLIProfile,
		}

		workflowRegisterErr := registerPoRWorkflow(t.Context(), workflowInput)
		require.NoError(t, workflowRegisterErr, "failed to register PoR workflow")

		workflowPauseErr := pausePoRWorkflow(workflowInput)
		require.NoError(t, workflowPauseErr, "failed to pause PoR workflow")

		workflowActivateErr := activatePoRWorkflow(workflowInput)
		require.NoError(t, workflowActivateErr, "failed to activate PoR workflow")
	}
	// Workflow-specific configuration -- END

	// TODO use address book to save the contract addresses

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

// config file to use: environment-one-don-multichain.toml
func TestCRE_OCR3_PoR_Workflow_SingleDon_MultipleWriters_MockedPrice(t *testing.T) {
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

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	capabilityFactoryFns := []keystonetypes.DONCapabilityWithConfigFactoryFn{
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

// config file to use: environment-gateway-don.toml
func TestCRE_OCR3_PoR_Workflow_GatewayDon_MockedPrice(t *testing.T) {
	t.Skip("Disabling until we discover how to fix its flakyness. Tracked as DX-625")
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

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, []string{in.WorkflowConfigs[0].FeedID})
	require.NoError(t, priceErr, "failed to create fake price provider")

	firstBlockchain := in.Blockchains[0]
	chainIDInt, chainErr := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{
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

	firstBlockchain := in.Blockchains[0]
	chainIDInt, chainErr := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	priceProvider := NewTrueUSDPriceProvider(testLogger, []string{in.WorkflowConfigs[0].FeedID})
	setupOutput := setupPoRTestEnvironment(t, testLogger, in, priceProvider, mustSetCapabilitiesFn, []keystonetypes.DONCapabilityWithConfigFactoryFn{
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

func waitForFeedUpdate(t *testing.T, testLogger zerolog.Logger, priceProvider PriceProvider, setupOutput *porSetupOutput, timeout time.Duration) {
	for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
		testLogger.Info().Msgf("Waiting for feed %s to update...", workflowConfig.FeedID)
		timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

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

		require.EqualValues(t, priceProvider.ExpectedPrices(workflowConfig.FeedID), priceProvider.ActualPrices(workflowConfig.FeedID), "prices do not match")
		testLogger.Info().Msgf("All %d prices were found in the feed %s", len(priceProvider.ExpectedPrices(workflowConfig.FeedID)), workflowConfig.FeedID)
	}
}

func debugTest(t *testing.T, testLogger zerolog.Logger, setupOutput *porSetupOutput, in *TestConfig) {
	if t.Failed() {
		counter := 0
		for chainSelector, workflowConfig := range setupOutput.chainSelectorToWorkflowConfig {
			dataFeedsCacheAddresses, dataFeedsCacheErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, df_changeset.DataFeedsCache.String())
			require.NoError(t, dataFeedsCacheErr, "failed to find data feeds cache address for chain %d", chainSelector)

			forwarderAddresses, forwarderErr := crecontracts.FindAddressesForChain(setupOutput.addressBook, chainSelector, keystone_changeset.KeystoneForwarder.String())
			require.NoError(t, forwarderErr, "failed to find forwarder address for chain %d", chainSelector)

			logTestInfo(testLogger, workflowConfig.FeedID, workflowConfig.WorkflowName, dataFeedsCacheAddresses.Hex(), forwarderAddresses.Hex())
			counter++
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
				BlockchainOutput: setupOutput.chainSelectorToBlockchainOutput[chainSelector],
				InfraInput:       in.Infra,
			}
			lidebug.PrintTestDebug(t.Context(), t.Name(), testLogger, debugInput)
		}
	}
}

func waitForWorkflowRegistrySyncer(nodeSetOutput []*keystonetypes.WrappedNodeOutput, topology *keystonetypes.DonTopology) error {
	for idx, nodeSetOut := range nodeSetOutput {
		if !flags.HasFlag(topology.DonsWithMetadata[idx].Flags, keystonetypes.WorkflowDON) {
			continue
		}

		workerNodesOutput := []*clnode.Output{}
		workerNodes, err := crenode.FindManyWithLabel(topology.DonsWithMetadata[idx].NodesMetadata, &types.Label{Key: crenode.NodeTypeKey, Value: types.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return errors.Wrap(err, "failed to find worker nodes")
		}

		for nodeIdx := range workerNodes {
			nodeIndexStr, findErr := crenode.FindLabelValue(workerNodes[nodeIdx], crenode.IndexKey)
			if findErr != nil {
				return errors.Wrapf(findErr, "failed to find node index for node %d in nodeset %s", nodeIdx, topology.DonsWithMetadata[idx].Name)
			}

			nodeIndex, convErr := strconv.Atoi(nodeIndexStr)
			if convErr != nil {
				return errors.Wrapf(convErr, "failed to convert node index '%s' to int for node %d in nodeset %s", nodeIndexStr, nodeIdx, topology.DonsWithMetadata[idx].Name)
			}

			workerNodesOutput = append(workerNodesOutput, nodeSetOut.CLNodes[nodeIndex])
		}

		nsClients, cErr := clclient.New(workerNodesOutput)
		if cErr != nil {
			return errors.Wrap(cErr, "failed to create node set clients")
		}
		eg1 := &errgroup.Group{}
		eg2 := &errgroup.Group{}
		eg3 := &errgroup.Group{}
		for _, c := range nsClients {
			eg1.Go(func() error {
				return c.WaitHealthy(".*WorkflowStore", "passing", 100)
			})
			eg2.Go(func() error {
				return c.WaitHealthy(".*WorkflowRegistrySyncer.FetcherService", "passing", 100)
			})
			eg3.Go(func() error {
				return c.WaitHealthy(".*RegistrySyncer", "passing", 100)
			})
		}
		if err := eg1.Wait(); err != nil {
			return errors.Wrap(err, "failed to wait for WorkflowStore health checks")
		}
		if err := eg2.Wait(); err != nil {
			return errors.Wrap(err, "failed to wait for FetcherService health checks")
		}
		if err := eg3.Wait(); err != nil {
			return errors.Wrap(err, "failed to wait for RegistrySyncer health checks")
		}
	}

	return nil
}
