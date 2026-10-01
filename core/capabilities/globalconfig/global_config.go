// Package globalconfig holds the offchain capabilities registry config delivered to a
// node via the cresettings job (config_type=capabilities_registry).
//
// The cresettings delegate writes the latest payload here; the LocalCapabilityManager reads a
// snapshot of it on every reconcile. Payloads are parsed into the OffchainCapabilitiesRegistry
// proto shared with chainlink-common so the offchain and on-chain shapes stay identical (see
// the Offchain Capabilities Registry design).
package globalconfig

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

// Update is a single offchain config delivery.
type Update struct {
	// Raw is the offchain_config payload (proto JSON) from the cresettings spec.
	Raw string
	// Hash is the spec hash; used to short-circuit idempotent re-applies. When empty, the
	// sha256 of Raw is used (the same hash the cresettings validator computes).
	Hash string
}

// GlobalConfig holds the latest applied offchain capabilities registry config.
// The zero value is ready to use.
type GlobalConfig struct {
	mu     sync.RWMutex
	raw    string
	parsed *capabilitiespb.OffchainCapabilitiesRegistry
	// version and hash describe the last applied payload. They are kept by Clear as a
	// high-water mark, so a payload older than one already applied stays rejected after the
	// job that carried it is deleted.
	version uint64
	hash    string

	subsMu sync.Mutex
	subs   map[chan struct{}]struct{}
}

// New returns an empty GlobalConfig.
func New() *GlobalConfig { return &GlobalConfig{} }

// Store validates and applies an update.
//
//   - The payload must parse and pass Validate (version >= 1, well-formed spec_config, ...).
//   - Re-applying the currently applied hash is a no-op (idempotent). After Clear, re-applying
//     the last applied hash restores it.
//   - Otherwise the version must be strictly greater than the applied version; a stale or
//     equal version with different content is rejected and the applied payload is kept.
//
// Subscribers are notified after every update that changes the applied payload.
func (g *GlobalConfig) Store(u Update) error {
	reg, err := parse(u.Raw)
	if err != nil {
		return err
	}
	hash := u.Hash
	if hash == "" {
		sum := sha256.Sum256([]byte(u.Raw))
		hash = hex.EncodeToString(sum[:])
	}
	version := reg.GetVersion()

	g.mu.Lock()
	if g.hash != "" && hash == g.hash {
		if g.parsed != nil {
			g.mu.Unlock()
			return nil // idempotent re-apply of the same payload
		}
		// Same payload as the last applied one, re-applied after Clear: restore it.
	} else if version <= g.version {
		applied := g.version
		g.mu.Unlock()
		return fmt.Errorf("offchain config version %d is not newer than applied version %d", version, applied)
	}
	g.raw = u.Raw
	g.parsed = reg
	g.version = version
	g.hash = hash
	g.mu.Unlock()

	g.notify()
	return nil
}

// Clear withdraws the applied payload, e.g. after the capabilities_registry job is deleted, so
// readers fall back to on-chain/TOML config. The version/hash high-water mark is kept, so a
// later payload must still be newer than the last one applied (or be that same payload).
// Subscribers are notified when a payload was actually withdrawn.
func (g *GlobalConfig) Clear() {
	g.mu.Lock()
	had := g.parsed != nil
	g.raw = ""
	g.parsed = nil
	g.mu.Unlock()
	if had {
		g.notify()
	}
}

// Load returns the current raw payload and its version. version is 0 when nothing is applied
// (never, or since Clear).
func (g *GlobalConfig) Load() (raw string, version uint64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.parsed == nil {
		return "", 0
	}
	return g.raw, g.version
}

// Info returns the domain, env and version of the applied payload; all zero values when
// nothing is applied.
func (g *GlobalConfig) Info() (domain, env string, version uint64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.parsed == nil {
		return "", "", 0
	}
	return g.parsed.GetDomain(), g.parsed.GetEnv(), g.version
}

// LoadParsed returns a snapshot of the current parsed registry and its version. The registry
// is nil and the version 0 when nothing is applied (never, or since Clear). The snapshot is a deep copy owned by the caller,
// so mutating it cannot affect the applied config or other readers.
func (g *GlobalConfig) LoadParsed() (reg *capabilitiespb.OffchainCapabilitiesRegistry, version uint64) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.parsed == nil {
		return nil, 0
	}
	return proto.Clone(g.parsed).(*capabilitiespb.OffchainCapabilitiesRegistry), g.version
}

// Subscribe returns a channel that is signalled after each update that changes the applied
// payload, and a function to unsubscribe. Signals are coalesced: the channel has a buffer of
// one, so a slow reader observes at least one signal after the latest change. Readers should
// call LoadParsed to obtain the current state rather than assume one signal per update.
func (g *GlobalConfig) Subscribe() (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1)
	g.subsMu.Lock()
	if g.subs == nil {
		g.subs = make(map[chan struct{}]struct{})
	}
	g.subs[ch] = struct{}{}
	g.subsMu.Unlock()
	return ch, func() {
		g.subsMu.Lock()
		delete(g.subs, ch)
		g.subsMu.Unlock()
	}
}

