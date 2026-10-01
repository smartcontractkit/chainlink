package globalconfig

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
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

	t.Run("applies first version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
		raw, v := g.Load()
		assert.Equal(t, `{"version":1}`, raw)
		assert.Equal(t, uint64(1), v)
	})

	t.Run("accepts strictly newer version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
		require.NoError(t, g.Store(Update{Raw: `{"version":2}`, Hash: "h2"}))
		_, v := g.Load()
		assert.Equal(t, uint64(2), v)
	})

	t.Run("rejects non-increasing version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		err := g.Store(Update{Raw: `{"version":5}`, Hash: "hDifferent"})
		require.Error(t, err)
		err = g.Store(Update{Raw: `{"version":4}`, Hash: "h4"})
		require.Error(t, err)
		_, v := g.Load()
		assert.Equal(t, uint64(5), v)
	})

	t.Run("idempotent re-apply of same hash", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
		_, v := g.Load()
		assert.Equal(t, uint64(5), v)
	})

	t.Run("proto-json string version", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":"42"}`, Hash: "h"}))
		_, v := g.Load()
		assert.Equal(t, uint64(42), v)
	})

	t.Run("rejects invalid payload", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.Error(t, g.Store(Update{Raw: ``, Hash: "h"}))
		require.Error(t, g.Store(Update{Raw: `not json`, Hash: "h"}))
		require.Error(t, g.Store(Update{Raw: `{"version":"nope"}`, Hash: "h"}))
	})
}

func TestGlobalConfig_Versioning(t *testing.T) {
	t.Parallel()

	t.Run("zero or missing version is rejected, even as the first payload", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.ErrorContains(t, g.Store(Update{Raw: `{"version":0}`, Hash: "h0"}), "version must be >= 1")
		require.ErrorContains(t, g.Store(Update{Raw: `{"dons":{}}`, Hash: "hx"}), "version must be >= 1")
		raw, v := g.Load()
		assert.Empty(t, raw)
		assert.Equal(t, uint64(0), v)
	})

	t.Run("stale replay of an older payload is rejected", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "a"}))
		require.NoError(t, g.Store(Update{Raw: `{"version":2}`, Hash: "b"}))
		// "a" is no longer the applied hash, so it is a stale update, not an idempotent one.
		require.ErrorContains(t, g.Store(Update{Raw: `{"version":1}`, Hash: "a"}), "not newer")
		raw, v := g.Load()
		assert.Equal(t, `{"version":2}`, raw)
		assert.Equal(t, uint64(2), v)
	})

	t.Run("empty hash falls back to sha256 of the payload", func(t *testing.T) {
		t.Parallel()
		g := New()
		raw := `{"version":3}`
		require.NoError(t, g.Store(Update{Raw: raw}))
		// Same payload again: idempotent (no version error).
		require.NoError(t, g.Store(Update{Raw: raw}))
		// Same payload with the explicit validator hash: also idempotent.
		sum := sha256.Sum256([]byte(raw))
		require.NoError(t, g.Store(Update{Raw: raw, Hash: hex.EncodeToString(sum[:])}))
	})

	t.Run("rejected update keeps the applied payload", func(t *testing.T) {
		t.Parallel()
		g := New()
		require.NoError(t, g.Store(Update{Raw: `{"version":5,"dons":{"1":{}}}`, Hash: "h5"}))
		require.Error(t, g.Store(Update{Raw: `{"version":6,"dons":{"1":{"donId":2}}}`, Hash: "h6"}))
		reg, v := g.LoadParsed()
		assert.Equal(t, uint64(5), v)
		assert.Contains(t, reg.GetDons(), uint32(1))
	})
}

func TestGlobalConfig_Clear(t *testing.T) {
	t.Parallel()

	g := New()
	ch, unsubscribe := g.Subscribe()
	defer unsubscribe()

	require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
	<-ch

	g.Clear()
	raw, v := g.Load()
	assert.Empty(t, raw)
	assert.Equal(t, uint64(0), v)
	reg, v := g.LoadParsed()
	assert.Nil(t, reg)
	assert.Equal(t, uint64(0), v)
	assert.Len(t, ch, 1, "Clear notifies when it withdraws a payload")
	<-ch
	g.Clear()
	assert.Empty(t, ch, "clearing an already cleared config does not notify")

	// The high-water mark survives Clear: older payloads stay rejected...
	require.ErrorContains(t, g.Store(Update{Raw: `{"version":4}`, Hash: "h4"}), "not newer than applied version 5")
	require.ErrorContains(t, g.Store(Update{Raw: `{"version":5,"dons":{}}`, Hash: "h5b"}), "not newer")
	// ...but the last applied payload can be re-applied (delete + recreate of the same job).
	require.NoError(t, g.Store(Update{Raw: `{"version":5}`, Hash: "h5"}))
	_, v = g.Load()
	assert.Equal(t, uint64(5), v)
	<-ch

	g.Clear()
	require.NoError(t, g.Store(Update{Raw: `{"version":6}`, Hash: "h6"}))
	_, v = g.Load()
	assert.Equal(t, uint64(6), v)
}

