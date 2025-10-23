package capabilities

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type DonNotifier struct {
	mu       sync.Mutex
	don      capabilities.DON
	notified bool
	ch       chan struct{}
}

func NewDonNotifier() *DonNotifier {
	return &DonNotifier{
		ch: make(chan struct{}),
	}
}

func (n *DonNotifier) NotifyDonSet(don capabilities.DON) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.don = don
	if !n.notified {
		n.notified = true
		close(n.ch)
	}
}

func (n *DonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	// First, check if we already have the DON without blocking
	n.mu.Lock()
	if n.notified {
		defer n.mu.Unlock()
		return n.don, nil
	}
	n.mu.Unlock()

	// Otherwise, wait for notification or context cancellation
	select {
	case <-ctx.Done():
		return capabilities.DON{}, ctx.Err()
	case <-n.ch:
		n.mu.Lock()
		defer n.mu.Unlock()
		return n.don, nil
	}
}
