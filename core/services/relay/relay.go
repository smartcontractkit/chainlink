package relay

type Network string

var (
	EVM             Network = "evm"
	Solana          Network = "solana"
	Terra           Network = "terra"
	StarkNet        Network = "starknet"
	SupportedRelays         = map[Network]struct{}{
		EVM:      {},
		Solana:   {},
		Terra:    {},
		StarkNet: {},
	}
)
