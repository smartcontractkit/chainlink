package capabilities_test

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-yaml/yaml"
	"github.com/google/go-github/v41/github"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	geth_types "github.com/ethereum/go-ethereum/core/types"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"

	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	cr_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/workflowregistry"
)

// Copying this to avoid dependency on the core repo
func GetChainType(chainType string) (uint8, error) {
	switch chainType {
	case "evm":
		return 1, nil
	// case Solana:
	// 	return 2, nil
	// case Cosmos:
	// 	return 3, nil
	// case StarkNet:
	// 	return 4, nil
	// case Aptos:
	// 	return 5, nil
	default:
		return 0, fmt.Errorf("unexpected chaintype.ChainType: %#v", chainType)
	}
}

// Copying this to avoid dependency on the core repo
func MarshalMultichainPublicKey(ost map[string]types.OnchainPublicKey) (types.OnchainPublicKey, error) {
	pubKeys := make([][]byte, 0, len(ost))
	for k, pubKey := range ost {
		typ, err := GetChainType(k)
		if err != nil {
			// skipping unknown key type
			continue
		}
		buf := new(bytes.Buffer)
		if err = binary.Write(buf, binary.LittleEndian, typ); err != nil {
			return nil, err
		}
		length := len(pubKey)
		if length < 0 || length > math.MaxUint16 {
			return nil, errors.New("pubKey doesn't fit into uint16")
		}
		if err = binary.Write(buf, binary.LittleEndian, uint16(length)); err != nil {
			return nil, err
		}
		_, _ = buf.Write(pubKey)
		pubKeys = append(pubKeys, buf.Bytes())
	}
	// sort keys based on encoded type to make encoding deterministic
	slices.SortFunc(pubKeys, func(a, b []byte) int { return cmp.Compare(a[0], b[0]) })
	return bytes.Join(pubKeys, nil), nil
}

type WorkflowConfig struct {
	UseChainlinkCLI bool                    `toml:"use_chainlink_cli"`
	ChainlinkCLI    *ChainlinkCLIConfig     `toml:"chainlink_cli"`
	UseExising      bool                    `toml:"use_existing"`
	Existing        *ExistingWorkflowConfig `toml:"existing"`
}

type ExistingWorkflowConfig struct {
	BinaryURL string `toml:"binary_url"`
	ConfigURL string `toml:"config_url"`
}

type ChainlinkCLIConfig struct {
	FolderLocation *string `toml:"folder_location"`
}

type WorkflowTestConfig struct {
	BlockchainA    *blockchain.Input `toml:"blockchain_a" validate:"required"`
	NodeSet        *ns.Input         `toml:"nodeset" validate:"required"`
	WorkflowConfig *WorkflowConfig   `toml:"workflow_config" validate:"required"`
}

type OCR3Config struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type NodeInfo struct {
	OcrKeyBundleID            string
	TransmitterAddress        string
	PeerID                    string
	Signer                    common.Address
	OffchainPublicKey         [32]byte
	OnchainPublicKey          types.OnchainPublicKey
	ConfigEncryptionPublicKey [32]byte
}

func extractKey(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return value
}

func downloadGHAssetFromLatestRelease(owner, repository, releaseType, assetName, ghToken string) ([]byte, error) {
	var content []byte
	if ghToken == "" {
		return content, errors.New("no github token provided")
	}

	if (releaseType == test_env.AUTOMATIC_LATEST_TAG) || (releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG) {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ghToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		ghClient := github.NewClient(tc)

		latestTags, _, err := ghClient.Repositories.ListReleases(context.Background(), owner, repository, &github.ListOptions{PerPage: 20})
		if err != nil {
			return content, errors.Wrapf(err, "failed to list releases for %s", repository)
		}

		var latestRelease *github.RepositoryRelease
		for _, tag := range latestTags {
			if releaseType == test_env.AUTOMATIC_STABLE_LATEST_TAG {
				if tag.Prerelease != nil && *tag.Prerelease {
					continue
				}
				if tag.Draft != nil && *tag.Draft {
					continue
				}
			}
			if tag.TagName != nil {
				latestRelease = tag
				break
			}
		}

		if latestRelease == nil {
			return content, errors.New("failed to find latest release with automatic tag: " + releaseType)
		}

		var assetID int64
		for _, asset := range latestRelease.Assets {
			if strings.Contains(asset.GetName(), assetName) {
				assetID = asset.GetID()
				break
			}
		}

		if assetID == 0 {
			return content, fmt.Errorf("failed to find asset %s for %s", assetName, *latestRelease.TagName)
		}

		asset, _, err := ghClient.Repositories.DownloadReleaseAsset(context.Background(), owner, repository, assetID, tc)
		if err != nil {
			return content, errors.Wrapf(err, "failed to download asset %s for %s", assetName, *latestRelease.TagName)
		}

		content, err = io.ReadAll(asset)
		if err != nil {
			return content, err
		}

		return content, nil
	}

	return content, errors.New("no automatic tag provided")
}

