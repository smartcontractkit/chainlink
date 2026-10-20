package globalconfig

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

func TestGlobalConfig_Store(t *testing.T) {
	t.Parallel()

	g := New()
	v, err := g.Store(`{"version":5,"dons":{"don-1":{}}}`, "h5")
	require.NoError(t, err)
	assert.Equal(t, uint64(5), v)

	// Same hash is a no-op.
	v, err = g.Store(`{"version":5,"dons":{"don-1":{}}}`, "h5")
	require.NoError(t, err)
	assert.Equal(t, uint64(5), v)

	// Equal or older versions with different content, and invalid payloads, are rejected and
	// keep the applied payload.
	for _, raw := range []string{`{"version":5}`, `{"version":4}`, `{"version":0}`, `{"version":6,"dons":{"":{}}}`, ``} {
		_, err = g.Store(raw, "other")
		require.Error(t, err, raw)
	}
	assert.Equal(t, uint64(5), g.LoadParsed().GetVersion())
	assert.Contains(t, g.LoadParsed().GetDons(), "don-1")

	v, err = g.Store(`{"version":6}`, "h6")
	require.NoError(t, err)
	assert.Equal(t, uint64(6), v)
}

func TestGlobalConfig_LoadParsedIsASnapshot(t *testing.T) {
	t.Parallel()

	assert.Nil(t, New().LoadParsed())

	g := New()
	_, err := g.Store(`{"version":1,"dons":{"don-7":{"capabilities":{"cron@1.0.0":{}}}}}`, "h1")
	require.NoError(t, err)

	snap := g.LoadParsed()
	snap.Version = 99
	snap.Dons["don-8"] = &capabilitiespb.DONConfig{}

	again := g.LoadParsed()
	assert.Equal(t, uint64(1), again.GetVersion())
	assert.NotContains(t, again.GetDons(), "don-8")
}

func TestGlobalConfig_Subscribe(t *testing.T) {
	t.Parallel()

	g := New()
	ch, unsubscribe := g.Subscribe()
	store := func(raw, hash string) error { _, err := g.Store(raw, hash); return err }

	require.NoError(t, store(`{"version":1}`, "h1"))
	require.NoError(t, store(`{"version":2}`, "h2"))
	// Coalesced: two changes, one pending signal.
	assert.Len(t, ch, 1)
	<-ch

	// No-op and rejected stores do not signal.
	require.NoError(t, store(`{"version":2}`, "h2"))
	require.Error(t, store(`{"version":1}`, "h1"))
	assert.Empty(t, ch)

	unsubscribe()
	require.NoError(t, store(`{"version":3}`, "h3"))
	assert.Empty(t, ch)
}

