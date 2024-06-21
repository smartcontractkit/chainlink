package relay

const (
	NetworkEVM      = "evm"
	NetworkCosmos   = "cosmos"
	NetworkSolana   = "solana"
	NetworkStarkNet = "starknet"
	NetworkAptos    = "aptos"
)

var SupportedNetworks = map[string]struct{}{
	NetworkEVM:      {},
	NetworkCosmos:   {},
	NetworkSolana:   {},
	NetworkStarkNet: {},
	NetworkAptos:    {},
}
