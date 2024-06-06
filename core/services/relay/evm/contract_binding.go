package evm

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type filterRegisterer struct {
	pollingFilter logpoller.Filter
	filterLock    sync.Mutex
	// registerCalled is used to determine if Register was called during Chain Reader service Start.
	// This is done to avoid calling Register while the service is not running because log poller is most likely also not running.
	registerCalled bool
}

// contractBinding stores read bindings and manages the common contract event filter.
type contractBinding struct {
	name string
	// filterRegisterer is used to manage polling filter registration for the common contract filter.
	// The common contract filter should be used by events that share filtering args.
	filterRegisterer
	// key is read name method, event or event keys used for queryKey.
	readBindings map[string]readBinding
	// bound determines if address is set to the contract binding.
	bound bool
}

// Bind binds contract addresses to contract binding and registers the common contract polling filter.
func (cb *contractBinding) Bind(ctx context.Context, lp logpoller.LogPoller, boundContract commontypes.BoundContract) error {
	if cb.bound {
		// we are changing contract address reference, so we need to unregister old filter it exists
		if err := cb.Unregister(ctx, lp); err != nil {
			return err
		}
	}

	cb.pollingFilter.Addresses = evmtypes.AddressArray{common.HexToAddress(boundContract.Address)}
	cb.pollingFilter.Name = logpoller.FilterName(boundContract.Name+"."+uuid.NewString(), boundContract.Address)
	cb.bound = true

	if cb.registerCalled {
		return cb.Register(ctx, lp)
	}

	return nil
}

// Register registers the common contract filter.
func (cb *contractBinding) Register(ctx context.Context, lp logpoller.LogPoller) error {
	cb.registerCalled = true
	if !cb.bound {
		return nil
	}

	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if len(cb.pollingFilter.EventSigs) > 0 && !lp.HasFilter(cb.pollingFilter.Name) {
		if err := lp.RegisterFilter(ctx, cb.pollingFilter); err != nil {
			return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
		}
	}

	return nil
}

// Unregister unregisters the common contract filter.
func (cb *contractBinding) Unregister(ctx context.Context, lp logpoller.LogPoller) error {
	if !cb.bound {
		return nil
	}

	cb.filterLock.Lock()
	defer cb.filterLock.Unlock()

	if !lp.HasFilter(cb.pollingFilter.Name) {
		return nil
	}

	if err := lp.UnregisterFilter(ctx, cb.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return nil
}
