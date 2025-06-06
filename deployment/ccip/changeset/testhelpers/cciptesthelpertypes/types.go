package cciptesthelpertypes

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// NewRandomTopology creates a new random topology with the given arguments.
func NewRandomTopology(args RandomTopologyArgs) RoleDONTopology {
	return &topology{
		impl: &randomTopology{RandomTopologyArgs: args},
	}
}

var _ RoleDONTopology = &topology{}

// topology wraps a RoleDONTopology implementation and ensures that the home chain is supported by all nodes.
// See the RoleDONTopology interface for more details.
type topology struct {
	impl RoleDONTopology
}

// ChainToNodeMapping implements RoleDONTopology.
func (t *topology) ChainToNodeMapping(nonBootstrapP2pIDs [][32]byte, nonHomeChainSelectors []cciptypes.ChainSelector, homeChainSelector cciptypes.ChainSelector) (map[cciptypes.ChainSelector][][32]byte, error) {
	chainToNodeMapping, err := t.impl.ChainToNodeMapping(nonBootstrapP2pIDs, nonHomeChainSelectors, homeChainSelector)
	if err != nil {
		return nil, err
	}

	// ensure that the home chain is supported by all nodes.
	chainToNodeMapping[homeChainSelector] = nonBootstrapP2pIDs

	return chainToNodeMapping, nil
}

// MinRoleDONSize represents the size of the smallest possible CCIP role DON.
const MinRoleDONSize = 4

// NChain returns the number of readers for a chain based on a given fChain value.
// The general formula is: NChain = 3 * fChain + 1
func NChain[T constraints.Integer](fChain T) T {
	return 3*fChain + 1
}

// RoleDONTopology is used to specify the chain sub-committees of the role DON.
// A chain sub-committee is the set of nodes that supports a specific chain.
// The size of the sub-committee is determined by the fChain value assigned to the chain.
// The general formula is: NChain = 3 * fChain + 1
// Finally, the set of all sub-committees is a "topology" of the role DON. In the most
// over-specified form, this can be represented as a mapping of node P2P IDs to a slice
// of chain selectors.
//
// Note that this does NOT affect the chain support of the home chain,
// which all nodes must always support. This should be ensured by all implementations
// of RoleDONTopology, otherwise fundamental assumptions about the role don will be violated.
//
// Examples of different implementations:
// * randomized topology, where the nodes are randomly assigned to chains, while honoring the fChain value of each chain.
// * exhaustively specified topology (i.e explicit node => [chain selector] mapping), informed by some JSON file after doing
// some analysis on node operator performance.
// * topology specified by some other means, e.g. a function that takes a node index and returns a chain selector based on some logic,
// maybe the nodes that support solana shouldn't also support aptos, etc.
// * maybe we want to impose a max # of chains on nodes, i.e. to cap costs per node.
type RoleDONTopology interface {
	// ChainToNodeMapping returns a mapping of chain selectors to the set of nodes that support that chain.
	// This can be plugged in directly into the CCIPHome ChainConfig.readers value for the chain.
	// the nonBootstrapP2pIDs array is the set of all non-bootstrap ccip nodes in the role don.
	// nonHomeChainSelectors is the set of all non-home chain selectors.
	// homeChainSelector is the home chain selector.
	// All implementations of RoleDONTopology must ensure that the home chain is supported by all nodes,
	// i.e mapping[homeChainSelector] = all non-bootstrap ccip nodes.
	ChainToNodeMapping(
		nonBootstrapP2pIDs [][32]byte,
		nonHomeChainSelectors []cciptypes.ChainSelector,
		homeChainSelector cciptypes.ChainSelector,
	) (map[cciptypes.ChainSelector][][32]byte, error)
}

var _ RoleDONTopology = &randomTopology{}

// RandomTopologyArgs are the arguments for the random topology.
type RandomTopologyArgs struct {
	// FChainToNumChains maps (the fChain value) to the number of chains that have that fChain value.
	// Note that the sum of all of the number of chains (i.e the values in the map) must be equal exactly to the number
	// of chains set up in the test MINUS the home chain.
	// For example, if we have 2 chains, and we want to have an fChain value of 2 for the first chain,
	// and an fChain value of 1 for the second chain, we would set the map to:
	// {2: 1, 1: 1}
	FChainToNumChains map[int]int
	// Seed is a value used to seed the PRNG that will assign nodes to chains.
	// It is used to ensure that the same topology is generated across multiple runs of the topology.
	Seed int64
}

// randomTopology generates a random chain->node mapping, while honoring the fChain value of each chain.
// This object is currently intended to be used primarily for tests, where
// we don't really care what specific node is assigned which specific chain,
// but rather that the overall network topology is honored.
// It can also be used in live ccip test environments to generate a random topology,
// such as staging and production-testnet.
type randomTopology struct {
	RandomTopologyArgs
}

