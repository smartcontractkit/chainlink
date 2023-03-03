// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

type Lazy interface {
	Produce() interface{}
}

// The jsii proxy struct for Lazy
type jsiiProxy_Lazy struct {
	_ byte // padding
}

func Lazy_Any(producer IAnyProducer) interface{} {
	_init_.Initialize()

	if err := validateLazy_AnyParameters(producer); err != nil {
		panic(err)
	}
	var returns interface{}

	_jsii_.StaticInvoke(
		"cdk8s.Lazy",
		"any",
		[]interface{}{producer},
		&returns,
	)

	return returns
}

func (l *jsiiProxy_Lazy) Produce() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		l,
		"produce",
		nil, // no parameters
		&returns,
	)

	return returns
}

