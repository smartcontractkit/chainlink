package automation

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	devenv_ocr2 "github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

type NodeDetail struct {
	P2PId                 string
	TransmitterAddresses  []string
	OCR2ConfigPublicKey   string
	OCR2OffchainPublicKey string
	OCR2OnChainPublicKey  string
	OCR2Id                string
}

type NodeDetails struct {
	NodeDetails     []NodeDetail
	P2PBootstrapper string
}

func collectNodeDetails(chainID uint64, nodes []*clclient.ChainlinkClient, clNodes []*clnode.Output) (*NodeDetails, error) {
	nodeDetails := NodeDetails{
		NodeDetails: make([]NodeDetail, 0),
	}

	for i, node := range nodes {
		nodeDetail := NodeDetail{}
		P2PIds, err := node.MustReadP2PKeys()
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to read P2P keys from node %d", i))
		}
		nodeDetail.P2PId = P2PIds.Data[0].Attributes.PeerID

		OCR2Keys, err := node.MustReadOCR2Keys()
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to read OCR2 keys from node %d", i))
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

		TransmitterKeys, err := node.EthAddressesForChain(strconv.FormatUint(chainID, 10))
		nodeDetail.TransmitterAddresses = make([]string, 0)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to read Transmitter keys from node %d", i))
		}
		nodeDetail.TransmitterAddresses = append(nodeDetail.TransmitterAddresses, TransmitterKeys...)
		nodeDetails.NodeDetails = append(nodeDetails.NodeDetails, nodeDetail)
	}

	L.Info().Msg("Collected Node Details")
	L.Debug().Interface("Node Details", nodeDetails.NodeDetails).Msg("Node Details")

	nodeDetails.P2PBootstrapper = fmt.Sprintf("%s@%s:%d", nodeDetails.NodeDetails[0].P2PId, clNodes[0].Node.ContainerName, 6690)
	return &nodeDetails, nil
}

func deployContracts(chainClient *seth.Client, config *Automation) error {
	l := framework.L
	if config.DeployedContracts.LinkToken == "" {
		addr, err := DeployLINK(l, chainClient)
		if err != nil {
			return fmt.Errorf("error deploying link token contract: %w", err)
		}
		config.DeployedContracts.LinkToken = addr
	}

	if config.DeployedContracts.Weth == "" {
		addr, err := DeployWETH(l, chainClient)
		if err != nil {
			return fmt.Errorf("error deploying weth token contract: %w", err)
		}
		config.DeployedContracts.Weth = addr
	}

	if config.DeployedContracts.LinkEthFeed == "" {
		addr, err := DeployLinkEthFeed(chainClient, config.RegistrySettings.FallbackLinkPrice)
		if err != nil {
			return fmt.Errorf("error deploying link eth feed contract: %w", err)
		}
		config.DeployedContracts.LinkEthFeed = addr
	}

	if config.DeployedContracts.EthGasFeed == "" {
		addr, err := DeployGasFeed(chainClient, config.RegistrySettings.FallbackGasPrice)
		if err != nil {
			return fmt.Errorf("error deploying gas feed contract: %w", err)
		}
		config.DeployedContracts.EthGasFeed = addr
	}

	if config.DeployedContracts.EthUSDFeed == "" {
		addr, err := DeployEthUSDFeed(chainClient, config.RegistrySettings.FallbackLinkPrice)
		if err != nil {
			return fmt.Errorf("error deploying eth usd feed contract: %w", err)
		}
		config.DeployedContracts.EthUSDFeed = addr
	}

	if config.DeployedContracts.LinkUSDFeed == "" {
		addr, err := DeployLinkUSDFeed(chainClient, config.RegistrySettings.FallbackLinkPrice)
		if err != nil {
			return fmt.Errorf("error deploying link usd feed contract: %w", err)
		}
		config.DeployedContracts.LinkUSDFeed = addr
	}

	if config.DeployedContracts.Transcoder == "" {
		addr, err := DeployTranscoder(chainClient)
		if err != nil {
			return fmt.Errorf("error deploying transcoder contract: %w", err)
		}
		config.DeployedContracts.Transcoder = addr
	}

	if config.DeployedContracts.Registry == "" {
		registryAddr, chainModuleAddr, err := DeployRegistry(chainClient, config.MustGetRegistryVersion(), config)
		if err != nil {
			return fmt.Errorf("error deploying registry contract: %w", err)
		}
		config.DeployedContracts.Registry = registryAddr
		config.DeployedContracts.ChainModule = chainModuleAddr
	}

	if config.DeployedContracts.Registrar == "" {
		addr, err := DeployRegistrar(chainClient, config.MustGetRegistryVersion(), config)
		if err != nil {
			return fmt.Errorf("error deploying registrar contract: %w", err)
		}
		config.DeployedContracts.Registrar = addr
	}

	if config.DeployedContracts.MultiCall == "" {
		addr, err := DeployMultiCall(chainClient)
		if err != nil {
			return fmt.Errorf("error deploying multi call contract: %w", err)
		}
		config.DeployedContracts.MultiCall = addr
	}

	return nil
}

