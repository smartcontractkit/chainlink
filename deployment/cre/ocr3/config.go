package ocr3

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"google.golang.org/protobuf/proto"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

var (
	OCR3Capability cldf.ContractType = "OCR3Capability" // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
)

type TopLevelConfigSource struct {
	OracleConfig OracleConfig
}

type MCMSConfig struct {
	MinDuration time.Duration
}

type NodeKeys struct {
	EthAddress             string `json:"EthAddress"`
	AptosAccount           string `json:"AptosAccount"`
	AptosBundleID          string `json:"AptosBundleID"`
	AptosOnchainPublicKey  string `json:"AptosOnchainPublicKey"`
	SolanaOnchainPublicKey string `json:"SolanaOnchainPublicKey"`
	SolanaBundleID         string `json:"SolanaBundleID"`
	P2PPeerID              string `json:"P2PPeerID"`             // p2p_<key>
	OCR2BundleID           string `json:"OCR2BundleID"`          // used only in job spec
	OCR2OnchainPublicKey   string `json:"OCR2OnchainPublicKey"`  // ocr2on_evm_<key>
	OCR2OffchainPublicKey  string `json:"OCR2OffchainPublicKey"` // ocr2off_evm_<key>
	OCR2ConfigPublicKey    string `json:"OCR2ConfigPublicKey"`   // ocr2cfg_evm_<key>
	CSAPublicKey           string `json:"CSAPublicKey"`
	EncryptionPublicKey    string `json:"EncryptionPublicKey"`
}

// OCR2OracleConfig is the input configuration for an OCR2/3 contract.
type OCR2OracleConfig struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (c OCR2OracleConfig) MarshalJSON() ([]byte, error) {
	alias := struct {
		Signers               []string
		Transmitters          []string
		F                     uint8
		OnchainConfig         string
		OffchainConfigVersion uint64
		OffchainConfig        string
	}{
		Signers:               make([]string, len(c.Signers)),
		Transmitters:          make([]string, len(c.Transmitters)),
		F:                     c.F,
		OnchainConfig:         "0x" + hex.EncodeToString(c.OnchainConfig),
		OffchainConfigVersion: c.OffchainConfigVersion,
		OffchainConfig:        "0x" + hex.EncodeToString(c.OffchainConfig),
	}

	for i, signer := range c.Signers {
		alias.Signers[i] = hex.EncodeToString(signer)
	}

	for i, transmitter := range c.Transmitters {
		alias.Transmitters[i] = transmitter.Hex()
	}

	return json.Marshal(alias)
}

