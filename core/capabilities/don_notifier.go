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
		close(n.ch)
	}
	n.notified = true
}

func (n *DonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	select {
	case <-ctx.Done():
		return capabilities.DON{}, ctx.Err()
	case <-n.ch:
	}
	<-n.ch
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.don, nil
}