// Validate implements RoleDONTopology.
// The chainSelectors must not include the home chain, as it is supported by all nodes.
func (r *randomTopology) validate(chainSelectors []cciptypes.ChainSelector, homeChainSelector cciptypes.ChainSelector) error {
	totalChains := 0
	for _, numChains := range r.FChainToNumChains {
		totalChains += numChains
	}
	if totalChains != len(chainSelectors) {
		return fmt.Errorf(
			"the sum of the number of chains in the chain topology must be equal to the number of chains provided in the config, got %d, expected %d",
			totalChains,
			len(chainSelectors),
		)
	}

	// ensure the home chain is not in the chainSelectors.
	if slices.Contains(chainSelectors, homeChainSelector) {
		return fmt.Errorf("the home chain selector %d is included in the chainSelectors, please remove it from the chainSelectors", homeChainSelector)
	}

	return nil
}

// ChainToNodeMapping implements RoleDONTopology by randomly sampling the set of
// nodes that support each chain, while honoring the fChain value of each chain.
// The chainSelectors must not include the home chain, as it is supported by all nodes.
func (r *randomTopology) ChainToNodeMapping(
	nonBootstrapP2pIDs [][32]byte,
	chainSelectors []cciptypes.ChainSelector,
	homeChainSelector cciptypes.ChainSelector,
) (map[cciptypes.ChainSelector][][32]byte, error) {
	if len(nonBootstrapP2pIDs) < MinRoleDONSize {
		return nil, fmt.Errorf("the number of non-bootstrap ccip nodes must be at least %d, got %d",
			MinRoleDONSize,
			len(nonBootstrapP2pIDs),
		)
	}

	if err := r.validate(chainSelectors, homeChainSelector); err != nil {
		return nil, err
	}

	// allocate for each chain in chainSelectors an fChain value, using the randomness seed provided.
	gen := rand.New(rand.NewSource(r.Seed))
	chainToFChain, err := r.getChainToFChainMapping(gen, chainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain to fChain mapping: %w", err)
	}

	chainToNodeMapping := make(map[cciptypes.ChainSelector][][32]byte)
	for _, chainSelector := range chainSelectors {
		nodes, err := r.getNodesForChain(gen, chainToFChain[chainSelector], nonBootstrapP2pIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get nodes for chain %s: %w", chainSelector, err)
		}

		chainToNodeMapping[chainSelector] = nodes
	}

	// the home chain is supported by all nodes.
	chainToNodeMapping[homeChainSelector] = nonBootstrapP2pIDs

	return chainToNodeMapping, nil
}

// getNodesForChain draws a set of nodes that will support a given fChain value.
// Due to the randomness involved in the setup, a node might end up supporting multiple chains.
// However, this setup ensures that at most 3 * fChain + 1 nodes will support a given chain.
func (r *randomTopology) getNodesForChain(gen *rand.Rand, fChain int, nonBootstrapP2pIDs [][32]byte) ([][32]byte, error) {
	numNodesToDraw := NChain(fChain)

	// should never happen but a good sanity check.
	if len(nonBootstrapP2pIDs) < numNodesToDraw {
		return nil, fmt.Errorf("the number of non-bootstrap ccip nodes must be at least %d, got %d",
			numNodesToDraw,
			len(nonBootstrapP2pIDs),
		)
	}

	nodes := make([][32]byte, 0, numNodesToDraw)

	alreadyDrawn := make(map[[32]byte]bool)
	for len(nodes) < numNodesToDraw {
		idx := gen.Intn(len(nonBootstrapP2pIDs))
		p2pID := nonBootstrapP2pIDs[idx]
		if !alreadyDrawn[p2pID] {
			nodes = append(nodes, p2pID)
			alreadyDrawn[p2pID] = true
		}
	}

	return nodes, nil
}

func (r *randomTopology) getChainToFChainMapping(gen *rand.Rand, chainSelectors []cciptypes.ChainSelector) (map[cciptypes.ChainSelector]int, error) {
	chainToFChain := make(map[cciptypes.ChainSelector]int)

	// for deterministic looping over the fChainToNumChains map.
	keys := maps.Keys(r.FChainToNumChains)
	sort.Ints(keys)

	var fChains []int
	for _, fChain := range keys {
		numChains := r.FChainToNumChains[fChain]
		for range numChains {
			fChains = append(fChains, fChain)
		}
	}

	// this is checked in the validate function, but we'll check it again here for good measure.
	if len(fChains) != len(chainSelectors) {
		return nil, fmt.Errorf("the number of fChains must be equal to the number of chainSelectors, len(fChains) = %d, len(chainSelectors) = %d", len(fChains), len(chainSelectors))
	}

	// shuffle the fChains to ensure that the order is random.
	gen.Shuffle(len(fChains), func(i, j int) { fChains[i], fChains[j] = fChains[j], fChains[i] })

	// assign the fChains to the chainSelectors to get a randomized assignment.
	for i, chainSelector := range chainSelectors {
		chainToFChain[chainSelector] = fChains[i]
	}

	return chainToFChain, nil
}