func (c *OCR2OracleConfig) UnmarshalJSON(data []byte) error {
	type aliasT struct {
		Signers               []string
		Transmitters          []string
		F                     uint8
		OnchainConfig         string
		OffchainConfigVersion uint64
		OffchainConfig        string
	}
	var alias aliasT
	err := json.Unmarshal(data, &alias)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OCR2OracleConfig alias: %w", err)
	}
	c.F = alias.F
	c.OffchainConfigVersion = alias.OffchainConfigVersion
	c.Signers = make([][]byte, len(alias.Signers))
	for i, signer := range alias.Signers {
		c.Signers[i], err = hex.DecodeString(strings.TrimPrefix(signer, "0x"))
		if err != nil {
			return fmt.Errorf("failed to decode signer: %w", err)
		}
	}
	c.Transmitters = make([]common.Address, len(alias.Transmitters))
	for i, transmitter := range alias.Transmitters {
		c.Transmitters[i] = common.HexToAddress(transmitter)
	}
	c.OnchainConfig, err = hex.DecodeString(strings.TrimPrefix(alias.OnchainConfig, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode onchain config: %w", err)
	}
	c.OffchainConfig, err = hex.DecodeString(strings.TrimPrefix(alias.OffchainConfig, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode offchain config: %w", err)
	}
	return nil
}

func GenerateOCR3Config(cfg OracleConfig, nca []NodeKeys, secrets focr.OCRSecrets, reportingPluginConfigOverride []byte) (OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
	}
	onchainPubKeys := [][]byte{}
	allPubKeys := map[string]any{}
	if secrets.IsEmpty() {
		return OCR2OracleConfig{}, errors.New("OCRSecrets is required")
	}
	for _, n := range nca {
		// evm keys always required
		if n.OCR2OnchainPublicKey == "" {
			return OCR2OracleConfig{}, errors.New("OCR2OnchainPublicKey is required")
		}
		ethPubKey := common.HexToAddress(n.OCR2OnchainPublicKey)
		pubKeys := map[string]types.OnchainPublicKey{
			string(chaintype.EVM): ethPubKey.Bytes(),
		}
		// add aptos key if present
		if n.AptosOnchainPublicKey != "" {
			aptosPubKey, err := hex.DecodeString(n.AptosOnchainPublicKey)
			if err != nil {
				return OCR2OracleConfig{}, fmt.Errorf("failed to decode AptosOnchainPublicKey: %w", err)
			}
			pubKeys[string(chaintype.Aptos)] = aptosPubKey
		}
		// add solana key if present
		if n.SolanaOnchainPublicKey != "" {
			solPubKey, err := hex.DecodeString(n.SolanaOnchainPublicKey)
			if err != nil {
				return OCR2OracleConfig{}, fmt.Errorf("failed to decode SolanaOnchainPublicKey: %w", err)
			}
			pubKeys[string(chaintype.Solana)] = solPubKey
		}

		// validate uniqueness of each individual key
		for _, key := range pubKeys {
			raw := hex.EncodeToString(key)
			_, exists := allPubKeys[raw]
			if exists {
				return OCR2OracleConfig{}, fmt.Errorf("Duplicate onchain public key: '%s'", raw)
			}
			allPubKeys[raw] = struct{}{}
		}
		pubKey, err := ocrcommon.MarshalMultichainPublicKey(pubKeys)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to marshal multichain public key: %w", err)
		}
		onchainPubKeys = append(onchainPubKeys, pubKey)
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2OffchainPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			return OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 offchain public key. expected %d but got %d", ed25519.PublicKeySize, nCopied)
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, pkBytesFixed)
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2ConfigPublicKey)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2ConfigPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 config public key. expected %d but got %d", ed25519.PublicKeySize, n)
		}

		configPubKeysBytes = append(configPubKeysBytes, pkBytesFixed)
	}

	identities := []confighelper.OracleIdentityExtra{}
	for index := range nca {
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(nca[index].EthAddress),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	cfgBytes := reportingPluginConfigOverride
	if cfgBytes == nil {
		offchainCfg, err := getOffchainCfg(cfg)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to get offchain config: %w", err)
		}
		if offchainCfg != nil {
			offchainCfgAsProto, err := offchainCfg.ToProto()
			if err != nil {
				return OCR2OracleConfig{}, fmt.Errorf("failed to convert offchainConfig to proto: %w", err)
			}
			cfgBytes, err = proto.Marshal(offchainCfgAsProto)
			if err != nil {
				return OCR2OracleConfig{}, fmt.Errorf("failed to marshal offchainConfig to proto: %w", err)
			}
		}
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsDeterministic(
		secrets.EphemeralSk,
		secrets.SharedSecret,
		time.Duration(cfg.DeltaProgressMillis)*time.Millisecond,
		time.Duration(cfg.DeltaResendMillis)*time.Millisecond,
		time.Duration(cfg.DeltaInitialMillis)*time.Millisecond,
		time.Duration(cfg.DeltaRoundMillis)*time.Millisecond,
		time.Duration(cfg.DeltaGraceMillis)*time.Millisecond,
		time.Duration(cfg.DeltaCertifiedCommitRequestMillis)*time.Millisecond,
		time.Duration(cfg.DeltaStageMillis)*time.Millisecond,
		cfg.MaxRoundsPerEpoch,
		cfg.TransmissionSchedule,
		identities,
		cfgBytes, // reportingPluginConfig
		nil,      // maxDurationInitialization
		time.Duration(cfg.MaxDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationObservationMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldAcceptMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldTransmitMillis)*time.Millisecond,
		cfg.MaxFaultyOracles,
		nil, // empty onChain config
	)
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to generate contract config args: %w", err)
	}

	var configSigners [][]byte
	for _, signer := range signers {
		configSigners = append(configSigners, signer)
	}

	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to convert transmitters to addresses: %w", err)
	}

	config := OCR2OracleConfig{
		Signers:               configSigners,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return config, nil
}

