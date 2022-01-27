package loader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
)

type jobProposalBatcher struct {
	app chainlink.Application
}

func (b *jobProposalBatcher) loadByManagersIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var mgrsIDs []int64
	for ix, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err == nil {
			mgrsIDs = append(mgrsIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	jps, err := b.app.GetFeedsService().ListJobProposalsByManagersIDs(mgrsIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of job proposals to feeds managers IDs
	jpsForMgr := map[string][]feeds.JobProposal{}
	for _, jp := range jps {
		mgrID := strconv.Itoa(int(jp.FeedsManagerID))
		jpsForMgr[mgrID] = append(jpsForMgr[mgrID], jp)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, ns := range jpsForMgr {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: ns, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any job proposals as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []feeds.JobProposal{}, Error: nil}
	}

	return results
}
