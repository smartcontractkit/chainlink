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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mem = append(c.mem, transfers...)
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
