package plugins

import (
	"errors"
	"sort"
	"sync"
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
}

func NewLoopRegistry() *LoopRegistry {
	return &LoopRegistry{
		registry: map[string]*RegisteredLoop{},
	}
}

// Register creates a port of the plugin. It is idempotent. Duplicate calls to Register will return the same port
func (m *LoopRegistry) Register(id string) (*RegisteredLoop, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.get(id)
	if !ok {
		return m.create(id)
	}
	return p, nil
}

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

func (m *LoopRegistry) Get(id string) (*RegisteredLoop, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.get(id)
}

// create returns a port number for the given plugin to use for prometheus handler.
// NOT safe for concurrent use.
func (m *LoopRegistry) create(pluginName string) (*RegisteredLoop, error) {
	if _, exists := m.registry[pluginName]; exists {
		return nil, ErrExists
	}
	nextPort := pluginDefaultPort + len(m.registry)
	envCfg := NewEnvConfig(nextPort)

	m.registry[pluginName] = &RegisteredLoop{Name: pluginName, EnvCfg: envCfg}
	return m.registry[pluginName], nil
}

// get is a helper to return the port assigned to the plugin, if any
// NOT safe for concurrent use.
func (m *LoopRegistry) get(pluginName string) (*RegisteredLoop, bool) {
	p, exists := m.registry[pluginName]
	return p, exists
}
