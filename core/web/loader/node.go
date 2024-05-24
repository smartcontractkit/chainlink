package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type nodeBatcher struct {
	app chainlink.Application
}

func (b *nodeBatcher) loadByChainIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	// note backward compatibility -- this only ever supported evm chains
	evmrelayIDs := make([]types.RelayID, 0, len(keys))

	for ix, key := range keys {
		rid := types.RelayID{Network: types.NetworkEVM, ChainID: key.String()}
		evmrelayIDs = append(evmrelayIDs, rid)
		keyOrder[key.String()] = ix
	}

	allNodes, _, err := b.app.GetRelayers().NodeStatuses(ctx, 0, -1, evmrelayIDs...)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}
	// Generate a map of nodes to chainIDs
	nodesForChain := map[string][]types.NodeStatus{}
	for _, n := range allNodes {
		nodesForChain[n.ChainID] = append(nodesForChain[n.ChainID], n)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, ns := range nodesForChain {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: ns, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any nodes as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []types.NodeStatus{}, Error: nil}
	}

	return results
}
