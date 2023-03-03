// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/internal"
)

type Chart interface {
	constructs.Construct
	// Labels applied to all resources in this chart.
	//
	// This is an immutable copy.
	Labels() *map[string]*string
	// The default namespace for all objects in this chart.
	Namespace() *string
	// The tree node.
	Node() constructs.Node
	// Create a dependency between this Chart and other constructs.
	//
	// These can be other ApiObjects, Charts, or custom.
	AddDependency(dependencies ...constructs.IConstruct)
	// Generates a app-unique name for an object given it's construct node path.
	//
	// Different resource types may have different constraints on names
	// (`metadata.name`). The previous version of the name generator was
	// compatible with DNS_SUBDOMAIN but not with DNS_LABEL.
	//
	// For example, `Deployment` names must comply with DNS_SUBDOMAIN while
	// `Service` names must comply with DNS_LABEL.
	//
	// Since there is no formal specification for this, the default name
	// generation scheme for kubernetes objects in cdk8s was changed to DNS_LABEL,
	// since it’s the common denominator for all kubernetes resources
	// (supposedly).
	//
	// You can override this method if you wish to customize object names at the
	// chart level.
	GenerateObjectName(apiObject ApiObject) *string
	// Renders this chart to a set of Kubernetes JSON resources.
	//
	// Returns: array of resource manifests.
	ToJson() *[]interface{}
	// Returns a string representation of this construct.
	ToString() *string
}

// The jsii proxy struct for Chart
type jsiiProxy_Chart struct {
	internal.Type__constructsConstruct
}

func (j *jsiiProxy_Chart) Labels() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Chart) Namespace() *string {
	var returns *string
	_jsii_.Get(
		j,
		"namespace",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Chart) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}


func NewChart(scope constructs.Construct, id *string, props *ChartProps) Chart {
	_init_.Initialize()

	if err := validateNewChartParameters(scope, id, props); err != nil {
		panic(err)
	}
	j := jsiiProxy_Chart{}

	_jsii_.Create(
		"cdk8s.Chart",
		[]interface{}{scope, id, props},
		&j,
	)

	return &j
}

func NewChart_Override(c Chart, scope constructs.Construct, id *string, props *ChartProps) {
	_init_.Initialize()

	_jsii_.Create(
		"cdk8s.Chart",
		[]interface{}{scope, id, props},
		c,
	)
}

// Return whether the given object is a Chart.
//
// We do attribute detection since we can't reliably use 'instanceof'.
func Chart_IsChart(x interface{}) *bool {
	_init_.Initialize()

	if err := validateChart_IsChartParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk8s.Chart",
		"isChart",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Checks if `x` is a construct.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
// Deprecated: use `x instanceof Construct` instead.
func Chart_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateChart_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"cdk8s.Chart",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Finds the chart in which a node is defined.
func Chart_Of(c constructs.IConstruct) Chart {
	_init_.Initialize()

	if err := validateChart_OfParameters(c); err != nil {
		panic(err)
	}
	var returns Chart

	_jsii_.StaticInvoke(
		"cdk8s.Chart",
		"of",
		[]interface{}{c},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Chart) AddDependency(dependencies ...constructs.IConstruct) {
	args := []interface{}{}
	for _, a := range dependencies {
		args = append(args, a)
	}

	_jsii_.InvokeVoid(
		c,
		"addDependency",
		args,
	)
}

func (c *jsiiProxy_Chart) GenerateObjectName(apiObject ApiObject) *string {
	if err := c.validateGenerateObjectNameParameters(apiObject); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		c,
		"generateObjectName",
		[]interface{}{apiObject},
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Chart) ToJson() *[]interface{} {
	var returns *[]interface{}

	_jsii_.Invoke(
		c,
		"toJson",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (c *jsiiProxy_Chart) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		c,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

