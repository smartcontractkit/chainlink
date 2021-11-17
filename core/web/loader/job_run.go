package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	stringutils "github.com/smartcontractkit/chainlink/core/utils/string_utils"
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
		id, err := stringutils.ToInt32(key.String())
		if err == nil {
			plnSpecIDs = append(plnSpecIDs, id)
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
		id := stringutils.FromInt32(jb.PipelineSpecID)

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
