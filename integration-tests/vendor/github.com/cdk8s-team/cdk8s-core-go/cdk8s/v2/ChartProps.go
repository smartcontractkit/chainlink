// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


type ChartProps struct {
	// Labels to apply to all resources in this chart.
	Labels *map[string]*string `field:"optional" json:"labels" yaml:"labels"`
	// The default namespace for all objects defined in this chart (directly or indirectly).
	//
	// This namespace will only apply to objects that don't have a
	// `namespace` explicitly defined for them.
	Namespace *string `field:"optional" json:"namespace" yaml:"namespace"`
}