func DeployLINK(logger zerolog.Logger, chainClient *seth.Client) (string, error) {
	linkToken, err := contracts.DeployLinkTokenContract(logger, chainClient)
	if err != nil {
		return "", err
	}
	return linkToken.Address(), nil
}

func DeployWETH(logger zerolog.Logger, chainClient *seth.Client) (string, error) {
	wethToken, err := contracts.DeployWETHTokenContract(logger, chainClient)
	if err != nil {
		return "", err
	}
	return wethToken.Address(), nil
}

func DeployTranscoder(chainClient *seth.Client) (string, error) {
	transcoder, err := contracts.DeployUpkeepTranscoder(chainClient)
	if err != nil {
		return "", err
	}
	return transcoder.Address(), nil
}

func DeployLinkEthFeed(chainClient *seth.Client, fallbackLinkPrice *big.Int) (string, error) {
	ethLinkFeed, err := contracts.DeployMockLINKETHFeed(chainClient, fallbackLinkPrice)
	if err != nil {
		return "", err
	}
	return ethLinkFeed.Address(), nil
}

func DeployEthUSDFeed(chainClient *seth.Client, fallbackPrice *big.Int) (string, error) {
	ethUSDFeed, err := contracts.DeployMockETHUSDFeed(chainClient, fallbackPrice)
	if err != nil {
		return "", err
	}
	return ethUSDFeed.Address(), nil
}

func DeployLinkUSDFeed(chainClient *seth.Client, fallbackPrice *big.Int) (string, error) {
	linkUSDFeed, err := contracts.DeployMockETHUSDFeed(chainClient, fallbackPrice)
	if err != nil {
		return "", err
	}
	return linkUSDFeed.Address(), nil
}

func DeployGasFeed(chainClient *seth.Client, fallbackGasPrice *big.Int) (string, error) {
	gasFeed, err := contracts.DeployMockGASFeed(chainClient, fallbackGasPrice)
	if err != nil {
		return "", err
	}
	return gasFeed.Address(), nil
}

func DeployRegistry(chainClient *seth.Client, registryVersion contracts.KeeperRegistryVersion, config *Automation) (registryAddr, chainModuleAddr string, err error) {
	registryOpts := &contracts.KeeperRegistryOpts{
		RegistryVersion:   registryVersion,
		LinkAddr:          config.DeployedContracts.LinkToken,
		ETHFeedAddr:       config.DeployedContracts.LinkEthFeed,
		GasFeedAddr:       config.DeployedContracts.EthGasFeed,
		TranscoderAddr:    config.DeployedContracts.Transcoder,
		RegistrarAddr:     utils.ZeroAddress.Hex(),
		Settings:          config.GetRegistryConfig(),
		LinkUSDFeedAddr:   config.DeployedContracts.EthUSDFeed,
		NativeUSDFeedAddr: config.DeployedContracts.LinkUSDFeed,
		WrappedNativeAddr: config.DeployedContracts.Weth,
	}
	registry, err := contracts.DeployKeeperRegistry(chainClient, registryOpts)
	if err != nil {
		return "", "", err
	}
	return registry.Address(), registry.ChainModuleAddress().Hex(), nil
}

func DeployRegistrar(chainClient *seth.Client, registryVersion contracts.KeeperRegistryVersion, config *Automation) (string, error) {
	if config.DeployedContracts.Registry == "" {
		return "", errors.New("registry must be deployed before registrar")
	}
	registrarSettings := contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: uint8(2),
		AutoApproveMaxAllowed: math.MaxUint16,
		MinLinkJuels:          big.NewInt(0),
		RegistryAddr:          config.DeployedContracts.Registry,
		WETHTokenAddr:         config.DeployedContracts.Weth,
	}

	registrar, err := contracts.DeployKeeperRegistrar(chainClient, registryVersion, config.DeployedContracts.LinkToken, registrarSettings)
	if err != nil {
		return "", err
	}
	return registrar.Address(), nil
}

func DeployMultiCall(chainClient *seth.Client) (string, error) {
	multiCall, err := contracts.DeployMultiCallContract(chainClient)
	if err != nil {
		return "", err
	}
	return multiCall.Hex(), nil
}

