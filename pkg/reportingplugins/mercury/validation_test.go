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
		t.Run("errors when validFrom == block number", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = paos[i].ValidFromBlockNum
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "only 0/4 attributed observations have currentBlockNum > validFromBlockNum, need at least f+1 (2/4) to make a new report; this is most likely a duplicate report for the block range; consensusCurrentBlock=16634355, consensusValidFromBlock=16634355")
		})
		t.Run("errors when block number < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = -1
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "only 0/4 attributed observations have currentBlockNum > validFromBlockNum, need at least f+1 (2/4) to make a new report; this is most likely a duplicate report for the block range; consensusCurrentBlock=-1, consensusValidFromBlock=16634355")
		})
		t.Run("when validFrom > block number", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
				paos[i].ValidFromBlockNum = 2
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "only 0/4 attributed observations have currentBlockNum > validFromBlockNum, need at least f+1 (2/4) to make a new report; this is most likely a duplicate report for the block range; consensusCurrentBlock=1, consensusValidFromBlock=2")
		})
		t.Run("when validFrom < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
				paos[i].ValidFromBlockNum = -1
			}
			err = ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "validFromBlockNum must be >= 0 (got: -1)")
		})
		t.Run("returns error if it cannot come to consensus about currentBlockNum", func(t *testing.T) {
			paos := NewParsedAttributedObservations()
			for i := range paos {
				paos[i].CurrentBlockNum = 500 + int64(i)
				paos[i].ValidFromBlockNum = 499
			}
			err := ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "GetConsensusCurrentBlock failed: no block number matching hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e with at least f+1 votes")
		})
		t.Run("returns error if it cannot come to consensus about validFromBlockNum", func(t *testing.T) {
			paos := NewParsedAttributedObservations()
			for i := range paos {
				paos[i].CurrentBlockNum = 500
				paos[i].ValidFromBlockNum = 499 - int64(i)
			}
			err := ValidateBlockValues(paos, f, 0)
			assert.EqualError(t, err, "GetConsensusValidFromBlock failed: no valid from block number with at least f+1 votes")
		})
	})
}
