package services

// Checkable should be implemented by any type requiring health checks.
// From the k8s docs:
// > ready means it’s initialized and healthy means that it can accept traffic in kubernetes
// See: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
type Checkable interface {
	// Ready should return nil if ready, or an error message otherwise.
	Ready() error
	// HealthReport returns a full health report of the callee including it's dependencies.
	// key is the dep name, value is nil if healthy, or error message otherwise.
	HealthReport() map[string]error
}
