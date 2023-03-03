// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
)

// Represents the dependency graph for a given Node.
//
// This graph includes the dependency relationships between all nodes in the
// node (construct) sub-tree who's root is this Node.
//
// Note that this means that lonely nodes (no dependencies and no dependants) are also included in this graph as
// childless children of the root node of the graph.
//
// The graph does not include cross-scope dependencies. That is, if a child on the current scope depends on a node
// from a different scope, that relationship is not represented in this graph.
type DependencyGraph interface {
	// Returns the root of the graph.
	//
	// Note that this vertex will always have `null` as its `.value` since it is an artifical root
	// that binds all the connected spaces of the graph.
	Root() DependencyVertex
	// See: Vertex.topology()
	//
	Topology() *[]constructs.IConstruct
}

// The jsii proxy struct for DependencyGraph
type jsiiProxy_DependencyGraph struct {
	_ byte // padding
}

func (j *jsiiProxy_DependencyGraph) Root() DependencyVertex {
	var returns DependencyVertex
	_jsii_.Get(
		j,
		"root",
		&returns,
	)
	return returns
}


func NewDependencyGraph(node constructs.Node) DependencyGraph {
	_init_.Initialize()

	if err := validateNewDependencyGraphParameters(node); err != nil {
		panic(err)
	}
	j := jsiiProxy_DependencyGraph{}

	_jsii_.Create(
		"cdk8s.DependencyGraph",
		[]interface{}{node},
		&j,
	)

	return &j
}

func NewDependencyGraph_Override(d DependencyGraph, node constructs.Node) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.DependencyGraph",
		[]interface{}{node},
		d,
	)
}

func (d *jsiiProxy_DependencyGraph) Topology() *[]constructs.IConstruct {
	var returns *[]constructs.IConstruct

	_jsii_.Invoke(
		d,
		"topology",
		nil, // no parameters
		&returns,
	)

	return returns
}

