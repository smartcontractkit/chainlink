package services

// Checkable should be implemented by any type requiring health checks.
type Checkable interface {
	// Checkables should return nil if ready, or an error message otherwise.
	Ready() error
	// Checkables should return nil if healthy, or an error message otherwise.
	Healthy() error
}
