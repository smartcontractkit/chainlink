package ethcontractconfig

import "github.com/ethereum/go-ethereum/common"

type SetConfigArgs struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}
