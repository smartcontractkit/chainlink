package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type chainWriterConfig struct {
	c toml.ChainWriter
}

func (b *chainWriterConfig) FromAddress() *types.EIP55Address {
	return b.c.FromAddress
}

func (b *chainWriterConfig) ForwarderAddress() *types.EIP55Address {
	return b.c.ForwarderAddress
}

func (b *chainWriterConfig) GasLimit() uint64 {
	return *b.c.GasLimit
}

func (b *chainWriterConfig) Checker() string {
	return *b.c.Checker
}

func (b *chainWriterConfig) ABI() string {
	return *b.c.ABI
}

func (b *chainWriterConfig) ContractFunction() string {
	return *b.c.ContractFunction
}
