package capabilities_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"os"
	"runtime"
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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	keystonecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/crib"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	keystoneporconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config/por"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	keystonepor "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/por"
	keystonesecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	libenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	keystoneporcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli/por"
	libfunding "github.com/smartcontractkit/chainlink/system-tests/lib/funding"
	libnix "github.com/smartcontractkit/chainlink/system-tests/lib/nix"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

const (
	cronCapabilityAssetFile            = "amd64_cron"
	ghReadTokenEnvVarName              = "GITHUB_READ_TOKEN"
	E2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	E2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
	cribConfigsDir                     = "crib-configs"
)

type TestConfig struct {
	BlockchainA                   *blockchain.Input                      `toml:"blockchain_a" validate:"required"`
	NodeSets                      []*ns.Input                            `toml:"nodesets" validate:"required"`
	WorkflowConfig                *WorkflowConfig                        `toml:"workflow_config" validate:"required"`
	JD                            *jd.Input                              `toml:"jd" validate:"required"`
	Fake                          *fake.Input                            `toml:"fake"`
	KeystoneContracts             *keystonetypes.KeystoneContractsInput  `toml:"keystone_contracts"`
	WorkflowRegistryConfiguration *keystonetypes.WorkflowRegistryInput   `toml:"workflow_registry_configuration"`
	FeedConsumer                  *keystonetypes.DeployFeedConsumerInput `toml:"feed_consumer"`
	Infra                         *libtypes.InfraInput                   `toml:"infra" validate:"required"`
}

type WorkflowConfig struct {
	UseCRECLI                bool `toml:"use_cre_cli"`
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

// Defines relases/versions of test dependencies that will be downloaded from Github
type DependenciesConfig struct {
	CapabiltiesVersion string `toml:"capabilities_version" validate:"required"`
	CRECLIVersion      string `toml:"cre_cli_version" validate:"required"`
}

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL string `toml:"binary_url" validate:"required"`
	ConfigURL string `toml:"config_url" validate:"required"`
}

func validateEnvVars(t *testing.T, in *TestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")

	var ghReadToken string
	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(libjobs.E2eJobDistributorImageEnvVarName), "missing env var: "+libjobs.E2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(libjobs.E2eJobDistributorVersionEnvVarName), "missing env var: "+libjobs.E2eJobDistributorVersionEnvVarName)

		// disabled until we can figure out how to generate a gist read:write token in CI
		/*
		   This test can be run in two modes:
		   1. `existing` mode: it uses a workflow binary (and configuration) file that is already uploaded to Gist
		   2. `compile` mode: it compiles a new workflow binary and uploads it to Gist

		   For the `new` mode to work, the `GITHUB_API_TOKEN` env var must be set to a token that has `gist:read` and `gist:write` permissions, but this permissions
		   are tied to account not to repository. Currently, we have no service account in the CI at all. And using a token that's tied to personal account of a developer
		   is not a good idea. So, for now, we are only allowing the `existing` mode in CI.
		*/

		// we use this special function to subsitute a placeholder env variable with the actual environment variable name
		// it is defined in .github/e2e-tests.yml as '{{ env.GITHUB_API_TOKEN }}'
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	require.NotEmpty(t, ghReadToken, ghReadTokenEnvVarName+" env var must be set")

	if in.WorkflowConfig.UseCRECLI {
		if in.WorkflowConfig.ShouldCompileNewWorkflow {
			gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
			require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
			err := os.Setenv("CRE_GITHUB_API_TOKEN", gistWriteToken)
			require.NoError(t, err, "failed to set CRE_GITHUB_API_TOKEN env var")
		}
	}
}

type binaryDownloadOutput struct {
	creCLIAbsPath string
}

