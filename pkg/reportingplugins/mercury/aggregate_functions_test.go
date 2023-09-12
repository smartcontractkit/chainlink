package mercury

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testParsedAttributedObservation struct {
	Timestamp                  uint32
	BenchmarkPrice             *big.Int
	BenchmarkPriceValid        bool
	Bid                        *big.Int
	BidValid                   bool
	Ask                        *big.Int
	AskValid                   bool
	MaxFinalizedTimestamp      int64
	MaxFinalizedTimestampValid bool
	LinkFee                    *big.Int
	LinkFeeValid               bool
	NativeFee                  *big.Int
	NativeFeeValid             bool
}

func (t testParsedAttributedObservation) GetObserver() commontypes.OracleID { return 0 }
func (t testParsedAttributedObservation) GetTimestamp() uint32              { return t.Timestamp }
func (t testParsedAttributedObservation) GetBenchmarkPrice() (*big.Int, bool) {
	return t.BenchmarkPrice, t.BenchmarkPriceValid
}
func (t testParsedAttributedObservation) GetBid() (*big.Int, bool) {
	return t.Bid, t.BidValid
}
func (t testParsedAttributedObservation) GetAsk() (*big.Int, bool) {
	return t.Ask, t.AskValid
}
func (t testParsedAttributedObservation) GetMaxFinalizedTimestamp() (int64, bool) {
	return t.MaxFinalizedTimestamp, t.MaxFinalizedTimestampValid
}
func (t testParsedAttributedObservation) GetLinkFee() (*big.Int, bool) {
	return t.LinkFee, t.LinkFeeValid
}
func (t testParsedAttributedObservation) GetNativeFee() (*big.Int, bool) {
	return t.NativeFee, t.NativeFeeValid
}
func newValidParsedAttributedObservations() []testParsedAttributedObservation {
	return []testParsedAttributedObservation{
		testParsedAttributedObservation{
			Timestamp: 1689648456,

			BenchmarkPrice:      big.NewInt(123),
			BenchmarkPriceValid: true,
			Bid:                 big.NewInt(120),
			BidValid:            true,
			Ask:                 big.NewInt(130),
			AskValid:            true,

			MaxFinalizedTimestamp:      1679448456,
			MaxFinalizedTimestampValid: true,

			LinkFee:        big.NewInt(1),
			LinkFeeValid:   true,
			NativeFee:      big.NewInt(1),
			NativeFeeValid: true,
		},
		testParsedAttributedObservation{
			Timestamp: 1689648456,

			BenchmarkPrice:      big.NewInt(456),
			BenchmarkPriceValid: true,
			Bid:                 big.NewInt(450),
			BidValid:            true,
			Ask:                 big.NewInt(460),
			AskValid:            true,

			MaxFinalizedTimestamp:      1679448456,
			MaxFinalizedTimestampValid: true,

			LinkFee:        big.NewInt(2),
			LinkFeeValid:   true,
			NativeFee:      big.NewInt(2),
			NativeFeeValid: true,
		},
		testParsedAttributedObservation{
			Timestamp: 1689648789,

			BenchmarkPrice:      big.NewInt(789),
			BenchmarkPriceValid: true,
			Bid:                 big.NewInt(780),
			BidValid:            true,
			Ask:                 big.NewInt(800),
			AskValid:            true,

			MaxFinalizedTimestamp:      1679448456,
			MaxFinalizedTimestampValid: true,

			LinkFee:        big.NewInt(3),
			LinkFeeValid:   true,
			NativeFee:      big.NewInt(3),
			NativeFeeValid: true,
		},
		testParsedAttributedObservation{
			Timestamp: 1689648789,

			BenchmarkPrice:      big.NewInt(456),
			BenchmarkPriceValid: true,
			Bid:                 big.NewInt(450),
			BidValid:            true,
			Ask:                 big.NewInt(460),
			AskValid:            true,

			MaxFinalizedTimestamp:      1679513477,
			MaxFinalizedTimestampValid: true,

			LinkFee:        big.NewInt(4),
			LinkFeeValid:   true,
			NativeFee:      big.NewInt(4),
			NativeFeeValid: true,
		},
	}
}
func NewValidParsedAttributedObservations(paos ...testParsedAttributedObservation) []testParsedAttributedObservation {
	if len(paos) == 0 {
		paos = newValidParsedAttributedObservations()
	}
	return []testParsedAttributedObservation{
		paos[0],
		paos[1],
		paos[2],
		paos[3],
	}
}

