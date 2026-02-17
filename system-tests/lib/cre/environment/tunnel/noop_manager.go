package tunnel

import "context"

type noopManager struct{}

func NewNoopManager() Manager {
	return &noopManager{}
}

func (n *noopManager) Start(_ context.Context, refs []EndpointRef) ([]TunnelBinding, error) {
	bindings := make([]TunnelBinding, 0, len(refs))
	for _, ref := range refs {
		bindings = append(bindings, TunnelBinding{
			EndpointRef: ref,
			LocalURL:    ref.OriginalURL,
		})
	}
	return bindings, nil
}

func (n *noopManager) Stop(_ context.Context) error { return nil }

func (n *noopManager) IsStarted() bool { return false }
