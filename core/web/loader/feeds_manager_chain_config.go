package loader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
)

type feedsManagerChainConfigBatcher struct {
	app chainlink.Application
}

func (b *feedsManagerChainConfigBatcher) loadByManagerIDs(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
	ids, keyOrder := keyOrderInt64(keys)

	cfgs, err := b.app.GetFeedsService().ListChainConfigsByManagerIDs(ids)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of specs to job proposal IDs
	cfgsForManager := map[string][]feeds.ChainConfig{}
	for _, cfg := range cfgs {
		mgrID := strconv.Itoa(int(cfg.FeedsManagerID))
		cfgsForManager[mgrID] = append(cfgsForManager[mgrID], cfg)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, ns := range cfgsForManager {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: ns, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any job proposals as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []feeds.ChainConfig{}, Error: nil}
	}

	return results
}
