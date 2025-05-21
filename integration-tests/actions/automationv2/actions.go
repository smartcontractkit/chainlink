package automationv2

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	ctf_concurrency "github.com/smartcontractkit/chainlink-testing-framework/lib/concurrency"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tt "github.com/smartcontractkit/chainlink/integration-tests/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type NodeDetails struct {
	P2PId                 string
	TransmitterAddresses  []string
	OCR2ConfigPublicKey   string
	OCR2OffchainPublicKey string
	OCR2OnChainPublicKey  string
	OCR2Id                string
}

type AutomationTest struct {
	ChainClient *seth.Client

	TestConfig tt.AutomationTestConfig

	LinkToken   contracts.LinkToken
	Transcoder  contracts.UpkeepTranscoder
	LINKETHFeed contracts.MockLINKETHFeed
	ETHUSDFeed  contracts.MockETHUSDFeed
	LINKUSDFeed contracts.MockETHUSDFeed
	WETHToken   contracts.WETHToken
	GasFeed     contracts.MockGasFeed
	Registry    contracts.KeeperRegistry
	Registrar   contracts.KeeperRegistrar

	RegistrySettings       contracts.KeeperRegistrySettings
	RegistrarSettings      contracts.KeeperRegistrarSettings
	PluginConfig           ocr2keepers30config.OffchainConfig
	PublicConfig           ocr3.PublicConfig
	UpkeepPrivilegeManager common.Address
	UpkeepIDs              []*big.Int

	IsOnk8s bool

	ChainlinkNodesk8s []*nodeclient.ChainlinkK8sClient
	ChainlinkNodes    []*nodeclient.ChainlinkClient

	DockerEnv *test_env.CLClusterTestEnv

	NodeDetails              []NodeDetails
	DefaultP2Pv2Bootstrapper string
	mercuryCredentialName    string
	TransmitterKeyIndex      int

	Logger zerolog.Logger
}

type UpkeepConfig struct {
	UpkeepName     string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	TriggerType    uint8
	CheckData      []byte
	TriggerConfig  []byte
	OffchainConfig []byte
	FundingAmount  *big.Int
}

func NewAutomationTestK8s(
	l zerolog.Logger,
	chainClient *seth.Client,
	chainlinkNodes []*nodeclient.ChainlinkK8sClient,
	config tt.AutomationTestConfig,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		TestConfig:             config,
		ChainlinkNodesk8s:      chainlinkNodes,
		IsOnk8s:                true,
		TransmitterKeyIndex:    0,
		UpkeepPrivilegeManager: chainClient.MustGetRootKeyAddress(),
		mercuryCredentialName:  "",
		Logger:                 l,
	}
}

func NewAutomationTestDocker(
	l zerolog.Logger,
	chainClient *seth.Client,
	chainlinkNodes []*nodeclient.ChainlinkClient,
	config tt.AutomationTestConfig,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		TestConfig:             config,
		ChainlinkNodes:         chainlinkNodes,
		IsOnk8s:                false,
		TransmitterKeyIndex:    0,
		UpkeepPrivilegeManager: chainClient.MustGetRootKeyAddress(),
		mercuryCredentialName:  "",
		Logger:                 l,
	}
}

func (a *AutomationTest) SetIsOnk8s(flag bool) {
	a.IsOnk8s = flag
}

func (a *AutomationTest) SetMercuryCredentialName(name string) {
	a.mercuryCredentialName = name
}

func (a *AutomationTest) SetTransmitterKeyIndex(index int) {
	a.TransmitterKeyIndex = index
}

func (a *AutomationTest) SetUpkeepPrivilegeManager(address string) {
	a.UpkeepPrivilegeManager = common.HexToAddress(address)
}

func (a *AutomationTest) SetDockerEnv(env *test_env.CLClusterTestEnv) {
	a.DockerEnv = env
}

func (a *AutomationTest) DeployLINK() error {
	linkToken, err := contracts.DeployLinkTokenContract(a.Logger, a.ChainClient)
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	return nil
}