func NewInvalidParsedAttributedObservations() []testParsedAttributedObservation {
	return []testParsedAttributedObservation{
		testParsedAttributedObservation{
			Timestamp: 1,

			BenchmarkPrice:      big.NewInt(123),
			BenchmarkPriceValid: false,
			Bid:                 big.NewInt(120),
			BidValid:            false,
			Ask:                 big.NewInt(130),
			AskValid:            false,

			MaxFinalizedTimestamp:      1679648456,
			MaxFinalizedTimestampValid: false,

			LinkFee:        big.NewInt(1),
			LinkFeeValid:   false,
			NativeFee:      big.NewInt(1),
			NativeFeeValid: false,
		},
		testParsedAttributedObservation{
			Timestamp: 2,

			BenchmarkPrice:      big.NewInt(456),
			BenchmarkPriceValid: false,
			Bid:                 big.NewInt(450),
			BidValid:            false,
			Ask:                 big.NewInt(460),
			AskValid:            false,

			MaxFinalizedTimestamp:      1679648456,
			MaxFinalizedTimestampValid: false,

			LinkFee:        big.NewInt(2),
			LinkFeeValid:   false,
			NativeFee:      big.NewInt(2),
			NativeFeeValid: false,
		},
		testParsedAttributedObservation{
			Timestamp: 2,

			BenchmarkPrice:      big.NewInt(789),
			BenchmarkPriceValid: false,
			Bid:                 big.NewInt(780),
			BidValid:            false,
			Ask:                 big.NewInt(800),
			AskValid:            false,

			MaxFinalizedTimestamp:      1679648456,
			MaxFinalizedTimestampValid: false,

			LinkFee:        big.NewInt(3),
			LinkFeeValid:   false,
			NativeFee:      big.NewInt(3),
			NativeFeeValid: false,
		},
		testParsedAttributedObservation{
			Timestamp: 3,

			BenchmarkPrice:      big.NewInt(456),
			BenchmarkPriceValid: true,
			Bid:                 big.NewInt(450),
			BidValid:            true,
			Ask:                 big.NewInt(460),
			AskValid:            true,

			MaxFinalizedTimestamp:      1679513477,
			MaxFinalizedTimestampValid: true,

			LinkFee:        big.NewInt(4),
			LinkFeeValid:   true,
			NativeFee:      big.NewInt(4),
			NativeFeeValid: true,
		},
	}
}

