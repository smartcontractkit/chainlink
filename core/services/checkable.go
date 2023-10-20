package services

import (
	"errors"
	"testing"
)

// Checkable should be implemented by any type requiring health checks.
// From the k8s docs:
// > ready means it’s initialized and healthy means that it can accept traffic in kubernetes
// See: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
type Checkable interface {
	// Ready should return nil if ready, or an error message otherwise.
	Ready() error
	// HealthReport returns a full health report of the callee including it's dependencies.
	// key is the dep name, value is nil if healthy, or error message otherwise.
	// See CopyHealth.
	HealthReport() map[string]error
	// Name returns the fully qualified name of the component. Usually the logger name.
	Name() string
}

// CopyHealth copies health statuses from src to dest.
// If duplicate names are encountered, the errors are joined, unless testing in which case a panic is thrown.
func CopyHealth(dest, src map[string]error) {
	for name, err := range src {
		errOrig, ok := dest[name]
		if ok {
			if testing.Testing() {
				panic("service names must be unique: duplicate name: " + name)
			}
			if errOrig != nil {
				dest[name] = errors.Join(errOrig, err)
				continue
			}
		}
		dest[name] = err
	}
}
