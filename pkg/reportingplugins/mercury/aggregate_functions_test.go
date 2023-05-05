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
	paos := NewParsedAttributedObservations()

	t.Run("GetConsensusTimestamp", func(t *testing.T) {
		ts := GetConsensusTimestamp(paos)

		assert.Equal(t, 1676484828, int(ts))
	})
	t.Run("GetConsensusBenchmarkPrice", func(t *testing.T) {
		bp := GetConsensusBenchmarkPrice(paos)

		assert.Equal(t, "346", bp.String())
	})
	t.Run("GetConsensusBid", func(t *testing.T) {
		bid := GetConsensusBid(paos)

		assert.Equal(t, "345", bid.String())
	})
	t.Run("GetConsensusAsk", func(t *testing.T) {
		ask := GetConsensusAsk(paos)

		assert.Equal(t, "350", ask.String())
	})

	t.Run("GetConsensusCurrentBlock", func(t *testing.T) {
		hash, num, ts, err := GetConsensusCurrentBlock(paos, f)

		require.NoError(t, err)
		assert.Equal(t, mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"), hash)
		assert.Equal(t, 16634365, int(num))
		assert.Equal(t, uint64(1682591344), ts)

		t.Run("if there are not at least f+1 in consensus about hash", func(t *testing.T) {
			_, _, _, err := GetConsensusCurrentBlock(paos, 3)
			assert.EqualError(t, err, "couldn't get consensus current block: no block hash with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := NewParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockNum = int64(i)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "coulnd't get consensus current block: no block number matching hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about timestamp", func(t *testing.T) {
			badPaos := NewParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockTimestamp = uint64(i * 100)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "coulnd't get consensus current block: no block timestamp matching block hash 0x40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e and block number 16634365 with at least f+1 votes")
		})
	})

	t.Run("GetConsensusValidFromBlock", func(t *testing.T) {
		num, err := GetConsensusValidFromBlock(paos, f)

		require.NoError(t, err)
		assert.Equal(t, 16634355, int(num))

		t.Run("if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := NewParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].ValidFromBlockNum = int64(i)
			}

			_, err := GetConsensusValidFromBlock(badPaos, f)
			assert.EqualError(t, err, "no valid from block number with at least f+1 votes")
		})
	})
}
