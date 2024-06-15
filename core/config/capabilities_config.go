package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type CapabilitiesExternalRegistry interface {
	Address() string
	NetworkID() string
	ChainID() string
	RelayID() types.RelayID
}

type Capabilities interface {
	Peering() P2P
	ExternalRegistry() CapabilitiesExternalRegistry
}
