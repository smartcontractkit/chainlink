package loop

import (
	"context"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// Plugin is a base layer for plugins to easily manage sub-[types.Service]s.
// Useful for implementing PluginRelayer and PluginMedian.
type Plugin struct {
	Logger logger.Logger

	mu sync.RWMutex
	ss []services.Service
}

func (p *Plugin) Ready() error { return nil }
func (p *Plugin) Name() string { return p.Logger.Name() }

func (p *Plugin) SubService(s services.Service) {
	p.mu.Lock()
	p.ss = append(p.ss, s)
	p.mu.Unlock()
}

func (p *Plugin) Start(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var ms services.MultiStart
	for _, s := range p.ss {
		if err := ms.Start(ctx, s); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) HealthReport() map[string]error {
	hr := map[string]error{p.Name(): nil}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.ss {
		services.CopyHealth(hr, s.HealthReport())
	}
	return hr
}

func (p *Plugin) Close() (err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return services.MultiCloser(p.ss).Close()
}