func TestGlobalConfig_LoadParsedIsASnapshot(t *testing.T) {
	t.Parallel()

	g := New()
	require.NoError(t, g.Store(Update{Raw: `{"version":1,"dons":{"7":{"capabilityConfigs":{"cron@1.0.0":{}}}}}`, Hash: "h1"}))

	snap, _ := g.LoadParsed()
	snap.Version = 99
	snap.Dons[8] = &capabilitiespb.OffchainDONConfig{}
	snap.Dons[7].CapabilityConfigs["evil@1.0.0"] = &capabilitiespb.CapabilityConfig{}

	again, v := g.LoadParsed()
	assert.Equal(t, uint64(1), v)
	assert.Equal(t, uint64(1), again.GetVersion())
	assert.NotContains(t, again.GetDons(), uint32(8))
	assert.NotContains(t, again.GetDons()[7].GetCapabilityConfigs(), "evil@1.0.0")

	empty, v := New().LoadParsed()
	assert.Nil(t, empty)
	assert.Equal(t, uint64(0), v)
}

func TestGlobalConfig_Subscribe(t *testing.T) {
	t.Parallel()

	g := New()
	ch, unsubscribe := g.Subscribe()

	require.NoError(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
	require.NoError(t, g.Store(Update{Raw: `{"version":2}`, Hash: "h2"}))
	// Coalesced: two changes, one pending signal.
	assert.Len(t, ch, 1)
	<-ch

	// Idempotent and rejected stores do not signal.
	require.NoError(t, g.Store(Update{Raw: `{"version":2}`, Hash: "h2"}))
	require.Error(t, g.Store(Update{Raw: `{"version":1}`, Hash: "h1"}))
	assert.Empty(t, ch)

	unsubscribe()
	require.NoError(t, g.Store(Update{Raw: `{"version":3}`, Hash: "h3"}))
	assert.Empty(t, ch)
}

func TestGlobalConfig_ConcurrentStoreAndLoad(t *testing.T) {
	t.Parallel()

	g := New()
	const n = 50
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for v := 1; v <= n; v++ {
			assert.NoError(t, g.Store(Update{Raw: fmt.Sprintf(`{"version":%d,"dons":{"1":{"capabilityConfigs":{"c@1.0.0":{}}}}}`, v)}))
		}
	}()
	for range 2 {
		go func() {
			defer wg.Done()
			for range n {
				reg, v := g.LoadParsed()
				if reg != nil {
					// The snapshot is internally consistent and caller-owned.
					assert.Equal(t, v, reg.GetVersion())
					reg.Dons[1].CapabilityConfigs["x"] = nil
				}
			}
		}()
	}
	wg.Wait()
	_, v := g.Load()
	assert.Equal(t, uint64(n), v)
}

func TestValidate(t *testing.T) {
	t.Parallel()
	require.NoError(t, Validate(`{"version":1,"dons":{}}`))
	require.NoError(t, Validate(`{"version":"1"}`))
	require.NoError(t, Validate(`{"version":1,"dons":{"7":{"donId":7,"capabilityConfigs":{"cron@1.0.0":{}}}}}`))
	require.NoError(t, Validate(`{"version":1,"unknownField":true}`), "unknown fields are tolerated")
	require.Error(t, Validate(``))
	require.Error(t, Validate(`{`))
	require.ErrorContains(t, Validate(`{"version":0}`), "version must be >= 1")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"7":{"donId":8}}}`), "does not match map key")
	require.ErrorContains(t, Validate(`{"version":1,"dons":{"7":{"capabilityConfigs":{"":{}}}}}`), "empty capability ID")

	t.Run("malformed spec_config is rejected at ingestion", func(t *testing.T) {
		t.Parallel()
		reg := &capabilitiespb.OffchainCapabilitiesRegistry{Version: 1, Dons: map[uint32]*capabilitiespb.OffchainDONConfig{
			7: {CapabilityConfigs: map[string]*capabilitiespb.CapabilityConfig{
				"cron@1.0.0": {SpecConfig: &valuespb.Map{Fields: map[string]*valuespb.Value{"broken": {}}}},
			}},
		}}
		b, err := protojson.Marshal(reg)
		require.NoError(t, err)
		require.ErrorContains(t, Validate(string(b)), `capability_configs["cron@1.0.0"]: invalid spec_config: key "broken": value has no type`)

		withSpec := func(spec string) string {
			return `{"version":1,"dons":{"7":{"capabilityConfigs":{"c@1.0.0":{"specConfig":` + spec + `}}}}}`
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
