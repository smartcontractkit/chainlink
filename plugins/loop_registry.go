package plugins

import (
	"errors"
	"sort"
	"sync"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

const (
	pluginDefaultPort = 2112
)

var ErrExists = errors.New("plugin already registered")

type RegisteredLoop struct {
	Name   string
	EnvCfg EnvConfig
}

// LoopRegistry is responsible for assigning ports to plugins that are to be used for the
// plugin's prometheus HTTP server
type LoopRegistry struct {
	mu       sync.Mutex
	registry map[string]*RegisteredLoop

	lggr logger.Logger
}

func NewLoopRegistry(lggr logger.Logger) *LoopRegistry {
	return &LoopRegistry{
		registry: map[string]*RegisteredLoop{},
		lggr:     lggr,
	}
}

// Register creates a port of the plugin. It is not idempotent. Duplicate calls to Register will return [ErrExists]
// Safe for concurrent use.
func (m *LoopRegistry) Register(id string) (*RegisteredLoop, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.registry[id]; exists {
		return nil, ErrExists
	}
	nextPort := pluginDefaultPort + len(m.registry)
	envCfg := NewEnvConfig(nextPort)

	m.registry[id] = &RegisteredLoop{Name: id, EnvCfg: envCfg}
	m.lggr.Debug("Registered loopp %q with config %v, port %d", id, envCfg, envCfg.PrometheusPort())
	return m.registry[id], nil
}

// Return slice sorted by plugin name. Safe for concurrent use.
func (m *LoopRegistry) List() []*RegisteredLoop {
	var registeredLoops []*RegisteredLoop
	m.mu.Lock()
	for _, known := range m.registry {
		registeredLoops = append(registeredLoops, known)
	}
	m.mu.Unlock()

	sort.Slice(registeredLoops, func(i, j int) bool {
		return registeredLoops[i].Name < registeredLoops[j].Name
	})
	return registeredLoops
}

// Get plugin by id. Safe for concurrent use.
func (m *LoopRegistry) Get(id string) (*RegisteredLoop, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.registry[id]
	return p, exists
}
