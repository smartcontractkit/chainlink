// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Options for name generation.
type NameOptions struct {
	// Delimiter to use between components.
	Delimiter *string `field:"optional" json:"delimiter" yaml:"delimiter"`
	// Extra components to include in the name.
	Extra *[]*string `field:"optional" json:"extra" yaml:"extra"`
	// Include a short hash as last part of the name.
	IncludeHash *bool `field:"optional" json:"includeHash" yaml:"includeHash"`
	// Maximum allowed length for the name.
	MaxLen *float64 `field:"optional" json:"maxLen" yaml:"maxLen"`
}

