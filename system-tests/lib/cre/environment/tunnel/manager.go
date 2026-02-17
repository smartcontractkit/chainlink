package tunnel

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type manager struct {
	provider Provider

	mu       sync.Mutex
	bindings map[string]TunnelBinding
}

func NewManager(provider Provider) Manager {
	return &manager{
		provider: provider,
		bindings: make(map[string]TunnelBinding),
	}
}

func (m *manager) Start(ctx context.Context, refs []EndpointRef) ([]TunnelBinding, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	started := make([]TunnelBinding, 0, len(refs))
	newlyOpened := make([]TunnelBinding, 0, len(refs))

	for _, ref := range refs {
		key := endpointKey(ref.ComponentID, ref.EndpointName)
		if existing, ok := m.bindings[key]; ok {
			started = append(started, existing)
			continue
		}

		if err := validateEndpointRef(ref); err != nil {
			_ = m.closeMany(ctx, newlyOpened)
			return nil, err
		}

		binding, err := m.provider.Open(ctx, ref)
		if err != nil {
			_ = m.closeMany(ctx, newlyOpened)
			return nil, fmt.Errorf("failed to open tunnel via %s for %s/%s: %w", m.provider.Name(), ref.ComponentID, ref.EndpointName, err)
		}

		m.bindings[key] = binding
		started = append(started, binding)
		newlyOpened = append(newlyOpened, binding)
	}

	return started, nil
}

func (m *manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	bindings := make([]TunnelBinding, 0, len(m.bindings))
	for _, b := range m.bindings {
		bindings = append(bindings, b)
	}
	clear(m.bindings)

	return m.closeMany(ctx, bindings)
}

func (m *manager) IsStarted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.bindings) > 0
}

func (m *manager) Snapshot() []TunnelBinding {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]TunnelBinding, 0, len(m.bindings))
	for _, b := range m.bindings {
		out = append(out, b)
	}
	return out
}

func (m *manager) closeMany(ctx context.Context, bindings []TunnelBinding) error {
	var joined error
	for _, b := range bindings {
		if err := m.provider.Close(ctx, b); err != nil {
			joined = errors.Join(joined, err)
		}
	}
	return joined
}

func validateEndpointRef(ref EndpointRef) error {
	if ref.ComponentID == "" {
		return errors.New("endpoint componentID is required")
	}
	if ref.EndpointName == "" {
		return errors.New("endpoint endpointName is required")
	}
	if ref.Host == "" {
		return errors.New("endpoint host is required")
	}
	if ref.Port <= 0 || ref.Port > 65535 {
		return fmt.Errorf("endpoint port %d is invalid", ref.Port)
	}
	return nil
}

func endpointKey(componentID, endpointName string) string {
	return componentID + ":" + endpointName
}
