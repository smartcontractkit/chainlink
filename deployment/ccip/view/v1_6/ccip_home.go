package v1_6

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	cciptypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
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

type CCIPHomeOCR3Config struct {
	PluginType            uint8              `json:",omitempty"`
	ChainSelector         uint64             `json:",omitempty"`
	FRoleDON              uint8              `json:",omitempty"`
	OffchainConfigVersion uint64             `json:",omitempty"`
	OfframpAddress        string             `json:",omitempty"`
	RmnHomeAddress        string             `json:",omitempty"`
	Nodes                 []CCIPHomeOCR3Node `json:",omitempty"`
	OffchainConfig        []byte             `json:",omitempty"`
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
	Readers [][]byte
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
		var readers [][]byte
		for _, r := range cfg.ChainConfig.Readers {
			readers = append(readers, r[:])
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
	dons, err := cr.GetDONs(nil)
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
		dvs = append(dvs, DonView{
			DonID:         d.Id,
			CommitConfigs: toGetAllConfigsView(commitConfigs),
			ExecConfigs:   toGetAllConfigsView(execConfigs),
		})
	}
	return CCIPHomeView{
		ContractMetaData:   meta,
		ChainConfigs:       chains,
		CapabilityRegistry: crAddr,
		Dons:               dvs,
	}, nil
}

func toGetAllConfigsView(cfg ccip_home.GetAllConfigs) GetAllConfigs {
	return GetAllConfigs{
		ActiveConfig:    toCCIPHomeVersionedConfig(cfg.ActiveConfig),
		CandidateConfig: toCCIPHomeVersionedConfig(cfg.CandidateConfig),
	}
}

func toCCIPHomeVersionedConfig(cfg ccip_home.CCIPHomeVersionedConfig) CCIPHomeVersionedConfig {
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
	return CCIPHomeVersionedConfig{
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
			OffchainConfig:        cfg.Config.OffchainConfig,
		},
	}
}
