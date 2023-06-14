package mercury

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
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
			assert.EqualError(t, err, "no unique block with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := NewValidParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockNum = int64(i)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "no unique block with at least f+1 votes")
		})
		t.Run("if there are not at least f+1 in consensus about timestamp", func(t *testing.T) {
			badPaos := NewValidParsedAttributedObservations()
			for i := range badPaos {
				badPaos[i].CurrentBlockTimestamp = uint64(i * 100)
			}
			_, _, _, err := GetConsensusCurrentBlock(badPaos, f)
			assert.EqualError(t, err, "no unique block with at least f+1 votes")
		})
		t.Run("in the event of an even split for block number/hash, take the higher block number", func(t *testing.T) {
			validFrom := int64(26014056)
			// values below are from a real observed case of this happening in the wild
			paos := []ParsedAttributedObservation{
				ParsedAttributedObservation{
					Timestamp:                    1686759784,
					Observer:                     commontypes.OracleID(2),
					BenchmarkPrice:               big.NewInt(90700),
					Bid:                          big.NewInt(26200),
					Ask:                          big.NewInt(17500),
					PricesValid:                  true,
					CurrentBlockNum:              26014055,
					CurrentBlockHash:             mustDecodeHex("1a2b96ef9a29614c9fc4341a5ca6690ed8ee1a2cd6b232c90ba8bea65a4b93b5"),
					CurrentBlockTimestamp:        1686759784,
					CurrentBlockValid:            true,
					MaxFinalizedBlockNumber:      0,
					MaxFinalizedBlockNumberValid: false,
				},
				ParsedAttributedObservation{
					Timestamp:                    1686759784,
					Observer:                     commontypes.OracleID(3),
					BenchmarkPrice:               big.NewInt(92000),
					Bid:                          big.NewInt(21300),
					Ask:                          big.NewInt(74700),
					PricesValid:                  true,
					CurrentBlockNum:              26014056,
					CurrentBlockHash:             mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"),
					CurrentBlockTimestamp:        1686759784,
					CurrentBlockValid:            true,
					MaxFinalizedBlockNumber:      0,
					MaxFinalizedBlockNumberValid: false,
				},
				ParsedAttributedObservation{
					Timestamp:                    1686759784,
					Observer:                     commontypes.OracleID(1),
					BenchmarkPrice:               big.NewInt(67300),
					Bid:                          big.NewInt(70100),
					Ask:                          big.NewInt(83200),
					PricesValid:                  true,
					CurrentBlockNum:              26014056,
					CurrentBlockHash:             mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"),
					CurrentBlockTimestamp:        1686759784,
					CurrentBlockValid:            true,
					MaxFinalizedBlockNumber:      0,
					MaxFinalizedBlockNumberValid: false,
				},
				ParsedAttributedObservation{
					Timestamp:                    1686759784,
					Observer:                     commontypes.OracleID(0),
					BenchmarkPrice:               big.NewInt(8600),
					Bid:                          big.NewInt(89100),
					Ask:                          big.NewInt(53300),
					PricesValid:                  true,
					CurrentBlockNum:              26014055,
					CurrentBlockHash:             mustDecodeHex("1a2b96ef9a29614c9fc4341a5ca6690ed8ee1a2cd6b232c90ba8bea65a4b93b5"),
					CurrentBlockTimestamp:        1686759784,
					CurrentBlockValid:            true,
					MaxFinalizedBlockNumber:      0,
					MaxFinalizedBlockNumberValid: false,
				},
			}
			assert.NoError(t, ValidateCurrentBlock(paos, f, validFrom))
			hash, num, _, err := GetConsensusCurrentBlock(paos, f)
			assert.NoError(t, err)
			assert.Equal(t, mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"), hash)
			assert.Equal(t, int64(26014056), num)
			assert.GreaterOrEqual(t, num, validFrom)
		})
		pao := func(num int64, hash string, ts uint64) ParsedAttributedObservation {
			return ParsedAttributedObservation{CurrentBlockNum: num, CurrentBlockHash: mustDecodeHex(hash), CurrentBlockTimestamp: ts, CurrentBlockValid: true}
		}
		t.Run("when there are multiple possible blocks meeting > f+1 hashes, takes the hash with the most block numbers in agreement", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
				pao(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
				pao(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
				pao(42, "1111111111111111111111111111111111111111111111111111111111111111", 1),
				pao(42, "1111111111111111111111111111111111111111111111111111111111111111", 1),
				pao(41, "1111111111111111111111111111111111111111111111111111111111111111", 1),
				pao(43, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
			}
			assert.NoError(t, ValidateCurrentBlock(paos, f, 41))
			hash, num, ts, err := GetConsensusCurrentBlock(paos, f)
			assert.NoError(t, err)
			assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
			assert.Equal(t, int64(42), num)
			assert.Equal(t, uint64(1), ts)
		})
		t.Run("in the event of an even split of numbers/hashes, takes the hash with the highest block number", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
			}
			assert.NoError(t, ValidateCurrentBlock(paos, f, 41))
			hash, num, ts, err := GetConsensusCurrentBlock(paos, f)
			assert.NoError(t, err)
			assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
			assert.Equal(t, int64(42), num)
			assert.Equal(t, uint64(1), ts)
		})
		t.Run("in the case where all block numbers are equal but timestamps differ, tie-breaks on latest timestamp", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
			}
			assert.NoError(t, ValidateCurrentBlock(paos, f, 41))
			hash, num, ts, err := GetConsensusCurrentBlock(paos, f)
			assert.NoError(t, err)
			assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
			assert.Equal(t, int64(42), num)
			assert.Equal(t, uint64(2), ts)
		})
		t.Run("in the case where all block numbers and timestamps are equal, tie-breaks by taking the 'lowest' hash", func(t *testing.T) {
			paos := []ParsedAttributedObservation{
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				pao(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
			}
			assert.NoError(t, ValidateCurrentBlock(paos, f, 41))
			hash, num, ts, err := GetConsensusCurrentBlock(paos, f)
			assert.NoError(t, err)
			assert.Equal(t, mustDecodeHex("2222222222222222222222222222222222222222222222222222222222222222"), hash)
			assert.Equal(t, int64(42), num)
			assert.Equal(t, uint64(1), ts)
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
