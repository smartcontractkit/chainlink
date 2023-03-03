// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
)

// Represents a vertex in the graph.
//
// The value of each vertex is an `IConstruct` that is accessible via the `.value` getter.
type DependencyVertex interface {
	// Returns the parents of the vertex (i.e dependants).
	Inbound() *[]DependencyVertex
	// Returns the children of the vertex (i.e dependencies).
	Outbound() *[]DependencyVertex
	// Returns the IConstruct this graph vertex represents.
	//
	// `null` in case this is the root of the graph.
	Value() constructs.IConstruct
	// Adds a vertex as a dependency of the current node.
	//
	// Also updates the parents of `dep`, so that it contains this node as a parent.
	//
	// This operation will fail in case it creates a cycle in the graph.
	AddChild(dep DependencyVertex)
	// Returns a topologically sorted array of the constructs in the sub-graph.
	Topology() *[]constructs.IConstruct
}

// The jsii proxy struct for DependencyVertex
type jsiiProxy_DependencyVertex struct {
	_ byte // padding
}

func (j *jsiiProxy_DependencyVertex) Inbound() *[]DependencyVertex {
	var returns *[]DependencyVertex
	_jsii_.Get(
		j,
		"inbound",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DependencyVertex) Outbound() *[]DependencyVertex {
	var returns *[]DependencyVertex
	_jsii_.Get(
		j,
		"outbound",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DependencyVertex) Value() constructs.IConstruct {
	var returns constructs.IConstruct
	_jsii_.Get(
		j,
		"value",
		&returns,
	)
	return returns
}


func NewDependencyVertex(value constructs.IConstruct) DependencyVertex {
	_init_.Initialize()

	j := jsiiProxy_DependencyVertex{}

	_jsii_.Create(
		"cdk8s.DependencyVertex",
		[]interface{}{value},
		&j,
	)

	return &j
}

func NewDependencyVertex_Override(d DependencyVertex, value constructs.IConstruct) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.DependencyVertex",
		[]interface{}{value},
		d,
	)
}

func (d *jsiiProxy_DependencyVertex) AddChild(dep DependencyVertex) {
	if err := d.validateAddChildParameters(dep); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"addChild",
		[]interface{}{dep},
	)
}

func (d *jsiiProxy_DependencyVertex) Topology() *[]constructs.IConstruct {
	var returns *[]constructs.IConstruct

	_jsii_.Invoke(
		d,
		"topology",
		nil, // no parameters
		&returns,
	)

	return returns
}

