package pipeline_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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

	g := pipeline.NewTree()
	err := g.UnmarshalText([]byte(pipeline.DotStr))
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
	p, err := pipeline.Parse([]byte(pipeline.DotStr))
	require.NoError(t, err)

	answer1 := &pipeline.MedianTask{
		AllowedFaults: "",
	}
	answer2 := &pipeline.BridgeTask{
		Name: "election_winner",
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times: "1.23",
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path: "one,two",
	}
	ds1 := &pipeline.BridgeTask{
		Name: "voter_turnout",
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times: "4.56",
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path: "three,four",
	}
	ds2 := &pipeline.HTTPTask{
		URL:         "https://chain.link/voter_turnout/USA-2020",
		Method:      "GET",
		RequestData: `{"hi": "hello"}`,
	}

	answer1.BaseTask = pipeline.NewBaseTask(6, "answer1", []pipeline.Task{ds1_multiply, ds2_multiply}, nil, 0)
	answer2.BaseTask = pipeline.NewBaseTask(7, "answer2", nil, nil, 1)
	ds1_multiply.BaseTask = pipeline.NewBaseTask(2, "ds1_multiply", []pipeline.Task{ds1_parse}, []pipeline.Task{answer1}, 0)
	ds2_multiply.BaseTask = pipeline.NewBaseTask(5, "ds2_multiply", []pipeline.Task{ds2_parse}, []pipeline.Task{answer1}, 0)
	ds1_parse.BaseTask = pipeline.NewBaseTask(1, "ds1_parse", []pipeline.Task{ds1}, []pipeline.Task{ds1_multiply}, 0)
	ds2_parse.BaseTask = pipeline.NewBaseTask(4, "ds2_parse", []pipeline.Task{ds2}, []pipeline.Task{ds2_multiply}, 0)
	ds1.BaseTask = pipeline.NewBaseTask(0, "ds1", nil, []pipeline.Task{ds1_parse}, 0)
	ds2.BaseTask = pipeline.NewBaseTask(3, "ds2", nil, []pipeline.Task{ds2_parse}, 0)

	for i, task := range p.Tasks {
		// Make sure inputs appear before the task, and outputs don't
		for _, input := range task.Inputs() {
			require.Contains(t, p.Tasks[:i], input)
		}
		for _, output := range task.Outputs() {
			require.NotContains(t, p.Tasks[:i], output)
		}
	}

	expected := []pipeline.Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	require.Len(t, p.Tasks, len(expected))

	require.Equal(t, expected, p.Tasks)
}

func TestGraph_HasCycles(t *testing.T) {
	_, err := pipeline.Parse([]byte(pipeline.DotStr))
	require.NoError(t, err)

	_, err = pipeline.Parse([]byte(`
        digraph {
            a [type=bridge];
            b [type=multiply times=1.23];
            a -> b -> a;
        }
    `))
	require.Error(t, err)
}
