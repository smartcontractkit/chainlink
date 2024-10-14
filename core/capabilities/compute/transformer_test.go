package compute

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		var gotValue string
		gotValue, err = popValue[string](m, "test")
		assert.NoError(t, err)
		assert.Equal(t, "value", gotValue)
	})

	t.Run("not found", func(t *testing.T) {
		_, err = popValue[string](m, "foo")
		var nfe *NotFoundError
		assert.ErrorAs(t, err, &nfe)
	})

	t.Run("type mismatch", func(t *testing.T) {
		_, err = popValue[string](m, "mismatch")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "could not unwrap value")
	})

	assert.Len(t, m.Underlying, 0)
}

func Test_popOptionalValue(t *testing.T) {
	m, err := values.NewMap(
		map[string]any{
			"test": "value",
			"buzz": "fizz",
		},
	)
	assert.NoError(t, err)
	t.Run("found value", func(t *testing.T) {
		var gotValue string
		gotValue, err = popOptionalValue[string](m, "test")
		assert.NoError(t, err)
		assert.Equal(t, "value", gotValue)
	})

	t.Run("not found returns nil error", func(t *testing.T) {
		var gotValue string
		gotValue, err = popOptionalValue[string](m, "foo")
		assert.NoError(t, err)
		assert.Zero(t, gotValue)
	})

	t.Run("some other error fails", func(t *testing.T) {
		var gotValue int
		gotValue, err = popOptionalValue[int](m, "buzz")
		assert.Error(t, err)
		assert.Zero(t, gotValue)
	})

	assert.Len(t, m.Underlying, 0)
}

func Test_transformer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		lgger := logger.TestLogger(t)
		giveMap, err := values.NewMap(map[string]any{
			"maxMemoryMBs": 1024,
			"timeout":      "4s",
			"tickInterval": "8s",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		assert.NoError(t, err)

		wantTO := 4 * time.Second
		wantConfig := &ParsedConfig{
			Binary: []byte{0x01, 0x02, 0x03},
			Config: []byte{0x04, 0x05, 0x06},
			ModuleConfig: &host.ModuleConfig{
				MaxMemoryMBs: 1024,
				Timeout:      &wantTO,
				TickInterval: 8 * time.Second,
				Logger:       lgger,
			},
		}

		tf := NewTransformer()
		gotConfig, err := tf.Transform(giveMap, WithLogger(lgger))

		assert.NoError(t, err)
		assert.Equal(t, wantConfig, gotConfig)
	})

	t.Run("success missing optional fields", func(t *testing.T) {
		lgger := logger.TestLogger(t)
		giveMap, err := values.NewMap(map[string]any{
			"binary": []byte{0x01, 0x02, 0x03},
			"config": []byte{0x04, 0x05, 0x06},
		})
		assert.NoError(t, err)

		wantConfig := &ParsedConfig{
			Binary: []byte{0x01, 0x02, 0x03},
			Config: []byte{0x04, 0x05, 0x06},
			ModuleConfig: &host.ModuleConfig{
				Logger: lgger,
			},
		}

		tf := NewTransformer()
		gotConfig, err := tf.Transform(giveMap, WithLogger(lgger))

		assert.NoError(t, err)
		assert.Equal(t, wantConfig, gotConfig)
	})

	t.Run("fails parsing timeout", func(t *testing.T) {
		lgger := logger.TestLogger(t)
		giveMap, err := values.NewMap(map[string]any{
			"timeout": "not a duration",
			"binary":  []byte{0x01, 0x02, 0x03},
			"config":  []byte{0x04, 0x05, 0x06},
		})
		assert.NoError(t, err)

		tf := NewTransformer()
		_, err = tf.Transform(giveMap, WithLogger(lgger))

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid request")
	})

	t.Run("fails parsing tick interval", func(t *testing.T) {
		lgger := logger.TestLogger(t)
		giveMap, err := values.NewMap(map[string]any{
			"tickInterval": "not a duration",
			"binary":       []byte{0x01, 0x02, 0x03},
			"config":       []byte{0x04, 0x05, 0x06},
		})
		assert.NoError(t, err)

		tf := NewTransformer()
		_, err = tf.Transform(giveMap, WithLogger(lgger))

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid request")
	})
}
