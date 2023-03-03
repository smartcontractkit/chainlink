//go:build no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

// Building without runtime type checking enabled, so all the below just return nil

func validateNames_ToDnsLabelParameters(scope constructs.Construct, options *NameOptions) error {
	return nil
}

func validateNames_ToLabelValueParameters(scope constructs.Construct, options *NameOptions) error {
	return nil
}