func LoadRegistry(chainClient *seth.Client, registryAddress, chainModuleAddress string, registryVersion contracts.KeeperRegistryVersion) (contracts.KeeperRegistry, error) {
	registry, err := contracts.LoadKeeperRegistry(L, chainClient, common.HexToAddress(registryAddress), registryVersion, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return nil, err
	}
	L.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", registryAddress).Msg("Successfully loaded Registry")
	return registry, nil
}

func createJobs(nodes []*clclient.ChainlinkClient, nodeDetails *NodeDetails, chainID int, registryVersion contracts.KeeperRegistryVersion, registryAddress string, mercuryCredentialName string) error {
	if err := addBootstrapJob(nodes[0], chainID, registryAddress); err != nil {
		return err
	}

	return addAutomationJobs(nodes, nodeDetails, chainID, registryVersion, registryAddress, mercuryCredentialName)
}

func addBootstrapJob(bootstrapNode *clclient.ChainlinkClient, chainID int, registryAddress string) error {
	bootstrapSpec := &devenv_ocr2.TaskJobSpec{
		Name:    "ocr2 bootstrap node " + registryAddress,
		JobType: "bootstrap",
		OCR2OracleSpec: devenv_ocr2.OracleSpec{
			ContractID: registryAddress,
			Relay:      "evm",
			RelayConfig: map[string]any{
				"chainID": chainID,
			},
			ContractConfigTrackerPollInterval: *devenv_ocr2.NewInterval(time.Second * 15),
		},
	}
	_, err := bootstrapNode.MustCreateJob(bootstrapSpec)
	if err != nil {
		return errors.Join(err, errors.New("failed to create bootstrap job on bootstrap node"))
	}
	return nil
}

func addAutomationJobs(nodes []*clclient.ChainlinkClient, nodeDetails *NodeDetails, chainID int, registryVersion contracts.KeeperRegistryVersion, registryAddress string, mercuryCredentialName string) error {
	var contractVersion string
	switch registryVersion {
	case contracts.RegistryVersion_2_2, contracts.RegistryVersion_2_3:
		contractVersion = "v2.1+"
	case contracts.RegistryVersion_2_1:
		contractVersion = "v2.1"
	case contracts.RegistryVersion_2_0:
		contractVersion = "v2.0"
	default:
		return errors.New("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
	}
	pluginCfg := map[string]any{
		"contractVersion": "\"" + contractVersion + "\"",
	}
	if strings.Contains(contractVersion, "v2.1") {
		if mercuryCredentialName != "" {
			pluginCfg["mercuryCredentialName"] = "\"" + mercuryCredentialName + "\""
		}
	}
	for i := 1; i < len(nodes); i++ {
		autoOCR2JobSpec := devenv_ocr2.TaskJobSpec{
			Name:    "automation-" + contractVersion + "-" + registryAddress,
			JobType: "offchainreporting2",
			OCR2OracleSpec: devenv_ocr2.OracleSpec{
				PluginType: "ocr2automation",
				ContractID: registryAddress,
				Relay:      "evm",
				RelayConfig: map[string]any{
					"chainID": chainID,
				},
				PluginConfig:                      pluginCfg,
				ContractConfigTrackerPollInterval: *devenv_ocr2.NewInterval(time.Second * 15),
				TransmitterID:                     null.StringFrom(nodeDetails.NodeDetails[i].TransmitterAddresses[0]), // TODO benchmark test might need to be set dynamically
				P2PV2Bootstrappers:                pq.StringArray{nodeDetails.P2PBootstrapper},
				OCRKeyBundleID:                    null.StringFrom(nodeDetails.NodeDetails[i].OCR2Id),
			},
		}
		_, err := nodes[i].MustCreateJob(&autoOCR2JobSpec)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to create OCR2 job on node %d", i+1))
		}
	}
	return nil
}

