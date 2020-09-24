package pipeline_test

import (
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func TestGraph_Decode(t *testing.T) {
	expected := map[string]map[string]bool{
		"ds1": {
			"ds1":          false,
			"ds1_parse":    true,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      false,
			"answer2":      false,
		},
		"ds1_parse": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": true,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      false,
			"answer2":      false,
		},
		"ds1_multiply": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      true,
			"answer2":      false,
		},
		"ds2": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    true,
			"ds2_multiply": false,
			"answer1":      false,
			"answer2":      false,
		},
		"ds2_parse": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": true,
			"answer1":      false,
			"answer2":      false,
		},
		"ds2_multiply": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      true,
			"answer2":      false,
		},
		"answer1": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      false,
			"answer2":      false,
		},
		"answer2": {
			"ds1":          false,
			"ds1_parse":    false,
			"ds1_multiply": false,
			"ds2":          false,
			"ds2_parse":    false,
			"ds2_multiply": false,
			"answer1":      false,
			"answer2":      false,
		},
	}

	g := pipeline.NewTaskDAG()
	err := g.UnmarshalText([]byte(dotStr))
	require.NoError(t, err)

	nodes := make(map[string]int64)
	iter := g.Nodes()
	for iter.Next() {
		n := iter.Node().(interface {
			graph.Node
			DOTID() string
		})
		nodes[n.DOTID()] = n.ID()
	}

	for from, connections := range expected {
		for to, connected := range connections {
			require.Equal(t, connected, g.HasEdgeFromTo(nodes[from], nodes[to]))
		}
	}
}

func TestGraph_TasksInDependencyOrder(t *testing.T) {
	g := pipeline.NewTaskDAG()
	err := g.UnmarshalText([]byte(dotStr))
	require.NoError(t, err)

	u, err := url.Parse("https://chain.link/voter_turnout/USA-2020")
	require.NoError(t, err)

	answer1 := &pipeline.MedianTask{
		BaseTask: pipeline.NewBaseTask("answer1", nil, 0),
	}
	answer2 := &pipeline.BridgeTask{
		Name:     "election_winner",
		BaseTask: pipeline.NewBaseTask("answer2", nil, 1),
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times:    decimal.NewFromFloat(1.23),
		BaseTask: pipeline.NewBaseTask("ds1_multiply", answer1, 0),
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path:     []string{"one", "two"},
		BaseTask: pipeline.NewBaseTask("ds1_parse", ds1_multiply, 0),
	}
	ds1 := &pipeline.BridgeTask{
		Name:     "voter_turnout",
		BaseTask: pipeline.NewBaseTask("ds1", ds1_parse, 0),
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times:    decimal.NewFromFloat(4.56),
		BaseTask: pipeline.NewBaseTask("ds2_multiply", answer1, 0),
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path:     []string{"three", "four"},
		BaseTask: pipeline.NewBaseTask("ds2_parse", ds2_multiply, 0),
	}
	ds2 := &pipeline.HTTPTask{
		URL:         models.WebURL(*u),
		Method:      "GET",
		RequestData: pipeline.HttpRequestData{"hi": "hello"},
		BaseTask:    pipeline.NewBaseTask("ds2", ds2_parse, 0),
	}

	tasks, err := g.TasksInDependencyOrder()
	require.NoError(t, err)

	// Make sure that no task appears in the array until its output task has already appeared
	for i, task := range tasks {
		if task.OutputTask() != nil {
			require.Contains(t, tasks[:i], task.OutputTask())
		}
	}

	expected := []pipeline.Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	require.Len(t, tasks, len(expected))

	for _, task := range expected {
		require.Contains(t, tasks, task)
	}
}

func TestGraph_HasCycles(t *testing.T) {
	g := pipeline.NewTaskDAG()
	err := g.UnmarshalText([]byte(dotStr))
	require.NoError(t, err)
	require.False(t, g.HasCycles())

	g = pipeline.NewTaskDAG()
	err = dot.Unmarshal([]byte(`
        digraph {
            a [type=bridge];
            b [type=multiply times=1.23];
            a -> b -> a;
        }
    `), g)
	require.NoError(t, err)
	require.True(t, g.HasCycles())
}
