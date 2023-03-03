// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2/jsii"
)

// Testing utilities for cdk8s applications.
type Testing interface {
}

// The jsii proxy struct for Testing
type jsiiProxy_Testing struct {
	_ byte // padding
}

// Returns an app for testing with the following properties: - Output directory is a temp dir.
func Testing_App(props *AppProps) App {
	_init_.Initialize()

	if err := validateTesting_AppParameters(props); err != nil {
		panic(err)
	}
	var returns App

	_jsii_.StaticInvoke(
		"cdk8s.Testing",
		"app",
		[]interface{}{props},
		&returns,
	)

	return returns
}

// Returns: a Chart that can be used for tests.
func Testing_Chart() Chart {
	_init_.Initialize()

	var returns Chart

	_jsii_.StaticInvoke(
		"cdk8s.Testing",
		"chart",
		nil, // no parameters
		&returns,
	)

	return returns
}

// Returns the Kubernetes manifest synthesized from this chart.
func Testing_Synth(chart Chart) *[]interface{} {
	_init_.Initialize()

	if err := validateTesting_SynthParameters(chart); err != nil {
		panic(err)
	}
	var returns *[]interface{}

	_jsii_.StaticInvoke(
		"cdk8s.Testing",
		"synth",
		[]interface{}{chart},
		&returns,
	)

	return returns
}

