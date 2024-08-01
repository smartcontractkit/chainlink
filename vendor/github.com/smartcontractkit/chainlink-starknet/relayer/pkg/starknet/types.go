package starknet

import (
	"github.com/NethermindEth/juno/core/felt"
)

type CallOps struct {
	ContractAddress *felt.Felt
	Selector        *felt.Felt
	Calldata        []*felt.Felt
}
