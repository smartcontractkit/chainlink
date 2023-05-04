package chainlink

import (
	"errors"
	"sync"
)

const (
	PluginDefaultPort = 2112
	invalidPort       = -1
)

var ErrExists = errors.New("plugin already registered")

// PluginPortManager is responsible for assigning ports to plugins that are to be used for the
// plugin's prometheus HTTP server
type PluginPortManager struct {
	mu      sync.Mutex
	portMap map[string]int
}

func NewPluginPortManager() *PluginPortManager {
	return &PluginPortManager{
		portMap: map[string]int{},
	}
}

// Register creates a port of the plugin. It is idempotent. Duplicate calls to Register will return the same port
func (m *PluginPortManager) Register(plugName string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.get(plugName)
	if !ok {
		// safe to ignore error because we are inside a lock
		// and non-existent has be checked.
		p, _ = m.create(plugName)
	}
	return p
}

// create returns a port number for the given plugin to use for prometheus handler.
func (m *PluginPortManager) create(pluginName string) (int, error) {
	if _, exists := m.portMap[pluginName]; exists {
		return invalidPort, ErrExists
	}
	p := PluginDefaultPort + len(m.portMap)
	m.portMap[pluginName] = p
	return p, nil
}

// get returns the port assigned to the plugin, if any
func (m *PluginPortManager) get(pluginName string) (int, bool) {
	p, exists := m.portMap[pluginName]
	return p, exists
}
