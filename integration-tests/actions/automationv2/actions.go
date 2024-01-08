package automationv2

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	ocr2keepers20config "github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registrar_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	ctfTestEnv "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
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
	ChainClient blockchain.EVMClient
	Deployer    contracts.ContractDeployer

	LinkToken   contracts.LinkToken
	Transcoder  contracts.UpkeepTranscoder
	EthLinkFeed contracts.MockETHLINKFeed
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

	ChainlinkNodesk8s []*client.ChainlinkK8sClient
	ChainlinkNodes    []*client.ChainlinkClient

	DockerEnv *test_env.CLClusterTestEnv

	NodeDetails              []NodeDetails
	DefaultP2Pv2Bootstrapper string
	MercuryCredentialName    string
	TransmitterKeyIndex      int
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
	chainClient blockchain.EVMClient,
	deployer contracts.ContractDeployer,
	chainlinkNodes []*client.ChainlinkK8sClient,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		Deployer:               deployer,
		ChainlinkNodesk8s:      chainlinkNodes,
		IsOnk8s:                true,
		TransmitterKeyIndex:    0,
		UpkeepPrivilegeManager: common.HexToAddress(chainClient.GetDefaultWallet().Address()),
	}
}

func NewAutomationTestDocker(
	chainClient blockchain.EVMClient,
	deployer contracts.ContractDeployer,
	chainlinkNodes []*client.ChainlinkClient,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		Deployer:               deployer,
		ChainlinkNodes:         chainlinkNodes,
		IsOnk8s:                false,
		TransmitterKeyIndex:    0,
		UpkeepPrivilegeManager: common.HexToAddress(chainClient.GetDefaultWallet().Address()),
	}
}

func (a *AutomationTest) SetIsOnk8s(flag bool) {
	a.IsOnk8s = flag
}