func getOffchainCfg(oracleCfg OracleConfig) (offchainConfig, error) {
	var result offchainConfig
	if oracleCfg.ConsensusCapOffchainConfig != nil {
		result = oracleCfg.ConsensusCapOffchainConfig
	}

	if oracleCfg.ChainCapOffchainConfig != nil {
		if result != nil {
			return nil, fmt.Errorf("multiple offchain configs specified: %+v. Only one allowed", oracleCfg)
		}

		result = oracleCfg.ChainCapOffchainConfig
	}

	return result, nil
}

type ConfigureOCR3Request struct {
	Cfg        *OracleConfig
	Chain      cldf_evm.Chain
	Contract   *ocr3_capability.OCR3Capability
	Nodes      []deployment.Node
	DryRun     bool
	OcrSecrets focr.OCRSecrets

	ReportingPluginConfigOverride []byte

	UseMCMS bool
}

func (r ConfigureOCR3Request) generateOCR3Config() (OCR2OracleConfig, error) {
	nks := makeNodeKeysSlice(r.Nodes, r.Chain.Selector)
	if r.Cfg == nil {
		return OCR2OracleConfig{}, errors.New("OCR3 config is required")
	}
	return GenerateOCR3Config(*r.Cfg, nks, r.OcrSecrets, r.ReportingPluginConfigOverride)
}

type ConfigureOCR3Response struct {
	OcrConfig OCR2OracleConfig
	Ops       *mcmstypes.BatchOperation
}

