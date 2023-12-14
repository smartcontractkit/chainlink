package evm

import (
	"github.com/ethereum/go-ethereum/common"
)

// Bindings Key being contract name.
type Bindings map[string]ContractBinding

// ContractBinding key being read name(event or method)
type ContractBinding map[string]common.Address
