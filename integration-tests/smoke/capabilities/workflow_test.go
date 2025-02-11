package capabilities_test

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-yaml/yaml"
	"github.com/google/go-github/v41/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials/insecure"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

type WorkflowConfig struct {
	UseCRECLI                bool `toml:"use_cre_cli"`
	ShouldCompileNewWorkflow bool `toml:"should_compile_new_workflow"`
	// Tells the test where the workflow to compile is located
	WorkflowFolderLocation *string             `toml:"workflow_folder_location"`
	CompiledWorkflowConfig *CompiledConfig     `toml:"compiled_config"`
	DependenciesConfig     *DependenciesConfig `toml:"dependencies"`
	// id, which will be used, when registering the DON with the workflow registry,
	// and when instructing the Gateway job on the bootstrap node as to which workflow to run.
	DonID        uint32 `toml:"don_id" validate:"required"`
	WorkflowName string `toml:"workflow_name" validate:"required" `
}

// Defines relases/versions of test dependencies that will be downloaded from Github
type DependenciesConfig struct {
	CapabiltiesVersion string `toml:"capabilities_version"`
	CRECLIVersion      string `toml:"cre_cli_version"`
}

// Defines the location of already compiled workflow binary and config files
// They will be used if WorkflowConfig.ShouldCompileNewWorkflow is `false`
// Otherwise test will compile and upload a new workflow
type CompiledConfig struct {
	BinaryURL string `toml:"binary_url"`
	ConfigURL string `toml:"config_url"`
}

type WorkflowTestConfig struct {
	BlockchainA    *blockchain.Input    `toml:"blockchain_a" validate:"required"`
	NodeSet        *ns.Input            `toml:"nodeset" validate:"required"`
	WorkflowConfig *WorkflowConfig      `toml:"workflow_config" validate:"required"`
	JD             *jd.Input            `toml:"jd" validate:"required"`
	PriceProvider  *PriceProviderConfig `toml:"price_provider"`
}

type FakeConfig struct {
	*fake.Input
	Prices []float64 `toml:"prices"`
}

type PriceProviderConfig struct {
	Fake   *FakeConfig `toml:"fake"`
	FeedID string      `toml:"feed_id" validate:"required"`
	URL    string      `toml:"url"`
}

func downloadGHAssetFromRelease(owner, repository, releaseTag, assetName, ghToken string) ([]byte, error) {
	var content []byte
	if ghToken == "" {
		return content, errors.New("no github token provided")
	}

	// assuming 180s is enough to fetch releases, find the asset we need and download it
	// some assets might be 30+ MB, so we need to give it some time (for really slow connections)
	ctx, cancelFn := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancelFn()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)

	ghReleases, _, err := ghClient.Repositories.ListReleases(ctx, owner, repository, &github.ListOptions{PerPage: 20})
	if err != nil {
		return content, errors.Wrapf(err, "failed to list releases for %s", repository)
	}

	var ghRelease *github.RepositoryRelease
	for _, release := range ghReleases {
		if release.TagName == nil {
			continue
		}

		if *release.TagName == releaseTag {
			ghRelease = release
			break
		}
	}

	if ghRelease == nil {
		return content, errors.New("failed to find release with tag: " + releaseTag)
	}

	var assetID int64
	for _, asset := range ghRelease.Assets {
		if strings.Contains(asset.GetName(), assetName) {
			assetID = asset.GetID()
			break
		}
	}

	if assetID == 0 {
		return content, fmt.Errorf("failed to find asset %s for %s", assetName, *ghRelease.TagName)
	}

	asset, _, err := ghClient.Repositories.DownloadReleaseAsset(ctx, owner, repository, assetID, tc)
	if err != nil {
		return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *ghRelease.TagName)
	}

	content, err = io.ReadAll(asset)
	if err != nil {
		return content, err
	}

	return content, nil
}

func GenerateWorkflowIDFromStrings(owner string, name string, workflow []byte, config []byte, secretsURL string) (string, error) {
	ownerWithoutPrefix := owner
	if strings.HasPrefix(owner, "0x") {
		ownerWithoutPrefix = owner[2:]
	}

	ownerb, err := hex.DecodeString(ownerWithoutPrefix)
	if err != nil {
		return "", err
	}

	wid, err := pkgworkflows.GenerateWorkflowID(ownerb, name, workflow, config, secretsURL)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(wid[:]), nil
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func download(url string) ([]byte, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func downloadAndDecode(url string) ([]byte, error) {
	data, err := download(url)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return decoded, nil
}

type CRECLISettings struct {
	DevPlatform  DevPlatform  `yaml:"dev-platform"`
	UserWorkflow UserWorkflow `yaml:"user-workflow"`
	Logging      Logging      `yaml:"logging"`
	McmsConfig   McmsConfig   `yaml:"mcms-config"`
	Contracts    Contracts    `yaml:"contracts"`
	Rpcs         []RPC        `yaml:"rpcs"`
}

type DevPlatform struct {
	CapabilitiesRegistryAddress string `yaml:"capabilities-registry-contract-address"`
	DonID                       uint32 `yaml:"don-id"`
	WorkflowRegistryAddress     string `yaml:"workflow-registry-contract-address"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries"`
}

type ContractRegistry struct {
	Name          string `yaml:"name"`
	Address       string `yaml:"address"`
	ChainSelector uint64 `yaml:"chain-selector"`
}

type RPC struct {
	ChainSelector uint64 `yaml:"chain-selector"`
	URL           string `yaml:"url"`
}

type PoRWorkflowConfig struct {
	FeedID          string `json:"feed_id"`
	URL             string `json:"url"`
	ConsumerAddress string `json:"consumer_address"`
}

const (
	CRECLISettingsFileName             = ".cre-cli-settings.yaml"
	cronCapabilityAssetFile            = "amd64_cron"
	e2eJobDistributorImageEnvVarName   = "E2E_JD_IMAGE"
	e2eJobDistributorVersionEnvVarName = "E2E_JD_VERSION"
	ghReadTokenEnvVarName              = "GITHUB_READ_TOKEN"
	GistIP                             = "185.199.108.133"
)

var (
	CRECLICommand string
)

func downloadAndInstallChainlinkCLI(ghToken, version string) error {
	system := runtime.GOOS
	arch := runtime.GOARCH

	switch system {
	case "darwin", "linux":
		// nothing to do, we have the binaries
	default:
		return fmt.Errorf("chainlnk-cli does not support OS: %s", system)
	}

	switch arch {
	case "amd64", "arm64":
		// nothing to do, we have the binaries
	default:
		return fmt.Errorf("chainlnk-cli does not support arch: %s", arch)
	}

	CRECLIAssetFile := fmt.Sprintf("cre_%s_%s_%s.tar.gz", version, system, arch)
	content, err := downloadGHAssetFromRelease("smartcontractkit", "dev-platform", version, CRECLIAssetFile, ghToken)
	if err != nil {
		return errors.Wrapf(err, "failed to download CRE CLI asset %s", CRECLIAssetFile)
	}

	tmpfile, err := os.CreateTemp("", CRECLIAssetFile)
	if err != nil {
		return errors.Wrapf(err, "failed to create temp file for CRE CLI asset %s", CRECLIAssetFile)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(content); err != nil {
		return errors.Wrapf(err, "failed to write content to temp file for CRE CLI asset %s", CRECLIAssetFile)
	}

	cmd := exec.Command("tar", "-xvf", tmpfile.Name(), "-C", ".") // #nosec G204
	if cmd.Run() != nil {
		return errors.Wrapf(err, "failed to extract CRE CLI asset %s", CRECLIAssetFile)
	}

	extractedFileName := fmt.Sprintf("cre_%s_%s_%s", version, system, arch)
	cmd = exec.Command("chmod", "+x", extractedFileName)
	if cmd.Run() != nil {
		return errors.Wrapf(err, "failed to make %s executable", extractedFileName)
	}

	// set it to absolute path, because some commands (e.g. compile) need to be executed in the context
	// of the workflow directory
	extractedFile, err := os.Open(extractedFileName)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", extractedFileName)
	}

	CRECLICommand, err = filepath.Abs(extractedFile.Name())
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for %s", tmpfile.Name())
	}

	if isInstalled := isInstalled(CRECLICommand); !isInstalled {
		return errors.New("failed to install CRE CLI or it is not available in the PATH")
	}

	return nil
}

func downloadCronCapability(ghToken, version string) (string, error) {
	content, err := downloadGHAssetFromRelease("smartcontractkit", "capabilities", version, cronCapabilityAssetFile, ghToken)
	if err != nil {
		return "", err
	}

	fileName := cronCapabilityAssetFile
	file, err := os.Create(cronCapabilityAssetFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return "", err
	}

	return fileName, nil
}

func validateInputsAndEnvVars(t *testing.T, in *WorkflowTestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")
	require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig, "dependencies config must be set")

	if !in.WorkflowConfig.UseCRECLI {
		require.False(t, in.WorkflowConfig.ShouldCompileNewWorkflow, "if you are not using CRE CLI you cannot compile a new workflow")
	}

	var ghReadToken string
	// this is a small hack to avoid changing the reusable workflow
	if os.Getenv("CI") == "true" {
		// This part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
		// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)
		require.NotEmpty(t, os.Getenv(e2eJobDistributorImageEnvVarName), "missing env var: "+e2eJobDistributorImageEnvVarName)
		require.NotEmpty(t, os.Getenv(e2eJobDistributorVersionEnvVarName), "missing env var: "+e2eJobDistributorVersionEnvVarName)

		// disabled until we can figure out how to generate a gist read:write token in CI
		/*
		 This test can be run in two modes:
		 1. `existing` mode: it uses a workflow binary (and configuration) file that is already uploaded to Gist
		 2. `compile` mode: it compiles a new workflow binary and uploads it to Gist

		 For the `new` mode to work, the `GITHUB_API_TOKEN` env var must be set to a token that has `gist:read` and `gist:write` permissions, but this permissions
		 are tied to account not to repository. Currently, we have no service account in the CI at all. And using a token that's tied to personal account of a developer
		 is not a good idea. So, for now, we are only allowing the `existing` mode in CI.
		*/
		require.False(t, in.WorkflowConfig.ShouldCompileNewWorkflow, "you cannot compile a new workflow in the CI as of now due to issues with generating a gist write token")

		// we use this special function to subsitute a placeholder env variable with the actual environment variable name
		// it is defined in .github/e2e-tests.yml as '{{ env.GITHUB_API_TOKEN }}'
		ghReadToken = ctfconfig.MustReadEnvVar_String(ghReadTokenEnvVarName)
	} else {
		ghReadToken = os.Getenv(ghReadTokenEnvVarName)
	}

	require.NotEmpty(t, ghReadToken, ghReadTokenEnvVarName+" env var must be set")
	require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion, "capabilities_version must be set in the dependencies config")

	_, err := downloadCronCapability(ghReadToken, in.WorkflowConfig.DependenciesConfig.CapabiltiesVersion)
	require.NoError(t, err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")

	if in.WorkflowConfig.UseCRECLI {
		require.NotEmpty(t, in.WorkflowConfig.DependenciesConfig.CRECLIVersion, "chainlink_cli_version must be set in the dependencies config")

		err = downloadAndInstallChainlinkCLI(ghReadToken, in.WorkflowConfig.DependenciesConfig.CRECLIVersion)
		require.NoError(t, err, "failed to download and install CRE CLI. Make sure token has content:read permissions to the dev-platform repo")

		if in.WorkflowConfig.ShouldCompileNewWorkflow {
			gistWriteToken := os.Getenv("GIST_WRITE_TOKEN")
			require.NotEmpty(t, gistWriteToken, "GIST_WRITE_TOKEN must be set to use CRE CLI to compile workflows. It requires gist:read and gist:write permissions")
			err := os.Setenv("GITHUB_API_TOKEN", gistWriteToken)
			require.NoError(t, err, "failed to set GITHUB_API_TOKEN env var")
			require.NotEmpty(t, in.WorkflowConfig.WorkflowFolderLocation, "workflow_folder_location must be set, when compiling new workflow")
		}
	}

	if in.PriceProvider.Fake == nil {
		require.NotEmpty(t, in.PriceProvider.URL, "URL must be set in the price provider config, if fake provider is not used")
	}

	// make sure the feed id is in the correct format
	in.PriceProvider.FeedID = strings.TrimPrefix(in.PriceProvider.FeedID, "0x")
}

