// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/internal"
)

// Reads a YAML manifest from a file or a URL and defines all resources as API objects within the defined scope.
//
// The names (`metadata.name`) of imported resources will be preserved as-is
// from the manifest.
type Include interface {
	constructs.Construct
	// Returns all the included API objects.
	ApiObjects() *[]ApiObject
	// The tree node.
	Node() constructs.Node
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for Include
type jsiiProxy_Include struct {
	internal.Type__constructsConstruct
}

func (j *jsiiProxy_Include) ApiObjects() *[]ApiObject {
	var returns *[]ApiObject
	_jsii_.Get(
		j,
		"apiObjects",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Include) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


func NewInclude(scope constructs.Construct, id *string, props *IncludeProps) Include {
	_init_.Initialize()

	if err := validateNewIncludeParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_Include{}

	_jsii_.Create(
		"cdk8s.Include",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewInclude_Override(i Include, scope constructs.Construct, id *string, props *IncludeProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.Include",
		[]interface{}{scope, id, props},
		i,
	)
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func Include_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateInclude_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk8s.Include",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_Include) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

