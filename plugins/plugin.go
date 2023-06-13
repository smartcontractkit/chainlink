package plugins

import (
	"sync"

	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

// Base is a base layer for plugins to easily manage sub-[types.Service]s.
type Base struct {
	Logger logger.Logger

	mu   sync.RWMutex
	srvs []types.Service
}

func (p *Base) Ready() error { return nil }
func (p *Base) Name() string { return p.Logger.Name() }

func (p *Base) SubService(s types.Service) {
	p.mu.Lock()
	p.srvs = append(p.srvs, s)
	p.mu.Unlock()
}

func (p *Base) HealthReport() map[string]error {
	hr := map[string]error{p.Name(): nil}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, s := range p.srvs {
		maps.Copy(s.HealthReport(), hr)
	}
	return hr
}

func (p *Base) Close() (err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return services.MultiCloser(p.srvs).Close()
}
