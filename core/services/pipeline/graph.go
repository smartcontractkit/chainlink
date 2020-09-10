package pipeline

import (
	"encoding/json"
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"

	"github.com/smartcontractkit/chainlink/core/store/models"
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

func (g *TaskDAG) ReverseWalkTasks(fn func(task Task) error) error {
	visited := make(map[int64]bool)
	stack := g.outputs()

	// var tasks []Task
	tasksByID := map[int64]Task{}
	for len(stack) > 0 {
		node := stack[0]
		stack = stack[1:]
		stack = append(stack, unwrapGraphNodes(g.To(node.ID()))...)
		if visited[node.ID()] {
			continue
		}

		task, err := NewTaskByType(TaskType(node.attrs["type"]))
		if err != nil {
			return err
		}

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToSliceHookFunc(","),
				func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
					switch f {
					case reflect.TypeOf(""):
						switch t {
						case reflect.TypeOf(models.WebURL{}):
							u, err := url.Parse(data.(string))
							if err != nil {
								return nil, err
							}
							return models.WebURL(*u), nil

						case reflect.TypeOf(HttpRequestData{}):
							var m map[string]interface{}
							err := json.Unmarshal([]byte(data.(string)), &m)
							return HttpRequestData(m), err

						case reflect.TypeOf(decimal.Decimal{}):
							return decimal.NewFromString(data.(string))
						}
					}
					return data, nil
				},
			),
			Result: task,
		})
		if err != nil {
			return err
		}

		err = decoder.Decode(node.attrs)
		if err != nil {
			return err
		}

		var outputTasks []Task
		for _, output := range node.outputs() {
			outputTasks = append(outputTasks, tasksByID[output.ID()])
		}
		task.SetOutputTasks(outputTasks)

		err = fn(task)
		if err != nil {
			return err
		}

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
