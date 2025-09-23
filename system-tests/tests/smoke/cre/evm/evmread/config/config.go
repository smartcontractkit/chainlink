package config

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Config struct {
	ChainSelector    uint64
	ContractAddress  []byte
	AccountAddress   []byte
	ExpectedBalance  *big.Int
	TxHash           []byte
	ExpectedReceipt  *types.Receipt
	ExpectedBinaryTx []byte
}
