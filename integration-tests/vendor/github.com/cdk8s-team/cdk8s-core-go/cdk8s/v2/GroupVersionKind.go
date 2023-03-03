// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


type GroupVersionKind struct {
	// The object's API version (e.g. `authorization.k8s.io/v1`).
	ApiVersion *string `field:"required" json:"apiVersion" yaml:"apiVersion"`
	// The object kind.
	Kind *string `field:"required" json:"kind" yaml:"kind"`
}

