package pipeline_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	configtest2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestTimeoutAttribute(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	timeout, set := p.Tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "10s"), timeout)
	assert.Equal(t, true, set)

	a = `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];`
	p, err = pipeline.Parse(a)
	require.NoError(t, err)
	timeout, set = p.Tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "0s"), timeout)
	assert.Equal(t, false, set)
}

func TestTaskHTTPUnmarshal(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=http allowunrestrictednetworkaccess=true method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	require.Len(t, p.Tasks, 1)

	task := p.Tasks[0].(*pipeline.HTTPTask)
	require.Equal(t, "true", task.AllowUnrestrictedNetworkAccess)
}

func TestTaskAnyUnmarshal(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=any failEarly=true];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	require.Len(t, p.Tasks, 1)
	_, ok := p.Tasks[0].(*pipeline.AnyTask)
	require.True(t, ok)
	require.Equal(t, true, p.Tasks[0].Base().FailEarly)
}

func TestRetryUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    string
		retries uint32
		min     time.Duration
		max     time.Duration
	}{
		{

			"nothing specified",
			`ds1 [type=any];`,
			0,
			time.Second * 5,
			time.Minute,
		},
		{

			"only retry specified",
			`ds1 [type=any retries=5];`,
			5,
			time.Second * 5,
			time.Minute,
		},
		{
			"all params set",
			`ds1 [type=http retries=10 minBackoff="1s" maxBackoff="30m"];`,
			10,
			time.Second,
			time.Minute * 30,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := pipeline.Parse(test.spec)
			require.NoError(t, err)
			require.Len(t, p.Tasks, 1)
			require.Equal(t, test.retries, p.Tasks[0].TaskRetries())
			require.Equal(t, test.min, p.Tasks[0].TaskMinBackoff())
			require.Equal(t, test.max, p.Tasks[0].TaskMaxBackoff())
		})
	}
}

