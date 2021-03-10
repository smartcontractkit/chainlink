package pipeline_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"

	"github.com/bmizerany/assert"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/require"
)

func TestTimeoutAttribute(t *testing.T) {
	t.Parallel()

	g := pipeline.NewTaskDAG()
	a := `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}" timeout="10s"];`
	err := g.UnmarshalText([]byte(a))
	require.NoError(t, err)
	tasks, err := g.TasksInDependencyOrder()
	require.NoError(t, err)
	timeout, set := tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "10s"), timeout)
	assert.Equal(t, true, set)

	g = pipeline.NewTaskDAG()
	a = `ds1 [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}"];`
	err = g.UnmarshalText([]byte(a))
	require.NoError(t, err)
	tasks, err = g.TasksInDependencyOrder()
	require.NoError(t, err)
	timeout, set = tasks[0].TaskTimeout()
	assert.Equal(t, cltest.MustParseDuration(t, "0s"), timeout)
	assert.Equal(t, false, set)
}

func Test_TaskHTTPUnmarshal(t *testing.T) {
	t.Parallel()

	g := pipeline.NewTaskDAG()
	a := `ds1 [type=http allowunrestrictednetworkaccess=true method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}" timeout="10s"];`
	err := g.UnmarshalText([]byte(a))
	require.NoError(t, err)
	tasks, err := g.TasksInDependencyOrder()
	require.NoError(t, err)
	require.Len(t, tasks, 1)

	task := tasks[0].(*pipeline.HTTPTask)
	require.Equal(t, pipeline.MaybeBoolTrue, task.AllowUnrestrictedNetworkAccess)
}

func Test_TaskAnyUnmarshal(t *testing.T) {
	t.Parallel()

	g := pipeline.NewTaskDAG()
	a := `ds1 [type=any];`
	err := g.UnmarshalText([]byte(a))
	require.NoError(t, err)
	tasks, err := g.TasksInDependencyOrder()
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	_, ok := tasks[0].(*pipeline.AnyTask)
	require.True(t, ok)
}

func Test_UnmarshalTaskFromMap(t *testing.T) {
	t.Parallel()

	t.Run("returns error if task is not the right type", func(t *testing.T) {
		taskMap := interface{}(nil)
		_, err := pipeline.UnmarshalTaskFromMap(pipeline.TaskType("http"), taskMap, "foo-dot-id", nil, nil, nil, 0)
		require.EqualError(t, err, "UnmarshalTaskFromMap: UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string. Got <nil> (<nil>) of type <nil>")

		taskMap = struct {
			foo time.Time
			bar int
		}{time.Unix(42, 42), 42}
		_, err = pipeline.UnmarshalTaskFromMap(pipeline.TaskType("http"), taskMap, "foo-dot-id", nil, nil, nil, 0)
		require.Error(t, err)
		require.Contains(t, err.Error(), "UnmarshalTaskFromMap: UnmarshalTaskFromMap only accepts a map[string]interface{} or a map[string]string")
	})
}