func getNodesInfo(
	t *testing.T,
	nodes []*clclient.ChainlinkClient,
) (nodesInfo []NodeInfo) {
	nodesInfo = make([]NodeInfo, len(nodes))

	for i, node := range nodes {
		// OCR Keys
		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err)
		nodesInfo[i].OcrKeyBundleID = ocr2Keys.Data[0].ID

		firstOCR2Key := ocr2Keys.Data[0].Attributes
		nodesInfo[i].Signer = common.HexToAddress(extractKey(firstOCR2Key.OnChainPublicKey))

		pubKeys := make(map[string]types.OnchainPublicKey)
		ethOnchainPubKey, err := hex.DecodeString(extractKey(firstOCR2Key.OnChainPublicKey))
		require.NoError(t, err)
		pubKeys["evm"] = ethOnchainPubKey

		multichainPubKey, err := MarshalMultichainPublicKey(pubKeys)
		require.NoError(t, err)
		nodesInfo[i].OnchainPublicKey = multichainPubKey

		offchainPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.OffChainPublicKey))
		require.NoError(t, err)
		var offchainPublicKey [32]byte
		copy(offchainPublicKey[:], offchainPublicKeyBytes)
		nodesInfo[i].OffchainPublicKey = offchainPublicKey

		sharedSecretEncryptionPublicKeyBytes, err := hex.DecodeString(extractKey(firstOCR2Key.ConfigPublicKey))
		require.NoError(t, err)
		var sharedSecretEncryptionPublicKey [32]byte
		copy(sharedSecretEncryptionPublicKey[:], sharedSecretEncryptionPublicKeyBytes)
		nodesInfo[i].ConfigEncryptionPublicKey = sharedSecretEncryptionPublicKey

		// ETH Keys
		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		nodesInfo[i].TransmitterAddress = ethKeys.Data[0].Attributes.Address

		// P2P Keys
		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err)
		nodesInfo[i].PeerID = p2pKeys.Data[0].Attributes.PeerID
	}

	return nodesInfo
}