func (a *AutomationTest) LoadLINK(address string) error {
	linkToken, err := contracts.LoadLinkTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	a.Logger.Info().Str("LINK Token Address", a.LinkToken.Address()).Msg("Successfully loaded LINK Token")
	return nil
}

func (a *AutomationTest) DeployTranscoder() error {
	transcoder, err := contracts.DeployUpkeepTranscoder(a.ChainClient)
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	return nil
}

func (a *AutomationTest) LoadTranscoder(address string) error {
	transcoder, err := contracts.LoadUpkeepTranscoder(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	a.Logger.Info().Str("Transcoder Address", a.Transcoder.Address()).Msg("Successfully loaded Transcoder")
	return nil
}

func (a *AutomationTest) DeployLinkEthFeed() error {
	ethLinkFeed, err := contracts.DeployMockLINKETHFeed(a.ChainClient, a.RegistrySettings.FallbackLinkPrice)
	if err != nil {
		return err
	}
	a.LINKETHFeed = ethLinkFeed
	return nil
}

func (a *AutomationTest) LoadLinkEthFeed(address string) error {
	ethLinkFeed, err := contracts.LoadMockLINKETHFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKETHFeed = ethLinkFeed
	a.Logger.Info().Str("LINK/ETH Feed Address", a.LINKETHFeed.Address()).Msg("Successfully loaded LINK/ETH Feed")
	return nil
}

func (a *AutomationTest) DeployEthUSDFeed() error {
	ethUSDFeed, err := contracts.DeployMockETHUSDFeed(a.ChainClient, a.RegistrySettings.FallbackLinkPrice)
	if err != nil {
		return err
	}
	a.ETHUSDFeed = ethUSDFeed
	return nil
}

func (a *AutomationTest) LoadEthUSDFeed(address string) error {
	ethUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.ETHUSDFeed = ethUSDFeed
	a.Logger.Info().Str("ETH/USD Feed Address", a.ETHUSDFeed.Address()).Msg("Successfully loaded ETH/USD Feed")
	return nil
}

func (a *AutomationTest) DeployLinkUSDFeed() error {
	linkUSDFeed, err := contracts.DeployMockETHUSDFeed(a.ChainClient, a.RegistrySettings.FallbackLinkPrice)
	if err != nil {
		return err
	}
	a.LINKUSDFeed = linkUSDFeed
	return nil
}

func (a *AutomationTest) LoadLinkUSDFeed(address string) error {
	linkUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKUSDFeed = linkUSDFeed
	a.Logger.Info().Str("LINK/USD Feed Address", a.LINKUSDFeed.Address()).Msg("Successfully loaded LINK/USD Feed")
	return nil
}

func (a *AutomationTest) DeployWETH() error {
	wethToken, err := contracts.DeployWETHTokenContract(a.Logger, a.ChainClient)
	if err != nil {
		return err
	}
	a.WETHToken = wethToken
	return nil
}

func (a *AutomationTest) LoadWETH(address string) error {
	wethToken, err := contracts.LoadWETHTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.WETHToken = wethToken
	a.Logger.Info().Str("WETH Token Address", a.WETHToken.Address()).Msg("Successfully loaded WETH Token")
	return nil
}

func (a *AutomationTest) DeployGasFeed() error {
	gasFeed, err := contracts.DeployMockGASFeed(a.ChainClient, a.RegistrySettings.FallbackGasPrice)
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	return nil
}

func (a *AutomationTest) LoadEthGasFeed(address string) error {
	gasFeed, err := contracts.LoadMockGASFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	a.Logger.Info().Str("Gas Feed Address", a.GasFeed.Address()).Msg("Successfully loaded Gas Feed")
	return nil
}

func (a *AutomationTest) DeployRegistry() error {
	registryOpts := &contracts.KeeperRegistryOpts{
		RegistryVersion:   a.RegistrySettings.RegistryVersion,
		LinkAddr:          a.LinkToken.Address(),
		ETHFeedAddr:       a.LINKETHFeed.Address(),
		GasFeedAddr:       a.GasFeed.Address(),
		TranscoderAddr:    a.Transcoder.Address(),
		RegistrarAddr:     utils.ZeroAddress.Hex(),
		Settings:          a.RegistrySettings,
		LinkUSDFeedAddr:   a.ETHUSDFeed.Address(),
		NativeUSDFeedAddr: a.LINKUSDFeed.Address(),
		WrappedNativeAddr: a.WETHToken.Address(),
	}
	registry, err := contracts.DeployKeeperRegistry(a.ChainClient, registryOpts)
	if err != nil {
		return err
	}
	a.Registry = registry
	return nil
}

func (a *AutomationTest) LoadRegistry(registryAddress, chainModuleAddress string) error {
	registry, err := contracts.LoadKeeperRegistry(a.Logger, a.ChainClient, common.HexToAddress(registryAddress), a.RegistrySettings.RegistryVersion, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return err
	}
	a.Registry = registry
	a.Logger.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", a.Registry.Address()).Msg("Successfully loaded Registry")
	return nil
}

func (a *AutomationTest) DeployRegistrar() error {
	if a.Registry == nil {
		return fmt.Errorf("registry must be deployed or loaded before registrar")
	}
	a.RegistrarSettings.RegistryAddr = a.Registry.Address()
	a.RegistrarSettings.WETHTokenAddr = a.WETHToken.Address()
	registrar, err := contracts.DeployKeeperRegistrar(a.ChainClient, a.RegistrySettings.RegistryVersion, a.LinkToken.Address(), a.RegistrarSettings)
	if err != nil {
		return err
	}
	a.Registrar = registrar
	return nil
}

func (a *AutomationTest) LoadRegistrar(address string) error {
	if a.Registry == nil {
		return fmt.Errorf("registry must be deployed or loaded before registrar")
	}
	a.RegistrarSettings.RegistryAddr = a.Registry.Address()
	registrar, err := contracts.LoadKeeperRegistrar(a.ChainClient, common.HexToAddress(address), a.RegistrySettings.RegistryVersion)
	if err != nil {
		return err
	}
	a.Logger.Info().Str("Registrar Address", registrar.Address()).Msg("Successfully loaded Registrar")
	a.Registrar = registrar
	return nil
}

func (a *AutomationTest) CollectNodeDetails() error {
	var (
		nodes []*nodeclient.ChainlinkClient
	)
	if a.IsOnk8s {
		for _, node := range a.ChainlinkNodesk8s[:] {
			nodes = append(nodes, node.ChainlinkClient)
		}
		a.ChainlinkNodes = nodes
	} else {
		nodes = a.ChainlinkNodes[:]
	}

	nodeDetails := make([]NodeDetails, 0)

	for i, node := range nodes {
		nodeDetail := NodeDetails{}
		P2PIds, err := node.MustReadP2PKeys()
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to read P2P keys from node %d", i))
		}
		nodeDetail.P2PId = P2PIds.Data[0].Attributes.PeerID

		OCR2Keys, err := node.MustReadOCR2Keys()
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to read OCR2 keys from node %d", i))
		}
		for _, key := range OCR2Keys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				nodeDetail.OCR2ConfigPublicKey = key.Attributes.ConfigPublicKey
				nodeDetail.OCR2OffchainPublicKey = key.Attributes.OffChainPublicKey
				nodeDetail.OCR2OnChainPublicKey = key.Attributes.OnChainPublicKey
				nodeDetail.OCR2Id = key.ID
				break
			}
		}

		TransmitterKeys, err := node.EthAddressesForChain(fmt.Sprint(a.ChainClient.ChainID))
		nodeDetail.TransmitterAddresses = make([]string, 0)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to read Transmitter keys from node %d", i))
		}
		nodeDetail.TransmitterAddresses = append(nodeDetail.TransmitterAddresses, TransmitterKeys...)
		nodeDetails = append(nodeDetails, nodeDetail)
	}
	a.NodeDetails = nodeDetails

	if a.IsOnk8s {
		a.DefaultP2Pv2Bootstrapper = fmt.Sprintf("%s@%s-node-1:%d", a.NodeDetails[0].P2PId, a.ChainlinkNodesk8s[0].Name(), 6690)
	} else {
		a.DefaultP2Pv2Bootstrapper = fmt.Sprintf("%s@%s:%d", a.NodeDetails[0].P2PId, a.ChainlinkNodes[0].InternalIP(), 6690)
	}
	return nil
}

