package ethgo

// Network is a chain id
type Network uint64

const (
	// Mainnet is the mainnet network
	Mainnet Network = 1

	// Ropsten is the POW testnet
	Ropsten Network = 3

	// Rinkeby is a POW testnet
	Rinkeby Network = 4

	// Goerli is the Clique testnet
	Goerli Network = 5
)
