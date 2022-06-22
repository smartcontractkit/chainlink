package relay

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
