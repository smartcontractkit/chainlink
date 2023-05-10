package plugins

import (
	"errors"
	"sort"
	"sync"
)

const (
	PluginDefaultPort = 2112
	invalidPort       = -1
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
}

func NewLoopRegistry() *LoopRegistry {
	return &LoopRegistry{
		registry: map[string]*RegisteredLoop{},
	}
}

// Register creates a port of the plugin. It is idempotent. Duplicate calls to Register will return the same port
func (m *LoopRegistry) Register(id string, staticCfg LoggingConfig) *RegisteredLoop {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.get(id)
	if !ok {
		// safe to ignore error because we are inside a lock
		// and non-existent has be checked.
		p, _ = m.create(id, staticCfg)
	}
	return p
}

func (m *LoopRegistry) List() []*RegisteredLoop {
	var registeredLoops []*RegisteredLoop
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, known := range m.registry {
		registeredLoops = append(registeredLoops, known)
	}
	sort.Slice(registeredLoops, func(i, j int) bool {
		return registeredLoops[i].Name < registeredLoops[j].Name
	})
	return registeredLoops
}

func (m *LoopRegistry) Get(id string) (*RegisteredLoop, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.get(id)
	return p, ok
}

// create returns a port number for the given plugin to use for prometheus handler.
func (m *LoopRegistry) create(pluginName string, staticCfg LoggingConfig) (*RegisteredLoop, error) {
	if _, exists := m.registry[pluginName]; exists {
		return nil, ErrExists
	}
	nextPort := PluginDefaultPort + len(m.registry)
	envCfg := NewEnvConfig(staticCfg, nextPort)

	m.registry[pluginName] = &RegisteredLoop{Name: pluginName, EnvCfg: envCfg}
	return m.registry[pluginName], nil
}

// get returns the port assigned to the plugin, if any
func (m *LoopRegistry) get(pluginName string) (*RegisteredLoop, bool) {
	p, exists := m.registry[pluginName]
	return p, exists
}
