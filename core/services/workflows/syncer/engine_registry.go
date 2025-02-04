package syncer

import (
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

const defaultGlobalCountLimit = 50

type EngineRegistryConfig struct {
	GlobalCountLimit uint
}

type EngineRegistry struct {
	engines          map[string]services.Service
	mu               sync.RWMutex
	globalCountLimit uint
}

func NewEngineRegistry(opts ...func(*EngineRegistry)) *EngineRegistry {
	er := &EngineRegistry{
		engines:          make(map[string]services.Service),
		globalCountLimit: defaultGlobalCountLimit,
	}

	for _, o := range opts {
		o(er)
	}

	return er
}

func WithGlobalCountLimit(gcl uint) func(*EngineRegistry) {
	return func(e *EngineRegistry) {
		e.globalCountLimit = gcl
	}
}

// Add adds an engine to the registry.
func (r *EngineRegistry) Add(id string, engine services.Service) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, found := r.engines[id]; found {
		return errors.New("attempting to register duplicate engine")
	}

	// validate if adding another engine would exceed the global count limit
	if len(r.engines)+1 > int(r.globalCountLimit) {
		return errors.New("global engine count limit reached")
	}
	r.engines[id] = engine
	return nil
}

// Get retrieves an engine from the registry.
func (r *EngineRegistry) Get(id string) (services.Service, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[id]
	if !found {
		return nil, errors.New("engine not found")
	}
	return engine, nil
}

// IsRunning is true if the engine exists and is ready.
func (r *EngineRegistry) IsRunning(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	engine, found := r.engines[id]
	if !found {
		return false
	}

	return engine.Ready() == nil
}

// Pop removes an engine from the registry and returns the engine if found.
func (r *EngineRegistry) Pop(id string) (services.Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	engine, ok := r.engines[id]
	if !ok {
		return nil, errors.New("remove failed: engine not found")
	}
	delete(r.engines, id)
	return engine, nil
}

// Close closes all engines in the registry.
func (r *EngineRegistry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	for id, engine := range r.engines {
		closeErr := engine.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
		delete(r.engines, id)
	}
	return err
}
