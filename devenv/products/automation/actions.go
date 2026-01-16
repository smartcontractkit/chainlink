package automation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	pkg_errors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"gopkg.in/guregu/null.v4"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	ctf_concurrency "github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	devenv_ocr2 "github.com/smartcontractkit/chainlink/devenv/products/ocr2"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

type NodeDetail struct {
	P2PId                 string
	TransmitterAddresses  []string
	OCR2ConfigPublicKey   string
	OCR2OffchainPublicKey string
	OCR2OnChainPublicKey  string
	OCR2Id                string
}

type AutomationTest struct {
	ChainClient *seth.Client

	Config Automation

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

	// ChainlinkNodesk8s []*nodeclient.ChainlinkK8sClient
	ChainlinkNodes []*clclient.ChainlinkClient

	NodeDetails              []NodeDetail
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
	// chainlinkNodes []*nodeclient.ChainlinkK8sClient,
	config Automation,
) *AutomationTest {
	return &AutomationTest{
		ChainClient: chainClient,
		Config:      config,
		// ChainlinkNodesk8s:      chainlinkNodes,
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
	chainlinkNodes []*clclient.ChainlinkClient,
	config Automation,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		Config:                 config,
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
		return errors.New("registry must be deployed or loaded before registrar")
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
		return errors.New("registry must be deployed or loaded before registrar")
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
		nodes []*clclient.ChainlinkClient
	)
	if a.IsOnk8s {
		// 	for _, node := range a.ChainlinkNodesk8s {
		// 		nodes = append(nodes, node.ChainlinkClient)
		// }
		// 	a.ChainlinkNodes = nodes
	} else {
		nodes = a.ChainlinkNodes
	}

	nodeDetails := make([]NodeDetail, 0)

	for i, node := range nodes {
		nodeDetail := NodeDetail{}
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
			if strings.EqualFold(key.Attributes.ChainType, "evm") {
				nodeDetail.OCR2ConfigPublicKey = key.Attributes.ConfigPublicKey
				nodeDetail.OCR2OffchainPublicKey = key.Attributes.OffChainPublicKey
				nodeDetail.OCR2OnChainPublicKey = key.Attributes.OnChainPublicKey
				nodeDetail.OCR2Id = key.ID
				break
			}
		}

		TransmitterKeys, err := node.EthAddressesForChain(strconv.FormatInt(a.ChainClient.ChainID, 10))
		nodeDetail.TransmitterAddresses = make([]string, 0)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to read Transmitter keys from node %d", i))
		}
		nodeDetail.TransmitterAddresses = append(nodeDetail.TransmitterAddresses, TransmitterKeys...)
		nodeDetails = append(nodeDetails, nodeDetail)
	}
	a.NodeDetails = nodeDetails

	if a.IsOnk8s {
		// a.DefaultP2Pv2Bootstrapper = fmt.Sprintf("%s@%s-node-1:%d", a.NodeDetails[0].P2PId, a.ChainlinkNodesk8s[0].Name(), 6690)
	} else {
		a.DefaultP2Pv2Bootstrapper = fmt.Sprintf("%s@%s:%d", a.NodeDetails[0].P2PId, a.ChainlinkNodes[0].InternalIP(), 6690)
	}
	return nil
}

func ChainlinkNodeAddressesLocal(nodes []*clclient.ChainlinkClient) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// ChainlinkNodeAddressesAtIndex will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddressesAtIndex(nodes []*clclient.ChainlinkClient, keyIndex int) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		nodeAddresses, err := node.EthAddresses()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(nodeAddresses[keyIndex]))
	}
	return addresses, nil
}

func (a *AutomationTest) AddBootstrapJob() error {
	bootstrapSpec := &devenv_ocr2.TaskJobSpec{
		Name:    "ocr2 bootstrap node " + a.Registry.Address(),
		JobType: "bootstrap",
		OCR2OracleSpec: devenv_ocr2.OracleSpec{
			ContractID: a.Registry.Address(),
			Relay:      "evm",
			RelayConfig: map[string]any{
				"chainID": int(a.ChainClient.ChainID),
			},
			ContractConfigTrackerPollInterval: *devenv_ocr2.NewInterval(time.Second * 15),
		},
	}
	_, err := a.ChainlinkNodes[0].MustCreateJob(bootstrapSpec)
	if err != nil {
		return errors.Join(err, errors.New("failed to create bootstrap job on bootstrap node"))
	}
	return nil
}

