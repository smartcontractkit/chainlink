package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type nodeBatcher struct {
	app chainlink.Application
}

func (b *nodeBatcher) loadByChainIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var chainIDs []utils.Big
	for ix, key := range keys {
		id := utils.Big{}
		if err := id.UnmarshalText([]byte(key.String())); err == nil {
			chainIDs = append(chainIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	nodes, err := b.app.GetChains().EVM.GetNodesByChainIDs(ctx, chainIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of nodes to chainIDs
	nodesForChain := map[string][]types.Node{}
	for _, n := range nodes {
		nodesForChain[n.EVMChainID.String()] = append(nodesForChain[n.EVMChainID.String()], n)
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
		results[ix] = &dataloader.Result{Data: []types.Node{}, Error: nil}
	}

	return results
}
