package v0_5_0

import "github.com/ethereum/go-ethereum/common"

type ConfiguratorConfig struct {
	ConfiguratorAddress   common.Address
	ConfigID              [32]byte
	Signers               [][]byte
	OffchainTransmitters  [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (sc ConfiguratorConfig) GetContractAddress() common.Address { return sc.ConfiguratorAddress }