func setConfigOnRegistry(nodeDetails *NodeDetails, config *Automation, chainClient *seth.Client) error {
	donNodes := nodeDetails.NodeDetails[1:]
	S := make([]int, len(donNodes))
	oracleIdentities := make([]ocr2.OracleIdentityExtra, len(donNodes))
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
				return errors.New("wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return errors.New("wrong number of elements copied")
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(chainlinkNode.OCR2OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = ocr2.OracleIdentityExtra{
				OracleIdentity: ocr2.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            chainlinkNode.P2PId,
					TransmitAccount:   types.Account(chainlinkNode.TransmitterAddresses[0]), // TODO benchmark test might need to be set dynamically
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return errors.Join(err, errors.New("failed to build oracle identities"))
	}

	registrySettings := config.GetRegistryConfig()
	switch registrySettings.RegistryVersion {
	case contracts.RegistryVersion_2_0:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR2ConfigArgs(config.GetPluginConfig(), config.GetPublicConfig(), S, oracleIdentities)
		if err != nil {
			return errors.Join(err, errors.New("failed to build config args"))
		}
	case contracts.RegistryVersion_2_1, contracts.RegistryVersion_2_2, contracts.RegistryVersion_2_3:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR3ConfigArgs(config.GetPluginConfig(), config.GetPublicConfig(), S, oracleIdentities)
		if err != nil {
			return errors.Join(err, errors.New("failed to build config args"))
		}
	default:
		return errors.New("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
	}

	signers := []common.Address{}
	for _, signer := range signerOnchainPublicKeys {
		if len(signer) != 20 {
			return fmt.Errorf("OnChainPublicKey '%v' has wrong length for address", signer)
		}
		signers = append(signers, common.BytesToAddress(signer))
	}

	transmitters := []common.Address{}
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

	registry, err := LoadRegistry(chainClient, config.DeployedContracts.Registry, config.DeployedContracts.ChainModule, registrySettings.RegistryVersion)
	if err != nil {
		return errors.Join(err, errors.New("failed to load registry"))
	}

	if registrySettings.RegistryVersion == contracts.RegistryVersion_2_0 {
		ocrConfig.OnchainConfig = registrySettings.Encode20OnchainConfig(config.DeployedContracts.Registrar)
		err = registry.SetConfig(registrySettings, ocrConfig)
		if err != nil {
			return errors.Join(err, errors.New("failed to set config on registry"))
		}
	} else {
		switch registrySettings.RegistryVersion {
		case contracts.RegistryVersion_2_1:
			ocrConfig.TypedOnchainConfig21 = registrySettings.Create21OnchainConfig(config.DeployedContracts.Registrar, chainClient.MustGetRootKeyAddress())
		case contracts.RegistryVersion_2_2:
			ocrConfig.TypedOnchainConfig22 = registrySettings.Create22OnchainConfig(config.DeployedContracts.Registrar, chainClient.MustGetRootKeyAddress(), common.HexToAddress(config.DeployedContracts.ChainModule), registry.ReorgProtectionEnabled())
		case contracts.RegistryVersion_2_3:
			ocrConfig.TypedOnchainConfig23 = registrySettings.Create23OnchainConfig(config.DeployedContracts.Registrar, chainClient.MustGetRootKeyAddress(), common.HexToAddress(config.DeployedContracts.ChainModule), registry.ReorgProtectionEnabled())
			ocrConfig.BillingTokens = []common.Address{
				common.HexToAddress(config.DeployedContracts.LinkToken),
				common.HexToAddress(config.DeployedContracts.Weth),
			}

			ocrConfig.BillingConfigs = []i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23BillingConfig{
				{
					GasFeePPB:         100,
					FlatFeeMilliCents: big.NewInt(500),
					PriceFeed:         common.HexToAddress(config.DeployedContracts.EthUSDFeed),
					Decimals:          18,
					FallbackPrice:     big.NewInt(1000),
					MinSpend:          big.NewInt(200),
				},
				{
					GasFeePPB:         100,
					FlatFeeMilliCents: big.NewInt(500),
					PriceFeed:         common.HexToAddress(config.DeployedContracts.LinkUSDFeed),
					Decimals:          18,
					FallbackPrice:     big.NewInt(1000),
					MinSpend:          big.NewInt(200),
				},
			}
		default:
			return fmt.Errorf("unsupported registry version: %s", registrySettings.RegistryVersion.String())
		}
		L.Debug().Interface("ocrConfig", ocrConfig).Msg("Setting OCR3 config")
		err = registry.SetConfigTypeSafe(ocrConfig)
		if err != nil {
			return errors.Join(err, errors.New("failed to set config on registry"))
		}
	}
	return nil
}

func calculateOCR2ConfigArgs(
	pluginConfig ocr2keepers30config.OffchainConfig,
	publicConfig ocr3.PublicConfig,
	S []int, //nolint:gocritic //S param name is capitalised on purpose
	oracleIdentities []ocr2.OracleIdentityExtra,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8, //nolint:revive //we want to
	onchainConfig_ []byte, //nolint:revive //we want to
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

//nolint:gocritic // S is capitalized intentionally
func calculateOCR3ConfigArgs(
	pluginConfig ocr2keepers30config.OffchainConfig,
	publicConfig ocr3.PublicConfig,
	S []int,
	oracleIdentities []ocr2.OracleIdentityExtra,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8, //nolint:revive //we want to use underscores
	onchainConfig_ []byte, //nolint:revive //we want to use underscores
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
