package relay

const (
	NetworkEVM      = "evm"
	NetworkCosmos   = "cosmos"
	NetworkSolana   = "solana"
	NetworkStarkNet = "starknet"
	NetworkAptos    = "aptos"

	NetworkDummy = "dummy"
)

var SupportedNetworks = map[string]struct{}{
	NetworkEVM:      {},
	NetworkCosmos:   {},
	NetworkSolana:   {},
	NetworkStarkNet: {},
	NetworkAptos:    {},

	NetworkDummy: {},
}
