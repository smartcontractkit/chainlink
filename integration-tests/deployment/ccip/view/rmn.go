package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type RMN struct {
	Contract
	IsCursed bool                     `json:"isCursed"`
	Config   RMNRemoteVersionedConfig `json:"config"`
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

func RMNSnapshot(rmnReader RMNReader) (RMN, error) {
	tv, err := rmnReader.TypeAndVersion(nil)
	if err != nil {
		return RMN{}, err
	}
	config, err := rmnReader.GetVersionedConfig(nil)
	if err != nil {
		return RMN{}, err
	}
	isCursed, err := rmnReader.IsCursed0(nil)
	if err != nil {
		return RMN{}, err
	}
	return RMN{
		Contract: Contract{
			Address:        rmnReader.Address().Hex(),
			TypeAndVersion: tv,
		},
		IsCursed: isCursed,
		Config:   config,
	}, nil
}

type RMNReader interface {
	ContractState
	GetVersionedConfig(opts *bind.CallOpts) (RMNRemoteVersionedConfig, error)
	IsCursed(opts *bind.CallOpts, subject [16]byte) (bool, error)
	IsCursed0(opts *bind.CallOpts) (bool, error)
}
