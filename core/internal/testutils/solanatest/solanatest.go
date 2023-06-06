package solanatest

import (
	"github.com/google/uuid"
)

// RandomChainID returns a random uuid id for testing. Use this instead of a constant to prevent DB collisions.
func RandomChainID() string {
	return uuid.New().String()
}