// copied from Bala's unmerged PR: https://github.com/smartcontractkit/chainlink/pull/15751
// TODO: remove this once the PR is merged and import his function
func getNodeInfo(nodeOut *ns.Output, bootstrapNodeCount int) ([]devenv.NodeInfo, error) {
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= len(nodeOut.CLNodes); i++ {
		p2pURL, err := url.Parse(nodeOut.CLNodes[i-1].Node.DockerP2PUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse p2p url: %w", err)
		}
		if i <= bootstrapNodeCount {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("bootstrap-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("node-%d", i),
				P2PPort:     p2pURL.Port(),
				CLConfig: nodeclient.ChainlinkConfig{
					URL:        nodeOut.CLNodes[i-1].Node.HostURL,
					Email:      nodeOut.CLNodes[i-1].Node.APIAuthUser,
					Password:   nodeOut.CLNodes[i-1].Node.APIAuthPassword,
					InternalIP: nodeOut.CLNodes[i-1].Node.InternalIP,
				},
			})
		}
	}
	return nodeInfo, nil
}

func buildChainlinkDeploymentEnv(t *testing.T, jdOutput *jd.Output, nodeOutput *ns.Output, bs *blockchain.Output, sc *seth.Client) (*deployment.Environment, *devenv.DON, uint64) {
	lgr := logger.TestLogger(t)

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	require.NoError(t, err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)

	nodeInfo, err := getNodeInfo(nodeOutput, 1)
	require.NoError(t, err, "failed to get node info")

	jdConfig := devenv.JDConfig{
		GRPC:     jdOutput.HostGRPCUrl,
		WSRPC:    jdOutput.DockerWSRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: nodeInfo,
	}

	require.GreaterOrEqual(t, len(bs.Nodes), 1, "expected at least one node in the blockchain output")

	devenvConfig := devenv.EnvironmentConfig{
		JDConfig: jdConfig,
		Chains: []devenv.ChainConfig{
			{
				ChainID:   sc.Cfg.Network.ChainID,
				ChainName: sc.Cfg.Network.Name,
				ChainType: strings.ToUpper(bs.Family),
				WSRPCs: []devenv.CribRPCs{{
					External: bs.Nodes[0].HostWSUrl,
					Internal: bs.Nodes[0].DockerInternalWSUrl,
				}},
				HTTPRPCs: []devenv.CribRPCs{{
					External: bs.Nodes[0].HostHTTPUrl,
					Internal: bs.Nodes[0].DockerInternalHTTPUrl,
				}},
				DeployerKey: sc.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
			},
		},
	}

	env, don, err := devenv.NewEnvironment(context.Background, lgr, devenvConfig)
	require.NoError(t, err, "failed to create environment")

	return env, don, chainSelector
}

func deployKeystoneContracts(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) keystone_changeset.ContractSet {
	// Deploy keystone forwarder contract
	_ = deployKeystoneForwarder(t, testLogger, ctfEnv, chainSelector)

	// Deploy OCR3 contract
	_ = deployOCR3(t, testLogger, ctfEnv, chainSelector)

	// Deploy capabilities registry contract
	_ = deployCapabilitiesRegistry(t, testLogger, ctfEnv, chainSelector)

	contractSetResponse, err := keystone_changeset.GetContractSets(nil, &keystone_changeset.GetContractSetsRequest{
		Chains:      ctfEnv.Chains,
		AddressBook: ctfEnv.ExistingAddresses,
	})
	require.NoError(t, err, "failed to get contract sets")

	contractSet, ok := contractSetResponse.ContractSets[chainSelector]
	require.True(t, ok, "failed to get contract set for chain %d", chainSelector)

	return contractSet
}

func deployOCR3(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployOCR3(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy OCR3 Capability contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "OCR3Capability") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed OCR3Capability contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func deployCapabilitiesRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployCapabilityRegistry(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy Capabilities Registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "CapabilitiesRegistry") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed Capabilities Registry contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func deployKeystoneForwarder(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64) common.Address {
	output, err := keystone_changeset.DeployForwarder(*ctfEnv, keystone_changeset.DeployForwarderRequest{
		ChainSelectors: []uint64{chainSelector},
	})
	require.NoError(t, err, "failed to deploy forwarder contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var forwarderAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "KeystoneForwarder") {
			forwarderAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed KeystoneForwarder contract at %s", forwarderAddress.Hex())
			break
		}
	}

	return forwarderAddress
}

