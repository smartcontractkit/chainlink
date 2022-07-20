package loader

import (
	"context"
	"errors"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type chainBatcher struct {
	app chainlink.Application
}

func (b *chainBatcher) loadByIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
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

	// Fetch the chains
	chains, err := b.app.EVMORM().GetChainsByIDs(chainIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, c := range chains {
		ix, ok := keyOrder[c.ID.String()]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: c, Error: nil}
			delete(keyOrder, c.ID.String())
		}
	}

	// fill array positions without any nodes
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("chain not found")}
	}

	return results
}
