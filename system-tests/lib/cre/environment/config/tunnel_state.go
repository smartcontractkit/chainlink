package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

const TunnelStateFilename = "tunnels.toml"

type TunnelProcess struct {
	PID         int    `toml:"pid"`
	Kind        string `toml:"kind"`
	InstanceID  string `toml:"instance_id"`
	Region      string `toml:"region"`
	RemotePort  int    `toml:"remote_port"`
	LocalPort   int    `toml:"local_port"`
	ComponentID string `toml:"component_id,omitempty"`
	Endpoint    string `toml:"endpoint,omitempty"`
	CreatedAt   string `toml:"created_at,omitempty"`
}

type TunnelState struct {
	Version int             `toml:"version"`
	Tunnels []TunnelProcess `toml:"tunnels"`
}

var tunnelStateMu sync.Mutex

func MustTunnelStateFileAbsPath(relativePathToRepoRoot string) string {
	absPath, err := filepath.Abs(filepath.Join(relativePathToRepoRoot, StateDirname, TunnelStateFilename))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path for tunnel state file: %w", err))
	}
	return absPath
}

func LoadTunnelState(relativePathToRepoRoot string) (*TunnelState, error) {
	tunnelStateMu.Lock()
	defer tunnelStateMu.Unlock()
	return loadTunnelStateUnlocked(MustTunnelStateFileAbsPath(relativePathToRepoRoot))
}

func StoreTunnelState(relativePathToRepoRoot string, state *TunnelState) error {
	tunnelStateMu.Lock()
	defer tunnelStateMu.Unlock()
	return storeTunnelStateUnlocked(MustTunnelStateFileAbsPath(relativePathToRepoRoot), state)
}

func ClearTunnelState(relativePathToRepoRoot string) error {
	tunnelStateMu.Lock()
	defer tunnelStateMu.Unlock()
	return storeTunnelStateUnlocked(MustTunnelStateFileAbsPath(relativePathToRepoRoot), &TunnelState{
		Version: 1,
		Tunnels: []TunnelProcess{},
	})
}

func loadTunnelStateUnlocked(path string) (*TunnelState, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &TunnelState{Version: 1, Tunnels: []TunnelProcess{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read tunnel state file: %w", err)
	}

	state := &TunnelState{}
	if err := toml.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tunnel state file: %w", err)
	}
	if state.Version == 0 {
		state.Version = 1
	}
	if state.Tunnels == nil {
		state.Tunnels = []TunnelProcess{}
	}
	return state, nil
}

func storeTunnelStateUnlocked(path string, state *TunnelState) error {
	if state == nil {
		state = &TunnelState{Version: 1, Tunnels: []TunnelProcess{}}
	}
	if state.Version == 0 {
		state.Version = 1
	}
	if state.Tunnels == nil {
		state.Tunnels = []TunnelProcess{}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create tunnel state directory: %w", err)
	}
	data, err := toml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal tunnel state: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write tunnel state file: %w", err)
	}
	return nil
}