func TestUnmarshalTaskFromMap(t *testing.T) {
	t.Parallel()

	t.Run("returns error if task is not the right type", func(t *testing.T) {
		taskMap := interface{}(nil)
		_, err := pipeline.UnmarshalTaskFromMap(pipeline.TaskType("http"), taskMap, 0, "foo-dot-id")
		require.EqualError(t, err, "UnmarshalTaskFromMap: UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string. Got <nil> (<nil>) of type <nil>")

		taskMap = struct {
			foo time.Time
			bar int
		}{time.Unix(42, 42), 42}
		_, err = pipeline.UnmarshalTaskFromMap(pipeline.TaskType("http"), taskMap, 0, "foo-dot-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "UnmarshalTaskFromMap: UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string")
	})

	t.Run("unknown task type", func(t *testing.T) {
		taskMap := map[string]string{}
		_, err := pipeline.UnmarshalTaskFromMap(pipeline.TaskType("xxx"), taskMap, 0, "foo-dot-id")
		require.EqualError(t, err, `UnmarshalTaskFromMap: unknown task type: "xxx"`)
	})

	tests := []struct {
		taskType         pipeline.TaskType
		expectedTaskType interface{}
	}{
		{pipeline.TaskTypeHTTP, &pipeline.HTTPTask{}},
		{pipeline.TaskTypeBridge, &pipeline.BridgeTask{}},
		{pipeline.TaskTypeMean, &pipeline.MeanTask{}},
		{pipeline.TaskTypeMedian, &pipeline.MedianTask{}},
		{pipeline.TaskTypeMode, &pipeline.ModeTask{}},
		{pipeline.TaskTypeSum, &pipeline.SumTask{}},
		{pipeline.TaskTypeMultiply, &pipeline.MultiplyTask{}},
		{pipeline.TaskTypeDivide, &pipeline.DivideTask{}},
		{pipeline.TaskTypeJSONParse, &pipeline.JSONParseTask{}},
		{pipeline.TaskTypeCBORParse, &pipeline.CBORParseTask{}},
		{pipeline.TaskTypeAny, &pipeline.AnyTask{}},
		{pipeline.TaskTypeVRF, &pipeline.VRFTask{}},
		{pipeline.TaskTypeVRFV2, &pipeline.VRFTaskV2{}},
		{pipeline.TaskTypeEstimateGasLimit, &pipeline.EstimateGasLimitTask{}},
		{pipeline.TaskTypeETHCall, &pipeline.ETHCallTask{}},
		{pipeline.TaskTypeETHTx, &pipeline.ETHTxTask{}},
		{pipeline.TaskTypeETHABIEncode, &pipeline.ETHABIEncodeTask{}},
		{pipeline.TaskTypeETHABIEncode2, &pipeline.ETHABIEncodeTask2{}},
		{pipeline.TaskTypeETHABIDecode, &pipeline.ETHABIDecodeTask{}},
		{pipeline.TaskTypeETHABIDecodeLog, &pipeline.ETHABIDecodeLogTask{}},
		{pipeline.TaskTypeMerge, &pipeline.MergeTask{}},
		{pipeline.TaskTypeLowercase, &pipeline.LowercaseTask{}},
		{pipeline.TaskTypeUppercase, &pipeline.UppercaseTask{}},
		{pipeline.TaskTypeConditional, &pipeline.ConditionalTask{}},
		{pipeline.TaskTypeHexDecode, &pipeline.HexDecodeTask{}},
		{pipeline.TaskTypeBase64Decode, &pipeline.Base64DecodeTask{}},
	}

	for _, test := range tests {
		t.Run(string(test.taskType), func(t *testing.T) {
			taskMap := map[string]string{}
			task, err := pipeline.UnmarshalTaskFromMap(test.taskType, taskMap, 0, "foo-dot-id")
			require.NoError(t, err)
			require.IsType(t, test.expectedTaskType, task)
		})
	}
}

func TestMarshalJSONSerializable_replaceBytesWithHex(t *testing.T) {
	t.Parallel()

	type jsm = map[string]interface{}

	toJSONSerializable := func(val jsm) *pipeline.JSONSerializable {
		return &pipeline.JSONSerializable{
			Valid: true,
			Val:   val,
		}
	}

	var (
		testAddr1 = common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406f111")
		testAddr2 = common.HexToAddress("0x2ab9a2Dc53736b361b72d900CdF9F78F9406f222")
		testHash1 = common.HexToHash("0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071eaf111")
		testHash2 = common.HexToHash("0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071eaf222")
	)

	tests := []struct {
		name     string
		input    *pipeline.JSONSerializable
		expected string
		err      error
	}{
		{"invalid input", &pipeline.JSONSerializable{Valid: false}, "null", nil},
		{"empty object", toJSONSerializable(jsm{}), "{}", nil},
		{"byte slice", toJSONSerializable(jsm{"slice": []byte{0x10, 0x20, 0x30}}),
			`{"slice":"0x102030"}`, nil},
		{"address", toJSONSerializable(jsm{"addr": testAddr1}),
			`{"addr":"0x2aB9a2dc53736B361B72d900cDF9f78f9406f111"}`, nil},
		{"hash", toJSONSerializable(jsm{"hash": testHash1}),
			`{"hash":"0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071eaf111"}`, nil},
		{"slice of byte slice", toJSONSerializable(jsm{"slices": [][]byte{{0x10, 0x11, 0x12}, {0x20, 0x21, 0x22}}}),
			`{"slices":["0x101112","0x202122"]}`, nil},
		{"slice of addresses", toJSONSerializable(jsm{"addresses": []common.Address{testAddr1, testAddr2}}),
			`{"addresses":["0x2aB9a2dc53736B361B72d900cDF9f78f9406f111","0x2aB9A2Dc53736b361b72D900CDf9f78f9406F222"]}`, nil},
		{"slice of hashes", toJSONSerializable(jsm{"hashes": []common.Hash{testHash1, testHash2}}),
			`{"hashes":["0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071eaf111","0x317cfd032b5d6657995f17fe768f7cc4ea0ada27ad421c4caa685a9071eaf222"]}`, nil},
		{"slice of interfaces", toJSONSerializable(jsm{"ifaces": []interface{}{[]byte{0x10, 0x11, 0x12}, []byte{0x20, 0x21, 0x22}}}),
			`{"ifaces":["0x101112","0x202122"]}`, nil},
		{"map", toJSONSerializable(jsm{"map": jsm{"slice": []byte{0x10, 0x11, 0x12}, "addr": testAddr1}}),
			`{"map":{"addr":"0x2aB9a2dc53736B361B72d900cDF9f78f9406f111","slice":"0x101112"}}`, nil},
		{"byte array 4", toJSONSerializable(jsm{"ba4": [4]byte{1, 2, 3, 4}}),
			`{"ba4":"0x01020304"}`, nil},
		{"byte array 8", toJSONSerializable(jsm{"ba8": [8]uint8{1, 2, 3, 4, 5, 6, 7, 8}}),
			`{"ba8":"0x0102030405060708"}`, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := test.input.MarshalJSON()
			assert.Equal(t, test.expected, string(bytes))
			assert.Equal(t, test.err, errors.Cause(err))
		})
	}
}

