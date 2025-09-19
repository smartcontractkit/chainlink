package ocr3

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

const offchainPublicKeyType byte = 0x8

func oCR3CapabilityCompatibleOnchainPublicKey(offchainPublicKey types.OffchainPublicKey) types.OnchainPublicKey {
	result := make([]byte, 0, 1+2+len(offchainPublicKey))
	result = append(result, offchainPublicKeyType)
	result = binary.LittleEndian.AppendUint16(result, uint16(len(offchainPublicKey)))
	result = append(result, offchainPublicKey[:]...)

	return result
}

func GenerateDKGConfig(cfg OracleConfig, nca []NodeKeys, secrets ocr.OCRSecrets, dkgCfg dkgocrtypes.ReportingPluginConfig) (OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
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

	onChainPublicKeys := make([]types.OnchainPublicKey, 0, len(offchainPubKeysBytes))
	for _, pk := range offchainPubKeysBytes {
		onChainPublicKeys = append(onChainPublicKeys, oCR3CapabilityCompatibleOnchainPublicKey(pk))
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
				OnchainPublicKey:  onChainPublicKeys[index],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(common.HexToAddress(fmt.Sprintf("0xc1c1c1c1%x", offchainPubKeysBytes[index][:16])).Hex()),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	cfgBytes, err := dkgCfg.MarshalBinary()
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to marshal ReportingPluginConfig: %w", err)
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

type ConfigureDKGRequest struct {
	Cfg                   *OracleConfig
	Chain                 cldf_evm.Chain
	Contract              *ocr3_capability.OCR3Capability
	Nodes                 []deployment.Node
	DryRun                bool
	OcrSecrets            ocr.OCRSecrets
	ReportingPluginConfig dkgocrtypes.ReportingPluginConfig
	UseMCMS               bool
}

func (r ConfigureDKGRequest) generateDKGConfig() (OCR2OracleConfig, error) {
	nks := makeNodeKeysSlice(r.Nodes, r.Chain.Selector)
	if r.Cfg == nil {
		return OCR2OracleConfig{}, errors.New("OCR3 config is required")
	}
	return GenerateDKGConfig(*r.Cfg, nks, r.OcrSecrets, r.ReportingPluginConfig)
}

type ConfigureDKGResponse struct {
	OcrConfig OCR2OracleConfig
	Ops       *mcmstypes.BatchOperation
}

type ConfigureDKGConfig struct {
	ChainSel              uint64
	NodeIDs               []string
	Contract              *ocr3_capability.OCR3Capability
	OCR3Config            *OracleConfig
	DryRun                bool
	ReportingPluginConfig dkgocrtypes.ReportingPluginConfig

	UseMCMS bool
}

func ConfigureDKGContract(req ConfigureDKGRequest) (*ConfigureDKGResponse, error) {
	if req.Contract == nil {
		return nil, errors.New("OCR3 contract is nil")
	}
	ocrConfig, err := req.generateDKGConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}
	if req.DryRun {
		return &ConfigureDKGResponse{ocrConfig, nil}, nil
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

	return &ConfigureDKGResponse{ocrConfig, &ops}, nil
}

type ConfigureDKGResp struct {
	OCR2OracleConfig
	Ops *mcmstypes.BatchOperation
}

func ConfigureDKGContractFromJD(env *cldf.Environment, cfg ConfigureDKGConfig) (*ConfigureDKGResp, error) {
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
	r, err := ConfigureDKGContract(ConfigureDKGRequest{
		Cfg:                   cfg.OCR3Config,
		Chain:                 registryChain,
		Contract:              contract,
		Nodes:                 nodes,
		DryRun:                cfg.DryRun,
		UseMCMS:               cfg.UseMCMS,
		OcrSecrets:            env.OCRSecrets,
		ReportingPluginConfig: cfg.ReportingPluginConfig,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureDKGResp{
		OCR2OracleConfig: r.OcrConfig,
		Ops:              r.Ops,
	}, nil
}