func Test_AggregateFunctions(t *testing.T) {
	f := 1
	validPaos := NewValidParsedAttributedObservations()
	invalidPaos := NewInvalidParsedAttributedObservations()

	t.Run("GetConsensusTimestamp", func(t *testing.T) {
		validMPaos := convert(validPaos)
		ts := GetConsensusTimestamp(validMPaos)

		assert.Equal(t, 1689648789, int(ts))
	})

	t.Run("GetConsensusBenchmarkPrice", func(t *testing.T) {
		t.Run("gets consensus price when prices are valid", func(t *testing.T) {
			validMPaos := convert(validPaos)
			bp, err := GetConsensusBenchmarkPrice(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, "456", bp.String())
		})

		t.Run("fails when fewer than f+1 prices are valid", func(t *testing.T) {
			invalidMPaos := convert(invalidPaos)
			_, err := GetConsensusBenchmarkPrice(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price (got: 1/4)")
		})
	})

	t.Run("GetConsensusBid", func(t *testing.T) {
		t.Run("gets consensus bid when prices are valid", func(t *testing.T) {
			validMPaos := convertBid(validPaos)
			bid, err := GetConsensusBid(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, "450", bid.String())
		})

		t.Run("fails when fewer than f+1 prices are valid", func(t *testing.T) {
			invalidMPaos := convertBid(invalidPaos)
			_, err := GetConsensusBid(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price (got: 1/4)")
		})
	})

	t.Run("GetConsensusAsk", func(t *testing.T) {
		t.Run("gets consensus ask when prices are valid", func(t *testing.T) {
			validMPaos := convertAsk(validPaos)
			bid, err := GetConsensusAsk(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, "460", bid.String())
		})

		t.Run("fails when fewer than f+1 prices are valid", func(t *testing.T) {
			invalidMPaos := convertAsk(invalidPaos)
			_, err := GetConsensusAsk(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid price (got: 1/4)")
		})
	})

	t.Run("GetConsensusMaxFinalizedTimestamp", func(t *testing.T) {
		t.Run("gets consensus on maxFinalizedTimestamp when valid", func(t *testing.T) {
			validMPaos := convertMaxFinalizedTimestamp(validPaos)
			ts, err := GetConsensusMaxFinalizedTimestamp(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, int64(1679448456), ts)
		})

		t.Run("uses highest value as tiebreaker", func(t *testing.T) {
			paos := newValidParsedAttributedObservations()
			(paos[0]).MaxFinalizedTimestamp = 1679513477
			validMPaos := convertMaxFinalizedTimestamp(NewValidParsedAttributedObservations(paos...))
			ts, err := GetConsensusMaxFinalizedTimestamp(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, int64(1679513477), ts)
		})

		t.Run("fails when fewer than f+1 maxFinalizedTimestamps are valid", func(t *testing.T) {
			invalidMPaos := convertMaxFinalizedTimestamp(invalidPaos)
			_, err := GetConsensusMaxFinalizedTimestamp(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid maxFinalizedTimestamp (got: 1/4)")
		})

		t.Run("fails when cannot come to consensus f+1 maxFinalizedTimestamps", func(t *testing.T) {
			paos := []PAOMaxFinalizedTimestamp{
				testParsedAttributedObservation{
					MaxFinalizedTimestamp:      1679648456,
					MaxFinalizedTimestampValid: true,
				},
				testParsedAttributedObservation{
					MaxFinalizedTimestamp:      1679648457,
					MaxFinalizedTimestampValid: true,
				},
				testParsedAttributedObservation{
					MaxFinalizedTimestamp:      1679648458,
					MaxFinalizedTimestampValid: true,
				},
				testParsedAttributedObservation{
					MaxFinalizedTimestamp:      1679513477,
					MaxFinalizedTimestampValid: true,
				},
			}
			_, err := GetConsensusMaxFinalizedTimestamp(paos, f)
			assert.EqualError(t, err, "no valid maxFinalizedTimestamp with at least f+1 votes (got counts: map[1679513477:1 1679648456:1 1679648457:1 1679648458:1])")
		})
	})

	t.Run("GetConsensusLinkFee", func(t *testing.T) {
		t.Run("gets consensus on linkFee when valid", func(t *testing.T) {
			validMPaos := convertLinkFee(validPaos)
			linkFee, err := GetConsensusLinkFee(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(3), linkFee)
		})
		t.Run("treats zero values as valid", func(t *testing.T) {
			paos := NewValidParsedAttributedObservations()
			for i := range paos {
				paos[i].LinkFee = big.NewInt(0)
			}
			linkFee, err := GetConsensusLinkFee(convertLinkFee(paos), f)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(0), linkFee)
		})
		t.Run("treats negative values as invalid", func(t *testing.T) {
			paos := NewValidParsedAttributedObservations()
			for i := range paos {
				paos[i].LinkFee = big.NewInt(int64(0 - i))
			}
			_, err := GetConsensusLinkFee(convertLinkFee(paos), f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid linkFee (got: 1/4)")
		})

		t.Run("fails when fewer than f+1 linkFees are valid", func(t *testing.T) {
			invalidMPaos := convertLinkFee(invalidPaos)
			_, err := GetConsensusLinkFee(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid linkFee (got: 1/4)")
		})
	})

	t.Run("GetConsensusNativeFee", func(t *testing.T) {
		t.Run("gets consensus on nativeFee when valid", func(t *testing.T) {
			validMPaos := convertNativeFee(validPaos)
			nativeFee, err := GetConsensusNativeFee(validMPaos, f)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(3), nativeFee)
		})
		t.Run("treats zero values as valid", func(t *testing.T) {
			paos := NewValidParsedAttributedObservations()
			for i := range paos {
				paos[i].NativeFee = big.NewInt(0)
			}
			nativeFee, err := GetConsensusNativeFee(convertNativeFee(paos), f)
			require.NoError(t, err)
			assert.Equal(t, big.NewInt(0), nativeFee)
		})
		t.Run("treats negative values as invalid", func(t *testing.T) {
			paos := NewValidParsedAttributedObservations()
			for i := range paos {
				paos[i].NativeFee = big.NewInt(int64(0 - i))
			}
			_, err := GetConsensusNativeFee(convertNativeFee(paos), f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid nativeFee (got: 1/4)")
		})
		t.Run("fails when fewer than f+1 nativeFees are valid", func(t *testing.T) {
			invalidMPaos := convertNativeFee(invalidPaos)
			_, err := GetConsensusNativeFee(invalidMPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid nativeFee (got: 1/4)")
		})
	})
}

// convert funcs are necessary because go is not smart enough to cast
// []interface1 to []interface2 even if interface1 is a superset of interface2
func convert(pao []testParsedAttributedObservation) (ret []PAO) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertMaxFinalizedTimestamp(pao []testParsedAttributedObservation) (ret []PAOMaxFinalizedTimestamp) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertAsk(pao []testParsedAttributedObservation) (ret []PAOAsk) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertBid(pao []testParsedAttributedObservation) (ret []PAOBid) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertLinkFee(pao []testParsedAttributedObservation) (ret []PAOLinkFee) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
func convertNativeFee(pao []testParsedAttributedObservation) (ret []PAONativeFee) {
	for _, v := range pao {
		ret = append(ret, v)
	}
	return ret
}
