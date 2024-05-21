package loader

import (
	"context"
	"slices"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type chainBatcher struct {
	app chainlink.Application
}

func (b *chainBatcher) loadByIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var chainIDs []string
	for ix, key := range keys {
		chainIDs = append(chainIDs, key.String())
		keyOrder[key.String()] = ix
	}

	var cs []types.ChainStatus
	relayers := b.app.GetRelayers().Slice()

	for _, r := range relayers {
		s, err := r.GetChainStatus(ctx)
		if err != nil {
			return []*dataloader.Result{{Data: nil, Error: err}}
		}

		if slices.Contains(chainIDs, s.ID) {
			cs = append(cs, s)
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for _, c := range cs {
		ix, ok := keyOrder[c.ID]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: c, Error: nil}
			delete(keyOrder, c.ID)
		}
	}

	// fill array positions without any nodes
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: chains.ErrNotFound}
	}

	return results
}
