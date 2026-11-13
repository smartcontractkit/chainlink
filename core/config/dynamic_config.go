package config

// DynamicConfig is the node's side of the Dynamic Config contract: a signed
// config file, rendered by the dynamic-config-sidecar, that lists the jobs the
// node must run. The node reconciles its jobs against that file on POST /reload
// instead of receiving them over the Job Distributor websocket.
type DynamicConfig interface {
	// Enabled turns the reconciler and the /reload endpoint on.
	Enabled() bool
	// Path is the rendered config file (DYNCONFIG_CONFIG_PATH on the sidecar).
	Path() string
	// ManagedJobPrefix scopes reconciliation: only jobs whose name starts with
	// it are created or deleted, so the node keeps every other job.
	ManagedJobPrefix() string
}
