package v1_6

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
)

type RMNHomeView struct {
	types.ContractMetaData
	CandidateConfig *RMNHomeVersionedConfig `json:"candidateConfig,omitempty"`
	ActiveConfig    *RMNHomeVersionedConfig `json:"activeConfig,omitempty"`
}

type RMNHomeVersionedConfig struct {
	Version       uint32               `json:"version"`
	StaticConfig  RMNHomeStaticConfig  `json:"staticConfig"`
	DynamicConfig RMNHomeDynamicConfig `json:"dynamicConfig"`
	Digest        [32]byte             `json:"digest"`
}

type RMNHomeStaticConfig struct {
	Nodes []RMNHomeNode `json:"nodes"`
}

type RMNHomeDynamicConfig struct {
	SourceChains []RMNHomeSourceChain `json:"sourceChains"`
}

type RMNHomeSourceChain struct {
	ChainSelector       uint64   `json:"selector"`
	F                   uint64   `json:"f"`
	ObserverNodesBitmap *big.Int `json:"observerNodesBitmap"`
}

type RMNHomeNode struct {
	PeerId            [32]byte `json:"peerId"`
	OffchainPublicKey [32]byte `json:"offchainPublicKey"`
}

type DigestFunc func(*bind.CallOpts) ([32]byte, error)

func mapNodes(nodes []rmn_home.RMNHomeNode) []RMNHomeNode {
	result := make([]RMNHomeNode, len(nodes))
	for i, node := range nodes {
		result[i] = RMNHomeNode{
			PeerId:            node.PeerId,
			OffchainPublicKey: node.OffchainPublicKey,
		}
	}
	return result
}

func mapSourceChains(chains []rmn_home.RMNHomeSourceChain) []RMNHomeSourceChain {
	result := make([]RMNHomeSourceChain, len(chains))
	for i, chain := range chains {
		result[i] = RMNHomeSourceChain{
			ChainSelector:       chain.ChainSelector,
			F:                   chain.F,
			ObserverNodesBitmap: chain.ObserverNodesBitmap,
		}
	}
	return result
}

func generateRmnHomeVersionedConfig(reader *rmn_home.RMNHome, digestFunc DigestFunc) (*RMNHomeVersionedConfig, error) {
	digest, err := digestFunc(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get digest: %w", err)
	}

	if digest == [32]byte{} {
		return nil, nil
	}

	config, err := reader.GetConfig(nil, digest)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	staticConfig := RMNHomeStaticConfig{
		Nodes: mapNodes(config.VersionedConfig.StaticConfig.Nodes),
	}

	dynamicConfig := RMNHomeDynamicConfig{
		SourceChains: mapSourceChains(config.VersionedConfig.DynamicConfig.SourceChains),
	}

	return &RMNHomeVersionedConfig{
		Version:       config.VersionedConfig.Version,
		Digest:        config.VersionedConfig.ConfigDigest,
		StaticConfig:  staticConfig,
		DynamicConfig: dynamicConfig,
	}, nil
}

func GenerateRMNHomeView(rmnReader *rmn_home.RMNHome) (RMNHomeView, error) {
	if rmnReader == nil {
		return RMNHomeView{}, nil
	}

	activeConfig, err := generateRmnHomeVersionedConfig(rmnReader, rmnReader.GetActiveDigest)
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to generate active config: %w", err)
	}

	candidateConfig, err := generateRmnHomeVersionedConfig(rmnReader, rmnReader.GetCandidateDigest)
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to generate candidate config: %w", err)
	}

	contractMetaData, err := types.NewContractMetaData(rmnReader, rmnReader.Address())
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to create contract metadata: %w", err)
	}

	return RMNHomeView{
		ContractMetaData: contractMetaData,
		CandidateConfig:  candidateConfig,
		ActiveConfig:     activeConfig,
	}, nil
}