func prepareWorkflowRegistry(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64, sc *seth.Client, donID uint32) common.Address {
	output, err := workflow_registry_changeset.Deploy(*ctfEnv, chainSelector)
	require.NoError(t, err, "failed to deploy workflow registry contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var workflowRegistryAddr common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "WorkflowRegistry") {
			workflowRegistryAddr = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed WorkflowRegistry contract at %s", workflowRegistryAddr.Hex())
		}
	}

	// Configure Workflow Registry contract
	_, err = workflow_registry_changeset.UpdateAllowedDons(*ctfEnv, &workflow_registry_changeset.UpdateAllowedDonsRequest{
		RegistryChainSel: chainSelector,
		DonIDs:           []uint32{donID},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update allowed Dons")

	_, err = workflow_registry_changeset.UpdateAuthorizedAddresses(*ctfEnv, &workflow_registry_changeset.UpdateAuthorizedAddressesRequest{
		RegistryChainSel: chainSelector,
		Addresses:        []string{sc.MustGetRootKeyAddress().Hex()},
		Allowed:          true,
	})
	require.NoError(t, err, "failed to update authorized addresses")

	return workflowRegistryAddr
}

func prepareFeedsConsumer(t *testing.T, testLogger zerolog.Logger, ctfEnv *deployment.Environment, chainSelector uint64, sc *seth.Client, forwarderAddress common.Address, workflowName string) common.Address {
	output, err := keystone_changeset.DeployFeedsConsumer(*ctfEnv, &keystone_changeset.DeployFeedsConsumerRequest{
		ChainSelector: chainSelector,
	})
	require.NoError(t, err, "failed to deploy feeds_consumer contract")

	err = ctfEnv.ExistingAddresses.Merge(output.AddressBook)
	require.NoError(t, err, "failed to merge address book")

	addresses, err := ctfEnv.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err, "failed to get addresses for chain %d from the address book", chainSelector)

	var feedsConsumerAddress common.Address
	for addrStr, tv := range addresses {
		if strings.Contains(tv.String(), "FeedConsumer") {
			feedsConsumerAddress = common.HexToAddress(addrStr)
			testLogger.Info().Msgf("Deployed FeedConsumer contract at %s", feedsConsumerAddress.Hex())
			break
		}
	}

	require.NotEmpty(t, feedsConsumerAddress, "failed to find FeedConsumer address in the address book")

	// configure Keystone Feeds Consumer contract, so it can accept reports from the forwarder contract,
	// that come from our workflow that is owned by the root private key
	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	// Prepare hex-encoded and truncated workflow name
	var workflowNameBytes [10]byte
	var HashTruncateName = func(name string) string {
		// Compute SHA-256 hash of the input string
		hash := sha256.Sum256([]byte(name))

		// Encode as hex to ensure UTF8
		var hashBytes []byte = hash[:]
		resultHex := hex.EncodeToString(hashBytes)

		// Truncate to 10 bytes
		truncated := []byte(resultHex)[:10]
		return string(truncated)
	}

	truncated := HashTruncateName(workflowName)
	copy(workflowNameBytes[:], []byte(truncated))

	_, decodeErr := sc.Decode(feedsConsumerInstance.SetConfig(
		sc.NewTXOpts(),
		[]common.Address{forwarderAddress},           // allowed senders
		[]common.Address{sc.MustGetRootKeyAddress()}, // allowed workflow owners
		// here we need to use hex-encoded workflow name converted to []byte
		[][10]byte{workflowNameBytes}, // allowed workflow names
	))
	require.NoError(t, decodeErr, "failed to set config for feeds consumer")

	return feedsConsumerAddress
}

func registerWorkflowDirectly(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName string) {
	require.NotEmpty(t, in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL)
	workFlowData, err := downloadAndDecode(in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL)
	require.NoError(t, err, "failed to download and decode workflow binary")

	var configData []byte
	if in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL != "" {
		configData, err = download(in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL)
		require.NoError(t, err, "failed to download workflow config")
	}

	// use non-encoded workflow name
	workflowID, idErr := GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
	require.NoError(t, idErr, "failed to generate workflow ID")

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	require.NoError(t, err, "failed to create workflow registry instance")

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL, in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL, ""))
	require.NoError(t, decodeErr, "failed to register workflow")
}

//revive:disable // ignore confusing-results
func compileWorkflowWithCRECLI(t *testing.T, in *WorkflowTestConfig, feedsConsumerAddress common.Address, feedID, dataURL string, settingsFile *os.File) (string, string) {
	configFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err, "failed to create workflow config file")

	cleanFeedId := strings.TrimPrefix(feedID, "0x")
	feedLength := len(cleanFeedId)

	require.GreaterOrEqual(t, feedLength, 32, "feed ID must be at least 32 characters long")

	if feedLength > 32 {
		cleanFeedId = cleanFeedId[:32]
	}

	feedIDToUse := "0x" + cleanFeedId

	workflowConfig := PoRWorkflowConfig{
		FeedID:          feedIDToUse,
		URL:             dataURL,
		ConsumerAddress: feedsConsumerAddress.Hex(),
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	require.NoError(t, err, "failed to marshal workflow config")

	_, err = configFile.Write(configMarshalled)
	require.NoError(t, err, "failed to write workflow config file")

	var outputBuffer bytes.Buffer

	// the CLI expects the workflow code to be located in the same directory as its `go.mod`` file. That's why we assume that the file, which
	// contains the entrypoint method is always named `main.go`. This is a limitation of the CLI, which we can't change.
	compileCmd := exec.Command(CRECLICommand, "workflow", "compile", "-S", settingsFile.Name(), "-c", configFile.Name(), "main.go") // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = *in.WorkflowConfig.WorkflowFolderLocation
	err = compileCmd.Start()
	require.NoError(t, err, "failed to start compile command")

	err = compileCmd.Wait()
	fmt.Println("Compile output:\n", outputBuffer.String())

	require.NoError(t, err, "failed to wait for compile command")

	re := regexp.MustCompile(`Gist URL=([^\s]+)`)
	matches := re.FindAllStringSubmatch(outputBuffer.String(), -1)
	require.Len(t, matches, 2, "failed to find 2 gist URLs in compile output")

	ansiEscapePattern := `\x1b\[[0-9;]*m`
	re = regexp.MustCompile(ansiEscapePattern)

	workflowGistURL := re.ReplaceAllString(matches[0][1], "")
	workflowConfigURL := re.ReplaceAllString(matches[1][1], "")

	require.NotEmpty(t, workflowGistURL, "failed to find workflow gist URL")
	require.NotEmpty(t, workflowConfigURL, "failed to find workflow config gist URL")

	return workflowGistURL, workflowConfigURL
}

func preapreCRECLISettingsFile(t *testing.T, sc *seth.Client, capRegAddr, workflowRegistryAddr common.Address, donID uint32, chainSelector uint64, rpcHTTPURL string) *os.File {
	// create CRE CLI settings file
	settingsFile, err := os.CreateTemp("", CRECLISettingsFileName)
	require.NoError(t, err, "failed to create CRE CLI settings file")

	settings := CRECLISettings{
		DevPlatform: DevPlatform{
			CapabilitiesRegistryAddress: capRegAddr.Hex(),
			DonID:                       donID,
			WorkflowRegistryAddress:     workflowRegistryAddr.Hex(),
		},
		UserWorkflow: UserWorkflow{
			WorkflowOwnerAddress: sc.MustGetRootKeyAddress().Hex(),
		},
		Logging: Logging{},
		McmsConfig: McmsConfig{
			ProposalsDirectory: "./",
		},
		Contracts: Contracts{
			ContractRegistry: []ContractRegistry{
				{
					Name:          "CapabilitiesRegistry",
					Address:       capRegAddr.Hex(),
					ChainSelector: chainSelector,
				},
				{
					Name:          "WorkflowRegistry",
					Address:       workflowRegistryAddr.Hex(),
					ChainSelector: chainSelector,
				},
			},
		},
		Rpcs: []RPC{
			{
				ChainSelector: chainSelector,
				URL:           rpcHTTPURL,
			},
		},
	}

	settingsMarshalled, err := yaml.Marshal(settings)
	require.NoError(t, err, "failed to marshal CRE CLI settings")

	_, err = settingsFile.Write(settingsMarshalled)
	require.NoError(t, err, "failed to write %s settings file", CRECLISettingsFileName)

	return settingsFile
}

