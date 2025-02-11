package v1_6

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
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

func decodeHexString(hexStr string, expectedLength int) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	if len(bytes) != expectedLength {
		return nil, fmt.Errorf("invalid length: expected %d, got %d", expectedLength, len(bytes))
	}
	return bytes, nil
}

func (c RMNHomeVersionedConfig) MarshalJSON() ([]byte, error) {
	type Alias RMNHomeVersionedConfig
	return json.Marshal(&struct {
		Digest string `json:"digest"`
		*Alias
	}{
		Digest: hex.EncodeToString(c.Digest[:]),
		Alias:  (*Alias)(&c),
	})
}

func (c *RMNHomeVersionedConfig) UnmarshalJSON(data []byte) error {
	type Alias RMNHomeVersionedConfig
	aux := &struct {
		Digest string `json:"digest"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	digestBytes, err := decodeHexString(aux.Digest, 32)
	if err != nil {
		return err
	}
	copy(c.Digest[:], digestBytes)
	return nil
}

type RMNHomeStaticConfig struct {
	Nodes []RMNHomeNode `json:"nodes"`
}

type RMNHomeDynamicConfig struct {
	SourceChains []RMNHomeSourceChain `json:"sourceChains"`
}

type RMNHomeSourceChain struct {
	ChainSelector       uint64   `json:"selector"`
	FObserve            uint64   `json:"fObserve"`
	ObserverNodesBitmap *big.Int `json:"observerNodesBitmap"`
}

type RMNHomeNode struct {
	PeerId            [32]byte `json:"peerId"`
	OffchainPublicKey [32]byte `json:"offchainPublicKey"`
}

func (n RMNHomeNode) MarshalJSON() ([]byte, error) {
	type Alias RMNHomeNode
	return json.Marshal(&struct {
		PeerId            string `json:"peerId"`
		OffchainPublicKey string `json:"offchainPublicKey"`
		*Alias
	}{
		PeerId:            hex.EncodeToString(n.PeerId[:]),
		OffchainPublicKey: hex.EncodeToString(n.OffchainPublicKey[:]),
		Alias:             (*Alias)(&n),
	})
}

func (n *RMNHomeNode) UnmarshalJSON(data []byte) error {
	type Alias RMNHomeNode
	aux := &struct {
		PeerId            string `json:"peerId"`
		OffchainPublicKey string `json:"offchainPublicKey"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	peerIdBytes, err := decodeHexString(aux.PeerId, 32)
	if err != nil {
		return err
	}
	copy(n.PeerId[:], peerIdBytes)

	offchainPublicKeyBytes, err := decodeHexString(aux.OffchainPublicKey, 32)
	if err != nil {
		return err
	}
	copy(n.OffchainPublicKey[:], offchainPublicKeyBytes)

	return nil
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
			FObserve:            chain.FObserve,
			ObserverNodesBitmap: chain.ObserverNodesBitmap,
		}
	}
	return result
}

func generateRmnHomeVersionedConfig(reader *rmn_home.RMNHome, digestFunc DigestFunc) (*RMNHomeVersionedConfig, error) {
	address := reader.Address()
	digest, err := digestFunc(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get digest for contract %s: %w", address, err)
	}

	if digest == [32]byte{} {
		return nil, nil
	}

	config, err := reader.GetConfig(nil, digest)
	if err != nil {
		return nil, fmt.Errorf("failed to get config for contract %s: %w", address, err)
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

	address := rmnReader.Address()

	activeConfig, err := generateRmnHomeVersionedConfig(rmnReader, rmnReader.GetActiveDigest)
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to generate active config for contract %s: %w", address, err)
	}

	candidateConfig, err := generateRmnHomeVersionedConfig(rmnReader, rmnReader.GetCandidateDigest)
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to generate candidate config for contract %s: %w", address, err)
	}

	contractMetaData, err := types.NewContractMetaData(rmnReader, rmnReader.Address())
	if err != nil {
		return RMNHomeView{}, fmt.Errorf("failed to create contract metadata for contract %s: %w", address, err)
	}

	return RMNHomeView{
		ContractMetaData: contractMetaData,
		CandidateConfig:  candidateConfig,
		ActiveConfig:     activeConfig,
	}, nil
}
