package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"

	FirstBlockAge                           = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency       = 30 * time.Minute
	TokenPriceBatchWriteFrequency           = 30 * time.Minute
	BatchGasLimit                           = 6_500_000
	RelativeBoostPerWaitHour                = 1.5
	InflightCacheExpiry                     = 10 * time.Minute
	RootSnoozeTime                          = 30 * time.Minute
	BatchingStrategyID                      = 0
	DeltaProgress                           = 30 * time.Second
	DeltaResend                             = 10 * time.Second
	DeltaInitial                            = 20 * time.Second
	DeltaRound                              = 2 * time.Second
	DeltaGrace                              = 2 * time.Second
	DeltaCertifiedCommitRequest             = 10 * time.Second
	DeltaStage                              = 10 * time.Second
	Rmax                                    = 3
	MaxDurationQuery                        = 500 * time.Millisecond
	MaxDurationObservation                  = 5 * time.Second
	MaxDurationShouldAcceptAttestedReport   = 10 * time.Second
	MaxDurationShouldTransmitAcceptedReport = 10 * time.Second
)

var (
	CCIPCapabilityID = utils.Keccak256Fixed(MustABIEncode(`[{"type": "string"}, {"type": "string"}]`, CapabilityLabelledName, CapabilityVersion))
	CCIPHomeABI      *abi.ABI
)

func init() {
	var err error
	CCIPHomeABI, err = ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
}

func MustABIEncode(abiString string, args ...interface{}) []byte {
	encoded, err := utils.ABIEncode(abiString, args...)
	if err != nil {
		panic(err)
	}
	return encoded
}

// getNodeOperatorIDMap returns a map of node operator names to their IDs
// If maxNops is greater than the number of node operators, it will return all node operators
// Unused now but could be useful in the future.
func getNodeOperatorIDMap(capReg *capabilities_registry.CapabilitiesRegistry, maxNops uint32) (map[string]uint32, error) {
	nopIdByName := make(map[string]uint32)
	operators, err := capReg.GetNodeOperators(nil)
	if err != nil {
		return nil, err
	}
	if len(operators) < int(maxNops) {
		maxNops = uint32(len(operators))
	}
	for i := uint32(1); i <= maxNops; i++ {
		operator, err := capReg.GetNodeOperator(nil, i)
		if err != nil {
			return nil, err
		}
		nopIdByName[operator.Name] = i
	}
	return nopIdByName, nil
}

func LatestCCIPDON(registry *capabilities_registry.CapabilitiesRegistry) (*capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	dons, err := registry.GetDONs(nil)
	if err != nil {
		return nil, err
	}
	var ccipDON capabilities_registry.CapabilitiesRegistryDONInfo
	for _, don := range dons {
		if len(don.CapabilityConfigurations) == 1 &&
			don.CapabilityConfigurations[0].CapabilityId == CCIPCapabilityID &&
			don.Id > ccipDON.Id {
			ccipDON = don
		}
	}
	return &ccipDON, nil
}

// DonIDForChain returns the DON ID for the chain with the given selector
// It looks up with the CCIPHome contract to find the OCR3 configs for the DONs, and returns the DON ID for the chain matching with the given selector from the OCR3 configs
func DonIDForChain(registry *capabilities_registry.CapabilitiesRegistry, ccipHome *ccip_home.CCIPHome, chainSelector uint64) (uint32, error) {
	dons, err := registry.GetDONs(nil)
	if err != nil {
		return 0, err
	}
	// TODO: what happens if there are multiple dons for one chain (accidentally?)
	for _, don := range dons {
		if len(don.CapabilityConfigurations) == 1 &&
			don.CapabilityConfigurations[0].CapabilityId == CCIPCapabilityID {
			configs, err := ccipHome.GetAllConfigs(nil, don.Id, uint8(types.PluginTypeCCIPCommit))
			if err != nil {
				return 0, err
			}
			if configs.ActiveConfig.Config.ChainSelector == chainSelector || configs.CandidateConfig.Config.ChainSelector == chainSelector {
				return don.Id, nil
			}
		}
	}
	return 0, fmt.Errorf("no DON found for chain %d", chainSelector)
}