func registerWorkflow(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, capRegAddr, workflowRegistryAddr, feedsConsumerAddress common.Address, donID uint32, chainSelector uint64, workflowName, pkey, rpcHTTPURL, dataURL string) {
	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon
	if !in.WorkflowConfig.ShouldCompileNewWorkflow && !in.WorkflowConfig.UseCRECLI {
		registerWorkflowDirectly(t, in, sc, workflowRegistryAddr, donID, workflowName)

		return
	}

	// These two env vars are required by the CRE CLI
	err := os.Setenv("WORKFLOW_OWNER_ADDRESS", sc.MustGetRootKeyAddress().Hex())
	require.NoError(t, err, "failed to set WORKFLOW_OWNER_ADDRESS env var")

	err = os.Setenv("ETH_PRIVATE_KEY", pkey)
	require.NoError(t, err, "failed to set ETH_PRIVATE_KEY env var")

	// create CRE CLI settings file
	settingsFile := preapreCRECLISettingsFile(t, sc, capRegAddr, workflowRegistryAddr, donID, chainSelector, rpcHTTPURL)

	var workflowGistURL string
	var workflowConfigURL string

	// compile and upload the workflow, if we are not using an existing one
	if in.WorkflowConfig.ShouldCompileNewWorkflow {
		workflowGistURL, workflowConfigURL = compileWorkflowWithCRECLI(t, in, feedsConsumerAddress, in.PriceProvider.FeedID, dataURL, settingsFile)
	} else {
		workflowGistURL = in.WorkflowConfig.CompiledWorkflowConfig.BinaryURL
		workflowConfigURL = in.WorkflowConfig.CompiledWorkflowConfig.ConfigURL
	}

	// register the workflow
	registerCmd := exec.Command(CRECLICommand, "workflow", "register", workflowName, "-b", workflowGistURL, "-c", workflowConfigURL, "-S", settingsFile.Name(), "-v")
	registerCmd.Stdout = os.Stdout
	registerCmd.Stderr = os.Stderr
	err = registerCmd.Run()
	require.NoError(t, err, "failed to register workflow using CRE CLI")
}

func startNodes(t *testing.T, in *WorkflowTestConfig, bc *blockchain.Output) *ns.Output {
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
		for _, nodeSpec := range in.NodeSet.NodeSpecs {
			nodeSpec.Node.Image = image
		}
	}

	nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err, "failed to deploy node set")

	return nodeset
}

// In order to whitelist host IP in the gateway, we need to resolve the host.docker.internal to the host IP,
// and since CL image doesn't have dig or nslookup, we need to use curl.
func resolveHostDockerInternaIp(testLogger zerolog.Logger, nsOutput *ns.Output) (string, error) {
	containerName := nsOutput.CLNodes[0].Node.ContainerName
	cmd := []string{"curl", "-v", "http://host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*Trying ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		testLogger.Error().Msgf("failed to extract IP address from curl output:\n%s", output)
		return "", errors.New("failed to extract IP address from curl output")
	}

	testLogger.Info().Msgf("Resolved host.docker.internal to %s", matches[1])

	return matches[1], nil
}

func fundNodes(t *testing.T, don *devenv.DON, sc *seth.Client) {
	for _, node := range don.Nodes {
		_, err := actions.SendFunds(zerolog.Logger{}, sc, actions.FundsToSendPayload{
			ToAddress:  common.HexToAddress(node.AccountAddr[sc.Cfg.Network.ChainID]),
			Amount:     big.NewInt(5000000000000000000),
			PrivateKey: sc.MustGetRootPrivateKey(),
		})
		require.NoError(t, err)
	}
}

func configureNodes(t *testing.T, don *devenv.DON, in *WorkflowTestConfig, bc *blockchain.Output, capRegAddr, workflowRegistryAddr, forwarderAddress common.Address) (*ns.Output, []*clclient.ChainlinkClient) {
	workflowNodeSet := don.Nodes[1:]

	bootstrapNodePeerId, err := nodeToP2PID(don.Nodes[0], keyExtractingTransformFn)
	require.NoError(t, err, "failed to get bootstrap node peer ID")

	// configure the bootstrap node
	in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'
				ContractPollInterval = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				DefaultBootstrappers = ['%s@localhost:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				DefaultBootstrappers = ['%s@localhost:6690']

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'
			`,
		bootstrapNodePeerId,
		bootstrapNodePeerId,
		bc.ChainID,
		bc.Nodes[0].DockerInternalWSUrl,
		bc.Nodes[0].DockerInternalHTTPUrl,
	)

	chainIDInt, err := strconv.Atoi(bc.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := mustSafeUint64(int64(chainIDInt))

	// configure worker nodes with p2p, peering capabilitity (for DON-2-DON communication),
	// capability (external) registry, workflow registry and gateway connector (required for reading from workflow registry and for external communication)
	for i := range workflowNodeSet {
		in.NodeSet.NodeSpecs[i+1].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'
				ContractPollInterval = '1s'

				[P2P.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:5001']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:5001']

				[Capabilities.Peering.V2]
				Enabled = true
				ListenAddresses = ['0.0.0.0:6690']
				# assuming that node0 is the bootstrap node
				DefaultBootstrappers = ['%s@node0:6690']

				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'

				# This is needed for the target capability to be initialized
				[EVM.Workflow]
				FromAddress = '%s'
				ForwarderAddress = '%s'
				GasLimitDefault = 400_000

				[Capabilities.ExternalRegistry]
				Address = '%s'
				NetworkID = 'evm'
				ChainID = '%s'

				[Capabilities.WorkflowRegistry]
				Address = "%s"
				NetworkID = "evm"
				ChainID = "%s"

				[Capabilities.GatewayConnector]
				DonID = "%s"
				ChainIDForNodeKey = "%s"
				NodeAddress = '%s'

				[[Capabilities.GatewayConnector.Gateways]]
				Id = "por_gateway"
				URL = "%s"
			`,
			bootstrapNodePeerId,
			bootstrapNodePeerId,
			bc.ChainID,
			bc.Nodes[0].DockerInternalWSUrl,
			bc.Nodes[0].DockerInternalHTTPUrl,
			workflowNodeSet[i].AccountAddr[chainIDUint64],
			forwarderAddress.Hex(),
			capRegAddr,
			bc.ChainID,
			workflowRegistryAddr.Hex(),
			bc.ChainID,
			strconv.FormatUint(uint64(in.WorkflowConfig.DonID), 10),
			bc.ChainID,
			workflowNodeSet[i].AccountAddr[chainIDUint64],
			"ws://node0:5003/node", // bootstrap node exposes gateway port on 5003
		)
	}

	// we need to restart all nodes for configuration changes to take effect
	nodeset, err := ns.UpgradeNodeSet(t, in.NodeSet, bc, 5*time.Second)
	require.NoError(t, err, "failed to upgrade node set")

	// we need to recreate chainlink clients after the nodes are restarted
	nodeClients, err := clclient.New(nodeset.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	return nodeset, nodeClients
}

func reinitialiseJDClient(t *testing.T, ctfEnv *deployment.Environment, jdOutput *jd.Output, nodeOutput *ns.Output) deployment.Environment {
	nodeInfo, err := getNodeInfo(nodeOutput, 1)
	require.NoError(t, err, "failed to get node info")

	jdConfig := devenv.JDConfig{
		GRPC:     jdOutput.HostGRPCUrl,
		WSRPC:    jdOutput.DockerWSRPCUrl,
		Creds:    insecure.NewCredentials(),
		NodeInfo: nodeInfo,
	}

	offChain, err := devenv.NewJDClient(context.Background(), jdConfig)
	require.NoError(t, err, "failed to create JD client")

	ctfEnv.Offchain = offChain

	return *ctfEnv
}

func mustSafeUint64(input int64) uint64 {
	if input < 0 {
		panic(fmt.Errorf("int64 %d is below uint64 min value", input))
	}
	return uint64(input)
}

