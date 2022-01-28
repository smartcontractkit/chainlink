package pipeline_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_TaskHTTPUnmarshal(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=http allowunrestrictednetworkaccess=true method=GET url="https://chain.link/voter_turnout/USA-2020" requestData=<{"hi": "hello"}> timeout="10s"];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	require.Len(t, p.Tasks, 1)

	task := p.Tasks[0].(*pipeline.HTTPTask)
	require.Equal(t, "true", task.AllowUnrestrictedNetworkAccess)
}

func Test_TaskAnyUnmarshal(t *testing.T) {
	t.Parallel()

	a := `ds1 [type=any failEarly=true];`
	p, err := pipeline.Parse(a)
	require.NoError(t, err)
	require.Len(t, p.Tasks, 1)
	_, ok := p.Tasks[0].(*pipeline.AnyTask)
	require.True(t, ok)
	require.Equal(t, true, p.Tasks[0].Base().FailEarly)
}

func Test_RetryUnmarshal(t *testing.T) {
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
		p, err := pipeline.Parse(test.spec)
		require.NoError(t, err)
		require.Len(t, p.Tasks, 1)
		require.Equal(t, test.retries, p.Tasks[0].TaskRetries())
		require.Equal(t, test.min, p.Tasks[0].TaskMinBackoff())
		require.Equal(t, test.max, p.Tasks[0].TaskMaxBackoff())
	}

}

func Test_UnmarshalTaskFromMap(t *testing.T) {
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
}

func TestUnmarshalJSONSerializable_Valid(t *testing.T) {
	tests := []struct {
		name, input string
		expected    interface{}
	}{
		{"bool", `true`, true},
		{"string", `"foo"`, "foo"},
		{"raw", `{"foo": 42}`, map[string]interface{}{"foo": float64(42)}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i pipeline.JSONSerializable
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.True(t, i.Valid)
			assert.Equal(t, test.expected, i.Val)
		})
	}
}

func TestUnmarshalJSONSerializable_Invalid(t *testing.T) {
	tests := []struct {
		name, input string
	}{
		{"null json", `null`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var i pipeline.JSONSerializable
			err := json.Unmarshal([]byte(test.input), &i)
			require.NoError(t, err)
			assert.False(t, i.Valid)
		})
	}
}

func TestCheckInputs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		inputs []pipeline.Result
		expErr bool
	}{
		{
			name:   "nil",
			inputs: nil,
		},
		{
			name:   "empty",
			inputs: []pipeline.Result{},
		},
		{
			name:   "empty-val",
			inputs: []pipeline.Result{{}},
		},
		{
			name:   "single",
			inputs: []pipeline.Result{{Value: "testval"}},
		},
		{
			name:   "multi",
			inputs: []pipeline.Result{{Value: "testval"}, {Value: float64(1.1)}},
		},
		{
			name:   "single-err",
			inputs: []pipeline.Result{{Error: errors.New("test error")}},
			expErr: true,
		},
		{
			name:   "multi-errs",
			inputs: []pipeline.Result{{Error: errors.New("test error")}, {Error: errors.New("test error")}},
			expErr: true,
		},
		{
			name:   "multi-val-and-errs",
			inputs: []pipeline.Result{{Value: "testval"}, {Error: errors.New("test error")}},
			expErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := pipeline.CheckInputs(tt.inputs)
			if err != nil && !tt.expErr {
				t.Errorf("unexpected error: %v", err)
			} else if tt.expErr {
				if err == nil {
					t.Error("error expected, but got none")
				} else if !errors.Is(err, pipeline.ErrTooManyErrors) {
					t.Errorf("%T error expected, but got none", pipeline.ErrTooManyErrors)
				}
			}
		})
	}
}

func TestCheckInputsLen(t *testing.T) {
	for _, tt := range []struct {
		name     string
		inputs   []pipeline.Result
		min, max int
		expErr   error
	}{
		{
			name:   "nil-optional",
			inputs: nil,
			min:    0, max: 1,
		},
		{
			name:   "nil-required",
			inputs: nil,
			min:    1, max: 1,
			expErr: pipeline.ErrWrongInputCardinality,
		},
		{
			name:   "empty-optional",
			inputs: []pipeline.Result{},
			min:    0, max: 1,
		},
		{
			name:   "empty-required",
			inputs: []pipeline.Result{},
			min:    1, max: 1,
			expErr: pipeline.ErrWrongInputCardinality,
		},
		{
			name:   "single-optional",
			inputs: []pipeline.Result{{Value: "testval"}},
			min:    0, max: 1,
		},
		{
			name:   "single-required",
			inputs: []pipeline.Result{{Value: "testval"}},
			min:    1, max: 1,
		},
		{
			name:   "multi-optional",
			inputs: []pipeline.Result{{Value: "testval"}, {Value: float64(1.1)}},
			min:    0, max: 2,
		},
		{
			name:   "multi-required",
			inputs: []pipeline.Result{{Value: "testval"}, {Value: float64(1.1)}},
			min:    2, max: 2,
		},
		{
			name:   "multi-optional-short",
			inputs: []pipeline.Result{{Value: "testval"}},
			min:    0, max: 2,
		},
		{
			name:   "multi-required-short",
			inputs: []pipeline.Result{{Value: "testval"}},
			min:    2, max: 2,
			expErr: pipeline.ErrWrongInputCardinality,
		},
		{
			name:   "single-long",
			inputs: []pipeline.Result{{Value: "testval"}, {Value: float64(1.1)}},
			min:    1, max: 1,
			expErr: pipeline.ErrWrongInputCardinality,
		},
		{
			name:   "single-error",
			inputs: []pipeline.Result{{Error: errors.New("test error")}},
			min:    1, max: 1,
			expErr: pipeline.ErrTooManyErrors,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := pipeline.CheckInputsLen(tt.inputs, tt.min, tt.max)
			if err != nil && tt.expErr == nil {
				t.Errorf("unexpected error: %v", err)
			} else if tt.expErr != nil {
				if err == nil {
					t.Error("error expected, but got none")
				} else if !errors.Is(err, tt.expErr) {
					t.Errorf("%T error expected, but got none", tt.expErr)
				}
			}
		})
	}
}
