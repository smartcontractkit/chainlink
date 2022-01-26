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

// tree fulfills the graph.DirectedGraph interface, which makes it possible
// for us to `dot.Unmarshal(...)` a DOT string directly into it.
type Graph struct {
	*simple.DirectedGraph
}

func NewGraph() *Graph {
	return &Graph{DirectedGraph: simple.NewDirectedGraph()}
}

func (g *Graph) NewNode() graph.Node {
	return &GraphNode{Node: g.DirectedGraph.NewNode()}
}

func (g *Graph) UnmarshalText(bs []byte) (err error) {
	if g.DirectedGraph == nil {
		g.DirectedGraph = simple.NewDirectedGraph()
	}
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	err = dot.Unmarshal(bs, g)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal DOT into a pipeline.Graph")
	}
	return nil
}

type GraphNode struct {
	graph.Node
	dotID string
	attrs map[string]string
}

func (n *GraphNode) DOTID() string {
	return n.dotID
}

func (n *GraphNode) SetDOTID(id string) {
	n.dotID = id
}

func (n *GraphNode) String() string {
	return n.dotID
}

var bracketQuotedAttrRegexp = regexp.MustCompile(`\A\s*<([^<>]+)>\s*\z`)

func (n *GraphNode) SetAttribute(attr encoding.Attribute) error {
	if n.attrs == nil {
		n.attrs = make(map[string]string)
	}

	// Strings quoted in angle brackets (supported natively by DOT) should
	// have those brackets removed before decoding to task parameter types
	sanitized := bracketQuotedAttrRegexp.ReplaceAllString(attr.Value, "$1")

	n.attrs[attr.Key] = sanitized
	return nil
}

func (n *GraphNode) Attributes() []encoding.Attribute {
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

type Pipeline struct {
	Tasks  []Task
	tree   *Graph
	Source string
}

func (p *Pipeline) UnmarshalText(bs []byte) (err error) {
	parsed, err := Parse(string(bs))
	if err != nil {
		return err
	}
	*p = *parsed
	return nil
}

func (p *Pipeline) MinTimeout() (time.Duration, bool, error) {
	var minTimeout time.Duration = 1<<63 - 1
	var aTimeoutSet bool
	for _, t := range p.Tasks {
		if timeout, set := t.TaskTimeout(); set && timeout < minTimeout {
			minTimeout = timeout
			aTimeoutSet = true
		}
	}
	return minTimeout, aTimeoutSet, nil
}

func (p *Pipeline) RequiresPreInsert() bool {
	for _, task := range p.Tasks {
		switch task.Type() {
		case TaskTypeBridge:
			if task.(*BridgeTask).Async == "true" {
				return true
			}
		case TaskTypeETHTx:
			// we want to pre-insert pipeline_task_runs always
			return true
		default:
		}
	}
	return false
}

func (p *Pipeline) ByDotID(id string) Task {
	for _, task := range p.Tasks {
		if task.DotID() == id {
			return task
		}
	}
	return nil
}

func Parse(text string) (*Pipeline, error) {
	g := NewGraph()
	err := g.UnmarshalText([]byte(text))

	if err != nil {
		return nil, err
	}

	p := &Pipeline{
		tree:   g,
		Tasks:  make([]Task, 0, g.Nodes().Len()),
		Source: text,
	}

	// toposort all the nodes: dependencies ordered before outputs. This also does cycle checking for us.
	nodes, err := topo.SortStabilized(g, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to topologically sort the graph, cycle detected")
	}

	// we need a temporary mapping of graph.IDs to positional ids after toposort
	ids := make(map[int64]int)

	// use the new ordering as the id so that we can easily reproduce the original toposort
	for id, node := range nodes {
		node, is := node.(*GraphNode)
		if !is {
			panic("unreachable")
		}

		if node.dotID == InputTaskKey {
			return nil, errors.Errorf("'%v' is a reserved keyword that cannot be used as a task's name", InputTaskKey)
		}

		task, err := UnmarshalTaskFromMap(TaskType(node.attrs["type"]), node.attrs, id, node.dotID)
		if err != nil {
			return nil, err
		}

		// re-link the edges
		for inputs := g.To(node.ID()); inputs.Next(); {
			from := p.Tasks[ids[inputs.Node().ID()]]

			from.Base().outputs = append(from.Base().outputs, task)
			task.Base().inputs = append(task.Base().inputs, from)
		}

		// This is subtle: g.To doesn't return nodes in deterministic order, which would occasionally swap the order
		// of inputs, therefore we manually sort. We don't need to sort outputs the same way because these appends happen
		// in p.Task order, which is deterministic via topo.SortStable.
		sort.Slice(task.Base().inputs, func(i, j int) bool {
			return task.Base().inputs[i].ID() < task.Base().inputs[j].ID()
		})

		p.Tasks = append(p.Tasks, task)
		ids[node.ID()] = id
	}

	return p, nil
}
