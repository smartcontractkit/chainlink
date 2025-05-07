package verification

import "context"

// Verifiable describes entities that can be verified.
type Verifiable interface {
	// Verify performs the verification process, returning true if successful.
	Verify(ctx context.Context) (bool, error)
	// String returns a string representation of the Verifiable entity.
	String() string
}