func createNodeJobsWithJd(t *testing.T, ctfEnv *deployment.Environment, don *devenv.DON, bc *blockchain.Output, keystoneContractSet keystone_changeset.ContractSet, extraAllowedPorts []int, extraAllowedIps []string) {
	// if there's only one OCR3 contract in the set, we can use `nil` as the address to get its instance
	ocr3Contract, err := keystoneContractSet.GetOCR3Contract(nil)
	require.NoError(t, err, "failed to get OCR3 contract address")

	ocr3CapabilityAddress := ocr3Contract.Address().Hex()

	chainIDInt, err := strconv.Atoi(bc.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := mustSafeUint64(int64(chainIDInt))

	bootstrapNodePeerId, err := nodeToP2PID(don.Nodes[0], keyExtractingTransformFn)
	require.NoError(t, err, "failed to get bootstrap node peer ID")

	jobCount := 2 + (len(don.Nodes)-1)*3
	errCh := make(chan error, jobCount)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		bootstrapJobSpec := fmt.Sprintf(`
				type = "bootstrap"
				schemaVersion = 1
				externalJobID = "%s"
				name = "Botostrap"
				contractID = "%s"
				contractConfigTrackerPollInterval = "1s"
				contractConfigConfirmations = 1
				relay = "evm"
				[relayConfig]
				chainID = %s
				providerType = "ocr3-capability"
			`, uuid.NewString(),
			ocr3CapabilityAddress,
			bc.ChainID)

		bootstrapJobRequest := &jobv1.ProposeJobRequest{
			NodeId: don.Nodes[0].NodeID,
			Spec:   bootstrapJobSpec,
		}

		_, bootErr := ctfEnv.Offchain.ProposeJob(context.Background(), bootstrapJobRequest)
		if bootErr != nil {
			errCh <- errors.Wrapf(bootErr, "failed to propose bootstrap job")
			return
		}

		gatewayJobSpec := fmt.Sprintf(`
				type = "gateway"
				schemaVersion = 1
				externalJobID = "%s"
				name = "Gateway"
				forwardingAllowed = false
				[gatewayConfig.ConnectionManagerConfig]
				AuthChallengeLen = 10
				AuthGatewayId = "por_gateway"
				AuthTimestampToleranceSec = 5
				HeartbeatIntervalSec = 20
				[[gatewayConfig.Dons]]
				DonId = "1"
				F = 1
				HandlerName = "web-api-capabilities"
					[gatewayConfig.Dons.HandlerConfig]
					MaxAllowedMessageAgeSec = 1_000
						[gatewayConfig.Dons.HandlerConfig.NodeRateLimiter]
						GlobalBurst = 10
						GlobalRPS = 50
						PerSenderBurst = 10
						PerSenderRPS = 10
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 1"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 2"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 3"
					[[gatewayConfig.Dons.Members]]
					Address = "%s"
					Name = "Workflow Node 4"
				[gatewayConfig.NodeServerConfig]
				HandshakeTimeoutMillis = 1_000
				MaxRequestBytes = 100_000
				Path = "/node"
				Port = 5_003 #this is the port the other nodes will use to connect to the gateway
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000
				[gatewayConfig.UserServerConfig]
				ContentTypeHeader = "application/jsonrpc"
				MaxRequestBytes = 100_000
				Path = "/"
				Port = 5_002
				ReadTimeoutMillis = 1_000
				RequestTimeoutMillis = 10_000
				WriteTimeoutMillis = 1_000
				[gatewayConfig.HTTPClientConfig]
				MaxResponseBytes = 100_000_000
			`,
			uuid.NewString(),
			// ETH keys of the workflow nodes
			don.Nodes[1].AccountAddr[chainIDUint64],
			don.Nodes[2].AccountAddr[chainIDUint64],
			don.Nodes[3].AccountAddr[chainIDUint64],
			don.Nodes[4].AccountAddr[chainIDUint64],
		)

		if len(extraAllowedPorts) != 0 {
			var allowedPorts string
			for _, port := range extraAllowedPorts {
				allowedPorts += fmt.Sprintf("%d, ", port)
			}

			// when we pass custom allowed IPs, defaults are not used and we need to
			// pass HTTP and HTTPS explicitly
			gatewayJobSpec += fmt.Sprintf(`
				AllowedPorts = [80, 443, %s]
				`,
				allowedPorts,
			)
		}

		if len(extraAllowedIps) != 0 {
			allowedIPs := strings.Join(extraAllowedIps, `", "`)

			gatewayJobSpec += fmt.Sprintf(`
			AllowedIps = ["%s"]
			`,
				allowedIPs,
			)
		}

		gatewayJobRequest := &jobv1.ProposeJobRequest{
			NodeId: don.Nodes[0].NodeID,
			Spec:   gatewayJobSpec,
		}

		_, gateErr := ctfEnv.Offchain.ProposeJob(context.Background(), gatewayJobRequest)
		if gateErr != nil {
			errCh <- errors.Wrapf(gateErr, "failed to propose gateway job for the bootstrap node")
		}
	}()

	// for each capability that's required by the workflow, create a job for workflow each node
	for i, node := range don.Nodes {
		// First node is a bootstrap node, so we skip it
		if i == 0 {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// since we are using a capability that is not bundled-in, we need to copy it to the Docker container
			// and point the job to the copied binary

			// failed to propose job. err: rpc error: code = Internal desc = failed to propose job to node: failed to generate a job based on spec: unknown job type: standardcapabilities
			cronJobSpec := fmt.Sprintf(`
					type = "standardcapabilities"
					schemaVersion = 1
					externalJobID = "%s"
					name = "cron-capabilities"
					forwardingAllowed = false
					command = "/home/capabilities/%s"
					config = ""
				`,
				uuid.NewString(),
				cronCapabilityAssetFile)

			cronJobRequest := &jobv1.ProposeJobRequest{
				NodeId: node.NodeID,
				Spec:   cronJobSpec,
			}

			_, cronErr := ctfEnv.Offchain.ProposeJob(context.Background(), cronJobRequest)
			if cronErr != nil {
				errCh <- errors.Wrapf(cronErr, "failed to propose cron job for node %s", node.NodeID)
				return
			}

			computeJobSpec := fmt.Sprintf(`
					type = "standardcapabilities"
					schemaVersion = 1
					externalJobID = "%s"
					name = "compute-capabilities"
					forwardingAllowed = false
					command = "__builtin_custom-compute-action"
					config = """
					NumWorkers = 3
						[rateLimiter]
						globalRPS = 20.0
						globalBurst = 30
						perSenderRPS = 1.0
						perSenderBurst = 5
					"""
				`,
				uuid.NewString())

			computeJobRequest := &jobv1.ProposeJobRequest{
				NodeId: node.NodeID,
				Spec:   computeJobSpec,
			}

			_, compErr := ctfEnv.Offchain.ProposeJob(context.Background(), computeJobRequest)
			if compErr != nil {
				errCh <- errors.Wrapf(compErr, "failed to propose compute job for node %s", node.NodeID)
				return
			}

			consensusJobSpec := fmt.Sprintf(`
					type = "offchainreporting2"
					schemaVersion = 1
					externalJobID = "%s"
					name = "Keystone OCR3 Consensus Capability"
					contractID = "%s"
					ocrKeyBundleID = "%s"
					p2pv2Bootstrappers = [
						"%s@%s",
					]
					relay = "evm"
					pluginType = "plugin"
					transmitterID = "%s"
					[relayConfig]
					chainID = "%s"
					[pluginConfig]
					command = "/usr/local/bin/chainlink-ocr3-capability"
					ocrVersion = 3
					pluginName = "ocr-capability"
					providerType = "ocr3-capability"
					telemetryType = "plugin"
					[onchainSigningStrategy]
					strategyName = 'multi-chain'
					[onchainSigningStrategy.config]
					evm = "%s"
					`,
				uuid.NewString(),
				ocr3CapabilityAddress,
				node.Ocr2KeyBundleID,
				bootstrapNodePeerId,
				"node0:5001",
				node.AccountAddr[chainIDUint64],
				bc.ChainID,
				node.Ocr2KeyBundleID,
			)

			consensusJobRequest := &jobv1.ProposeJobRequest{
				NodeId: node.NodeID,
				Spec:   consensusJobSpec,
			}

			_, consErr := ctfEnv.Offchain.ProposeJob(context.Background(), consensusJobRequest)
			if consErr != nil {
				errCh <- errors.Wrapf(consErr, "failed to propose consensus job for node %s ", node.NodeID)
			}
		}()
	}
	wg.Wait()

	close(errCh)

	errFound := false
	for err := range errCh {
		errFound = true
		//nolint:testifylint // we want to assert here to catch all errors
		assert.NoError(t, err, "job creation/acception failed")
	}

	require.False(t, errFound, "failed to create at least one job")
}

func noOpTransformFn(value string) string {
	return value
}