func (a *AutomationTest) AddBootstrapJob() error {
	bootstrapSpec := &nodeclient.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node " + a.Registry.Address(),
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: a.Registry.Address(),
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(a.ChainClient.ChainID),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
	_, err := a.ChainlinkNodes[0].MustCreateJob(bootstrapSpec)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to create bootstrap job on bootstrap node"))
	}
	return nil
}

func (a *AutomationTest) AddAutomationJobs() error {
	var contractVersion string
	if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_2 || a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_3 {
		contractVersion = "v2.1+"
	} else if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_1 {
		contractVersion = "v2.1"
	} else if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
		contractVersion = "v2.0"
	} else {
		return fmt.Errorf("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
	}
	pluginCfg := map[string]interface{}{
		"contractVersion": "\"" + contractVersion + "\"",
	}
	if strings.Contains(contractVersion, "v2.1") {
		if a.mercuryCredentialName != "" {
			pluginCfg["mercuryCredentialName"] = "\"" + a.mercuryCredentialName + "\""
		}
	}
	for i := 1; i < len(a.ChainlinkNodes); i++ {
		autoOCR2JobSpec := nodeclient.OCR2TaskJobSpec{
			Name:    "automation-" + contractVersion + "-" + a.Registry.Address(),
			JobType: "offchainreporting2",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "ocr2automation",
				ContractID: a.Registry.Address(),
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": int(a.ChainClient.ChainID),
				},
				PluginConfig:                      pluginCfg,
				ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
				TransmitterID:                     null.StringFrom(a.NodeDetails[i].TransmitterAddresses[a.TransmitterKeyIndex]),
				P2PV2Bootstrappers:                pq.StringArray{a.DefaultP2Pv2Bootstrapper},
				OCRKeyBundleID:                    null.StringFrom(a.NodeDetails[i].OCR2Id),
			},
		}
		_, err := a.ChainlinkNodes[i].MustCreateJob(&autoOCR2JobSpec)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to create OCR2 job on node %d", i+1))
		}
	}
	return nil
}

