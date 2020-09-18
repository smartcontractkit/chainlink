package pipeline

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gopkg.in/guregu/null.v4"
)

// TaskDAG fulfills the graph.DirectedGraph interface, which makes it possible
// for us to `dot.Unmarshal(...)` a DOT string directly into it.  Once unmarshalled,
// calling `TaskDAG#Tasks()` will return
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

func (g *TaskDAG) UnmarshalText(bs []byte) error {
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	err := dot.Unmarshal(bs, g)
	if err != nil {
		return err
	}
	g.DOTSource = string(bs)
	return nil
}

func (g *TaskDAG) HasCycles() bool {
	return len(topo.DirectedCyclesIn(g)) > 0
}

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

		task, err := UnmarshalTask(TaskType(node.attrs["type"]), node.attrs, nil, nil)
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

func (g TaskDAG) ToPipelineSpec() (Spec, error) {
	tasks, err := g.TasksInDependencyOrder()
	if err != nil {
		return Spec{}, err
	}

	// Convert the task DAG into TaskSpec DB rows.  We walk the TaskDAG backwards,
	// from final outputs to inputs, to ensure that each task's successor is
	// already in the `taskSpecIDs` map.
	taskSpecs := []TaskSpec{}
	taskSpecIDs := make(map[Task]int32)
	for _, task := range tasks {
		var successorID null.Int
		if task.OutputTask() != nil {
			successor := task.OutputTask()
			successorID = null.IntFrom(int64(taskSpecIDs[successor]))
		}

		taskSpec := TaskSpec{
			Type:        task.Type(),
			JSON:        JSONSerializable{Value: task},
			SuccessorID: successorID,
		}

		taskSpecIDs[task] = taskSpec.ID
		taskSpecs = append(taskSpecs, taskSpec)
	}

	return Spec{
		DotDagSource: g.DOTSource,
		TaskSpecs:    taskSpecs,
	}, nil
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
