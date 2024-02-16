package types

import (
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type LLOProvider interface {
	ConfigProvider
	ContractTransmitter() llotypes.Transmitter
	ChannelDefinitionCache() llotypes.ChannelDefinitionCache
}
