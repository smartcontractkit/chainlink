package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type chainBatcher struct {
	app chainlink.Application
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