func (a *AutomationTest) SetMercuryCredentialName(name string) {
	a.MercuryCredentialName = name
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
	linkToken, err := a.Deployer.DeployLinkTokenContract()
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for link token contract to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadLINK(address string) error {
	linkToken, err := a.Deployer.LoadLinkToken(common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	return nil
}

func (a *AutomationTest) DeployTranscoder() error {
	transcoder, err := a.Deployer.DeployUpkeepTranscoder()
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for transcoder contract to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadTranscoder(address string) error {
	transcoder, err := a.Deployer.LoadUpkeepTranscoder(common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	return nil
}

func (a *AutomationTest) DeployEthLinkFeed() error {
	ethLinkFeed, err := a.Deployer.DeployMockETHLINKFeed(a.RegistrySettings.FallbackLinkPrice)
	if err != nil {
		return err
	}
	a.EthLinkFeed = ethLinkFeed
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for Mock ETH LINK feed to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadEthLinkFeed(address string) error {
	ethLinkFeed, err := a.Deployer.LoadETHLINKFeed(common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.EthLinkFeed = ethLinkFeed
	return nil
}

func (a *AutomationTest) DeployGasFeed() error {
	gasFeed, err := a.Deployer.DeployMockGasFeed(a.RegistrySettings.FallbackGasPrice)
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for mock gas feed to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadEthGasFeed(address string) error {
	gasFeed, err := a.Deployer.LoadGasFeed(common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	return nil
}

func (a *AutomationTest) DeployRegistry() error {
	registryOpts := &contracts.KeeperRegistryOpts{
		RegistryVersion: a.RegistrySettings.RegistryVersion,
		LinkAddr:        a.LinkToken.Address(),
		ETHFeedAddr:     a.EthLinkFeed.Address(),
		GasFeedAddr:     a.GasFeed.Address(),
		TranscoderAddr:  a.Transcoder.Address(),
		RegistrarAddr:   utils.ZeroAddress.Hex(),
		Settings:        a.RegistrySettings,
	}
	registry, err := a.Deployer.DeployKeeperRegistry(registryOpts)
	if err != nil {
		return err
	}
	a.Registry = registry
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for registry contract to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadRegistry(address string) error {
	registry, err := a.Deployer.LoadKeeperRegistry(common.HexToAddress(address), a.RegistrySettings.RegistryVersion)
	if err != nil {
		return err
	}
	a.Registry = registry
	return nil
}

func (a *AutomationTest) DeployRegistrar() error {
	if a.Registry == nil {
		return fmt.Errorf("registry must be deployed or loaded before registrar")
	}
	a.RegistrarSettings.RegistryAddr = a.Registry.Address()
	registrar, err := a.Deployer.DeployKeeperRegistrar(a.RegistrySettings.RegistryVersion, a.LinkToken.Address(), a.RegistrarSettings)
	if err != nil {
		return err
	}
	a.Registrar = registrar
	err = a.ChainClient.WaitForEvents()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed waiting for registrar contract to deploy"))
	}
	return nil
}

func (a *AutomationTest) LoadRegistrar(address string) error {
	if a.Registry == nil {
		return fmt.Errorf("registry must be deployed or loaded before registrar")
	}
	a.RegistrarSettings.RegistryAddr = a.Registry.Address()
	registrar, err := a.Deployer.LoadKeeperRegistrar(common.HexToAddress(address), a.RegistrySettings.RegistryVersion)
	if err != nil {
		return err
	}
	a.Registrar = registrar
	return nil
}

func (a *AutomationTest) CollectNodeDetails() error {
	var (
		nodes []*client.ChainlinkClient
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

		TransmitterKeys, err := node.EthAddressesForChain(a.ChainClient.GetChainID().String())
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
	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node " + a.Registry.Address(),
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: a.Registry.Address(),
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(a.ChainClient.GetChainID().Int64()),
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
	if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_1 {
		contractVersion = "v2.1"
	} else if a.RegistrySettings.RegistryVersion == ethereum.RegistryVersion_2_0 {
		contractVersion = "v2.0"
	} else {
		return fmt.Errorf("v2.0 and v2.1 are the only supported versions")
	}
	for i := 1; i < len(a.ChainlinkNodes); i++ {
		autoOCR2JobSpec := client.OCR2TaskJobSpec{
			Name:    "automation-" + contractVersion + "-" + a.Registry.Address(),
			JobType: "offchainreporting2",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "ocr2automation",
				ContractID: a.Registry.Address(),
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": int(a.ChainClient.GetChainID().Int64()),
				},
				PluginConfig: map[string]interface{}{
					"mercuryCredentialName": "\"" + a.MercuryCredentialName + "\"",
					"contractVersion":       "\"" + contractVersion + "\"",
				},
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
	var offC []byte
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
		offC, err = json.Marshal(ocr2keepers20config.OffchainConfig{
			TargetProbability:    a.PluginConfig.TargetProbability,
			TargetInRounds:       a.PluginConfig.TargetInRounds,
			PerformLockoutWindow: a.PluginConfig.PerformLockoutWindow,
			GasLimitPerReport:    a.PluginConfig.GasLimitPerReport,
			GasOverheadPerUpkeep: a.PluginConfig.GasOverheadPerUpkeep,
			MinConfirmations:     a.PluginConfig.MinConfirmations,
			MaxUpkeepBatchSize:   a.PluginConfig.MaxUpkeepBatchSize,
		})
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to marshal plugin config"))
		}

		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = ocr2.ContractSetConfigArgsForTests(
			a.PublicConfig.DeltaProgress, a.PublicConfig.DeltaResend,
			a.PublicConfig.DeltaRound, a.PublicConfig.DeltaGrace,
			a.PublicConfig.DeltaStage, uint8(a.PublicConfig.RMax),
			S, oracleIdentities, offC,
			a.PublicConfig.MaxDurationQuery, a.PublicConfig.MaxDurationObservation,
			1200*time.Millisecond,
			a.PublicConfig.MaxDurationShouldAcceptAttestedReport,
			a.PublicConfig.MaxDurationShouldTransmitAcceptedReport,
			a.PublicConfig.F, a.PublicConfig.OnchainConfig,
		)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to build config args"))
		}

	case ethereum.RegistryVersion_2_1:
		offC, err = json.Marshal(ocr2keepers30config.OffchainConfig{
			TargetProbability:    a.PluginConfig.TargetProbability,
			TargetInRounds:       a.PluginConfig.TargetInRounds,
			PerformLockoutWindow: a.PluginConfig.PerformLockoutWindow,
			GasLimitPerReport:    a.PluginConfig.GasLimitPerReport,
			GasOverheadPerUpkeep: a.PluginConfig.GasOverheadPerUpkeep,
			MinConfirmations:     a.PluginConfig.MinConfirmations,
			MaxUpkeepBatchSize:   a.PluginConfig.MaxUpkeepBatchSize,
		})
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to marshal plugin config"))
		}

		signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err = ocr3.ContractSetConfigArgsForTests(
			a.PublicConfig.DeltaProgress, a.PublicConfig.DeltaResend, a.PublicConfig.DeltaInitial,
			a.PublicConfig.DeltaRound, a.PublicConfig.DeltaGrace, a.PublicConfig.DeltaCertifiedCommitRequest,
			a.PublicConfig.DeltaStage, a.PublicConfig.RMax,
			S, oracleIdentities, offC,
			a.PublicConfig.MaxDurationQuery, a.PublicConfig.MaxDurationObservation,
			a.PublicConfig.MaxDurationShouldAcceptAttestedReport,
			a.PublicConfig.MaxDurationShouldTransmitAcceptedReport,
			a.PublicConfig.F, a.PublicConfig.OnchainConfig,
		)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to build config args"))
		}
	default:
		return fmt.Errorf("v2.0 and v2.1 are the only supported versions")
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

	onchainConfig, err := a.RegistrySettings.EncodeOnChainConfig(a.Registrar.Address(), a.UpkeepPrivilegeManager)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to encode onchain config"))
	}

	ocrConfig := contracts.OCRv2Config{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	err = a.Registry.SetConfig(a.RegistrySettings, ocrConfig)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to set config on registry"))
	}
	return nil
}

func (a *AutomationTest) RegisterUpkeeps(upkeepConfigs []UpkeepConfig) ([]common.Hash, error) {
	var registrarABI *abi.ABI
	var err error
	var registrationRequest []byte
	registrationTxHashes := make([]common.Hash, 0)

	for _, upkeepConfig := range upkeepConfigs {
		switch a.RegistrySettings.RegistryVersion {
		case ethereum.RegistryVersion_2_0:
			registrarABI, err = keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.GetAbi()
			if err != nil {
				return nil, errors.Join(err, fmt.Errorf("failed to get registrar abi"))
			}
			registrationRequest, err = registrarABI.Pack(
				"register", upkeepConfig.UpkeepName, upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract, upkeepConfig.GasLimit, upkeepConfig.AdminAddress,
				upkeepConfig.CheckData,
				upkeepConfig.OffchainConfig, upkeepConfig.FundingAmount,
				common.HexToAddress(a.ChainClient.GetDefaultWallet().Address()))
			if err != nil {
				return nil, errors.Join(err, fmt.Errorf("failed to pack registrar request"))
			}
		case ethereum.RegistryVersion_2_1:
			registrarABI, err = automation_registrar_wrapper2_1.AutomationRegistrarMetaData.GetAbi()
			if err != nil {
				return nil, errors.Join(err, fmt.Errorf("failed to get registrar abi"))
			}
			registrationRequest, err = registrarABI.Pack(
				"register", upkeepConfig.UpkeepName, upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract, upkeepConfig.GasLimit, upkeepConfig.AdminAddress,
				upkeepConfig.TriggerType, upkeepConfig.CheckData, upkeepConfig.TriggerConfig,
				upkeepConfig.OffchainConfig, upkeepConfig.FundingAmount,
				common.HexToAddress(a.ChainClient.GetDefaultWallet().Address()))
			if err != nil {
				return nil, errors.Join(err, fmt.Errorf("failed to pack registrar request"))
			}
		default:
			return nil, fmt.Errorf("v2.0 and v2.1 are the only supported versions")
		}
		tx, err := a.LinkToken.TransferAndCall(a.Registrar.Address(), upkeepConfig.FundingAmount, registrationRequest)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to register upkeep"))
		}
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
	}
	return registrationTxHashes, nil
}

func (a *AutomationTest) ConfirmUpkeepsRegistered(registrationTxHashes []common.Hash) ([]*big.Int, error) {
	upkeepIds := make([]*big.Int, 0)
	for _, txHash := range registrationTxHashes {
		receipt, err := a.ChainClient.GetTxReceipt(txHash)
		if err != nil {
			return nil, errors.Join(err, fmt.Errorf("failed to confirm upkeep registration"))
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
			return nil, fmt.Errorf("failed to parse upkeep id from registration receipt")
		}
		upkeepIds = append(upkeepIds, upkeepId)
	}
	a.UpkeepIDs = upkeepIds
	return upkeepIds, nil
}

func (a *AutomationTest) AddJobsAndSetConfig(t *testing.T) {
	l := logging.GetTestLogger(t)
	err := a.AddBootstrapJob()
	require.NoError(t, err, "Error adding bootstrap job")
	err = a.AddAutomationJobs()
	require.NoError(t, err, "Error adding automation jobs")

	l.Debug().
		Interface("Plugin Config", a.PluginConfig).
		Interface("Public Config", a.PublicConfig).
		Interface("Registry Settings", a.RegistrySettings).
		Interface("Registrar Settings", a.RegistrarSettings).
		Msg("Configuring registry")
	err = a.SetConfigOnRegistry()
	require.NoError(t, err, "Error setting config on registry")
	l.Info().Str("Registry Address", a.Registry.Address()).Msg("Successfully setConfig on registry")
}

func (a *AutomationTest) SetupMercuryMock(t *testing.T, imposters []ctfTestEnv.KillgraveImposter) {
	if a.IsOnk8s {
		t.Error("mercury mock is not supported on k8s")
	}
	if a.DockerEnv == nil {
		t.Error("docker env is not set")
	}
	err := a.DockerEnv.MockAdapter.AddImposter(imposters)
	if err != nil {
		require.NoError(t, err, "Error adding mock imposter")
	}
}

func (a *AutomationTest) SetupAutomationDeployment(t *testing.T) {
	l := logging.GetTestLogger(t)
	err := a.CollectNodeDetails()
	require.NoError(t, err, "Error collecting node details")
	l.Info().Msg("Collected Node Details")
	l.Debug().Interface("Node Details", a.NodeDetails).Msg("Node Details")

	err = a.DeployLINK()
	require.NoError(t, err, "Error deploying link token contract")

	err = a.DeployEthLinkFeed()
	require.NoError(t, err, "Error deploying eth link feed contract")
	err = a.DeployGasFeed()
	require.NoError(t, err, "Error deploying gas feed contract")

	err = a.DeployTranscoder()
	require.NoError(t, err, "Error deploying transcoder contract")

	err = a.DeployRegistry()
	require.NoError(t, err, "Error deploying registry contract")
	err = a.DeployRegistrar()
	require.NoError(t, err, "Error deploying registrar contract")

	a.AddJobsAndSetConfig(t)
}

func (a *AutomationTest) LoadAutomationDeployment(t *testing.T, linkTokenAddress,
	ethLinkFeedAddress, gasFeedAddress, transcoderAddress, registryAddress, registrarAddress string) {
	l := logging.GetTestLogger(t)
	err := a.CollectNodeDetails()
	require.NoError(t, err, "Error collecting node details")
	l.Info().Msg("Collected Node Details")
	l.Debug().Interface("Node Details", a.NodeDetails).Msg("Node Details")

	err = a.LoadLINK(linkTokenAddress)
	require.NoError(t, err, "Error loading link token contract")

	err = a.LoadEthLinkFeed(ethLinkFeedAddress)
	require.NoError(t, err, "Error loading eth link feed contract")
	err = a.LoadEthGasFeed(gasFeedAddress)
	require.NoError(t, err, "Error loading gas feed contract")
	err = a.LoadTranscoder(transcoderAddress)
	require.NoError(t, err, "Error loading transcoder contract")
	err = a.LoadRegistry(registryAddress)
	require.NoError(t, err, "Error loading registry contract")
	err = a.LoadRegistrar(registrarAddress)
	require.NoError(t, err, "Error loading registrar contract")

	a.AddJobsAndSetConfig(t)

}