// TestValidate_StrictFields: unknown or misspelled fields are rejected instead of silently
// dropped, so a payload can never be accepted with its capability config quietly missing.
func TestValidate_StrictFields(t *testing.T) {
	t.Parallel()
	for name, raw := range map[string]string{
		"unknown top-level field":            `{"version":1,"unknownField":true}`,
		"old schema field capabilityConfigs": `{"version":1,"dons":{"don-7":{"capabilityConfigs":{"cron@1.0.0":{}}}}}`,
		"old schema field donId":             `{"version":1,"dons":{"don-7":{"donId":7}}}`,
		"unknown CapabilityConfig field":     `{"version":1,"dons":{"don-7":{"capabilities":{"cron@1.0.0":{"ocrConfig":{"f":3}}}}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.ErrorContains(t, Validate(raw), "unknown field")
		})
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	require.NoError(t, Validate(`{"version":1,"dons":{}}`))
	require.NoError(t, Validate(`{"version":"1"}`))
	require.NoError(t, Validate(`{"domain":"cre","env":"prod","version":201,"dons":{"workflow_1_zone-a":{"capabilities":{"cron@1.0.0":{}}}}}`))
	require.Error(t, Validate(``))
	require.Error(t, Validate(`{`))
	require.ErrorContains(t, Validate(`{"version":0}`), "version must be >= 1")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"":{}}}`), "empty DON name")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"don-7":{"capabilities":{"":{}}}}}`), "empty capability ID")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"don-7":{"capabilities":{"c@1.0.0":{"ocr3Configs":{"default":{"f":1}}}}}}}`), "only spec_config and method_configs are supported")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"don-7":{"capabilities":{"c@1.0.0":{"localOnly":true}}}}}`), "only spec_config and method_configs are supported")

	t.Run("malformed spec_config is rejected at ingestion", func(t *testing.T) {
		t.Parallel()
		reg := &capabilitiespb.OffchainCapabilitiesRegistry{Version: 1, Dons: map[string]*capabilitiespb.DONConfig{
			"don-7": {Capabilities: map[string]*capabilitiespb.CapabilityConfig{
				"cron@1.0.0": {SpecConfig: &valuespb.Map{Fields: map[string]*valuespb.Value{"broken": {}}}},
			}},
		}}
		b, err := protojson.Marshal(reg)
		require.NoError(t, err)
		require.ErrorContains(t, Validate(string(b)), `capabilities["cron@1.0.0"]: invalid spec_config: key "broken": value has no type`)

		withSpec := func(spec string) string {
			return `{"version":1,"dons":{"don-7":{"capabilities":{"c@1.0.0":{"specConfig":` + spec + `}}}}}`
		}
		// Would panic in values.FromProto (decimal.NewFromBigInt(nil, ...)) if not rejected first.
		require.ErrorContains(t, Validate(withSpec(`{"fields":{"d":{"decimalValue":{}}}}`)), "decimal value has no coefficient")
		require.ErrorContains(t, Validate(withSpec(`{"fields":{"l":{"listValue":{"fields":[{"stringValue":"a"},{}]}}}}`)), "[1]: value has no type")
		require.ErrorContains(t, Validate(withSpec(`{"fields":{"m":{"mapValue":{"fields":{"x":{}}}}}}`)), `"x": value has no type`)
		require.ErrorContains(t, Validate(withSpec(`{"fields":{"s":"not-a-value"}}`)), "invalid payload")
		require.NoError(t, Validate(withSpec(`{"fields":{"s":{"stringValue":"ok"},"n":{"int64Value":"3"}}}`)))
	})

	t.Run("well-formed spec_config converts to a flat map", func(t *testing.T) {
		t.Parallel()
		vm, err := values.NewMap(map[string]any{"interval": "30", "n": int64(2)})
		require.NoError(t, err)
		got, err := SpecConfigMap(values.ProtoMap(vm))
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"interval": "30", "n": int64(2)}, got)

		got, err = SpecConfigMap(nil)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestMethodConfigsFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil map yields nil", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(nil)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("converts trigger config", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
			"LogTrigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(10 * time.Second),
						RegistrationExpiry:      durationpb.New(20 * time.Second),
						MinResponsesToAggregate: 3,
						MessageExpiry:           durationpb.New(30 * time.Second),
						MaxBatchSize:            5,
						BatchCollectionPeriod:   durationpb.New(40 * time.Second),
					},
				},
				AggregatorConfig: &capabilitiespb.AggregatorConfig{AggregatorType: capabilitiespb.AggregatorType_SignedReport},
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		cfg := got["LogTrigger"]
		require.NotNil(t, cfg.RemoteTriggerConfig)
		assert.Equal(t, 10*time.Second, cfg.RemoteTriggerConfig.RegistrationRefresh)
		assert.Equal(t, 20*time.Second, cfg.RemoteTriggerConfig.RegistrationExpiry)
		assert.Equal(t, uint32(3), cfg.RemoteTriggerConfig.MinResponsesToAggregate)
		assert.Equal(t, 30*time.Second, cfg.RemoteTriggerConfig.MessageExpiry)
		assert.Equal(t, uint32(5), cfg.RemoteTriggerConfig.MaxBatchSize)
		assert.Equal(t, 40*time.Second, cfg.RemoteTriggerConfig.BatchCollectionPeriod)
		require.NotNil(t, cfg.AggregatorConfig)
		assert.Equal(t, capabilities.AggregatorType_SignedReport, cfg.AggregatorConfig.AggregatorType)
		assert.Nil(t, cfg.RemoteExecutableConfig)
	})

	t.Run("converts executable config", func(t *testing.T) {
		t.Parallel()
		got, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
						DeltaStage:                durationpb.New(time.Second),
						RequestTimeout:            durationpb.New(15 * time.Second),
						ServerMaxParallelRequests: 7,
						RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
						MinResponsesToAggregate:   4,
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		cfg := got["Write"]
		require.NotNil(t, cfg.RemoteExecutableConfig)
		assert.Equal(t, capabilities.Schedule_OneAtATime, cfg.RemoteExecutableConfig.TransmissionSchedule)
		assert.Equal(t, time.Second, cfg.RemoteExecutableConfig.DeltaStage)
		assert.Equal(t, 15*time.Second, cfg.RemoteExecutableConfig.RequestTimeout)
		assert.Equal(t, uint32(7), cfg.RemoteExecutableConfig.ServerMaxParallelRequests)
		assert.Equal(t, capabilities.RequestHasherType_WriteReportExcludeSignatures, cfg.RemoteExecutableConfig.RequestHasherType)
		assert.Equal(t, uint32(4), cfg.RemoteExecutableConfig.MinResponsesToAggregate)
		assert.Nil(t, cfg.RemoteTriggerConfig)
		assert.Nil(t, cfg.AggregatorConfig)
	})

	t.Run("unset oneof is rejected like on-chain", func(t *testing.T) {
		t.Parallel()
		in := map[string]*capabilitiespb.CapabilityMethodConfig{"View": {}}
		_, err := MethodConfigsFromProto(in)
		require.ErrorContains(t, err, "remote_config must be set")
		encoded, err := proto.Marshal(&capabilitiespb.CapabilityConfig{MethodConfigs: in})
		require.NoError(t, err)
		_, err = (registry.CapabilityConfiguration{Config: encoded}).Unmarshal()
		require.Error(t, err)
	})

	t.Run("matches on-chain conversion", func(t *testing.T) {
		t.Parallel()
		// The same pb shape the on-chain path produces via
		// registry.CapabilityConfiguration.Unmarshal must convert identically here.
		in := map[string]*capabilitiespb.CapabilityMethodConfig{
			"Trigger": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
					RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
						RegistrationRefresh:     durationpb.New(time.Second),
						RegistrationExpiry:      durationpb.New(20 * time.Second),
						MinResponsesToAggregate: 3,
						MessageExpiry:           durationpb.New(30 * time.Second),
						MaxBatchSize:            5,
						BatchCollectionPeriod:   durationpb.New(100 * time.Millisecond),
					},
				},
				AggregatorConfig: &capabilitiespb.AggregatorConfig{AggregatorType: capabilitiespb.AggregatorType_SignedReport},
			},
			"Write": {
				RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
					RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
						TransmissionSchedule:      capabilitiespb.TransmissionSchedule_OneAtATime,
						DeltaStage:                durationpb.New(time.Second),
						RequestTimeout:            durationpb.New(30 * time.Second),
						ServerMaxParallelRequests: 10,
						RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
						MinResponsesToAggregate:   4,
					},
				},
			},
		}
		got, err := MethodConfigsFromProto(in)
		require.NoError(t, err)
		encoded, err := proto.Marshal(&capabilitiespb.CapabilityConfig{MethodConfigs: in})
		require.NoError(t, err)
		onchain, err := (registry.CapabilityConfiguration{Config: encoded}).Unmarshal()
		require.NoError(t, err)
		assert.Equal(t, onchain.CapabilityMethodConfig, got)
		require.Len(t, got, 2)
		require.NotNil(t, got["Write"].RemoteExecutableConfig)
		assert.Equal(t, 30*time.Second, got["Write"].RemoteExecutableConfig.RequestTimeout)
	})
}

func TestValidate_MethodConfigsRejectsInvalidUpdates(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		method  string
		config  string
		wantErr string
	}{
		{"unset remote config", "Write", `{}`, "remote_config must be set"},
		{"empty method", "", `{}`, "empty method name"},
		{"missing timeout", "Write", `{"remoteExecutableConfig":{"serverMaxParallelRequests":1}}`, "request_timeout must be positive"},
		{"zero timeout", "Write", `{"remoteExecutableConfig":{"requestTimeout":"0s","serverMaxParallelRequests":1}}`, "request_timeout must be positive"},
		{"negative timeout", "Write", `{"remoteExecutableConfig":{"requestTimeout":"-1s","serverMaxParallelRequests":1}}`, "request_timeout must not be negative"},
		{"missing parallelism", "Write", `{"remoteExecutableConfig":{"requestTimeout":"1s"}}`, "server_max_parallel_requests must be positive"},
		{"zero parallelism", "Write", `{"remoteExecutableConfig":{"requestTimeout":"1s","serverMaxParallelRequests":0}}`, "server_max_parallel_requests must be positive"},
		{"negative delta", "Write", `{"remoteExecutableConfig":{"requestTimeout":"1s","deltaStage":"-1s","serverMaxParallelRequests":1}}`, "delta_stage must not be negative"},
		{"duration overflow", "Write", `{"remoteExecutableConfig":{"requestTimeout":"9223372036.854775808s","serverMaxParallelRequests":1}}`, "request_timeout exceeds the Go duration range"},
		{"unknown schedule", "Write", `{"remoteExecutableConfig":{"requestTimeout":"1s","serverMaxParallelRequests":1,"transmissionSchedule":99}}`, "unknown transmission_schedule"},
		{"unknown hasher", "Write", `{"remoteExecutableConfig":{"requestTimeout":"1s","serverMaxParallelRequests":1,"requestHasherType":99}}`, "unknown request_hasher_type"},
		{"unknown aggregator", "Trigger", `{"remoteTriggerConfig":{},"aggregatorConfig":{"aggregatorType":99}}`, "unknown aggregator_type"},
		{"negative refresh", "Trigger", `{"remoteTriggerConfig":{"registrationRefresh":"-1s"}}`, "registration_refresh must not be negative"},
		{"negative expiry", "Trigger", `{"remoteTriggerConfig":{"registrationExpiry":"-1s"}}`, "registration_expiry must not be negative"},
		{"negative message expiry", "Trigger", `{"remoteTriggerConfig":{"messageExpiry":"-1s"}}`, "message_expiry must not be negative"},
		{"negative batch period", "Trigger", `{"remoteTriggerConfig":{"batchCollectionPeriod":"-1s"}}`, "batch_collection_period must not be negative"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := fmt.Sprintf(`{"version":2,"dons":{"don-1":{"capabilities":{"evm@1.0.0":{"methodConfigs":{%q:%s}}}}}}`, tc.method, tc.config)
			require.ErrorContains(t, Validate(raw), tc.wantErr)

			gc := New()
			_, err := gc.Store(`{"version":1}`, "valid")
			require.NoError(t, err)
			updates, unsubscribe := gc.Subscribe()
			defer unsubscribe()
			_, err = gc.Store(raw, "invalid")
			require.ErrorContains(t, err, tc.wantErr)
			assert.Equal(t, uint64(1), gc.LoadParsed().GetVersion(), "a rejected update must preserve the applied payload")
			assert.Empty(t, updates, "a rejected update must not trigger reconciliation")
		})
	}
}

func TestMethodConfigsFromProto_InvalidDuration(t *testing.T) {
	t.Parallel()
	// Direct protobuf callers also get validation; AsDuration must not normalize malformed input.
	_, err := MethodConfigsFromProto(map[string]*capabilitiespb.CapabilityMethodConfig{
		"Trigger": {RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh: &durationpb.Duration{Seconds: 1, Nanos: -1},
			},
		}},
	})
	require.ErrorContains(t, err, "invalid registration_refresh")
}

func TestValidate_MethodConfigsTriggerDefaults(t *testing.T) {
	t.Parallel()
	for _, config := range []string{`{"remoteTriggerConfig":{}}`, `{"remoteTriggerConfig":{"registrationRefresh":"0s","batchCollectionPeriod":"0s"}}`} {
		raw := fmt.Sprintf(`{"version":1,"dons":{"don-1":{"capabilities":{"evm@1.0.0":{"methodConfigs":{"Trigger":%s}}}}}}`, config)
		require.NoError(t, Validate(raw))
	}
}

func TestValidate_MethodConfigs(t *testing.T) {
	t.Parallel()

	t.Run("accepts payload with method_configs", func(t *testing.T) {
		t.Parallel()
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"write-chain_evm_1@1.0.0":{"methodConfigs":{"Write":{"remoteExecutableConfig":{"requestTimeout":"30s","serverMaxParallelRequests":10}}}}}}}}`
		require.NoError(t, Validate(raw))
	})

	t.Run("accepts payload without method_configs", func(t *testing.T) {
		t.Parallel()
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"cron@1.0.0":{}}}}}`
		require.NoError(t, Validate(raw))
	})

	t.Run("rejects unknown field inside method config", func(t *testing.T) {
		t.Parallel()
		// Strict protojson parsing rejects unknown fields anywhere in the payload
		// (a misspelled method config key cannot silently become "no config").
		raw := `{"version":1,"dons":{"don-1":{"capabilities":{"write-chain_evm_1@1.0.0":{"methodConfigs":{"Write":{"remoteExecutableConfigTypo":{"requestTimeout":"30s","serverMaxParallelRequests":10}}}}}}}}`
		require.Error(t, Validate(raw))
	})
}