// this is a small hack to avoid changing the reusable workflow, which doesn't allow to run any pre-execution hooks
func downloadBinaryFiles(in *TestConfig) (*binaryDownloadOutput, error) {
	var ghReadToken string
	if os.Getenv("CI") == "true" {
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	_, err := keystonecapabilities.DownloadCapabilityFromRelease(ghReadToken, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, cronCapabilityAssetFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")
	}

	output := &binaryDownloadOutput{}

	if in.WorkflowConfig.UseCRECLI {
		output.creCLIAbsPath, err = libcrecli.DownloadAndInstallChainlinkCLI(ghReadToken, in.WorkflowConfig.DependenciesConfig.CRECLIVersion)
		if err != nil {
			return nil, errors.Wrap(err, "failed to download and install CRE CLI. Make sure token has content:read permissions to the dev-platform repo")
		}
	}

	return output, nil
}

type registerPoRWorkflowInput struct {
	*WorkflowConfig
	chainSelector               uint64
	workflowDonID               uint32
	feedID                      string
	workflowRegistryAddress     common.Address
	feedConsumerAddress         common.Address
	capabilitiesRegistryAddress common.Address
	priceProvider               PriceProvider
	sethClient                  *seth.Client
	deployerPrivateKey          string
	blockchain                  *blockchain.Output
	binaryDownloadOutput        binaryDownloadOutput
}

func registerPoRWorkflow(input registerPoRWorkflowInput) error {
	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon, but there's still quite a lot of people
	// who have no access to dev-platform repo, so they cannot use the CRE CLI
	if !input.WorkflowConfig.ShouldCompileNewWorkflow && !input.WorkflowConfig.UseCRECLI {
		err := libcontracts.RegisterWorkflow(input.sethClient, input.workflowRegistryAddress, input.workflowDonID, input.WorkflowConfig.WorkflowName, input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL, input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL)
		if err != nil {
			return errors.Wrap(err, "failed to register workflow")
		}

		return nil
	}

	// These two env vars are required by the CRE CLI
	err := os.Setenv("CRE_ETH_PRIVATE_KEY", input.deployerPrivateKey)
	if err != nil {
		return errors.Wrap(err, "failed to set CRE_ETH_PRIVATE_KEY")
	}

	// create CRE CLI settings file
	settingsFile, settingsErr := libcrecli.PrepareCRECLISettingsFile(input.sethClient.MustGetRootKeyAddress(), input.capabilitiesRegistryAddress, input.workflowRegistryAddress, input.workflowDonID, input.chainSelector, input.blockchain.Nodes[0].HostHTTPUrl)
	if settingsErr != nil {
		return errors.Wrap(settingsErr, "failed to create CRE CLI settings file")
	}

	var workflowURL string
	var workflowConfigURL string

	workflowConfigFile, configErr := keystoneporcrecli.CreateConfigFile(input.feedConsumerAddress, input.feedID, input.priceProvider.URL())
	if configErr != nil {
		return errors.Wrap(configErr, "failed to create workflow config file")
	}

	// compile and upload the workflow, if we are not using an existing one
	if input.WorkflowConfig.ShouldCompileNewWorkflow {
		compilationResult, err := libcrecli.CompileWorkflow(input.binaryDownloadOutput.creCLIAbsPath, *input.WorkflowConfig.WorkflowFolderLocation, workflowConfigFile, settingsFile)
		if err != nil {
			return errors.Wrap(err, "failed to compile workflow")
		}

		workflowURL = compilationResult.WorkflowURL
		workflowConfigURL = compilationResult.ConfigURL
	} else {
		workflowURL = input.WorkflowConfig.CompiledWorkflowConfig.BinaryURL
		workflowConfigURL = input.WorkflowConfig.CompiledWorkflowConfig.ConfigURL
	}

	registerErr := libcrecli.DeployWorkflow(input.binaryDownloadOutput.creCLIAbsPath, input.WorkflowName, workflowURL, workflowConfigURL, settingsFile)
	if registerErr != nil {
		return errors.Wrap(registerErr, "failed to register workflow")
	}

	return nil
}

func logTestInfo(l zerolog.Logger, feedID, workflowName, feedConsumerAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedID)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("FeedConsumer address: %s", feedConsumerAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func extraAllowedPortsAndIps(testLogger zerolog.Logger, fakePort int, containerName string) ([]string, []int, error) {
	// we need to explicitly allow the port used by the fake data provider
	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
	// because that's where the fake data provider is running
	var hostIP string
	var err error

	// TODO add handling for CRIB as none of the current cases will work in k8s
	system := runtime.GOOS
	switch system {
	case "darwin":
		hostIP, err = libdon.ResolveHostDockerInternaIP(testLogger, containerName)
	case "linux":
		// for linux framework already returns an IP, so we don't need to resolve it,
		// but we need to remove the http:// prefix
		hostIP = strings.ReplaceAll(framework.HostDockerInternal(), "http://", "")
	default:
		err = fmt.Errorf("unsupported OS: %s", system)
	}
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve host.docker.internal IP")
	}

	testLogger.Info().Msgf("Will allow IP %s and port %d for the fake data provider", hostIP, fakePort)

	ips, err := net.LookupIP("gist.githubusercontent.com")
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to resolve IP for gist.githubusercontent.com")
	}

	gistIPs := make([]string, len(ips))
	for i, ip := range ips {
		gistIPs[i] = ip.To4().String()
		testLogger.Debug().Msgf("Resolved IP for gist.githubusercontent.com: %s", gistIPs[i])
	}

	// we also need to explicitly allow Gist's IP
	return append(gistIPs, hostIP), []int{fakePort}, nil
}

