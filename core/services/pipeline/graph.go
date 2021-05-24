package pipeline

import (
	"regexp"
	"sort"
	"time"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

// TaskDAG fulfills the graph.DirectedGraph interface, which makes it possible
// for us to `dot.Unmarshal(...)` a DOT string directly into it.  Once unmarshalled,
// calling `TaskDAG#TasksInDependencyOrder()` will return the unmarshaled tasks.
// NOTE: We only permit one child
type TaskDAG struct {
	*simple.DirectedGraph
	DOTSource string
}

func NewTaskDAG() *TaskDAG {
	return &TaskDAG{DirectedGraph: simple.NewDirectedGraph()}
}

func (g *TaskDAG) NewNode() graph.Node {
	return &TaskDAGNode{Node: g.DirectedGraph.NewNode(), g: g}
}

func (g *TaskDAG) UnmarshalText(bs []byte) (err error) {
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

		numPredecessors := g.To(node.ID()).Len()
		task, err := UnmarshalTaskFromMap(TaskType(node.attrs["type"]), node.attrs, node.dotID, nil, nil, nil, numPredecessors)
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

func (g TaskDAG) MinTimeout() (time.Duration, bool, error) {
	var minTimeout time.Duration = 1<<63 - 1
	var aTimeoutSet bool
	tasks, err := g.TasksInDependencyOrder()
	if err != nil {
		return minTimeout, aTimeoutSet, err
	}
	for _, t := range tasks {
		if timeout, set := t.TaskTimeout(); set && timeout < minTimeout {
			minTimeout = timeout
			aTimeoutSet = true
		}
	}
	return minTimeout, aTimeoutSet, nil
}

func (g TaskDAG) outputs() []*TaskDAGNode {
	var outputs []*TaskDAGNode
	iter := g.Nodes()
	for iter.Next() {
		node, is := iter.Node().(*TaskDAGNode)
		if !is {
			panic("this is impossible but we must appease go staticcheck")
		}
		if g.From(node.ID()) == graph.Empty {
			outputs = append(outputs, node)
		}
	}
	return outputs
}

type TaskDAGNode struct {
	graph.Node
	g     *TaskDAG
	dotID string
	attrs map[string]string
}

func NewTaskDAGNode(n graph.Node, dotID string, attrs map[string]string) *TaskDAGNode {
	return &TaskDAGNode{
		Node:  n,
		attrs: attrs,
		dotID: dotID,
	}
}

func (n *TaskDAGNode) SetDAG(g *TaskDAG) {
	n.g = g
}

func (n *TaskDAGNode) DOTID() string {
	return n.dotID
}

func (n *TaskDAGNode) SetDOTID(id string) {
	n.dotID = id
}

func (n *TaskDAGNode) String() string {
	return n.dotID
}

var bracketQuotedAttrRegexp = regexp.MustCompile(`^\s*<([^<>]+)>\s*$`)

func (n *TaskDAGNode) SetAttribute(attr encoding.Attribute) error {
	if n.attrs == nil {
		n.attrs = make(map[string]string)
	}

	// Strings quoted in angle brackets (supported natively by DOT) should
	// have those brackets removed before decoding to task parameter types
	sanitized := bracketQuotedAttrRegexp.ReplaceAllString(attr.Value, "$1")

	n.attrs[attr.Key] = sanitized
	return nil
}

func (n *TaskDAGNode) Attributes() []encoding.Attribute {
	var r []encoding.Attribute
	for k, v := range n.attrs {
		r = append(r, encoding.Attribute{Key: k, Value: v})
	}
	// Ensure the slice returned is deterministic.
	sort.Slice(r, func(i, j int) bool {
		return r[i].Key < r[j].Key
	})
	return r
}

func (n *TaskDAGNode) outputs() []*TaskDAGNode {
	var nodes []*TaskDAGNode
	ns := n.g.From(n.ID())
	for ns.Next() {
		nodes = append(nodes, ns.Node().(*TaskDAGNode))
	}
	return nodes
}

func unwrapGraphNodes(nodes graph.Nodes) []*TaskDAGNode {
	var out []*TaskDAGNode
	for nodes.Next() {
		out = append(out, nodes.Node().(*TaskDAGNode))
	}
	return out
}
