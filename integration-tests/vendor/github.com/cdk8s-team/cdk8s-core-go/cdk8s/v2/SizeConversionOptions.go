// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Options for how to convert time to a different unit.
type SizeConversionOptions struct {
	// How conversions should behave when it encounters a non-integer result.
	Rounding SizeRoundingBehavior `field:"optional" json:"rounding" yaml:"rounding"`
}

