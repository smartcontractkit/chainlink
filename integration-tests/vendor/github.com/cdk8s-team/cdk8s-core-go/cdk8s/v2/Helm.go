// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
)

// Represents a Helm deployment.
//
// Use this construct to import an existing Helm chart and incorporate it into your constructs.
type Helm interface {
	Include
	// Returns all the included API objects.
	ApiObjects() *[]ApiObject
	// The tree node.
	Node() constructs.Node
	// The helm release name.
	ReleaseName() *string
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for Helm
type jsiiProxy_Helm struct {
	jsiiProxy_Include
}

func (j *jsiiProxy_Helm) ApiObjects() *[]ApiObject {
	var returns *[]ApiObject
	_jsii_.Get(
		j,
		"apiObjects",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Helm) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Helm) ReleaseName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"releaseName",
		&returns,
	)
	return returns
}


func NewHelm(scope constructs.Construct, id *string, props *HelmProps) Helm {
	_init_.Initialize()

	if err := validateNewHelmParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_Helm{}

	_jsii_.Create(
		"cdk8s.Helm",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewHelm_Override(h Helm, scope constructs.Construct, id *string, props *HelmProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.Helm",
		[]interface{}{scope, id, props},
		h,
	)
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func Helm_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateHelm_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk8s.Helm",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (h *jsiiProxy_Helm) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		h,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

