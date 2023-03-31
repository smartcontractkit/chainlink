package cosmostest

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// RandomChainID returns a random chain id for testing. Use this instead of a constant to prevent DB collisions.
func RandomChainID() string {
	return fmt.Sprintf("Chainlinktest-%s", uuid.NewV4())
}
