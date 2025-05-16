package verification

import "context"

// Verifiable describes entities that can be verified.
type Verifiable interface {
	// Verify performs the verification process.
	Verify(ctx context.Context) error
	// IsVerified checks if the entity is already verified.
	IsVerified(ctx context.Context) (bool, error)
	// String returns a string representation of the Verifiable entity.
	String() string
}