func BuildSetOCR3ConfigArgs(
	donID uint32,
	ccipHome *ccip_home.CCIPHome,
	destSelector uint64,
) ([]offramp.MultiOCR3BaseOCRConfigArgs, error) {
	var offrampOCR3Configs []offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range []types.PluginType{types.PluginTypeCCIPCommit, types.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipHome.GetAllConfigs(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return nil, err2
		}

		fmt.Printf("pluginType: %s, destSelector: %d, donID: %d, activeConfig digest: %x, candidateConfig digest: %x\n",
			pluginType.String(), destSelector, donID, ocrConfig.ActiveConfig.ConfigDigest, ocrConfig.CandidateConfig.ConfigDigest)

		// we expect only an active config and no candidate config.
		if ocrConfig.ActiveConfig.ConfigDigest == [32]byte{} || ocrConfig.CandidateConfig.ConfigDigest != [32]byte{} {
			return nil, fmt.Errorf("invalid OCR3 config state, expected active config and no candidate config, donID: %d", donID)
		}

		activeConfig := ocrConfig.ActiveConfig
		var signerAddresses []common.Address
		var transmitterAddresses []common.Address
		for _, node := range activeConfig.Config.Nodes {
			signerAddresses = append(signerAddresses, common.BytesToAddress(node.SignerKey))
			transmitterAddresses = append(transmitterAddresses, common.BytesToAddress(node.TransmitterKey))
		}

		offrampOCR3Configs = append(offrampOCR3Configs, offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   activeConfig.ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              activeConfig.Config.FRoleDON,
			IsSignatureVerificationEnabled: pluginType == types.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return offrampOCR3Configs, nil
}

func SetupExecDON(
	donID uint32,
	execConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	home deployment.Chain,
	nodes deployment.Nodes,
	ccipHome *ccip_home.CCIPHome,
) error {
	encodedSetCandidateCall, err := CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		execConfig.PluginType,
		execConfig,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack set candidate call: %w", err)
	}

	// set candidate call
	tx, err := capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ exec config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm update don w/ exec config: %w", err)
	}

	execCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, execConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get exec candidate digest 1st time: %w", err)
	}

	if execCandidateDigest == [32]byte{} {
		return fmt.Errorf("candidate digest is empty, expected nonempty")
	}

	// promote candidate call
	encodedPromotionCall, err := CCIPHomeABI.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		execConfig.PluginType,
		execCandidateDigest,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack promotion call: %w", err)
	}

	tx, err = capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ exec config: %w", err)
	}
	bn, err := deployment.ConfirmIfNoError(home, tx, err)
	if err != nil {
		return fmt.Errorf("confirm update don w/ exec config: %w", err)
	}
	if bn == 0 {
		return fmt.Errorf("UpdateDON tx not confirmed")
	}
	// check if candidate digest is promoted
	pEvent, err := ccipHome.FilterConfigPromoted(&bind.FilterOpts{
		Context: context.Background(),
		Start:   bn,
	}, [][32]byte{execCandidateDigest})
	if err != nil {
		return fmt.Errorf("filter exec config promoted: %w", err)
	}
	if !pEvent.Next() {
		return fmt.Errorf("exec config not promoted")
	}
	// check that candidate digest is empty.
	execCandidateDigest, err = ccipHome.GetCandidateDigest(nil, donID, execConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get exec candidate digest 2nd time: %w", err)
	}

	if execCandidateDigest != [32]byte{} {
		return fmt.Errorf("candidate digest is nonempty after promotion, expected empty")
	}

	// check that active digest is non-empty.
	execActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(types.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get active exec digest: %w", err)
	}

	if execActiveDigest == [32]byte{} {
		return fmt.Errorf("active exec digest is empty, expected nonempty")
	}

	execConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(types.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get all exec configs 2nd time: %w", err)
	}

	// print the above info
	fmt.Printf("completed exec DON creation and promotion: donID: %d execCandidateDigest: %x, execActiveDigest: %x, execCandidateDigestFromGetAllConfigs: %x, execActiveDigestFromGetAllConfigs: %x\n",
		donID, execCandidateDigest, execActiveDigest, execConfigs.CandidateConfig.ConfigDigest, execConfigs.ActiveConfig.ConfigDigest)

	return nil
}

func SetupCommitDON(
	donID uint32,
	commitConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	home deployment.Chain,
	nodes deployment.Nodes,
	ccipHome *ccip_home.CCIPHome,
) error {
	encodedSetCandidateCall, err := CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		commitConfig.PluginType,
		commitConfig,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack set candidate call: %w", err)
	}
	tx, err := capReg.AddDON(home.DeployerKey, nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: CCIPCapabilityID,
			Config:       encodedSetCandidateCall,
		},
	}, false, false, nodes.DefaultF())
	if err != nil {
		return fmt.Errorf("add don w/ commit config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm add don w/ commit config: %w", err)
	}

	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, commitConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}

	if commitCandidateDigest == [32]byte{} {
		return fmt.Errorf("candidate digest is empty, expected nonempty")
	}
	fmt.Printf("commit candidate digest after setCandidate: %x\n", commitCandidateDigest)

	encodedPromotionCall, err := CCIPHomeABI.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		commitConfig.PluginType,
		commitCandidateDigest,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack promotion call: %w", err)
	}

	tx, err = capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ commit config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm update don w/ commit config: %w", err)
	}

	// check that candidate digest is empty.
	commitCandidateDigest, err = ccipHome.GetCandidateDigest(nil, donID, commitConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get commit candidate digest 2nd time: %w", err)
	}

	if commitCandidateDigest != [32]byte{} {
		return fmt.Errorf("candidate digest is nonempty after promotion, expected empty")
	}

	// check that active digest is non-empty.
	commitActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(types.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get active commit digest: %w", err)
	}

	if commitActiveDigest == [32]byte{} {
		return fmt.Errorf("active commit digest is empty, expected nonempty")
	}

	commitConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(types.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get all commit configs 2nd time: %w", err)
	}

	// print the above information
	fmt.Printf("completed commit DON creation and promotion: donID: %d, commitCandidateDigest: %x, commitActiveDigest: %x, commitCandidateDigestFromGetAllConfigs: %x, commitActiveDigestFromGetAllConfigs: %x\n",
		donID, commitCandidateDigest, commitActiveDigest, commitConfigs.CandidateConfig.ConfigDigest, commitConfigs.ActiveConfig.ConfigDigest)

	return nil
}

