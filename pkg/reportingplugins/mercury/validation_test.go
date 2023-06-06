package mercury

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	f := 1
	paos := NewValidParsedAttributedObservations()
	min := big.NewInt(0)
	max := big.NewInt(10_000)

	badMin := big.NewInt(9_000)
	badMax := big.NewInt(10)

	t.Run("ValidateBenchmarkPrice", func(t *testing.T) {
		err := ValidateBenchmarkPrice(paos, f, min, max)
		assert.NoError(t, err)

		err = ValidateBenchmarkPrice(paos, f, min, badMax)
		assert.EqualError(t, err, "median benchmark price 346 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateBenchmarkPrice(paos, f, badMin, max)
		assert.EqualError(t, err, "median benchmark price 346 is outside of allowable range (Min: 9000, Max: 10000)")
	})

	t.Run("ValidateBid", func(t *testing.T) {
		err := ValidateBid(paos, f, min, max)
		assert.NoError(t, err)

		err = ValidateBid(paos, f, min, badMax)
		assert.EqualError(t, err, "median bid price 345 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateBid(paos, f, badMin, max)
		assert.EqualError(t, err, "median bid price 345 is outside of allowable range (Min: 9000, Max: 10000)")
	})
	t.Run("ValidateAsk", func(t *testing.T) {
		err := ValidateAsk(paos, f, min, max)
		assert.NoError(t, err)

		err = ValidateAsk(paos, f, min, badMax)
		assert.EqualError(t, err, "median ask price 350 is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateAsk(paos, f, badMin, max)
		assert.EqualError(t, err, "median ask price 350 is outside of allowable range (Min: 9000, Max: 10000)")
	})
	t.Run("ValidateCurrentBlock", func(t *testing.T) {
		t.Run("succeeds when validFromBlockNum < current block num and currentBlockNum has consensus", func(t *testing.T) {
			err := ValidateCurrentBlock(paos, f, 16634364)
			assert.NoError(t, err)
		})
		t.Run("succeeds when validFromBlockNum is equal to current block number", func(t *testing.T) {
			err := ValidateCurrentBlock(paos, f, 16634365)
			assert.NoError(t, err)
		})

		t.Run("errors when block number < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = -1
			}
			err := ValidateCurrentBlock(paos, f, 2)
			assert.EqualError(t, err, "only 0/4 attributed observations have currentBlockNum >= validFromBlockNum, need at least f+1 (2/4) to make a new report; consensusCurrentBlock=-1, validFromBlockNum=2")
		})
		t.Run("errors when validFrom > block number", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
			}
			err := ValidateCurrentBlock(paos, f, 16634366)
			assert.EqualError(t, err, "only 0/4 attributed observations have currentBlockNum >= validFromBlockNum, need at least f+1 (2/4) to make a new report; consensusCurrentBlock=1, validFromBlockNum=16634366")
		})
		t.Run("errors when validFrom < 0", func(t *testing.T) {
			for i := range paos {
				paos[i].CurrentBlockNum = 1
			}
			err := ValidateCurrentBlock(paos, f, -1)
			assert.EqualError(t, err, "validFromBlockNum must be >= 0 (got: -1)")
		})
		t.Run("returns error if it cannot come to consensus about currentBlockNum", func(t *testing.T) {
			paos := NewValidParsedAttributedObservations()
			for i := range paos {
				paos[i].CurrentBlockNum = 500 + int64(i)
			}
			err := ValidateCurrentBlock(paos, f, 0)
			assert.EqualError(t, err, "GetConsensusCurrentBlock failed: couldn't get consensus current block: no block number matching hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e with at least f+1 votes")
		})
	})
}
