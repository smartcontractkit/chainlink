package pipeline_test

import (
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
