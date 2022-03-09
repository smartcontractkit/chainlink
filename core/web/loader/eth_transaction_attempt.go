package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

type ethTransactionAttemptBatcher struct {
	app chainlink.Application
}

func (b *ethTransactionAttemptBatcher) loadByEthTransactionIDs(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// Create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// Collect the keys to search for
	var ethTxsIDs []int64
	for ix, key := range keys {
		id, err := stringutils.ToInt64(key.String())
		if err == nil {
			ethTxsIDs = append(ethTxsIDs, id)
		}

		keyOrder[key.String()] = ix
	}

	attempts, err := b.app.TxmORM().FindEthTxAttemptsByEthTxIDs(ethTxsIDs)
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}

	// Generate a map of attempts to txIDs
	attemptsForTx := map[string][]txmgr.EthTxAttempt{}
	for _, a := range attempts {
		id := stringutils.FromInt64(a.EthTxID)

		attemptsForTx[id] = append(attemptsForTx[id], a)
	}

	// Construct the output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	for k, ns := range attemptsForTx {
		ix, ok := keyOrder[k]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result{Data: ns, Error: nil}
			delete(keyOrder, k)
		}
	}

	// fill array positions without any attempts as an empty slice
	for _, ix := range keyOrder {
		results[ix] = &dataloader.Result{Data: []txmgr.EthTxAttempt{}, Error: nil}
	}

	return results
}