func ConfigureOCR3contract(req ConfigureOCR3Request) (*ConfigureOCR3Response, error) {
	if req.Contract == nil {
		return nil, errors.New("OCR3 contract is nil")
	}
	ocrConfig, err := req.generateOCR3Config()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}
	if req.DryRun {
		return &ConfigureOCR3Response{ocrConfig, nil}, nil
	}

	txOpt := req.Chain.DeployerKey
	if req.UseMCMS {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := req.Contract.SetConfig(txOpt,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	if err != nil {
		err = cldf.DecodeErr(ocr3_capability.OCR3CapabilityABI, err)
		return nil, fmt.Errorf("failed to call SetConfig for OCR3 contract %s using mcms: %T: %w", req.Contract.Address().String(), req.UseMCMS, err)
	}

	var ops mcmstypes.BatchOperation
	if !req.UseMCMS {
		_, err = req.Chain.Confirm(tx)
		if err != nil {
			err = cldf.DecodeErr(ocr3_capability.OCR3CapabilityABI, err)
			return nil, fmt.Errorf("failed to confirm SetConfig for OCR3 contract %s: %w", req.Contract.Address().String(), err)
		}
	} else {
		ops, err = proposalutils.BatchOperationForChain(req.Chain.Selector, req.Contract.Address().Hex(), tx.Data(), big.NewInt(0), string(OCR3Capability), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create batch operation: %w", err)
		}
	}

	return &ConfigureOCR3Response{ocrConfig, &ops}, nil
}

type ConfigureOCR3Resp struct {
	OCR2OracleConfig
	Ops *mcmstypes.BatchOperation
}

type ConfigureOCR3Config struct {
	ChainSel   uint64
	NodeIDs    []string
	Contract   *ocr3_capability.OCR3Capability
	OCR3Config *OracleConfig
	DryRun     bool

	ReportingPluginConfigOverride []byte

	UseMCMS bool
}

func ConfigureOCR3ContractFromJD(env *cldf.Environment, cfg ConfigureOCR3Config) (*ConfigureOCR3Resp, error) {
	prefix := ""
	if cfg.DryRun {
		prefix = "DRY RUN: "
	}
	env.Logger.Infof("%sconfiguring OCR3 contract for chain %d", prefix, cfg.ChainSel)
	if cfg.Contract == nil {
		return nil, errors.New("OCR3 contract is required")
	}

	evmChains := env.BlockChains.EVMChains()
	registryChain, ok := evmChains[cfg.ChainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found in environment", cfg.ChainSel)
	}

	contract := cfg.Contract

	nodes, err := deployment.NodeInfo(cfg.NodeIDs, env.Offchain)
	if err != nil {
		return nil, err
	}
	r, err := ConfigureOCR3contract(ConfigureOCR3Request{
		Cfg:                           cfg.OCR3Config,
		Chain:                         registryChain,
		Contract:                      contract,
		Nodes:                         nodes,
		DryRun:                        cfg.DryRun,
		UseMCMS:                       cfg.UseMCMS,
		OcrSecrets:                    env.OCRSecrets,
		ReportingPluginConfigOverride: cfg.ReportingPluginConfigOverride,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureOCR3Resp{
		OCR2OracleConfig: r.OcrConfig,
		Ops:              r.Ops,
	}, nil
}

func makeNodeKeysSlice(nodes []deployment.Node, registryChainSel uint64) []NodeKeys {
	var out []NodeKeys
	for _, n := range nodes {
		out = append(out, toNodeKeys(&n, registryChainSel))
	}
	return out
}

func toNodeKeys(o *deployment.Node, registryChainSel uint64) NodeKeys {
	var aptosOcr2KeyBundleID string
	var aptosOnchainPublicKey string
	var aptosCC *deployment.OCRConfig
	var solanaOcr2KeyBundleID string
	var solanaCC *deployment.OCRConfig
	var solanaOnchainPublickey string
	for details, cfg := range o.SelToOCRConfig {
		if family, err := chainsel.GetSelectorFamily(details.ChainSelector); err == nil {
			if family == chainsel.FamilyAptos {
				aptosCC = &cfg
			}
			if family == chainsel.FamilySolana {
				solanaCC = &cfg
			}
		}
	}

	if aptosCC != nil {
		aptosOcr2KeyBundleID = aptosCC.KeyBundleID
		aptosOnchainPublicKey = fmt.Sprintf("%x", aptosCC.OnchainPublicKey[:])
	}

	if solanaCC != nil {
		solanaOcr2KeyBundleID = solanaCC.KeyBundleID
		solanaOnchainPublickey = fmt.Sprintf("%x", solanaCC.OnchainPublicKey[:])
	}

	evmCC, exists := o.OCRConfigForChainSelector(registryChainSel)
	if !exists {
		panic(fmt.Sprintf("ocr2 config not found for chain selector %d, node %s", registryChainSel, o.NodeID))
	}
	return NodeKeys{
		EthAddress:            string(evmCC.TransmitAccount),
		P2PPeerID:             strings.TrimPrefix(o.PeerID.String(), "p2p_"),
		OCR2BundleID:          evmCC.KeyBundleID,
		OCR2OffchainPublicKey: hex.EncodeToString(evmCC.OffchainPublicKey[:]),
		OCR2OnchainPublicKey:  fmt.Sprintf("%x", evmCC.OnchainPublicKey[:]),
		OCR2ConfigPublicKey:   hex.EncodeToString(evmCC.ConfigEncryptionPublicKey[:]),
		CSAPublicKey:          o.CSAKey,
		// default value of encryption public key is the CSA public key
		// TODO: DEVSVCS-760
		EncryptionPublicKey: strings.TrimPrefix(o.CSAKey, "csa_"),
		// TODO Aptos support. How will that be modeled in clo data?
		// TODO: AptosAccount is unset but probably unused
		AptosBundleID:          aptosOcr2KeyBundleID,
		AptosOnchainPublicKey:  aptosOnchainPublicKey,
		SolanaOnchainPublicKey: solanaOnchainPublickey,
		SolanaBundleID:         solanaOcr2KeyBundleID,
	}
}
