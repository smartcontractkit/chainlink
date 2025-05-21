package pipeline_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestAtrributesAttribute(t *testing.T) {
	a := `ds1 [type=http method=GET tags=<{"attribute1":"value1", "attribute2":42}>];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	task := p.Tasks[0]
	assert.JSONEq(t, "{\"attribute1\":\"value1\", \"attribute2\":42}", task.TaskTags())
}

func TestTimeoutAttribute(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	timeout, set := p.Tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "10s"), timeout)
	assert.True(t, set)

	a = `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}>];`
	p, err = pipeline.Parse(a)
	require.NoError(t, err)
	timeout, set = p.Tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "0s"), timeout)
	assert.False(t, set)
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
	require.True(t, p.Tasks[0].Base().FailEarly)
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
		{pipeline.TaskTypeVRFV2Plus, &pipeline.VRFTaskV2Plus{}},
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
				assert.Len(t, outputs, test.outputsLen)
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

	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(999))
		c.EVM[0].GasEstimator.LimitJobType = toml.GasLimitJobType{
			DR:     ptr(uint32(100)),
			VRF:    ptr(uint32(101)),
			FM:     ptr(uint32(102)),
			OCR:    ptr(uint32(103)),
			Keeper: ptr(uint32(104)),
			OCR2:   ptr(uint32(105)),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	t.Run("spec defined gas limit", func(t *testing.T) {
		var specGasLimit uint32 = 1
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.DirectRequestJobType, &specGasLimit)
		assert.Equal(t, uint64(1), gasLimit)
	})

	t.Run("direct request specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.DirectRequestJobType, nil)
		assert.Equal(t, uint64(100), gasLimit)
	})

	t.Run("OCR specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.OffchainReportingJobType, nil)
		assert.Equal(t, uint64(103), gasLimit)
	})

	t.Run("OCR2 specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.OffchainReporting2JobType, nil)
		assert.Equal(t, uint64(105), gasLimit)
	})

	t.Run("VRF specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.VRFJobType, nil)
		assert.Equal(t, uint64(101), gasLimit)
	})

	t.Run("flux monitor specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.FluxMonitorJobType, nil)
		assert.Equal(t, uint64(102), gasLimit)
	})

	t.Run("keeper specific gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.KeeperJobType, nil)
		assert.Equal(t, uint64(104), gasLimit)
	})

	t.Run("fallback to default gas limit", func(t *testing.T) {
		gasLimit := pipeline.SelectGasLimit(cfg.EVM().GasEstimator(), pipeline.WebhookJobType, nil)
		assert.Equal(t, uint64(999), gasLimit)
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
	assert.Equal(t, 2, nextTask.Task.ID())

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Equal(t, 3, nextTask.Task.ID())

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Equal(t, 4, nextTask.Task.ID())

	nextTask = trrs.GetNextTaskOf(*nextTask)
	assert.Empty(t, nextTask)
}

func TestGetDescendantTasks(t *testing.T) {
	t.Parallel()

	t.Run("GetDescendantTasks with multiple levels of tasks", func(t *testing.T) {
		l3T2 := pipeline.AnyTask{
			BaseTask: pipeline.NewBaseTask(6, "l3T2", nil, nil, 1),
		}
		l3T1 := pipeline.MedianTask{
			BaseTask: pipeline.NewBaseTask(5, "l3T1", nil, nil, 1),
		}
		l2T1 := pipeline.MultiplyTask{
			BaseTask: pipeline.NewBaseTask(4, "l2T1", nil, []pipeline.Task{&l3T1, &l3T2}, 1),
		}
		l1T1 := pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(3, "l1T1", nil, []pipeline.Task{&l2T1}, 2),
		}
		l1T2 := pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(2, "l1T2", nil, nil, 3),
		}
		l1T3 := pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "l1T3", nil, nil, 4),
		}

		baseTask := pipeline.BridgeTask{
			Name:     "bridge-task",
			BaseTask: pipeline.NewBaseTask(0, "baseTask", nil, []pipeline.Task{&l1T1, &l1T2, &l1T3}, 0),
		}

		descendents := baseTask.GetDescendantTasks()
		assert.Len(t, descendents, 6)
	})

	t.Run("GetDescendantTasks with duplicate tasks defined", func(t *testing.T) {
		l2T1 := pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(2, "l1T2", nil, nil, 3),
		}
		l1T1 := pipeline.JSONParseTask{
			BaseTask: pipeline.NewBaseTask(1, "l1T2", nil, []pipeline.Task{&l2T1, &l2T1, &l2T1}, 3),
		}
		taskWithRepeats := pipeline.BridgeTask{
			Name:     "bridge-task",
			BaseTask: pipeline.NewBaseTask(0, "taskWithRepeats", nil, []pipeline.Task{&l1T1, &l1T1, &l1T1}, 0),
		}
		descendents := taskWithRepeats.GetDescendantTasks()
		assert.Len(t, descendents, 2)
	})

	t.Run("GetDescendantTasks with nil output tasks", func(t *testing.T) {
		taskWithRepeats := pipeline.BridgeTask{
			Name:     "bridge-task",
			BaseTask: pipeline.NewBaseTask(0, "taskWithRepeats", nil, nil, 0),
		}
		descendents := taskWithRepeats.GetDescendantTasks()
		assert.Empty(t, descendents)
	})

	t.Run("GetDescendantTasks with empty list of output tasks", func(t *testing.T) {
		taskWithRepeats := pipeline.BridgeTask{
			Name:     "bridge-task",
			BaseTask: pipeline.NewBaseTask(0, "taskWithRepeats", nil, []pipeline.Task{}, 0),
		}
		descendents := taskWithRepeats.GetDescendantTasks()
		assert.Empty(t, descendents)
	})
}
