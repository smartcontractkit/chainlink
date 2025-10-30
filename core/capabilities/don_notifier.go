package capabilities

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type DonNotifier interface {
	// Notify all subscribers of a new DON without blocking for a subscriber.
	NotifyDonSet(don capabilities.DON)
}

type DonSubscriber interface {
	// Subscribe returns a channel that will receive the latest DON.  Unsubscribe
	// by calling the returned function.
	Subscribe(ctx context.Context) (<-chan capabilities.DON, func(), error)
}

// DonNotifyWaitSubscriber handles the lifecyle of a Workflow DON update.  A node may
// only belong to a single workflow DON, but multiple capabilities DONs.  In practice,
// this interface is used to update subscribers with the current workflow DON
// state.
type DonNotifyWaitSubscriber interface {
	DonNotifier
	DonSubscriber

	// Block until a new DON is received or the context is canceled.  The current
	// DON, if set, is returned immediately.
	WaitForDon(ctx context.Context) (capabilities.DON, error)
}

type donNotifier struct {
	don         atomic.Pointer[capabilities.DON]
	subscribers sync.Map
}

func NewDonNotifier() *donNotifier {
	return &donNotifier{}
}

func (n *donNotifier) NotifyDonSet(don capabilities.DON) {
	n.don.Store(&don)

	// Broadcast the new DON to all subscriber channels.
	n.subscribers.Range(func(key, _ any) bool {
		s := key.(chan capabilities.DON)

		select {
		case s <- don:
		default:
		}

		return true
	})
}

// Subscribe returns a listen only channel that will return the latest value
// state of the workflow DON for this node until the cleanup function is called.
// The current state is buffered into the returned channel.
func (n *donNotifier) Subscribe(ctx context.Context) (<-chan capabilities.DON, func(), error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	// Buffered so as not to block.
	s := make(chan capabilities.DON, 1)
	unsubscribe := func() {
		n.subscribers.Delete(s)
	}

	n.subscribers.Store(s, struct{}{})

	if n.don.Load() != nil {
		s <- *n.don.Load()
	}

	return s, unsubscribe, nil
}

func (n *donNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	if n.don.Load() != nil {
		return *n.don.Load(), nil
	}

	subCh, unsubscribe, err := n.Subscribe(ctx)
	if err != nil {
		return capabilities.DON{}, fmt.Errorf("failed to subscribe to DON updates: %w", err)
	}
	defer unsubscribe()

	select {
	case <-ctx.Done():
		return capabilities.DON{}, fmt.Errorf("failed to wait for don: %w", ctx.Err())
	case don := <-subCh:
		return don, nil
	}
}
