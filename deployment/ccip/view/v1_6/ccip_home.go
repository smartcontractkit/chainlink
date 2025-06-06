package v1_6

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	cciptypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type DonView struct {
	DonID         uint32        `json:"donID"`
	CommitConfigs GetAllConfigs `json:"commitConfigs"`
	ExecConfigs   GetAllConfigs `json:"execConfigs"`
}

type GetAllConfigs struct {
	ActiveConfig    CCIPHomeVersionedConfig `json:",omitempty"`
	CandidateConfig CCIPHomeVersionedConfig `json:",omitempty"`
}

type OracleIdentity struct {
	OffchainPublicKey string
	OnchainPublicKey  string
	PeerID            string
	TransmitAccount   string
}

type CCIPHomeOCR3Config struct {
	PluginType                              uint8                               `json:",omitempty"`
	ChainSelector                           uint64                              `json:",omitempty"`
	FRoleDON                                uint8                               `json:",omitempty"`
	OffchainConfigVersion                   uint64                              `json:",omitempty"`
	OfframpAddress                          string                              `json:",omitempty"`
	RmnHomeAddress                          string                              `json:",omitempty"`
	Nodes                                   []CCIPHomeOCR3Node                  `json:",omitempty"`
	DeltaProgress                           string                              `json:",omitempty"`
	DeltaResend                             string                              `json:",omitempty"`
	DeltaInitial                            string                              `json:",omitempty"`
	DeltaRound                              string                              `json:",omitempty"`
	DeltaGrace                              string                              `json:",omitempty"`
	DeltaCertifiedCommitRequest             string                              `json:",omitempty"`
	DeltaStage                              string                              `json:",omitempty"`
	RMax                                    uint64                              `json:",omitempty"`
	MaxDurationInitialization               string                              `json:",omitempty"`
	MaxDurationQuery                        string                              `json:",omitempty"`
	MaxDurationObservation                  string                              `json:",omitempty"`
	MaxDurationShouldAcceptAttestedReport   string                              `json:",omitempty"`
	MaxDurationShouldTransmitAcceptedReport string                              `json:",omitempty"`
	F                                       int                                 `json:",omitempty"`
	CommitOffChainConfig                    *pluginconfig.CommitOffchainConfig  `json:",omitempty"`
	ExecuteOffChainConfig                   *pluginconfig.ExecuteOffchainConfig `json:",omitempty"`
	S                                       []int                               `json:",omitempty"`
	OracleIdentities                        []OracleIdentity                    `json:",omitempty"`
}

type CCIPHomeOCR3Node struct {
	P2pID          string `json:",omitempty"`
	SignerKey      string `json:",omitempty"`
	TransmitterKey string `json:",omitempty"`
}

type CCIPHomeVersionedConfig struct {
	Version      uint32             `json:",omitempty"`
	ConfigDigest []byte             `json:",omitempty"`
	Config       CCIPHomeOCR3Config `json:",omitempty"`
}

type CCIPHomeView struct {
	types.ContractMetaData
	ChainConfigs       []CCIPHomeChainConfigArgView `json:"chainConfigs"`
	CapabilityRegistry common.Address               `json:"capabilityRegistry"`
	Dons               []DonView                    `json:"dons"`
}

type CCIPHomeChainConfigArgView struct {
	ChainSelector uint64
	ChainConfig   CCIPHomeChainConfigView
}

type CCIPHomeChainConfigView struct {
	Readers []string
	FChain  uint8
	Config  chainconfig.ChainConfig
}

