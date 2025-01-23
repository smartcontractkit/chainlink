package compute

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NotFoundError(t *testing.T) {
	nfe := NewNotFoundError("test")
	assert.Equal(t, "could not find \"test\" in map", nfe.Error())
}

func Test_popValue(t *testing.T) {
	m, err := values.NewMap(
		map[string]any{
			"test":     "value",
			"mismatch": 42,
		},
	)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		var gotValue string
		gotValue, err = popValue[string](m, "test")
		require.NoError(t, err)
		assert.Equal(t, "value", gotValue)
	})

	t.Run("not found", func(t *testing.T) {
		_, err = popValue[string](m, "foo")
		var nfe *NotFoundError
		require.ErrorAs(t, err, &nfe)
	})

	t.Run("type mismatch", func(t *testing.T) {
		_, err = popValue[string](m, "mismatch")
		require.Error(t, err)
		require.ErrorContains(t, err, "could not unwrap value")
	})

	assert.Empty(t, m.Underlying)
}

func Test_popOptionalValue(t *testing.T) {
	m, err := values.NewMap(
		map[string]any{
			"test": "value",
			"buzz": "fizz",
		},
	)
	require.NoError(t, err)
	t.Run("found value", func(t *testing.T) {
		var gotValue string
		gotValue, err = popOptionalValue[string](m, "test")
		require.NoError(t, err)
		assert.Equal(t, "value", gotValue)
	})

	t.Run("not found returns nil error", func(t *testing.T) {
		var gotValue string
		gotValue, err = popOptionalValue[string](m, "foo")
		require.NoError(t, err)
		assert.Zero(t, gotValue)
	})

	t.Run("some other error fails", func(t *testing.T) {
		var gotValue int
		gotValue, err = popOptionalValue[int](m, "buzz")
		require.Error(t, err)
		assert.Zero(t, gotValue)
	})

	assert.Empty(t, m.Underlying)
}

func Test_transformer(t *testing.T) {
	var (
		lgger   = logger.TestLogger(t)
		emitter = custmsg.NewLabeler()
	)
	t.Run("success", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"maxMemoryMBs": 1024,
			"timeout":      "4s",
			"tickInterval": "8s",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		wantTO := 4 * time.Second
		wantConfig := &ParsedConfig{
			Binary: []byte{0x01, 0x02, 0x03},
			Config: []byte{0x04, 0x05, 0x06},
			ModuleConfig: &host.ModuleConfig{
				MaxMemoryMBs: 1024,
				Timeout:      &wantTO,
				TickInterval: 8 * time.Second,
				Logger:       lgger,
				Labeler:      emitter,
			},
		}

		config := Config{
			MaxMemoryMBs:    2048,
			MaxTimeout:      20 * time.Second,
			MaxTickInterval: 10 * time.Second,
		}
		tf := NewTransformer(lgger, emitter, config)
		_, gotConfig, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, wantConfig, gotConfig)
	})

	t.Run("success missing optional fields", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"binary": []byte{0x01, 0x02, 0x03},
			"config": []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		timeout := defaultMaxTimeout
		wantConfig := &ParsedConfig{
			Binary: []byte{0x01, 0x02, 0x03},
			Config: []byte{0x04, 0x05, 0x06},
			ModuleConfig: &host.ModuleConfig{
				Logger:       lgger,
				Labeler:      emitter,
				TickInterval: defaultMaxTickInterval,
				Timeout:      &timeout,
				MaxMemoryMBs: defaultMaxMemoryMBs,
			},
		}

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, gotConfig, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, wantConfig, gotConfig)
	})

	t.Run("fails parsing timeout", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"timeout": "not a duration",
			"binary":  []byte{0x01, 0x02, 0x03},
			"config":  []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, _, err = tf.Transform(giveReq)

		require.Error(t, err)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("fails parsing tick interval", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"tickInterval": "not a duration",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, _, err = tf.Transform(giveReq)

		require.Error(t, err)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("invalid tickInterval, applies default", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"tickInterval": "-50ms",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, pc, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, defaultMaxTickInterval, pc.ModuleConfig.TickInterval)
	})

	t.Run("invalid timeout, applies default", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"timeout": "-50ms",
			"binary":  []byte{0x01, 0x02, 0x03},
			"config":  []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, pc, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, defaultMaxTimeout, *pc.ModuleConfig.Timeout)
	})

	t.Run("timeout too high, applies default", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"timeout": "1h",
			"binary":  []byte{0x01, 0x02, 0x03},
			"config":  []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, pc, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, defaultMaxTimeout, *pc.ModuleConfig.Timeout)
	})

	t.Run("tickInterval too high, applies default", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"tickInterval": "1h",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, pc, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, defaultMaxTickInterval, pc.ModuleConfig.TickInterval)
	})

	t.Run("applies default tick interval if missing", func(t *testing.T) {
		giveMap, err := values.NewMap(map[string]any{
			"binary": []byte{0x01, 0x02, 0x03},
			"config": []byte{0x04, 0x05, 0x06},
		})
		giveReq := capabilities.CapabilityRequest{
			Config: giveMap,
		}
		require.NoError(t, err)

		config := Config{}
		config.ApplyDefaults()
		tf := NewTransformer(lgger, emitter, config)
		_, pc, err := tf.Transform(giveReq)

		require.NoError(t, err)
		assert.Equal(t, defaultMaxTickInterval, pc.ModuleConfig.TickInterval)
	})
}
