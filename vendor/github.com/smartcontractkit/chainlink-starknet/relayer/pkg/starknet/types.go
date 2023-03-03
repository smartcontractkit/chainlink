package starknet

import (
	caigotypes "github.com/dontpanicdao/caigo/types"
)

type CallOps struct {
	ContractAddress caigotypes.Hash
	Selector        string
	Calldata        []string
}
