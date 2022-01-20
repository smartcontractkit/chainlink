package services

// Types requiring health checks should implement the Checkable interface.
type Checkable interface {
	// Checkables should return nil if ready, or an error message otherwise.
	Ready() error
	// Checkables should return nil if healthy, or an error message otherwise.
	Healthy() error
}
