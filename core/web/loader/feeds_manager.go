package loader

import (
	"context"
	"errors"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

type feedsBatcher struct {
	app chainlink.Application
}

func (b *feedsBatcher) loadByIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var managersIDs []int64
	for ix, key := range keys {
		id, err := stringutils.ToInt64(key.String())
		if err == nil {
			managersIDs = append(managersIDs, id)
		}
		keyOrder[key.String()] = ix
	}

	// Fetch the feeds managers
	managers, err := b.app.GetFeedsService().ListManagersByIDs(managersIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, c := range managers {
		id := stringutils.FromInt64(c.ID)

		ix, ok := keyOrder[id]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: c, Error: nil}
			delete(keyOrder, id)
		}
	}

	// fill array positions without any feeds managers
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("feeds manager not found")}
	}

	return results
}
