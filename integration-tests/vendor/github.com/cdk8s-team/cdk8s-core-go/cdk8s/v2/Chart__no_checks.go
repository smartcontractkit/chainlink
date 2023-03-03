//go:build no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

// Building without runtime type checking enabled, so all the below just return nil

func (c *jsiiProxy_Chart) validateGenerateObjectNameParameters(apiObject ApiObject) error {
	return nil
}

func validateChart_IsChartParameters(x interface{}) error {
	return nil
}

func validateChart_IsConstructParameters(x interface{}) error {
	return nil
}

func validateChart_OfParameters(c constructs.IConstruct) error {
	return nil
}

func validateNewChartParameters(scope constructs.Construct, id *string, props *ChartProps) error {
	return nil
}

