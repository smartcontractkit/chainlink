package cciptesthelpertypes

import (
	"errors"
	"fmt"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// NewDistributedTopology creates a new rotating topology with the given arguments.
func NewDistributedTopology(args DistributedTopologyArgs) RoleDONTopology {
	return &topology{
		impl: &DistributedTopology{DistributedTopologyArgs: args},
	}
}

type DistributedTopologyArgs struct {
	FValues map[cciptypes.ChainSelector]uint8
}

type DistributedTopology struct {
	DistributedTopologyArgs
}

// check if every chain selector has a defined 'f' value
func (t *DistributedTopologyArgs) validate(nonHomeChainSelectors []cciptypes.ChainSelector) error {
	if len(t.FValues) == 0 {
		return errors.New("the topology must have at least one f value defined")
	}

	for _, chainSelector := range nonHomeChainSelectors {
		if _, ok := t.FValues[chainSelector]; !ok {
			return fmt.Errorf("the topology must have a f value defined for every non-home chain selector, missing fvalue for %s", chainSelector)
		}
	}

	return nil
}

func (t *DistributedTopologyArgs) ChainToNodeMapping(nonBootstrapP2pIDs [][32]byte, chainSelectors []cciptypes.ChainSelector, homeChainSelector cciptypes.ChainSelector) (map[cciptypes.ChainSelector][][32]byte, error) {
	if len(nonBootstrapP2pIDs) < MinRoleDONSize {
		return nil, fmt.Errorf("the number of non-bootstrap ccip nodes must be at least %d, got %d",
			MinRoleDONSize,
			len(nonBootstrapP2pIDs),
		)
	}

	if err := t.validate(chainSelectors); err != nil {
		return nil, err
	}

	chainToNodeMapping := make(map[cciptypes.ChainSelector][][32]byte)
	chainToNodeMapping[homeChainSelector] = nonBootstrapP2pIDs

	// rotate through all nodes to ensure each node has equal load across all chains
	currIndex := 0
	for _, chainSelector := range chainSelectors {
		numNodes := NChain(t.FValues[chainSelector])
		for i := 0; i < int(numNodes); i++ {
			if currIndex >= len(nonBootstrapP2pIDs) {
				currIndex = 0 // wrap around if we reach the end of the list
			}
			chainToNodeMapping[chainSelector] = append(chainToNodeMapping[chainSelector], nonBootstrapP2pIDs[currIndex])
			currIndex++
		}
	}
	return chainToNodeMapping, nil
}
