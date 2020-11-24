package pipeline_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/require"
)

func Test_TaskHTTPUnmarshal(t *testing.T) {
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
