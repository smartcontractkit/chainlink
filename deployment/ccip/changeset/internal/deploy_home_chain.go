package internal

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

var (
	CCIPHomeABI *abi.ABI
)

func init() {
	var err error
	CCIPHomeABI, err = ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
}

// LatestCCIPDON returns the latest CCIP DON from the capabilities registry
// Keeping this function for reference
func LatestCCIPDON(registry *capabilities_registry.CapabilitiesRegistry) (*capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	dons, err := registry.GetDONs(nil)
	if err != nil {
		return nil, err
	}
	var ccipDON capabilities_registry.CapabilitiesRegistryDONInfo
	for _, don := range dons {
		if len(don.CapabilityConfigurations) == 1 &&
			don.CapabilityConfigurations[0].CapabilityId == shared.CCIPCapabilityID &&
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
			don.CapabilityConfigurations[0].CapabilityId == shared.CCIPCapabilityID {
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
		return 0, fmt.Errorf("more than one DON found for (chain selector %d, ccip capability id %x) pair", chainSelector, shared.CCIPCapabilityID[:])
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
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
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
		if err := validateOCR3Config(destSelector, configForOCR3.Config, &chainCfg); err != nil {
			return nil, err
		}

		var signerAddresses []common.Address
		var transmitterAddresses []common.Address
		for _, node := range configForOCR3.Config.Nodes {
			signerAddresses = append(signerAddresses, common.BytesToAddress(node.SignerKey))

			// Not all nodes support the destination chain, if the transmitter key is empty in the CCIPHome OCR3 config,
			// it means that we can omit it from the transmitter whitelist on the OCR3 contract
			// on the destination chain.
			if len(node.TransmitterKey) > 0 {
				transmitterAddresses = append(transmitterAddresses, common.BytesToAddress(node.TransmitterKey))
			}
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

func validateOCR3Config(chainSel uint64, configForOCR3 ccip_home.CCIPHomeOCR3Config, chainConfig *ccip_home.CCIPHomeChainConfig) error {
	if chainConfig != nil {
		// chainConfigs must be set before OCR3 configs due to the added fChain == F validation
		if chainConfig.FChain == 0 || bytes.IsEmpty(chainConfig.Config) || len(chainConfig.Readers) == 0 {
			return fmt.Errorf("chain config is not set for chain selector %d", chainSel)
		}
		for _, reader := range chainConfig.Readers {
			if bytes.IsEmpty(reader[:]) {
				return fmt.Errorf("reader is empty, chain selector %d", chainSel)
			}
		}
		// FRoleDON >= fChain is a requirement
		if configForOCR3.FRoleDON < chainConfig.FChain {
			return fmt.Errorf("OCR3 config FRoleDON is lower than chainConfig FChain, chain %d", chainSel)
		}

		if len(configForOCR3.Nodes) < 3*int(chainConfig.FChain)+1 {
			return fmt.Errorf("number of nodes %d is less than 3 * fChain + 1 %d", len(configForOCR3.Nodes), 3*int(chainConfig.FChain)+1)
		}

		// check that we have enough transmitters for the destination chain.
		// note that this is done onchain, but we'll do it here for good measure to avoid reverts.
		// see https://github.com/smartcontractkit/chainlink-ccip/blob/8529b8c89093d0cd117b73645ea64b2d2a8092f4/chains/evm/contracts/capability/CCIPHome.sol#L511-L514.
		minTransmitterReq := 3*int(chainConfig.FChain) + 1
		var numNonzeroTransmitters int
		for _, node := range configForOCR3.Nodes {
			if len(node.TransmitterKey) > 0 {
				numNonzeroTransmitters++
			}
		}
		if numNonzeroTransmitters < minTransmitterReq {
			return fmt.Errorf("number of transmitters (%d) is less than 3 * fChain + 1 (%d), chain selector %d",
				numNonzeroTransmitters, minTransmitterReq, chainSel)
		}
	}

	// check if there is any zero byte address
	// The reason for this is that the MultiOCR3Base disallows zero addresses and duplicates
	if bytes.IsEmpty(configForOCR3.OfframpAddress) {
		return fmt.Errorf("zero address found in offramp address,  chain %d", chainSel)
	}
	if bytes.IsEmpty(configForOCR3.RmnHomeAddress) {
		return fmt.Errorf("zero address found in rmn home address,  chain %d", chainSel)
	}
	mapSignerKey := make(map[string]struct{})
	mapTransmitterKey := make(map[string]struct{})
	for _, node := range configForOCR3.Nodes {
		if bytes.IsEmpty(node.SignerKey) {
			return fmt.Errorf("zero address found in signer key, chain %d", chainSel)
		}

		// NOTE: We don't check for empty/zero transmitter address because the node can have a zero transmitter address if it does not support the destination chain.

		if bytes.IsEmpty(node.P2pId[:]) {
			return fmt.Errorf("empty p2p id, chain %d", chainSel)
		}

		// Signer and non-zero transmitter duplication must be checked
		if _, ok := mapSignerKey[hexutil.Encode(node.SignerKey)]; ok {
			return fmt.Errorf("duplicate signer key found, chain %d", chainSel)
		}

		// If len(node.TransmitterKey) == 0, the node does not support the destination chain, and we can definitely
		// have more than one node not supporting the destination chain, so the duplicate check doesn't make sense
		// for those.
		if _, ok := mapTransmitterKey[hexutil.Encode(node.TransmitterKey)]; ok && len(node.TransmitterKey) != 0 {
			return fmt.Errorf("duplicate transmitter key found, chain %d", chainSel)
		}
		mapSignerKey[hexutil.Encode(node.SignerKey)] = struct{}{}
		mapTransmitterKey[hexutil.Encode(node.TransmitterKey)] = struct{}{}
	}
	return nil
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
	configType globals.ConfigType,
) ([]MultiOCR3BaseOCRConfigArgsSolana, error) {
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
	ocr3Configs := make([]MultiOCR3BaseOCRConfigArgsSolana, 0)
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
		if err := validateOCR3Config(destSelector, configForOCR3.Config, &chainCfg); err != nil {
			return nil, err
		}

		var signerAddresses [][20]byte
		var transmitterAddresses []solana.PublicKey
		for _, node := range configForOCR3.Config.Nodes {
			var signer [20]uint8
			if len(node.SignerKey) != 20 {
				return nil, fmt.Errorf("node signer key not 20 bytes long, got: %d", len(node.SignerKey))
			}
			copy(signer[:], node.SignerKey)
			signerAddresses = append(signerAddresses, signer)
			key := solana.PublicKeyFromBytes(node.TransmitterKey)
			transmitterAddresses = append(transmitterAddresses, key)
		}

		ocr3Configs = append(ocr3Configs, MultiOCR3BaseOCRConfigArgsSolana{
			ConfigDigest:                   configForOCR3.ConfigDigest,
			OCRPluginType:                  uint8(pluginType),
			F:                              configForOCR3.Config.FRoleDON,
			IsSignatureVerificationEnabled: pluginType == types.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return ocr3Configs, nil
}

// we can't use the EVM one because we need the 32 byte transmitter address
type MultiOCR3BaseOCRConfigArgsAptos struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  uint8
	F                              uint8
	IsSignatureVerificationEnabled bool
	Signers                        [][]byte
	Transmitters                   [][]byte
}

// BuildSetOCR3ConfigArgsSolana builds OCR3 config for Aptos chains
func BuildSetOCR3ConfigArgsAptos(
	donID uint32,
	ccipHome *ccip_home.CCIPHome,
	destSelector uint64,
	configType globals.ConfigType,
) ([]MultiOCR3BaseOCRConfigArgsAptos, error) {
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
	var offrampOCR3Configs []MultiOCR3BaseOCRConfigArgsAptos
	for _, pluginType := range []types.PluginType{types.PluginTypeCCIPCommit, types.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipHome.GetAllConfigs(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return nil, err2
		}

		configForOCR3 := ocrConfig.ActiveConfig
		// we expect only an active config
		switch configType {
		case globals.ConfigTypeActive:
			if ocrConfig.ActiveConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("invalid OCR3 config state, expected active config, donID: %d, activeConfig: %v, candidateConfig: %v",
					donID, hexutil.Encode(ocrConfig.ActiveConfig.ConfigDigest[:]), hexutil.Encode(ocrConfig.CandidateConfig.ConfigDigest[:]))
			}
		case globals.ConfigTypeCandidate:
			if ocrConfig.CandidateConfig.ConfigDigest == [32]byte{} {
				return nil, fmt.Errorf("invalid OCR3 config state, expected candidate config, donID: %d, activeConfig: %v, candidateConfig: %v",
					donID, hexutil.Encode(ocrConfig.ActiveConfig.ConfigDigest[:]), hexutil.Encode(ocrConfig.CandidateConfig.ConfigDigest[:]))
			}
			configForOCR3 = ocrConfig.CandidateConfig
		}

		if err := validateOCR3Config(destSelector, configForOCR3.Config, &chainCfg); err != nil {
			return nil, err
		}

		var signerAddresses [][]byte
		var transmitterAddresses [][]byte
		for _, node := range configForOCR3.Config.Nodes {
			signerAddresses = append(signerAddresses, node.SignerKey)
			transmitterAddresses = append(transmitterAddresses, node.TransmitterKey)
		}

		offrampOCR3Configs = append(offrampOCR3Configs, MultiOCR3BaseOCRConfigArgsAptos{
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

func BuildOCR3ConfigForCCIPHome(
	ccipHome *ccip_home.CCIPHome,
	ocrSecrets cldf.OCRSecrets,
	offRampAddress []byte,
	destSelector uint64,
	nodes deployment.Nodes,
	rmnHomeAddress common.Address,
	ocrParams commontypes.OCRParameters,
	commitOffchainCfg *pluginconfig.CommitOffchainConfig,
	execOffchainCfg *pluginconfig.ExecuteOffchainConfig,
	skipChainConfigValidation bool,
) (map[types.PluginType]ccip_home.CCIPHomeOCR3Config, error) {
	addressCodec := ccipcommon.NewAddressCodec(map[string]ccipcommon.ChainSpecificAddressCodec{
		chain_selectors.FamilyEVM:    ccipevm.AddressCodec{},
		chain_selectors.FamilySolana: ccipsolana.AddressCodec{},
	})

	// check if we have info from this node for another chain in the same destFamily
	destFamily, err := chain_selectors.GetSelectorFamily(destSelector)
	if err != nil {
		return nil, err
	}

	var p2pIDs [][32]byte
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)

		// TODO: not every node supports the destination chain, but nodes must have an OCR identity for the
		// destination chain, in order to be able to participate in the OCR protocol, sign reports, etc.
		// However, JD currently only returns the "OCRConfig" for chains that are explicitly supported by the node,
		// presumably in the TOML config.
		// JD should instead give us the OCR identity for the destination chain, and, if the node does NOT
		// actually support the chain (in terms of TOML config), then return an empty transmitter address,
		// which is what we're supposed to set anyway if that particular node doesn't support the destination chain.
		// The current workaround is to check if we have the OCR identity for the destination chain based off of
		// the node's OCR identity for another chain in the same family.
		// This is a HACK, because it is entirely possible that the destination chain is a unique family,
		// and no other supported chain by the node has the same family, e.g. Solana.
		cfg, exists := node.OCRConfigForChainSelector(destSelector)
		if !exists {
			// check if we have an oracle identity for another chain in the same family as destFamily.
			allOCRConfigs := node.AllOCRConfigs()
			for chainDetails, ocrConfig := range allOCRConfigs {
				chainFamily, err := chain_selectors.GetSelectorFamily(chainDetails.ChainSelector)
				if err != nil {
					return nil, err
				}

				if chainFamily == destFamily {
					cfg = ocrConfig
					break
				}
			}

			if cfg.OffchainPublicKey == [32]byte{} {
				return nil, fmt.Errorf(
					"no OCR config for chain %d (family %s) from node %s (peer id %s) and no other OCR config for another chain in the same family",
					destSelector, destFamily, node.Name, node.PeerID.String(),
				)
			}
		}

		var transmitAccount ocrtypes.Account
		if !exists {
			// empty account means that the node cannot transmit for this chain
			// we replace this with a canonical address with the oracle ID as the address when doing the ocr config validation below, but it should remain empty
			// in the CCIPHome OCR config and it should not be included in the destination chain transmitters whitelist.
			transmitAccount = ocrtypes.Account("")
		} else {
			transmitAccount = cfg.TransmitAccount
		}
		p2pIDs = append(p2pIDs, node.PeerID)
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  cfg.OnchainPublicKey,    // should be the same for all chains within the same family
				TransmitAccount:   transmitAccount,         // different per chain (!) can be empty if the node does not support the destination chain
				OffchainPublicKey: cfg.OffchainPublicKey,   // should be the same for all chains within the same family
				PeerID:            cfg.PeerID.String()[4:], // should be the same for all oracle identities
			},
			ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey, // should be the same for all chains within the same family
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
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(*commitOffchainCfg)
		} else {
			if execOffchainCfg == nil {
				return nil, errors.New("execOffchainCfg is nil")
			}
			encodedOffchainConfig, err2 = pluginconfig.EncodeExecuteOffchainConfig(*execOffchainCfg)
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
			// TODO: this should just use the addresscodec
			family, err := chain_selectors.GetSelectorFamily(destSelector)
			if err != nil {
				return nil, err
			}
			var parsed []byte

			// if the node does not support the destination chain, the transmitter address is empty.
			if len(transmitter) == 0 {
				transmittersBytes[i] = []byte{}
				continue
			}

			switch family {
			case chain_selectors.FamilyEVM:
				parsed, err2 = common.ParseHexOrString(string(transmitter))
				if err2 != nil {
					return nil, err2
				}
			case chain_selectors.FamilySolana:
				pk, err := solana.PublicKeyFromBase58(string(transmitter))
				if err != nil {
					return nil, fmt.Errorf("failed to decode SVM address '%s': %w", transmitter, err)
				}
				parsed = pk.Bytes()
			case chain_selectors.FamilyAptos:
				parsed, err = hex.DecodeString(strings.TrimPrefix(string(transmitter), "0x"))
				if err != nil {
					return nil, fmt.Errorf("failed to decode Aptos address '%s': %w", transmitter, err)
				}
			}

			transmittersBytes[i] = parsed
		}

		// validate ocr3 params correctness
		// TODO: this is super hacky, should not have to do this.
		transmitters, err := replaceEmptyTransmitters(transmitters, addressCodec, destSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to replace empty transmitters in transmitters list before validating ocr3 params: %w", err)
		}

		_, err = ocr3confighelper.PublicConfigFromContractConfig(false, ocrtypes.ContractConfig{
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

		if !skipChainConfigValidation {
			chainConfig, err := ccipHome.GetChainConfig(nil, destSelector)
			if err != nil {
				return nil, fmt.Errorf("can't get chain config for %d: %w", destSelector, err)
			}
			if err := validateOCR3Config(destSelector, ocr3Configs[pluginType], &chainConfig); err != nil {
				return nil, fmt.Errorf("failed to validate ocr3 config: %w", err)
			}
		}
	}

	return ocr3Configs, nil
}

// replaceEmptyTransmitters replaces empty transmitters with a canonical address, using the oracle ID as the address in order to pass OCR config validation.
// TODO: this is super hacky, should not have to do this.
func replaceEmptyTransmitters(transmitters []ocrtypes.Account, addressCodec ccipcommon.AddressCodec, destSelector uint64) ([]ocrtypes.Account, error) {
	var ret []ocrtypes.Account
	for oracleID, transmitter := range transmitters {
		acct := transmitter
		if len(acct) == 0 {
			// #nosec G115 - Overflow is not a concern in this test scenario
			canonicalAddress, err := addressCodec.OracleIDAsAddressBytes(uint8(oracleID), ccipocr3.ChainSelector(destSelector))
			if err != nil {
				return nil, err
			}

			acctString, err := addressCodec.AddressBytesToString(canonicalAddress, ccipocr3.ChainSelector(destSelector))
			if err != nil {
				return nil, err
			}

			acct = ocrtypes.Account(acctString)
		}
		ret = append(ret, acct)
	}

	return ret, nil
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