func GenerateCCIPHomeView(cr *capabilities_registry.CapabilitiesRegistry, ch *ccip_home.CCIPHome) (CCIPHomeView, error) {
	if ch == nil {
		return CCIPHomeView{}, errors.New("cannot generate view for nil CCIPHome")
	}
	meta, err := types.NewContractMetaData(ch, ch.Address())
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to generate contract metadata for CCIPHome %s: %w", ch.Address(), err)
	}
	numChains, err := ch.GetNumChainConfigurations(nil)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get number of chain configurations for CCIPHome %s: %w", ch.Address(), err)
	}
	// Pagination shouldn't be required here, but we can add it if needed.
	chainCfg, err := ch.GetAllChainConfigs(nil, big.NewInt(0), numChains)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get all chain configs for CCIPHome %s: %w", ch.Address(), err)
	}
	var chains []CCIPHomeChainConfigArgView
	for _, cfg := range chainCfg {
		decodedChainCfg, err := chainconfig.DecodeChainConfig(cfg.ChainConfig.Config)
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to decode chain config for CCIPHome %s: %w", ch.Address(), err)
		}
		var readers []string
		for _, r := range cfg.ChainConfig.Readers {
			readers = append(readers, libocrtypes.PeerID(r).String())
		}
		chains = append(chains, CCIPHomeChainConfigArgView{
			ChainSelector: cfg.ChainSelector,
			ChainConfig: CCIPHomeChainConfigView{
				FChain:  cfg.ChainConfig.FChain,
				Config:  decodedChainCfg,
				Readers: readers,
			},
		})
	}

	crAddr, err := ch.GetCapabilityRegistry(nil)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get capability registry for CCIPHome %s: %w", ch.Address(), err)
	}
	if crAddr != cr.Address() {
		return CCIPHomeView{}, fmt.Errorf("capability registry address mismatch for CCIPHome %s: %w", ch.Address(), err)
	}
	dons, err := shared.GetCCIPDonsFromCapRegistry(context.Background(), cr)
	if err != nil {
		return CCIPHomeView{}, fmt.Errorf("failed to get DONs for CCIPHome %s: %w", ch.Address(), err)
	}
	// Get every don's configuration.
	var dvs []DonView
	for _, d := range dons {
		commitConfigs, err := ch.GetAllConfigs(nil, d.Id, uint8(cciptypes.PluginTypeCCIPCommit))
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to get commit config for CCIPHome %s: %w", ch.Address(), err)
		}
		execConfigs, err := ch.GetAllConfigs(nil, d.Id, uint8(cciptypes.PluginTypeCCIPExec))
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to get exec config for CCIPHome %s: %w", ch.Address(), err)
		}
		commitCfg, err := toGetAllConfigsView(commitConfigs, cciptypes.PluginTypeCCIPCommit)
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to convert commit config for CCIPHome %s: %w", ch.Address(), err)
		}
		execCfg, err := toGetAllConfigsView(execConfigs, cciptypes.PluginTypeCCIPExec)
		if err != nil {
			return CCIPHomeView{}, fmt.Errorf("failed to convert exec config for CCIPHome %s: %w", ch.Address(), err)
		}
		dvs = append(dvs, DonView{
			DonID:         d.Id,
			CommitConfigs: commitCfg,
			ExecConfigs:   execCfg,
		})
	}
	return CCIPHomeView{
		ContractMetaData:   meta,
		ChainConfigs:       chains,
		CapabilityRegistry: crAddr,
		Dons:               dvs,
	}, nil
}

func toGetAllConfigsView(cfg ccip_home.GetAllConfigs, pluginType cciptypes.PluginType) (GetAllConfigs, error) {
	active, err := toCCIPHomeVersionedConfig(cfg.ActiveConfig, pluginType)
	if err != nil {
		return GetAllConfigs{}, fmt.Errorf("failed to convert active config: %w", err)
	}
	candidate, err := toCCIPHomeVersionedConfig(cfg.CandidateConfig, pluginType)
	if err != nil {
		return GetAllConfigs{}, fmt.Errorf("failed to convert candidate config: %w", err)
	}
	return GetAllConfigs{
		ActiveConfig:    active,
		CandidateConfig: candidate,
	}, nil
}

func toCCIPHomeVersionedConfig(cfg ccip_home.CCIPHomeVersionedConfig, pluginType cciptypes.PluginType) (CCIPHomeVersionedConfig, error) {
	var nodes []CCIPHomeOCR3Node
	for _, n := range cfg.Config.Nodes {
		peerID := p2pkey.PeerID(n.P2pId)
		nodes = append(nodes, CCIPHomeOCR3Node{
			P2pID:          peerID.String(),
			SignerKey:      ccipocr3.UnknownAddress(n.SignerKey).String(),
			TransmitterKey: ccipocr3.UnknownAddress(n.TransmitterKey).String(),
		})
	}
	offRampAddr := ccipocr3.UnknownAddress(cfg.Config.OfframpAddress).String()
	rmnAddr := ccipocr3.UnknownAddress(cfg.Config.RmnHomeAddress).String()
	c := CCIPHomeVersionedConfig{
		Version:      cfg.Version,
		ConfigDigest: cfg.ConfigDigest[:],
		Config: CCIPHomeOCR3Config{
			PluginType:            cfg.Config.PluginType,
			ChainSelector:         cfg.Config.ChainSelector,
			FRoleDON:              cfg.Config.FRoleDON,
			OffchainConfigVersion: cfg.Config.OffchainConfigVersion,
			OfframpAddress:        offRampAddr,
			RmnHomeAddress:        rmnAddr,
			Nodes:                 nodes,
		},
	}
	if err := populateDecodedOCRParams(&c.Config, cfg, pluginType); err != nil {
		return c, fmt.Errorf("failed to populate decoded OCR params: %w", err)
	}
	return c, nil
}

