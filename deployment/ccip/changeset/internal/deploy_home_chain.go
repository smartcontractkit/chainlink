package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"
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
	configType globals.ConfigType,
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

		configForOCR3 := ocrConfig.ActiveConfig
		// we expect only an active config
		if configType == globals.ConfigTypeActive {
			if ocrConfig.ActiveConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("invalid OCR3 config state, expected active config, donID: %d, activeConfig: %v, candidateConfig: %v",
					donID, hexutil.Encode(ocrConfig.ActiveConfig.ConfigDigest[:]), hexutil.Encode(ocrConfig.CandidateConfig.ConfigDigest[:]))
			}
		} else if configType == globals.ConfigTypeCandidate {
			if ocrConfig.CandidateConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("invalid OCR3 config state, expected candidate config, donID: %d, activeConfig: %v, candidateConfig: %v",
					donID, hexutil.Encode(ocrConfig.ActiveConfig.ConfigDigest[:]), hexutil.Encode(ocrConfig.CandidateConfig.ConfigDigest[:]))
			}
			configForOCR3 = ocrConfig.CandidateConfig
		}

		var signerAddresses []common.Address
		var transmitterAddresses []common.Address
		for _, node := range configForOCR3.Config.Nodes {
			signerAddresses = append(signerAddresses, common.BytesToAddress(node.SignerKey))
			transmitterAddresses = append(transmitterAddresses, common.BytesToAddress(node.TransmitterKey))
		}

		offrampOCR3Configs = append(offrampOCR3Configs, offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   configForOCR3.ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              configForOCR3.Config.FRoleDON,
			IsSignatureVerificationEnabled: pluginType == types.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return offrampOCR3Configs, nil
}

// https://github.com/smartcontractkit/chainlink-ccip/blob/bdbfcc588847d70817333487a9883e94c39a332e/chains/solana/gobindings/ccip_router/SetOcrConfig.go#L23
type MultiOCR3BaseOCRConfigArgsSolana struct {
	ConfigDigest                   [32]byte
	OCRPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        [][20]byte
	Transmitters                   []solana.PublicKey
}

// BuildSetOCR3ConfigArgsSolana builds OCR3 config for Solana chains
func BuildSetOCR3ConfigArgsSolana(
	donID uint32,
	ccipHome *ccip_home.CCIPHome,
	destSelector uint64,
) ([]MultiOCR3BaseOCRConfigArgsSolana, error) {
	ocr3Configs := make([]MultiOCR3BaseOCRConfigArgsSolana, 0)
	for _, pluginType := range []types.PluginType{types.PluginTypeCCIPCommit, types.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipHome.GetAllConfigs(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return nil, err2
		}

		// we expect only an active config and no candidate config.
		if ocrConfig.ActiveConfig.ConfigDigest == [32]byte{} || ocrConfig.CandidateConfig.ConfigDigest != [32]byte{} {
			return nil, fmt.Errorf("invalid OCR3 config state, expected active config and no candidate config, donID: %d", donID)
		}

		activeConfig := ocrConfig.ActiveConfig
		var signerAddresses [][20]byte
		var transmitterAddresses []solana.PublicKey
		for _, node := range activeConfig.Config.Nodes {
			var signer [20]uint8
			if len(node.SignerKey) != 20 {
				return nil, fmt.Errorf("node signer key not 20 bytes long, got: %d", len(node.SignerKey))
			}
			copy(signer[:], node.SignerKey)
			signerAddresses = append(signerAddresses, signer)
			// https://smartcontract-it.atlassian.net/browse/NONEVM-1254
			key, err := solana.PublicKeyFromBase58(string(node.TransmitterKey))
			if err != nil {
				return nil, err
			}
			transmitterAddresses = append(transmitterAddresses, key)
		}

		ocr3Configs = append(ocr3Configs, MultiOCR3BaseOCRConfigArgsSolana{
			ConfigDigest:                   activeConfig.ConfigDigest,
			OCRPluginType:                  uint8(pluginType),
			F:                              activeConfig.Config.FRoleDON,
			IsSignatureVerificationEnabled: pluginType == types.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return ocr3Configs, nil
}

func BuildOCR3ConfigForCCIPHome(
	ocrSecrets deployment.OCRSecrets,
	offRampAddress []byte,
	destSelector uint64,
	nodes deployment.Nodes,
	rmnHomeAddress common.Address,
	ocrParams commontypes.OCRParameters,
	commitOffchainCfg *pluginconfig.CommitOffchainConfig,
	execOffchainCfg *pluginconfig.ExecuteOffchainConfig,
) (map[types.PluginType]ccip_home.CCIPHomeOCR3Config, error) {
	p2pIDs := nodes.PeerIDs()
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)
		cfg, exists := node.OCRConfigForChainSelector(destSelector)
		if !exists {
			return nil, fmt.Errorf("no OCR config for chain %d", destSelector)
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
				MerkleRootAsyncObserverSyncFreq:    commitOffchainCfg.MerkleRootAsyncObserverSyncFreq,
				MerkleRootAsyncObserverSyncTimeout: commitOffchainCfg.MerkleRootAsyncObserverSyncTimeout,
				InflightPriceCheckRetries:          commitOffchainCfg.InflightPriceCheckRetries,
				TransmissionDelayMultiplier:        commitOffchainCfg.TransmissionDelayMultiplier,
				FeeInfo:                            commitOffchainCfg.FeeInfo,
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
			ChainSelector:         destSelector,
			FRoleDON:              configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        offRampAddress,
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
