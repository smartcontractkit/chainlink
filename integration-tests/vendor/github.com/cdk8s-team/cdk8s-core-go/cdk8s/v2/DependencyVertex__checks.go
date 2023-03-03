//go:build !no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	"fmt"
)

func (d *jsiiProxy_DependencyVertex) validateAddChildParameters(dep DependencyVertex) error {
	if dep == nil {
		return fmt.Errorf("parameter dep is required, but nil was provided")
	}

	return nil
}