type BlockchainsInput struct {
	blockchainInput *blockchain.Input
	infraInput      *libtypes.InfraInput
	nixShell        *libnix.Shell
}

type BlockchainOutput struct {
	chainSelector      uint64
	blockchainOutput   *blockchain.Output
	sethClient         *seth.Client
	deployerPrivateKey string
}

func CreateBlockchains(
	cldLogger logger.Logger,
	testLogger zerolog.Logger,
	input BlockchainsInput,
) (*BlockchainOutput, error) {
	if input.blockchainInput == nil {
		return nil, errors.New("blockchain input is nil")
	}

	if input.infraInput.InfraType == libtypes.CRIB {
		if input.nixShell == nil {
			return nil, errors.New("nix shell is nil")
		}

		deployCribBlockchainInput := &keystonetypes.DeployCribBlockchainInput{
			BlockchainInput: input.blockchainInput,
			NixShell:        input.nixShell,
			CribConfigsDir:  cribConfigsDir,
		}

		var blockchainErr error
		input.blockchainInput.Out, blockchainErr = crib.DeployBlockchain(deployCribBlockchainInput)
		if blockchainErr != nil {
			return nil, errors.Wrap(blockchainErr, "failed to deploy blockchain")
		}
	}

	// Create a new blockchain network and Seth client to interact with it
	blockchainOutput, err := blockchain.NewBlockchainNetwork(input.blockchainInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blockchain network")
	}

	pkey := os.Getenv("PRIVATE_KEY")
	if pkey == "" {
		return nil, errors.New("PRIVATE_KEY env var must be set")
	}

	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(blockchainOutput.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create seth client")
	}

	chainSelector, err := chainselectors.SelectorFromChainId(sethClient.Cfg.Network.ChainID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chain selector for chain id %d", sethClient.Cfg.Network.ChainID)
	}

	return &BlockchainOutput{
		chainSelector:      chainSelector,
		blockchainOutput:   blockchainOutput,
		sethClient:         sethClient,
		deployerPrivateKey: pkey,
	}, nil
}

func CreateJobDistributor(input *jd.Input) (*jd.Output, error) {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(E2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(E2eJobDistributorVersionEnvVarName)
		input.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}

	jdOutput, err := jd.NewJD(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new job distributor")
	}

	return jdOutput, nil
}

