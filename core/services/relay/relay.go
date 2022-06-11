package relay

import relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

type Network string

var (
	EVM             Network = "evm"
	Solana          Network = "solana"
	Terra           Network = "terra"
	SupportedRelays         = map[Network]struct{}{
		EVM:    {},
		Solana: {},
		Terra:  {},
	}
)

type VRFRelayer interface {
	relaytypes.Relayer
	NewDKGProvider(rargs relaytypes.RelayArgs, transmitterID string) (DKGProvider, error)
	NewVRFProvider(rargs relaytypes.RelayArgs, transmitterID string) (VRFProvider, error)
}

type DKGProvider interface {
	relaytypes.Plugin
}

type VRFProvider interface {
	relaytypes.Plugin
}
