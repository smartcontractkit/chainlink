// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


type IncludeProps struct {
	// Local file path or URL which includes a Kubernetes YAML manifest.
	//
	// Example:
	//   mymanifest.yaml
	//
	Url *string `field:"required" json:"url" yaml:"url"`
}