type setupOutput struct {
	priceProvider        PriceProvider
	feedsConsumerAddress common.Address
	forwarderAddress     common.Address
	sethClient           *seth.Client
	blockchainOutput     *blockchain.Output
	donTopology          *keystonetypes.DonTopology
	nodeOutput           []*keystonetypes.WrappedNodeOutput
}

func setupTestEnvironment(t *testing.T, testLogger zerolog.Logger, in *TestConfig, priceProvider PriceProvider, binaryDownloadOutput binaryDownloadOutput, mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet) *setupOutput {
	// Universal setup -- START

	nodeSetInput := mustSetCapabilitiesFn(in.NodeSets)
	topologyErr := libdon.ValidateTopology(nodeSetInput, *in.Infra)
	require.NoError(t, topologyErr, "failed to validate would-be topology")

	// Shell is only required, when using CRIB, because we want to run commands in the same "nix develop" context
	// We need to have this reference in the outer scope, because subsequent functions will need it
	var nixShell *libnix.Shell
	if in.Infra.InfraType == libtypes.CRIB {
		startNixShellInput := &keystonetypes.StartNixShellInput{
			InfraInput:     in.Infra,
			CribConfigsDir: cribConfigsDir,
		}

		var nixErr error
		nixShell, nixErr = crib.StartNixShell(startNixShellInput)
		require.NoError(t, nixErr, "failed to start Nix shell")

		t.Cleanup(func() {
			_ = nixShell.Close()
		})
	}

	blockchainsInput := BlockchainsInput{
		blockchainInput: in.BlockchainA,
		infraInput:      in.Infra,
		nixShell:        nixShell,
	}

	singeFileLogger := cldlogger.NewSingleFileLogger(t)
	blockchainsOutput, err := CreateBlockchains(singeFileLogger, testLogger, blockchainsInput)
	require.NoError(t, err, "failed to start environment")

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability, workflow registry)
	// but first, we need to create deployment.Environment that will contain only chain information in order to deploy contracts with the CLD
	chainsConfig := []devenv.ChainConfig{
		{
			ChainID:   blockchainsOutput.sethClient.Cfg.Network.ChainID,
			ChainName: blockchainsOutput.sethClient.Cfg.Network.Name,
			ChainType: strings.ToUpper(blockchainsOutput.blockchainOutput.Family),
			WSRPCs: []devenv.CribRPCs{{
				External: blockchainsOutput.blockchainOutput.Nodes[0].HostWSUrl,
				Internal: blockchainsOutput.blockchainOutput.Nodes[0].DockerInternalWSUrl,
			}},
			HTTPRPCs: []devenv.CribRPCs{{
				External: blockchainsOutput.blockchainOutput.Nodes[0].HostHTTPUrl,
				Internal: blockchainsOutput.blockchainOutput.Nodes[0].DockerInternalHTTPUrl,
			}},
			DeployerKey: blockchainsOutput.sethClient.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the RPC node
		},
	}

	chains, err := devenv.NewChains(singeFileLogger, chainsConfig)
	require.NoError(t, err, "failed to create chains")

	chainsOnlyCld := &deployment.Environment{
		Logger:            singeFileLogger,
		Chains:            chains,
		ExistingAddresses: deployment.NewMemoryAddressBook(),
		GetContext: func() context.Context {
			return testcontext.Get(t)
		},
	}

	keystoneContractsInput := &keystonetypes.KeystoneContractsInput{
		ChainSelector: blockchainsOutput.chainSelector,
		CldEnv:        chainsOnlyCld,
	}

	keystoneContractsOutput, err := libcontracts.DeployKeystone(testLogger, keystoneContractsInput)
	require.NoError(t, err, "failed to deploy keystone contracts")

	// Translate node input to structure required further down the road and put as much information
	// as we have at this point in labels. It will be used to generate node configs
	topology, err := libdon.BuildTopology(nodeSetInput, *blockchainsInput.infraInput)
	require.NoError(t, err, "failed to build input DON topology")

	// Generate EVM and P2P keys, which are needed to prepare the node configs
	// That way we can pass them final configs and do away with restarting the nodes
	var keys *keystonetypes.GenerateKeysOutput
	chainIDInt, err := strconv.Atoi(blockchainsOutput.blockchainOutput.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")

	generateKeysInput := &keystonetypes.GenerateKeysInput{
		GenerateEVMKeysForChainIDs: []int{chainIDInt},
		GenerateP2PKeys:            true,
		Topology:                   topology,
		Password:                   "", // since the test runs on private ephemeral blockchain we don't use real keys and do not care a lot about the password
	}
	keys, err = libdon.GenereteKeys(generateKeysInput)
	require.NoError(t, err, "failed to generate keys")

	topology, err = libdon.AddKeysToTopology(topology, keys)
	require.NoError(t, err, "failed to add keys to topology")

	// Configure Workflow Registry contract
	workflowRegistryInput := &keystonetypes.WorkflowRegistryInput{
		ChainSelector:  blockchainsOutput.chainSelector,
		CldEnv:         chainsOnlyCld,
		AllowedDonIDs:  []uint32{topology.WorkflowDONID},
		WorkflowOwners: []common.Address{blockchainsOutput.sethClient.MustGetRootKeyAddress()},
	}

	_, err = libcontracts.ConfigureWorkflowRegistry(testLogger, workflowRegistryInput)
	require.NoError(t, err, "failed to configure workflow registry")

	// Allow extra IPs and ports for the fake data provider, which is running on host machine and requires explicit whitelisting
	// If using live endpoint, we don't need to do this
	var extraAllowedIPs []string
	var extraAllowedPorts []int

	if _, ok := priceProvider.(*FakePriceProvider); ok {
		// In the future we might need to have a way to deploy fake price provider to CRIB, now we don't as it is not Dockerised and there are no Helm charts for it
		require.Equal(t, libtypes.Docker, in.Infra.InfraType, "fake data provider is only supported in Docker infra")

		// it doesn't really matter which container we will use to resolve the host.docker.internal IP, it will be the same for all of them
		// here we will blokchain container, because by that time it will be running
		extraAllowedIPs, extraAllowedPorts, err = extraAllowedPortsAndIps(testLogger, in.Fake.Port, blockchainsOutput.blockchainOutput.ContainerName)
		require.NoError(t, err, "failed to get extra allowed ports and IPs")
	}

	peeringData, err := libdon.FindPeeringData(topology)
	require.NoError(t, err, "failed to get peering data")

	for i, donMetadata := range topology.DonsMetadata {
		config, configErr := keystoneporconfig.GenerateConfigs(
			keystonetypes.GeneratePoRConfigsInput{
				DonMetadata:                 donMetadata,
				BlockchainOutput:            blockchainsOutput.blockchainOutput,
				DonID:                       donMetadata.ID,
				Flags:                       donMetadata.Flags,
				PeeringData:                 peeringData,
				CapabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
				WorkflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
				ForwarderAddress:            keystoneContractsOutput.ForwarderAddress,
				GatewayConnectorOutput:      topology.GatewayConnectorOutput,
			},
		)
		require.NoError(t, configErr, "failed to define config for DON %d", donMetadata.ID)

		secretsInput := &keystonetypes.GenerateSecretsInput{
			DonMetadata: donMetadata,
		}

		if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
			secretsInput.EVMKeys = evmKeys
		}

		if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
			secretsInput.P2PKeys = p2pKeys
		}

		// EVM and P2P keys will be provided to nodes as secrets
		secrets, secretsErr := keystonesecrets.GenerateSecrets(
			secretsInput,
		)
		require.NoError(t, secretsErr, "failed to define secrets for DON %d", donMetadata.ID)

		for j := range donMetadata.NodesMetadata {
			nodeSetInput[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			nodeSetInput[i].NodeSpecs[j].Node.TestSecretsOverrides = secrets[j]
		}
	}

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range nodeSetInput {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range nodeSetInput[i].NodeSpecs {
				nodeSetInput[i].NodeSpecs[j].Node.Image = image
			}
		}
	}

	if in.Infra.InfraType == libtypes.CRIB {
		testLogger.Info().Msg("Saving node configs and secret overrides")

		deployCribDonsInput := &keystonetypes.DeployCribDonsInput{
			Topology:       topology,
			NodeSetInputs:  nodeSetInput,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var devspaceErr error
		nodeSetInput, devspaceErr = crib.DeployDons(deployCribDonsInput)
		require.NoError(t, devspaceErr, "failed to deploy Dons with devspace")

		deployCribJdInput := &keystonetypes.DeployCribJdInput{
			JDInput:        in.JD,
			NixShell:       nixShell,
			CribConfigsDir: cribConfigsDir,
		}

		var jdErr error
		in.JD.Out, jdErr = crib.DeployJd(deployCribJdInput)
		require.NoError(t, jdErr, "failed to deploy JD with devspace")
	}

	jdOutput, err := CreateJobDistributor(in.JD)
	require.NoError(t, err, "failed to create new job distributor")

	nodeOutput := make([]*keystonetypes.WrappedNodeOutput, 0, len(nodeSetInput))
	for _, nodeSetInput := range nodeSetInput {
		nodeset, nodesetErr := ns.NewSharedDBNodeSet(nodeSetInput.Input, blockchainsOutput.blockchainOutput)
		require.NoError(t, nodesetErr, "failed to deploy node set named %s", nodeSetInput.Name)

		nodeOutput = append(nodeOutput, &keystonetypes.WrappedNodeOutput{
			Output:       nodeset,
			NodeSetName:  nodeSetInput.Name,
			Capabilities: nodeSetInput.Capabilities,
		})
	}

	// Prepare the CLD environment that's required by the keystone changeset
	// Ugly glue hack ¯\_(ツ)_/¯
	fullCldInput := &keystonetypes.FullCLDEnvironmentInput{
		JdOutput:          jdOutput,
		BlockchainOutput:  blockchainsOutput.blockchainOutput,
		SethClient:        blockchainsOutput.sethClient,
		NodeSetOutput:     nodeOutput,
		ExistingAddresses: chainsOnlyCld.ExistingAddresses,
		Topology:          topology,
	}

	// We need to use TLS for CRIB, because it exposes HTTPS endpoints
	var creds credentials.TransportCredentials
	if in.Infra.InfraType == libtypes.CRIB {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	fullCldOutput, err := libenv.BuildFullCLDEnvironment(singeFileLogger, fullCldInput, creds)
	require.NoError(t, err, "failed to build chainlink deployment environment")

	// Fund the nodes
	for _, metaDon := range fullCldOutput.DonTopology.DonsWithMetadata {
		for _, node := range metaDon.DON.Nodes {
			_, fundingErr := libfunding.SendFunds(zerolog.Logger{}, blockchainsOutput.sethClient, libtypes.FundsToSend{
				ToAddress:  common.HexToAddress(node.AccountAddr[blockchainsOutput.sethClient.Cfg.Network.ChainID]),
				Amount:     big.NewInt(5000000000000000000),
				PrivateKey: blockchainsOutput.sethClient.MustGetRootPrivateKey(),
			})
			require.NoError(t, fundingErr, "failed to send funds to node %s", node.AccountAddr[blockchainsOutput.sethClient.Cfg.Network.ChainID])
		}
	}

	// Generate and propose jobs (they will auto-accepted)
	donToJobSpecs, jobSpecsErr := keystonepor.GenerateJobSpecs(
		&keystonetypes.GeneratePoRJobSpecsInput{
			BlockchainOutput:       blockchainsOutput.blockchainOutput,
			DonsWithMetadata:       fullCldOutput.DonTopology.DonsWithMetadata,
			OCR3CapabilityAddress:  keystoneContractsOutput.OCR3CapabilityAddress,
			ExtraAllowedPorts:      extraAllowedPorts,
			ExtraAllowedIPs:        extraAllowedIPs,
			CronCapBinName:         cronCapabilityAssetFile,
			GatewayConnectorOutput: *topology.GatewayConnectorOutput,
		},
	)
	require.NoError(t, jobSpecsErr, "failed to define job specs for DONs")

	createJobsInput := keystonetypes.CreateJobsInput{
		CldEnv:        fullCldOutput.Environment,
		DonTopology:   fullCldOutput.DonTopology,
		DonToJobSpecs: donToJobSpecs,
	}

	// TODO in the future, maybe should we remove all jobs first, if it's running in CRIB? or at least jobs of certain types?
	// that would allow us to run the same test multiple times without the need to restart the whole environment
	err = libdon.CreateJobs(testLogger, createJobsInput)
	require.NoError(t, err, "failed to configure nodes and create jobs")

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO: workflow/core team should expose a way for us to check if the OCR listener is ready
	testLogger.Info().Msg("Waiting 45s for OCR listeners to be ready...")
	time.Sleep(45 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 and Keystone configuration...")

	// Configure the Forwarder, OCR3 and Capabilities contracts
	configureKeystoneInput := keystonetypes.ConfigureKeystoneInput{
		ChainSelector: blockchainsOutput.chainSelector,
		CldEnv:        fullCldOutput.Environment,
		Topology:      topology,
	}
	err = libcontracts.ConfigureKeystone(configureKeystoneInput)
	require.NoError(t, err, "failed to configure keystone contracts")

	// Universal setup -- END
	// Workflow-specific configuration -- START
	deployFeedConsumerInput := &keystonetypes.DeployFeedConsumerInput{
		ChainSelector: blockchainsOutput.chainSelector,
		CldEnv:        chainsOnlyCld,
	}
	deployFeedsConsumerOutput, err := libcontracts.DeployFeedsConsumer(testLogger, deployFeedConsumerInput)
	require.NoError(t, err, "failed to deploy feeds consumer")

	configureFeedConsumerInput := &keystonetypes.ConfigureFeedConsumerInput{
		SethClient:            blockchainsOutput.sethClient,
		FeedConsumerAddress:   deployFeedsConsumerOutput.FeedConsumerAddress,
		AllowedSenders:        []common.Address{keystoneContractsOutput.ForwarderAddress},
		AllowedWorkflowOwners: []common.Address{blockchainsOutput.sethClient.MustGetRootKeyAddress()},
		AllowedWorkflowNames:  []string{in.WorkflowConfig.WorkflowName},
	}
	_, err = libcontracts.ConfigureFeedsConsumer(testLogger, configureFeedConsumerInput)
	require.NoError(t, err, "failed to configure feeds consumer")

	registerInput := registerPoRWorkflowInput{
		WorkflowConfig:              in.WorkflowConfig,
		chainSelector:               blockchainsOutput.chainSelector,
		workflowDonID:               fullCldOutput.DonTopology.WorkflowDonID,
		feedID:                      in.WorkflowConfig.FeedID,
		workflowRegistryAddress:     keystoneContractsOutput.WorkflowRegistryAddress,
		feedConsumerAddress:         deployFeedsConsumerOutput.FeedConsumerAddress,
		capabilitiesRegistryAddress: keystoneContractsOutput.CapabilitiesRegistryAddress,
		priceProvider:               priceProvider,
		sethClient:                  blockchainsOutput.sethClient,
		deployerPrivateKey:          blockchainsOutput.deployerPrivateKey,
		blockchain:                  blockchainsOutput.blockchainOutput,
		binaryDownloadOutput:        binaryDownloadOutput,
	}

	err = registerPoRWorkflow(registerInput)
	require.NoError(t, err, "failed to register PoR workflow")
	// Workflow-specific configuration -- END

	// Set inputs in the test config, so that they can be saved
	in.KeystoneContracts = keystoneContractsInput
	in.FeedConsumer = deployFeedConsumerInput
	in.WorkflowRegistryConfiguration = workflowRegistryInput

	return &setupOutput{
		priceProvider:        priceProvider,
		feedsConsumerAddress: deployFeedsConsumerOutput.FeedConsumerAddress,
		forwarderAddress:     keystoneContractsOutput.ForwarderAddress,
		sethClient:           blockchainsOutput.sethClient,
		blockchainOutput:     blockchainsOutput.blockchainOutput,
		donTopology:          fullCldOutput.DonTopology,
		nodeOutput:           nodeOutput,
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

	binaryDownloadOutput, err := downloadBinaryFiles(in)
	require.NoError(t, err, "failed to download binary files")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       keystonetypes.SingleDonFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake)
	require.NoError(t, priceErr, "failed to create fake price provider")

	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, *binaryDownloadOutput, mustSetCapabilitiesFn)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

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
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, _, err := feedsConsumerInstance.GetPrice(
			setupOutput.sethClient.NewCallOpts(),
			feedBytes,
		)
		require.NoError(t, err, "failed to get price from Keystone Consumer contract")

		hasNextPrice := setupOutput.priceProvider.NextPrice(price, elapsed)
		if !hasNextPrice {
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}

		return !hasNextPrice
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

	binaryDownloadOutput, err := downloadBinaryFiles(in)
	require.NoError(t, err, "failed to download binary files")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       keystonetypes.SingleDonFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				GatewayNodeIndex:   0,
				BootstrapNodeIndex: -1, // <----- it's crucial to indicate there's no bootstrap node
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake)
	require.NoError(t, priceErr, "failed to create fake price provider")

	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, *binaryDownloadOutput, mustSetCapabilitiesFn)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

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
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, _, err := feedsConsumerInstance.GetPrice(
			setupOutput.sethClient.NewCallOpts(),
			feedBytes,
		)
		require.NoError(t, err, "failed to get price from Keystone Consumer contract")

		hasNextPrice := setupOutput.priceProvider.NextPrice(price, elapsed)
		if !hasNextPrice {
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}

		return !hasNextPrice
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

	binaryDownloadOutput, err := downloadBinaryFiles(in)
	require.NoError(t, err, "failed to download binary files")

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
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[2],
				Capabilities:       []string{},
				DONTypes:           []string{keystonetypes.GatewayDON}, // <----- it's crucial to set the correct DON type
				GatewayNodeIndex:   0,
				BootstrapNodeIndex: -1, // <----- it's crucial to indicate there's no bootstrap node
			},
		}
	}

	priceProvider := NewTrueUSDPriceProvider(testLogger)
	setupOutput := setupTestEnvironment(t, testLogger, in, priceProvider, *binaryDownloadOutput, mustSetCapabilitiesFn)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

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
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	testLogger.Info().Msg("Waiting for feed to update...")
	timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(setupOutput.feedsConsumerAddress, setupOutput.sethClient.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	feedBytes := common.HexToHash(in.WorkflowConfig.FeedID)

	assert.Eventually(t, func() bool {
		elapsed := time.Since(startTime).Round(time.Second)
		price, _, err := feedsConsumerInstance.GetPrice(
			setupOutput.sethClient.NewCallOpts(),
			feedBytes,
		)
		require.NoError(t, err, "failed to get price from Keystone Consumer contract")

		hasNextPrice := setupOutput.priceProvider.NextPrice(price, elapsed)
		if !hasNextPrice {
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}

		return !hasNextPrice
	}, timeout, 10*time.Second, "feed did not update, timeout after: %s", timeout)

	require.EqualValues(t, priceProvider.ExpectedPrices(), priceProvider.ActualPrices(), "prices do not match")
	testLogger.Info().Msgf("All %d prices were found in the feed", len(priceProvider.ExpectedPrices()))
}
