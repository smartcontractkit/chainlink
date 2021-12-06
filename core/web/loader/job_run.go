package loader

import (
	"context"
	"errors"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

type jobRunBatcher struct {
	app chainlink.Application
}

func (b *jobRunBatcher) loadByIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var runIDs []int64
	for ix, key := range keys {
		id, err := stringutils.ToInt64(key.String())
		if err == nil {
			runIDs = append(runIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	// Fetch the runs
	runs, err := b.app.JobORM().FindPipelineRunsByIDs(runIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, r := range runs {
		idStr := stringutils.FromInt64(r.ID)

		ix, ok := keyOrder[idStr]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: r, Error: nil}
			delete(keyOrder, idStr)
		}
	}

	// fill array positions without any job runs
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("run not found")}
	}

	return results
}
