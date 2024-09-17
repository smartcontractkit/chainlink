package pipeline

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
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

func (g *Graph) NewEdge(from, to graph.Node) graph.Edge {
	return &GraphEdge{Edge: g.DirectedGraph.NewEdge(from, to)}
}

func (g *Graph) UnmarshalText(bs []byte) (err error) {
	if g.DirectedGraph == nil {
		g.DirectedGraph = simple.NewDirectedGraph()
	}
	defer func() {
		if rerr := recover(); rerr != nil {
			err = fmt.Errorf("could not unmarshal DOT into a pipeline.Graph: %v", rerr)
		}
	}()
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	err = dot.Unmarshal(bs, g)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal DOT into a pipeline.Graph")
	}
	g.AddImplicitDependenciesAsEdges()
	return nil
}

// Looks at node attributes and searches for implicit dependencies on other nodes
// expressed as attribute values. Adds those dependencies as implicit edges in the graph.
func (g *Graph) AddImplicitDependenciesAsEdges() {
	for nodesIter := g.Nodes(); nodesIter.Next(); {
		graphNode := nodesIter.Node().(*GraphNode)

		params := make(map[string]bool)
		// Walk through all attributes and find all params which this node depends on
		for _, attr := range graphNode.Attributes() {
			for _, item := range variableRegexp.FindAll([]byte(attr.Value), -1) {
				expr := strings.TrimSpace(string(item[2 : len(item)-1]))
				param := strings.Split(expr, ".")[0]
				params[param] = true
			}
		}
		// Iterate through all nodes and add a new edge if node belongs to params set, and there already isn't an edge.
		for nodesIter2 := g.Nodes(); nodesIter2.Next(); {
			gn := nodesIter2.Node().(*GraphNode)
			if params[gn.DOTID()] {
				// If these are distinct nodes with no existing edge between them, then add an implicit edge.
				if gn.ID() != graphNode.ID() && !g.HasEdgeFromTo(gn.ID(), graphNode.ID()) {
					edge := g.NewEdge(gn, graphNode).(*GraphEdge)
					// Setting isImplicit indicates that this edge wasn't specified via the TOML spec,
					// but rather added automatically here.
					// This distinction is needed, as we don't want to propagate results of a task to its dependent
					// tasks along implicit edge, as some tasks can't handle unexpected inputs from implicit edges.
					edge.SetIsImplicit(true)
					g.SetEdge(edge)
				}
			}
		}
	}
}

// Indicates whether there's an implicit edge from uid -> vid.
// Implicit edged are ones that weren't added via the TOML spec, but via the pipeline parsing code
func (g *Graph) IsImplicitEdge(uid, vid int64) bool {
	edge := g.Edge(uid, vid).(*GraphEdge)
	if edge == nil {
		return false
	}
	return edge.IsImplicit()
}

type GraphEdge struct {
	graph.Edge

	// Indicates that this edge was implicitly added by the pipeline parser, and not via the TOML specs.
	isImplicit bool
}

func (e *GraphEdge) IsImplicit() bool {
	return e.isImplicit
}

func (e *GraphEdge) SetIsImplicit(isImplicit bool) {
	e.isImplicit = isImplicit
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
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("empty pipeline")
	}
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

	resultIdxs := make(map[int32]struct{})

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

		if task.OutputIndex() > 0 {
			_, exists := resultIdxs[task.OutputIndex()]
			if exists {
				return nil, errors.New("duplicate sorting indexes detected")
			}

			resultIdxs[task.OutputIndex()] = struct{}{}
		}

		// re-link the edges
		for inputs := g.To(node.ID()); inputs.Next(); {
			isImplicitEdge := g.IsImplicitEdge(inputs.Node().ID(), node.ID())
			from := p.Tasks[ids[inputs.Node().ID()]]

			from.Base().outputs = append(from.Base().outputs, task)
			task.Base().inputs = append(task.Base().inputs, TaskDependency{!isImplicitEdge, from})
		}

		// This is subtle: g.To doesn't return nodes in deterministic order, which would occasionally swap the order
		// of inputs, therefore we manually sort. We don't need to sort outputs the same way because these appends happen
		// in p.Task order, which is deterministic via topo.SortStable.
		sort.Slice(task.Base().inputs, func(i, j int) bool {
			return task.Base().inputs[i].InputTask.ID() < task.Base().inputs[j].InputTask.ID()
		})

		p.Tasks = append(p.Tasks, task)
		ids[node.ID()] = id
	}

	return p, nil
}
