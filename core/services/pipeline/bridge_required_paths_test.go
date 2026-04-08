package pipeline

import (
	"errors"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/require"
)

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

	paths := RequiredJSONPathsFromBridge(bt)
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

	paths := RequiredJSONPathsFromBridge(bt)
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

	paths := RequiredJSONPathsFromBridge(bt)
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

	paths := RequiredJSONPathsFromBridge(bt)
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

	paths := RequiredJSONPathsFromBridge(bt)
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

	require.Equal(t, [][]string{{"only", "b0"}}, RequiredJSONPathsFromBridge(b0))
	require.Nil(t, RequiredJSONPathsFromBridge(b1))
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

	paths := RequiredJSONPathsFromBridge(bt)
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

	paths := RequiredJSONPathsFromBridge(bt)
	require.Nil(t, paths)
}

func TestJSONDecodeValidateRequiredPaths(t *testing.T) {
	t.Parallel()

	// body := []byte(`{"data":{"result":"ok"}}`)
	// require.NoError(t, jsonDecodeValidateRequiredPaths(body, [][]string{{"data", "result"}}))

	numBody := []byte(`{"data":{"result":42}}`)
	require.NoError(t, jsonDecodeValidateRequiredPaths(numBody, [][]string{{"data", "result"}}))

	// Large but finite JSON number: decimal accepts values float64 would overflow.
	largeNumBody := []byte(`{"data":{"result":1e309}}`)
	require.NoError(t, jsonDecodeValidateRequiredPaths(largeNumBody, [][]string{{"data", "result"}}))

	err := jsonDecodeValidateRequiredPaths(largeNumBody, [][]string{{"data", "missing"}})
	require.Error(t, err)
	require.True(t, errors.Is(err, jsonparser.KeyPathNotFoundError))

	nullBody := []byte(`{"data":{"result":null}}`)
	err = jsonDecodeValidateRequiredPaths(nullBody, [][]string{{"data", "result"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "is null")

	// Exponent too large for shopspring/decimal (ParseInt exponent as int32).
	invalidDecimalBody := []byte(`{"data":{"result":1e10000000000}}`)
	err = jsonDecodeValidateRequiredPaths(invalidDecimalBody, [][]string{{"data", "result"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid decimal")

	nanStrBody := []byte(`{"data":{"result":"NaN"}}`)
	err = jsonDecodeValidateRequiredPaths(nanStrBody, [][]string{{"data", "result"}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "NaN")

	nanStrBodyLower := []byte(`{"data":{"result":"nan"}}`)
	require.Error(t, jsonDecodeValidateRequiredPaths(nanStrBodyLower, [][]string{{"data", "result"}}))
}

func TestRequiredJSONPathsFromBridge_NilTask(t *testing.T) {
	t.Parallel()
	require.Nil(t, RequiredJSONPathsFromBridge(nil))
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
