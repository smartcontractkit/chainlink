package loader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type jobProposalSpecBatcher struct {
	app chainlink.Application
}

func (b *jobProposalSpecBatcher) loadByJobProposalsIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := keyOrderInt64(keys)

	specs, err := b.app.GetFeedsService().ListSpecsByJobProposalIDs(ids)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of specs to job proposal IDs
	specsForJP := map[string][]feeds.JobProposalSpec{}
	for _, spec := range specs {
		jpID := strconv.Itoa(int(spec.JobProposalID))
		specsForJP[jpID] = append(specsForJP[jpID], spec)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, ns := range specsForJP {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: ns, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any job proposals as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []feeds.JobProposalSpec{}, Error: nil}
	}

	return results
}
