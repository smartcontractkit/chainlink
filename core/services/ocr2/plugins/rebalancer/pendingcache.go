package rebalancer

import (
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

type PendingTransfersCache struct {
	mem []models.PendingTransfer
	mu  *sync.RWMutex
}

func NewPendingTransfersCache() *PendingTransfersCache {
	return &PendingTransfersCache{
		mem: make([]models.PendingTransfer, 0),
		mu:  &sync.RWMutex{},
	}
}

func (c *PendingTransfersCache) Add(transfers []models.PendingTransfer) {
	// exclude transfers that are already present
	newTransfers := make([]models.PendingTransfer, 0, len(transfers))
	for _, tr := range transfers {
		if !c.ContainsTransfer(tr.Transfer) {
			newTransfers = append(newTransfers, tr)
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.mem = append(c.mem, newTransfers...)
}

func (c *PendingTransfersCache) Set(transfers []models.PendingTransfer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mem = transfers
}

func (c *PendingTransfersCache) ContainsTransfer(tr models.Transfer) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, pt := range c.mem {
		if pt.Transfer.Equals(tr) {
			return true
		}
	}
	return false
}

func (c *PendingTransfersCache) LatestNetworkTransfer(net models.NetworkSelector) (models.PendingTransfer, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.mem) == 0 {
		return models.PendingTransfer{}, false
	}

	var mostRecentTransfer models.PendingTransfer
	for _, tr := range c.mem {
		if tr.From == net && tr.Date.After(mostRecentTransfer.Date) {
			mostRecentTransfer = tr
		}
	}

	found := !mostRecentTransfer.Equals(models.Transfer{})
	return mostRecentTransfer, found
}
