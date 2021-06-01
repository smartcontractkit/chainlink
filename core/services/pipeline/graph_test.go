package pipeline_test

import (
	"fmt"
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
	t.Skip()
	p, err := pipeline.Parse([]byte(pipeline.DotStr))
	require.NoError(t, err)

	answer1 := &pipeline.MedianTask{
		BaseTask:      pipeline.NewBaseTask(6, "answer1", nil, 0),
		AllowedFaults: "",
	}
	answer2 := &pipeline.BridgeTask{
		Name:     "election_winner",
		BaseTask: pipeline.NewBaseTask(7, "answer2", nil, 1),
	}
	ds1_multiply := &pipeline.MultiplyTask{
		Times:    "1.23",
		BaseTask: pipeline.NewBaseTask(2, "ds1_multiply", answer1, 0),
	}
	ds1_parse := &pipeline.JSONParseTask{
		Path:     "one,two",
		BaseTask: pipeline.NewBaseTask(1, "ds1_parse", ds1_multiply, 0),
	}
	ds1 := &pipeline.BridgeTask{
		Name:     "voter_turnout",
		BaseTask: pipeline.NewBaseTask(0, "ds1", ds1_parse, 0),
	}
	ds2_multiply := &pipeline.MultiplyTask{
		Times:    "4.56",
		BaseTask: pipeline.NewBaseTask(5, "ds2_multiply", answer1, 0),
	}
	ds2_parse := &pipeline.JSONParseTask{
		Path:     "three,four",
		BaseTask: pipeline.NewBaseTask(4, "ds2_parse", ds2_multiply, 0),
	}
	ds2 := &pipeline.HTTPTask{
		URL:         "https://chain.link/voter_turnout/USA-2020",
		Method:      "GET",
		RequestData: `{"hi": "hello"}`,
		BaseTask:    pipeline.NewBaseTask(3, "ds2", ds2_parse, 0),
	}

	tasks := p.TopoSort()

	for i, task := range tasks {
		// Make sure inputs appear before the task, and outputs don't
		for _, input := range task.Inputs() {
			require.Contains(t, tasks[:i], input)
		}
		for _, output := range task.Outputs() {
			require.NotContains(t, tasks[:i], output)
		}
	}

	for _, t := range tasks {
		fmt.Printf("%v %v", t.ID(), t.DotID())
	}

	expected := []pipeline.Task{ds1, ds1_parse, ds1_multiply, ds2, ds2_parse, ds2_multiply, answer1, answer2}
	require.Len(t, tasks, len(expected))

	require.Equal(t, tasks, expected)
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