func BuildOCR3ConfigForCCIPHome(
	ocrSecrets deployment.OCRSecrets,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	feedChainSel uint64,
	tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	nodes deployment.Nodes,
	rmnHomeAddress common.Address,
	configs []pluginconfig.TokenDataObserverConfig,
) (map[types.PluginType]ccip_home.CCIPHomeOCR3Config, error) {
	p2pIDs := nodes.PeerIDs()
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)
		cfg, exists := node.OCRConfigForChainSelector(dest.Selector)
		if !exists {
			return nil, fmt.Errorf("no OCR config for chain %d", dest.Selector)
		}
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  cfg.OnchainPublicKey,
				TransmitAccount:   cfg.TransmitAccount,
				OffchainPublicKey: cfg.OffchainPublicKey,
				PeerID:            cfg.PeerID.String()[4:],
			}, ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey,
		})
	}

	// Add DON on capability registry contract
	ocr3Configs := make(map[types.PluginType]ccip_home.CCIPHomeOCR3Config)
	for _, pluginType := range []types.PluginType{types.PluginTypeCCIPCommit, types.PluginTypeCCIPExec} {
		var encodedOffchainConfig []byte
		var err2 error
		if pluginType == types.PluginTypeCCIPCommit {
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(RemoteGasPriceBatchWriteFrequency),
				TokenPriceBatchWriteFrequency:      *config.MustNewDuration(TokenPriceBatchWriteFrequency),
				PriceFeedChainSelector:             ccipocr3.ChainSelector(feedChainSel),
				TokenInfo:                          tokenInfo,
				NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
				MaxReportTransmissionCheckAttempts: 5,
				MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
				SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
				RMNEnabled:                         os.Getenv("ENABLE_RMN") == "true", // only enabled in manual test
			})
		} else {
			encodedOffchainConfig, err2 = pluginconfig.EncodeExecuteOffchainConfig(pluginconfig.ExecuteOffchainConfig{
				BatchGasLimit:             BatchGasLimit,
				RelativeBoostPerWaitHour:  RelativeBoostPerWaitHour,
				MessageVisibilityInterval: *config.MustNewDuration(FirstBlockAge),
				InflightCacheExpiry:       *config.MustNewDuration(InflightCacheExpiry),
				RootSnoozeTime:            *config.MustNewDuration(RootSnoozeTime),
				BatchingStrategyID:        BatchingStrategyID,
				TokenDataObservers:        configs,
			})
		}
		if err2 != nil {
			return nil, err2
		}
		signers, transmitters, configF, _, offchainConfigVersion, offchainConfig, err2 := ocr3confighelper.ContractSetConfigArgsDeterministic(
			ocrSecrets.EphemeralSk,
			ocrSecrets.SharedSecret,
			DeltaProgress,
			DeltaResend,
			DeltaInitial,
			DeltaRound,
			DeltaGrace,
			DeltaCertifiedCommitRequest,
			DeltaStage,
			Rmax,
			schedule,
			oracles,
			encodedOffchainConfig,
			nil, // maxDurationInitialization
			MaxDurationQuery,
			MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport,
			int(nodes.DefaultF()),
			[]byte{}, // empty OnChainConfig
		)
		if err2 != nil {
			return nil, err2
		}

		signersBytes := make([][]byte, len(signers))
		for i, signer := range signers {
			signersBytes[i] = signer
		}

		transmittersBytes := make([][]byte, len(transmitters))
		for i, transmitter := range transmitters {
			parsed, err2 := common.ParseHexOrString(string(transmitter))
			if err2 != nil {
				return nil, err2
			}
			transmittersBytes[i] = parsed
		}

		var ocrNodes []ccip_home.CCIPHomeOCR3Node
		for i := range nodes {
			ocrNodes = append(ocrNodes, ccip_home.CCIPHomeOCR3Node{
				P2pId:          p2pIDs[i],
				SignerKey:      signersBytes[i],
				TransmitterKey: transmittersBytes[i],
			})
		}

		_, ok := ocr3Configs[pluginType]
		if ok {
			return nil, fmt.Errorf("pluginType %s already exists in ocr3Configs", pluginType.String())
		}

		ocr3Configs[pluginType] = ccip_home.CCIPHomeOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         dest.Selector,
			FRoleDON:              configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        offRamp.Address().Bytes(),
			Nodes:                 ocrNodes,
			OffchainConfig:        offchainConfig,
			RmnHomeAddress:        rmnHomeAddress.Bytes(),
		}
	}

	return ocr3Configs, nil
}
