package vrfcommon

import (
	"math/big"
	"sync"
)

const (
	// MaxForceFulfillments is the maximum number of force fulfillments that we will
	// perform on a subscription.
	MaxForceFulfillments = 10

	// PruneInterval is the block interval at which the force fulfillment rate limiter
	// will prune it's cache of force-fulfill counts per sub.
	PruneInterval = 10_000
)

type ForceFulfillRateLimiter struct {
	// forceFulfillsCount is a mapping from subscription id (as a string)
	// to the number of force fulfillments performed for that sub.
	forceFulfillsCount map[string]int
	mu                 sync.Mutex
	latestHead         uint64
}

func NewForceFulfillRateLimiter() *ForceFulfillRateLimiter {
	return &ForceFulfillRateLimiter{
		forceFulfillsCount: make(map[string]int),
		mu:                 sync.Mutex{},
	}
}

// ShouldFulfill returns true if and only if a force fulfillment should be performed
// for the provided sub id.
func (f *ForceFulfillRateLimiter) ShouldFulfill(subId *big.Int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.forceFulfillsCount[subId.String()] < MaxForceFulfillments
}

func (f *ForceFulfillRateLimiter) NumFulfilled(subId *big.Int) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.forceFulfillsCount[subId.String()]
}

// FulfillmentPerformed indicates that a force-fulfillment was performed for the
// given sub id.
func (f *ForceFulfillRateLimiter) FulfillmentPerformed(subId *big.Int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.forceFulfillsCount[subId.String()]++
}

// SetLatestHead sets the latest head in the rate limiter to the provided value
// and prunes the cache if necessary.
func (f *ForceFulfillRateLimiter) SetLatestHead(head uint64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if head < f.latestHead {
		// may happen on re-orgs
		return
	}
	oldLatest := f.latestHead
	f.latestHead = head
	if (f.latestHead - oldLatest) >= PruneInterval {
		f.prune()
	}
}

// assumes that the caller holds the lock already
func (f *ForceFulfillRateLimiter) prune() {
	f.forceFulfillsCount = make(map[string]int)
}
