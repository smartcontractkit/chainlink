//go:build !no_runtime_type_checking

// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s

import (
	"fmt"
)

func validateLazy_AnyParameters(producer IAnyProducer) error {
	if producer == nil {
		return fmt.Errorf("parameter producer is required, but nil was provided")
	}

	return nil
}

