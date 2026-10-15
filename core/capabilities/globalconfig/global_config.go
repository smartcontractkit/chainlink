// Package globalconfig holds the offchain capabilities registry config delivered to a
// node via the cresettings job (config_type=capabilities_registry).
//
// The cresettings delegate writes the latest payload here; the LocalCapabilityManager reads a
// snapshot of it on every reconcile. Payloads are parsed into the OffchainCapabilitiesRegistry
// proto shared with chainlink-common so the offchain and on-chain shapes stay identical (see
// the Offchain Capabilities Registry design).
package globalconfig

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

// GlobalConfig holds the latest applied offchain capabilities registry config.
// The zero value is ready to use.
type GlobalConfig struct {
	mu     sync.RWMutex
	parsed *capabilitiespb.OffchainCapabilitiesRegistry
	hash   string

	subsMu sync.Mutex
	subs   map[chan struct{}]struct{}
}

// New returns an empty GlobalConfig.
func New() *GlobalConfig { return &GlobalConfig{} }

// Store validates and applies a payload and returns the applied version.
//
// Re-applying the currently applied hash is a no-op. Otherwise the version must be strictly
// greater than the applied one; a rejected payload leaves the applied one in place.
// Subscribers are notified after every update that changes the applied payload.
//
// There is no way to withdraw a payload: deleting the capabilities_registry job keeps the last
// one applied until a newer one is stored or the node restarts.
func (g *GlobalConfig) Store(raw, hash string) (uint64, error) {
	reg, err := parse(raw)
	if err != nil {
		return 0, err
	}

	g.mu.Lock()
	applied := g.parsed.GetVersion()
	if g.parsed != nil && hash == g.hash {
		g.mu.Unlock()
		return applied, nil
	}
	if reg.GetVersion() <= applied {
		g.mu.Unlock()
		return 0, fmt.Errorf("offchain config version %d is not newer than applied version %d", reg.GetVersion(), applied)
	}
	g.parsed = reg
	g.hash = hash
	g.mu.Unlock()

	g.notify()
	return reg.GetVersion(), nil
}

// LoadParsed returns a deep copy of the applied registry, or nil when nothing is applied.
func (g *GlobalConfig) LoadParsed() *capabilitiespb.OffchainCapabilitiesRegistry {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.parsed == nil {
		return nil
	}
	return proto.Clone(g.parsed).(*capabilitiespb.OffchainCapabilitiesRegistry)
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
//   - capabilities may only set spec_config and method_configs,
//   - every spec_config must convert to a config map,
//   - method_configs must pass the offchain method validation before conversion.
//
// It does not perform on-chain cross-validation, which happens in the LocalCapabilityManager
// where the on-chain DON set is available.
func Validate(raw string) error {
	_, err := parse(raw)
	return err
}

// SpecConfigMap converts an offchain spec_config into the flat map a capability is launched
// with. A nil spec_config yields a nil map.
//
// The value tree is checked before conversion: a value with no type would otherwise become a
// JSON null in the launched config, and some empty
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
			// Only spec_config and method_configs are sourced offchain in this slice;
			// ocr3_configs and oracle_factory_configs remain on-chain.
			allowed := &capabilitiespb.CapabilityConfig{SpecConfig: capCfg.GetSpecConfig(), MethodConfigs: capCfg.GetMethodConfigs()}
			if !proto.Equal(capCfg, allowed) {
				return fmt.Errorf("dons[%q].capabilities[%q]: only spec_config and method_configs are supported", donName, capID)
			}
			if _, err := SpecConfigMap(capCfg.GetSpecConfig()); err != nil {
				return fmt.Errorf("dons[%q].capabilities[%q]: invalid spec_config: %w", donName, capID, err)
			}
			if _, err := MethodConfigsFromProto(capCfg.GetMethodConfigs()); err != nil {
				return fmt.Errorf("dons[%q].capabilities[%q]: invalid method_configs: %w", donName, capID, err)
			}
		}
	}
	return nil
}

