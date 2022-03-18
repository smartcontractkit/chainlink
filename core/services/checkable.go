package services

// Checkable should be implemented by any type requiring health checks.
type Checkable interface {
	// Ready should return nil if ready, or an error message otherwise.
	Ready() error
	// Healthy should return nil if healthy, or an error message otherwise.
	Healthy() error
}