func generateOCR3Config(
	t *testing.T,
	nodesInfo []NodeInfo,
) (config *OCR3Config) {
	oracleIdentities := []confighelper.OracleIdentityExtra{}
	transmissionSchedule := []int{}

	for _, nodeInfo := range nodesInfo {
		transmissionSchedule = append(transmissionSchedule, 1)
		oracleIdentity := confighelper.OracleIdentityExtra{}
		oracleIdentity.OffchainPublicKey = nodeInfo.OffchainPublicKey
		oracleIdentity.OnchainPublicKey = nodeInfo.OnchainPublicKey
		oracleIdentity.ConfigEncryptionPublicKey = nodeInfo.ConfigEncryptionPublicKey
		oracleIdentity.PeerID = nodeInfo.PeerID
		oracleIdentity.TransmitAccount = types.Account(nodeInfo.TransmitterAddress)
		oracleIdentities = append(oracleIdentities, oracleIdentity)
	}

	maxDurationInitialization := 10 * time.Second

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		5*time.Second,              // DeltaProgress: Time between rounds
		5*time.Second,              // DeltaResend: Time between resending unconfirmed transmissions
		5*time.Second,              // DeltaInitial: Initial delay before starting the first round
		2*time.Second,              // DeltaRound: Time between rounds within an epoch
		500*time.Millisecond,       // DeltaGrace: Grace period for delayed transmissions
		1*time.Second,              // DeltaCertifiedCommitRequest: Time between certified commit requests
		30*time.Second,             // DeltaStage: Time between stages of the protocol
		uint64(10),                 // MaxRoundsPerEpoch: Maximum number of rounds per epoch
		transmissionSchedule,       // TransmissionSchedule: Transmission schedule
		oracleIdentities,           // Oracle identities with their public keys
		nil,                        // Plugin config (empty for now)
		&maxDurationInitialization, // MaxDurationInitialization: ???
		1*time.Second,              // MaxDurationQuery: Maximum duration for querying
		1*time.Second,              // MaxDurationObservation: Maximum duration for observation
		1*time.Second,              // MaxDurationAccept: Maximum duration for acceptance
		1*time.Second,              // MaxDurationTransmit: Maximum duration for transmission
		1,                          // F: Maximum number of faulty oracles
		nil,                        // OnChain config (empty for now)
	)
	require.NoError(t, err)

	signerAddresses := [][]byte{}
	for _, signer := range signers {
		signerAddresses = append(signerAddresses, signer)
	}

	transmitterAddresses := []common.Address{}
	for _, transmitter := range transmitters {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(transmitter)))
	}

	return &OCR3Config{
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
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

type ChainlinkCliSettings struct {
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
	chainlinkCliAssetFile   = "cre_v1.0.2_linux_amd64.tar.gz"
	cronCapabilityAssetFile = "amd64_cron"
)

func downloadAndInstallChainlinkCLI(ghToken string) error {
	content, err := downloadGHAssetFromLatestRelease("smartcontractkit", "dev-platform", test_env.AUTOMATIC_LATEST_TAG, chainlinkCliAssetFile, ghToken)
	if err != nil {
		return err
	}

	tmpfile, err := os.CreateTemp("", chainlinkCliAssetFile)
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write(content); err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xvf", tmpfile.Name()) // #nosec G204
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("chmod", "+x", "chainlink-cli")
	err = cmd.Run()

	if err != nil {
		return err
	}

	if isInstalled := isInstalled("chainlink-cli"); !isInstalled {
		return errors.New("failed to install chainlink-cli or it is not available in the PATH")
	}

	return nil
}

func downloadCronCapability(ghToken string) (string, error) {
	content, err := downloadGHAssetFromLatestRelease("smartcontractkit", "capabilities", test_env.AUTOMATIC_LATEST_TAG, cronCapabilityAssetFile, ghToken)
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

func validateInputsAndEnvVars(t *testing.T, testConfig *WorkflowTestConfig) {
	require.NotEmpty(t, os.Getenv("PRIVATE_KEY"), "PRIVATE_KEY env var must be set")
	if !testConfig.WorkflowConfig.UseChainlinkCLI {
		require.True(t, testConfig.WorkflowConfig.UseExising, "if you are not using chainlink-cli you must use an existing workflow")
	}

	ghToken := os.Getenv("GITHUB_API_TOKEN")
	_, err := downloadCronCapability(ghToken)
	require.NoError(t, err, "failed to download cron capability. Make sure token has content:read permissions to the capabilities repo")

	// TODO this part should ideally happen outside of the test, but due to how our reusable e2e test workflow is structured now
	// we cannot execute this part in workflow steps (it doesn't support any pre-execution hooks)
	if os.Getenv("IS_CI") == "true" {
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV)
		require.NotEmpty(t, os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV), "missing env var: "+ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV)

		if testConfig.WorkflowConfig.UseChainlinkCLI {
			err = downloadAndInstallChainlinkCLI(ghToken)
			require.NoError(t, err, "failed to download and install chainlink-cli. Make sure token has content:read permissions to the dev-platform repo")
		}
	}

	if testConfig.WorkflowConfig.UseChainlinkCLI {
		require.True(t, isInstalled("chainlink-cli"), "chainlink-cli is required for this test. Please install it, add to path and run again")

		if !testConfig.WorkflowConfig.UseExising {
			require.NotEmpty(t, os.Getenv("GITHUB_API_TOKEN"), "GITHUB_API_TOKEN must be set to use chainlink-cli. It requires gist:read and gist:write permissions")
		} else {
			require.NotEmpty(t, testConfig.WorkflowConfig.ChainlinkCLI.FolderLocation, "folder_location must be set in the chainlink_cli config")
		}
	}
}

func buildChainlinkDeploymentEnv(t *testing.T, sc *seth.Client) (*deployment.Environment, uint64) {
	lgr := logger.TestLogger(t)

	addressBook := deployment.NewMemoryAddressBook()
	chainMap := make(map[uint64]deployment.Chain)
	ctx := context.Background()

	chainSelector, err := chainselectors.SelectorFromChainId(sc.Cfg.Network.ChainID)
	require.NoError(t, err, "failed to get chain selector for chain id %d", sc.Cfg.Network.ChainID)
	chainMap[chainSelector] = deployment.Chain{
		Selector:    chainSelector,
		Client:      sc.Client,
		DeployerKey: sc.NewTXOpts(seth.WithNonce(nil)), // set nonce to nil, so that it will be fetched from the chain
		Confirm: func(tx *geth_types.Transaction) (uint64, error) {
			decoded, revertErr := sc.DecodeTx(tx)
			if revertErr != nil {
				return 0, revertErr
			}
			if decoded.Receipt == nil {
				return 0, fmt.Errorf("no receipt found for transaction %s even though it wasn't reverted. This should not happen", tx.Hash().String())
			}
			return decoded.Receipt.BlockNumber.Uint64(), nil
		},
	}

	return deployment.NewEnvironment("ctfV2", lgr, addressBook, chainMap, nil, nil, nil, func() context.Context { return ctx }, deployment.OCRSecrets{}), chainSelector
}

