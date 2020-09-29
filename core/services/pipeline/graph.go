package pipeline

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// TaskDAG fulfills the graph.DirectedGraph interface, which makes it possible
// for us to `dot.Unmarshal(...)` a DOT string directly into it.  Once unmarshalled,
// calling `TaskDAG#TasksInDependencyOrder()` will return the unmarshaled tasks.
type TaskDAG struct {
	*simple.DirectedGraph
	DOTSource string
}

func NewTaskDAG() *TaskDAG {
	return &TaskDAG{DirectedGraph: simple.NewDirectedGraph()}
}

func (g *TaskDAG) NewNode() graph.Node {
	return &taskDAGNode{Node: g.DirectedGraph.NewNode(), g: g}
}

func (g *TaskDAG) UnmarshalText(bs []byte) (err error) {
	defer utils.LogIfError(&err, "TaskDAG#UnmarshalText: %+v")
	defer utils.WrapIfError(&err, "MedianTask errored")
	if g.DirectedGraph == nil {
		g.DirectedGraph = simple.NewDirectedGraph()
	}
	g.DOTSource = string(bs)
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	err = dot.Unmarshal(bs, g)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal DOT into a pipeline.TaskDAG")
	}
	return nil
}

func (g *TaskDAG) HasCycles() bool {
	return len(topo.DirectedCyclesIn(g)) > 0
}

// Returns a slice of Tasks starting at the outputs of the DAG and ending at
// the inputs.  As you iterate through this slice, you can expect that any individual
// Task's outputs will already have been traversed.
func (g TaskDAG) TasksInDependencyOrder() ([]Task, error) {
	visited := make(map[int64]bool)
	stack := g.outputs()

	tasksByID := map[int64]Task{}
	var tasks []Task
	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]
		stack = append(stack, unwrapGraphNodes(g.To(node.ID()))...)
		if visited[node.ID()] {
			continue
		}

		task, err := UnmarshalTaskFromMap(TaskType(node.attrs["type"]), node.attrs, node.dotID, nil, nil)
		if err != nil {
			return nil, err
		}

		var outputTasks []Task
		for _, output := range node.outputs() {
			outputTasks = append(outputTasks, tasksByID[output.ID()])
		}
		if len(outputTasks) > 1 {
			return nil, errors.New("task has > 1 output task")
		} else if len(outputTasks) == 1 {
			task.SetOutputTask(outputTasks[0])
		}

		tasks = append(tasks, task)

		tasksByID[node.ID()] = task
		visited[node.ID()] = true
	}
	return tasks, nil
}

func (g *TaskDAG) inputs() []*taskDAGNode {
	var inputs []*taskDAGNode
	iter := g.Nodes()
	for iter.Next() {
		node := iter.Node().(*taskDAGNode)
		if g.To(node.ID()) == graph.Empty {
			inputs = append(inputs, node)
		}
	}
	return inputs
}

func (g *TaskDAG) outputs() []*taskDAGNode {
	var outputs []*taskDAGNode
	iter := g.Nodes()
	for iter.Next() {
		node := iter.Node().(*taskDAGNode)
		if g.From(node.ID()) == graph.Empty {
			outputs = append(outputs, node)
		}
	}
	return outputs
}

type taskDAGNode struct {
	graph.Node
	g     *TaskDAG
	dotID string
	attrs map[string]string
}

func (n *taskDAGNode) DOTID() string {
	return n.dotID
}

func (n *taskDAGNode) SetDOTID(id string) {
	n.dotID = id
}

func (n *taskDAGNode) String() string {
	return n.dotID
}

func (n *taskDAGNode) SetAttribute(attr encoding.Attribute) error {
	if n.attrs == nil {
		n.attrs = make(map[string]string)
	}
	n.attrs[attr.Key] = attr.Value
	return nil
}

func (n *taskDAGNode) inputs() []*taskDAGNode {
	var nodes []*taskDAGNode
	ns := n.g.To(n.ID())
	for ns.Next() {
		nodes = append(nodes, ns.Node().(*taskDAGNode))
	}
	return nodes
}

func (n *taskDAGNode) outputs() []*taskDAGNode {
	var nodes []*taskDAGNode
	ns := n.g.From(n.ID())
	for ns.Next() {
		nodes = append(nodes, ns.Node().(*taskDAGNode))
	}
	return nodes
}

func unwrapGraphNodes(nodes graph.Nodes) []*taskDAGNode {
	var out []*taskDAGNode
	for nodes.Next() {
		out = append(out, nodes.Node().(*taskDAGNode))
	}
	return out
}