func populateDecodedOCRParams(config *CCIPHomeOCR3Config, ccipHomeCfg ccip_home.CCIPHomeVersionedConfig, pluginType cciptypes.PluginType) error {
	// empty config means no config is set, so we skip it
	if ccipHomeCfg.ConfigDigest == [32]byte{} {
		return nil
	}
	signers := make([]ocrtypes.OnchainPublicKey, 0)
	transmitters := make([]ocrtypes.Account, 0)
	for _, node := range ccipHomeCfg.Config.Nodes {
		signers = append(signers, node.SignerKey)
		transmitters = append(transmitters, ocrtypes.Account(node.TransmitterKey))
	}
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(false, ocrtypes.ContractConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.FRoleDON,
		OnchainConfig:         []byte{}, // empty OnChainConfig
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        ccipHomeCfg.Config.OffchainConfig,
	})
	if err != nil {
		return fmt.Errorf("failed to decode offchain config: %w", err)
	}
	identities := make([]OracleIdentity, 0, len(publicConfig.OracleIdentities))
	for _, o := range publicConfig.OracleIdentities {
		identities = append(identities, OracleIdentity{
			OffchainPublicKey: ccipocr3.UnknownAddress(o.OffchainPublicKey[:]).String(),
			OnchainPublicKey:  ccipocr3.UnknownAddress(o.OnchainPublicKey[:]).String(),
			PeerID:            o.PeerID,
			TransmitAccount:   ccipocr3.UnknownAddress(o.TransmitAccount[:]).String(),
		})
	}
	if pluginType == cciptypes.PluginTypeCCIPCommit {
		commitOffChain, err := pluginconfig.DecodeCommitOffchainConfig(publicConfig.ReportingPluginConfig)
		if err != nil {
			return fmt.Errorf("failed to decode commit offchain config: %w", err)
		}
		config.CommitOffChainConfig = &commitOffChain
	}
	if pluginType == cciptypes.PluginTypeCCIPExec {
		executeOffChainConfig, err := pluginconfig.DecodeExecuteOffchainConfig(publicConfig.ReportingPluginConfig)
		if err != nil {
			return fmt.Errorf("failed to decode execute offchain config: %w", err)
		}
		config.ExecuteOffChainConfig = &executeOffChainConfig
	}
	config.DeltaProgress = publicConfig.DeltaProgress.String()
	config.DeltaResend = publicConfig.DeltaResend.String()
	config.DeltaInitial = publicConfig.DeltaInitial.String()
	config.DeltaRound = publicConfig.DeltaRound.String()
	config.DeltaGrace = publicConfig.DeltaGrace.String()
	config.DeltaCertifiedCommitRequest = publicConfig.DeltaCertifiedCommitRequest.String()
	config.DeltaStage = publicConfig.DeltaStage.String()
	config.RMax = publicConfig.RMax
	config.S = publicConfig.S
	config.OracleIdentities = identities
	if publicConfig.MaxDurationInitialization != nil {
		d := *publicConfig.MaxDurationInitialization
		config.MaxDurationInitialization = d.String()
	}
	config.MaxDurationQuery = publicConfig.MaxDurationQuery.String()
	config.MaxDurationObservation = publicConfig.MaxDurationObservation.String()
	config.MaxDurationShouldAcceptAttestedReport = publicConfig.MaxDurationShouldAcceptAttestedReport.String()
	config.MaxDurationShouldTransmitAcceptedReport = publicConfig.MaxDurationShouldTransmitAcceptedReport.String()
	config.F = publicConfig.F
	return nil
}
