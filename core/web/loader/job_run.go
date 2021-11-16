package loader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type jobRunBatcher struct {
	app chainlink.Application
}

func (b *jobRunBatcher) loadByPipelineSpecIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
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

	// Generate a map of pipeline runs to pipeline spec id
	runsForJob := map[string][]pipeline.Run{}
	for _, jb := range jbRuns {
		id := strconv.Itoa(int(jb.PipelineSpecID))

		runsForJob[id] = append(runsForJob[id], jb)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, rs := range runsForJob {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: rs, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any job runs
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("job run not found")}
	}

	return results
}
