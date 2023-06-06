package mercury

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func Test_AggregateFunctions(t *testing.T) {
	f := 1
	validPaos := NewValidParsedAttributedObservations()
	invalidPaos := NewInvalidParsedAttributedObservations()

	t.Run("GetConsensusTimestamp", func(t *testing.T) {
		ts := GetConsensusTimestamp(validPaos)

		assert.Equal(t, 1676484828, int(ts))
	})
	t.Run("GetConsensusBenchmarkPrice", func(t *testing.T) {
		t.Run("when prices valid, gets median price", func(t *testing.T) {
			bp, err := GetConsensusBenchmarkPrice(validPaos, f)
			require.NoError(t, err)
			assert.Equal(t, "346", bp.String())
		})

		t.Run("if more than f+1 are invalid, fails", func(t *testing.T) {
			_, err := GetConsensusBenchmarkPrice(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price")
		})
	})
	t.Run("GetConsensusBid", func(t *testing.T) {
		t.Run("when prices valid, gets median bid", func(t *testing.T) {
			bid, err := GetConsensusBid(validPaos, f)
			require.NoError(t, err)
			assert.Equal(t, "345", bid.String())
		})

		t.Run("if more than f+1 are invalid, fails", func(t *testing.T) {
			_, err := GetConsensusBid(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price")
		})
	})
	t.Run("GetConsensusAsk", func(t *testing.T) {
		t.Run("when prices valid, gets median bid", func(t *testing.T) {
			ask, err := GetConsensusAsk(validPaos, f)
			require.NoError(t, err)

			assert.Equal(t, "350", ask.String())
		})

		t.Run("if invalid, fails", func(t *testing.T) {
			_, err := GetConsensusAsk(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price")
		})
	})

	t.Run("GetConsensusCurrentBlock", func(t *testing.T) {
		t.Run("succeeds in the valid case", func(t *testing.T) {
			hash, num, ts, err := GetConsensusCurrentBlock(validPaos, f)

			require.NoError(t, err)
			assert.Equal(t, mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"), hash)
			assert.Equal(t, 16634365, int(num))
			assert.Equal(t, uint64(1682591344), ts)
		})

		t.Run("if invalid, fails", func(t *testing.T) {
			_, _, _, err := GetConsensusCurrentBlock(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid current block (got: 0/4)")
		})
		t.Run("if there are not at least f+1 in consensus about hash", func(t *testing.T) {
			_, _, _, err := GetConsensusCurrentBlock(validPaos, 3)
			assert.EqualError(t, err, "couldn't get consensus current block: no block hash with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := NewValidParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockNum = int64(i)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "couldn't get consensus current block: no block number matching hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about timestamp", func(t *testing.T) {
			badPaos := NewValidParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockTimestamp = uint64(i * 100)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "couldn't get consensus current block: no block timestamp matching block hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e and block number 16634365 with at least f+1 votes")
		})
	})

	t.Run("GetConsensusMaxFinalizedBlockNum", func(t *testing.T) {
		t.Run("in the valid case", func(t *testing.T) {
			num, err := GetConsensusMaxFinalizedBlockNum(validPaos, f)

			require.NoError(t, err)
			assert.Equal(t, 16634355, int(num))
		})

		t.Run("errors if there are not at least f+1 valid", func(t *testing.T) {
			_, err := GetConsensusMaxFinalizedBlockNum(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid maxFinalizedBlockNumber (got: 0/4)")
		})

		t.Run("errors if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := NewValidParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].MaxFinalizedBlockNumber = int64(i)
			}

			_, err := GetConsensusMaxFinalizedBlockNum(badPaos, f)
			assert.EqualError(t, err, "no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: map[0:1 1:1 2:1 3:1])")
		})
	})
}
