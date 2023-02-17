package mercury

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	f := 1
	paos := NewParsedAttributedObservations()
	min := big.NewInt(0)
	max := big.NewInt(10_000)

	badMin := big.NewInt(9_000)
	badMax := big.NewInt(10)

	t.Run("ValidateBenchmarkPrice", func(t *testing.T) {
		err := ValidateBenchmarkPrice(paos, min, max)
		assert.NoError(t, err)

		err = ValidateBenchmarkPrice(paos, min, badMax)
		assert.EqualError(t, err, "median benchmark price 346 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateBenchmarkPrice(paos, badMin, max)
		assert.EqualError(t, err, "median benchmark price 346 is outside of allowable range (Min: 9000, Max: 10000)")
	})

	t.Run("ValidateBid", func(t *testing.T) {
		err := ValidateBid(paos, min, max)
		assert.NoError(t, err)

		err = ValidateBid(paos, min, badMax)
		assert.EqualError(t, err, "median bid price 345 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateBid(paos, badMin, max)
		assert.EqualError(t, err, "median bid price 345 is outside of allowable range (Min: 9000, Max: 10000)")
	})
	t.Run("ValidateAsk", func(t *testing.T) {
		err := ValidateAsk(paos, min, max)
		assert.NoError(t, err)

		err = ValidateAsk(paos, min, badMax)
		assert.EqualError(t, err, "median ask price 350 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateAsk(paos, badMin, max)
		assert.EqualError(t, err, "median ask price 350 is outside of allowable range (Min: 9000, Max: 10000)")
	})
	t.Run("ValidateBlockValues", func(t *testing.T) {
		err := ValidateBlockValues(paos, f, 0)
		assert.NoError(t, err)

		t.Run("errors when maxFinalizedBlockNumber is equal to or larger than current block number", func(t *testing.T) {
			err := ValidateBlockValues(paos, f, 16634365)
			assert.EqualError(t, err, "maxFinalizedBlockNumber (16634365) must be less than current block number (16634365)")
		})
		t.Run("valid when validFrom == block number", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = paos[i].ValidFromBlockNum
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.NoError(t, err)
		})
		t.Run("errors when block number < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = -1
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "block number must be >= 0 (got: -1)")
		})
		t.Run("when validFrom > block number", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
				paos[i].ValidFromBlockNum = 2
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "validFromBlockNum (2) must be less than or equal to current block number (1)")
		})
		t.Run("when validFrom < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
				paos[i].ValidFromBlockNum = -1
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "validFromBlockNum must be >= 0 (got: -1)")
		})
	})
}
