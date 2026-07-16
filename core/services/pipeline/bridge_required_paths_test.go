package pipeline

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"
)

// TestingSetBridgeRequiredJSONPaths sets required JSON paths on a bridge task
// for tests in external test packages (e.g. pipeline_test).
func TestingSetBridgeRequiredJSONPaths(t *BridgeTask, paths [][]string) {
	t.requiredJSONPaths = paths
}

func TestRequiredJSONPathsFromBridge_DirectEdge(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="data,result"];
b -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Equal(t, [][]string{{"data", "result"}}, paths)
}

func TestRequiredJSONPathsFromBridge_DataVarRoot(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="a,b" data="$(b)"];
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Equal(t, [][]string{{"a", "b"}}, paths)
}

func TestRequiredJSONPathsFromBridge_CustomSeparator(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="a|b|c" separator="|"];
b -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Equal(t, [][]string{{"a", "b", "c"}}, paths)
}

func TestRequiredJSONPathsFromBridge_SkipsLax(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="data,result" lax="true"];
b -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Nil(t, paths)
}

func TestRequiredJSONPathsFromBridge_SkipsDynamicPath(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="$(jobRun.x)"];
b -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Nil(t, paths)
}

func TestRequiredJSONPathsFromBridge_LowestOutputIndexWins(t *testing.T) {
	t.Parallel()

	dot := `
b0 [type=bridge name=a index=0];
b1 [type=bridge name=b index=1];
p [type=jsonparse path="only,b0"];
b0 -> p;
b1 -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var b0, b1 *BridgeTask
	for _, tk := range pl.Tasks {
		switch tk.DotID() {
		case "b0":
			b0 = tk.(*BridgeTask)
		case "b1":
			b1 = tk.(*BridgeTask)
		}
	}
	require.NotNil(t, b0)
	require.NotNil(t, b1)

	require.Equal(t, [][]string{{"only", "b0"}}, b0.getRequiredJSONPaths())
	require.Nil(t, b1.getRequiredJSONPaths())
}

func TestRequiredJSONPathsFromBridge_DedupesPaths(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p1 [type=jsonparse path="data,result"];
p2 [type=jsonparse path="data,result"];
b -> p1;
b -> p2;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Equal(t, [][]string{{"data", "result"}}, paths)
}

func TestRequiredJSONPathsFromBridge_EmptyPathSkipped(t *testing.T) {
	t.Parallel()

	dot := `
b [type=bridge name=testbridge];
p [type=jsonparse path="" ];
b -> p;
`
	pl, err := Parse(dot)
	require.NoError(t, err)

	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)

	paths := bt.getRequiredJSONPaths()
	require.Nil(t, paths)
}

func TestJSONDecodeValidateRequiredPaths(t *testing.T) {
	t.Parallel()

	t.Run("empty_paths_noop", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, jsonDecodeValidateRequiredPaths([]byte(`{}`), nil))
		require.NoError(t, jsonDecodeValidateRequiredPaths([]byte(`{}`), [][]string{}))
	})

	t.Run("json_number_and_string_ok", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":42}}`), [][]string{{"data", "result"}}))
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"123.45"}}`), [][]string{{"data", "result"}}))
	})

	t.Run("large_finite_decimal_json_number", func(t *testing.T) {
		t.Parallel()
		// Decimal accepts magnitudes float64 cannot represent.
		body := []byte(`{"data":{"result":1e309}}`)
		require.NoError(t, jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "result"}}))
	})

	t.Run("missing_path", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"data":{"result":1}}`)
		err := jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "missing"}})
		require.Error(t, err)
		require.ErrorIs(t, err, jsonparser.KeyPathNotFoundError)
	})

	t.Run("null_value", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"data":{"result":null}}`)
		err := jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "result"}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "is null")
	})

	t.Run("scientific_notation_rejected_after_decimal_fails", func(t *testing.T) {
		t.Parallel()
		// Decimal rejects huge exponent; must not be accepted via a hex fallback (digits+e are valid hex).
		body := []byte(`{"data":{"result":1e10000000000}}`)
		err := jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "result"}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid value")
	})

	t.Run("nan_string_rejected", func(t *testing.T) {
		t.Parallel()
		err := jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"NaN"}}`), [][]string{{"data", "result"}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "NaN")

		require.Error(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"nan"}}`), [][]string{{"data", "result"}}))
	})

	t.Run("hex_literals", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"0xff"}}`), [][]string{{"data", "result"}}))
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"0XAbCd"}}`), [][]string{{"data", "result"}}))
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"-0x10"}}`), [][]string{{"data", "result"}}))
	})

	t.Run("large_hex_integer", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"data":{"result":"0x1fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"}}`)
		require.NoError(t, jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "result"}}))
	})

	t.Run("invalid_hex", func(t *testing.T) {
		t.Parallel()
		err := jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"0x"}}`), [][]string{{"data", "result"}})
		require.Error(t, err)
		err = jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"0xGG"}}`), [][]string{{"data", "result"}})
		require.Error(t, err)
	})

	t.Run("non_numeric_string_rejected", func(t *testing.T) {
		t.Parallel()
		err := jsonDecodeValidateRequiredPaths(
			[]byte(`{"data":{"result":"ok"}}`), [][]string{{"data", "result"}})
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid value")
	})

	t.Run("multiple_paths", func(t *testing.T) {
		t.Parallel()
		body := []byte(`{"a":1,"b":"0x2a"}`)
		require.NoError(t, jsonDecodeValidateRequiredPaths(body, [][]string{{"a"}, {"b"}}))
	})

	t.Run("empty_path_segment_skipped", func(t *testing.T) {
		t.Parallel()
		// jsonDecodeValidateRequiredPaths ignores empty segment lists (same as no required paths).
		require.NoError(t, jsonDecodeValidateRequiredPaths(
			[]byte(`{}`), [][]string{{}}))
	})
}

func TestRequiredJSONPathsFromBridge_NilTask(t *testing.T) {
	t.Parallel()
	require.Nil(t, ((*BridgeTask)(nil)).getRequiredJSONPaths())
}

func TestParseBridgeTask_CheckRequired(t *testing.T) {
	t.Parallel()

	dot := `b [type=bridge name=testbridge checkRequired=true];`
	pl, err := Parse(dot)
	require.NoError(t, err)
	var bt *BridgeTask
	for _, tk := range pl.Tasks {
		if tk.DotID() == "b" {
			bt = tk.(*BridgeTask)
			break
		}
	}
	require.NotNil(t, bt)
	require.Equal(t, "true", bt.CheckRequired)
}
