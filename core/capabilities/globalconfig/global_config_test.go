package globalconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
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
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"don-7":{"capabilities":{"c@1.0.0":{"ocr3Configs":{"default":{"f":1}}}}}}}`), "only spec_config is supported")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"don-7":{"capabilities":{"c@1.0.0":{"localOnly":true}}}}}`), "only spec_config is supported")

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
