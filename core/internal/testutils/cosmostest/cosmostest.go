package cosmostest

import (
	"fmt"
	"math/rand"
)

// RandomChainID returns a random chain id for testing. Use this instead of a constant to prevent DB collisions.
func RandomChainID() string {
	return fmt.Sprintf("Chainlinktest-%d", rand.Int31n(999999))
}
