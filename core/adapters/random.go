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
func (ra *Random) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	var uint256Max big.Int
	one := big.NewInt(1)
	base := big.NewInt(2)
	exp := big.NewInt(256)
	uint256Max.Sub(base.Exp(base, exp, nil), one)
	randVal, err := rand.Int(rand.Reader, &uint256Max)
	if err != nil {
		input.SetError(err)
		return input
	}
	input.ApplyResult(randVal.String())
	return input
}
