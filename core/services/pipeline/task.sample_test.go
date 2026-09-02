package pipeline_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestSampleTask(t *testing.T) {
	t.Parallel()

	inputMap := map[string]any{
		"data": map[string]any{
			"binance": map[string]any{
				"mid": 399.30,
				"bid": 399.20,
				"ask": 399.40,
			},
		},
		"timestamps": map[string]any{
			"providerIndicatedTimeUnixMs": 1725000000000,
		},
	}

	t.Run("from map input", func(t *testing.T) {
		task := pipeline.SampleTask{
			BaseTask:  pipeline.NewBaseTask(0, "s", nil, nil, 0),
			Source:    "binance",
			Weight:    "0.333",
			Unit:      "USDT",
			ValuePath: "data,binance,mid",
			TsPath:    "timestamps,providerIndicatedTimeUnixMs",
		}
		out, runInfo := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: inputMap}})
		assert.False(t, runInfo.IsPending)
		require.NoError(t, out.Error)
		s, ok := out.Value.(pipeline.Sample)
		require.True(t, ok)
		assert.Equal(t, "binance", s.Source)
		assert.True(t, s.Value.Equal(decimal.RequireFromString("399.30")))
		assert.Equal(t, int64(1725000000000), s.TsMs)
		assert.True(t, s.Weight.Equal(decimal.RequireFromString("0.333")))
		assert.Equal(t, "USDT", s.Unit)
	})

	t.Run("from JSON string input", func(t *testing.T) {
		task := pipeline.SampleTask{
			BaseTask:  pipeline.NewBaseTask(0, "s", nil, nil, 0),
			Source:    "binance",
			Weight:    "1",
			ValuePath: "data,binance,mid",
			TsPath:    "timestamps,providerIndicatedTimeUnixMs",
		}
		inputJSON := `{"data":{"binance":{"mid":399.30}},"timestamps":{"providerIndicatedTimeUnixMs":1725000000000}}`
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: inputJSON}})
		require.NoError(t, out.Error)
		s := out.Value.(pipeline.Sample)
		assert.True(t, s.Value.Equal(decimal.RequireFromString("399.30")))
		assert.Equal(t, int64(1725000000000), s.TsMs)
	})

	t.Run("missing tsPath errors", func(t *testing.T) {
		task := pipeline.SampleTask{
			BaseTask:  pipeline.NewBaseTask(0, "s", nil, nil, 0),
			Source:    "binance",
			Weight:    "1",
			ValuePath: "data,binance,mid",
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: inputMap}})
		require.Error(t, out.Error)
	})

	t.Run("tsUnit seconds", func(t *testing.T) {
		task := pipeline.SampleTask{
			BaseTask:  pipeline.NewBaseTask(0, "s", nil, nil, 0),
			Source:    "binance",
			Weight:    "1",
			ValuePath: "data,binance,mid",
			TsPath:    "timestamps,ts",
			TsUnit:    "s",
		}
		input := map[string]any{
			"data":       map[string]any{"binance": map[string]any{"mid": 100.0}},
			"timestamps": map[string]any{"ts": 1725000000},
		}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: input}})
		require.NoError(t, out.Error)
		s := out.Value.(pipeline.Sample)
		assert.Equal(t, int64(1725000000000), s.TsMs)
	})

	t.Run("default weight is 1", func(t *testing.T) {
		task := pipeline.SampleTask{
			BaseTask:  pipeline.NewBaseTask(0, "s", nil, nil, 0),
			Source:    "x",
			ValuePath: "data,mid",
			TsPath:    "ts",
		}
		input := map[string]any{"data": map[string]any{"mid": 100.0}, "ts": 1725000000000}
		out, _ := task.Run(t.Context(), logger.TestLogger(t), pipeline.NewVarsFrom(nil), []pipeline.Result{{Value: input}})
		require.NoError(t, out.Error)
		s := out.Value.(pipeline.Sample)
		assert.True(t, s.Weight.Equal(decimal.NewFromInt(1)))
	})
}
