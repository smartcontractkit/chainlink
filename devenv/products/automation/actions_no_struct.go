package automation

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/lib/pq"
	pkg_errors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	devenv_ocr2 "github.com/smartcontractkit/chainlink/devenv/products/ocr2"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"
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

func CollectNodeDetails(chainID uint64, nodes []*clclient.ChainlinkClient, clNodes []*clnode.Output) (*NodeDetails, error) {
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

func DeployRegistry(chainClient *seth.Client, registryVersion ethereum.KeeperRegistryVersion, config *Automation) (registryAddr, chainModuleAddr string, err error) {
	registryOpts := &contracts.KeeperRegistryOpts{
		RegistryVersion:   registryVersion,
		LinkAddr:          config.DeployedContracts.LinkToken,
		ETHFeedAddr:       config.DeployedContracts.LinkEthFeed,
		GasFeedAddr:       config.DeployedContracts.EthGasFeed,
		TranscoderAddr:    config.DeployedContracts.Transcoder,
		RegistrarAddr:     utils.ZeroAddress.Hex(),
		Settings:          ReadRegistryConfig(config),
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

func DeployRegistrar(chainClient *seth.Client, registryVersion ethereum.KeeperRegistryVersion, config *Automation) (string, error) {
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

func LoadRegistry(chainClient *seth.Client, registryAddress, chainModuleAddress string, registryVersion ethereum.KeeperRegistryVersion) (contracts.KeeperRegistry, error) {
	registry, err := contracts.LoadKeeperRegistry(L, chainClient, common.HexToAddress(registryAddress), registryVersion, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return nil, err
	}
	L.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", registryAddress).Msg("Successfully loaded Registry")
	return registry, nil
}

func createJobs(nodes []*clclient.ChainlinkClient, nodeDetails *NodeDetails, chainID int, registryVersion ethereum.KeeperRegistryVersion, registryAddress string, mercuryCredentialName string) error {
	if err := AddBootstrapJob(nodes[0], chainID, registryAddress); err != nil {
		return err
	}

	return AddAutomationJobs(nodes, nodeDetails, chainID, registryVersion, registryAddress, mercuryCredentialName)
}

func AddBootstrapJob(bootstrapNode *clclient.ChainlinkClient, chainID int, registryAddress string) error {
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

func AddAutomationJobs(nodes []*clclient.ChainlinkClient, nodeDetails *NodeDetails, chainID int, registryVersion ethereum.KeeperRegistryVersion, registryAddress string, mercuryCredentialName string) error {
	var contractVersion string
	switch registryVersion {
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

func SetConfigOnRegistry(nodeDetails *NodeDetails, config *Automation, chainClient *seth.Client) error {
	donNodes := nodeDetails.NodeDetails[1:]
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
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
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

	registrySettings := ReadRegistryConfig(config)
	switch registrySettings.RegistryVersion {
	case ethereum.RegistryVersion_2_0:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR2ConfigArgs(ReadPluginConfig(config.PluginConfig), ReadPublicConfig(config.PublicConfig), S, oracleIdentities)
		if err != nil {
			return errors.Join(err, errors.New("failed to build config args"))
		}
	case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2, ethereum.RegistryVersion_2_3:
		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = calculateOCR3ConfigArgs(ReadPluginConfig(config.PluginConfig), ReadPublicConfig(config.PublicConfig), S, oracleIdentities)
		if err != nil {
			return errors.Join(err, errors.New("failed to build config args"))
		}
	default:
		return errors.New("v2.0, v2.1, v2.2 and v2.3 are the only supported versions")
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

	registry, err := LoadRegistry(chainClient, config.DeployedContracts.Registry, config.DeployedContracts.ChainModule, registrySettings.RegistryVersion)
	if err != nil {
		return errors.Join(err, errors.New("failed to load registry"))
	}

	if registrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
		ocrConfig.OnchainConfig = registrySettings.Encode20OnchainConfig(config.DeployedContracts.Registrar)
		err = registry.SetConfig(registrySettings, ocrConfig)
		if err != nil {
			return errors.Join(err, errors.New("failed to set config on registry"))
		}
	} else {
		switch registrySettings.RegistryVersion {
		case ethereum.RegistryVersion_2_1:
			ocrConfig.TypedOnchainConfig21 = registrySettings.Create21OnchainConfig(config.DeployedContracts.Registrar, chainClient.MustGetRootKeyAddress())
		case ethereum.RegistryVersion_2_2:
			ocrConfig.TypedOnchainConfig22 = registrySettings.Create22OnchainConfig(config.DeployedContracts.Registrar, chainClient.MustGetRootKeyAddress(), common.HexToAddress(config.DeployedContracts.ChainModule), registry.ReorgProtectionEnabled())
		case ethereum.RegistryVersion_2_3:
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
		}
		L.Debug().Interface("ocrConfig", ocrConfig).Msg("Setting OCR3 config")
		err = registry.SetConfigTypeSafe(ocrConfig)
		if err != nil {
			return errors.Join(err, errors.New("failed to set config on registry"))
		}
	}
	return nil
}

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

// GenerateUpkeepReport generates a report of performed, successful, reverted and stale upkeeps for a given registry contract based on transaction logs. In case of test failure it can help us
// to triage the issue by providing more context.
func GenerateUpkeepReport(t *testing.T, chainClient *seth.Client, startBlock, endBlock *big.Int, instance contracts.KeeperRegistry, registryVersion ethereum.KeeperRegistryVersion) (performedUpkeeps, successfulUpkeeps, revertedUpkeeps, staleUpkeeps int, err error) {
	registryLogs := []gethtypes.Log{}
	l := framework.L

	var (
		blockBatchSize  int64 = 100
		logs            []gethtypes.Log
		timeout         = 5 * time.Second
		addr            = common.HexToAddress(instance.Address())
		queryStartBlock = startBlock
	)

	// Gather logs from the registry in 100 block chunks to avoid read limits
	for queryStartBlock.Cmp(endBlock) < 0 {
		filterQuery := geth.FilterQuery{
			Addresses: []common.Address{addr},
			FromBlock: queryStartBlock,
			ToBlock:   big.NewInt(0).Add(queryStartBlock, big.NewInt(blockBatchSize)),
		}

		// This RPC call can possibly time out or otherwise die. Failure is not an option, keep retrying to get our stats.
		err = errors.New("initial error") // to ensure our for loop runs at least once
		for err != nil {
			ctx, cancel := context.WithTimeout(t.Context(), timeout)
			logs, err = chainClient.Client.FilterLogs(ctx, filterQuery)
			cancel()
			if err != nil {
				l.Error().
					Err(err).
					Interface("Filter Query", filterQuery).
					Str("Timeout", timeout.String()).
					Msg("Error getting logs from chain, trying again")
				timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
				continue
			}
			l.Info().
				Uint64("From Block", queryStartBlock.Uint64()).
				Uint64("To Block", filterQuery.ToBlock.Uint64()).
				Int("Log Count", len(logs)).
				Str("Registry Address", addr.Hex()).
				Msg("Collected logs")
			queryStartBlock.Add(queryStartBlock, big.NewInt(blockBatchSize))
			registryLogs = append(registryLogs, logs...)
		}
	}

	var contractABI *abi.ABI
	contractABI, err = contracts.GetRegistryContractABI(registryVersion)
	if err != nil {
		return
	}

	for _, allLogs := range registryLogs {
		log := allLogs
		var eventDetails *abi.Event
		eventDetails, err = contractABI.EventByID(log.Topics[0])
		if err != nil {
			l.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error getting event details for log, report data inaccurate")
			break
		}
		if eventDetails.Name == "UpkeepPerformed" {
			performedUpkeeps++
			var parsedLog *contracts.UpkeepPerformedLog
			parsedLog, err = instance.ParseUpkeepPerformedLog(&log)
			if err != nil {
				l.Error().Err(err).Str("Log Hash", log.TxHash.Hex()).Msg("Error parsing upkeep performed log, report data inaccurate")
				break
			}
			if !parsedLog.Success {
				revertedUpkeeps++
			} else {
				successfulUpkeeps++
			}
		} else if eventDetails.Name == "StaleUpkeepReport" {
			staleUpkeeps++
		}
	}

	return
}

func GetStalenessReportCleanupFn(t *testing.T, logger zerolog.Logger, chainClient *seth.Client, startBlock uint64, registry contracts.KeeperRegistry, registryVersion ethereum.KeeperRegistryVersion) func() {
	return func() {
		if t.Failed() {
			endBlock, err := chainClient.Client.BlockNumber(t.Context())
			require.NoError(t, err, "Failed to get end block")

			total, ok, reverted, stale, err := GenerateUpkeepReport(t, chainClient, new(big.Int).SetUint64(startBlock), new(big.Int).SetUint64(endBlock), registry, registryVersion)
			require.NoError(t, err, "Failed to get staleness data")
			if stale > 0 || reverted > 0 {
				logger.Warn().Int("Total upkeeps", total).Int("Successful upkeeps", ok).Int("Reverted Upkeeps", reverted).Int("Stale Upkeeps", stale).Msg("Staleness data")
			} else {
				logger.Info().Int("Total upkeeps", total).Int("Successful upkeeps", ok).Int("Reverted Upkeeps", reverted).Int("Stale Upkeeps", stale).Msg("Staleness data")
			}
		}
	}
}

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

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []*clclient.ChainlinkClient) ([]common.Address, error) {
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
