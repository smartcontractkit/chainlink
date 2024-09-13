package v1_6

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
)

type RMN struct {
	types.ContractMetaData
	IsCursed bool                     `json:"isCursed"`
	Config   RMNRemoteVersionedConfig `json:"config,omitempty"`
}

type RMNRemoteVersionedConfig struct {
	Version    uint32            `json:"version"`
	Signers    []RMNRemoteSigner `json:"signers"`
	MinSigners uint64            `json:"minSigners"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"`
	NodeIndex        uint64 `json:"node_index"`
}

func RMNSnapshot(rmnReader *rmn_remote.RMNRemote) (RMN, error) {
	tv, err := rmnReader.TypeAndVersion(nil)
	if err != nil {
		return RMN{}, err
	}
	config, err := rmnReader.GetVersionedConfig(nil)
	if err != nil {
		return RMN{}, err
	}
	rmnConfig := RMNRemoteVersionedConfig{
		Version:    config.Version,
		Signers:    make([]RMNRemoteSigner, 0, len(config.Config.Signers)),
		MinSigners: config.Config.MinSigners,
	}
	for _, signer := range config.Config.Signers {
		rmnConfig.Signers = append(rmnConfig.Signers, RMNRemoteSigner{
			OnchainPublicKey: signer.OnchainPublicKey.Hex(),
			NodeIndex:        signer.NodeIndex,
		})
	}
	isCursed, err := rmnReader.IsCursed0(nil)
	if err != nil {
		return RMN{}, err
	}
	return RMN{
		ContractMetaData: types.ContractMetaData{
			Address:        rmnReader.Address(),
			TypeAndVersion: tv,
		},
		IsCursed: isCursed,
		Config:   rmnConfig,
	}, nil
}
