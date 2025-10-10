package internal

import (
	"context"
	"encoding/binary"
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
	"github.com/xssnick/tonutils-go/address"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

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
	pluginTypes []types.PluginType,
) ([]offramp.MultiOCR3BaseOCRConfigArgs, error) {
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
	var offrampOCR3Configs []offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range pluginTypes {
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

// we can't use the EVM one because we need the 32 byte transmitter address
type MultiOCR3BaseOCRConfigArgsSui struct {
	ConfigDigest                   [32]byte
	OcrPluginType                  byte
	F                              byte
	IsSignatureVerificationEnabled bool
	Signers                        [][]byte
	Transmitters                   []string
}

// BuildSetOCR3ConfigArgsSolana builds OCR3 config for Aptos chains
func BuildSetOCR3ConfigArgsSui(
	donID uint32,
	ccipHome *ccip_home.CCIPHome,
	destSelector uint64,
	configType globals.ConfigType,
) ([]MultiOCR3BaseOCRConfigArgsSui, error) {
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
	var offrampOCR3Configs []MultiOCR3BaseOCRConfigArgsSui
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
		var transmitterAddresses []string
		for _, node := range configForOCR3.Config.Nodes {
			signerAddresses = append(signerAddresses, node.SignerKey)

			transmitterAddress := "0x" + hex.EncodeToString(node.TransmitterKey)
			transmitterAddresses = append(transmitterAddresses, transmitterAddress)
		}

		offrampOCR3Configs = append(offrampOCR3Configs, MultiOCR3BaseOCRConfigArgsSui{
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
	pluginTypes []types.PluginType,
) ([]MultiOCR3BaseOCRConfigArgsSolana, error) {
	chainCfg, err := ccipHome.GetChainConfig(nil, destSelector)
	if err != nil {
		return nil, fmt.Errorf("error getting chain config for chain selector %d it must be set before OCR3Config set up: %w", destSelector, err)
	}
	ocr3Configs := make([]MultiOCR3BaseOCRConfigArgsSolana, 0)
	for _, pluginType := range pluginTypes {
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
	ocrSecrets focr.OCRSecrets,
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

		if destFamily == chain_selectors.FamilyAptos {
			transmitAccount = replaceAptosPublicKeys(transmitAccount)
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
			case chain_selectors.FamilySui:
				parsed, err = hex.DecodeString(strings.TrimPrefix(string(transmitter), "0x"))
				if err != nil {
					return nil, fmt.Errorf("failed to decode SUI address '%s': %w", transmitter, err)
				}
			case chain_selectors.FamilyTon:
				pk := address.MustParseAddr(string(transmitter))
				if pk == nil || pk.IsAddrNone() {
					return nil, fmt.Errorf("failed to parse TON address '%s'", transmitter)
				}
				// TODO: this reimplements addrCodec's ToRawAddr helper
				parsed = binary.BigEndian.AppendUint32(nil, uint32(pk.Workchain())) //nolint:gosec // G115
				parsed = append(parsed, pk.Data()...)
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

func replaceAptosPublicKeys(transmitterKey ocrtypes.Account) ocrtypes.Account {
	// Due to missing support in Operator UI, nodes will currently submit their account address to JD, when we actually need their pubkeys
	// As a temporary fix, hardcoding a mapping of address->pubkey, until we've got this fixed with https://github.com/smartcontractkit/operator-ui/pull/105
	// at which point we can remove this mapping and use the AccountAddressPublicKey returned by JD directly
	if pubkey, ok := map[ocrtypes.Account]ocrtypes.Account{
		// address -> pubkey
		"abf583b6a78104352571ac4c26c0c21cc300ac4721f20e7ccafe8fef0165994c": "df0abd3db66ac53143414283d80a989ab7952977a1829eb411a1866ec59b2e63",
		"b5fd635b553acd5edf10d5be495200d0e25dbcc49f533ccac0b3ae2a9002fc61": "49078b68242df84bb5934384c54fd2bedb781aea9d51b39b54d0a6a4e5b815d1",
		"e5c2b78461dbc0a65ec91acc7ebcd58c5174c76f260723ebbbf0e8b27b996138": "442016f35c766196688c8c334d559e319cc0f6eba514df8809d552c6c02d2410",
		"fc569c01021cdfb8e6e2aebf24a13c95f81e8d93f45380e0e584e5302a9c0700": "d2d0c410670c170a50ce6f52feb0eb579a6c0a8d8985b0097c5f192ffc204278",
		"1f396f4d814dde49e30735cc60ad3401a970b431984013ee8cc642af12780237": "412aef9e285a7226e9c54229595bc5124cc97fd4a7b744e3eb9be335ced5b3bf",
		"ef29451782a5a02f06588e8d5def7808c8160f381e7c1bfec6811305c70ea118": "0d82ef11f54b1ccf468e945eba8b05618343e015827f8fd36f3c425aeacb5c20",
		"7cd085db836299c483dd8ba74de6026d157dea13fd723fdfb1edd480129004ec": "21a6f41c79831274d5aad2eeb6fbf92d6c8fd91b3ef5de778f978a21f4d353fb",
		"9f967c2c31bccc06f679431f02c0134ff42aad339eb6257263aee72c707078c8": "54c9b0677196c0c409147023f63e21bc249ec86e34e3d31a84e14ec70f1689d8",
		"fc36546994ee99285dbb8a505ef2e1b5e1e65e583679c15b1ed8dde84be1a141": "f86e882182325f169677c07305833dbfdfc47c782acdad4378d0e9ae54197c26",
		"1e3308f5344a3da5e635c6ccbc1a384edb48615d68b1be0ee7531f3a35e53c87": "261eabc31932ba74ba6aa36377010faff4944e4cee06dbd17f9b4b7509dd0db8",
		"c3032bbe90b8d87444fed7bc1271287e9117f8dd043492a179371d712750d680": "e4933e51e556b1216bf2746add853a88dbee1031a8991e32b70a44ab64a83199",
		"a9f1212c6ca3fbc5ce842e807545670d95351b0c25dd7b05ed0e24662152b10d": "1ada224320f7e43fc45d5534fcaccc3fb1c88c364bfa9d36f72228b72069a28e",
		"abcc084a1a0719e219dc75203f6fc3a2e76a43eec9f8915f4481b7d89a8e47e3": "3b23c5703bbf86cc8620649f9f4f88890bd9098f4d8de951bb3c1897a883d7f2",
		"5fa9348f0d5d7410832e08296f82bab8c5e6b7705d5c97574f2c5e3fd5b185e2": "05bdcbed6359af1b47db1e17a2ace65288021658b72e4050412802b00474316c",
		"b0b6f9d21c73141fae33b8d889ac8eccb5e55b0ff5391a0e64c94e8b2d3370b0": "de3155f797d22fdbc8300bd4793b1f4dcfef8eeb581dbc76d426998d54396c0a",
		"1b601e9727eb8a4feed06c578e8cdc05405ec849b7259c5ec107dfdc1250aa9e": "87bfbb26a7b01c6802152a59a682cb42f0a1fa9844963ef15bd57fd7b8257189",
	}[transmitterKey]; ok {
		return pubkey
	}
	return transmitterKey
}
