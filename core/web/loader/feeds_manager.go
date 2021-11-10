package loader

import (
	"context"
	"errors"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

type feedsBatcher struct {
	app chainlink.Application
}

func (b *feedsBatcher) loadByIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var managersIDs []int64
	for ix, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err == nil {
			managersIDs = append(managersIDs, id)
		}
		keyOrder[key.String()] = ix
	}

	// Fetch the feeds managers
	managers, err := b.app.GetFeedsService().GetManagers(managersIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, c := range managers {
		id := strconv.FormatInt(c.ID, 10)

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