// MethodConfigsFromProto converts an offchain method_configs map (pb) into the
// capabilities.CapabilityMethodConfig map the launcher's don2don shims consume. For every set
// remote_config oneof it mirrors the conversion in registry.CapabilityConfiguration.Unmarshal in
// chainlink-common field-for-field, so a trigger/executable method converts identically to the
// on-chain path. A nil map yields nil. Invalid durations, unknown enum values and executable
// configs that cannot be applied by the shims are rejected before conversion.
//
// An unset remote_config oneof has a deliberately different meaning offchain. On-chain,
// Unmarshal rejects that entry (its switch has no such case) and the whole
// capability config fails to load. Here it converts to a zero CapabilityMethodConfig (both
// RemoteTriggerConfig and RemoteExecutableConfig nil) and is accepted at ingestion. The launcher
// then treats a zero config as "no remote config found" for that method, so it is neither added to
// the wanted shim set nor created — the effect is to DISABLE that method's shim (any existing one is
// pruned and not recreated). This is an intentional offchain-only lever (disable a method without an
// on-chain change); the cross-check reports it as mismatch. If strict on-chain parity is ever
// required instead, make the nil case below return an error so such entries are rejected at
// ingestion like the on-chain path.
func MethodConfigsFromProto(in map[string]*capabilitiespb.CapabilityMethodConfig) (map[string]capabilities.CapabilityMethodConfig, error) {
	if in == nil {
		return nil, nil
	}
	out := make(map[string]capabilities.CapabilityMethodConfig, len(in))
	for method, mc := range in {
		if method == "" {
			return nil, errors.New("empty method name")
		}
		if err := validateMethodConfig(mc); err != nil {
			return nil, fmt.Errorf("method %q: %w", method, err)
		}
		var config capabilities.CapabilityMethodConfig
		switch remoteCfg := mc.GetRemoteConfig().(type) {
		case nil:
			// unset oneof: zero config, handled by the launcher as no remote config
		case *capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig:
			rt := remoteCfg.RemoteTriggerConfig
			config = capabilities.CapabilityMethodConfig{
				RemoteTriggerConfig: &capabilities.RemoteTriggerConfig{
					RegistrationRefresh:     rt.GetRegistrationRefresh().AsDuration(),
					RegistrationExpiry:      rt.GetRegistrationExpiry().AsDuration(),
					MinResponsesToAggregate: rt.GetMinResponsesToAggregate(),
					MessageExpiry:           rt.GetMessageExpiry().AsDuration(),
					MaxBatchSize:            rt.GetMaxBatchSize(),
					BatchCollectionPeriod:   rt.GetBatchCollectionPeriod().AsDuration(),
				},
			}
		case *capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig:
			re := remoteCfg.RemoteExecutableConfig
			config = capabilities.CapabilityMethodConfig{
				RemoteExecutableConfig: &capabilities.RemoteExecutableConfig{
					TransmissionSchedule:      capabilities.TransmissionSchedule(re.GetTransmissionSchedule()),
					DeltaStage:                re.GetDeltaStage().AsDuration(),
					RequestTimeout:            re.GetRequestTimeout().AsDuration(),
					ServerMaxParallelRequests: re.GetServerMaxParallelRequests(),
					RequestHasherType:         capabilities.RequestHasherType(re.GetRequestHasherType()),
					MinResponsesToAggregate:   re.GetMinResponsesToAggregate(),
				},
			}
		default:
			return nil, fmt.Errorf("method %q: unknown method config type %T", method, mc.GetRemoteConfig())
		}
		if agg := mc.GetAggregatorConfig(); agg != nil {
			config.AggregatorConfig = &capabilities.AggregatorConfig{
				AggregatorType: capabilities.AggregatorType(agg.GetAggregatorType()),
			}
		}
		out[method] = config
	}
	return out, nil
}

// validateMethodConfig checks values independently of the DON that will consume them. Trigger
// zero durations retain the shims' defaults; executable configs must supply the positive
// timeout and parallelism required by both sides of the don2don connection.
func validateMethodConfig(mc *capabilitiespb.CapabilityMethodConfig) error {
	if agg := mc.GetAggregatorConfig(); agg != nil {
		if _, ok := capabilitiespb.AggregatorType_name[int32(agg.GetAggregatorType())]; !ok {
			return fmt.Errorf("unknown aggregator_type %d", agg.GetAggregatorType())
		}
	}
	switch remote := mc.GetRemoteConfig().(type) {
	case nil:
		return nil // explicitly disabled
	case *capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig:
		rt := remote.RemoteTriggerConfig
		for _, field := range []struct {
			name  string
			value *durationpb.Duration
		}{
			{"registration_refresh", rt.GetRegistrationRefresh()},
			{"registration_expiry", rt.GetRegistrationExpiry()},
			{"message_expiry", rt.GetMessageExpiry()},
			{"batch_collection_period", rt.GetBatchCollectionPeriod()},
		} {
			if err := validateMethodDuration(field.name, field.value); err != nil {
				return err
			}
		}
	case *capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig:
		re := remote.RemoteExecutableConfig
		if _, ok := capabilitiespb.TransmissionSchedule_name[int32(re.GetTransmissionSchedule())]; !ok {
			return fmt.Errorf("unknown transmission_schedule %d", re.GetTransmissionSchedule())
		}
		if _, ok := capabilitiespb.RequestHasherType_name[int32(re.GetRequestHasherType())]; !ok {
			return fmt.Errorf("unknown request_hasher_type %d", re.GetRequestHasherType())
		}
		if err := validateMethodDuration("delta_stage", re.GetDeltaStage()); err != nil {
			return err
		}
		if err := validateMethodDuration("request_timeout", re.GetRequestTimeout()); err != nil {
			return err
		}
		if re.GetRequestTimeout().AsDuration() <= 0 {
			return errors.New("request_timeout must be positive")
		}
		if re.GetServerMaxParallelRequests() == 0 {
			return errors.New("server_max_parallel_requests must be positive")
		}
	default:
		return fmt.Errorf("unknown method config type %T", mc.GetRemoteConfig())
	}
	return nil
}

func validateMethodDuration(name string, value *durationpb.Duration) error {
	if value == nil {
		return nil
	}
	if err := value.CheckValid(); err != nil {
		return fmt.Errorf("invalid %s: %w", name, err)
	}
	d := value.AsDuration()
	if d < 0 {
		return fmt.Errorf("%s must not be negative", name)
	}
	// AsDuration saturates on overflow. Reject values that would silently change in conversion.
	roundTrip := durationpb.New(d)
	if roundTrip.Seconds != value.Seconds || roundTrip.Nanos != value.Nanos {
		return fmt.Errorf("%s exceeds the Go duration range", name)
	}
	return nil
}
