package evm

import (
	"context"
	"fmt"
	"sync"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

type filterRegisterer struct {
	pollingFilter logpoller.Filter
	filterLock    sync.Mutex
	isRegistered  bool
}

// contractBinding stores read bindings and manages the common contract event filter.
type contractBinding struct {
	// filterRegisterer is used to manage polling filter registration for the common contract filter.
	// The common contract filter should be used by events that share filtering args.
	filterRegisterer
	// key is read name method, event or event keys used for queryKey.
	readBindings map[string]readBinding
	bound        bool
}

// Register registers the common contract filter.
func (cb *contractBinding) Register(ctx context.Context, lp logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if cb.bound && len(cb.pollingFilter.EventSigs) > 0 && !lp.HasFilter(cb.pollingFilter.Name) {
		if err := lp.RegisterFilter(ctx, cb.pollingFilter); err != nil {
			return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
		}
		cb.isRegistered = true
	}

	return nil
}

// Unregister unregisters the common contract filter.
func (cb *contractBinding) Unregister(ctx context.Context, lp logpoller.LogPoller) error {
	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if !lp.HasFilter(cb.pollingFilter.Name) {
		return nil
	}

	if err := lp.UnregisterFilter(ctx, cb.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	cb.isRegistered = false

	return nil
}
