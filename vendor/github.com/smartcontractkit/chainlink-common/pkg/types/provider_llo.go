package types

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type LLOProvider interface {
	ConfigProvider
	ContractTransmitter() llo.Transmitter
	ChannelDefinitionCache() llo.ChannelDefinitionCache
}
