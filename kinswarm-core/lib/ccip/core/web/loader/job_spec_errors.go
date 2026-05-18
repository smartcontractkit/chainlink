package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

type jobSpecErrorsBatcher struct {
	app chainlink.Application
}

func (b *jobSpecErrorsBatcher) loadByJobIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var jobIDs []int32
	for ix, key := range keys {
		id, err := stringutils.ToInt32(key.String())
		if err == nil {
			jobIDs = append(jobIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	specErrors, err := b.app.JobORM().FindSpecErrorsByJobIDs(ctx, jobIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of jobIDs to spec errors
	specErrorsForJobs := map[string][]job.SpecError{}
	for _, s := range specErrors {
		jobID := stringutils.FromInt32(s.JobID)
		specErrorsForJobs[jobID] = append(specErrorsForJobs[jobID], s)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, s := range specErrorsForJobs {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: s, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any nodes as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []job.SpecError{}, Error: nil}
	}

	return results
}
