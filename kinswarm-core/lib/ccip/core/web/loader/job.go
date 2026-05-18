package loader

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
)

type jobBatcher struct {
	app chainlink.Application
}

func (b *jobBatcher) loadByExternalJobIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var jobIDs []uuid.UUID
	for ix, key := range keys {
		id, err := uuid.Parse(key.String())
		if err == nil {
			jobIDs = append(jobIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	// Fetch the jobs
	var jobs []job.Job
	for _, id := range jobIDs {
		job, err := b.app.JobORM().FindJobByExternalJobID(ctx, id)

		if err != nil {
			return []*dataloader.Result{{Data: nil, Error: err}}
		}

		jobs = append(jobs, job)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, j := range jobs {
		id := j.ExternalJobID.String()

		ix, ok := keyOrder[id]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: j, Error: nil}
			delete(keyOrder, id)
		}
	}

	// fill array positions without any feeds managers
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("feeds manager not found")}
	}

	return results
}

func (b *jobBatcher) loadByPipelineSpecIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var plSpecIDs []int32
	for ix, key := range keys {
		id, err := stringutils.ToInt32(key.String())
		if err == nil {
			plSpecIDs = append(plSpecIDs, id)
		}
		keyOrder[key.String()] = ix
	}

	// Fetch the jobs
	jobs, err := b.app.JobORM().FindJobsByPipelineSpecIDs(ctx, plSpecIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for _, j := range jobs {
		id := stringutils.FromInt32(j.PipelineSpecID)

		ix, ok := keyOrder[id]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: j, Error: nil}
			delete(keyOrder, id)
		}
	}

	// fill array positions without any jobs
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: errors.New("job not found")}
	}

	return results
}
