package loader

import (
	"context"
	"slices"

	"github.com/graph-gophers/dataloader"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type chainBatcher struct {
	app chainlink.Application
}

// Deprecated: use loadByChainIDs is deprecated and we should be using loadByRelayIDs.
func (b *chainBatcher) loadByIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var chainIDs []string
	for ix, key := range keys {
		chainIDs = append(chainIDs, key.String())
		keyOrder[key.String()] = ix
	}

	var cs []chainlink.NetworkChainStatus
	relayersMap := b.app.GetRelayers().GetIDToRelayerMap()

	for k, v := range relayersMap {
		s, err := v.GetChainStatus(ctx)
		if err != nil {
			return []*dataloader.Result{{Data: nil, Error: err}}
		}

		if slices.Contains(chainIDs, s.ID) {
			cs = append(cs, chainlink.NetworkChainStatus{
				ChainStatus: s,
				Network:     k.Network,
			})
		}
	}

	// todo: future improvements to handle multiple chains with same id
	if len(cs) > len(keys) {
		b.app.GetLogger().Warn("Found multiple chain with same id")
		return []*dataloader.Result{{Data: nil, Error: chains.ErrMultipleChainFound}}
	}

	results := make([]*dataloader.Result, len(keys))
	for _, c := range cs {
		ix, ok := keyOrder[c.ID]
		// if found, remove from index lookup map, so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: c, Error: nil}
			delete(keyOrder, c.ID)
		}
	}

	// fill array positions without any nodes
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: nil, Error: chains.ErrNotFound}
	}

	return results
}

func (b *chainBatcher) loadByRelayIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, 0, len(keys))
	for _, key := range keys {
		var relay commonTypes.RelayID
		err := relay.UnmarshalString(key.String())
		if err != nil {
			results = append(results, &dataloader.Result{Data: nil, Error: err})
			continue
		}

		relayer, err := b.app.GetRelayers().Get(relay)
		if err != nil {
			results = append(results, &dataloader.Result{Data: nil, Error: chains.ErrNotFound})
			continue
		}

		status, err := relayer.GetChainStatus(ctx)
		if err != nil {
			results = append(results, &dataloader.Result{Data: nil, Error: err})
			continue
		}

		results = append(results, &dataloader.Result{Data: chainlink.NetworkChainStatus{
			ChainStatus: status,
			Network:     relay.Network,
		}, Error: err})
	}

	return results
}