func TestUnmarshalJSONSerializable(t *testing.T) {
	t.Parallel()

	big, ok := new(big.Int).SetString("18446744073709551616", 10)
	assert.True(t, ok)

	tests := []struct {
		name, input string
		expected    interface{}
	}{
		{"null json", `null`, nil},
		{"bool", `true`, true},
		{"string", `"foo"`, "foo"},
		{"object with int", `{"foo": 42}`, map[string]interface{}{"foo": int64(42)}},
		{"object with float", `{"foo": 3.14}`, map[string]interface{}{"foo": float64(3.14)}},
		{"object with big int", `{"foo": 18446744073709551616}`, map[string]interface{}{"foo": big}},
		{"slice", `[42, 3.14]`, []interface{}{int64(42), float64(3.14)}},
		{"nested map", `{"m": {"foo": 42}}`, map[string]interface{}{"m": map[string]interface{}{"foo": int64(42)}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i pipeline.JSONSerializable
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			if test.expected != nil {
				assert.True(t, i.Valid)
				assert.Equal(t, test.expected, i.Val)
			}
		})
	}
}

func TestCheckInputs(t *testing.T) {
	t.Parallel()

	emptyPR := []pipeline.Result{}
	nonEmptyPR := []pipeline.Result{
		{
			Value: "foo",
			Error: nil,
		},
		{
			Value: "err",
			Error: errors.New("bar"),
		},
	}

	tests := []struct {
		name                      string
		pr                        []pipeline.Result
		minLen, maxLen, maxErrors int
		err                       error
		outputsLen                int
	}{
		{"minLen violation", emptyPR, 1, 0, 0, pipeline.ErrWrongInputCardinality, 0},
		{"maxLen violation", nonEmptyPR, 1, 1, 0, pipeline.ErrWrongInputCardinality, 0},
		{"maxErrors violation", nonEmptyPR, 1, 2, 0, pipeline.ErrTooManyErrors, 0},
		{"ok", nonEmptyPR, 1, 2, 1, nil, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputs, err := pipeline.CheckInputs(test.pr, test.minLen, test.maxLen, test.maxErrors)
			if test.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.outputsLen, len(outputs))
			} else {
				assert.Equal(t, test.err, errors.Cause(err))
			}
		})
	}
}

func TestTaskRunResult_IsPending(t *testing.T) {
	t.Parallel()

	trr := &pipeline.TaskRunResult{}
	assert.True(t, trr.IsPending())

	trrWithResult := &pipeline.TaskRunResult{Result: pipeline.Result{Value: "foo"}}
	assert.False(t, trrWithResult.IsPending())

	trrWithFinishedAt := &pipeline.TaskRunResult{FinishedAt: null.NewTime(time.Now(), true)}
	assert.False(t, trrWithFinishedAt.IsPending())
}

func TestSelectGasLimit(t *testing.T) {
	t.Parallel()

	gcfg := configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(999))
		c.EVM[0].GasEstimator.LimitJobType = v2.GasLimitJobType{
			DR:     ptr(uint32(100)),
			VRF:    ptr(uint32(101)),
			FM:     ptr(uint32(102)),
			OCR:    ptr(uint32(103)),
			Keeper: ptr(uint32(103)),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	t.Run("spec defined gas limit", func(t *testing.T) {
		var specGasLimit uint32 = 1
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.DirectRequestJobType, &specGasLimit)
		assert.Equal(t, uint32(1), gasLimit)
	})

	t.Run("direct request specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.DirectRequestJobType, nil)
		assert.Equal(t, uint32(100), gasLimit)
	})

	t.Run("OCR specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.OffchainReportingJobType, nil)
		assert.Equal(t, uint32(103), gasLimit)
	})

	t.Run("VRF specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.VRFJobType, nil)
		assert.Equal(t, uint32(101), gasLimit)
	})

	t.Run("flux monitor specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.FluxMonitorJobType, nil)
		assert.Equal(t, uint32(102), gasLimit)
	})

	t.Run("keeper specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.KeeperJobType, nil)
		assert.Equal(t, uint32(103), gasLimit)
	})

	t.Run("fallback to default gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg, pipeline.WebhookJobType, nil)
		assert.Equal(t, uint32(999), gasLimit)
	})
}
func TestGetNextTaskOf(t *testing.T) {
	trrs := pipeline.TaskRunResults{
		{
			Task: &pipeline.BridgeTask{
				BaseTask: pipeline.NewBaseTask(1, "t1", nil, nil, 0),
			},
		},
		{
			Task: &pipeline.HTTPTask{
				BaseTask: pipeline.NewBaseTask(2, "t2", nil, nil, 0),
			},
		},
		{
			Task: &pipeline.ETHABIDecodeTask{
				BaseTask: pipeline.NewBaseTask(3, "t3", nil, nil, 0),
			},
		},
		{
			Task: &pipeline.JSONParseTask{
				BaseTask: pipeline.NewBaseTask(4, "t4", nil, nil, 0),
			},
		},
	}

	firstTask := trrs[0]
	nextTask := trrs.GetNextTaskOf(firstTask)
	assert.Equal(t, nextTask.Task.ID(), 2)

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Equal(t, nextTask.Task.ID(), 3)

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Equal(t, nextTask.Task.ID(), 4)

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Empty(t, nextTask)

}