func (g *GlobalConfig) notify() {
	g.subsMu.Lock()
	defer g.subsMu.Unlock()
	for ch := range g.subs {
		select {
		case ch <- struct{}{}:
		default: // a signal is already pending
		}
	}
}

// Validate checks that a payload parses into a well-formed OffchainCapabilitiesRegistry:
//   - the payload must contain only fields known to this node (see parse),
//   - version must be >= 1 (0 is the proto default, i.e. "unset", and would defeat the
//     monotonic version check),
//   - DON names (map keys) and capability IDs must be non-empty,
//   - every spec_config must convert to a config map (the shape the launcher merges).
//
// It does not perform on-chain cross-validation, which happens in the LocalCapabilityManager
// where the on-chain DON set is available.
func Validate(raw string) error {
	_, err := parse(raw)
	return err
}

// Parse decodes and validates an offchain_config payload (proto JSON).
func Parse(raw string) (*capabilitiespb.OffchainCapabilitiesRegistry, error) {
	return parse(raw)
}

// SpecConfigMap converts an offchain spec_config into the flat config map merged into the
// config a capability is launched with. A nil spec_config yields a nil map.
//
// The value tree is checked before conversion: a value with no type would otherwise become a
// JSON null that silently overrides the on-chain/TOML value for that key, and some empty
// messages (e.g. a decimal without a coefficient) panic during conversion.
func SpecConfigMap(sc *valuespb.Map) (map[string]any, error) {
	if sc == nil {
		return nil, nil
	}
	for k, v := range sc.GetFields() {
		if err := checkValue(v); err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
	}
	m, err := values.FromMapValueProto(sc)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	unwrapped, err := m.Unwrap()
	if err != nil {
		return nil, err
	}
	out, ok := unwrapped.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("spec_config unwrapped to %T, expected map", unwrapped)
	}
	return out, nil
}

// checkValue rejects value trees that values.FromProto would convert to a JSON null or panic on.
func checkValue(v *valuespb.Value) error {
	switch k := v.GetValue().(type) {
	case nil:
		return errors.New("value has no type")
	case *valuespb.Value_DecimalValue:
		if k.DecimalValue == nil || k.DecimalValue.GetCoefficient() == nil {
			return errors.New("decimal value has no coefficient")
		}
	case *valuespb.Value_BigintValue:
		if k.BigintValue == nil {
			return errors.New("bigint value is empty")
		}
	case *valuespb.Value_TimeValue:
		if k.TimeValue == nil {
			return errors.New("time value is empty")
		}
	case *valuespb.Value_ListValue:
		if k.ListValue == nil {
			return errors.New("list value is empty")
		}
		for i, el := range k.ListValue.GetFields() {
			if err := checkValue(el); err != nil {
				return fmt.Errorf("[%d]: %w", i, err)
			}
		}
	case *valuespb.Value_MapValue:
		if k.MapValue == nil {
			return errors.New("map value is empty")
		}
		for key, el := range k.MapValue.GetFields() {
			if err := checkValue(el); err != nil {
				return fmt.Errorf("%q: %w", key, err)
			}
		}
	}
	return nil
}

// parse decodes proto-JSON into the registry proto and validates it.
//
// Parsing is strict: a field this node does not know is an error, not silently dropped. A
// misspelled or wrong-schema field (e.g. "capabilityConfigs" instead of "capabilities") would
// otherwise be accepted as a payload with no config, and the node would quietly fall back to
// legacy config. Payloads using new schema fields must therefore only be rolled out to nodes
// that understand them; older nodes reject them visibly (job error).
func parse(raw string) (*capabilitiespb.OffchainCapabilitiesRegistry, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("offchain config: empty payload")
	}
	var reg capabilitiespb.OffchainCapabilitiesRegistry
	if err := protojson.Unmarshal([]byte(raw), &reg); err != nil {
		return nil, fmt.Errorf("offchain config: invalid payload: %w", err)
	}
	if err := validate(&reg); err != nil {
		return nil, fmt.Errorf("offchain config: %w", err)
	}
	return &reg, nil
}

func validate(reg *capabilitiespb.OffchainCapabilitiesRegistry) error {
	if reg.GetVersion() == 0 {
		return errors.New("version must be >= 1")
	}
	for donName, don := range reg.GetDons() {
		if donName == "" {
			return errors.New("dons: empty DON name")
		}
		for capID, capCfg := range don.GetCapabilities() {
			if capID == "" {
				return fmt.Errorf("dons[%q]: empty capability ID", donName)
			}
			if _, err := SpecConfigMap(capCfg.GetSpecConfig()); err != nil {
				return fmt.Errorf("dons[%q].capabilities[%q]: invalid spec_config: %w", donName, capID, err)
			}
		}
	}
	return nil
}
