// Package globalconfig holds the offchain capabilities registry config delivered to a
// node via the cresettings job (config_type=capabilities_registry).
//
// The cresettings delegate writes the latest payload here; the Launcher reads it (Phase 2)
// and merges it with the on-chain registry. For now GlobalConfig stores the raw payload and
// its monotonic version only — full parsing into per-DON capability config lands with the
// OffchainCapabilitiesRegistry proto (see the Offchain Capabilities Registry design).
package globalconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Update is a single offchain config delivery.
type Update struct {
	// Raw is the offchain_config payload (proto JSON) from the cresettings spec.
	Raw string
	// Hash is the spec hash; used to short-circuit idempotent re-applies.
	Hash string
}

// GlobalConfig holds the latest applied offchain capabilities registry config.
// The zero value is ready to use.
type GlobalConfig struct {
	mu      sync.RWMutex
	raw     string
	version uint64
	hash    string
}

// New returns an empty GlobalConfig.
func New() *GlobalConfig { return &GlobalConfig{} }

// Store validates and applies an update. It is idempotent for a repeated hash and rejects
// any payload whose version is not strictly greater than the currently applied version
// (monotonic; version 0 means "unset" and is always accepted as the first value).
func (g *GlobalConfig) Store(u Update) error {
	version, err := peekVersion(u.Raw)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.hash != "" && u.Hash == g.hash {
		return nil // idempotent re-apply of the same payload
	}
	if g.version != 0 && version <= g.version {
		return fmt.Errorf("offchain config version %d is not newer than applied version %d", version, g.version)
	}

	g.raw = u.Raw
	g.version = version
	g.hash = u.Hash
	return nil
}

// Load returns the current raw payload and its version. version is 0 when nothing has been
// applied yet.
func (g *GlobalConfig) Load() (raw string, version uint64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.raw, g.version
}

// Validate checks that a payload is well-formed enough to accept (valid JSON with a
// readable version). It does not yet validate the full config structure.
func Validate(raw string) error {
	_, err := peekVersion(raw)
	return err
}

// peekVersion extracts the top-level "version" field from the offchain config JSON without
// depending on the full OffchainCapabilitiesRegistry proto type. Proto JSON encodes uint64
// as a string, but the human-authored form may use a number, so both are accepted.
func peekVersion(raw string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, errors.New("offchain config: empty payload")
	}
	var head struct {
		Version any `json:"version"`
	}
	if err := json.Unmarshal([]byte(raw), &head); err != nil {
		return 0, fmt.Errorf("offchain config: invalid JSON: %w", err)
	}
	switch v := head.Version.(type) {
	case nil:
		return 0, nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("offchain config: negative version %v", v)
		}
		return uint64(v), nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("offchain config: invalid version %q: %w", v, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("offchain config: invalid version type %T", v)
	}
}
