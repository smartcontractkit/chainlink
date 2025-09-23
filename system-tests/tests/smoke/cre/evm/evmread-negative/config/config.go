package config

import "github.com/ethereum/go-ethereum/common"

type Config struct {
	ChainSelector  uint64
	FunctionToTest string
	InvalidInput   string
	BalanceReader
}

type BalanceReader struct {
	BalanceReaderAddress common.Address
}
