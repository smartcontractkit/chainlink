// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// OwnerReference contains enough information to let you identify an owning object.
//
// An owning object must be in the same namespace as the dependent, or
// be cluster-scoped, so there is no namespace field.
type OwnerReference struct {
	// API version of the referent.
	ApiVersion *string `field:"required" json:"apiVersion" yaml:"apiVersion"`
	// Kind of the referent.
	// See: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	//
	Kind *string `field:"required" json:"kind" yaml:"kind"`
	// Name of the referent.
	// See: http://kubernetes.io/docs/user-guide/identifiers#names
	//
	Name *string `field:"required" json:"name" yaml:"name"`
	// UID of the referent.
	// See: http://kubernetes.io/docs/user-guide/identifiers#uids
	//
	Uid *string `field:"required" json:"uid" yaml:"uid"`
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed.
	//
	// Defaults to false. To set this field, a user needs "delete"
	// permission of the owner, otherwise 422 (Unprocessable Entity) will be
	// returned.
	BlockOwnerDeletion *bool `field:"optional" json:"blockOwnerDeletion" yaml:"blockOwnerDeletion"`
	// If true, this reference points to the managing controller.
	Controller *bool `field:"optional" json:"controller" yaml:"controller"`
}

