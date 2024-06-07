package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type CapabilitiesRegistry interface {
	RemoteAddress() string
	NetworkID() string
	ChainID() string
	RelayID() types.RelayID
}

type Capabilities interface {
	Peering() P2P
	Registry() CapabilitiesRegistry
}
