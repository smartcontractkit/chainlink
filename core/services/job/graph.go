package job

import (
	"github.com/mitchellh/mapstructure"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type TaskDAG struct {
	*simple.DirectedGraph
}

func NewTaskDAG() *TaskDAG {
	return &TaskDAG{DirectedGraph: simple.NewDirectedGraph()}
}

func (g *TaskDAG) NewNode() graph.Node {
	return &taskDAGNode{Node: g.DirectedGraph.NewNode(), g: g}
}

func (g *TaskDAG) UnmarshalText(bs []byte) error {
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	return dot.Unmarshal(bs, g)
}

func (g *TaskDAG) HasCycles() bool {
	return len(topo.DirectedCyclesIn(g)) == 0
}

func (g *TaskDAG) Tasks() ([]Task, error) {
	visited := make(map[int64]bool)
	stack := g.inputs()

	var tasks []Task
	tasksByID := map[int64]Task{}
	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]
		stack = append(stack, unwrapGraphNodes(g.From(node.ID()))...)
		if visited[node.ID()] {
			continue
		}

		task, err := NewTaskByType(TaskType(node.attrs["type"]))
		if err != nil {
			return nil, err
		}

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.StringToSliceHookFunc(","),
			Result:     task,
		})
		if err != nil {
			return nil, err
		}

		err = decoder.Decode(node.attrs)
		if err != nil {
			return nil, err
		}

		var inputTasks []Task
		for _, input := range node.inputs() {
			inputTasks = append(inputTasks, tasksByID[input.ID()])
		}
		task.SetInputTasks(inputTasks)

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

func unwrapGraphNodes(nodes graph.Nodes) []*taskDAGNode {
	var out []*taskDAGNode
	for nodes.Next() {
		out = append(out, nodes.Node().(*taskDAGNode))
	}
	return out
}