func (a *AutomationTest) SetConfigOnRegistry() error {
	donNodes := a.NodeDetails[1:]
	S := make([]int, len(donNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(donNodes))
	var signerOnchainPublicKeys []types.OnchainPublicKey
	var transmitterAccounts []types.Account
	var f uint8
	var offchainConfigVersion uint64
	var offchainConfig []byte
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(donNodes))
	eg := &errgroup.Group{}
	for i, donNode := range donNodes {
		index, chainlinkNode := i, donNode
		eg.Go(func() error {
			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2OffchainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            chainlinkNode.P2PId,
					TransmitAccount:   types.Account(chainlinkNode.TransmitterAddresses[a.TransmitterKeyIndex]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to build oracle identities"))
	}

	switch a.RegistrySettings.RegistryVersion {
	case ethereum.RegistryVersion_2_0:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR2ConfigArgs(a, S, oracleIdentities)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to build config args"))
		}
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR3ConfigArgs(a, S, oracleIdentities)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to build config args"))
		}
	default:
		return fmt.Errorf("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
	}

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		if len(signer) != 20 {
			return fmt.Errorf("OnChainPublicKey '%v' has wrong length for address", signer)
		}
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		if !common.IsHexAddress(string(transmitter)) {
			return fmt.Errorf("TransmitAccount '%s' is not a valid Ethereum address", string(transmitter))
		}
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	ocrConfig := contracts.OCRv2Config{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
		ocrConfig.OnchainConfig = a.RegistrySettings.Encode20OnchainConfig(a.Registrar.Address())
		err = a.Registry.SetConfig(a.RegistrySettings, ocrConfig)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to set config on registry"))
		}
	} else {
		if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_1 {
			ocrConfig.TypedOnchainConfig21 = a.RegistrySettings.Create21OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager)
		} else if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_2 {
			ocrConfig.TypedOnchainConfig22 = a.RegistrySettings.Create22OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager, a.Registry.ChainModuleAddress(), a.Registry.ReorgProtectionEnabled())
		} else if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_3 {
			ocrConfig.TypedOnchainConfig23 = a.RegistrySettings.Create23OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager, a.Registry.ChainModuleAddress(), a.Registry.ReorgProtectionEnabled())
			ocrConfig.BillingTokens = []common.Address{
				common.HexToAddress(a.LinkToken.Address()),
				common.HexToAddress(a.WETHToken.Address()),
			}

			ocrConfig.BillingConfigs = []i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23BillingConfig{
				{
					GasFeePPB:         100,
					FlatFeeMilliCents: big.NewInt(500),
					PriceFeed:         common.HexToAddress(a.ETHUSDFeed.Address()),
					Decimals:          18,
					FallbackPrice:     big.NewInt(1000),
					MinSpend:          big.NewInt(200),
				},
				{
					GasFeePPB:         100,
					FlatFeeMilliCents: big.NewInt(500),
					PriceFeed:         common.HexToAddress(a.LINKUSDFeed.Address()),
					Decimals:          18,
					FallbackPrice:     big.NewInt(1000),
					MinSpend:          big.NewInt(200),
				},
			}
		}
		a.Logger.Debug().Interface("ocrConfig", ocrConfig).Msg("Setting OCR3 config")
		err = a.Registry.SetConfigTypeSafe(ocrConfig)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to set config on registry"))
		}
	}
	return nil
}

