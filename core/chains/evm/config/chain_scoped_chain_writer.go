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
