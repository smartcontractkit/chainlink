package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type chainWriterConfig struct {
	c toml.ChainWriter
}

func (b *chainWriterConfig) FromAddress() *ethkey.EIP55Address {
	return b.c.FromAddress
}

func (b *chainWriterConfig) ForwarderAddress() *ethkey.EIP55Address {
	return b.c.ForwarderAddress
}