func calculateOCR2ConfigArgs(a *AutomationTest, S []int, oracleIdentities []confighelper.OracleIdentityExtra) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	offC, _ := json.Marshal(ocr2keepers20config.OffchainConfig{
		TargetProbability:    a.PluginConfig.TargetProbability,
		TargetInRounds:       a.PluginConfig.TargetInRounds,
		PerformLockoutWindow: a.PluginConfig.PerformLockoutWindow,
		GasLimitPerReport:    a.PluginConfig.GasLimitPerReport,
		GasOverheadPerUpkeep: a.PluginConfig.GasOverheadPerUpkeep,
		MinConfirmations:     a.PluginConfig.MinConfirmations,
		MaxUpkeepBatchSize:   a.PluginConfig.MaxUpkeepBatchSize,
	})

	rMax := a.PublicConfig.RMax
	if rMax > math.MaxUint8 {
		panic(fmt.Errorf("rmax overflows uint8: %d", rMax))
	}

	return ocr2.ContractSetConfigArgsForTests(
		a.PublicConfig.DeltaProgress, a.PublicConfig.DeltaResend,
		a.PublicConfig.DeltaRound, a.PublicConfig.DeltaGrace,
		a.PublicConfig.DeltaStage, uint8(rMax),
		S, oracleIdentities, offC,
		nil,
		a.PublicConfig.MaxDurationQuery, a.PublicConfig.MaxDurationObservation,
		1200*time.Millisecond,
		a.PublicConfig.MaxDurationShouldAcceptAttestedReport,
		a.PublicConfig.MaxDurationShouldTransmitAcceptedReport,
		a.PublicConfig.F, a.PublicConfig.OnchainConfig,
	)
}

func calculateOCR3ConfigArgs(a *AutomationTest, S []int, oracleIdentities []confighelper.OracleIdentityExtra) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	offC, _ := json.Marshal(a.PluginConfig)

	return ocr3.ContractSetConfigArgsForTests(
		a.PublicConfig.DeltaProgress, a.PublicConfig.DeltaResend, a.PublicConfig.DeltaInitial,
		a.PublicConfig.DeltaRound, a.PublicConfig.DeltaGrace, a.PublicConfig.DeltaCertifiedCommitRequest,
		a.PublicConfig.DeltaStage, a.PublicConfig.RMax,
		S, oracleIdentities, offC,
		nil, a.PublicConfig.MaxDurationQuery, a.PublicConfig.MaxDurationObservation,
		a.PublicConfig.MaxDurationShouldAcceptAttestedReport,
		a.PublicConfig.MaxDurationShouldTransmitAcceptedReport,
		a.PublicConfig.F, a.PublicConfig.OnchainConfig,
	)
}