func (a *AutomationTest) AddAutomationJobs() error {
	var contractVersion string
	switch a.RegistrySettings.RegistryVersion {
	case ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		contractVersion = "v2.1+"
	case ethereum.RegistryVersion_2_1:
		contractVersion = "v2.1"
	case ethereum.RegistryVersion_2_0:
		contractVersion = "v2.0"
	default:
		return errors.New("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
	}
	pluginCfg := map[string]any{
		"contractVersion": "\"" + contractVersion + "\"",
	}
	if strings.Contains(contractVersion, "v2.1") {
		if a.mercuryCredentialName != "" {
			pluginCfg["mercuryCredentialName"] = "\"" + a.mercuryCredentialName + "\""
		}
	}
	for i := 1; i < len(a.ChainlinkNodes); i++ {
		autoOCR2JobSpec := devenv_ocr2.TaskJobSpec{
			Name:    "automation-" + contractVersion + "-" + a.Registry.Address(),
			JobType: "offchainreporting2",
			OCR2OracleSpec: devenv_ocr2.OracleSpec{
				PluginType: "ocr2automation",
				ContractID: a.Registry.Address(),
				Relay:      "evm",
				RelayConfig: map[string]any{
					"chainID": int(a.ChainClient.ChainID),
				},
				PluginConfig:                      pluginCfg,
				ContractConfigTrackerPollInterval: *devenv_ocr2.NewInterval(time.Second * 15),
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

// func (a *AutomationTest) SetConfigOnRegistry() error {
// 	donNodes := a.NodeDetails[1:]
// 	S := make([]int, len(donNodes))
// 	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(donNodes))
// 	var signerOnchainPublicKeys []types.OnchainPublicKey
// 	var transmitterAccounts []types.Account
// 	var f uint8
// 	var offchainConfigVersion uint64
// 	var offchainConfig []byte
// 	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(donNodes))
// 	eg := &errgroup.Group{}
// 	for i, donNode := range donNodes {
// 		index, chainlinkNode := i, donNode
// 		eg.Go(func() error {
// 			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2OffchainPublicKey, "ocr2off_evm_"))
// 			if err != nil {
// 				return err
// 			}

// 			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
// 			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
// 			if n != ed25519.PublicKeySize {
// 				return errors.New("wrong number of elements copied")
// 			}

// 			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2ConfigPublicKey, "ocr2cfg_evm_"))
// 			if err != nil {
// 				return err
// 			}

// 			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
// 			n = copy(configPkBytesFixed[:], configPkBytes)
// 			if n != ed25519.PublicKeySize {
// 				return errors.New("wrong number of elements copied")
// 			}

// 			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2OnChainPublicKey, "ocr2on_evm_"))
// 			if err != nil {
// 				return err
// 			}

// 			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
// 			oracleIdentities[index] = confighelper.OracleIdentityExtra{
// 				OracleIdentity: confighelper.OracleIdentity{
// 					OnchainPublicKey:  onchainPkBytes,
// 					OffchainPublicKey: offchainPkBytesFixed,
// 					PeerID:            chainlinkNode.P2PId,
// 					TransmitAccount:   types.Account(chainlinkNode.TransmitterAddresses[a.TransmitterKeyIndex]),
// 				},
// 				ConfigEncryptionPublicKey: configPkBytesFixed,
// 			}
// 			S[index] = 1
// 			return nil
// 		})
// 	}
// 	err := eg.Wait()
// 	if err != nil {
// 		return errors.Join(err, errors.New("failed to build oracle identities"))
// 	}

// 	switch a.RegistrySettings.RegistryVersion {
// 	case ethereum.RegistryVersion_2_0:
// 		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR2ConfigArgs(a, S, oracleIdentities)
// 		if err != nil {
// 			return errors.Join(err, errors.New("failed to build config args"))
// 		}
// 	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
// 		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR3ConfigArgs(a, S, oracleIdentities)
// 		if err != nil {
// 			return errors.Join(err, errors.New("failed to build config args"))
// 		}
// 	default:
// 		return errors.New("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
// 	}

// 	var signers []common.Address
// 	for _, signer := range signerOnchainPublicKeys {
// 		if len(signer) != 20 {
// 			return fmt.Errorf("OnChainPublicKey '%v' has wrong length for address", signer)
// 		}
// 		signers = append(signers, common.BytesToAddress(signer))
// 	}

// 	var transmitters []common.Address
// 	for _, transmitter := range transmitterAccounts {
// 		if !common.IsHexAddress(string(transmitter)) {
// 			return fmt.Errorf("TransmitAccount '%s' is not a valid Ethereum address", string(transmitter))
// 		}
// 		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
// 	}

// 	ocrConfig := contracts.OCRv2Config{
// 		Signers:               signers,
// 		Transmitters:          transmitters,
// 		F:                     f,
// 		OffchainConfigVersion: offchainConfigVersion,
// 		OffchainConfig:        offchainConfig,
// 	}

// 	if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
// 		ocrConfig.OnchainConfig = a.RegistrySettings.Encode20OnchainConfig(a.Registrar.Address())
// 		err = a.Registry.SetConfig(a.RegistrySettings, ocrConfig)
// 		if err != nil {
// 			return errors.Join(err, errors.New("failed to set config on registry"))
// 		}
// 	} else {
// 		switch a.RegistrySettings.RegistryVersion {
// 		case ethereum.RegistryVersion_2_1:
// 			ocrConfig.TypedOnchainConfig21 = a.RegistrySettings.Create21OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager)
// 		case ethereum.RegistryVersion_2_2:
// 			ocrConfig.TypedOnchainConfig22 = a.RegistrySettings.Create22OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager, a.Registry.ChainModuleAddress(), a.Registry.ReorgProtectionEnabled())
// 		case ethereum.RegistryVersion_2_3:
// 			ocrConfig.TypedOnchainConfig23 = a.RegistrySettings.Create23OnchainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager, a.Registry.ChainModuleAddress(), a.Registry.ReorgProtectionEnabled())
// 			ocrConfig.BillingTokens = []common.Address{
// 				common.HexToAddress(a.LinkToken.Address()),
// 				common.HexToAddress(a.WETHToken.Address()),
// 			}

// 			ocrConfig.BillingConfigs = []i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23BillingConfig{
// 				{
// 					GasFeePPB:         100,
// 					FlatFeeMilliCents: big.NewInt(500),
// 					PriceFeed:         common.HexToAddress(a.ETHUSDFeed.Address()),
// 					Decimals:          18,
// 					FallbackPrice:     big.NewInt(1000),
// 					MinSpend:          big.NewInt(200),
// 				},
// 				{
// 					GasFeePPB:         100,
// 					FlatFeeMilliCents: big.NewInt(500),
// 					PriceFeed:         common.HexToAddress(a.LINKUSDFeed.Address()),
// 					Decimals:          18,
// 					FallbackPrice:     big.NewInt(1000),
// 					MinSpend:          big.NewInt(200),
// 				},
// 			}
// 		}
// 		a.Logger.Debug().Interface("ocrConfig", ocrConfig).Msg("Setting OCR3 config")
// 		err = a.Registry.SetConfigTypeSafe(ocrConfig)
// 		if err != nil {
// 			return errors.Join(err, errors.New("failed to set config on registry"))
// 		}
// 	}
// 	return nil
// }

func calculateOCR2ConfigArgs(pluginConfig ocr2keepers30config.OffchainConfig, publicConfig ocr3.PublicConfig, S []int, oracleIdentities []confighelper.OracleIdentityExtra) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	offC, _ := json.Marshal(ocr2keepers20config.OffchainConfig{
		TargetProbability:    pluginConfig.TargetProbability,
		TargetInRounds:       pluginConfig.TargetInRounds,
		PerformLockoutWindow: pluginConfig.PerformLockoutWindow,
		GasLimitPerReport:    pluginConfig.GasLimitPerReport,
		GasOverheadPerUpkeep: pluginConfig.GasOverheadPerUpkeep,
		MinConfirmations:     pluginConfig.MinConfirmations,
		MaxUpkeepBatchSize:   pluginConfig.MaxUpkeepBatchSize,
	})

	rMax := publicConfig.RMax
	if rMax > math.MaxUint8 {
		panic(fmt.Errorf("rmax overflows uint8: %d", rMax))
	}

	return ocr2.ContractSetConfigArgsForTests(
		publicConfig.DeltaProgress, publicConfig.DeltaResend,
		publicConfig.DeltaRound, publicConfig.DeltaGrace,
		publicConfig.DeltaStage, uint8(rMax),
		S, oracleIdentities, offC,
		nil,
		publicConfig.MaxDurationQuery, publicConfig.MaxDurationObservation,
		1200*time.Millisecond,
		publicConfig.MaxDurationShouldAcceptAttestedReport,
		publicConfig.MaxDurationShouldTransmitAcceptedReport,
		publicConfig.F, publicConfig.OnchainConfig,
	)
}

func calculateOCR3ConfigArgs(pluginConfig ocr2keepers30config.OffchainConfig, publicConfig ocr3.PublicConfig, S []int, oracleIdentities []confighelper.OracleIdentityExtra) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	offC, _ := json.Marshal(pluginConfig)

	return ocr3.ContractSetConfigArgsForTests(
		publicConfig.DeltaProgress, publicConfig.DeltaResend, publicConfig.DeltaInitial,
		publicConfig.DeltaRound, publicConfig.DeltaGrace, publicConfig.DeltaCertifiedCommitRequest,
		publicConfig.DeltaStage, publicConfig.RMax,
		S, oracleIdentities, offC,
		nil, publicConfig.MaxDurationQuery, publicConfig.MaxDurationObservation,
		publicConfig.MaxDurationShouldAcceptAttestedReport,
		publicConfig.MaxDurationShouldTransmitAcceptedReport,
		publicConfig.F, publicConfig.OnchainConfig,
	)
}

type registrationResult struct {
	txHash common.Hash
}

func (r registrationResult) GetResult() common.Hash {
	return r.txHash
}

func (a *AutomationTest) RegisterUpkeeps(upkeepConfigs []UpkeepConfig, maxConcurrency int) ([]common.Hash, error) {
	concurrency, err := GetAndAssertCorrectConcurrency(a.ChainClient, 1)
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
				errorCh <- errors.Join(err, errors.New("failed to get registrar abi"))
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
				errorCh <- errors.Join(err, errors.New("failed to pack registrar request"))
				return
			}
		case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2: // 2.1 and 2.2 use the same registrar
			registrarABI, err = automation_registrar_wrapper2_1.AutomationRegistrarMetaData.GetAbi()
			if err != nil {
				errorCh <- errors.Join(err, errors.New("failed to get registrar abi"))
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
				errorCh <- errors.Join(err, errors.New("failed to pack registrar request"))
				return
			}
		default:
			errorCh <- errors.New("v2.0, v2.1, and v2.2 are the only supported versions")
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
	concurrency, err := GetAndAssertCorrectConcurrency(a.ChainClient, 1)
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
			errorCh <- errors.Join(err, errors.New("failed to confirm upkeep registration"))
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
			errorCh <- errors.New("failed to parse upkeep id from registration receipt")
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

// func (a *AutomationTest) AddJobsAndSetConfig(t *testing.T) {
// 	l := framework.L
// 	err := a.AddBootstrapJob()
// 	require.NoError(t, err, "Error adding bootstrap job")
// 	err = a.AddAutomationJobs()
// 	require.NoError(t, err, "Error adding automation jobs")

// 	l.Info().
// 		Interface("Plugin Config", a.PluginConfig).
// 		Interface("Public Config", a.PublicConfig).
// 		Interface("Registry Settings", a.RegistrySettings).
// 		Interface("Registrar Settings", a.RegistrarSettings).
// 		Msg("Configuring registry")
// 	err = a.SetConfigOnRegistry()
// 	require.NoError(t, err, "Error setting config on registry")
// 	l.Info().Str("Registry Address", a.Registry.Address()).Msg("Successfully setConfig on registry")
// }

// func (a *AutomationTest) SetupAutomationDeployment() error {
// 	return a.setupDeployment(true)
// }

// func (a *AutomationTest) SetupAutomationDeploymentWithoutJobs() error {
// 	return a.setupDeployment(false)
// }

func (a *AutomationTest) LoadContracts() error {
	if err := a.LoadLINK(a.Config.DeployedContracts.LinkToken); err != nil {
		return fmt.Errorf("error loading link token contract: %w", err)
	}

	if err := a.LoadWETH(a.Config.DeployedContracts.Weth); err != nil {
		return fmt.Errorf("error loading weth token contract: %w", err)
	}

	if err := a.LoadLinkEthFeed(a.Config.DeployedContracts.LinkEthFeed); err != nil {
		return fmt.Errorf("error loading link eth feed contract: %w", err)
	}

	if err := a.LoadEthGasFeed(a.Config.DeployedContracts.EthGasFeed); err != nil {
		return fmt.Errorf("error loading gas feed contract: %w", err)
	}

	if err := a.LoadEthUSDFeed(a.Config.DeployedContracts.EthUSDFeed); err != nil {
		return fmt.Errorf("error loading eth usd feed contract: %w", err)
	}

	if err := a.LoadLinkUSDFeed(a.Config.DeployedContracts.LinkUSDFeed); err != nil {
		return fmt.Errorf("error loading link usd feed contract: %w", err)
	}

	if err := a.LoadTranscoder(a.Config.DeployedContracts.Transcoder); err != nil {
		return fmt.Errorf("error loading transcoder contract: %w", err)
	}

	if err := a.LoadRegistry(a.Config.DeployedContracts.Registry, a.Config.DeployedContracts.ChainModule); err != nil {
		return fmt.Errorf("error loading registry contract: %w", err)
	}

	if a.Registry.RegistryOwnerAddress().String() != a.ChainClient.MustGetRootKeyAddress().String() {
		return fmt.Errorf("registry owner address is not the root key address")
	}

	if err := a.LoadRegistrar(a.Config.DeployedContracts.Registrar); err != nil {
		return fmt.Errorf("error loading registrar contract: %w", err)
	}

	return nil
}

// func (a *AutomationTest) setupDeployment(addJobs bool) error {
// 	l := framework.L
// 	err := a.CollectNodeDetails()
// 	if err != nil {
// 		return fmt.Errorf("error collecting node details: %w", err)
// 	}
// 	l.Info().Msg("Collected Node Details")
// 	l.Debug().Interface("Node Details", a.NodeDetails).Msg("Node Details")

// 	if a.Config.DeployedContracts.LinkToken != "" {
// 		linkAddress := a.Config.DeployedContracts.LinkToken
// 		err := a.LoadLINK(linkAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading link token contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployLINK()
// 		if err != nil {
// 			return fmt.Errorf("error deploying link token contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.Weth != "" {
// 		wethAddress := a.Config.DeployedContracts.Weth
// 		err := a.LoadWETH(wethAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading weth token contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployWETH()
// 		if err != nil {
// 			return fmt.Errorf("error deploying weth token contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.LinkEthFeed != "" {
// 		linkEthFeedAddress := a.Config.DeployedContracts.LinkEthFeed
// 		err := a.LoadLinkEthFeed(linkEthFeedAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading link eth feed contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployLinkEthFeed()
// 		if err != nil {
// 			return fmt.Errorf("error deploying link eth feed contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.EthGasFeed != "" {
// 		gasFeedAddress := a.Config.DeployedContracts.EthGasFeed
// 		err := a.LoadEthGasFeed(gasFeedAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading gas feed contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployGasFeed()
// 		if err != nil {
// 			return fmt.Errorf("error deploying gas feed contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.EthUSDFeed != "" {
// 		ethUsdFeedAddress := a.Config.DeployedContracts.EthUSDFeed
// 		err := a.LoadEthUSDFeed(ethUsdFeedAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading eth usd feed contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployEthUSDFeed()
// 		if err != nil {
// 			return fmt.Errorf("error deploying eth usd feed contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.LinkUSDFeed != "" {
// 		linkUsdFeedAddress := a.Config.DeployedContracts.LinkUSDFeed
// 		err := a.LoadLinkUSDFeed(linkUsdFeedAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading link usd feed contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployLinkUSDFeed()
// 		if err != nil {
// 			return fmt.Errorf("error deploying link usd feed contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.Transcoder != "" {
// 		transcoderAddress := a.Config.DeployedContracts.Transcoder
// 		err := a.LoadTranscoder(transcoderAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading transcoder contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployTranscoder()
// 		if err != nil {
// 			return fmt.Errorf("error deploying transcoder contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.Registry != "" && a.Config.DeployedContracts.ChainModule != "" {
// 		chainModuleAddress := a.Config.DeployedContracts.ChainModule
// 		registryAddress := a.Config.DeployedContracts.Registry
// 		err = a.LoadRegistry(registryAddress, chainModuleAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading registry contract: %w", err)
// 		}
// 		if a.Registry.RegistryOwnerAddress().String() != a.ChainClient.MustGetRootKeyAddress().String() {
// 			l.Debug().Str("RootKeyAddress", a.ChainClient.MustGetRootKeyAddress().String()).Str("Registry Owner Address", a.Registry.RegistryOwnerAddress().String()).Msg("Registry owner address is not the root key address")
// 			return fmt.Errorf("registry owner address is not the root key address")
// 		}
// 	} else {
// 		err = a.DeployRegistry()
// 		if err != nil {
// 			return fmt.Errorf("error deploying registry contract: %w", err)
// 		}
// 	}

// 	if a.Config.DeployedContracts.Registrar != "" {
// 		registrarAddress := a.Config.DeployedContracts.Registrar
// 		err = a.LoadRegistrar(registrarAddress)
// 		if err != nil {
// 			return fmt.Errorf("error loading registrar contract: %w", err)
// 		}
// 	} else {
// 		err = a.DeployRegistrar()
// 		if err != nil {
// 			return fmt.Errorf("error deploying registrar contract: %w", err)
// 		}
// 	}

// 	if addJobs {
// 		err = a.AddJobsAndSetConfig()
// 		if err != nil {
// 			return fmt.Errorf("error adding jobs and setting config: %w", err)
// 		}
// 	}

// 	return nil
// }

// SendLinkFundsToDeploymentAddresses sends LINK token to all addresses, but the root one, from the root address. It uses
// Multicall contract to batch all transfers in a single transaction. It also checks if the funds were transferred correctly.
// It's primary use case is to fund addresses that will be used for Upkeep registration (as that requires LINK balance) during
// Automation/Keeper test setup.
func SendLinkFundsToDeploymentAddresses(
	chainClient *seth.Client,
	concurrency,
	totalUpkeeps,
	operationsPerAddress int,
	multicallAddress common.Address,
	linkAmountPerUpkeep *big.Int,
	linkToken contracts.LinkToken,
) error {
	var generateCallData = func(receiver common.Address, amount *big.Int) ([]byte, error) {
		abi, err := link_token_interface.LinkTokenMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
		data, err := abi.Pack("transfer", receiver, amount)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	toTransferToMultiCallContract := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(totalUpkeeps+concurrency)))
	toTransferPerClient := big.NewInt(0).Mul(linkAmountPerUpkeep, big.NewInt(int64(operationsPerAddress+1)))

	// As a hack we use the geth wrapper directly, because we need to access receipt to get block number, which we will use to query the balance
	// This is needed as querying with 'latest' block number very rarely, but still, return stale balance. That's happening even though we wait for
	// the transaction to be mined.
	linkInstance, err := link_token_interface.NewLinkToken(common.HexToAddress(linkToken.Address()), contracts.MustNewWrappedContractBackend(nil, chainClient))
	if err != nil {
		return err
	}
	// TODO: 2:32PM WRN No matching event with valid indexed parameter count found for log Signature=0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef Transaction=0xde63b61cb6f22db882cee1f28c7a11159be3f4e5f07c71998eada723ccf9e6cd
	tx, err := chainClient.Decode(linkInstance.Transfer(chainClient.NewTXOpts(), multicallAddress, toTransferToMultiCallContract))
	if err != nil {
		return err
	}

	if tx.Receipt == nil {
		return errors.New("transaction receipt for LINK transfer to multicall contract is nil")
	}

	multiBalance, err := linkInstance.BalanceOf(&bind.CallOpts{From: chainClient.Addresses[0], BlockNumber: tx.Receipt.BlockNumber}, multicallAddress)
	if err != nil {
		return pkg_errors.Wrapf(err, "Error getting LINK balance of multicall contract")
	}

	// Old code that's querying latest block
	// err := linkToken.Transfer(multicallAddress.Hex(), toTransferToMultiCallContract)
	// if err != nil {
	//	return errors.Wrapf(err, "Error transferring LINK to multicall contract")
	//}
	//
	// balance, err := linkToken.BalanceOf(context.Background(), multicallAddress.Hex())
	// if err != nil {
	//	return errors.Wrapf(err, "Error getting LINK balance of multicall contract")
	//}

	if multiBalance.Cmp(toTransferToMultiCallContract) < 0 {
		return fmt.Errorf("Incorrect LINK balance of multicall contract. Expected at least: %s. Got: %s", toTransferToMultiCallContract.String(), multiBalance.String())
	}

	// Transfer LINK to ephemeral keys
	multiCallData := make([][]byte, 0)
	for i := 1; i <= concurrency; i++ {
		data, err := generateCallData(chainClient.Addresses[i], toTransferPerClient)
		if err != nil {
			return pkg_errors.Wrapf(err, "Error generating call data for LINK transfer")
		}
		multiCallData = append(multiCallData, data)
	}

	var call []contracts.Call
	for _, d := range multiCallData {
		data := contracts.Call{Target: common.HexToAddress(linkToken.Address()), AllowFailure: false, CallData: d}
		call = append(call, data)
	}

	multiCallABI, err := abi.JSON(strings.NewReader(contracts.MultiCallABI))
	if err != nil {
		return pkg_errors.Wrapf(err, "Error getting Multicall contract ABI")
	}
	boundContract := bind.NewBoundContract(multicallAddress, multiCallABI, chainClient.Client, chainClient.Client, chainClient.Client)
	// call aggregate3 to group all msg call data and send them in a single transaction // TODO: 2:33PM WRN No matching event with valid indexed parameter count found for log Signature=0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef Transaction=0xe277bace889e96bfba919283b6f1bfaaadc89cda0cb254845394704fff5bc2b8
	ephemeralTx, err := chainClient.Decode(boundContract.Transact(chainClient.NewTXOpts(), "aggregate3", call))
	if err != nil {
		return pkg_errors.Wrapf(err, "Error calling Multicall contract")
	}

	if ephemeralTx.Receipt == nil {
		return pkg_errors.New("transaction receipt for LINK transfer to ephemeral keys is nil")
	}

	for i := 1; i <= concurrency; i++ {
		ephemeralBalance, err := linkInstance.BalanceOf(&bind.CallOpts{From: chainClient.Addresses[0], BlockNumber: ephemeralTx.Receipt.BlockNumber}, chainClient.Addresses[i])
		// Old code that's querying latest block, for now we prefer to use block number from the transaction receipt
		// balance, err := linkToken.BalanceOf(context.Background(), chainClient.Addresses[i].Hex())
		if err != nil {
			return pkg_errors.Wrapf(err, "Error getting LINK balance of ephemeral key %d", i)
		}
		if ephemeralBalance.Cmp(toTransferPerClient) < 0 {
			return fmt.Errorf("Incorrect LINK balance after transfer. Ephemeral key %d. Expected: %s. Got: %s", i, toTransferPerClient.String(), ephemeralBalance.String())
		}
	}

	return nil
}
