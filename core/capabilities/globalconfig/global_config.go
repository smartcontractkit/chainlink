// Package globalconfig holds the offchain capabilities registry config delivered to a
// node via the cresettings job (config_type=capabilities_registry).
//
// The cresettings delegate writes the latest payload here; the Launcher reads it (Phase 2)
// and cross-validates it against the on-chain registry. Payloads are parsed into the
// OffchainCapabilitiesRegistry proto shared with chainlink-common so the offchain and
// on-chain shapes stay identical (see the Offchain Capabilities Registry design).
package globalconfig

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
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
	parsed  *capabilitiespb.OffchainCapabilitiesRegistry
	version uint64
	hash    string
}

// New returns an empty GlobalConfig.
func New() *GlobalConfig { return &GlobalConfig{} }

// Store validates and applies an update. It is idempotent for a repeated hash and rejects
// any payload whose version is not strictly greater than the currently applied version
// (monotonic; version 0 means "unset" and is always accepted as the first value).
func (g *GlobalConfig) Store(u Update) error {
	reg, err := parse(u.Raw)
	if err != nil {
		return err
	}
	version := reg.GetVersion()

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.hash != "" && u.Hash == g.hash {
		return nil // idempotent re-apply of the same payload
	}
	if g.version != 0 && version <= g.version {
		return fmt.Errorf("offchain config version %d is not newer than applied version %d", version, g.version)
	}

	g.raw = u.Raw
	g.parsed = reg
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

// LoadParsed returns the current parsed registry and its version. The registry is nil when
// nothing has been applied yet. The returned message must not be mutated.
func (g *GlobalConfig) LoadParsed() (reg *capabilitiespb.OffchainCapabilitiesRegistry, version uint64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.parsed, g.version
}

// Validate checks that a payload parses into an OffchainCapabilitiesRegistry. It does not
// perform on-chain cross-validation, which happens in the Launcher where the on-chain DON
// set is available.
func Validate(raw string) error {
	_, err := parse(raw)
	return err
}

// Parse decodes an offchain_config payload (proto JSON) into an OffchainCapabilitiesRegistry.
func Parse(raw string) (*capabilitiespb.OffchainCapabilitiesRegistry, error) {
	return parse(raw)
}

// parse decodes proto-JSON into the registry proto. Unknown fields are discarded so that a
// node running an older schema tolerates payloads authored against a newer one.
func parse(raw string) (*capabilitiespb.OffchainCapabilitiesRegistry, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("offchain config: empty payload")
	}
	var reg capabilitiespb.OffchainCapabilitiesRegistry
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal([]byte(raw), &reg); err != nil {
		return nil, fmt.Errorf("offchain config: invalid payload: %w", err)
	}
	return &reg, nil
}
