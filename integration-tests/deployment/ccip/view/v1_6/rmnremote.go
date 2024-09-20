package v1_6

import (
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
)

type RMNRemoteView struct {
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

func GenerateRMNRemoteView(rmnReader *rmn_remote.RMNRemote) (RMNRemoteView, error) {
	tv, err := types.NewContractMetaData(rmnReader, rmnReader.Address())
	if err != nil {
		return RMNRemoteView{}, err
	}
	config, err := rmnReader.GetVersionedConfig(nil)
	if err != nil {
		return RMNRemoteView{}, err
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
		return RMNRemoteView{}, err
	}
	return RMNRemoteView{
		ContractMetaData: tv,
		IsCursed:         isCursed,
		Config:           rmnConfig,
	}, nil
}
