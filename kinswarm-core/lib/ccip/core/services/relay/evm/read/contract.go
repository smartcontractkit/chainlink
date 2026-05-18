package read

import (
	"context"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

// contractBinding stores read bindings and manages the common contract event filter.
type contractBinding struct {
	// filterRegistrar is used to manage polling filter registration for the common contract filter.
	// The common contract filter should be used by events that share filtering args.
	registrar *syncedFilter

	// registered is used to determine if Register was called during Chain Reader service Start.
	// This is done to avoid calling Register while the service is not running because log poller is most likely also not running.
	registerCalled bool

	// internal properties
	name    string
	readers map[string]Reader       // key is read name method, event or event keys used for queryKey.
	bound   map[common.Address]bool // bound determines if address is set to the contract binding.
	mu      sync.RWMutex
}

func newContractBinding(name string) *contractBinding {
	return &contractBinding{
		name:      name,
		readers:   make(map[string]Reader),
		bound:     make(map[common.Address]bool),
		registrar: newSyncedFilter(),
	}
}

// GetReaderNamed returns a reader for the provided contract name. This method is thread safe.
func (cb *contractBinding) GetReaderNamed(name string) (Reader, bool) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	binding, exists := cb.readers[name]

	return binding, exists
}

// AddReaderNamed adds a new reader to the contract bindings for the provided contract name. This
// method is thread safe.
func (cb *contractBinding) AddReaderNamed(name string, rdr Reader) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.readers[name] = rdr
}

// Bind binds contract addresses to contract binding and registers the common contract polling filter.
func (cb *contractBinding) Bind(ctx context.Context, registrar Registrar, bindings ...common.Address) error {
	if cb.isBound() {
		if err := cb.Unregister(ctx, registrar); err != nil {
			return err
		}
	}

	for _, binding := range bindings {
		if cb.bindingExists(binding) {
			continue
		}

		cb.registrar.SetName(logpoller.FilterName(cb.name + "." + uuid.NewString()))
		cb.registrar.AddAddress(binding)
		cb.addBinding(binding)
	}

	// registerCalled during ChainReader start
	if cb.registered() {
		return cb.Register(ctx, registrar)
	}

	return nil
}

func (cb *contractBinding) BindReaders(ctx context.Context, addresses ...common.Address) error {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	var err error

	for _, reader := range cb.readers {
		err = errors.Join(err, reader.Bind(ctx, addresses...))
	}

	return err
}

// Unbind unbinds contract addresses from contract binding and unregisters the common contract polling filter.
func (cb *contractBinding) Unbind(ctx context.Context, registrar Registrar, bindings ...common.Address) error {
	for _, binding := range bindings {
		if !cb.bindingExists(binding) {
			continue
		}

		cb.registrar.RemoveAddress(binding)
		cb.removeBinding(binding)
	}

	// we are changing contract address reference, so we need to unregister old filter or re-register existing filter
	if cb.registrar.Count() == 0 {
		cb.registrar.SetName("")

		return cb.Unregister(ctx, registrar)
	} else if cb.registered() {
		return cb.Register(ctx, registrar)
	}

	return nil
}

func (cb *contractBinding) UnbindReaders(ctx context.Context, addresses ...common.Address) error {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	var err error

	for _, reader := range cb.readers {
		err = errors.Join(reader.Unbind(ctx, addresses...))
	}

	return err
}

func (cb *contractBinding) SetCodecAll(codec commontypes.RemoteCodec) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	for _, binding := range cb.readers {
		binding.SetCodec(codec)
	}
}

// Register registers the common contract filter.
func (cb *contractBinding) Register(ctx context.Context, registrar Registrar) error {
	cb.setRegistered()

	if !cb.isBound() {
		return nil
	}

	if cb.registrar.HasEventSigs() {
		if err := cb.registrar.Register(ctx, registrar); err != nil {
			return err
		}
	}

	return nil
}

func (cb *contractBinding) RegisterReaders(ctx context.Context) error {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	for _, binding := range cb.readers {
		if err := binding.Register(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Unregister unregisters the common contract filter.
func (cb *contractBinding) Unregister(ctx context.Context, registrar Registrar) error {
	if !cb.isBound() {
		return nil
	}

	return cb.registrar.Unregister(ctx, registrar)
}

func (cb *contractBinding) UnregisterReaders(ctx context.Context) error {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	for _, binding := range cb.readers {
		if err := binding.Unregister(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (cb *contractBinding) SetFilter(filter logpoller.Filter) {
	cb.registrar.SetFilter(filter)
}

func (cb *contractBinding) isBound() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return len(cb.bound) > 0
}

func (cb *contractBinding) bindingExists(binding common.Address) bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	bound, exists := cb.bound[binding]

	return exists && bound
}

func (cb *contractBinding) addBinding(binding common.Address) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.bound[binding] = true
}

func (cb *contractBinding) removeBinding(binding common.Address) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	delete(cb.bound, binding)
}

func (cb *contractBinding) registered() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return cb.registerCalled
}

func (cb *contractBinding) setRegistered() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.registerCalled = true
}
