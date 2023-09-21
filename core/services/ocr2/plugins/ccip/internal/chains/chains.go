package chains

type EVMChain uint64

const (
	Ethereum  EVMChain = 1
	Optimism  EVMChain = 10
	Arbitrum  EVMChain = 42161
	Avalanche EVMChain = 43114

	GoerliTestnet         EVMChain = 5
	OptimismGoerliTestnet EVMChain = 420
	AvalancheFujiTestnet  EVMChain = 43113
	ArbitrumGoerliTestnet EVMChain = 421613
)
