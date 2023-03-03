// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Options for defining API objects.
type ApiObjectProps struct {
	// API version.
	ApiVersion *string `field:"required" json:"apiVersion" yaml:"apiVersion"`
	// Resource kind.
	Kind *string `field:"required" json:"kind" yaml:"kind"`
	// Object metadata.
	//
	// If `name` is not specified, an app-unique name will be allocated by the
	// framework based on the path of the construct within thes construct tree.
	Metadata *ApiObjectMetadata `field:"optional" json:"metadata" yaml:"metadata"`
}

