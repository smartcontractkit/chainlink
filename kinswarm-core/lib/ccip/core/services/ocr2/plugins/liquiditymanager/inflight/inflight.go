package inflight

import (
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

// Container stores transfers that are in-flight.
// Transfers are expired when they are confirmed on-chain.
type Container interface {
	// Add adds a transfer to the inflight container.
	Add(t models.Transfer)
	// Expire removes any transfers from the inflight container that are in the pending list.
	Expire(pending []models.PendingTransfer) (numExpired int)
	// GetAll returns all transfers in the inflight container.
	GetAll() []models.Transfer
	// IsInflight returns true if the transfer is in the inflight container.
	IsInflight(t models.Transfer) bool
}

// transferID uniquely identifies a transfer for a short period of time.
type transferID struct {
	From   models.NetworkSelector
	To     models.NetworkSelector
	Amount string
}

type inflight struct {
	transfers map[transferID]models.Transfer
	mu        sync.RWMutex
}

func New() *inflight {
	return &inflight{
		transfers: make(map[transferID]models.Transfer),
	}
}

func (i *inflight) Add(t models.Transfer) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.transfers[transferID{
		From:   t.From,
		To:     t.To,
		Amount: t.Amount.String(),
	}] = t
}

func (i *inflight) Expire(pending []models.PendingTransfer) int {
	i.mu.Lock()
	defer i.mu.Unlock()

	var numExpired int
	for _, p := range pending {
		k := transferID{
			From:   p.From,
			To:     p.To,
			Amount: p.Amount.String(),
		}
		t, ok := i.transfers[k]
		// This check handles scenarios where the same transfer can be in-flight multiple times.
		// This only arises in bridges that have multi-stage finalization, like Optimism, or the
		// testonly bridge, which emulates Optimism. In that case, this is the scenario that arises:
		// 1. Initial rebalance out of Optimism is sent to the chain, and the transfer is added to the inflight container.
		// 2. The transfer is confirmed on-chain, and the inflight transfer is expired. This is regular so far.
		// 3. The bridge transfer is now in a pending state, where we are waiting for it to be ready to prove on L1.
		// 4. Once the transfer is ready to prove on L1, we generate the proof data and submit it to the rebalancer on L1. This is a new inflight transfer.
		// 5. In the next round, we still notice that the transfer is ready to prove, so we return it from the bridge. However, it is not ready to remove from the
		// inflight container, because its still inflight to L1. This stage check would return false because the stage of the inflight transfer is the same
		// as the stage of the pending transfer provided.
		// 6. The prove withdrawal tx is confirmed on L1, and the bridge notices that, so the transfer is then expired because it has a higher stage.
		if ok && p.Stage > t.Stage {
			numExpired++
			delete(i.transfers, k)
		}
	}

	return numExpired
}

func (i *inflight) GetAll() []models.Transfer {
	i.mu.RLock()
	defer i.mu.RUnlock()

	transfers := make([]models.Transfer, 0, len(i.transfers))
	for k := range i.transfers {
		transfers = append(transfers, i.transfers[k])
	}

	// Sort the transfers so that they are always in the same order.
	sort.Slice(transfers, func(i, j int) bool {
		return transfers[i].From < transfers[j].From
	})

	return transfers
}

func (i *inflight) IsInflight(t models.Transfer) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	_, ok := i.transfers[transferID{
		From:   t.From,
		To:     t.To,
		Amount: t.Amount.String(),
	}]
	return ok
}
