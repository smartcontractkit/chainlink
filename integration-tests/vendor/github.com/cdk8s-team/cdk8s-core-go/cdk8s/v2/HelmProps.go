// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Options for `Helm`.
type HelmProps struct {
	// The chart name to use. It can be a chart from a helm repository or a local directory.
	//
	// This name is passed to `helm template` and has all the relevant semantics.
	//
	// Example:
	//   "bitnami/redis"
	//
	Chart *string `field:"required" json:"chart" yaml:"chart"`
	// The local helm executable to use in order to create the manifest the chart.
	HelmExecutable *string `field:"optional" json:"helmExecutable" yaml:"helmExecutable"`
	// Additional flags to add to the `helm` execution.
	HelmFlags *[]*string `field:"optional" json:"helmFlags" yaml:"helmFlags"`
	// Scope all resources in to a given namespace.
	Namespace *string `field:"optional" json:"namespace" yaml:"namespace"`
	// The release name.
	// See: https://helm.sh/docs/intro/using_helm/#three-big-concepts
	//
	ReleaseName *string `field:"optional" json:"releaseName" yaml:"releaseName"`
	// Chart repository url where to locate the requested chart.
	Repo *string `field:"optional" json:"repo" yaml:"repo"`
	// Values to pass to the chart.
	Values *map[string]interface{} `field:"optional" json:"values" yaml:"values"`
	// Version constraint for the chart version to use.
	//
	// This constraint can be a specific tag (e.g. 1.1.1)
	// or it may reference a valid range (e.g. ^2.0.0).
	// If this is not specified, the latest version is used
	//
	// This name is passed to `helm template --version` and has all the relevant semantics.
	//
	// Example:
	//   "^2.0.0"
	//
	Version *string `field:"optional" json:"version" yaml:"version"`
}

