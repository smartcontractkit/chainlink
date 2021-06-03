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
type Tree struct {
	*simple.DirectedGraph
}

func NewTree() *Tree {
	return &Tree{DirectedGraph: simple.NewDirectedGraph()}
}

func (g *Tree) NewNode() graph.Node {
	return &TreeNode{Node: g.DirectedGraph.NewNode()}
}

func (g *Tree) UnmarshalText(bs []byte) (err error) {
	if g.DirectedGraph == nil {
		g.DirectedGraph = simple.NewDirectedGraph()
	}
	bs = append([]byte("digraph {\n"), bs...)
	bs = append(bs, []byte("\n}")...)
	err = dot.Unmarshal(bs, g)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal DOT into a pipeline.Tree")
	}
	return nil
}

//

type TreeNode struct {
	graph.Node
	dotID string
	attrs map[string]string
}

func NewTreeNode(n graph.Node, dotID string, attrs map[string]string) *TreeNode {
	return &TreeNode{
		Node:  n,
		attrs: attrs,
		dotID: dotID,
	}
}

func (n *TreeNode) DOTID() string {
	return n.dotID
}

func (n *TreeNode) SetDOTID(id string) {
	n.dotID = id
}

func (n *TreeNode) String() string {
	return n.dotID
}

var bracketQuotedAttrRegexp = regexp.MustCompile(`^\s*<([^<>]+)>\s*$`)

func (n *TreeNode) SetAttribute(attr encoding.Attribute) error {
	if n.attrs == nil {
		n.attrs = make(map[string]string)
	}

	// Strings quoted in angle brackets (supported natively by DOT) should
	// have those brackets removed before decoding to task parameter types
	sanitized := bracketQuotedAttrRegexp.ReplaceAllString(attr.Value, "$1")

	n.attrs[attr.Key] = sanitized
	return nil
}

func (n *TreeNode) Attributes() []encoding.Attribute {
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

//

type Pipeline struct {
	Tasks  map[int64]Task
	tree   *Tree
	Source string
}

func (p *Pipeline) UnmarshalText(bs []byte) (err error) {
	parsed, err := Parse(bs)
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

func (p *Pipeline) TopoSort() []Task {
	result := make([]Task, len(p.Tasks))

	for id, task := range p.Tasks {
		result[id] = task
	}
	return result
}

func Parse(bs []byte) (*Pipeline, error) {
	g := NewTree()
	err := g.UnmarshalText(bs)

	if err != nil {
		return nil, err
	}

	p := &Pipeline{
		tree:   g,
		Tasks:  make(map[int64]Task, g.Nodes().Len()),
		Source: string(bs),
	}

	// toposort all the nodes: dependencies ordered before outputs. This also does cycle checking for us.
	nodes, err := topo.SortStabilized(g, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to topologically sort the graph, cycle detected")
	}

	for id, node := range nodes {
		// use the new ordering as the id so that we can easily reproduce the original toposort
		id := int64(id)

		node, is := node.(*TreeNode)
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

		p.Tasks[id] = task
	}

	// re-link the edges
	for edges := g.Edges(); edges.Next(); {
		edge := edges.Edge()

		from := p.Tasks[edge.From().ID()]
		to := p.Tasks[edge.To().ID()]

		from.Base().outputs = append(from.Base().outputs, to)
		to.Base().inputs = append(to.Base().inputs, from)
	}

	return p, nil
}
