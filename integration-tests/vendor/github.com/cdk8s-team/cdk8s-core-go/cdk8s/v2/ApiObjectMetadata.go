// This is the core library of Cloud Development Kit (CDK) for Kubernetes (cdk8s). cdk8s apps synthesize into standard Kubernetes manifests which can be applied to any Kubernetes cluster.
package cdk8s


// Metadata associated with this object.
type ApiObjectMetadata struct {
	// Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata.
	//
	// They are not queryable and should be
	// preserved when modifying objects.
	// See: http://kubernetes.io/docs/user-guide/annotations
	//
	Annotations *map[string]*string `field:"optional" json:"annotations" yaml:"annotations"`
	// Namespaced keys that tell Kubernetes to wait until specific conditions are met before it fully deletes resources marked for deletion.
	//
	// Must be empty before the object is deleted from the registry. Each entry is
	// an identifier for the responsible component that will remove the entry from
	// the list. If the deletionTimestamp of the object is non-nil, entries in
	// this list can only be removed. Finalizers may be processed and removed in
	// any order.  Order is NOT enforced because it introduces significant risk of
	// stuck finalizers. finalizers is a shared field, any actor with permission
	// can reorder it. If the finalizer list is processed in order, then this can
	// lead to a situation in which the component responsible for the first
	// finalizer in the list is waiting for a signal (field value, external
	// system, or other) produced by a component responsible for a finalizer later
	// in the list, resulting in a deadlock. Without enforced ordering finalizers
	// are free to order amongst themselves and are not vulnerable to ordering
	// changes in the list.
	// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/
	//
	Finalizers *[]*string `field:"optional" json:"finalizers" yaml:"finalizers"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	//
	// May match selectors of replication controllers and services.
	// See: http://kubernetes.io/docs/user-guide/labels
	//
	Labels *map[string]*string `field:"optional" json:"labels" yaml:"labels"`
	// The unique, namespace-global, name of this object inside the Kubernetes cluster.
	//
	// Normally, you shouldn't specify names for objects and let the CDK generate
	// a name for you that is application-unique. The names CDK generates are
	// composed from the construct path components, separated by dots and a suffix
	// that is based on a hash of the entire path, to ensure uniqueness.
	//
	// You can supply custom name allocation logic by overriding the
	// `chart.generateObjectName` method.
	//
	// If you use an explicit name here, bear in mind that this reduces the
	// composability of your construct because it won't be possible to include
	// more than one instance in any app. Therefore it is highly recommended to
	// leave this unspecified.
	Name *string `field:"optional" json:"name" yaml:"name"`
	// Namespace defines the space within each name must be unique.
	//
	// An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty. Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces
	Namespace *string `field:"optional" json:"namespace" yaml:"namespace"`
	// List of objects depended by this object.
	//
	// If ALL objects in the list have
	// been deleted, this object will be garbage collected. If this object is
	// managed by a controller, then an entry in this list will point to this
	// controller, with the controller field set to true. There cannot be more
	// than one managing controller.
	//
	// Kubernetes sets the value of this field automatically for objects that are
	// dependents of other objects like ReplicaSets, DaemonSets, Deployments, Jobs
	// and CronJobs, and ReplicationControllers. You can also configure these
	// relationships manually by changing the value of this field. However, you
	// usually don't need to and can allow Kubernetes to automatically manage the
	// relationships.
	// See: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	//
	OwnerReferences *[]*OwnerReference `field:"optional" json:"ownerReferences" yaml:"ownerReferences"`
}