func keyExtractingTransformFn(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func nodeToP2PID(node devenv.Node, transformFn func(string) string) (string, error) {
	for _, label := range node.Labels() {
		if label.Key == devenv.NodeLabelP2PIDType {
			if label.Value == nil {
				return "", fmt.Errorf("p2p label value is nil for node %s", node.Name)
			}
			return transformFn(*label.Value), nil
		}
	}

	return "", fmt.Errorf("p2p label not found for node %s", node.Name)
}

func configureWorkflowDON(t *testing.T, ctfEnv *deployment.Environment, don *devenv.DON, chainSelector uint64) {
	kcrAllCaps := []keystone_changeset.DONCapabilityWithConfig{
		{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "offchain_reporting",
				Version:        "1.0.0",
				CapabilityType: 2, // CONSENSUS
				ResponseType:   0, // REPORT
			},
			Config: &capabilitiespb.CapabilityConfig{},
		},
		{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "write_geth-testnet",
				Version:        "1.0.0",
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		},
		{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "cron-trigger",
				Version:        "1.0.0",
				CapabilityType: uint8(0), // trigger
			},
			Config: &capabilitiespb.CapabilityConfig{},
		},
		{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "custom-compute",
				Version:        "1.0.0",
				CapabilityType: uint8(1), // action
			},
			Config: &capabilitiespb.CapabilityConfig{},
		},
	}

	peerIds := make([]string, len(don.Nodes)-1)
	for i, node := range don.Nodes {
		if i == 0 {
			continue
		}

		p2pId, err := nodeToP2PID(node, noOpTransformFn)
		require.NoError(t, err, "failed to get p2p id for node %s", node.Name)

		peerIds[i-1] = p2pId
	}

	nop := keystone_changeset.NOP{
		Name:  "NOP 1",
		Nodes: peerIds,
	}

	donName := "keystone-don"
	donCap := keystone_changeset.DonCapabilities{
		Name:         donName,
		F:            1,
		Nops:         []keystone_changeset.NOP{nop},
		Capabilities: kcrAllCaps,
	}

	transmissionSchedule := []int{len(don.Nodes) - 1}

	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              transmissionSchedule,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationAcceptMillis:           1000,
		MaxDurationTransmitMillis:         1000,
		MaxFaultyOracles:                  1,
		MaxQueryLengthBytes:               1000000,
		MaxObservationLengthBytes:         1000000,
		MaxReportLengthBytes:              1000000,
		MaxRequestBatchSize:               1000,
		UniqueReports:                     true,
	}

	cfg := keystone_changeset.InitialContractsCfg{
		RegistryChainSel: chainSelector,
		Dons:             []keystone_changeset.DonCapabilities{donCap},
		OCR3Config:       &oracleConfig,
	}

	_, err := keystone_changeset.ConfigureInitialContractsChangeset(*ctfEnv, cfg)
	require.NoError(t, err, "failed to configure initial contracts")
}

func startJobDistributor(t *testing.T, in *WorkflowTestConfig) *jd.Output {
	if os.Getenv("CI") == "true" {
		jdImage := ctfconfig.MustReadEnvVar_String(e2eJobDistributorImageEnvVarName)
		jdVersion := os.Getenv(e2eJobDistributorVersionEnvVarName)
		in.JD.Image = fmt.Sprintf("%s:%s", jdImage, jdVersion)
	}
	jdOutput, err := jd.NewJD(in.JD)
	require.NoError(t, err, "failed to create new job distributor")

	return jdOutput
}

func getLogFileHandles(t *testing.T, l zerolog.Logger, ns *ns.Output) ([]*os.File, error) {
	var logFiles []*os.File

	var belongsToCurrentEnv = func(filePath string) bool {
		for i, clNode := range ns.CLNodes {
			if clNode == nil {
				continue
			}

			// skip the first node, as it's the bootstrap node
			if i == 0 {
				continue
			}

			if strings.EqualFold(filePath, clNode.Node.ContainerName+".log") {
				return true
			}
		}
		return false
	}

	logsDir := "logs/docker-" + t.Name()

	fileWalkErr := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && belongsToCurrentEnv(info.Name()) {
			file, fileErr := os.Open(path)
			if fileErr != nil {
				return fmt.Errorf("failed to open file %s: %w", path, fileErr)
			}
			logFiles = append(logFiles, file)
		}
		return nil
	})

	expectedLogCount := len(ns.CLNodes) - 1
	if len(logFiles) != expectedLogCount {
		l.Warn().Int("Expected", expectedLogCount).Int("Got", len(logFiles)).Msg("Number of log files does not match number of worker nodes. Some logs might be missing.")
	}

	if fileWalkErr != nil {
		l.Error().Err(fileWalkErr).Msg("Error walking through log files. Will not look for report transmission transaction hashes")
		return nil, fileWalkErr
	}

	return logFiles, nil
}

