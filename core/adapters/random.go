package adapters

import (
	"crypto/rand"
	"math/big"

	"chainlink/core/store"
	"chainlink/core/store/models"
)

// Random adapter type holds no fields
type Random struct{}

// Perform returns a random uint256 number in 0 | 2**256-1 range
func (ra *Random) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return models.RunResultError(err)
	}
	ran := new(big.Int).SetBytes(b)
	return models.RunResultComplete(ran.String())
}
