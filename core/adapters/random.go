package adapters

import (
	"crypto/rand"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Random adapter type holds no fields
type Random struct{}

// Perform returns a random uint256 number in 0 | 2**256-1 range
func (ra *Random) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return models.NewRunOutputError(err)
	}
	ran := new(big.Int).SetBytes(b)
	return models.NewRunOutputCompleteWithResult(ran.String())
}
