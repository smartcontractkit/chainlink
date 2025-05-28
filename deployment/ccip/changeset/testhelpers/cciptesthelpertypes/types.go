package cciptesthelpertypes

// RoleDONTopology is used to specify the chain sub-committees of the role DON.
// A chain sub-committee is the set of nodes that supports a specific chain.
// The size of the sub-committee is determined by the fChain value assigned to the chain.
// The general formula is: NChain = 3 * fChain + 1
// Finally, the set of all sub-committees is a "topology" of the role DON. In the most
// over-specified form, this can be represented as a mapping of node P2P IDs to a slice
// of chain selectors.
//
// Note that this does NOT affect the chain support of the home chain,
// which all nodes must always support.
//
// This object is currently intended to be used primarily for tests, where
// we don't really care what specific node is assigned which specific chain,
// but rather that the overall network topology is honored. As such, the FChainToNumChains
// map field is intentionally simple, mapping fChain values to the number of chains that have that fChain value.
//
// TODO: figure out if this can be an interface instead, so that we can have different implementations. e.g.
// * exhaustively specified topology (i.e explicit node => [chain selector] mapping)
// * topology specified by some other means (e.g. a function that takes a node index and returns a chain selector based on some logic,
// maybe the solana node shouldn't also support aptos, etc.)
type RoleDONTopology struct {
	// FChainToNumChains maps (the fChain value) to the number of chains that have that fChain value.
	// Note that the sum of all of the number of chains (i.e the values in the map) must be equal exactly to the number
	// of chains set up in the test MINUS the home chain.
	// For example, if we have 2 chains, and we want to have an fChain value of 2 for the first chain,
	// and an fChain value of 1 for the second chain, we would set the map to:
	// {2: 1, 1: 1}
	FChainToNumChains map[int]int

	// Seed is a value used to seed the PRNG that will assign nodes to chains.
	// It is used to ensure that the same topology is generated across multiple runs of the test.
	Seed int64
}