func prepareCapabilitiesRegistry(t *testing.T, sc *seth.Client, allCaps []cr_wrapper.CapabilitiesRegistryCapability) (common.Address, [][32]byte) {
	capRegAddr, tx, capabilitiesRegistryInstance, err := cr_wrapper.DeployCapabilitiesRegistry(sc.NewTXOpts(), sc.Client)
	_, decodeErr := sc.Decode(tx, err)
	require.NoError(t, decodeErr, "failed to deploy capabilities registry contract")

	_, decodeErr = sc.Decode(capabilitiesRegistryInstance.AddCapabilities(
		sc.NewTXOpts(),
		allCaps,
	))
	require.NoError(t, decodeErr, "failed to add capabilities to capabilities registry")

	hashedCapabilities := make([][32]byte, len(allCaps))
	for i, capability := range allCaps {
		hashed, err := capabilitiesRegistryInstance.GetHashedCapabilityId(
			sc.NewCallOpts(),
			capability.LabelledName,
			capability.Version,
		)
		require.NoError(t, err, "failed to get hashed capability ID for %s", capability.LabelledName)
		hashedCapabilities[i] = hashed
	}

	return capRegAddr, hashedCapabilities
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

func configureKeystoneForwarder(t *testing.T, forwarderAddress common.Address, sc *seth.Client, nodesInfo []NodeInfo) {
	forwarderInstance, err := forwarder.NewKeystoneForwarder(forwarderAddress, sc.Client)
	require.NoError(t, err, "failed to create forwarder instance")

	signers := make([]common.Address, len(nodesInfo)-1)

	for i, node := range nodesInfo {
		// skip the first node, as it's the bootstrap node
		// it doesn't have any capabilities that are required by the workflow
		if i == 0 {
			continue
		}
		signers[i-1] = node.Signer
	}

	_, err = sc.Decode(forwarderInstance.SetConfig(
		sc.NewTXOpts(),
		1, // donID
		1, // configVersion -- wonder what it does
		1, // maximum number of faulty nodes
		signers))
	require.NoError(t, err, "failed to set config for forwarder")
}

func configureOCR3Capability(t *testing.T, ocr3CapabilityAddress common.Address, sc *seth.Client, nodeInfo []NodeInfo) {
	workflowNodesetInfo := nodeInfo[1:]

	ocr3CapabilityContract, err := ocr3_capability.NewOCR3Capability(ocr3CapabilityAddress, sc.Client)
	require.NoError(t, err, "failed to create OCR3 capability contract instance")

	ocr3Config := generateOCR3Config(t, workflowNodesetInfo)
	_, decodeErr := sc.Decode(ocr3CapabilityContract.SetConfig(
		sc.NewTXOpts(),
		ocr3Config.Signers,
		ocr3Config.Transmitters,
		ocr3Config.F,
		ocr3Config.OnchainConfig,
		ocr3Config.OffchainConfigVersion,
		ocr3Config.OffchainConfig,
	))
	require.NoError(t, decodeErr, "failed to set OCR3 configuration")
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
			testLogger.Info().Msgf("Deployed FeedConsumer contract at %s", feedsConsumerAddress.Hex())
			feedsConsumerAddress = common.HexToAddress(addrStr)
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

func deployOCR3Capability(t *testing.T, testLogger zerolog.Logger, sc *seth.Client) common.Address {
	ocr3CapabilityAddress, tx, _, err := ocr3_capability.DeployOCR3Capability(
		sc.NewTXOpts(),
		sc.Client,
	)
	_, decodeErr := sc.Decode(tx, err)
	require.NoError(t, decodeErr, "failed to deploy OCR Capability contract")

	testLogger.Info().Msgf("Deployed OCR3 Capability contract at %s", ocr3CapabilityAddress.Hex())

	return ocr3CapabilityAddress
}
func registerWorkflowDirectly(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, workflowRegistryAddr common.Address, donID uint32, workflowName string) {
	require.NotEmpty(t, in.WorkflowConfig.Existing.BinaryURL)
	workFlowData, err := downloadAndDecode(in.WorkflowConfig.Existing.BinaryURL)
	require.NoError(t, err, "failed to download and decode workflow binary")

	var configData []byte
	if in.WorkflowConfig.Existing.ConfigURL != "" {
		configData, err = download(in.WorkflowConfig.Existing.ConfigURL)
		require.NoError(t, err, "failed to download workflow config")
	}

	// use non-encoded workflow name
	workflowID, idErr := GenerateWorkflowIDFromStrings(sc.MustGetRootKeyAddress().Hex(), workflowName, workFlowData, configData, "")
	require.NoError(t, idErr, "failed to generate workflow ID")

	workflowRegistryInstance, err := workflow_registry_wrapper.NewWorkflowRegistry(workflowRegistryAddr, sc.Client)
	require.NoError(t, err, "failed to create workflow registry instance")

	// use non-encoded workflow name
	_, decodeErr := sc.Decode(workflowRegistryInstance.RegisterWorkflow(sc.NewTXOpts(), workflowName, [32]byte(common.Hex2Bytes(workflowID)), donID, uint8(0), in.WorkflowConfig.Existing.BinaryURL, in.WorkflowConfig.Existing.ConfigURL, ""))
	require.NoError(t, decodeErr, "failed to register workflow")
}

//revive:disable // ignore confusing-results
func compileWorkflowWithChainlinkCli(t *testing.T, in *WorkflowTestConfig, feedsConsumerAddress common.Address, settingsFile *os.File) (string, string) {
	feedID := "0x018BFE88407000400000000000000000"

	configFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err, "failed to create workflow config file")

	workflowConfig := PoRWorkflowConfig{
		FeedID:          feedID,
		URL:             "https://api.real-time-reserves.verinumus.io/v1/chainlink/proof-of-reserves/TrueUSD",
		ConsumerAddress: feedsConsumerAddress.Hex(),
	}

	configMarshalled, err := json.Marshal(workflowConfig)
	require.NoError(t, err, "failed to marshal workflow config")

	_, err = configFile.Write(configMarshalled)
	require.NoError(t, err, "failed to write workflow config file")

	var outputBuffer bytes.Buffer

	compileCmd := exec.Command("chainlink-cli", "workflow", "compile", "-S", settingsFile.Name(), "-c", configFile.Name(), "main.go") // #nosec G204
	compileCmd.Stdout = &outputBuffer
	compileCmd.Stderr = &outputBuffer
	compileCmd.Dir = *in.WorkflowConfig.ChainlinkCLI.FolderLocation
	err = compileCmd.Start()
	require.NoError(t, err, "failed to start compile command")

	err = compileCmd.Wait()
	require.NoError(t, err, "failed to wait for compile command")

	fmt.Println("Compile output:\n", outputBuffer.String())

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

func preapreChainlinkCliSettingsFile(t *testing.T, sc *seth.Client, capRegAddr, workflowRegistryAddr common.Address, donID uint32, chainSelector uint64, rpcHTTPURL string) *os.File {
	// create chainlink-cli settings file
	settingsFile, err := os.CreateTemp("", ".chainlink-cli-settings.yaml")
	require.NoError(t, err, "failed to create chainlink-cli settings file")

	settings := ChainlinkCliSettings{
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
	require.NoError(t, err, "failed to marshal chainlink-cli settings")

	_, err = settingsFile.Write(settingsMarshalled)
	require.NoError(t, err, "failed to write chainlink-cli settings file")

	return settingsFile
}

func registerWorkflow(t *testing.T, in *WorkflowTestConfig, sc *seth.Client, capRegAddr, workflowRegistryAddr, feedsConsumerAddress common.Address, donID uint32, chainSelector uint64, workflowName, pkey, rpcHTTPURL string) {
	// Register workflow directly using the provided binary and config URLs
	// This is a legacy solution, probably we can remove it soon
	if in.WorkflowConfig.UseExising && !in.WorkflowConfig.UseChainlinkCLI {
		registerWorkflowDirectly(t, in, sc, workflowRegistryAddr, donID, workflowName)

		return
	}

	// These two env vars are required by the chainlink-cli
	err := os.Setenv("WORKFLOW_OWNER_ADDRESS", sc.MustGetRootKeyAddress().Hex())
	require.NoError(t, err, "failed to set WORKFLOW_OWNER_ADDRESS env var")

	err = os.Setenv("ETH_PRIVATE_KEY", pkey)
	require.NoError(t, err, "failed to set ETH_PRIVATE_KEY env var")

	// create chainlink-cli settings file
	settingsFile := preapreChainlinkCliSettingsFile(t, sc, capRegAddr, workflowRegistryAddr, donID, chainSelector, rpcHTTPURL)

	var workflowGistURL string
	var workflowConfigURL string

	// compile and upload the workflow, if we are not using an existing one
	if !in.WorkflowConfig.UseExising {
		workflowGistURL, workflowConfigURL = compileWorkflowWithChainlinkCli(t, in, feedsConsumerAddress, settingsFile)
	} else {
		workflowGistURL = in.WorkflowConfig.Existing.BinaryURL
		workflowConfigURL = in.WorkflowConfig.Existing.ConfigURL
	}

	// register the workflow
	registerCmd := exec.Command("chainlink-cli", "workflow", "register", workflowName, "-b", workflowGistURL, "-c", workflowConfigURL, "-S", settingsFile.Name(), "-v")
	registerCmd.Stdout = os.Stdout
	registerCmd.Stderr = os.Stderr
	err = registerCmd.Run()
	require.NoError(t, err, "failed to register workflow using chainlink-cli")
}

func starAndFundNodes(t *testing.T, in *WorkflowTestConfig, bc *blockchain.Output, sc *seth.Client) (*ns.Output, []NodeInfo) {
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("IS_CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
		for _, nodeSpec := range in.NodeSet.NodeSpecs {
			nodeSpec.Node.Image = image
		}
	}

	nodeset, err := ns.NewSharedDBNodeSet(in.NodeSet, bc)
	require.NoError(t, err, "failed to deploy node set")

	nodeClients, err := clclient.New(nodeset.CLNodes)
	require.NoError(t, err, "failed to create chainlink clients")

	nodesInfo := getNodesInfo(t, nodeClients)

	// Fund all nodes
	for _, nodeInfo := range nodesInfo {
		_, err := actions.SendFunds(zerolog.Logger{}, sc, actions.FundsToSendPayload{
			ToAddress:  common.HexToAddress(nodeInfo.TransmitterAddress),
			Amount:     big.NewInt(5000000000000000000),
			PrivateKey: sc.MustGetRootPrivateKey(),
		})
		require.NoError(t, err)
	}

	return nodeset, nodesInfo
}

func configureNodes(t *testing.T, nodesInfo []NodeInfo, in *WorkflowTestConfig, bc *blockchain.Output, capRegAddr common.Address, workflowRegistryAddr common.Address, forwarderAddress common.Address) (*ns.Output, []*clclient.ChainlinkClient) {
	bootstrapNodeInfo := nodesInfo[0]
	workflowNodesetInfo := nodesInfo[1:]

	// configure the bootstrap node
	in.NodeSet.NodeSpecs[0].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

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
		bootstrapNodeInfo.PeerID,
		bootstrapNodeInfo.PeerID,
		bc.ChainID,
		bc.Nodes[0].DockerInternalWSUrl,
		bc.Nodes[0].DockerInternalHTTPUrl,
	)

	// configure worker nodes with p2p, peering capabilitity (for DON-2-DON communication),
	// capability (external) registry, workflow registry and gateway connector (required for reading from workflow registry and for external communication)
	for i := range workflowNodesetInfo {
		in.NodeSet.NodeSpecs[i+1].Node.TestConfigOverrides = fmt.Sprintf(`
				[Feature]
				LogPoller = true

				[OCR2]
				Enabled = true
				DatabaseTimeout = '1s'

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

				# This is needed for the target capability to be initialized
				[[EVM]]
				ChainID = '%s'

				[[EVM.Nodes]]
				Name = 'anvil'
				WSURL = '%s'
				HTTPURL = '%s'

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
				DonID = "1"
				ChainIDForNodeKey = "%s"
				NodeAddress = '%s'

				[[Capabilities.GatewayConnector.Gateways]]
				Id = "por_gateway"
				URL = "%s"
			`,
			bootstrapNodeInfo.PeerID,
			bootstrapNodeInfo.PeerID,
			bc.ChainID,
			bc.Nodes[0].DockerInternalWSUrl,
			bc.Nodes[0].DockerInternalHTTPUrl,
			workflowNodesetInfo[i].TransmitterAddress,
			forwarderAddress.Hex(),
			capRegAddr,
			bc.ChainID,
			workflowRegistryAddr.Hex(),
			bc.ChainID,
			bc.ChainID,
			workflowNodesetInfo[i].TransmitterAddress,
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

func createNodeJobs(t *testing.T, nodeClients []*clclient.ChainlinkClient, nodesInfo []NodeInfo, bc *blockchain.Output, ocr3CapabilityAddress common.Address) {
	bootstrapNodeInfo := nodesInfo[0]
	workflowNodesetInfo := nodesInfo[1:]

	// Create gateway and bootstrap (ocr3) jobs for the bootstrap node
	bootstrapNode := nodeClients[0]
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		bootstrapJobSpec := fmt.Sprintf(`
				type = "bootstrap"
				schemaVersion = 1
				name = "Botostrap"
				contractID = "%s"
				contractConfigTrackerPollInterval = "1s"
				contractConfigConfirmations = 1
				relay = "evm"

				[relayConfig]
				chainID = %s
				providerType = "ocr3-capability"
			`, ocr3CapabilityAddress, bc.ChainID)
		r, _, bootErr := bootstrapNode.CreateJobRaw(bootstrapJobSpec)
		assert.NoError(t, bootErr, "failed to create bootstrap job for the bootstrap node")
		assert.Empty(t, r.Errors, "failed to create bootstrap job for the bootstrap node")

		gatewayJobSpec := fmt.Sprintf(`
				type = "gateway"
				schemaVersion = 1
				name = "PoR Gateway"
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
			// ETH keys of the workflow nodes
			workflowNodesetInfo[0].TransmitterAddress,
			workflowNodesetInfo[1].TransmitterAddress,
			workflowNodesetInfo[2].TransmitterAddress,
			workflowNodesetInfo[3].TransmitterAddress,
		)

		r, _, gatewayErr := bootstrapNode.CreateJobRaw(gatewayJobSpec)
		assert.NoError(t, gatewayErr, "failed to create gateway job for the bootstrap node")
		assert.Empty(t, r.Errors, "failed to create gateway job for the bootstrap node")
	}()

	// for each capability that's required by the workflow, create a job for workflow each node
	for i, nodeClient := range nodeClients {
		// First node is a bootstrap node, so we skip it
		if i == 0 {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			// since we are using a capability that is not bundled-in, we need to copy it to the Docker container
			// and point the job to the copied binary
			cronJobSpec := fmt.Sprintf(`
					type = "standardcapabilities"
					schemaVersion = 1
					name = "cron-capabilities"
					forwardingAllowed = false
					command = "/home/capabilities/%s"
					config = ""
				`, cronCapabilityAssetFile)

			response, _, errCron := nodeClient.CreateJobRaw(cronJobSpec)
			assert.NoError(t, errCron, "failed to create cron job")
			assert.Empty(t, response.Errors, "failed to create cron job")

			computeJobSpec := `
					type = "standardcapabilities"
					schemaVersion = 1
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
				`

			response, _, errCompute := nodeClient.CreateJobRaw(computeJobSpec)
			assert.NoError(t, errCompute, "failed to create compute job")
			assert.Empty(t, response.Errors, "failed to create compute job")

			consensusJobSpec := fmt.Sprintf(`
					type = "offchainreporting2"
					schemaVersion = 1
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
				ocr3CapabilityAddress,
				nodesInfo[i].OcrKeyBundleID,
				bootstrapNodeInfo.PeerID,
				"node0:5001",
				nodesInfo[i].TransmitterAddress,
				bc.ChainID,
				nodesInfo[i].OcrKeyBundleID,
			)
			fmt.Println("consensusJobSpec", consensusJobSpec)
			response, _, errCons := nodeClient.CreateJobRaw(consensusJobSpec)
			assert.NoError(t, errCons, "failed to create consensus job")
			assert.Empty(t, response.Errors, "failed to create consensus job")
		}()
	}
	wg.Wait()
}

func registerDONAndCapabilities(t *testing.T, capRegAddr common.Address, hashedCapabilities [][32]byte, nodesInfo []NodeInfo, sc *seth.Client) {
	// Register node operators, nodes and DON in the Capabilities registry
	nopsToAdd := make([]cr_wrapper.CapabilitiesRegistryNodeOperator, len(nodesInfo)-1)
	nodesToAdd := make([]cr_wrapper.CapabilitiesRegistryNodeParams, len(nodesInfo)-1)
	donNodes := make([][32]byte, len(nodesInfo)-1)

	for i, node := range nodesInfo {
		// skip the first node, as it's the bootstrap node
		// it doesn't have any capabilities that are required by the workflow
		if i == 0 {
			continue
		}
		nopsToAdd[i-1] = cr_wrapper.CapabilitiesRegistryNodeOperator{
			Admin: common.HexToAddress(node.TransmitterAddress),
			Name:  fmt.Sprintf("NOP %d", i),
		}

		var peerID ragetypes.PeerID
		err := peerID.UnmarshalText([]byte(node.PeerID))
		require.NoError(t, err, "failed to unmarshal peer ID")

		nodesToAdd[i-1] = cr_wrapper.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      uint32(i), //nolint:gosec // disable G115
			Signer:              common.BytesToHash(node.Signer.Bytes()),
			P2pId:               peerID,
			EncryptionPublicKey: [32]byte{1, 2, 3, 4, 5},
			HashedCapabilityIds: hashedCapabilities,
		}

		donNodes[i-1] = peerID
	}

	capabilitiesRegistryInstance, err := cr_wrapper.NewCapabilitiesRegistry(capRegAddr, sc.Client)
	require.NoError(t, err, "failed to create capabilities registry instance")

	// Add NOPs to capabilities registry
	_, decodeErr := sc.Decode(capabilitiesRegistryInstance.AddNodeOperators(
		sc.NewTXOpts(),
		nopsToAdd,
	))
	require.NoError(t, decodeErr, "failed to add NOPs to capabilities registry")

	// Add nodes to capabilities registry
	_, decodeErr = sc.Decode(capabilitiesRegistryInstance.AddNodes(
		sc.NewTXOpts(),
		nodesToAdd,
	))
	require.NoError(t, decodeErr, "failed to add nodes to capabilities registry")

	capRegConfig := make([]cr_wrapper.CapabilitiesRegistryCapabilityConfiguration, len(hashedCapabilities))
	for i, hashed := range hashedCapabilities {
		capRegConfig[i] = cr_wrapper.CapabilitiesRegistryCapabilityConfiguration{
			CapabilityId: hashed,
			Config:       []byte(""),
		}
	}

	// Add nodeset to capabilities registry
	_, decodeErr = sc.Decode(capabilitiesRegistryInstance.AddDON(
		sc.NewTXOpts(),
		donNodes,
		capRegConfig,
		true,     // is public
		true,     // accepts workflows
		uint8(1), // max number of malicious nodes
	))
	require.NoError(t, decodeErr, "failed to add DON to capabilities registry")
}

/*
!!! ATTENTION !!!

Do not use this test as a template for your tests. It's hacky, since we were working under time pressure. We will soon refactor it follow best practices
and a golden example. Apart from its structure what is currently missing is:
- using `chainlink/deployment` to deploy and configure all the contracts
- using Job Distribution to create jobs for the nodes
- using only `chainlink-cli` to register the workflow
- using a mock service to provide the feed data
*/
func TestKeystoneWithOCR3Workflow(t *testing.T) {
	testLogger := logging.GetTestLogger(t)

	// Define and load the test configuration
	donID := uint32(1)
	workflowName := "abcdefgasd"
	feedID := "018bfe8840700040000000000000000000000000000000000000000000000000" // without 0x prefix!
	feedBytes := common.HexToHash(feedID)

	in, err := framework.Load[WorkflowTestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateInputsAndEnvVars(t, in)

	pkey := os.Getenv("PRIVATE_KEY")

	// Create a new blockchain network
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	sc, err := seth.NewClientBuilder().
		WithRpcUrl(bc.Nodes[0].HostWSUrl).
		WithPrivateKeys([]string{pkey}).
		Build()
	require.NoError(t, err, "failed to create seth client")

	// Prepare the chainlink/deployment environment
	ctfEnv, chainSelector := buildChainlinkDeploymentEnv(t, sc)

	// Define required capabilities
	// These need to match the capabilities that are required by the workflow,
	// which in our case is a Proof-of-Reserves workflow
	allCaps := []cr_wrapper.CapabilitiesRegistryCapability{
		{
			LabelledName:   "offchain_reporting",
			Version:        "1.0.0",
			CapabilityType: 2, // CONSENSUS
			ResponseType:   0, // REPORT
		},
		{
			LabelledName:   "write_geth-testnet",
			Version:        "1.0.0",
			CapabilityType: 3, // TARGET
			ResponseType:   1, // OBSERVATION_IDENTICAL
		},
		{
			LabelledName:   "cron-trigger",
			Version:        "1.0.0",
			CapabilityType: uint8(0), // trigger
		},
		{
			LabelledName:   "custom-compute",
			Version:        "1.0.0",
			CapabilityType: uint8(1), // action
		},
	}
	capRegAddr, hashedCapabilities := prepareCapabilitiesRegistry(t, sc, allCaps)

	// Deploy keystone forwarder contract
	forwarderAddress := deployKeystoneForwarder(t, testLogger, ctfEnv, chainSelector)

	// Deploy and pre-configure workflow registry contract
	workflowRegistryAddr := prepareWorkflowRegistry(t, testLogger, ctfEnv, chainSelector, sc, donID)

	// Deploy and configure Keystone Feeds Consumer contract
	feedsConsumerAddress := prepareFeedsConsumer(t, testLogger, ctfEnv, chainSelector, sc, forwarderAddress, workflowName)

	// Register the workflow (either via chainlink-cli or by calling the workflow registry directly)
	registerWorkflow(t, in, sc, capRegAddr, workflowRegistryAddr, feedsConsumerAddress, donID, chainSelector, workflowName, pkey, bc.Nodes[0].HostHTTPUrl)

	// Deploy and fund the DON
	_, nodesInfo := starAndFundNodes(t, in, bc, sc)
	_, nodeClients := configureNodes(t, nodesInfo, in, bc, capRegAddr, workflowRegistryAddr, forwarderAddress)

	// Deploy OCR3 Capability contract
	ocr3CapabilityAddress := deployOCR3Capability(t, testLogger, sc)

	// Create OCR3 and capability jobs for each node
	createNodeJobs(t, nodeClients, nodesInfo, bc, ocr3CapabilityAddress)

	// Register DON and capabilities
	registerDONAndCapabilities(t, capRegAddr, hashedCapabilities, nodesInfo, sc)

	// configure Keystone Forwarder contract
	configureKeystoneForwarder(t, forwarderAddress, sc, nodesInfo)

	// Wait for OCR listeners to be ready before setting the configuration.
	// If the ConfigSet event is missed, OCR protocol will not start.
	// TODO make it fluent!
	testLogger.Info().Msg("Waiting 30s for OCR listeners to be ready...")
	time.Sleep(30 * time.Second)
	testLogger.Info().Msg("Proceeding to set OCR3 configuration.")

	// Configure OCR3 capability contract
	configureOCR3Capability(t, ocr3CapabilityAddress, sc, nodesInfo)

	// It can take a while before the first report is produced, particularly on CI.
	timeout := 10 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	feedsConsumerInstance, err := feeds_consumer.NewKeystoneFeedsConsumer(feedsConsumerAddress, sc.Client)
	require.NoError(t, err, "failed to create feeds consumer instance")

	startTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			t.Fatalf("feed did not update, timeout after %s", timeout)
		case <-time.After(10 * time.Second):
			elapsed := time.Since(startTime).Round(time.Second)
			price, _, err := feedsConsumerInstance.GetPrice(
				sc.NewCallOpts(),
				feedBytes,
			)
			require.NoError(t, err, "failed to get price from Keystone Consumer contract")

			if price.String() != "0" {
				testLogger.Info().Msgf("Feed updated after %s - price set, price=%s", elapsed, price)
				return
			}
			testLogger.Info().Msgf("Feed not updated yet, waiting for %s", elapsed)
		}
	}
}
