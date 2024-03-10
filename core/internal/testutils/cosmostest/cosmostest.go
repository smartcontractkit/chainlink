package cosmostest

import (
	"fmt"

	"github.com/google/uuid"
)

// RandomChainID returns a random chain id for testing. Use this instead of a constant to prevent DB collisions.
func RandomChainID() string {
	return fmt.Sprintf("Chainlinktest-%s", uuid.New())
}