type registrationResult struct {
	txHash common.Hash
}

func (r registrationResult) GetResult() common.Hash {
	return r.txHash
}

func (a *AutomationTest) RegisterUpkeeps(upkeepConfigs []UpkeepConfig, maxConcurrency int) ([]common.Hash, error) {
	concurrency, err := actions.GetAndAssertCorrectConcurrency(a.ChainClient, 1)
	if err != nil {
		return nil, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		a.Logger.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	var registerUpkeep = func(resultCh chan registrationResult, errorCh chan error, executorNum int, upkeepConfig UpkeepConfig) {
		keyNum := executorNum + 1 // key 0 is the root key
		var registrationRequest []byte
		var registrarABI *abi.ABI
		var err error
		switch a.RegistrySettings.RegistryVersion {
		case ethereum.RegistryVersion_2_0:
			registrarABI, err = keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.GetAbi()
			if err != nil {
				errorCh <- errors.Join(err, fmt.Errorf("failed to get registrar abi"))
				return
			}
			registrationRequest, err = registrarABI.Pack(
				"register",
				upkeepConfig.UpkeepName,
				upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract,
				upkeepConfig.GasLimit,
				upkeepConfig.AdminAddress,
				upkeepConfig.CheckData,
				upkeepConfig.OffchainConfig,
				upkeepConfig.FundingAmount,
				a.ChainClient.Addresses[keyNum])
			if err != nil {
				errorCh <- errors.Join(err, fmt.Errorf("failed to pack registrar request"))
				return
			}
		case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2: // 2.1 and 2.2 use the same registrar
			registrarABI, err = automation_registrar_wrapper2_1.AutomationRegistrarMetaData.GetAbi()
			if err != nil {
				errorCh <- errors.Join(err, fmt.Errorf("failed to get registrar abi"))
				return
			}
			registrationRequest, err = registrarABI.Pack(
				"register",
				upkeepConfig.UpkeepName,
				upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract,
				upkeepConfig.GasLimit,
				upkeepConfig.AdminAddress,
				upkeepConfig.TriggerType,
				upkeepConfig.CheckData,
				upkeepConfig.TriggerConfig,
				upkeepConfig.OffchainConfig,
				upkeepConfig.FundingAmount,
				a.ChainClient.Addresses[keyNum])
			if err != nil {
				errorCh <- errors.Join(err, fmt.Errorf("failed to pack registrar request"))
				return
			}
		default:
			errorCh <- fmt.Errorf("v2.0, v2.1, and v2.2 are the only supported versions")
			return
		}

		tx, err := a.LinkToken.TransferAndCallFromKey(a.Registrar.Address(), upkeepConfig.FundingAmount, registrationRequest, keyNum)
		if err != nil {
			errorCh <- errors.Join(err, fmt.Errorf("client number %d failed to register upkeep %s", keyNum, upkeepConfig.UpkeepContract.Hex()))
			return
		}

		resultCh <- registrationResult{txHash: tx.Hash()}
	}

	executor := ctf_concurrency.NewConcurrentExecutor[common.Hash, registrationResult, UpkeepConfig](a.Logger)
	results, err := executor.Execute(concurrency, upkeepConfigs, registerUpkeep)
	if err != nil {
		return nil, err
	}

	if len(results) != len(upkeepConfigs) {
		return nil, fmt.Errorf("failed to register all upkeeps. Expected %d, got %d", len(upkeepConfigs), len(results))
	}

	a.Logger.Info().Msg("Successfully registered all upkeeps")

	return results, nil
}

type UpkeepId = *big.Int

type confirmationResult struct {
	upkeepID UpkeepId
}

func (c confirmationResult) GetResult() UpkeepId {
	return c.upkeepID
}

func (a *AutomationTest) ConfirmUpkeepsRegistered(registrationTxHashes []common.Hash, maxConcurrency int) ([]*big.Int, error) {
	concurrency, err := actions.GetAndAssertCorrectConcurrency(a.ChainClient, 1)
	if err != nil {
		return nil, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		a.Logger.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	var confirmUpkeep = func(resultCh chan confirmationResult, errorCh chan error, _ int, txHash common.Hash) {
		receipt, err := a.ChainClient.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			errorCh <- errors.Join(err, fmt.Errorf("failed to confirm upkeep registration"))
			return
		}

		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := a.Registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		if upkeepId == nil {
			errorCh <- fmt.Errorf("failed to parse upkeep id from registration receipt")
			return
		}
		resultCh <- confirmationResult{upkeepID: upkeepId}
	}

	executor := ctf_concurrency.NewConcurrentExecutor[UpkeepId, confirmationResult, common.Hash](a.Logger)
	results, err := executor.Execute(concurrency, registrationTxHashes, confirmUpkeep)

	if err != nil {
		return nil, fmt.Errorf("failed confirmations: %d | successful confirmations: %d", len(executor.GetErrors()), len(results))
	}

	if len(registrationTxHashes) != len(results) {
		return nil, fmt.Errorf("failed to confirm all upkeeps. Expected %d, got %d", len(registrationTxHashes), len(results))
	}

	seen := make(map[*big.Int]bool)
	for _, upkeepId := range results {
		if seen[upkeepId] {
			return nil, fmt.Errorf("duplicate upkeep id: %s. Something went wrong during upkeep confirmation. Please check the test code", upkeepId.String())
		}
		seen[upkeepId] = true
	}

	a.Logger.Info().Msg("Successfully confirmed all upkeeps")
	a.UpkeepIDs = results

	return results, nil
}

func (a *AutomationTest) AddJobsAndSetConfig(t *testing.T) {
	l := logging.GetTestLogger(t)
	err := a.AddBootstrapJob()
	require.NoError(t, err, "Error adding bootstrap job")
	err = a.AddAutomationJobs()
	require.NoError(t, err, "Error adding automation jobs")

	l.Info().
		Interface("Plugin Config", a.PluginConfig).
		Interface("Public Config", a.PublicConfig).
		Interface("Registry Settings", a.RegistrySettings).
		Interface("Registrar Settings", a.RegistrarSettings).
		Msg("Configuring registry")
	err = a.SetConfigOnRegistry()
	require.NoError(t, err, "Error setting config on registry")
	l.Info().Str("Registry Address", a.Registry.Address()).Msg("Successfully setConfig on registry")
}

func (a *AutomationTest) SetupAutomationDeployment(t *testing.T) {
	a.setupDeployment(t, true)
}

func (a *AutomationTest) SetupAutomationDeploymentWithoutJobs(t *testing.T) {
	a.setupDeployment(t, false)
}

func (a *AutomationTest) setupDeployment(t *testing.T, addJobs bool) {
	l := logging.GetTestLogger(t)
	err := a.CollectNodeDetails()
	require.NoError(t, err, "Error collecting node details")
	l.Info().Msg("Collected Node Details")
	l.Debug().Interface("Node Details", a.NodeDetails).Msg("Node Details")

	if a.TestConfig.GetAutomationConfig().UseExistingLinkTokenContract() {
		linkAddress, err := a.TestConfig.GetAutomationConfig().LinkTokenContractAddress()
		require.NoError(t, err, "Error getting link token contract address")
		err = a.LoadLINK(linkAddress.String())
		require.NoError(t, err, "Error loading link token contract")
	} else {
		err = a.DeployLINK()
		require.NoError(t, err, "Error deploying link token contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingWethContract() {
		wethAddress, err := a.TestConfig.GetAutomationConfig().WethContractAddress()
		require.NoError(t, err, "Error getting weth token contract address")
		err = a.LoadWETH(wethAddress.String())
		require.NoError(t, err, "Error loading weth token contract")
	} else {
		err = a.DeployWETH()
		require.NoError(t, err, "Error deploying weth token contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingLinkEthFeedContract() {
		linkEthFeedAddress, err := a.TestConfig.GetAutomationConfig().LinkEthFeedContractAddress()
		require.NoError(t, err, "Error getting link eth feed contract address")
		err = a.LoadLinkEthFeed(linkEthFeedAddress.String())
		require.NoError(t, err, "Error loading link eth feed contract")
	} else {
		err = a.DeployLinkEthFeed()
		require.NoError(t, err, "Error deploying link eth feed contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingEthGasFeedContract() {
		gasFeedAddress, err := a.TestConfig.GetAutomationConfig().EthGasFeedContractAddress()
		require.NoError(t, err, "Error getting gas feed contract address")
		err = a.LoadEthGasFeed(gasFeedAddress.String())
		require.NoError(t, err, "Error loading gas feed contract")
	} else {
		err = a.DeployGasFeed()
		require.NoError(t, err, "Error deploying gas feed contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingEthUSDFeedContract() {
		ethUsdFeedAddress, err := a.TestConfig.GetAutomationConfig().EthUSDFeedContractAddress()
		require.NoError(t, err, "Error getting eth usd feed contract address")
		err = a.LoadEthUSDFeed(ethUsdFeedAddress.String())
		require.NoError(t, err, "Error loading eth usd feed contract")
	} else {
		err = a.DeployEthUSDFeed()
		require.NoError(t, err, "Error deploying eth usd feed contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingLinkUSDFeedContract() {
		linkUsdFeedAddress, err := a.TestConfig.GetAutomationConfig().LinkUSDFeedContractAddress()
		require.NoError(t, err, "Error getting link usd feed contract address")
		err = a.LoadLinkUSDFeed(linkUsdFeedAddress.String())
		require.NoError(t, err, "Error loading link usd feed contract")
	} else {
		err = a.DeployLinkUSDFeed()
		require.NoError(t, err, "Error deploying link usd feed contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingTranscoderContract() {
		transcoderAddress, err := a.TestConfig.GetAutomationConfig().TranscoderContractAddress()
		require.NoError(t, err, "Error getting transcoder contract address")
		err = a.LoadTranscoder(transcoderAddress.String())
		require.NoError(t, err, "Error loading transcoder contract")
	} else {
		err = a.DeployTranscoder()
		require.NoError(t, err, "Error deploying transcoder contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingRegistryContract() {
		chainModuleAddress, err := a.TestConfig.GetAutomationConfig().ChainModuleContractAddress()
		require.NoError(t, err, "Error getting chain module contract address")
		registryAddress, err := a.TestConfig.GetAutomationConfig().RegistryContractAddress()
		require.NoError(t, err, "Error getting registry contract address")
		err = a.LoadRegistry(registryAddress.String(), chainModuleAddress.String())
		require.NoError(t, err, "Error loading registry contract")
		if a.Registry.RegistryOwnerAddress().String() != a.ChainClient.MustGetRootKeyAddress().String() {
			l.Debug().Str("RootKeyAddress", a.ChainClient.MustGetRootKeyAddress().String()).Str("Registry Owner Address", a.Registry.RegistryOwnerAddress().String()).Msg("Registry owner address is not the root key address")
			t.Error("Registry owner address is not the root key address")
			t.FailNow()
		}
	} else {
		err = a.DeployRegistry()
		require.NoError(t, err, "Error deploying registry contract")
	}

	if a.TestConfig.GetAutomationConfig().UseExistingRegistrarContract() {
		registrarAddress, err := a.TestConfig.GetAutomationConfig().RegistrarContractAddress()
		require.NoError(t, err, "Error getting registrar contract address")
		err = a.LoadRegistrar(registrarAddress.String())
		require.NoError(t, err, "Error loading registrar contract")
	} else {
		err = a.DeployRegistrar()
		require.NoError(t, err, "Error deploying registrar contract")
	}

	if addJobs {
		a.AddJobsAndSetConfig(t)
	}
}
