package loader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

type jobRunBatcher struct {
	app chainlink.Application
}

func (b *jobRunBatcher) loadByPipelineSpecIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var plnSpecIDs []int32
	for ix, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 32)
		if err == nil {
			plnSpecIDs = append(plnSpecIDs, int32(id))
		}
		keyOrder[key.String()] = ix
	}

	// Fetch the job runs
	jbRuns, err := b.app.JobORM().PipelineRunsByJobsIDs(plnSpecIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, c := range jbRuns {
		id := strconv.FormatInt(c.ID, 10)

		ix, ok := keyOrder[id]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: c, Error: nil}
			delete(keyOrder, id)
		}
	}

	// fill array positions without any job runs
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("job run not found")}
	}

	return results
}
