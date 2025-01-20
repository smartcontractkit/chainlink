package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink/deployment"
	types2 "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"

	FirstBlockAge                           = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency       = 30 * time.Minute
	TokenPriceBatchWriteFrequency           = 30 * time.Minute
	BatchGasLimit                           = 6_500_000
	RelativeBoostPerWaitHour                = 0.5
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

	GasPriceDeviationPPB    = 1000
	DAGasPriceDeviationPPB  = 0
	OptimisticConfirmations = 1
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
		return 0, fmt.Errorf("get Dons from capability registry: %w", err)
	}
	var donIDs []uint32
	for _, don := range dons {
		if len(don.CapabilityConfigurations) == 1 &&
			don.CapabilityConfigurations[0].CapabilityId == CCIPCapabilityID {
			configs, err := ccipHome.GetAllConfigs(nil, don.Id, uint8(types.PluginTypeCCIPCommit))
			if err != nil {
				return 0, fmt.Errorf("get all commit configs from cciphome: %w", err)
			}
			if configs.ActiveConfig.ConfigDigest == [32]byte{} && configs.CandidateConfig.ConfigDigest == [32]byte{} {
				configs, err = ccipHome.GetAllConfigs(nil, don.Id, uint8(types.PluginTypeCCIPExec))
				if err != nil {
					return 0, fmt.Errorf("get all exec configs from cciphome: %w", err)
				}
			}
			if configs.ActiveConfig.Config.ChainSelector == chainSelector || configs.CandidateConfig.Config.ChainSelector == chainSelector {
				donIDs = append(donIDs, don.Id)
			}
		}
	}

	// more than one DON is an error
	if len(donIDs) > 1 {
		return 0, fmt.Errorf("more than one DON found for (chain selector %d, ccip capability id %x) pair", chainSelector, CCIPCapabilityID[:])
	}

	// no DON found - don ID of 0 indicates that (this is the case in the CR as well).
	if len(donIDs) == 0 {
		return 0, nil
	}

	// DON found - return it.
	return donIDs[0], nil
}

// BuildSetOCR3ConfigArgs builds the OCR3 config arguments for the OffRamp contract
// using the donID's OCR3 configs from the CCIPHome contract.
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
			return nil, fmt.Errorf("invalid OCR3 config state, expected active config and no candidate config, donID: %d, activeConfig: %v, candidateConfig: %v",
				donID, hexutil.Encode(ocrConfig.ActiveConfig.ConfigDigest[:]), hexutil.Encode(ocrConfig.CandidateConfig.ConfigDigest[:]))
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

func BuildOCR3ConfigForCCIPHome(
	ocrSecrets deployment.OCRSecrets,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	nodes deployment.Nodes,
	rmnHomeAddress common.Address,
	ocrParams types2.OCRParameters,
	commitOffchainCfg *pluginconfig.CommitOffchainConfig,
	execOffchainCfg *pluginconfig.ExecuteOffchainConfig,
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
	pluginTypes := make([]types.PluginType, 0)
	if commitOffchainCfg != nil {
		pluginTypes = append(pluginTypes, types.PluginTypeCCIPCommit)
	}
	if execOffchainCfg != nil {
		pluginTypes = append(pluginTypes, types.PluginTypeCCIPExec)
	}
	for _, pluginType := range pluginTypes {
		var encodedOffchainConfig []byte
		var err2 error
		if pluginType == types.PluginTypeCCIPCommit {
			if commitOffchainCfg == nil {
				return nil, errors.New("commitOffchainCfg is nil")
			}
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency:  commitOffchainCfg.RemoteGasPriceBatchWriteFrequency,
				TokenPriceBatchWriteFrequency:      commitOffchainCfg.TokenPriceBatchWriteFrequency,
				PriceFeedChainSelector:             commitOffchainCfg.PriceFeedChainSelector,
				TokenInfo:                          commitOffchainCfg.TokenInfo,
				NewMsgScanBatchSize:                commitOffchainCfg.NewMsgScanBatchSize,
				MaxReportTransmissionCheckAttempts: commitOffchainCfg.MaxReportTransmissionCheckAttempts,
				MaxMerkleTreeSize:                  commitOffchainCfg.MaxMerkleTreeSize,
				SignObservationPrefix:              commitOffchainCfg.SignObservationPrefix,
				RMNEnabled:                         commitOffchainCfg.RMNEnabled,
				RMNSignaturesTimeout:               commitOffchainCfg.RMNSignaturesTimeout,
			})
		} else {
			if execOffchainCfg == nil {
				return nil, errors.New("execOffchainCfg is nil")
			}
			encodedOffchainConfig, err2 = pluginconfig.EncodeExecuteOffchainConfig(pluginconfig.ExecuteOffchainConfig{
				BatchGasLimit:             execOffchainCfg.BatchGasLimit,
				RelativeBoostPerWaitHour:  execOffchainCfg.RelativeBoostPerWaitHour,
				MessageVisibilityInterval: execOffchainCfg.MessageVisibilityInterval,
				InflightCacheExpiry:       execOffchainCfg.InflightCacheExpiry,
				RootSnoozeTime:            execOffchainCfg.RootSnoozeTime,
				BatchingStrategyID:        execOffchainCfg.BatchingStrategyID,
				TokenDataObservers:        execOffchainCfg.TokenDataObservers,
			})
		}
		if err2 != nil {
			return nil, err2
		}
		signers, transmitters, configF, onchainConfig, offchainConfigVersion, offchainConfig, err2 := ocr3confighelper.ContractSetConfigArgsDeterministic(
			ocrSecrets.EphemeralSk,
			ocrSecrets.SharedSecret,
			ocrParams.DeltaProgress,
			ocrParams.DeltaResend,
			ocrParams.DeltaInitial,
			ocrParams.DeltaRound,
			ocrParams.DeltaGrace,
			ocrParams.DeltaCertifiedCommitRequest,
			ocrParams.DeltaStage,
			ocrParams.Rmax,
			schedule,
			oracles,
			encodedOffchainConfig,
			nil, // maxDurationInitialization
			ocrParams.MaxDurationQuery,
			ocrParams.MaxDurationObservation,
			ocrParams.MaxDurationShouldAcceptAttestedReport,
			ocrParams.MaxDurationShouldTransmitAcceptedReport,
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
		// validate ocr3 params correctness
		_, err := ocr3confighelper.PublicConfigFromContractConfig(false, ocrtypes.ContractConfig{
			Signers:               signers,
			Transmitters:          transmitters,
			F:                     configF,
			OnchainConfig:         onchainConfig,
			OffchainConfigVersion: offchainConfigVersion,
			OffchainConfig:        offchainConfig,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to validate ocr3 params: %w", err)
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

func DONIdExists(cr *capabilities_registry.CapabilitiesRegistry, donIDs []uint32) error {
	// DON ids must exist
	dons, err := cr.GetDONs(nil)
	if err != nil {
		return fmt.Errorf("failed to get dons: %w", err)
	}
	for _, donID := range donIDs {
		exists := false
		for _, don := range dons {
			if don.Id == donID {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf("don id %d does not exist", donID)
		}
	}
	return nil
}
