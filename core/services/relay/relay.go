package relay

type Network string

var (
	EVM             Network = "evm"
	Solana          Network = "solana"
	StarkNet        Network = "starknet"
	SupportedRelays         = map[Network]struct{}{
		EVM:      {},
		Solana:   {},
		StarkNet: {},
	}
)
