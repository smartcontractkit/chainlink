package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type nodeBatcher struct {
	app chainlink.Application
}

func (b *nodeBatcher) loadByChainIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var ids []string
	for ix, key := range keys {
		ids = append(ids, key.String())
		keyOrder[key.String()] = ix
	}

	nodes, _, err := b.app.GetChains().EVM.NodeStatuses(ctx, 0, -1, ids...)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of nodes to chainIDs
	nodesForChain := map[string][]types.NodeStatus{}
	for _, n := range nodes {
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
