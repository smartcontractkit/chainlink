package adapters

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Bounds represents the start/end param for the Random adapter
type Bounds int64

// UnmarshalJSON implements json.Unmarshaler for Random
func (r *Bounds) UnmarshalJSON(input []byte) error {
	input = utils.RemoveQuotes(input)
	bound, err := strconv.ParseInt(string(input), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse into float: %s", input)
	}
	*r = Bounds(bound)
	return nil
}

// Random holds the range given by start and end parameters
type Random struct {
	Start Bounds `json:"start"`
	End   Bounds `json:"end"`
}

// Perform returns a random integer in the inclusive range as
// specified by Start and End parameters.
func (ra *Random) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	start := int64(ra.Start)
	end := int64(ra.End)
	if start == 0 && end == 0 {
		input.SetError(fmt.Errorf("Both start and end ranges must be specified as parameters"))
		return input
	}
	if start >= end {
		input.SetError(fmt.Errorf("End must be strictly greater than start"))
		return input
	}
	diff := big.NewInt(end - start)
	val, err := rand.Int(rand.Reader, diff)
	if err != nil {
		input.SetError(err)
		return input
	}
	randNumInRange := val.Int64() + start
	input.ApplyResult(strconv.FormatInt(randNumInRange, 10))
	return input
}