// This function is used to go through Chainlink Node logs and look for entries related to report transmissions.
// Once such a log entry is found, it looks for transaction hash and then it tries to decode the transaction and print the result.
func debugReportTransmissions(logFiles []*os.File, l zerolog.Logger, wsRPCURL string) {
	/*
	 Example log entry:
	 2025-01-28T14:44:48.080Z [DEBUG] Node sent transaction                              multinode@v0.0.0-20250121205514-f73e2f86c23b/transaction_sender.go:180 chainID=1337 logger=EVM.1337.TransactionSender tx={"type":"0x0","chainId":"0x539","nonce":"0x0","to":"0xcf7ed3acca5a467e9e704c703e8d87f634fb0fc9","gas":"0x61a80","gasPrice":"0x3b9aca00","maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x0","input":"0x11289565000000000000000000000000a513e6e4b8f2a923d98304ec87f64353c4d5c853000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000010d010f715db03509d388f706e16137722000e26aa650a64ac826ae8e5679cdf57fd96798ed50000000010000000100000a9c593aaed2f5371a5bc0779d1b8ea6f9c7d37bfcbb876a0a9444dbd36f64306466323239353031f39fd6e51aad88f6f4ce6ab8827279cfffb92266000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001018bfe88407000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bb5c162c8000000000000000000000000000000000000000000000000000000006798ed37000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000e700d4c57250eac9dc925c951154c90c1b6017944322fb2075055d8bdbe19000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000041561c171b7465e8efef35572ef82adedb49ea71b8344a34a54ce5e853f80ca1ad7d644ebe710728f21ebfc3e2407bd90173244f744faa011c3a57213c8c585de90000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004165e6f3623acc43f163a58761655841bfebf3f6b4ea5f8d34c64188036b0ac23037ebbd3854b204ca26d828675395c4b9079ca068d9798326eb8c93f26570a1080100000000000000000000000000000000000000000000000000000000000000","v":"0xa96","r":"0x168547e96e7088c212f85a4e8dddce044bbb2abfd5ccc8a5451fdfcb812c94e5","s":"0x2a735a3df046632c2aaa7e583fe161113f3345002e6c9137bbfa6800a63f28a4","hash":"0x3fc5508310f8deef09a46ad594dcc5dc9ba415319ef1dfa3136335eb9e87ff4d"} version=2.19.0@05c05a9

	 What we are looking for:
	 "hash":"0x3fc5508310f8deef09a46ad594dcc5dc9ba415319ef1dfa3136335eb9e87ff4d"
	*/
	reportTransmissionTxHashPattern := regexp.MustCompile(`"hash":"(0x[0-9a-fA-F]+)"`)

	// let's be prudent and assume that in extreme scenario when feed price isn't updated, but
	// transmission is still sent, we might have multiple transmissions per node, and if we want
	// to avoid blocking on the channel, we need to have a higher buffer
	resultsCh := make(chan string, len(logFiles)*4)

	wg := &sync.WaitGroup{}
	for _, f := range logFiles {
		wg.Add(1)
		file := f

		go func() {
			defer wg.Done()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()

				if !strings.Contains(jsonLogLine, "Node sent transaction") {
					continue
				}

				match := reportTransmissionTxHashPattern.MatchString(jsonLogLine)
				if match {
					resultsCh <- reportTransmissionTxHashPattern.FindStringSubmatch(jsonLogLine)[1]
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	if len(resultsCh) == 0 {
		l.Error().Msg("❌ No report transmissions found in Chainlink Node logs.")
		return
	}

	// required as Seth prints transaction traces to stdout with debug level
	_ = os.Setenv(seth.LogLevelEnvVar, "debug")

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(wsRPCURL).
		WithReadOnlyMode().
		WithGethWrappersFolders([]string{"../../../core/gethwrappers/keystone/generated"}). // point Seth to the folder with keystone geth wrappers, so that it can load contract ABIs
		Build()

	if err != nil {
		l.Error().Err(err).Msg("Failed to create seth client")
		return
	}

	for txHash := range resultsCh {
		l.Info().Msgf("🔍 Tracing report transmission transaction %s", txHash)
		// set tracing level to all to trace also successful transactions
		sc.Cfg.TracingLevel = seth.TracingLevel_All
		tx, _, err := sc.Client.TransactionByHash(context.Background(), common.HexToHash(txHash))
		if err != nil {
			l.Warn().Err(err).Msgf("Failed to get transaction by hash %s", txHash)
			continue
		}
		_, decodedErr := sc.DecodeTx(tx)

		if decodedErr != nil {
			l.Error().Err(decodedErr).Msgf("Transmission transaction %s failed due to %s", txHash, decodedErr.Error())
			continue
		}
	}
}

// this function is used to print debug information from Chainlink Node logs
// it checks whether workflow was executing, OCR was executing and whether reports were sent
// and if they were, it traces each report transmission transaction
func printTestDebug(t *testing.T, l zerolog.Logger, ns *ns.Output, wsRPCURL string) {
	logFiles, err := getLogFileHandles(t, l, ns)
	if err != nil {
		l.Error().Err(err).Msg("Failed to get log file handles. No debug information will be printed")
		return
	}

	defer func() {
		for _, f := range logFiles {
			_ = f.Close()
		}
	}()

	l.Info().Msg("🔍 Debug information from Chainlink Node logs:")

	// assuming one bootstrap node
	workflowNodeCount := len(ns.CLNodes) - 1
	if !checkIfWorkflowWasExecuting(logFiles, workflowNodeCount) {
		l.Error().Msg("❌ Workflow was not executing")
		return
	} else {
		l.Info().Msg("✅ Workflow was executing")
	}

	if !checkIfOCRWasExecuting(logFiles, workflowNodeCount) {
		l.Error().Msg("❌ OCR was not executing")
		return
	} else {
		l.Info().Msg("✅ OCR was executing")
	}

	if !checkIfAtLeastOneReportWasSent(logFiles, workflowNodeCount) {
		l.Error().Msg("❌ Reports were not sent")
		return
	} else {
		l.Info().Msg("✅ Reports were sent")

		// debug report transmissions
		debugReportTransmissions(logFiles, l, wsRPCURL)
	}
}

func checkIfLogsHaveText(logFiles []*os.File, bufferSize int, expectedText string, validationFn func(int) bool) bool {
	wg := &sync.WaitGroup{}

	resultsCh := make(chan struct{}, bufferSize)

	for _, f := range logFiles {
		wg.Add(1)
		file := f

		go func() {
			defer func() {
				wg.Done()
				// reset file pointer to the beginning of the file
				// so that subsequent reads start from the beginning
				_, _ = file.Seek(0, io.SeekStart)
			}()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				jsonLogLine := scanner.Text()

				if strings.Contains(jsonLogLine, expectedText) {
					resultsCh <- struct{}{}
					return
				}
			}
		}()
	}

	wg.Wait()
	close(resultsCh)

	var found int
	for range resultsCh {
		found++
	}

	return validationFn(found)
}

func exactCountValidationFn(expected int) func(int) bool {
	return func(found int) bool {
		return found == expected
	}
}

func checkIfWorkflowWasExecuting(logFiles []*os.File, workflowNodeCount int) bool {
	return checkIfLogsHaveText(logFiles, workflowNodeCount, "step request enqueued", exactCountValidationFn(workflowNodeCount))
}

func checkIfOCRWasExecuting(logFiles []*os.File, workflowNodeCount int) bool {
	return checkIfLogsHaveText(logFiles, workflowNodeCount, "✅ committed outcome", exactCountValidationFn(workflowNodeCount))
}

func checkIfAtLeastOneReportWasSent(logFiles []*os.File, workflowNodeCount int) bool {
	// we are looking for "Node sent transaction" log entry, which might appear various times in the logs
	// but most probably not in the logs of all nodes, since they take turns in sending reports
	// our buffer must be large enough to capture all the possible log entries in order to avoid channel blocking
	bufferSize := workflowNodeCount * 4

	return checkIfLogsHaveText(logFiles, bufferSize, "Node sent transaction", func(found int) bool { return found > 0 })
}

func logTestInfo(l zerolog.Logger, feedId, workflowName, feedConsumerAddr, forwarderAddr string) {
	l.Info().Msg("------ Test configuration:")
	l.Info().Msgf("Feed ID: %s", feedId)
	l.Info().Msgf("Workflow name: %s", workflowName)
	l.Info().Msgf("FeedConsumer address: %s", feedConsumerAddr)
	l.Info().Msgf("KeystoneForwarder address: %s", forwarderAddr)
}

func float64ToBigInt(f float64) *big.Int {
	f *= 100

	bigFloat := new(big.Float).SetFloat64(f)

	bigInt := new(big.Int)
	bigFloat.Int(bigInt) // Truncate towards zero

	return bigInt
}

func setupFakeDataProvider(t *testing.T, testLogger zerolog.Logger, in *WorkflowTestConfig, priceIndex *int) string {
	_, err := fake.NewFakeDataProvider(in.PriceProvider.Fake.Input)
	require.NoError(t, err)
	fakeApiPath := "/fake/api/price"
	fakeFinalUrl := fmt.Sprintf("%s:%d%s", framework.HostDockerInternal(), in.PriceProvider.Fake.Port, fakeApiPath)

	getPriceResponseFn := func() map[string]interface{} {
		response := map[string]interface{}{
			"accountName": "TrueUSD",
			"totalTrust":  in.PriceProvider.Fake.Prices[*priceIndex],
			"ripcord":     false,
			"updatedAt":   time.Now().Format(time.RFC3339),
		}

		marshalled, err := json.Marshal(response)
		if err == nil {
			testLogger.Info().Msgf("Returning response: %s", string(marshalled))
		} else {
			testLogger.Info().Msgf("Returning response: %v", response)
		}

		return response
	}

	err = fake.Func("GET", fakeApiPath, func(c *gin.Context) {
		c.JSON(200, getPriceResponseFn())
	})

	require.NoError(t, err, "failed to set up fake data provider")

	return fakeFinalUrl
}

func setupPriceProvider(t *testing.T, testLogger zerolog.Logger, in *WorkflowTestConfig) PriceProvider {
	if in.PriceProvider.Fake != nil {
		return NewFakePriceProvider(t, testLogger, in)
	}

	return NewLivePriceProvider(t, testLogger, in)
}

// PriceProvider abstracts away the logic of checking whether the feed has been correctly updated
// and it also returns port and URL of the price provider. This is so, because when using a mocked
// price provider we need start a separate service and whitelist its port and IP with the gateway job.
// Also, since it's a mocked price provider we can now check whether the feed has been correctly updated
// instead of only checking whether it has some price that's != 0.
type PriceProvider interface {
	URL() string
	NextPrice(price *big.Int, elapsed time.Duration) bool
	CheckPrices()
}

// LivePriceProvider is a PriceProvider implementation that uses a live feed to get the price, typically http://api.real-time-reserves.verinumus.io
type LivePriceProvider struct {
	t            *testing.T
	testLogger   zerolog.Logger
	url          string
	actualPrices []*big.Int
}

func NewLivePriceProvider(t *testing.T, testLogger zerolog.Logger, in *WorkflowTestConfig) PriceProvider {
	return &LivePriceProvider{
		testLogger: testLogger,
		url:        in.PriceProvider.URL,
		t:          t,
	}
}

func (l *LivePriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	l.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
	l.actualPrices = append(l.actualPrices, price)

	// no other price to return, we are done
	return false
}

func (l *LivePriceProvider) URL() string {
	return l.url
}

func (l *LivePriceProvider) CheckPrices() {
	// we don't have a way to check the price in the live feed, so we always assume it's correct
	// as long as it's != 0. And we only wait for the first price to be set.
	require.NotEmpty(l.t, l.actualPrices, "no prices found in the feed")
	require.NotEqual(l.t, l.actualPrices[0], big.NewInt(0), "price found in the feed is 0")
}

// FakePriceProvider is a PriceProvider implementation that uses a mocked feed to get the price
// It returns a configured price sequence and makes sure that the feed has been correctly updated
type FakePriceProvider struct {
	t              *testing.T
	testLogger     zerolog.Logger
	priceIndex     *int
	url            string
	expectedPrices []*big.Int
	actualPrices   []*big.Int
}

func NewFakePriceProvider(t *testing.T, testLogger zerolog.Logger, in *WorkflowTestConfig) PriceProvider {
	priceIndex := ptr.Ptr(0)
	expectedPrices := make([]*big.Int, len(in.PriceProvider.Fake.Prices))
	for i, p := range in.PriceProvider.Fake.Prices {
		// convert float64 to big.Int by multiplying by 100
		// just like the PoR workflow does
		expectedPrices[i] = float64ToBigInt(p)
	}

	return &FakePriceProvider{
		t:              t,
		testLogger:     testLogger,
		expectedPrices: expectedPrices,
		priceIndex:     priceIndex,
		url:            setupFakeDataProvider(t, testLogger, in, priceIndex),
	}
}

func (f *FakePriceProvider) priceAlreadyFound(price *big.Int) bool {
	for _, p := range f.actualPrices {
		if p.Cmp(price) == 0 {
			return true
		}
	}

	return false
}

func (f *FakePriceProvider) NextPrice(price *big.Int, elapsed time.Duration) bool {
	// if price is nil or 0 it means that the feed hasn't been updated yet
	if price == nil || price.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	if !f.priceAlreadyFound(price) {
		f.testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
		f.actualPrices = append(f.actualPrices, price)

		if len(f.actualPrices) == len(f.expectedPrices) {
			// all prices found, nothing more to check
			return false
		} else {
			require.Less(f.t, len(f.actualPrices), len(f.expectedPrices), "more prices found than expected")
			f.testLogger.Info().Msgf("Changing price provider price to %f", f.expectedPrices[len(f.actualPrices)])
			*f.priceIndex = len(f.actualPrices)

			// set new price and continue checking
			return true
		}
	}

	// continue checking, price not updated yet
	return true
}

func (f *FakePriceProvider) CheckPrices() {
	require.EqualValues(f.t, f.expectedPrices, f.actualPrices, "prices found in the feed do not match prices set in the mock")
	f.testLogger.Info().Msgf("All %d mocked prices were found in the feed", len(f.expectedPrices))
}

func (f *FakePriceProvider) URL() string {
	return f.url
}

func extraAllowedPortsAndIps(t *testing.T, testLogger zerolog.Logger, in *WorkflowTestConfig, nodeOutput *ns.Output) ([]string, []int) {
	// no need to allow anything, if we are using live feed
	if in.PriceProvider.Fake == nil {
		return nil, nil
	}

	// we need to explicitly allow the port used by the fake data provider
	// and IP corresponding to host.docker.internal or the IP of the host machine, if we are running on Linux,
	// because that's where the fake data provider is running
	var hostIp string
	var err error

	system := runtime.GOOS
	switch system {
	case "darwin":
		hostIp, err = resolveHostDockerInternaIp(testLogger, nodeOutput)
		require.NoError(t, err, "failed to resolve host.docker.internal IP")
	case "linux":
		// for linux framework already returns an IP, so we don't need to resolve it,
		// but we need to remove the http:// prefix
		hostIp = strings.ReplaceAll(framework.HostDockerInternal(), "http://", "")
	default:
		err = fmt.Errorf("unsupported OS: %s", system)
	}
	require.NoError(t, err, "failed to resolve host.docker.internal IP")

	testLogger.Info().Msgf("Will allow IP %s and port %d for the fake data provider", hostIp, in.PriceProvider.Fake.Port)

	// we also need to explicitly allow Gist's IP
	return []string{hostIp, GistIP}, []int{in.PriceProvider.Fake.Port}
}

/*
!!! ATTENTION !!!

Do not use this test as a template for your tests. It's hacky, since we were working under time pressure. We will soon refactor it follow best practices
and a golden example. Apart from its structure what is currently missing is:
- DON-2-DON support
- better structured and reusable methods
*/
func TestKeystoneWithOCR3Workflow(t *testing.T) {
	testLogger := framework.L

	// we need to use double-pointers, so that what's captured in the cleanup function is a pointer, not the actual object,
	// which is only set later in the test, after the cleanup function is defined
	var nodes **ns.Output
	var wsRPCURL *string

	// clean up is LIFO, so we need to make sure we execute the debug report transmission after logs are written down
	// by function added to clean up by framework.Load() method
	t.Cleanup(func() {
		if t.Failed() {
			if nodes == nil {
				testLogger.Warn().Msg("nodeset output is nil, skipping debug report transmission")
				return
			}
			printTestDebug(t, testLogger, *nodes, *wsRPCURL)
		}
	})

	// Load test configuration
	in, err := framework.Load[WorkflowTestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateInputsAndEnvVars(t, in)

	pkey := os.Getenv("PRIVATE_KEY")

	// Create a new blockchain network and Seth client to interact with it
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err, "failed to create seth client")

	// Get either a no-op price provider (for live endpoint)
	// or a fake price provider (for mock endpoint)
	priceProvider := setupPriceProvider(t, testLogger, in)

	// Start job distributor
	jdOutput := startJobDistributor(t, in)

	// Deploy the DON
	nodeOutput := startNodes(t, in, bc)

	// Prepare the chainlink/deployment environment
	ctfEnv, don, chainSelector := buildChainlinkDeploymentEnv(t, jdOutput, nodeOutput, bc, sc)

	// Fund the nodes
	fundNodes(t, don, sc)

	// Deploy keystone contracts (forwarder, capability registry, ocr3 capability)
	keystoneContractSet := deployKeystoneContracts(t, testLogger, ctfEnv, chainSelector)

	// Deploy and pre-configure workflow registry contract
	workflowRegistryAddr := prepareWorkflowRegistry(t, testLogger, ctfEnv, chainSelector, sc, in.WorkflowConfig.DonID)

	// Deploy and configure Keystone Feeds Consumer contract
	feedsConsumerAddress := prepareFeedsConsumer(t, testLogger, ctfEnv, chainSelector, sc, keystoneContractSet.Forwarder.Address(), in.WorkflowConfig.WorkflowName)

	// Register the workflow (either via CRE CLI or by calling the workflow registry directly)
	registerWorkflow(t, in, sc, keystoneContractSet.CapabilitiesRegistry.Address(), workflowRegistryAddr, feedsConsumerAddress, in.WorkflowConfig.DonID, chainSelector, in.WorkflowConfig.WorkflowName, pkey, bc.Nodes[0].HostHTTPUrl, priceProvider.URL())

	// Create OCR3 and capability jobs for each node JD
	ns, _ := configureNodes(t, don, in, bc, keystoneContractSet.CapabilitiesRegistry.Address(), workflowRegistryAddr, keystoneContractSet.Forwarder.Address())
	// JD client needs to be reinitialised after restarting nodes
	ctfEnv = ptr.Ptr(reinitialiseJDClient(t, ctfEnv, jdOutput, nodeOutput))

	ips, ports := extraAllowedPortsAndIps(t, testLogger, in, ns)
	createNodeJobsWithJd(t, ctfEnv, don, bc, keystoneContractSet, ports, ips)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.PriceProvider.FeedID, in.WorkflowConfig.WorkflowName, feedsConsumerAddress.Hex(), keystoneContractSet.Forwarder.Address().Hex())
		}
	})

	// set variables that are needed for the cleanup function, which debugs report transmissions
	nodes = &ns
	wsRPCURL = &bc.Nodes[0].HostWSUrl

	// CAUTION: It is crucial to configure OCR3 jobs on nodes before configuring the workflow contracts.
	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure the workflow DON and contracts
	configureWorkflowDON(t, ctfEnv, don, chainSelector)

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	testLogger.Info().Msg("Waiting for feed to update...")
	startTime := time.Now()
	feedBytes := common.HexToHash(in.PriceProvider.FeedID)

	for {
		select {
		case <-ctx.Done():
			testLogger.Error().Msgf("feed did not update, timeout after %s", timeout)
			t.FailNow()
		case <-time.After(10 * time.Second):
			elapsed := time.Since(startTime).Round(time.Second)
			price, _, err := feedsConsumerInstance.GetPrice(
				sc.NewCallOpts(),
				feedBytes,
			)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			if !priceProvider.NextPrice(price, elapsed) {
				// check if all expected prices were found and finish the test
				priceProvider.CheckPrices()
				return
			}
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}
	}
}
